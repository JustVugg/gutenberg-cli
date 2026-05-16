import path from "node:path";
import { assertCanWriteDirectory, writeJson, writeText } from "./fs.js";
import { envPrefix, slugify } from "./sanitize.js";

export function generateGoProject(blueprint, outDir, options = {}) {
  const slug = slugify(options.name || blueprint.slug);
  const displayName = options.displayName || blueprint.name;
  const manifest = {
    ...blueprint,
    name: displayName,
    slug,
    envPrefix: envPrefix(slug),
    packageName: `gutenberg.local/${slug}`,
    generatedBy: "gutenberg",
    generatedByVersion: "0.3.0",
    language: "go"
  };

  assertCanWriteDirectory(outDir, Boolean(options.force));

  writeJson(path.join(outDir, "gutenberg.manifest.json"), manifest);
  writeJson(path.join(outDir, "blackforge.manifest.json"), manifest);
  writeText(path.join(outDir, "go.mod"), goMod(manifest));
  writeText(path.join(outDir, ".env.example"), `${manifest.envPrefix}_API_KEY=\n${manifest.envPrefix}_BASE_URL=${manifest.baseUrls[0] || "https://api.example.com"}\n${manifest.envPrefix}_TOKEN_FILE=\n${manifest.envPrefix}_SQLITE_FILE=\n`);
  writeText(path.join(outDir, "README.md"), readme(manifest));
  writeText(path.join(outDir, "docs", "COOKBOOK.md"), cookbook(manifest));
  writeText(path.join(outDir, "cmd", manifest.slug, "main.go"), cliMain(manifest));
  writeText(path.join(outDir, "internal", "forge", "manifest.go"), manifestGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "client.go"), clientGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "store.go"), storeGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "auth.go"), authGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "mcp.go"), mcpGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "resilience.go"), resilienceGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "snapshot.go"), snapshotGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "telemetry.go"), telemetryGo(manifest));
  writeText(path.join(outDir, "internal", "forge", "forge_test.go"), testGo(manifest));

  return {
    outDir,
    manifest,
    entrypoint: path.join(outDir, "cmd", manifest.slug)
  };
}

function goMod(manifest) {
  return `module ${manifest.packageName}

go 1.24

require modernc.org/sqlite v1.39.1
`;
}

function manifestGo(manifest) {
  return `package forge

import "encoding/json"

const manifestJSON = ${goRawString(JSON.stringify(manifest, null, 2))}

type Manifest struct {
\tSchemaVersion  string            \`json:"schemaVersion"\`
\tSource         string            \`json:"source"\`
\tName           string            \`json:"name"\`
\tSlug           string            \`json:"slug"\`
\tEnvPrefix      string            \`json:"envPrefix"\`
\tDescription    string            \`json:"description"\`
\tVersion        string            \`json:"version"\`
\tBaseURLs       []string          \`json:"baseUrls"\`
\tAuth           Auth              \`json:"auth"\`
\tTags           []string          \`json:"tags"\`
\tOperations     []Operation       \`json:"operations"\`
\tHeroes         []Hero            \`json:"heroes"\`
\tDefaultHeaders map[string]string \`json:"defaultHeaders"\`
\tPolicy         Policy            \`json:"policy"\`
\tInsights       Insights          \`json:"insights"\`
}

type Hero struct {
\tAlias       string \`json:"alias"\`
\tOperationID string \`json:"operationId"\`
\tSummary     string \`json:"summary"\`
\tMethod      string \`json:"method"\`
\tPath        string \`json:"path"\`
}

type Auth struct {
\tMode    string       \`json:"mode"\`
\tSchemes []AuthScheme \`json:"schemes"\`
\tEnv     string       \`json:"env"\`
\tOAuth   bool         \`json:"oauth"\`
}

type AuthScheme struct {
\tName   string               \`json:"name"\`
\tType   string               \`json:"type"\`
\tIn     string               \`json:"in"\`
\tHeader string               \`json:"header"\`
\tScheme string               \`json:"scheme"\`
\tFlows  map[string]OAuthFlow \`json:"flows"\`
}

type OAuthFlow struct {
\tAuthorizationURL string            \`json:"authorizationUrl"\`
\tTokenURL         string            \`json:"tokenUrl"\`
\tRefreshURL       string            \`json:"refreshUrl"\`
\tScopes           map[string]string \`json:"scopes"\`
}

type Operation struct {
\tID             string       \`json:"id"\`
\tMethod         string       \`json:"method"\`
\tPath           string       \`json:"path"\`
\tTag            string       \`json:"tag"\`
\tSummary        string       \`json:"summary"\`
\tDescription    string       \`json:"description"\`
\tParameters     []Parameter  \`json:"parameters"\`
\tHasRequestBody bool         \`json:"hasRequestBody"\`
\tRisk           string       \`json:"risk"\`
\tCacheable      bool         \`json:"cacheable"\`
\tResponseCodes  []string     \`json:"responseCodes"\`
\tGraphQL        *GraphQLSpec \`json:"graphql,omitempty"\`
\tKind           string       \`json:"kind,omitempty"\`
\tPagination     *Pagination  \`json:"pagination,omitempty"\`
}

type Pagination struct {
\tStyle        string \`json:"style"\`
\tOffsetParam  string \`json:"offsetParam,omitempty"\`
\tLimitParam   string \`json:"limitParam,omitempty"\`
\tCursorParam  string \`json:"cursorParam,omitempty"\`
\tPageParam    string \`json:"pageParam,omitempty"\`
\tPerPageParam string \`json:"perPageParam,omitempty"\`
}

type GraphQLSpec struct {
\tKind  string       \`json:"kind"\`
\tField string       \`json:"field"\`
\tArgs  []GraphQLArg \`json:"args"\`
}

type GraphQLArg struct {
\tName string \`json:"name"\`
\tType string \`json:"type"\`
}

type Parameter struct {
\tName        string \`json:"name"\`
\tIn          string \`json:"in"\`
\tRequired    bool   \`json:"required"\`
\tDescription string \`json:"description"\`
\tType        string \`json:"type"\`
}

type Insights struct {
\tThesis              string   \`json:"thesis"\`
\tGeneratedAdvantages []string \`json:"generatedAdvantages"\`
\tRecommendedCommands []string \`json:"recommendedCommands"\`
\tRiskNotes           []string \`json:"riskNotes"\`
}

type Policy struct {
\tSchemaVersion string       \`json:"schemaVersion"\`
\tRules         []PolicyRule \`json:"rules"\`
\tRedaction     []string     \`json:"redaction"\`
}

type PolicyRule struct {
\tRisk        string \`json:"risk"\`
\tAction      string \`json:"action"\`
\tRequiresYes bool   \`json:"requiresYes"\`
}

func LoadManifest() Manifest {
\tvar manifest Manifest
\tif err := json.Unmarshal([]byte(manifestJSON), &manifest); err != nil {
\t\tpanic(err)
\t}
\treturn manifest
}

func Operations() []Operation {
\treturn LoadManifest().Operations
}

func GetOperation(id string) (Operation, bool) {
\tfor _, operation := range Operations() {
\t\tif operation.ID == id {
\t\t\treturn operation, true
\t\t}
\t}
\treturn Operation{}, false
}

func Heroes() []Hero {
\treturn LoadManifest().Heroes
}

func FindHero(alias string) *Hero {
\tfor _, hero := range Heroes() {
\t\tif hero.Alias == alias {
\t\t\treturn &hero
\t\t}
\t}
\treturn nil
}

func PolicyFor(operation Operation) PolicyRule {
\tmanifest := LoadManifest()
\tfor _, rule := range manifest.Policy.Rules {
\t\tif rule.Risk == operation.Risk {
\t\t\treturn rule
\t\t}
\t}
\tif operation.Risk == "read" {
\t\treturn PolicyRule{Risk: "read", Action: "allow"}
\t}
\treturn PolicyRule{Risk: operation.Risk, Action: "confirm", RequiresYes: true}
}

func RequiresConfirmation(operation Operation) bool {
\trule := PolicyFor(operation)
\treturn rule.RequiresYes || rule.Action == "confirm"
}

func PolicyDenies(operation Operation) bool {
\treturn PolicyFor(operation).Action == "deny"
}
`;
}

function clientGo(manifest) {
  return `package forge

import (
\t"bytes"
\t"context"
\t"encoding/json"
\t"errors"
\t"fmt"
\t"io"
\t"net/http"
\t"net/url"
\t"os"
\t"strings"
\t"time"
)

type CallOptions struct {
\tBaseURL     string
\tAPIKey      string
\tPathParams  map[string]string
\tQueryParams map[string]string
\tHeaders     map[string]string
\tBody        any
\tYes         bool
\tTimeout     time.Duration
}

type RequestPlan struct {
\tMethod  string         \`json:"method"\`
\tURL     string         \`json:"url"\`
\tHeaders map[string]string \`json:"headers"\`
\tBody    any            \`json:"body,omitempty"\`
}

type ResponseEnvelope struct {
\tOK         bool   \`json:"ok"\`
\tStatus     int    \`json:"status"\`
\tStatusText string \`json:"statusText"\`
\tData       any    \`json:"data"\`
}

type CallResult struct {
\tDryRun    bool              \`json:"dryRun"\`
\tOperation Operation         \`json:"operation"\`
\tRequest   RequestPlan       \`json:"request"\`
\tResponse  *ResponseEnvelope \`json:"response,omitempty"\`
\tNote      string            \`json:"note,omitempty"\`
}

func BuildURL(operation Operation, options CallOptions) (string, error) {
\tmanifest := LoadManifest()
\tbase := strings.TrimRight(options.BaseURL, "/")
\tif base == "" {
\t\tbase = strings.TrimRight(os.Getenv(manifest.EnvPrefix+"_BASE_URL"), "/")
\t}
\tif base == "" && len(manifest.BaseURLs) > 0 {
\t\tbase = strings.TrimRight(manifest.BaseURLs[0], "/")
\t}
\tif base == "" {
\t\treturn "", errors.New("missing base URL")
\t}
\tif operation.GraphQL != nil {
\t\treturn base, nil
\t}

\tapiPath := operation.Path
\tfor _, parameter := range operation.Parameters {
\t\tif parameter.In != "path" {
\t\t\tcontinue
\t\t}
\t\tvalue := ""
\t\tif options.PathParams != nil {
\t\t\tvalue = options.PathParams[parameter.Name]
\t\t}
\t\tif value == "" && options.QueryParams != nil {
\t\t\tvalue = options.QueryParams[parameter.Name]
\t\t}
\t\tif value == "" {
\t\t\treturn "", fmt.Errorf("missing path parameter: %s", parameter.Name)
\t\t}
\t\tapiPath = strings.ReplaceAll(apiPath, "{"+parameter.Name+"}", url.PathEscape(value))
\t}

\tparsed, err := url.Parse(base + apiPath)
\tif err != nil {
\t\treturn "", err
\t}
\tquery := parsed.Query()
\tfor _, parameter := range operation.Parameters {
\t\tif parameter.In != "query" {
\t\t\tcontinue
\t\t}
\t\tvalue := ""
\t\tif options.QueryParams != nil {
\t\t\tvalue = options.QueryParams[parameter.Name]
\t\t}
\t\tif value != "" {
\t\t\tquery.Set(parameter.Name, value)
\t\t} else if parameter.Required {
\t\t\treturn "", fmt.Errorf("missing query parameter: %s", parameter.Name)
\t\t}
\t}
\tparsed.RawQuery = query.Encode()
\treturn parsed.String(), nil
}

func AuthHeaders(options CallOptions) map[string]string {
\tmanifest := LoadManifest()
\theaders := map[string]string{}
\tapiKey := options.APIKey
\tif apiKey == "" {
\t\tapiKey = os.Getenv(manifest.EnvPrefix + "_API_KEY")
\t}
\ttoken, _ := loadTokenWithMaybeRefresh()
\tif apiKey == "" || manifest.Auth.Mode == "none" {
\t\tif token.AccessToken != "" {
\t\t\theaders["Authorization"] = "Bearer " + token.AccessToken
\t\t}
\t\treturn headers
\t}
\tif token.AccessToken != "" && manifest.Auth.OAuth {
\t\theaders["Authorization"] = "Bearer " + token.AccessToken
\t\treturn headers
\t}
\tscheme := AuthScheme{}
\tif len(manifest.Auth.Schemes) > 0 {
\t\tscheme = manifest.Auth.Schemes[0]
\t}
\tif scheme.Type == "apiKey" && scheme.In == "header" && scheme.Header != "" {
\t\theaders[scheme.Header] = apiKey
\t} else {
\t\theaders["Authorization"] = "Bearer " + apiKey
\t}
\treturn headers
}

func CallOperation(ctx context.Context, operationID string, options CallOptions) (CallResult, error) {
\tstartedAt := time.Now()
\toperation, ok := GetOperation(operationID)
\tif !ok {
\t\terr := fmt.Errorf("unknown operation: %s", operationID)
\t\tLogCall(operationID, 0, time.Since(startedAt), false, err)
\t\treturn CallResult{}, err
\t}
\trequestURL, err := BuildURL(operation, options)
\tif err != nil {
\t\treturn CallResult{}, err
\t}

\theaders := map[string]string{"Accept": "application/json"}
\tmanifest := LoadManifest()
\tfor key, value := range manifest.DefaultHeaders {
\t\theaders[key] = value
\t}
\tfor key, value := range AuthHeaders(options) {
\t\theaders[key] = value
\t}
\tfor key, value := range options.Headers {
\t\theaders[key] = value
\t}

\tvar bodyBytes []byte
\tvar bodyValue any
\tif options.Body == nil && operation.GraphQL != nil {
\t\toptions.Body = GraphQLPayload(operation, options.QueryParams)
\t}
\tif options.Body != nil {
\t\tbodyBytes, err = json.Marshal(options.Body)
\t\tif err != nil {
\t\t\treturn CallResult{}, err
\t\t}
\t\tbodyValue = options.Body
\t\theaders["Content-Type"] = "application/json"
\t}

\tplan := RequestPlan{Method: operation.Method, URL: requestURL, Headers: RedactHeaders(headers), Body: bodyValue}
\tif PolicyDenies(operation) {
\t\tLogCall(operationID, 0, time.Since(startedAt), true, nil)
\t\treturn CallResult{DryRun: true, Operation: operation, Request: plan, Note: "Policy denies this operation."}, nil
\t}
\tif RequiresConfirmation(operation) && !options.Yes {
\t\tLogCall(operationID, 0, time.Since(startedAt), true, nil)
\t\treturn CallResult{DryRun: true, Operation: operation, Request: plan, Note: "Policy requires --yes for this operation."}, nil
\t}

\ttimeout := options.Timeout
\tif timeout == 0 {
\t\ttimeout = 30 * time.Second
\t}
\tctx, cancel := context.WithTimeout(ctx, timeout)
\tdefer cancel()

\treq, err := http.NewRequestWithContext(ctx, operation.Method, requestURL, bytes.NewReader(bodyBytes))
\tif err != nil {
\t\treturn CallResult{}, err
\t}
\tfor key, value := range headers {
\t\treq.Header.Set(key, value)
\t}
\tresponse, err := DoWithResilience(req)
\tif err != nil {
\t\treturn CallResult{}, err
\t}
\tdefer response.Body.Close()
\tdata, err := DecodeResponse(response.Body, response.Header.Get("Content-Type"))
\tif err != nil {
\t\tLogCall(operationID, response.StatusCode, time.Since(startedAt), false, err)
\t\treturn CallResult{}, err
\t}
\tLogCall(operationID, response.StatusCode, time.Since(startedAt), false, nil)
\treturn CallResult{
\t\tDryRun: false,
\t\tOperation: operation,
\t\tRequest: plan,
\t\tResponse: &ResponseEnvelope{OK: response.StatusCode >= 200 && response.StatusCode < 300, Status: response.StatusCode, StatusText: response.Status, Data: data},
\t}, nil
}

func DecodeResponse(reader io.Reader, contentType string) (any, error) {
\tcontent, err := io.ReadAll(reader)
\tif err != nil {
\t\treturn nil, err
\t}
\tif len(content) == 0 {
\t\treturn nil, nil
\t}
\tvar value any
\tif strings.Contains(contentType, "json") || json.Valid(content) {
\t\tif err := json.Unmarshal(content, &value); err == nil {
\t\t\treturn value, nil
\t\t}
\t}
\treturn string(content), nil
}

func RedactHeaders(headers map[string]string) map[string]string {
\tredacted := map[string]string{}
\tfor key, value := range headers {
\t\tlower := strings.ToLower(key)
\t\tif strings.Contains(lower, "authorization") ||
\t\t\tstrings.Contains(lower, "token") ||
\t\t\tstrings.Contains(lower, "secret") ||
\t\t\tstrings.Contains(lower, "api-key") ||
\t\t\tstrings.Contains(lower, "apikey") ||
\t\t\tstrings.Contains(lower, "subscription-key") ||
\t\t\tstrings.Contains(lower, "subscription_key") ||
\t\t\tstrings.Contains(lower, "ocp-apim-subscription-key") ||
\t\t\tlower == "key" {
\t\t\tredacted[key] = "[redacted]"
\t\t} else {
\t\t\tredacted[key] = value
\t\t}
\t}
\treturn redacted
}

func GraphQLPayload(operation Operation, variables map[string]string) map[string]any {
\tif operation.GraphQL == nil {
\t\treturn map[string]any{}
\t}
\tvariableDefs := []string{}
\tfieldArgs := []string{}
\tfor _, arg := range operation.GraphQL.Args {
\t\targType := arg.Type
\t\tif argType == "" {
\t\t\targType = "String"
\t\t}
\t\tvariableDefs = append(variableDefs, "$"+arg.Name+": "+argType)
\t\tfieldArgs = append(fieldArgs, arg.Name+": $"+arg.Name)
\t}
\toperationName := operation.GraphQL.Field
\tprefix := operation.GraphQL.Kind
\tif prefix == "" {
\t\tprefix = "query"
\t}
\tquery := prefix + " " + operationName
\tif len(variableDefs) > 0 {
\t\tquery += "(" + strings.Join(variableDefs, ", ") + ")"
\t}
\tquery += " { " + operation.GraphQL.Field
\tif len(fieldArgs) > 0 {
\t\tquery += "(" + strings.Join(fieldArgs, ", ") + ")"
\t}
\tquery += " { __typename } }"
\tvars := map[string]any{}
\tfor key, value := range variables {
\t\tvars[key] = value
\t}
\treturn map[string]any{"query": query, "variables": vars}
}
`;
}

