package forge

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type snapshotResponse struct {
	Status  int                 `json:"status"`
	Header  map[string][]string `json:"header"`
	BodyB64 string              `json:"body"`
}

type snapshotEntry struct {
	Key      string           `json:"key"`
	Method   string           `json:"method"`
	URL      string           `json:"url"`
	Response snapshotResponse `json:"response"`
}

// snapshotMode returns the mode from env: "" (off), "record", or "replay".
func snapshotMode() string {
	return os.Getenv("STRIPE_SNAPSHOT_MODE")
}

func snapshotDir() string {
	if dir := os.Getenv("STRIPE_SNAPSHOT_DIR"); dir != "" {
		return dir
	}
	return ""
}

func snapshotKey(req *http.Request) string {
	hash := sha256.New()
	hash.Write([]byte(req.Method))
	hash.Write([]byte("\n"))
	hash.Write([]byte(req.URL.String()))
	if req.Body != nil {
		buf, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(buf))
		hash.Write(buf)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func snapshotPath(key string) string {
	return filepath.Join(snapshotDir(), key+".json")
}

// snapshotIntercept returns a synthesized response if mode=replay and snapshot exists.
func snapshotIntercept(req *http.Request) (*http.Response, bool, error) {
	dir := snapshotDir()
	if dir == "" || snapshotMode() != "replay" {
		return nil, false, nil
	}
	key := snapshotKey(req)
	content, err := os.ReadFile(snapshotPath(key))
	if err != nil {
		return nil, false, nil
	}
	var entry snapshotEntry
	if err := json.Unmarshal(content, &entry); err != nil {
		return nil, false, err
	}
	body, err := base64.StdEncoding.DecodeString(entry.Response.BodyB64)
	if err != nil {
		return nil, false, err
	}
	resp := &http.Response{
		Status:     http.StatusText(entry.Response.Status),
		StatusCode: entry.Response.Status,
		Header:     entry.Response.Header,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Request:    req,
	}
	return resp, true, nil
}

// snapshotPersist persists the response body to disk if mode=record. It also
// rewinds resp.Body so the caller can still read it.
func snapshotPersist(req *http.Request, resp *http.Response) {
	dir := snapshotDir()
	if dir == "" || snapshotMode() != "record" || resp == nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	entry := snapshotEntry{
		Key:    snapshotKey(req),
		Method: req.Method,
		URL:    req.URL.String(),
		Response: snapshotResponse{
			Status:  resp.StatusCode,
			Header:  resp.Header,
			BodyB64: base64.StdEncoding.EncodeToString(body),
		},
	}
	content, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	_ = os.WriteFile(snapshotPath(entry.Key), content, 0o644)
}
