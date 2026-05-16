package forge

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	ExpiresIn    int       `json:"expires_in,omitempty"`
	ObtainedAt   time.Time `json:"obtained_at"`
	TokenURL     string    `json:"token_url,omitempty"`
	ClientID     string    `json:"client_id,omitempty"`
}

func (t OAuthToken) ExpiresAt() time.Time {
	if t.ExpiresIn == 0 || t.ObtainedAt.IsZero() {
		return time.Time{}
	}
	return t.ObtainedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
}

func (t OAuthToken) NeedsRefresh() bool {
	if t.RefreshToken == "" || t.TokenURL == "" || t.ClientID == "" {
		return false
	}
	expiresAt := t.ExpiresAt()
	if expiresAt.IsZero() {
		return false
	}
	return time.Until(expiresAt) < 60*time.Second
}

type PendingPKCE struct {
	CodeVerifier string    `json:"code_verifier"`
	ClientID     string    `json:"client_id"`
	TokenURL     string    `json:"token_url"`
	RedirectURI  string    `json:"redirect_uri"`
	Scope        string    `json:"scope,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

func TokenFile() string {
	manifest := LoadManifest()
	if override := os.Getenv(manifest.EnvPrefix + "_TOKEN_FILE"); override != "" {
		return override
	}
	return filepath.Join(".gutenberg", manifest.Slug+"-token.json")
}

// VaultFile is the shared token vault used by every Gutenberg-generated tool.
// Override with GUTENBERG_VAULT_FILE. Default: ~/.gutenberg/vault.json (or .gutenberg/vault.json if home is unavailable).
func VaultFile() string {
	if override := os.Getenv("GUTENBERG_VAULT_FILE"); override != "" {
		return override
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".gutenberg", "vault.json")
	}
	return filepath.Join(home, ".gutenberg", "vault.json")
}

// vaultEnvelope is the on-disk format of the vault. When Encrypted is true,
// Payload contains base64(nonce || ciphertext) sealed with AES-GCM using
// the 32-byte key from GUTENBERG_VAULT_KEY (hex-encoded).
type vaultEnvelope struct {
	SchemaVersion string `json:"schemaVersion"`
	Encrypted     bool   `json:"encrypted"`
	Payload       string `json:"payload"`
}

type vaultData struct {
	Tokens map[string]OAuthToken `json:"tokens"`
}

func vaultKey() ([]byte, error) {
	hex := os.Getenv("GUTENBERG_VAULT_KEY")
	if hex == "" {
		return nil, nil
	}
	key, err := decodeHex(hex)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("GUTENBERG_VAULT_KEY must decode to 32 bytes (use 64 hex chars)")
	}
	return key, nil
}

func decodeHex(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if len(value)%2 != 0 {
		return nil, errors.New("hex value must have even length")
	}
	out := make([]byte, len(value)/2)
	for i := 0; i < len(out); i++ {
		hi, err := hexNibble(value[i*2])
		if err != nil {
			return nil, err
		}
		lo, err := hexNibble(value[i*2+1])
		if err != nil {
			return nil, err
		}
		out[i] = hi<<4 | lo
	}
	return out, nil
}

func hexNibble(c byte) (byte, error) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', nil
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, nil
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, nil
	}
	return 0, fmt.Errorf("invalid hex char: %c", c)
}

func loadVault() (vaultData, error) {
	data := vaultData{Tokens: map[string]OAuthToken{}}
	content, err := os.ReadFile(VaultFile())
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return data, err
	}
	var envelope vaultEnvelope
	if err := json.Unmarshal(content, &envelope); err != nil {
		return data, err
	}
	if !envelope.Encrypted {
		if err := json.Unmarshal([]byte(envelope.Payload), &data); err != nil && envelope.Payload != "" {
			return data, err
		}
		if data.Tokens == nil {
			data.Tokens = map[string]OAuthToken{}
		}
		return data, nil
	}
	key, err := vaultKey()
	if err != nil || key == nil {
		return data, errors.New("vault is encrypted but GUTENBERG_VAULT_KEY is missing or invalid")
	}
	raw, err := base64.StdEncoding.DecodeString(envelope.Payload)
	if err != nil {
		return data, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return data, err
	}
	if len(raw) < gcm.NonceSize() {
		return data, errors.New("vault ciphertext too short")
	}
	nonce, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return data, err
	}
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return data, err
	}
	if data.Tokens == nil {
		data.Tokens = map[string]OAuthToken{}
	}
	return data, nil
}

func saveVault(data vaultData) error {
	if data.Tokens == nil {
		data.Tokens = map[string]OAuthToken{}
	}
	plaintext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	envelope := vaultEnvelope{SchemaVersion: "gutenberg.vault.v1"}
	key, err := vaultKey()
	if err != nil {
		return err
	}
	if key == nil {
		envelope.Encrypted = false
		envelope.Payload = string(plaintext)
	} else {
		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			return err
		}
		ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
		combined := append(nonce, ciphertext...)
		envelope.Encrypted = true
		envelope.Payload = base64.StdEncoding.EncodeToString(combined)
	}
	content, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(VaultFile()), 0o755); err != nil {
		return err
	}
	return os.WriteFile(VaultFile(), append(content, '\n'), 0o600)
}

func LoadStoredToken() (OAuthToken, error) {
	slug := LoadManifest().Slug
	vault, err := loadVault()
	if err == nil {
		if token, ok := vault.Tokens[slug]; ok && token.AccessToken != "" {
			return token, nil
		}
	}
	// Legacy fallback: per-tool token file.
	content, err := os.ReadFile(TokenFile())
	if err != nil {
		return OAuthToken{}, err
	}
	var token OAuthToken
	if err := json.Unmarshal(content, &token); err != nil {
		return OAuthToken{}, err
	}
	return token, nil
}

func SaveStoredToken(token OAuthToken) error {
	if token.ObtainedAt.IsZero() {
		token.ObtainedAt = time.Now().UTC()
	}
	slug := LoadManifest().Slug
	vault, err := loadVault()
	if err != nil {
		return err
	}
	vault.Tokens[slug] = token
	return saveVault(vault)
}

func Logout() error {
	slug := LoadManifest().Slug
	vault, err := loadVault()
	if err == nil {
		delete(vault.Tokens, slug)
		if saveErr := saveVault(vault); saveErr != nil {
			return saveErr
		}
	}
	if err := os.Remove(TokenFile()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func OAuthConfig() map[string]any {
	manifest := LoadManifest()
	return map[string]any{"tokenFile": TokenFile(), "auth": manifest.Auth}
}

func OAuthStatus() map[string]any {
	token, err := LoadStoredToken()
	if err != nil {
		return map[string]any{"authenticated": false, "tokenFile": TokenFile()}
	}
	return map[string]any{
		"authenticated": true,
		"tokenFile": TokenFile(),
		"tokenType": token.TokenType,
		"scope": token.Scope,
		"obtainedAt": token.ObtainedAt,
		"expiresIn": token.ExpiresIn,
	}
}

func ClientCredentials(ctx context.Context, tokenURL, clientID, clientSecret, scope string) (OAuthToken, error) {
	if tokenURL == "" || clientID == "" || clientSecret == "" {
		return OAuthToken{}, errors.New("token-url, client-id and client-secret are required")
	}
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", clientID)
	values.Set("client_secret", clientSecret)
	if scope != "" {
		values.Set("scope", scope)
	}
	token, err := postTokenForm(ctx, tokenURL, values)
	if err != nil {
		return token, err
	}
	token.TokenURL = tokenURL
	token.ClientID = clientID
	return token, nil
}

// GenerateCodeVerifier returns a cryptographically random PKCE code_verifier (RFC 7636).
func GenerateCodeVerifier() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// CodeChallengeS256 returns the PKCE S256 code_challenge for a verifier.
func CodeChallengeS256(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func pendingPKCEFile() string {
	return TokenFile() + ".pkce"
}

func savePendingPKCE(p PendingPKCE) error {
	if err := os.MkdirAll(filepath.Dir(pendingPKCEFile()), 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pendingPKCEFile(), content, 0o600)
}

func loadPendingPKCE() (PendingPKCE, error) {
	content, err := os.ReadFile(pendingPKCEFile())
	if err != nil {
		return PendingPKCE{}, err
	}
	var p PendingPKCE
	if err := json.Unmarshal(content, &p); err != nil {
		return PendingPKCE{}, err
	}
	return p, nil
}

// PKCEStart generates a verifier and returns the authorization URL the user should visit.
func PKCEStart(authURL, tokenURL, clientID, redirectURI, scope string) (string, string, error) {
	if authURL == "" || tokenURL == "" || clientID == "" || redirectURI == "" {
		return "", "", errors.New("auth-url, token-url, client-id and redirect-uri are required")
	}
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		return "", "", err
	}
	parsed, err := url.Parse(authURL)
	if err != nil {
		return "", "", err
	}
	query := parsed.Query()
	query.Set("response_type", "code")
	query.Set("client_id", clientID)
	query.Set("redirect_uri", redirectURI)
	query.Set("code_challenge", CodeChallengeS256(verifier))
	query.Set("code_challenge_method", "S256")
	if scope != "" {
		query.Set("scope", scope)
	}
	parsed.RawQuery = query.Encode()
	if err := savePendingPKCE(PendingPKCE{
		CodeVerifier: verifier,
		ClientID:     clientID,
		TokenURL:     tokenURL,
		RedirectURI:  redirectURI,
		Scope:        scope,
		CreatedAt:    time.Now().UTC(),
	}); err != nil {
		return verifier, "", err
	}
	return verifier, parsed.String(), nil
}

// PKCEFinish exchanges the authorization code for a token using the saved verifier.
func PKCEFinish(ctx context.Context, code string) (OAuthToken, error) {
	pending, err := loadPendingPKCE()
	if err != nil {
		return OAuthToken{}, fmt.Errorf("no pending PKCE session: %w", err)
	}
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", pending.RedirectURI)
	values.Set("client_id", pending.ClientID)
	values.Set("code_verifier", pending.CodeVerifier)
	token, err := postTokenForm(ctx, pending.TokenURL, values)
	if err != nil {
		return token, err
	}
	token.TokenURL = pending.TokenURL
	token.ClientID = pending.ClientID
	_ = os.Remove(pendingPKCEFile())
	return token, nil
}

// RefreshAccessToken exchanges a refresh_token for a new access_token.
func RefreshAccessToken(ctx context.Context, tokenURL, clientID, refreshToken, scope string) (OAuthToken, error) {
	if tokenURL == "" || clientID == "" || refreshToken == "" {
		return OAuthToken{}, errors.New("token-url, client-id and refresh-token are required")
	}
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", refreshToken)
	values.Set("client_id", clientID)
	if scope != "" {
		values.Set("scope", scope)
	}
	token, err := postTokenForm(ctx, tokenURL, values)
	if err != nil {
		return token, err
	}
	token.TokenURL = tokenURL
	token.ClientID = clientID
	if token.RefreshToken == "" {
		token.RefreshToken = refreshToken
	}
	return token, nil
}

// loadTokenWithMaybeRefresh loads the token; if it's near expiry and refresh
// metadata is set, it refreshes (controlled by OPEN_METEO_AUTO_REFRESH or default true).
func loadTokenWithMaybeRefresh() (OAuthToken, error) {
	token, err := LoadStoredToken()
	if err != nil {
		return token, err
	}
	disable := os.Getenv("OPEN_METEO_AUTO_REFRESH") == "0"
	if disable || !token.NeedsRefresh() {
		return token, nil
	}
	refreshed, refreshErr := MaybeRefreshStored(context.Background())
	if refreshErr != nil {
		return token, refreshErr
	}
	return refreshed, nil
}

// MaybeRefreshStored loads the stored token and refreshes it if it's near expiry.
func MaybeRefreshStored(ctx context.Context) (OAuthToken, error) {
	token, err := LoadStoredToken()
	if err != nil {
		return token, err
	}
	if !token.NeedsRefresh() {
		return token, nil
	}
	refreshed, err := RefreshAccessToken(ctx, token.TokenURL, token.ClientID, token.RefreshToken, token.Scope)
	if err != nil {
		return token, err
	}
	if err := SaveStoredToken(refreshed); err != nil {
		return refreshed, err
	}
	return refreshed, nil
}

func DeviceCode(ctx context.Context, deviceURL, tokenURL, clientID, scope string) (DeviceCodeResponse, OAuthToken, error) {
	if deviceURL == "" || tokenURL == "" || clientID == "" {
		return DeviceCodeResponse{}, OAuthToken{}, errors.New("device-url, token-url and client-id are required")
	}
	values := url.Values{}
	values.Set("client_id", clientID)
	if scope != "" {
		values.Set("scope", scope)
	}
	request, err := http.NewRequestWithContext(ctx, "POST", deviceURL, strings.NewReader(values.Encode()))
	if err != nil {
		return DeviceCodeResponse{}, OAuthToken{}, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return DeviceCodeResponse{}, OAuthToken{}, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return DeviceCodeResponse{}, OAuthToken{}, fmt.Errorf("device authorization failed: %s", response.Status)
	}
	var device DeviceCodeResponse
	if err := json.NewDecoder(response.Body).Decode(&device); err != nil {
		return DeviceCodeResponse{}, OAuthToken{}, err
	}
	interval := time.Duration(device.Interval) * time.Second
	if interval == 0 {
		interval = 5 * time.Second
	}
	deadline := time.Now().Add(time.Duration(device.ExpiresIn) * time.Second)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return device, OAuthToken{}, ctx.Err()
		case <-time.After(interval):
		}
		poll := url.Values{}
		poll.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		poll.Set("device_code", device.DeviceCode)
		poll.Set("client_id", clientID)
		token, err := postTokenForm(ctx, tokenURL, poll)
		if err == nil {
			return device, token, nil
		}
		if !strings.Contains(err.Error(), "authorization_pending") && !strings.Contains(err.Error(), "slow_down") {
			return device, OAuthToken{}, err
		}
		if strings.Contains(err.Error(), "slow_down") {
			interval += 5 * time.Second
		}
	}
	return device, OAuthToken{}, errors.New("device code expired")
}

func postTokenForm(ctx context.Context, tokenURL string, values url.Values) (OAuthToken, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return OAuthToken{}, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return OAuthToken{}, err
	}
	defer response.Body.Close()
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return OAuthToken{}, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if value, ok := payload["error"].(string); ok {
			return OAuthToken{}, fmt.Errorf("oauth token error: %s", value)
		}
		return OAuthToken{}, fmt.Errorf("oauth token error: %s", response.Status)
	}
	token := OAuthToken{ObtainedAt: time.Now().UTC()}
	token.AccessToken, _ = payload["access_token"].(string)
	token.RefreshToken, _ = payload["refresh_token"].(string)
	token.TokenType, _ = payload["token_type"].(string)
	token.Scope, _ = payload["scope"].(string)
	if expires, ok := payload["expires_in"].(float64); ok {
		token.ExpiresIn = int(expires)
	}
	if token.AccessToken == "" {
		return OAuthToken{}, errors.New("oauth token response did not include access_token")
	}
	return token, nil
}