function storeGo(manifest) {
  return `package forge

import (
\t"database/sql"
\t"encoding/json"
\t"fmt"
\t"os"
\t"path/filepath"
\t"strings"
\t"time"

\t_ "modernc.org/sqlite"
)

type Cache struct {
\tVersion int           \`json:"version"\`
\tRecords []CacheRecord \`json:"records"\`
}

type CacheRecord struct {
\tOperationID string           \`json:"operationId"\`
\tRequest     RequestPlan      \`json:"request"\`
\tResponse    *ResponseEnvelope \`json:"response,omitempty"\`
\tCachedAt    string           \`json:"cachedAt"\`
}

type ResourceRecord struct {
\tResource    string \`json:"resource"\`
\tOperationID string \`json:"operationId"\`
\tKey         string \`json:"key"\`
\tJSON        string \`json:"json"\`
\tCachedAt    string \`json:"cachedAt"\`
}

type CacheStats struct {
\tFile        string         \`json:"file"\`
\tRecords     int            \`json:"records"\`
\tResources   int            \`json:"resources"\`
\tByOperation map[string]int \`json:"byOperation"\`
\tByResource  map[string]int \`json:"byResource"\`
\tFTS5        bool           \`json:"fts5"\`
}

func cacheFile() string {
\tmanifest := LoadManifest()
\tif override := os.Getenv(manifest.EnvPrefix + "_SQLITE_FILE"); override != "" {
\t\treturn override
\t}
\tif override := os.Getenv(manifest.EnvPrefix + "_CACHE_FILE"); override != "" {
\t\treturn override
\t}
\treturn filepath.Join(".gutenberg", manifest.Slug+".sqlite")
}

func openStore() (*sql.DB, error) {
\tfile := cacheFile()
\tif err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
\t\treturn nil, err
\t}
\tdb, err := sql.Open("sqlite", file)
\tif err != nil {
\t\treturn nil, err
\t}
\tdb.SetMaxOpenConns(1)
\tdb.SetMaxIdleConns(1)
\tif err := migrateStore(db); err != nil {
\t\tdb.Close()
\t\treturn nil, err
\t}
\treturn db, nil
}

// migrations is an append-only list. Each entry is applied exactly once,
// in order. NEVER renumber, edit, or delete a past migration — add a new one.
var migrations = []struct {
\tID   int
\tName string
\tSQL  string
}{
\t{
\t\tID:   1,
\t\tName: "initial",
\t\tSQL: ` + goRawString(`CREATE TABLE records (
\t\t\tid INTEGER PRIMARY KEY AUTOINCREMENT,
\t\t\toperation_id TEXT NOT NULL,
\t\t\trequest_json TEXT NOT NULL,
\t\t\tresponse_json TEXT,
\t\t\tcached_at TEXT NOT NULL
\t\t);
\t\tCREATE VIRTUAL TABLE records_fts USING fts5(operation_id, request_json, response_json, cached_at);`) + `,
\t},
\t{
\t\tID:   2,
\t\tName: "resource-projections",
\t\tSQL: ` + goRawString(`CREATE TABLE resources (
\t\t\tid INTEGER PRIMARY KEY AUTOINCREMENT,
\t\t\tresource TEXT NOT NULL,
\t\t\toperation_id TEXT NOT NULL,
\t\t\tkey TEXT NOT NULL,
\t\t\tjson TEXT NOT NULL,
\t\t\tcached_at TEXT NOT NULL
\t\t);
\t\tCREATE INDEX resources_resource_idx ON resources(resource);
\t\tCREATE VIRTUAL TABLE resources_fts USING fts5(resource, operation_id, key, json, cached_at);`) + `,
\t},
}

func migrateStore(db *sql.DB) error {
\tpragmas := []string{
\t\t"PRAGMA busy_timeout=5000",
\t\t"PRAGMA journal_mode=WAL",
\t\t"PRAGMA synchronous=NORMAL",
\t}
\tfor _, statement := range pragmas {
\t\tif _, err := db.Exec(statement); err != nil {
\t\t\treturn err
\t\t}
\t}
\tif _, err := db.Exec(` + goRawString("CREATE TABLE IF NOT EXISTS gutenberg_migrations (id INTEGER PRIMARY KEY, name TEXT NOT NULL, applied_at TEXT NOT NULL)") + `); err != nil {
\t\treturn err
\t}
\tapplied := map[int]bool{}
\trows, err := db.Query("SELECT id FROM gutenberg_migrations")
\tif err != nil {
\t\treturn err
\t}
\tfor rows.Next() {
\t\tvar id int
\t\tif err := rows.Scan(&id); err != nil {
\t\t\trows.Close()
\t\t\treturn err
\t\t}
\t\tapplied[id] = true
\t}
\trows.Close()
\tfor _, migration := range migrations {
\t\tif applied[migration.ID] {
\t\t\tcontinue
\t\t}
\t\ttx, err := db.Begin()
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tif _, err := tx.Exec(migration.SQL); err != nil {
\t\t\ttx.Rollback()
\t\t\treturn fmt.Errorf("migration %d (%s) failed: %w", migration.ID, migration.Name, err)
\t\t}
\t\tif _, err := tx.Exec("INSERT INTO gutenberg_migrations (id, name, applied_at) VALUES (?, ?, ?)", migration.ID, migration.Name, time.Now().UTC().Format(time.RFC3339)); err != nil {
\t\t\ttx.Rollback()
\t\t\treturn err
\t\t}
\t\tif err := tx.Commit(); err != nil {
\t\t\treturn err
\t\t}
\t}
\treturn nil
}

// AppliedMigrations returns the list of applied migration IDs (for diagnostics).
func AppliedMigrations() ([]int, error) {
\tdb, err := openStore()
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer db.Close()
\trows, err := db.Query("SELECT id FROM gutenberg_migrations ORDER BY id")
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer rows.Close()
\tids := []int{}
\tfor rows.Next() {
\t\tvar id int
\t\tif err := rows.Scan(&id); err != nil {
\t\t\treturn nil, err
\t\t}
\t\tids = append(ids, id)
\t}
\treturn ids, rows.Err()
}

func ReadCache() (Cache, error) {
\tdb, err := openStore()
\tif err != nil {
\t\treturn Cache{}, err
\t}
\tdefer db.Close()
\trows, err := db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records ORDER BY id")
\tif err != nil {
\t\treturn Cache{}, err
\t}
\tdefer rows.Close()
\tcache := Cache{Version: 1, Records: []CacheRecord{}}
\tfor rows.Next() {
\t\trecord, err := scanCacheRecord(rows)
\t\tif err != nil {
\t\t\treturn Cache{}, err
\t\t}
\t\tcache.Records = append(cache.Records, record)
\t}
\treturn cache, rows.Err()
}

func WriteCache(cache Cache) error {
\tdb, err := openStore()
\tif err != nil {
\t\treturn err
\t}
\tdefer db.Close()
\tif _, err := db.Exec("DELETE FROM records; DELETE FROM records_fts; DELETE FROM resources; DELETE FROM resources_fts"); err != nil {
\t\treturn err
\t}
\tfor _, record := range cache.Records {
\t\tif _, err := saveRecordWithDB(db, record); err != nil {
\t\t\treturn err
\t\t}
\t}
\treturn nil
}

func SaveRecord(record CacheRecord) (CacheRecord, error) {
\tdb, err := openStore()
\tif err != nil {
\t\treturn CacheRecord{}, err
\t}
\tdefer db.Close()
\treturn saveRecordWithDB(db, record)
}

func saveRecordWithDB(db *sql.DB, record CacheRecord) (CacheRecord, error) {
\trecord.CachedAt = time.Now().UTC().Format(time.RFC3339)
\trequestJSON, err := json.Marshal(record.Request)
\tif err != nil {
\t\treturn CacheRecord{}, err
\t}
\tresponseJSON, err := json.Marshal(record.Response)
\tif err != nil {
\t\treturn CacheRecord{}, err
\t}
\tif _, err := db.Exec(
\t\t"INSERT INTO records (operation_id, request_json, response_json, cached_at) VALUES (?, ?, ?, ?)",
\t\trecord.OperationID, string(requestJSON), string(responseJSON), record.CachedAt,
\t); err != nil {
\t\treturn CacheRecord{}, err
\t}
\t_, _ = db.Exec(
\t\t"INSERT INTO records_fts (operation_id, request_json, response_json, cached_at) VALUES (?, ?, ?, ?)",
\t\trecord.OperationID, string(requestJSON), string(responseJSON), record.CachedAt,
\t)
\t_ = saveResourceRows(db, record)
\treturn record, nil
}

func saveResourceRows(db *sql.DB, record CacheRecord) error {
\tif record.Response == nil || record.Response.Data == nil {
\t\treturn nil
\t}
\tresource := resourceName(record.OperationID)
\titems := resourceItems(record.Response.Data)
\tfor index, item := range items {
\t\tpayload, err := json.Marshal(item)
\t\tif err != nil {
\t\t\tcontinue
\t\t}
\t\tkey := resourceKey(item, index)
\t\tif _, err := db.Exec(
\t\t\t"INSERT INTO resources (resource, operation_id, key, json, cached_at) VALUES (?, ?, ?, ?, ?)",
\t\t\tresource, record.OperationID, key, string(payload), record.CachedAt,
\t\t); err != nil {
\t\t\treturn err
\t\t}
\t\t_, _ = db.Exec(
\t\t\t"INSERT INTO resources_fts (resource, operation_id, key, json, cached_at) VALUES (?, ?, ?, ?, ?)",
\t\t\tresource, record.OperationID, key, string(payload), record.CachedAt,
\t\t)
\t}
\treturn nil
}

func resourceName(operationID string) string {
\tif operation, ok := GetOperation(operationID); ok {
\t\tif strings.TrimSpace(operation.Tag) != "" {
\t\t\treturn strings.ToLower(strings.ReplaceAll(operation.Tag, " ", "_"))
\t\t}
\t}
\treturn strings.ToLower(strings.ReplaceAll(operationID, " ", "_"))
}

func resourceItems(data any) []any {
\tswitch value := data.(type) {
\tcase []any:
\t\treturn value
\tcase map[string]any:
\t\tfor _, key := range []string{"data", "items", "results", "records", "events"} {
\t\t\tif nested, ok := value[key].([]any); ok {
\t\t\t\treturn nested
\t\t\t}
\t\t}
\t\treturn []any{value}
\tdefault:
\t\treturn []any{value}
\t}
}

func resourceKey(item any, index int) string {
\tif object, ok := item.(map[string]any); ok {
\t\tfor _, key := range []string{"id", "key", "name", "slug", "uuid", "url"} {
\t\t\tif value, ok := object[key]; ok && value != nil {
\t\t\t\treturn fmt.Sprint(value)
\t\t\t}
\t\t}
\t}
\treturn fmt.Sprintf("%06d", index)
}

func SearchCache(query string) ([]CacheRecord, error) {
\tdb, err := openStore()
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer db.Close()
\tif strings.TrimSpace(query) == "" {
\t\treturn latestRecords(db)
\t}
\trows, err := db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records_fts WHERE records_fts MATCH ? LIMIT 50", ftsQuery(query))
\tif err != nil {
\t\trows, err = db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records WHERE lower(operation_id || ' ' || request_json || ' ' || response_json) LIKE ? LIMIT 50", "%"+strings.ToLower(query)+"%")
\t}
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer rows.Close()
\tresults := []CacheRecord{}
\tfor rows.Next() {
\t\trecord, err := scanCacheRecord(rows)
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tresults = append(results, record)
\t}
\treturn results, rows.Err()
}

func SearchResources(query string) ([]ResourceRecord, error) {
\tdb, err := openStore()
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer db.Close()
\tif strings.TrimSpace(query) == "" {
\t\treturn latestResources(db)
\t}
\trows, err := db.Query("SELECT resource, operation_id, key, json, cached_at FROM resources_fts WHERE resources_fts MATCH ? LIMIT 50", ftsQuery(query))
\tif err != nil {
\t\trows, err = db.Query("SELECT resource, operation_id, key, json, cached_at FROM resources WHERE lower(resource || ' ' || operation_id || ' ' || key || ' ' || json) LIKE ? LIMIT 50", "%"+strings.ToLower(query)+"%")
\t}
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer rows.Close()
\tresults := []ResourceRecord{}
\tfor rows.Next() {
\t\trecord, err := scanResourceRecord(rows)
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tresults = append(results, record)
\t}
\treturn results, rows.Err()
}

func GetCacheStats() (CacheStats, error) {
\tdb, err := openStore()
\tif err != nil {
\t\treturn CacheStats{}, err
\t}
\tdefer db.Close()
\tstats := CacheStats{File: cacheFile(), ByOperation: map[string]int{}, ByResource: map[string]int{}, FTS5: true}
\tif err := db.QueryRow("SELECT COUNT(*) FROM records").Scan(&stats.Records); err != nil {
\t\treturn CacheStats{}, err
\t}
\tif err := db.QueryRow("SELECT COUNT(*) FROM resources").Scan(&stats.Resources); err != nil {
\t\treturn CacheStats{}, err
\t}
\trows, err := db.Query("SELECT operation_id, COUNT(*) FROM records GROUP BY operation_id ORDER BY operation_id")
\tif err != nil {
\t\treturn CacheStats{}, err
\t}
\tdefer rows.Close()
\tfor rows.Next() {
\t\tvar operation string
\t\tvar count int
\t\tif err := rows.Scan(&operation, &count); err != nil {
\t\t\treturn CacheStats{}, err
\t\t}
\t\tstats.ByOperation[operation] = count
\t}
\trows.Close()
\trows, err = db.Query("SELECT resource, COUNT(*) FROM resources GROUP BY resource ORDER BY resource")
\tif err != nil {
\t\treturn CacheStats{}, err
\t}
\tdefer rows.Close()
\tfor rows.Next() {
\t\tvar resource string
\t\tvar count int
\t\tif err := rows.Scan(&resource, &count); err != nil {
\t\t\treturn CacheStats{}, err
\t\t}
\t\tstats.ByResource[resource] = count
\t}
\treturn stats, rows.Err()
}

func latestRecords(db *sql.DB) ([]CacheRecord, error) {
\trows, err := db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records ORDER BY id DESC LIMIT 50")
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer rows.Close()
\tresults := []CacheRecord{}
\tfor rows.Next() {
\t\trecord, err := scanCacheRecord(rows)
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tresults = append(results, record)
\t}
\treturn results, rows.Err()
}

func latestResources(db *sql.DB) ([]ResourceRecord, error) {
\trows, err := db.Query("SELECT resource, operation_id, key, json, cached_at FROM resources ORDER BY id DESC LIMIT 50")
\tif err != nil {
\t\treturn nil, err
\t}
\tdefer rows.Close()
\tresults := []ResourceRecord{}
\tfor rows.Next() {
\t\trecord, err := scanResourceRecord(rows)
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tresults = append(results, record)
\t}
\treturn results, rows.Err()
}

func scanCacheRecord(scanner interface{ Scan(dest ...any) error }) (CacheRecord, error) {
\tvar operationID, requestJSON, responseJSON, cachedAt string
\tif err := scanner.Scan(&operationID, &requestJSON, &responseJSON, &cachedAt); err != nil {
\t\treturn CacheRecord{}, err
\t}
\trecord := CacheRecord{OperationID: operationID, CachedAt: cachedAt}
\tif err := json.Unmarshal([]byte(requestJSON), &record.Request); err != nil {
\t\treturn CacheRecord{}, err
\t}
\tif responseJSON != "" && responseJSON != "null" {
\t\tvar response ResponseEnvelope
\t\tif err := json.Unmarshal([]byte(responseJSON), &response); err != nil {
\t\t\treturn CacheRecord{}, err
\t\t}
\t\trecord.Response = &response
\t}
\treturn record, nil
}

func scanResourceRecord(scanner interface{ Scan(dest ...any) error }) (ResourceRecord, error) {
\tvar record ResourceRecord
\tif err := scanner.Scan(&record.Resource, &record.OperationID, &record.Key, &record.JSON, &record.CachedAt); err != nil {
\t\treturn ResourceRecord{}, err
\t}
\treturn record, nil
}

func ftsQuery(query string) string {
\tparts := strings.Fields(query)
\tif len(parts) == 0 {
\t\treturn "*"
\t}
\tterms := []string{}
\tfor _, part := range parts {
\t\tclean := strings.Map(func(r rune) rune {
\t\t\tif r == '_' || r == '-' || r == '.' || r == '/' || r == ':' {
\t\t\t\treturn ' '
\t\t\t}
\t\t\tif r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
\t\t\t\treturn r
\t\t\t}
\t\t\treturn -1
\t\t}, part)
\t\tfor _, token := range strings.Fields(clean) {
\t\t\tterms = append(terms, fmt.Sprintf("%q*", token))
\t\t}
\t}
\tif len(terms) == 0 {
\t\treturn "*"
\t}
\treturn strings.Join(terms, " OR ")
}
`;
}

