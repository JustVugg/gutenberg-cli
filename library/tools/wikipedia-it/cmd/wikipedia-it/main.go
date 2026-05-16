package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gutenberg.local/wikipedia-it/internal/forge"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "wikipedia-it:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	command := "help"
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}
	options, positionals := parseArgs(args)
	switch command {
	case "help", "--help", "-h":
		printHelp()
		return nil
	case "info":
		return printJSON(forge.LoadManifest())
	case "operations":
		if options["json"] == "true" {
			return printJSON(forge.Operations())
		}
		printOperations(forge.Operations())
		return nil
	case "call":
		if len(positionals) == 0 {
			return fmt.Errorf("missing operation id")
		}
		body, err := parseBody(options["data"])
		if err != nil {
			return err
		}
		result, err := forge.CallOperation(context.Background(), positionals[0], forge.CallOptions{
			BaseURL: options["base-url"],
			APIKey: options["api-key"],
			PathParams: parsePairs(options["path"]),
			QueryParams: parsePairs(options["param"]),
			Headers: parseHeaders(args),
			Body: body,
			Yes: options["yes"] == "true",
		})
		if err != nil {
			return err
		}
		if options["cache"] == "true" && result.Response != nil {
			_, err := forge.SaveRecord(forge.CacheRecord{OperationID: positionals[0], Request: result.Request, Response: result.Response})
			if err != nil {
				return err
			}
		}
		if options["stream"] == "true" {
			return streamCall(context.Background(), positionals[0], options, args, parsePairs(options["path"]), parsePairs(options["param"]), body)
		}
		if options["select"] != "" && result.Response != nil {
			selected := jsonpathSelect(result.Response.Data, options["select"])
			return printJSON(selected)
		}
		if options["digest"] == "true" {
			return printDigest(result)
		}
		if options["compact"] == "true" {
			return printCompact(result)
		}
		return printJSON(result)
	case "sync":
		count := 0
		for _, operation := range forge.Operations() {
			if !operation.Cacheable {
				continue
			}
			if len(positionals) > 0 && operation.ID != positionals[0] {
				continue
			}
			result, err := forge.CallOperation(context.Background(), operation.ID, forge.CallOptions{
				BaseURL: options["base-url"],
				APIKey: options["api-key"],
				QueryParams: parsePairs(options["param"]),
				Headers: parseHeaders(args),
			})
			if err != nil {
				return err
			}
			if result.Response != nil {
				_, err := forge.SaveRecord(forge.CacheRecord{OperationID: operation.ID, Request: result.Request, Response: result.Response})
				if err != nil {
					return err
				}
				count++
			}
		}
		return printJSON(map[string]any{"synced": count})
	case "walk":
		if len(positionals) == 0 {
			return fmt.Errorf("missing operation id")
		}
		max := 5
		if v, err := strconv.Atoi(options["max"]); err == nil && v > 0 {
			max = v
		}
		return walkPaginated(context.Background(), positionals[0], options, args, max)
	case "search":
		results, err := forge.SearchCache(strings.Join(positionals, " "))
		if err != nil {
			return err
		}
		return printJSON(map[string]any{"results": results})
	case "cache":
		stats, err := forge.GetCacheStats()
		if err != nil {
			return err
		}
		return printJSON(stats)
	case "auth":
		return runAuth(positionals, options)
	case "mcp":
		return forge.RunMCP(os.Stdin, os.Stdout)
	case "heroes":
		if options["json"] == "true" {
			return printJSON(forge.Heroes())
		}
		printHeroes(forge.Heroes())
		return nil
	default:
		if hero := forge.FindHero(command); hero != nil {
			body, err := parseBody(options["data"])
			if err != nil {
				return err
			}
			result, err := forge.CallOperation(context.Background(), hero.OperationID, forge.CallOptions{
				BaseURL:     options["base-url"],
				APIKey:      options["api-key"],
				PathParams:  parsePairs(options["path"]),
				QueryParams: parsePairs(options["param"]),
				Headers:     parseHeaders(args),
				Body:        body,
				Yes:         options["yes"] == "true",
			})
			if err != nil {
				return err
			}
			if options["select"] != "" && result.Response != nil {
				return printJSON(jsonpathSelect(result.Response.Data, options["select"]))
			}
			if options["json"] == "true" {
				return printJSON(result)
			}
			if options["compact"] == "true" {
				return printCompact(result)
			}
			return printDigest(result)
		}
		return fmt.Errorf("unknown command: %s", command)
	}
}