function cliMain(manifest) {
  return `package main

import (
\t"bufio"
\t"context"
\t"encoding/json"
\t"errors"
\t"fmt"
\t"io"
\t"net/http"
\t"os"
\t"strconv"
\t"strings"

\t"gutenberg.local/${manifest.slug}/internal/forge"
)

func main() {
\tif err := run(os.Args[1:]); err != nil {
\t\tfmt.Fprintln(os.Stderr, "${manifest.slug}:", err)
\t\tos.Exit(1)
\t}
}

func run(args []string) error {
\tcommand := "help"
\tif len(args) > 0 {
\t\tcommand = args[0]
\t\targs = args[1:]
\t}
\toptions, positionals := parseArgs(args)
\tswitch command {
\tcase "help", "--help", "-h":
\t\tprintHelp()
\t\treturn nil
\tcase "info":
\t\treturn printJSON(forge.LoadManifest())
\tcase "operations":
\t\tif options["json"] == "true" {
\t\t\treturn printJSON(forge.Operations())
\t\t}
\t\tprintOperations(forge.Operations())
\t\treturn nil
\tcase "call":
\t\tif len(positionals) == 0 {
\t\t\treturn fmt.Errorf("missing operation id")
\t\t}
\t\tbody, err := parseBody(options["data"])
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tresult, err := forge.CallOperation(context.Background(), positionals[0], forge.CallOptions{
\t\t\tBaseURL: options["base-url"],
\t\t\tAPIKey: options["api-key"],
\t\t\tPathParams: parsePairs(options["path"]),
\t\t\tQueryParams: parsePairs(options["param"]),
\t\t\tHeaders: parseHeaders(args),
\t\t\tBody: body,
\t\t\tYes: options["yes"] == "true",
\t\t})
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tif options["cache"] == "true" && result.Response != nil {
\t\t\t_, err := forge.SaveRecord(forge.CacheRecord{OperationID: positionals[0], Request: result.Request, Response: result.Response})
\t\t\tif err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t}
\t\tif options["stream"] == "true" {
\t\t\treturn streamCall(context.Background(), positionals[0], options, args, parsePairs(options["path"]), parsePairs(options["param"]), body)
\t\t}
\t\tif options["select"] != "" && result.Response != nil {
\t\t\tselected := jsonpathSelect(result.Response.Data, options["select"])
\t\t\treturn printJSON(selected)
\t\t}
\t\tif options["digest"] == "true" {
\t\t\treturn printDigest(result)
\t\t}
\t\tif options["compact"] == "true" {
\t\t\treturn printCompact(result)
\t\t}
\t\treturn printJSON(result)
\tcase "sync":
\t\tcount := 0
\t\tfor _, operation := range forge.Operations() {
\t\t\tif !operation.Cacheable {
\t\t\t\tcontinue
\t\t\t}
\t\t\tif len(positionals) > 0 && operation.ID != positionals[0] {
\t\t\t\tcontinue
\t\t\t}
\t\t\tresult, err := forge.CallOperation(context.Background(), operation.ID, forge.CallOptions{
\t\t\t\tBaseURL: options["base-url"],
\t\t\t\tAPIKey: options["api-key"],
\t\t\t\tQueryParams: parsePairs(options["param"]),
\t\t\t\tHeaders: parseHeaders(args),
\t\t\t})
\t\t\tif err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tif result.Response != nil {
\t\t\t\t_, err := forge.SaveRecord(forge.CacheRecord{OperationID: operation.ID, Request: result.Request, Response: result.Response})
\t\t\t\tif err != nil {
\t\t\t\t\treturn err
\t\t\t\t}
\t\t\t\tcount++
\t\t\t}
\t\t}
\t\treturn printJSON(map[string]any{"synced": count})
\tcase "walk":
\t\tif len(positionals) == 0 {
\t\t\treturn fmt.Errorf("missing operation id")
\t\t}
\t\tmax := 5
\t\tif v, err := strconv.Atoi(options["max"]); err == nil && v > 0 {
\t\t\tmax = v
\t\t}
\t\treturn walkPaginated(context.Background(), positionals[0], options, args, max)
\tcase "search":
\t\tresults, err := forge.SearchCache(strings.Join(positionals, " "))
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(map[string]any{"results": results})
\tcase "resources":
\t\tresults, err := forge.SearchResources(strings.Join(positionals, " "))
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(map[string]any{"results": results})
\tcase "cache":
\t\tstats, err := forge.GetCacheStats()
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(stats)
\tcase "auth":
\t\treturn runAuth(positionals, options)
\tcase "mcp":
\t\treturn forge.RunMCP(os.Stdin, os.Stdout)
\tcase "heroes":
\t\tif options["json"] == "true" {
\t\t\treturn printJSON(forge.Heroes())
\t\t}
\t\tprintHeroes(forge.Heroes())
\t\treturn nil
\tdefault:
\t\tif hero := forge.FindHero(command); hero != nil {
\t\t\tbody, err := parseBody(options["data"])
\t\t\tif err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tresult, err := forge.CallOperation(context.Background(), hero.OperationID, forge.CallOptions{
\t\t\t\tBaseURL:     options["base-url"],
\t\t\t\tAPIKey:      options["api-key"],
\t\t\t\tPathParams:  parsePairs(options["path"]),
\t\t\t\tQueryParams: parsePairs(options["param"]),
\t\t\t\tHeaders:     parseHeaders(args),
\t\t\t\tBody:        body,
\t\t\t\tYes:         options["yes"] == "true",
\t\t\t})
\t\t\tif err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tif options["select"] != "" && result.Response != nil {
\t\t\t\treturn printJSON(jsonpathSelect(result.Response.Data, options["select"]))
\t\t\t}
\t\t\tif options["json"] == "true" {
\t\t\t\treturn printJSON(result)
\t\t\t}
\t\t\tif options["compact"] == "true" {
\t\t\t\treturn printCompact(result)
\t\t\t}
\t\t\treturn printDigest(result)
\t\t}
\t\treturn fmt.Errorf("unknown command: %s", command)
\t}
}

func printHeroes(heroes []forge.Hero) {
\tfmt.Printf("%-20s %-32s %s\\n", "alias", "operationId", "summary")
\tfmt.Printf("%-20s %-32s %s\\n", strings.Repeat("-", 20), strings.Repeat("-", 32), strings.Repeat("-", 8))
\tfor _, hero := range heroes {
\t\tfmt.Printf("%-20s %-32s %s\\n", hero.Alias, hero.OperationID, hero.Summary)
\t}
}

func printHelp() {
\tfmt.Print(\`${manifest.name} (${manifest.slug})

Commands:
  help                         Show this help
  info                         Show manifest
  operations [--json]          List operations
  call <operation> [options]   Call an operation
  sync [operation] [options]   Cache read operations locally
  walk <operation> [--max N]   Iterate paginated GET endpoint until empty/limit
  search <query>               Search cached responses
  resources [query]            Search projected domain resource rows
  cache                        Show cache stats
  auth <status|config|logout|client-credentials|device|pkce-start|pkce-finish|refresh>
  mcp                          Start MCP stdio server
  heroes [--json]              List auto-detected zero-friction commands
  <alias> [options]            Shortcut: <tool> nba == call <op-id> --digest

Options:
  --base-url <url>
  --api-key <key>
  --param name=value
  --path name=value
  --header 'k: v'    Add a request header (repeatable, takes priority over defaultHeaders)
  --select 'a.b[*].c' JSONPath subset to project response (jq-lite)
  --stream            For text/event-stream and ndjson: print line-by-line
  --data '{"key":"value"}'
  --cache
  --compact          One-line per top-level field (count + first-item hint for lists)
  --digest           Structured JSON, first 3 items per list, strings truncated to 80 chars
  --yes
\`)
}

func runAuth(positionals []string, options map[string]string) error {
\taction := "status"
\tif len(positionals) > 0 {
\t\taction = positionals[0]
\t}
\tswitch action {
\tcase "status":
\t\treturn printJSON(forge.OAuthStatus())
\tcase "config":
\t\treturn printJSON(forge.OAuthConfig())
\tcase "logout":
\t\treturn forge.Logout()
\tcase "client-credentials":
\t\ttoken, err := forge.ClientCredentials(context.Background(), options["token-url"], options["client-id"], options["client-secret"], options["scope"])
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tif err := forge.SaveStoredToken(token); err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(forge.OAuthStatus())
\tcase "device":
\t\tdevice, token, err := forge.DeviceCode(context.Background(), options["device-url"], options["token-url"], options["client-id"], options["scope"])
\t\tif device.UserCode != "" {
\t\t\tfmt.Fprintf(os.Stderr, "Open %s and enter code %s\\n", firstNonEmpty(device.VerificationURIComplete, device.VerificationURI), device.UserCode)
\t\t}
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\ttoken.TokenURL = options["token-url"]
\t\ttoken.ClientID = options["client-id"]
\t\tif err := forge.SaveStoredToken(token); err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(forge.OAuthStatus())
\tcase "pkce-start":
\t\tverifier, authURL, err := forge.PKCEStart(options["auth-url"], options["token-url"], options["client-id"], options["redirect-uri"], options["scope"])
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(map[string]any{"authorizationUrl": authURL, "codeVerifierStored": verifier != "", "next": "Open the URL, complete login, copy the 'code' query parameter, then run 'auth pkce-finish --code <code>'."})
\tcase "pkce-finish":
\t\tcode := options["code"]
\t\tif code == "" {
\t\t\treturn errors.New("missing --code")
\t\t}
\t\ttoken, err := forge.PKCEFinish(context.Background(), code)
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tif err := forge.SaveStoredToken(token); err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(forge.OAuthStatus())
\tcase "refresh":
\t\ttoken, err := forge.MaybeRefreshStored(context.Background())
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\treturn printJSON(map[string]any{"refreshed": true, "expiresAt": token.ExpiresAt(), "tokenFile": forge.TokenFile()})
\tdefault:
\t\treturn fmt.Errorf("unknown auth action: %s", action)
\t}
}

func firstNonEmpty(values ...string) string {
\tfor _, value := range values {
\t\tif value != "" {
\t\t\treturn value
\t\t}
\t}
\treturn ""
}

func parseArgs(args []string) (map[string]string, []string) {
\toptions := map[string]string{}
\tpositionals := []string{}
\tfor index := 0; index < len(args); index++ {
\t\titem := args[index]
\t\tif !strings.HasPrefix(item, "--") {
\t\t\tpositionals = append(positionals, item)
\t\t\tcontinue
\t\t}
\t\tkeyValue := strings.TrimPrefix(item, "--")
\t\tkey := keyValue
\t\tvalue := "true"
\t\tif strings.Contains(keyValue, "=") {
\t\t\tparts := strings.SplitN(keyValue, "=", 2)
\t\t\tkey = parts[0]
\t\t\tvalue = parts[1]
\t\t} else if index+1 < len(args) && !strings.HasPrefix(args[index+1], "--") {
\t\t\tvalue = args[index+1]
\t\t\tindex++
\t\t}
\t\tif existing, ok := options[key]; ok && existing != "" && existing != "true" {
\t\t\toptions[key] = existing + "," + value
\t\t} else {
\t\t\toptions[key] = value
\t\t}
\t}
\treturn options, positionals
}

// walkPaginated iterates an operation across pages using its declared pagination metadata.
func walkPaginated(ctx context.Context, operationID string, options map[string]string, args []string, max int) error {
\toperation, ok := forge.GetOperation(operationID)
\tif !ok {
\t\treturn fmt.Errorf("unknown operation: %s", operationID)
\t}
\tif operation.Pagination == nil {
\t\treturn fmt.Errorf("operation %s has no detected pagination metadata", operationID)
\t}
\tpathParams := parsePairs(options["path"])
\tbaseQuery := parsePairs(options["param"])
\theaders := parseHeaders(args)

\tpages := 0
\tswitch operation.Pagination.Style {
\tcase "offset-limit":
\t\toffset := 0
\t\tlimit := 50
\t\tif v, err := strconv.Atoi(baseQuery[operation.Pagination.LimitParam]); err == nil {
\t\t\tlimit = v
\t\t}
\t\tfor pages = 0; pages < max; pages++ {
\t\t\tq := cloneMap(baseQuery)
\t\t\tq[operation.Pagination.OffsetParam] = strconv.Itoa(offset)
\t\t\tq[operation.Pagination.LimitParam] = strconv.Itoa(limit)
\t\t\tresult, err := forge.CallOperation(ctx, operationID, forge.CallOptions{BaseURL: options["base-url"], APIKey: options["api-key"], PathParams: pathParams, QueryParams: q, Headers: headers})
\t\t\tif err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tif err := printJSON(result.Response.Data); err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tif isEmptyPage(result.Response.Data) {
\t\t\t\treturn nil
\t\t\t}
\t\t\toffset += limit
\t\t}
\tcase "page":
\t\tpage := 1
\t\tfor pages = 0; pages < max; pages++ {
\t\t\tq := cloneMap(baseQuery)
\t\t\tq[operation.Pagination.PageParam] = strconv.Itoa(page)
\t\t\tresult, err := forge.CallOperation(ctx, operationID, forge.CallOptions{BaseURL: options["base-url"], APIKey: options["api-key"], PathParams: pathParams, QueryParams: q, Headers: headers})
\t\t\tif err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tif err := printJSON(result.Response.Data); err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tif isEmptyPage(result.Response.Data) {
\t\t\t\treturn nil
\t\t\t}
\t\t\tpage++
\t\t}
\tcase "cursor":
\t\tcursor := baseQuery[operation.Pagination.CursorParam]
\t\tfor pages = 0; pages < max; pages++ {
\t\t\tq := cloneMap(baseQuery)
\t\t\tif cursor != "" {
\t\t\t\tq[operation.Pagination.CursorParam] = cursor
\t\t\t}
\t\t\tresult, err := forge.CallOperation(ctx, operationID, forge.CallOptions{BaseURL: options["base-url"], APIKey: options["api-key"], PathParams: pathParams, QueryParams: q, Headers: headers})
\t\t\tif err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tif err := printJSON(result.Response.Data); err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t\tnext := extractCursor(result.Response.Data)
\t\t\tif next == "" || next == cursor {
\t\t\t\treturn nil
\t\t\t}
\t\t\tcursor = next
\t\t}
\tdefault:
\t\treturn fmt.Errorf("unsupported pagination style: %s", operation.Pagination.Style)
\t}
\treturn nil
}

func cloneMap(m map[string]string) map[string]string {
\tout := map[string]string{}
\tfor k, v := range m {
\t\tout[k] = v
\t}
\treturn out
}

func isEmptyPage(data any) bool {
\tswitch v := data.(type) {
\tcase []any:
\t\treturn len(v) == 0
\tcase map[string]any:
\t\tfor _, key := range []string{"items", "results", "data", "entries", "list"} {
\t\t\tif arr, ok := v[key].([]any); ok {
\t\t\t\treturn len(arr) == 0
\t\t\t}
\t\t}
\t}
\treturn false
}

func extractCursor(data any) string {
\tobj, ok := data.(map[string]any)
\tif !ok {
\t\treturn ""
\t}
\tfor _, key := range []string{"next_cursor", "nextCursor", "next_page_token", "nextPageToken", "next", "cursor"} {
\t\tif s, ok := obj[key].(string); ok {
\t\t\treturn s
\t\t}
\t}
\treturn ""
}

// streamCall executes an operation expecting text/event-stream or application/x-ndjson
// and prints each event/line as it arrives.
func streamCall(ctx context.Context, operationID string, options map[string]string, args []string, pathParams, queryParams map[string]string, body any) error {
\toperation, ok := forge.GetOperation(operationID)
\tif !ok {
\t\treturn fmt.Errorf("unknown operation: %s", operationID)
\t}
\trequestURL, err := forge.BuildURL(operation, forge.CallOptions{BaseURL: options["base-url"], PathParams: pathParams, QueryParams: queryParams})
\tif err != nil {
\t\treturn err
\t}
\treq, err := http.NewRequestWithContext(ctx, operation.Method, requestURL, nil)
\tif err != nil {
\t\treturn err
\t}
\treq.Header.Set("Accept", "text/event-stream, application/x-ndjson, application/json")
\tfor key, value := range forge.AuthHeaders(forge.CallOptions{APIKey: options["api-key"]}) {
\t\treq.Header.Set(key, value)
\t}
\tfor key, value := range parseHeaders(args) {
\t\treq.Header.Set(key, value)
\t}
\tresp, err := http.DefaultClient.Do(req)
\tif err != nil {
\t\treturn err
\t}
\tdefer resp.Body.Close()
\tif resp.StatusCode >= 400 {
\t\treturn fmt.Errorf("HTTP %d %s", resp.StatusCode, resp.Status)
\t}
\treader := bufio.NewReader(resp.Body)
\tfor {
\t\tline, readErr := reader.ReadBytes('\\n')
\t\tif len(line) > 0 {
\t\t\ttrimmed := strings.TrimSpace(string(line))
\t\t\tif trimmed != "" {
\t\t\t\tfmt.Println(trimmed)
\t\t\t}
\t\t}
\t\tif readErr == io.EOF {
\t\t\treturn nil
\t\t}
\t\tif readErr != nil {
\t\t\treturn readErr
\t\t}
\t}
}

// jsonpathSelect walks 'value' along a dotted/bracket path with [*] wildcard support.
// Examples: "data.items", "data.items[0]", "data.items[*].name".
func jsonpathSelect(value any, expr string) any {
\texpr = strings.TrimPrefix(strings.TrimPrefix(expr, "$"), ".")
\tcursor := []any{value}
\ttoken := strings.Builder{}
\tflush := func() {
\t\tname := token.String()
\t\ttoken.Reset()
\t\tif name == "" {
\t\t\treturn
\t\t}
\t\tnext := []any{}
\t\tfor _, item := range cursor {
\t\t\tif obj, ok := item.(map[string]any); ok {
\t\t\t\tnext = append(next, obj[name])
\t\t\t}
\t\t}
\t\tcursor = next
\t}
\tfor i := 0; i < len(expr); i++ {
\t\tch := expr[i]
\t\tif ch == '.' {
\t\t\tflush()
\t\t\tcontinue
\t\t}
\t\tif ch == '[' {
\t\t\tflush()
\t\t\tend := strings.Index(expr[i:], "]")
\t\t\tif end == -1 {
\t\t\t\tbreak
\t\t\t}
\t\t\tinner := expr[i+1 : i+end]
\t\t\ti += end
\t\t\tif inner == "*" {
\t\t\t\tnext := []any{}
\t\t\t\tfor _, item := range cursor {
\t\t\t\t\tif arr, ok := item.([]any); ok {
\t\t\t\t\t\tnext = append(next, arr...)
\t\t\t\t\t}
\t\t\t\t}
\t\t\t\tcursor = next
\t\t\t\tcontinue
\t\t\t}
\t\t\tidx, err := strconv.Atoi(inner)
\t\t\tif err != nil {
\t\t\t\tcontinue
\t\t\t}
\t\t\tnext := []any{}
\t\t\tfor _, item := range cursor {
\t\t\t\tif arr, ok := item.([]any); ok && idx >= 0 && idx < len(arr) {
\t\t\t\t\tnext = append(next, arr[idx])
\t\t\t\t}
\t\t\t}
\t\t\tcursor = next
\t\t\tcontinue
\t\t}
\t\ttoken.WriteByte(ch)
\t}
\tflush()
\tif len(cursor) == 1 {
\t\treturn cursor[0]
\t}
\treturn cursor
}

func parseHeaders(args []string) map[string]string {
\tout := map[string]string{}
\tfor i := 0; i < len(args); i++ {
\t\tvar value string
\t\tif args[i] == "--header" || args[i] == "-H" {
\t\t\tif i+1 >= len(args) {
\t\t\t\tcontinue
\t\t\t}
\t\t\tvalue = args[i+1]
\t\t\ti++
\t\t} else if strings.HasPrefix(args[i], "--header=") {
\t\t\tvalue = strings.TrimPrefix(args[i], "--header=")
\t\t} else if strings.HasPrefix(args[i], "-H=") {
\t\t\tvalue = strings.TrimPrefix(args[i], "-H=")
\t\t} else {
\t\t\tcontinue
\t\t}
\t\tcolon := strings.Index(value, ":")
\t\tif colon == -1 {
\t\t\tcontinue
\t\t}
\t\tname := strings.TrimSpace(value[:colon])
\t\tval := strings.TrimSpace(value[colon+1:])
\t\tif name != "" {
\t\t\tout[name] = val
\t\t}
\t}
\treturn out
}

func parsePairs(value string) map[string]string {
\toutput := map[string]string{}
\tif value == "" {
\t\treturn output
\t}
\tfor _, item := range strings.Split(value, ",") {
\t\tparts := strings.SplitN(item, "=", 2)
\t\tif len(parts) == 2 {
\t\t\toutput[parts[0]] = parts[1]
\t\t}
\t}
\treturn output
}

func parseBody(value string) (any, error) {
\tif value == "" {
\t\treturn nil, nil
\t}
\tvar body any
\tif err := json.Unmarshal([]byte(value), &body); err != nil {
\t\treturn nil, err
\t}
\treturn body, nil
}

func printJSON(value any) error {
\tencoder := json.NewEncoder(os.Stdout)
\tencoder.SetIndent("", "  ")
\treturn encoder.Encode(value)
}

func printCompact(result forge.CallResult) error {
\tfmt.Printf("%s %s\\n", result.Operation.Method, result.Request.URL)
\tif result.DryRun {
\t\tfmt.Println("status: dry-run")
\t\tif result.Note != "" {
\t\t\tfmt.Println("note: " + result.Note)
\t\t}
\t\treturn nil
\t}
\tif result.Response == nil {
\t\tfmt.Println("status: no response")
\t\treturn nil
\t}
\tfmt.Printf("status: %d %s ok=%t\\n", result.Response.Status, result.Response.StatusText, result.Response.OK)
\tprintCompactValue("data", result.Response.Data)
\treturn nil
}

func printCompactValue(label string, value any) {
\tswitch data := value.(type) {
\tcase map[string]any:
\t\tpreferred := []string{"message", "errors", "data", "events", "results", "items", "leagues", "day"}
\t\tprinted := 0
\t\tfor _, key := range preferred {
\t\t\tif item, ok := data[key]; ok {
\t\t\t\tfmt.Printf("%s.%s: %s\\n", label, key, compactDescription(item))
\t\t\t\tprinted++
\t\t\t}
\t\t}
\t\tif printed == 0 {
\t\t\tcount := 0
\t\t\tfor key, item := range data {
\t\t\t\tif count >= 6 {
\t\t\t\t\tbreak
\t\t\t\t}
\t\t\t\tfmt.Printf("%s.%s: %s\\n", label, key, compactDescription(item))
\t\t\t\tcount++
\t\t\t}
\t\t}
\tcase []any:
\t\tfmt.Printf("%s: %d item(s)\\n", label, len(data))
\t\tfor index, item := range data {
\t\t\tif index >= 5 {
\t\t\t\tbreak
\t\t\t}
\t\t\tfmt.Printf("%s[%d]: %s\\n", label, index, compactDescription(item))
\t\t}
\tcase string:
\t\tif len(data) > 180 {
\t\t\tdata = data[:180] + "..."
\t\t}
\t\tfmt.Printf("%s: %s\\n", label, data)
\tcase nil:
\t\tfmt.Printf("%s: null\\n", label)
\tdefault:
\t\tfmt.Printf("%s: %v\\n", label, data)
\t}
}

func compactDescription(value any) string {
\tswitch data := value.(type) {
\tcase []any:
\t\tif len(data) == 0 {
\t\t\treturn "[]"
\t\t}
\t\treturn fmt.Sprintf("%d item(s); first: %s", len(data), compactDescription(data[0]))
\tcase map[string]any:
\t\tfor _, key := range []string{"name", "shortName", "displayName", "title", "message", "id", "key"} {
\t\t\tif text, ok := data[key].(string); ok && text != "" {
\t\t\t\treturn text
\t\t\t}
\t\t}
\t\treturn fmt.Sprintf("%d field(s)", len(data))
\tcase string:
\t\tif len(data) > 80 {
\t\t\treturn data[:80] + "..."
\t\t}
\t\treturn data
\tcase nil:
\t\treturn "null"
\tdefault:
\t\treturn fmt.Sprintf("%v", data)
\t}
}

// printDigest emits a structured but truncated view of a response.
// For each list it pretty-prints the first 3 items, truncating strings to 80 chars.
func printDigest(result forge.CallResult) error {
\tfmt.Printf("%s %s\\n", result.Operation.Method, result.Request.URL)
\tif result.DryRun {
\t\tfmt.Println("status: dry-run")
\t\tif result.Note != "" {
\t\t\tfmt.Println("note: " + result.Note)
\t\t}
\t\treturn nil
\t}
\tif result.Response == nil {
\t\tfmt.Println("status: no response")
\t\treturn nil
\t}
\tfmt.Printf("status: %d %s ok=%t\\n", result.Response.Status, result.Response.StatusText, result.Response.OK)
\tdigest := digestValue(result.Response.Data, 0)
\tencoder := json.NewEncoder(os.Stdout)
\tencoder.SetIndent("", "  ")
\treturn encoder.Encode(digest)
}

func digestValue(value any, depth int) any {
\tif depth > 6 {
\t\treturn "[truncated: depth limit]"
\t}
\tswitch data := value.(type) {
\tcase []any:
\t\tlimit := 3
\t\tout := []any{}
\t\tfor index, item := range data {
\t\t\tif index >= limit {
\t\t\t\tout = append(out, fmt.Sprintf("…+%d more", len(data)-limit))
\t\t\t\tbreak
\t\t\t}
\t\t\tout = append(out, digestValue(item, depth+1))
\t\t}
\t\treturn out
\tcase map[string]any:
\t\tout := map[string]any{}
\t\tfor key, item := range data {
\t\t\tout[key] = digestValue(item, depth+1)
\t\t}
\t\treturn out
\tcase string:
\t\tif len(data) > 80 {
\t\t\treturn data[:80] + "…"
\t\t}
\t\treturn data
\tdefault:
\t\treturn data
\t}
}

func printOperations(operations []forge.Operation) {
\tfmt.Printf("%-18s %-7s %-32s %-12s %s\\n", "id", "method", "path", "risk", "summary")
\tfmt.Printf("%-18s %-7s %-32s %-12s %s\\n", strings.Repeat("-", 18), strings.Repeat("-", 7), strings.Repeat("-", 32), strings.Repeat("-", 12), strings.Repeat("-", 16))
\tfor _, operation := range operations {
\t\tfmt.Printf("%-18s %-7s %-32s %-12s %s\\n", operation.ID, operation.Method, operation.Path, operation.Risk, operation.Summary)
\t}
}
`;
}

function authGo(manifest) {
  return `package forge

import (
\t"context"
\t"crypto/aes"
\t"crypto/cipher"
\t"crypto/rand"
\t"crypto/sha256"
\t"encoding/base64"
\t"encoding/json"
\t"errors"
\t"fmt"
\t"net/http"
\t"net/url"
\t"os"
\t"path/filepath"
\t"strings"
\t"time"
)

type OAuthToken struct {
\tAccessToken  string    \`json:"access_token"\`
\tRefreshToken string    \`json:"refresh_token,omitempty"\`
\tTokenType    string    \`json:"token_type,omitempty"\`
\tScope        string    \`json:"scope,omitempty"\`
\tExpiresIn    int       \`json:"expires_in,omitempty"\`
\tObtainedAt   time.Time \`json:"obtained_at"\`
\tTokenURL     string    \`json:"token_url,omitempty"\`
\tClientID     string    \`json:"client_id,omitempty"\`
}

func (t OAuthToken) ExpiresAt() time.Time {
\tif t.ExpiresIn == 0 || t.ObtainedAt.IsZero() {
\t\treturn time.Time{}
\t}
\treturn t.ObtainedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
}

func (t OAuthToken) NeedsRefresh() bool {
\tif t.RefreshToken == "" || t.TokenURL == "" || t.ClientID == "" {
\t\treturn false
\t}
\texpiresAt := t.ExpiresAt()
\tif expiresAt.IsZero() {
\t\treturn false
\t}
\treturn time.Until(expiresAt) < 60*time.Second
}

type PendingPKCE struct {
\tCodeVerifier string    \`json:"code_verifier"\`
\tClientID     string    \`json:"client_id"\`
\tTokenURL     string    \`json:"token_url"\`
\tRedirectURI  string    \`json:"redirect_uri"\`
\tScope        string    \`json:"scope,omitempty"\`
\tCreatedAt    time.Time \`json:"created_at"\`
}

type DeviceCodeResponse struct {
\tDeviceCode              string \`json:"device_code"\`
\tUserCode                string \`json:"user_code"\`
\tVerificationURI         string \`json:"verification_uri"\`
\tVerificationURIComplete string \`json:"verification_uri_complete"\`
\tExpiresIn               int    \`json:"expires_in"\`
\tInterval                int    \`json:"interval"\`
}

func TokenFile() string {
\tmanifest := LoadManifest()
\tif override := os.Getenv(manifest.EnvPrefix + "_TOKEN_FILE"); override != "" {
\t\treturn override
\t}
\treturn filepath.Join(".gutenberg", manifest.Slug+"-token.json")
}

// VaultFile is the shared token vault used by every Gutenberg-generated tool.
// Override with GUTENBERG_VAULT_FILE. Default: ~/.gutenberg/vault.json (or .gutenberg/vault.json if home is unavailable).
func VaultFile() string {
\tif override := os.Getenv("GUTENBERG_VAULT_FILE"); override != "" {
\t\treturn override
\t}
\thome, err := os.UserHomeDir()
\tif err != nil || home == "" {
\t\treturn filepath.Join(".gutenberg", "vault.json")
\t}
\treturn filepath.Join(home, ".gutenberg", "vault.json")
}

// vaultEnvelope is the on-disk format of the vault. When Encrypted is true,
// Payload contains base64(nonce || ciphertext) sealed with AES-GCM using
// the 32-byte key from GUTENBERG_VAULT_KEY (hex-encoded).
type vaultEnvelope struct {
\tSchemaVersion string \`json:"schemaVersion"\`
\tEncrypted     bool   \`json:"encrypted"\`
\tPayload       string \`json:"payload"\`
}

type vaultData struct {
\tTokens map[string]OAuthToken \`json:"tokens"\`
}

func vaultKey() ([]byte, error) {
\thex := os.Getenv("GUTENBERG_VAULT_KEY")
\tif hex == "" {
\t\treturn nil, nil
\t}
\tkey, err := decodeHex(hex)
\tif err != nil {
\t\treturn nil, err
\t}
\tif len(key) != 32 {
\t\treturn nil, errors.New("GUTENBERG_VAULT_KEY must decode to 32 bytes (use 64 hex chars)")
\t}
\treturn key, nil
}

func decodeHex(value string) ([]byte, error) {
\tvalue = strings.TrimSpace(value)
\tif len(value)%2 != 0 {
\t\treturn nil, errors.New("hex value must have even length")
\t}
\tout := make([]byte, len(value)/2)
\tfor i := 0; i < len(out); i++ {
\t\thi, err := hexNibble(value[i*2])
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tlo, err := hexNibble(value[i*2+1])
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tout[i] = hi<<4 | lo
\t}
\treturn out, nil
}

func hexNibble(c byte) (byte, error) {
\tswitch {
\tcase c >= '0' && c <= '9':
\t\treturn c - '0', nil
\tcase c >= 'a' && c <= 'f':
\t\treturn c - 'a' + 10, nil
\tcase c >= 'A' && c <= 'F':
\t\treturn c - 'A' + 10, nil
\t}
\treturn 0, fmt.Errorf("invalid hex char: %c", c)
}

func loadVault() (vaultData, error) {
\tdata := vaultData{Tokens: map[string]OAuthToken{}}
\tcontent, err := os.ReadFile(VaultFile())
\tif err != nil {
\t\tif os.IsNotExist(err) {
\t\t\treturn data, nil
\t\t}
\t\treturn data, err
\t}
\tvar envelope vaultEnvelope
\tif err := json.Unmarshal(content, &envelope); err != nil {
\t\treturn data, err
\t}
\tif !envelope.Encrypted {
\t\tif err := json.Unmarshal([]byte(envelope.Payload), &data); err != nil && envelope.Payload != "" {
\t\t\treturn data, err
\t\t}
\t\tif data.Tokens == nil {
\t\t\tdata.Tokens = map[string]OAuthToken{}
\t\t}
\t\treturn data, nil
\t}
\tkey, err := vaultKey()
\tif err != nil || key == nil {
\t\treturn data, errors.New("vault is encrypted but GUTENBERG_VAULT_KEY is missing or invalid")
\t}
\traw, err := base64.StdEncoding.DecodeString(envelope.Payload)
\tif err != nil {
\t\treturn data, err
\t}
\tblock, err := aes.NewCipher(key)
\tif err != nil {
\t\treturn data, err
\t}
\tgcm, err := cipher.NewGCM(block)
\tif err != nil {
\t\treturn data, err
\t}
\tif len(raw) < gcm.NonceSize() {
\t\treturn data, errors.New("vault ciphertext too short")
\t}
\tnonce, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
\tplaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
\tif err != nil {
\t\treturn data, err
\t}
\tif err := json.Unmarshal(plaintext, &data); err != nil {
\t\treturn data, err
\t}
\tif data.Tokens == nil {
\t\tdata.Tokens = map[string]OAuthToken{}
\t}
\treturn data, nil
}

func saveVault(data vaultData) error {
\tif data.Tokens == nil {
\t\tdata.Tokens = map[string]OAuthToken{}
\t}
\tplaintext, err := json.MarshalIndent(data, "", "  ")
\tif err != nil {
\t\treturn err
\t}
\tenvelope := vaultEnvelope{SchemaVersion: "gutenberg.vault.v1"}
\tkey, err := vaultKey()
\tif err != nil {
\t\treturn err
\t}
\tif key == nil {
\t\tenvelope.Encrypted = false
\t\tenvelope.Payload = string(plaintext)
\t} else {
\t\tblock, err := aes.NewCipher(key)
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tgcm, err := cipher.NewGCM(block)
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tnonce := make([]byte, gcm.NonceSize())
\t\tif _, err := rand.Read(nonce); err != nil {
\t\t\treturn err
\t\t}
\t\tciphertext := gcm.Seal(nil, nonce, plaintext, nil)
\t\tcombined := append(nonce, ciphertext...)
\t\tenvelope.Encrypted = true
\t\tenvelope.Payload = base64.StdEncoding.EncodeToString(combined)
\t}
\tcontent, err := json.MarshalIndent(envelope, "", "  ")
\tif err != nil {
\t\treturn err
\t}
\tif err := os.MkdirAll(filepath.Dir(VaultFile()), 0o755); err != nil {
\t\treturn err
\t}
\treturn os.WriteFile(VaultFile(), append(content, '\\n'), 0o600)
}

func LoadStoredToken() (OAuthToken, error) {
\tslug := LoadManifest().Slug
\tvault, err := loadVault()
\tif err == nil {
\t\tif token, ok := vault.Tokens[slug]; ok && token.AccessToken != "" {
\t\t\treturn token, nil
\t\t}
\t}
\t// Legacy fallback: per-tool token file.
\tcontent, err := os.ReadFile(TokenFile())
\tif err != nil {
\t\treturn OAuthToken{}, err
\t}
\tvar token OAuthToken
\tif err := json.Unmarshal(content, &token); err != nil {
\t\treturn OAuthToken{}, err
\t}
\treturn token, nil
}

func SaveStoredToken(token OAuthToken) error {
\tif token.ObtainedAt.IsZero() {
\t\ttoken.ObtainedAt = time.Now().UTC()
\t}
\tslug := LoadManifest().Slug
\tvault, err := loadVault()
\tif err != nil {
\t\treturn err
\t}
\tvault.Tokens[slug] = token
\treturn saveVault(vault)
}

func Logout() error {
\tslug := LoadManifest().Slug
\tvault, err := loadVault()
\tif err == nil {
\t\tdelete(vault.Tokens, slug)
\t\tif saveErr := saveVault(vault); saveErr != nil {
\t\t\treturn saveErr
\t\t}
\t}
\tif err := os.Remove(TokenFile()); err != nil && !os.IsNotExist(err) {
\t\treturn err
\t}
\treturn nil
}

func OAuthConfig() map[string]any {
\tmanifest := LoadManifest()
\treturn map[string]any{"tokenFile": TokenFile(), "auth": manifest.Auth}
}

func OAuthStatus() map[string]any {
\ttoken, err := LoadStoredToken()
\tif err != nil {
\t\treturn map[string]any{"authenticated": false, "tokenFile": TokenFile()}
\t}
\treturn map[string]any{
\t\t"authenticated": true,
\t\t"tokenFile": TokenFile(),
\t\t"tokenType": token.TokenType,
\t\t"scope": token.Scope,
\t\t"obtainedAt": token.ObtainedAt,
\t\t"expiresIn": token.ExpiresIn,
\t}
}

func ClientCredentials(ctx context.Context, tokenURL, clientID, clientSecret, scope string) (OAuthToken, error) {
\tif tokenURL == "" || clientID == "" || clientSecret == "" {
\t\treturn OAuthToken{}, errors.New("token-url, client-id and client-secret are required")
\t}
\tvalues := url.Values{}
\tvalues.Set("grant_type", "client_credentials")
\tvalues.Set("client_id", clientID)
\tvalues.Set("client_secret", clientSecret)
\tif scope != "" {
\t\tvalues.Set("scope", scope)
\t}
\ttoken, err := postTokenForm(ctx, tokenURL, values)
\tif err != nil {
\t\treturn token, err
\t}
\ttoken.TokenURL = tokenURL
\ttoken.ClientID = clientID
\treturn token, nil
}

// GenerateCodeVerifier returns a cryptographically random PKCE code_verifier (RFC 7636).
func GenerateCodeVerifier() (string, error) {
\tbuf := make([]byte, 32)
\tif _, err := rand.Read(buf); err != nil {
\t\treturn "", err
\t}
\treturn base64.RawURLEncoding.EncodeToString(buf), nil
}

// CodeChallengeS256 returns the PKCE S256 code_challenge for a verifier.
func CodeChallengeS256(verifier string) string {
\tsum := sha256.Sum256([]byte(verifier))
\treturn base64.RawURLEncoding.EncodeToString(sum[:])
}

func pendingPKCEFile() string {
\treturn TokenFile() + ".pkce"
}

func savePendingPKCE(p PendingPKCE) error {
\tif err := os.MkdirAll(filepath.Dir(pendingPKCEFile()), 0o755); err != nil {
\t\treturn err
\t}
\tcontent, err := json.MarshalIndent(p, "", "  ")
\tif err != nil {
\t\treturn err
\t}
\treturn os.WriteFile(pendingPKCEFile(), content, 0o600)
}

func loadPendingPKCE() (PendingPKCE, error) {
\tcontent, err := os.ReadFile(pendingPKCEFile())
\tif err != nil {
\t\treturn PendingPKCE{}, err
\t}
\tvar p PendingPKCE
\tif err := json.Unmarshal(content, &p); err != nil {
\t\treturn PendingPKCE{}, err
\t}
\treturn p, nil
}

// PKCEStart generates a verifier and returns the authorization URL the user should visit.
func PKCEStart(authURL, tokenURL, clientID, redirectURI, scope string) (string, string, error) {
\tif authURL == "" || tokenURL == "" || clientID == "" || redirectURI == "" {
\t\treturn "", "", errors.New("auth-url, token-url, client-id and redirect-uri are required")
\t}
\tverifier, err := GenerateCodeVerifier()
\tif err != nil {
\t\treturn "", "", err
\t}
\tparsed, err := url.Parse(authURL)
\tif err != nil {
\t\treturn "", "", err
\t}
\tquery := parsed.Query()
\tquery.Set("response_type", "code")
\tquery.Set("client_id", clientID)
\tquery.Set("redirect_uri", redirectURI)
\tquery.Set("code_challenge", CodeChallengeS256(verifier))
\tquery.Set("code_challenge_method", "S256")
\tif scope != "" {
\t\tquery.Set("scope", scope)
\t}
\tparsed.RawQuery = query.Encode()
\tif err := savePendingPKCE(PendingPKCE{
\t\tCodeVerifier: verifier,
\t\tClientID:     clientID,
\t\tTokenURL:     tokenURL,
\t\tRedirectURI:  redirectURI,
\t\tScope:        scope,
\t\tCreatedAt:    time.Now().UTC(),
\t}); err != nil {
\t\treturn verifier, "", err
\t}
\treturn verifier, parsed.String(), nil
}

// PKCEFinish exchanges the authorization code for a token using the saved verifier.
func PKCEFinish(ctx context.Context, code string) (OAuthToken, error) {
\tpending, err := loadPendingPKCE()
\tif err != nil {
\t\treturn OAuthToken{}, fmt.Errorf("no pending PKCE session: %w", err)
\t}
\tvalues := url.Values{}
\tvalues.Set("grant_type", "authorization_code")
\tvalues.Set("code", code)
\tvalues.Set("redirect_uri", pending.RedirectURI)
\tvalues.Set("client_id", pending.ClientID)
\tvalues.Set("code_verifier", pending.CodeVerifier)
\ttoken, err := postTokenForm(ctx, pending.TokenURL, values)
\tif err != nil {
\t\treturn token, err
\t}
\ttoken.TokenURL = pending.TokenURL
\ttoken.ClientID = pending.ClientID
\t_ = os.Remove(pendingPKCEFile())
\treturn token, nil
}

// RefreshAccessToken exchanges a refresh_token for a new access_token.
func RefreshAccessToken(ctx context.Context, tokenURL, clientID, refreshToken, scope string) (OAuthToken, error) {
\tif tokenURL == "" || clientID == "" || refreshToken == "" {
\t\treturn OAuthToken{}, errors.New("token-url, client-id and refresh-token are required")
\t}
\tvalues := url.Values{}
\tvalues.Set("grant_type", "refresh_token")
\tvalues.Set("refresh_token", refreshToken)
\tvalues.Set("client_id", clientID)
\tif scope != "" {
\t\tvalues.Set("scope", scope)
\t}
\ttoken, err := postTokenForm(ctx, tokenURL, values)
\tif err != nil {
\t\treturn token, err
\t}
\ttoken.TokenURL = tokenURL
\ttoken.ClientID = clientID
\tif token.RefreshToken == "" {
\t\ttoken.RefreshToken = refreshToken
\t}
\treturn token, nil
}

// loadTokenWithMaybeRefresh loads the token; if it's near expiry and refresh
// metadata is set, it refreshes (controlled by ${manifest.envPrefix || ""}_AUTO_REFRESH or default true).
func loadTokenWithMaybeRefresh() (OAuthToken, error) {
\ttoken, err := LoadStoredToken()
\tif err != nil {
\t\treturn token, err
\t}
\tdisable := os.Getenv("${manifest.envPrefix || "GUTENBERG"}_AUTO_REFRESH") == "0"
\tif disable || !token.NeedsRefresh() {
\t\treturn token, nil
\t}
\trefreshed, refreshErr := MaybeRefreshStored(context.Background())
\tif refreshErr != nil {
\t\treturn token, refreshErr
\t}
\treturn refreshed, nil
}

// MaybeRefreshStored loads the stored token and refreshes it if it's near expiry.
func MaybeRefreshStored(ctx context.Context) (OAuthToken, error) {
\ttoken, err := LoadStoredToken()
\tif err != nil {
\t\treturn token, err
\t}
\tif !token.NeedsRefresh() {
\t\treturn token, nil
\t}
\trefreshed, err := RefreshAccessToken(ctx, token.TokenURL, token.ClientID, token.RefreshToken, token.Scope)
\tif err != nil {
\t\treturn token, err
\t}
\tif err := SaveStoredToken(refreshed); err != nil {
\t\treturn refreshed, err
\t}
\treturn refreshed, nil
}

func DeviceCode(ctx context.Context, deviceURL, tokenURL, clientID, scope string) (DeviceCodeResponse, OAuthToken, error) {
\tif deviceURL == "" || tokenURL == "" || clientID == "" {
\t\treturn DeviceCodeResponse{}, OAuthToken{}, errors.New("device-url, token-url and client-id are required")
\t}
\tvalues := url.Values{}
\tvalues.Set("client_id", clientID)
\tif scope != "" {
\t\tvalues.Set("scope", scope)
\t}
\trequest, err := http.NewRequestWithContext(ctx, "POST", deviceURL, strings.NewReader(values.Encode()))
\tif err != nil {
\t\treturn DeviceCodeResponse{}, OAuthToken{}, err
\t}
\trequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
\tresponse, err := http.DefaultClient.Do(request)
\tif err != nil {
\t\treturn DeviceCodeResponse{}, OAuthToken{}, err
\t}
\tdefer response.Body.Close()
\tif response.StatusCode < 200 || response.StatusCode >= 300 {
\t\treturn DeviceCodeResponse{}, OAuthToken{}, fmt.Errorf("device authorization failed: %s", response.Status)
\t}
\tvar device DeviceCodeResponse
\tif err := json.NewDecoder(response.Body).Decode(&device); err != nil {
\t\treturn DeviceCodeResponse{}, OAuthToken{}, err
\t}
\tinterval := time.Duration(device.Interval) * time.Second
\tif interval == 0 {
\t\tinterval = 5 * time.Second
\t}
\tdeadline := time.Now().Add(time.Duration(device.ExpiresIn) * time.Second)
\tfor time.Now().Before(deadline) {
\t\tselect {
\t\tcase <-ctx.Done():
\t\t\treturn device, OAuthToken{}, ctx.Err()
\t\tcase <-time.After(interval):
\t\t}
\t\tpoll := url.Values{}
\t\tpoll.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
\t\tpoll.Set("device_code", device.DeviceCode)
\t\tpoll.Set("client_id", clientID)
\t\ttoken, err := postTokenForm(ctx, tokenURL, poll)
\t\tif err == nil {
\t\t\treturn device, token, nil
\t\t}
\t\tif !strings.Contains(err.Error(), "authorization_pending") && !strings.Contains(err.Error(), "slow_down") {
\t\t\treturn device, OAuthToken{}, err
\t\t}
\t\tif strings.Contains(err.Error(), "slow_down") {
\t\t\tinterval += 5 * time.Second
\t\t}
\t}
\treturn device, OAuthToken{}, errors.New("device code expired")
}

func postTokenForm(ctx context.Context, tokenURL string, values url.Values) (OAuthToken, error) {
\trequest, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(values.Encode()))
\tif err != nil {
\t\treturn OAuthToken{}, err
\t}
\trequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
\trequest.Header.Set("Accept", "application/json")
\tresponse, err := http.DefaultClient.Do(request)
\tif err != nil {
\t\treturn OAuthToken{}, err
\t}
\tdefer response.Body.Close()
\tvar payload map[string]any
\tif err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
\t\treturn OAuthToken{}, err
\t}
\tif response.StatusCode < 200 || response.StatusCode >= 300 {
\t\tif value, ok := payload["error"].(string); ok {
\t\t\treturn OAuthToken{}, fmt.Errorf("oauth token error: %s", value)
\t\t}
\t\treturn OAuthToken{}, fmt.Errorf("oauth token error: %s", response.Status)
\t}
\ttoken := OAuthToken{ObtainedAt: time.Now().UTC()}
\ttoken.AccessToken, _ = payload["access_token"].(string)
\ttoken.RefreshToken, _ = payload["refresh_token"].(string)
\ttoken.TokenType, _ = payload["token_type"].(string)
\ttoken.Scope, _ = payload["scope"].(string)
\tif expires, ok := payload["expires_in"].(float64); ok {
\t\ttoken.ExpiresIn = int(expires)
\t}
\tif token.AccessToken == "" {
\t\treturn OAuthToken{}, errors.New("oauth token response did not include access_token")
\t}
\treturn token, nil
}
`;
}