func printHeroes(heroes []forge.Hero) {
	fmt.Printf("%-20s %-32s %s\n", "alias", "operationId", "summary")
	fmt.Printf("%-20s %-32s %s\n", strings.Repeat("-", 20), strings.Repeat("-", 32), strings.Repeat("-", 8))
	for _, hero := range heroes {
		fmt.Printf("%-20s %-32s %s\n", hero.Alias, hero.OperationID, hero.Summary)
	}
}

func printHelp() {
	fmt.Print(`wikipedia-it (wikipedia-it)

Commands:
  help                         Show this help
  info                         Show manifest
  operations [--json]          List operations
  call <operation> [options]   Call an operation
  sync [operation] [options]   Cache read operations locally
  walk <operation> [--max N]   Iterate paginated GET endpoint until empty/limit
  search <query>               Search cached responses
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
`)
}

func runAuth(positionals []string, options map[string]string) error {
	action := "status"
	if len(positionals) > 0 {
		action = positionals[0]
	}
	switch action {
	case "status":
		return printJSON(forge.OAuthStatus())
	case "config":
		return printJSON(forge.OAuthConfig())
	case "logout":
		return forge.Logout()
	case "client-credentials":
		token, err := forge.ClientCredentials(context.Background(), options["token-url"], options["client-id"], options["client-secret"], options["scope"])
		if err != nil {
			return err
		}
		if err := forge.SaveStoredToken(token); err != nil {
			return err
		}
		return printJSON(forge.OAuthStatus())
	case "device":
		device, token, err := forge.DeviceCode(context.Background(), options["device-url"], options["token-url"], options["client-id"], options["scope"])
		if device.UserCode != "" {
			fmt.Fprintf(os.Stderr, "Open %s and enter code %s\n", firstNonEmpty(device.VerificationURIComplete, device.VerificationURI), device.UserCode)
		}
		if err != nil {
			return err
		}
		token.TokenURL = options["token-url"]
		token.ClientID = options["client-id"]
		if err := forge.SaveStoredToken(token); err != nil {
			return err
		}
		return printJSON(forge.OAuthStatus())
	case "pkce-start":
		verifier, authURL, err := forge.PKCEStart(options["auth-url"], options["token-url"], options["client-id"], options["redirect-uri"], options["scope"])
		if err != nil {
			return err
		}
		return printJSON(map[string]any{"authorizationUrl": authURL, "codeVerifierStored": verifier != "", "next": "Open the URL, complete login, copy the 'code' query parameter, then run 'auth pkce-finish --code <code>'."})
	case "pkce-finish":
		code := options["code"]
		if code == "" {
			return errors.New("missing --code")
		}
		token, err := forge.PKCEFinish(context.Background(), code)
		if err != nil {
			return err
		}
		if err := forge.SaveStoredToken(token); err != nil {
			return err
		}
		return printJSON(forge.OAuthStatus())
	case "refresh":
		token, err := forge.MaybeRefreshStored(context.Background())
		if err != nil {
			return err
		}
		return printJSON(map[string]any{"refreshed": true, "expiresAt": token.ExpiresAt(), "tokenFile": forge.TokenFile()})
	default:
		return fmt.Errorf("unknown auth action: %s", action)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func parseArgs(args []string) (map[string]string, []string) {
	options := map[string]string{}
	positionals := []string{}
	for index := 0; index < len(args); index++ {
		item := args[index]
		if !strings.HasPrefix(item, "--") {
			positionals = append(positionals, item)
			continue
		}
		keyValue := strings.TrimPrefix(item, "--")
		key := keyValue
		value := "true"
		if strings.Contains(keyValue, "=") {
			parts := strings.SplitN(keyValue, "=", 2)
			key = parts[0]
			value = parts[1]
		} else if index+1 < len(args) && !strings.HasPrefix(args[index+1], "--") {
			value = args[index+1]
			index++
		}
		if existing, ok := options[key]; ok && existing != "" && existing != "true" {
			options[key] = existing + "," + value
		} else {
			options[key] = value
		}
	}
	return options, positionals
}

// walkPaginated iterates an operation across pages using its declared pagination metadata.
func walkPaginated(ctx context.Context, operationID string, options map[string]string, args []string, max int) error {
	operation, ok := forge.GetOperation(operationID)
	if !ok {
		return fmt.Errorf("unknown operation: %s", operationID)
	}
	if operation.Pagination == nil {
		return fmt.Errorf("operation %s has no detected pagination metadata", operationID)
	}
	pathParams := parsePairs(options["path"])
	baseQuery := parsePairs(options["param"])
	headers := parseHeaders(args)

	pages := 0
	switch operation.Pagination.Style {
	case "offset-limit":
		offset := 0
		limit := 50
		if v, err := strconv.Atoi(baseQuery[operation.Pagination.LimitParam]); err == nil {
			limit = v
		}
		for pages = 0; pages < max; pages++ {
			q := cloneMap(baseQuery)
			q[operation.Pagination.OffsetParam] = strconv.Itoa(offset)
			q[operation.Pagination.LimitParam] = strconv.Itoa(limit)
			result, err := forge.CallOperation(ctx, operationID, forge.CallOptions{BaseURL: options["base-url"], APIKey: options["api-key"], PathParams: pathParams, QueryParams: q, Headers: headers})
			if err != nil {
				return err
			}
			if err := printJSON(result.Response.Data); err != nil {
				return err
			}
			if isEmptyPage(result.Response.Data) {
				return nil
			}
			offset += limit
		}
	case "page":
		page := 1
		for pages = 0; pages < max; pages++ {
			q := cloneMap(baseQuery)
			q[operation.Pagination.PageParam] = strconv.Itoa(page)
			result, err := forge.CallOperation(ctx, operationID, forge.CallOptions{BaseURL: options["base-url"], APIKey: options["api-key"], PathParams: pathParams, QueryParams: q, Headers: headers})
			if err != nil {
				return err
			}
			if err := printJSON(result.Response.Data); err != nil {
				return err
			}
			if isEmptyPage(result.Response.Data) {
				return nil
			}
			page++
		}
	case "cursor":
		cursor := baseQuery[operation.Pagination.CursorParam]
		for pages = 0; pages < max; pages++ {
			q := cloneMap(baseQuery)
			if cursor != "" {
				q[operation.Pagination.CursorParam] = cursor
			}
			result, err := forge.CallOperation(ctx, operationID, forge.CallOptions{BaseURL: options["base-url"], APIKey: options["api-key"], PathParams: pathParams, QueryParams: q, Headers: headers})
			if err != nil {
				return err
			}
			if err := printJSON(result.Response.Data); err != nil {
				return err
			}
			next := extractCursor(result.Response.Data)
			if next == "" || next == cursor {
				return nil
			}
			cursor = next
		}
	default:
		return fmt.Errorf("unsupported pagination style: %s", operation.Pagination.Style)
	}
	return nil
}

func cloneMap(m map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range m {
		out[k] = v
	}
	return out
}

func isEmptyPage(data any) bool {
	switch v := data.(type) {
	case []any:
		return len(v) == 0
	case map[string]any:
		for _, key := range []string{"items", "results", "data", "entries", "list"} {
			if arr, ok := v[key].([]any); ok {
				return len(arr) == 0
			}
		}
	}
	return false
}

func extractCursor(data any) string {
	obj, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	for _, key := range []string{"next_cursor", "nextCursor", "next_page_token", "nextPageToken", "next", "cursor"} {
		if s, ok := obj[key].(string); ok {
			return s
		}
	}
	return ""
}

// streamCall executes an operation expecting text/event-stream or application/x-ndjson
// and prints each event/line as it arrives.
func streamCall(ctx context.Context, operationID string, options map[string]string, args []string, pathParams, queryParams map[string]string, body any) error {
	operation, ok := forge.GetOperation(operationID)
	if !ok {
		return fmt.Errorf("unknown operation: %s", operationID)
	}
	requestURL, err := forge.BuildURL(operation, forge.CallOptions{BaseURL: options["base-url"], PathParams: pathParams, QueryParams: queryParams})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, operation.Method, requestURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream, application/x-ndjson, application/json")
	for key, value := range forge.AuthHeaders(forge.CallOptions{APIKey: options["api-key"]}) {
		req.Header.Set(key, value)
	}
	for key, value := range parseHeaders(args) {
		req.Header.Set(key, value)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d %s", resp.StatusCode, resp.Status)
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, readErr := reader.ReadBytes('\n')
		if len(line) > 0 {
			trimmed := strings.TrimSpace(string(line))
			if trimmed != "" {
				fmt.Println(trimmed)
			}
		}
		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}

// jsonpathSelect walks 'value' along a dotted/bracket path with [*] wildcard support.
// Examples: "data.items", "data.items[0]", "data.items[*].name".
func jsonpathSelect(value any, expr string) any {
	expr = strings.TrimPrefix(strings.TrimPrefix(expr, "$"), ".")
	cursor := []any{value}
	token := strings.Builder{}
	flush := func() {
		name := token.String()
		token.Reset()
		if name == "" {
			return
		}
		next := []any{}
		for _, item := range cursor {
			if obj, ok := item.(map[string]any); ok {
				next = append(next, obj[name])
			}
		}
		cursor = next
	}
	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if ch == '.' {
			flush()
			continue
		}
		if ch == '[' {
			flush()
			end := strings.Index(expr[i:], "]")
			if end == -1 {
				break
			}
			inner := expr[i+1 : i+end]
			i += end
			if inner == "*" {
				next := []any{}
				for _, item := range cursor {
					if arr, ok := item.([]any); ok {
						next = append(next, arr...)
					}
				}
				cursor = next
				continue
			}
			idx, err := strconv.Atoi(inner)
			if err != nil {
				continue
			}
			next := []any{}
			for _, item := range cursor {
				if arr, ok := item.([]any); ok && idx >= 0 && idx < len(arr) {
					next = append(next, arr[idx])
				}
			}
			cursor = next
			continue
		}
		token.WriteByte(ch)
	}
	flush()
	if len(cursor) == 1 {
		return cursor[0]
	}
	return cursor
}

func parseHeaders(args []string) map[string]string {
	out := map[string]string{}
	for i := 0; i < len(args); i++ {
		var value string
		if args[i] == "--header" || args[i] == "-H" {
			if i+1 >= len(args) {
				continue
			}
			value = args[i+1]
			i++
		} else if strings.HasPrefix(args[i], "--header=") {
			value = strings.TrimPrefix(args[i], "--header=")
		} else if strings.HasPrefix(args[i], "-H=") {
			value = strings.TrimPrefix(args[i], "-H=")
		} else {
			continue
		}
		colon := strings.Index(value, ":")
		if colon == -1 {
			continue
		}
		name := strings.TrimSpace(value[:colon])
		val := strings.TrimSpace(value[colon+1:])
		if name != "" {
			out[name] = val
		}
	}
	return out
}

func parsePairs(value string) map[string]string {
	output := map[string]string{}
	if value == "" {
		return output
	}
	for _, item := range strings.Split(value, ",") {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			output[parts[0]] = parts[1]
		}
	}
	return output
}

func parseBody(value string) (any, error) {
	if value == "" {
		return nil, nil
	}
	var body any
	if err := json.Unmarshal([]byte(value), &body); err != nil {
		return nil, err
	}
	return body, nil
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func printCompact(result forge.CallResult) error {
	fmt.Printf("%s %s\n", result.Operation.Method, result.Request.URL)
	if result.DryRun {
		fmt.Println("status: dry-run")
		if result.Note != "" {
			fmt.Println("note: " + result.Note)
		}
		return nil
	}
	if result.Response == nil {
		fmt.Println("status: no response")
		return nil
	}
	fmt.Printf("status: %d %s ok=%t\n", result.Response.Status, result.Response.StatusText, result.Response.OK)
	printCompactValue("data", result.Response.Data)
	return nil
}

func printCompactValue(label string, value any) {
	switch data := value.(type) {
	case map[string]any:
		preferred := []string{"message", "errors", "data", "events", "results", "items", "leagues", "day"}
		printed := 0
		for _, key := range preferred {
			if item, ok := data[key]; ok {
				fmt.Printf("%s.%s: %s\n", label, key, compactDescription(item))
				printed++
			}
		}
		if printed == 0 {
			count := 0
			for key, item := range data {
				if count >= 6 {
					break
				}
				fmt.Printf("%s.%s: %s\n", label, key, compactDescription(item))
				count++
			}
		}
	case []any:
		fmt.Printf("%s: %d item(s)\n", label, len(data))
		for index, item := range data {
			if index >= 5 {
				break
			}
			fmt.Printf("%s[%d]: %s\n", label, index, compactDescription(item))
		}
	case string:
		if len(data) > 180 {
			data = data[:180] + "..."
		}
		fmt.Printf("%s: %s\n", label, data)
	case nil:
		fmt.Printf("%s: null\n", label)
	default:
		fmt.Printf("%s: %v\n", label, data)
	}
}

func compactDescription(value any) string {
	switch data := value.(type) {
	case []any:
		if len(data) == 0 {
			return "[]"
		}
		return fmt.Sprintf("%d item(s); first: %s", len(data), compactDescription(data[0]))
	case map[string]any:
		for _, key := range []string{"name", "shortName", "displayName", "title", "message", "id", "key"} {
			if text, ok := data[key].(string); ok && text != "" {
				return text
			}
		}
		return fmt.Sprintf("%d field(s)", len(data))
	case string:
		if len(data) > 80 {
			return data[:80] + "..."
		}
		return data
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", data)
	}
}

// printDigest emits a structured but truncated view of a response.
// For each list it pretty-prints the first 3 items, truncating strings to 80 chars.
func printDigest(result forge.CallResult) error {
	fmt.Printf("%s %s\n", result.Operation.Method, result.Request.URL)
	if result.DryRun {
		fmt.Println("status: dry-run")
		if result.Note != "" {
			fmt.Println("note: " + result.Note)
		}
		return nil
	}
	if result.Response == nil {
		fmt.Println("status: no response")
		return nil
	}
	fmt.Printf("status: %d %s ok=%t\n", result.Response.Status, result.Response.StatusText, result.Response.OK)
	digest := digestValue(result.Response.Data, 0)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(digest)
}

func digestValue(value any, depth int) any {
	if depth > 6 {
		return "[truncated: depth limit]"
	}
	switch data := value.(type) {
	case []any:
		limit := 3
		out := []any{}
		for index, item := range data {
			if index >= limit {
				out = append(out, fmt.Sprintf("…+%d more", len(data)-limit))
				break
			}
			out = append(out, digestValue(item, depth+1))
		}
		return out
	case map[string]any:
		out := map[string]any{}
		for key, item := range data {
			out[key] = digestValue(item, depth+1)
		}
		return out
	case string:
		if len(data) > 80 {
			return data[:80] + "…"
		}
		return data
	default:
		return data
	}
}

func printOperations(operations []forge.Operation) {
	fmt.Printf("%-18s %-7s %-32s %-12s %s\n", "id", "method", "path", "risk", "summary")
	fmt.Printf("%-18s %-7s %-32s %-12s %s\n", strings.Repeat("-", 18), strings.Repeat("-", 7), strings.Repeat("-", 32), strings.Repeat("-", 12), strings.Repeat("-", 16))
	for _, operation := range operations {
		fmt.Printf("%-18s %-7s %-32s %-12s %s\n", operation.ID, operation.Method, operation.Path, operation.Risk, operation.Summary)
	}
}