function mcpGo(manifest) {
  return `package forge

import (
\t"bufio"
\t"context"
\t"encoding/json"
\t"fmt"
\t"io"
\t"strconv"
\t"strings"
)

type rpcMessage struct {
\tJSONRPC string         \`json:"jsonrpc,omitempty"\`
\tID      any            \`json:"id,omitempty"\`
\tMethod  string         \`json:"method,omitempty"\`
\tParams  map[string]any \`json:"params,omitempty"\`
}

func RunMCP(input io.Reader, output io.Writer) error {
\treader := bufio.NewReader(input)
\tfor {
\t\tmessage, err := readRPC(reader)
\t\tif err == io.EOF {
\t\t\treturn nil
\t\t}
\t\tif err != nil {
\t\t\treturn err
\t\t}
\t\tresponse, err := handleRPC(message)
\t\tif err != nil {
\t\t\tresponse = map[string]any{"jsonrpc": "2.0", "id": message.ID, "error": map[string]any{"code": -32000, "message": err.Error()}}
\t\t}
\t\tif response != nil {
\t\t\tif err := writeRPC(output, response); err != nil {
\t\t\t\treturn err
\t\t\t}
\t\t}
\t}
}

func handleRPC(message rpcMessage) (map[string]any, error) {
\tswitch message.Method {
\tcase "initialize":
\t\treturn map[string]any{"jsonrpc": "2.0", "id": message.ID, "result": map[string]any{
\t\t\t"protocolVersion": "2024-11-05",
\t\t\t"capabilities": map[string]any{"tools": map[string]any{}},
\t\t\t"serverInfo": map[string]any{"name": "${manifest.slug}-mcp", "version": "0.1.0"},
\t\t}}, nil
\tcase "notifications/initialized":
\t\treturn nil, nil
\tcase "tools/list":
\t\treturn map[string]any{"jsonrpc": "2.0", "id": message.ID, "result": map[string]any{"tools": []map[string]any{
\t\t\t{"name": "${manifest.slug}_operations", "description": "List operations", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
\t\t\t{"name": "${manifest.slug}_call", "description": "Call an operation", "inputSchema": map[string]any{"type": "object", "required": []string{"operationId"}, "properties": map[string]any{"operationId": map[string]any{"type": "string"}, "params": map[string]any{"type": "object"}, "path": map[string]any{"type": "object"}, "yes": map[string]any{"type": "boolean"}}}},
\t\t\t{"name": "${manifest.slug}_search_cache", "description": "Search cache", "inputSchema": map[string]any{"type": "object", "required": []string{"query"}, "properties": map[string]any{"query": map[string]any{"type": "string"}}}},
\t\t}}}, nil
\tcase "tools/call":
\t\tname, _ := message.Params["name"].(string)
\t\targs, _ := message.Params["arguments"].(map[string]any)
\t\tif name == "${manifest.slug}_operations" {
\t\t\treturn toolResult(message.ID, Operations()), nil
\t\t}
\t\tif name == "${manifest.slug}_search_cache" {
\t\t\tquery, _ := args["query"].(string)
\t\t\tresults, err := SearchCache(query)
\t\t\tif err != nil {
\t\t\t\treturn nil, err
\t\t\t}
\t\t\treturn toolResult(message.ID, results), nil
\t\t}
\t\tif name == "${manifest.slug}_call" {
\t\t\toperationID, _ := args["operationId"].(string)
\t\t\tresult, err := CallOperation(context.Background(), operationID, CallOptions{QueryParams: mapAnyToString(args["params"]), PathParams: mapAnyToString(args["path"]), Body: args["body"], Yes: args["yes"] == true})
\t\t\tif err != nil {
\t\t\t\treturn nil, err
\t\t\t}
\t\t\treturn toolResult(message.ID, result), nil
\t\t}
\t\treturn nil, fmt.Errorf("unknown tool: %s", name)
\tdefault:
\t\treturn map[string]any{"jsonrpc": "2.0", "id": message.ID, "error": map[string]any{"code": -32601, "message": "method not found"}}, nil
\t}
}

func toolResult(id any, value any) map[string]any {
\tcontent, _ := json.MarshalIndent(value, "", "  ")
\treturn map[string]any{"jsonrpc": "2.0", "id": id, "result": map[string]any{"content": []map[string]any{{"type": "text", "text": string(content)}}}}
}

func readRPC(reader *bufio.Reader) (rpcMessage, error) {
\tline, err := reader.ReadString('\\n')
\tif err != nil {
\t\treturn rpcMessage{}, err
\t}
\tline = strings.TrimSpace(line)
\tif strings.HasPrefix(line, "{") {
\t\tvar message rpcMessage
\t\treturn message, json.Unmarshal([]byte(line), &message)
\t}
\tif !strings.HasPrefix(strings.ToLower(line), "content-length:") {
\t\treturn rpcMessage{}, fmt.Errorf("unexpected MCP frame header: %s", line)
\t}
\tlength, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:")))
\tif err != nil {
\t\treturn rpcMessage{}, err
\t}
\tfor {
\t\theader, err := reader.ReadString('\\n')
\t\tif err != nil {
\t\t\treturn rpcMessage{}, err
\t\t}
\t\tif strings.TrimSpace(header) == "" {
\t\t\tbreak
\t\t}
\t}
\tbody := make([]byte, length)
\tif _, err := io.ReadFull(reader, body); err != nil {
\t\treturn rpcMessage{}, err
\t}
\tvar message rpcMessage
\treturn message, json.Unmarshal(body, &message)
}

func writeRPC(output io.Writer, message map[string]any) error {
\tbody, err := json.Marshal(message)
\tif err != nil {
\t\treturn err
\t}
\t_, err = fmt.Fprintf(output, "Content-Length: %d\\r\\n\\r\\n%s", len(body), body)
\treturn err
}

func mapAnyToString(value any) map[string]string {
\toutput := map[string]string{}
\tinput, ok := value.(map[string]any)
\tif !ok {
\t\treturn output
\t}
\tfor key, item := range input {
\t\toutput[key] = fmt.Sprint(item)
\t}
\treturn output
}
`;
}

function testGo(manifest) {
  return `package forge

import (
\t"context"
\t"encoding/json"
\t"os"
\t"path/filepath"
\t"strings"
\t"testing"
)

func TestManifestOperations(t *testing.T) {
\tmanifest := LoadManifest()
\tif manifest.Slug != "${manifest.slug}" {
\t\tt.Fatalf("slug = %s", manifest.Slug)
\t}
\tif len(manifest.Operations) == 0 {
\t\tt.Fatal("expected operations")
\t}
}

func TestGoldenSnapshotsParse(t *testing.T) {
\tdir := "testdata/golden"
\tif _, err := os.Stat(dir); os.IsNotExist(err) {
\t\tt.Skip("no golden snapshots yet")
\t}
\tentries, err := os.ReadDir(dir)
\tif err != nil {
\t\tt.Fatalf("read dir: %v", err)
\t}
\tfor _, e := range entries {
\t\tif !strings.HasSuffix(e.Name(), ".json") {
\t\t\tcontinue
\t\t}
\t\tdata, err := os.ReadFile(filepath.Join(dir, e.Name()))
\t\tif err != nil {
\t\t\tt.Errorf("%s: %v", e.Name(), err)
\t\t\tcontinue
\t\t}
\t\tvar value any
\t\tif err := json.Unmarshal(data, &value); err != nil {
\t\t\tt.Errorf("%s: invalid JSON: %v", e.Name(), err)
\t\t}
\t}
}

func TestWriteOperationsDryRun(t *testing.T) {
\tt.Setenv("${manifest.envPrefix}_BASE_URL", "https://example.com")
\tt.Setenv("GUTENBERG_AUDIT_FILE", filepath.Join(t.TempDir(), "audit.jsonl"))
\tfor _, operation := range Operations() {
\t\tif operation.Risk == "read" {
\t\t\tcontinue
\t\t}
\t\tpathParams := map[string]string{}
\t\tqueryParams := map[string]string{}
\t\tfor _, parameter := range operation.Parameters {
\t\t\tif parameter.In == "path" {
\t\t\t\tpathParams[parameter.Name] = "test"
\t\t\t} else if parameter.In == "query" && parameter.Required {
\t\t\t\tqueryParams[parameter.Name] = "test"
\t\t\t}
\t\t}
\t\tresult, err := CallOperation(context.Background(), operation.ID, CallOptions{
\t\t\tBody:        map[string]any{"example": true},
\t\t\tPathParams:  pathParams,
\t\t\tQueryParams: queryParams,
\t\t})
\t\tif err != nil {
\t\t\tt.Fatalf("%s: %v", operation.ID, err)
\t\t}
\t\tif !result.DryRun {
\t\t\tt.Fatalf("%s: expected dry-run", operation.ID)
\t\t}
\t\treturn
\t}
}
`;
}

function readme(manifest) {
  return `# ${manifest.name}

Generated by Gutenberg as a Go-first agent-native tool.

## Quick Start

\`\`\`bash
go test ./...
go run ./cmd/${manifest.slug} operations
go run ./cmd/${manifest.slug} call ${manifest.operations[0]?.id || "operationId"}
go run ./cmd/${manifest.slug} mcp
\`\`\`

## Product Thesis

${manifest.insights.thesis}

Write and destructive operations dry-run until \`--yes\` is passed.
`;
}

function cookbook(manifest) {
  const firstRead = manifest.operations.find((operation) => operation.cacheable);
  const firstWrite = manifest.operations.find((operation) => operation.risk !== "read");
  return `# ${manifest.name} Cookbook

\`\`\`bash
go run ./cmd/${manifest.slug} operations
go run ./cmd/${manifest.slug} call ${manifest.operations[0]?.id || "operationId"}
go run ./cmd/${manifest.slug} sync ${firstRead?.id || ""}
go run ./cmd/${manifest.slug} search status
go run ./cmd/${manifest.slug} resources
go run ./cmd/${manifest.slug} call ${firstWrite?.id || manifest.operations[0]?.id || "operationId"} --data '{"example":true}'
go run ./cmd/${manifest.slug} mcp
\`\`\`
`;
}

export function goRawString(value) {
  return "`" + value.replaceAll("`", "` + \"`\" + `") + "`";
}

function telemetryGo(manifest) {
  return `package forge

import (
\t"encoding/json"
\t"os"
\t"path/filepath"
\t"time"
)

// Telemetry: opt-in local JSONL log of operation calls. NEVER sent anywhere.
// Enabled when GUTENBERG_TELEMETRY=1. Override path via GUTENBERG_TELEMETRY_FILE.
//
// Audit: risky operations (write/destructive), including dry-runs, are always
// written to a local JSONL audit file. Override path via GUTENBERG_AUDIT_FILE.

type telemetryEvent struct {
\tTimestamp   string \`json:"ts"\`
\tTool        string \`json:"tool"\`
\tOperationID string \`json:"operationId"\`
\tRisk        string \`json:"risk,omitempty"\`
\tStatus      int    \`json:"status,omitempty"\`
\tElapsedMs   int64  \`json:"elapsedMs,omitempty"\`
\tDryRun      bool   \`json:"dryRun,omitempty"\`
\tError       string \`json:"error,omitempty"\`
}

func telemetryFile() string {
\tif override := os.Getenv("GUTENBERG_TELEMETRY_FILE"); override != "" {
\t\treturn override
\t}
\thome, err := os.UserHomeDir()
\tif err != nil || home == "" {
\t\treturn filepath.Join(".gutenberg", "usage.jsonl")
\t}
\treturn filepath.Join(home, ".gutenberg", "usage.jsonl")
}

func auditFile() string {
\tif override := os.Getenv("GUTENBERG_AUDIT_FILE"); override != "" {
\t\treturn override
\t}
\treturn filepath.Join(".gutenberg", "audit.jsonl")
}

func LogCall(operationID string, status int, elapsed time.Duration, dryRun bool, callErr error) {
\toperation, ok := GetOperation(operationID)
\trisk := ""
\tif ok {
\t\trisk = operation.Risk
\t}
\tevent := telemetryEvent{
\t\tTimestamp:   time.Now().UTC().Format(time.RFC3339),
\t\tTool:        "${manifest.slug}",
\t\tOperationID: operationID,
\t\tRisk:        risk,
\t\tStatus:      status,
\t\tElapsedMs:   elapsed.Milliseconds(),
\t\tDryRun:      dryRun,
\t}
\tif callErr != nil {
\t\tevent.Error = callErr.Error()
\t}
\tif risk != "" && risk != "read" {
\t\twriteEvent(auditFile(), event)
\t}
\tif os.Getenv("GUTENBERG_TELEMETRY") != "1" {
\t\treturn
\t}
\twriteEvent(telemetryFile(), event)
}

func writeEvent(path string, event telemetryEvent) {
\tline, err := json.Marshal(event)
\tif err != nil {
\t\treturn
\t}
\tif err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
\t\treturn
\t}
\tfile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
\tif err != nil {
\t\treturn
\t}
\tdefer file.Close()
\t_, _ = file.Write(append(line, '\\n'))
}
`;
}

function snapshotGo(manifest) {
  return `package forge

import (
\t"bytes"
\t"crypto/sha256"
\t"encoding/base64"
\t"encoding/hex"
\t"encoding/json"
\t"io"
\t"net/http"
\t"os"
\t"path/filepath"
)

type snapshotResponse struct {
\tStatus  int                 \`json:"status"\`
\tHeader  map[string][]string \`json:"header"\`
\tBodyB64 string              \`json:"body"\`
}

type snapshotEntry struct {
\tKey      string           \`json:"key"\`
\tMethod   string           \`json:"method"\`
\tURL      string           \`json:"url"\`
\tResponse snapshotResponse \`json:"response"\`
}

// snapshotMode returns the mode from env: "" (off), "record", or "replay".
func snapshotMode() string {
\treturn os.Getenv("${manifest.envPrefix}_SNAPSHOT_MODE")
}

func snapshotDir() string {
\tif dir := os.Getenv("${manifest.envPrefix}_SNAPSHOT_DIR"); dir != "" {
\t\treturn dir
\t}
\treturn ""
}

func snapshotKey(req *http.Request) string {
\thash := sha256.New()
\thash.Write([]byte(req.Method))
\thash.Write([]byte("\\n"))
\thash.Write([]byte(req.URL.String()))
\tif req.Body != nil {
\t\tbuf, _ := io.ReadAll(req.Body)
\t\treq.Body = io.NopCloser(bytes.NewReader(buf))
\t\thash.Write(buf)
\t}
\treturn hex.EncodeToString(hash.Sum(nil))
}

func snapshotPath(key string) string {
\treturn filepath.Join(snapshotDir(), key+".json")
}

// snapshotIntercept returns a synthesized response if mode=replay and snapshot exists.
func snapshotIntercept(req *http.Request) (*http.Response, bool, error) {
\tdir := snapshotDir()
\tif dir == "" || snapshotMode() != "replay" {
\t\treturn nil, false, nil
\t}
\tkey := snapshotKey(req)
\tcontent, err := os.ReadFile(snapshotPath(key))
\tif err != nil {
\t\treturn nil, false, nil
\t}
\tvar entry snapshotEntry
\tif err := json.Unmarshal(content, &entry); err != nil {
\t\treturn nil, false, err
\t}
\tbody, err := base64.StdEncoding.DecodeString(entry.Response.BodyB64)
\tif err != nil {
\t\treturn nil, false, err
\t}
\tresp := &http.Response{
\t\tStatus:     http.StatusText(entry.Response.Status),
\t\tStatusCode: entry.Response.Status,
\t\tHeader:     entry.Response.Header,
\t\tBody:       io.NopCloser(bytes.NewReader(body)),
\t\tRequest:    req,
\t}
\treturn resp, true, nil
}

// snapshotPersist persists the response body to disk if mode=record. It also
// rewinds resp.Body so the caller can still read it.
func snapshotPersist(req *http.Request, resp *http.Response) {
\tdir := snapshotDir()
\tif dir == "" || snapshotMode() != "record" || resp == nil {
\t\treturn
\t}
\tbody, err := io.ReadAll(resp.Body)
\tif err != nil {
\t\treturn
\t}
\t_ = resp.Body.Close()
\tresp.Body = io.NopCloser(bytes.NewReader(body))
\tentry := snapshotEntry{
\t\tKey:    snapshotKey(req),
\t\tMethod: req.Method,
\t\tURL:    req.URL.String(),
\t\tResponse: snapshotResponse{
\t\t\tStatus:  resp.StatusCode,
\t\t\tHeader:  resp.Header,
\t\t\tBodyB64: base64.StdEncoding.EncodeToString(body),
\t\t},
\t}
\tcontent, err := json.MarshalIndent(entry, "", "  ")
\tif err != nil {
\t\treturn
\t}
\tif err := os.MkdirAll(dir, 0o755); err != nil {
\t\treturn
\t}
\t_ = os.WriteFile(snapshotPath(entry.Key), content, 0o644)
}
`;
}

function resilienceGo(manifest) {
  const env = manifest.envPrefix;
  return `package forge

import (
\t"errors"
\t"math"
\t"math/rand"
\t"net/http"
\t"os"
\t"strconv"
\t"sync"
\t"time"
)

// Resilience: token-bucket rate limiter + retry-with-backoff + circuit breaker.
// Configurable via env:
//   ${env}_RATE_LIMIT_RPS   tokens per second (0 disables; default 0)
//   ${env}_RATE_BURST       burst size (default rps or 1)
//   ${env}_RETRY_MAX        max attempts including the first (default 3)
//   ${env}_RETRY_BASE_MS    base backoff in ms (default 200)
//   ${env}_CB_THRESHOLD     consecutive failures to open the breaker (default 5)
//   ${env}_CB_COOLDOWN_MS   cooldown when open (default 30000)

type tokenBucket struct {
\tmu       sync.Mutex
\trate     float64
\tburst    float64
\ttokens   float64
\tlast     time.Time
}

func newTokenBucket(rps, burst float64) *tokenBucket {
\tif burst < 1 {
\t\tburst = 1
\t}
\treturn &tokenBucket{rate: rps, burst: burst, tokens: burst, last: time.Now()}
}

func (b *tokenBucket) Take() {
\tif b == nil || b.rate <= 0 {
\t\treturn
\t}
\tfor {
\t\tb.mu.Lock()
\t\tnow := time.Now()
\t\telapsed := now.Sub(b.last).Seconds()
\t\tb.tokens = math.Min(b.burst, b.tokens+elapsed*b.rate)
\t\tb.last = now
\t\tif b.tokens >= 1 {
\t\t\tb.tokens--
\t\t\tb.mu.Unlock()
\t\t\treturn
\t\t}
\t\tneeded := (1 - b.tokens) / b.rate
\t\tb.mu.Unlock()
\t\ttime.Sleep(time.Duration(needed * float64(time.Second)))
\t}
}

type circuitBreaker struct {
\tmu        sync.Mutex
\tthreshold int
\tcooldown  time.Duration
\tfailures  int
\topenedAt  time.Time
}

func (c *circuitBreaker) Allow() error {
\tif c == nil || c.threshold <= 0 {
\t\treturn nil
\t}
\tc.mu.Lock()
\tdefer c.mu.Unlock()
\tif c.failures >= c.threshold {
\t\tif time.Since(c.openedAt) < c.cooldown {
\t\t\treturn errors.New("${manifest.slug}: circuit breaker open")
\t\t}
\t\tc.failures = 0
\t}
\treturn nil
}

func (c *circuitBreaker) onSuccess() {
\tif c == nil {
\t\treturn
\t}
\tc.mu.Lock()
\tdefer c.mu.Unlock()
\tc.failures = 0
}

func (c *circuitBreaker) onFailure() {
\tif c == nil {
\t\treturn
\t}
\tc.mu.Lock()
\tdefer c.mu.Unlock()
\tc.failures++
\tif c.failures >= c.threshold {
\t\tc.openedAt = time.Now()
\t}
}

var (
\tresilienceOnce   sync.Once
\tresilienceBucket *tokenBucket
\tresilienceCB     *circuitBreaker
\tresilienceMax    int
\tresilienceBaseMs int
)

func initResilience() {
\trps := envFloat("${env}_RATE_LIMIT_RPS", 0)
\tburst := envFloat("${env}_RATE_BURST", rps)
\tif burst <= 0 {
\t\tburst = 1
\t}
\tresilienceBucket = newTokenBucket(rps, burst)
\tresilienceMax = envInt("${env}_RETRY_MAX", 3)
\tif resilienceMax < 1 {
\t\tresilienceMax = 1
\t}
\tresilienceBaseMs = envInt("${env}_RETRY_BASE_MS", 200)
\tcbThreshold := envInt("${env}_CB_THRESHOLD", 5)
\tcbCooldownMs := envInt("${env}_CB_COOLDOWN_MS", 30000)
\tresilienceCB = &circuitBreaker{threshold: cbThreshold, cooldown: time.Duration(cbCooldownMs) * time.Millisecond}
}

// DoWithResilience wraps an HTTP call with rate limiting, retries, and a circuit breaker.
func DoWithResilience(req *http.Request) (*http.Response, error) {
\tresilienceOnce.Do(initResilience)
\tif resp, ok, err := snapshotIntercept(req); ok {
\t\treturn resp, err
\t}
\tif err := resilienceCB.Allow(); err != nil {
\t\treturn nil, err
\t}
\tvar lastErr error
\tfor attempt := 1; attempt <= resilienceMax; attempt++ {
\t\tresilienceBucket.Take()
\t\tresp, err := http.DefaultClient.Do(req)
\t\tif resp != nil && err == nil && !shouldRetry(resp.StatusCode) {
\t\t\tsnapshotPersist(req, resp)
\t\t}
\t\tif resp != nil && resp.StatusCode == 429 {
\t\t\twait := parseRetryAfter(resp.Header.Get("Retry-After"))
\t\t\tif wait > 0 && wait <= 5*time.Minute {
\t\t\t\tresp.Body.Close()
\t\t\t\ttime.Sleep(wait)
\t\t\t\tcontinue
\t\t\t}
\t\t}
\t\tif err == nil && !shouldRetry(resp.StatusCode) {
\t\t\tresilienceCB.onSuccess()
\t\t\treturn resp, nil
\t\t}
\t\tif resp != nil && resp.Body != nil {
\t\t\tresp.Body.Close()
\t\t}
\t\tlastErr = err
\t\tresilienceCB.onFailure()
\t\tif attempt == resilienceMax {
\t\t\tbreak
\t\t}
\t\tsleep := backoff(attempt, resilienceBaseMs)
\t\ttime.Sleep(sleep)
\t}
\tif lastErr == nil {
\t\tlastErr = errors.New("${manifest.slug}: request failed after retries")
\t}
\treturn nil, lastErr
}

func shouldRetry(status int) bool {
\treturn status == 429 || (status >= 500 && status <= 599)
}

// parseRetryAfter accepts the HTTP Retry-After header: delta-seconds or HTTP-date.
func parseRetryAfter(value string) time.Duration {
\tif value == "" {
\t\treturn 0
\t}
\tif secs, err := strconv.Atoi(value); err == nil {
\t\treturn time.Duration(secs) * time.Second
\t}
\tif when, err := http.ParseTime(value); err == nil {
\t\tdiff := time.Until(when)
\t\tif diff > 0 {
\t\t\treturn diff
\t\t}
\t}
\treturn 0
}

func backoff(attempt, baseMs int) time.Duration {
\texp := math.Pow(2, float64(attempt-1))
\tjitter := rand.Float64() * 0.5
\treturn time.Duration(float64(baseMs)*exp*(1+jitter)) * time.Millisecond
}

func envFloat(key string, fallback float64) float64 {
\tvalue := os.Getenv(key)
\tif value == "" {
\t\treturn fallback
\t}
\tif parsed, err := strconv.ParseFloat(value, 64); err == nil {
\t\treturn parsed
\t}
\treturn fallback
}

func envInt(key string, fallback int) int {
\tvalue := os.Getenv(key)
\tif value == "" {
\t\treturn fallback
\t}
\tif parsed, err := strconv.Atoi(value); err == nil {
\t\treturn parsed
\t}
\treturn fallback
}
`;
}
