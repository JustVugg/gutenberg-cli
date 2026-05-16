package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type CallOptions struct {
	BaseURL     string
	APIKey      string
	PathParams  map[string]string
	QueryParams map[string]string
	Headers     map[string]string
	Body        any
	Yes         bool
	Timeout     time.Duration
}

type RequestPlan struct {
	Method  string         `json:"method"`
	URL     string         `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    any            `json:"body,omitempty"`
}

type ResponseEnvelope struct {
	OK         bool   `json:"ok"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
	Data       any    `json:"data"`
}

type CallResult struct {
	DryRun    bool              `json:"dryRun"`
	Operation Operation         `json:"operation"`
	Request   RequestPlan       `json:"request"`
	Response  *ResponseEnvelope `json:"response,omitempty"`
	Note      string            `json:"note,omitempty"`
}

func BuildURL(operation Operation, options CallOptions) (string, error) {
	manifest := LoadManifest()
	base := strings.TrimRight(options.BaseURL, "/")
	if base == "" {
		base = strings.TrimRight(os.Getenv(manifest.EnvPrefix+"_BASE_URL"), "/")
	}
	if base == "" && len(manifest.BaseURLs) > 0 {
		base = strings.TrimRight(manifest.BaseURLs[0], "/")
	}
	if base == "" {
		return "", errors.New("missing base URL")
	}
	if operation.GraphQL != nil {
		return base, nil
	}

	apiPath := operation.Path
	for _, parameter := range operation.Parameters {
		if parameter.In != "path" {
			continue
		}
		value := ""
		if options.PathParams != nil {
			value = options.PathParams[parameter.Name]
		}
		if value == "" && options.QueryParams != nil {
			value = options.QueryParams[parameter.Name]
		}
		if value == "" {
			return "", fmt.Errorf("missing path parameter: %s", parameter.Name)
		}
		apiPath = strings.ReplaceAll(apiPath, "{"+parameter.Name+"}", url.PathEscape(value))
	}

	parsed, err := url.Parse(base + apiPath)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	for _, parameter := range operation.Parameters {
		if parameter.In != "query" {
			continue
		}
		value := ""
		if options.QueryParams != nil {
			value = options.QueryParams[parameter.Name]
		}
		if value != "" {
			query.Set(parameter.Name, value)
		} else if parameter.Required {
			return "", fmt.Errorf("missing query parameter: %s", parameter.Name)
		}
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func AuthHeaders(options CallOptions) map[string]string {
	manifest := LoadManifest()
	headers := map[string]string{}
	apiKey := options.APIKey
	if apiKey == "" {
		apiKey = os.Getenv(manifest.EnvPrefix + "_API_KEY")
	}
	token, _ := loadTokenWithMaybeRefresh()
	if apiKey == "" || manifest.Auth.Mode == "none" {
		if token.AccessToken != "" {
			headers["Authorization"] = "Bearer " + token.AccessToken
		}
		return headers
	}
	if token.AccessToken != "" && manifest.Auth.OAuth {
		headers["Authorization"] = "Bearer " + token.AccessToken
		return headers
	}
	scheme := AuthScheme{}
	if len(manifest.Auth.Schemes) > 0 {
		scheme = manifest.Auth.Schemes[0]
	}
	if scheme.Type == "apiKey" && scheme.In == "header" && scheme.Header != "" {
		headers[scheme.Header] = apiKey
	} else {
		headers["Authorization"] = "Bearer " + apiKey
	}
	return headers
}

func CallOperation(ctx context.Context, operationID string, options CallOptions) (CallResult, error) {
	startedAt := time.Now()
	operation, ok := GetOperation(operationID)
	if !ok {
		err := fmt.Errorf("unknown operation: %s", operationID)
		LogCall(operationID, 0, time.Since(startedAt), false, err)
		return CallResult{}, err
	}
	requestURL, err := BuildURL(operation, options)
	if err != nil {
		return CallResult{}, err
	}

	headers := map[string]string{"Accept": "application/json"}
	manifest := LoadManifest()
	for key, value := range manifest.DefaultHeaders {
		headers[key] = value
	}
	for key, value := range AuthHeaders(options) {
		headers[key] = value
	}
	for key, value := range options.Headers {
		headers[key] = value
	}

	var bodyBytes []byte
	var bodyValue any
	if options.Body == nil && operation.GraphQL != nil {
		options.Body = GraphQLPayload(operation, options.QueryParams)
	}
	if options.Body != nil {
		bodyBytes, err = json.Marshal(options.Body)
		if err != nil {
			return CallResult{}, err
		}
		bodyValue = options.Body
		headers["Content-Type"] = "application/json"
	}

	plan := RequestPlan{Method: operation.Method, URL: requestURL, Headers: RedactHeaders(headers), Body: bodyValue}
	if PolicyDenies(operation) {
		LogCall(operationID, 0, time.Since(startedAt), true, nil)
		return CallResult{DryRun: true, Operation: operation, Request: plan, Note: "Policy denies this operation."}, nil
	}
	if RequiresConfirmation(operation) && !options.Yes {
		LogCall(operationID, 0, time.Since(startedAt), true, nil)
		return CallResult{DryRun: true, Operation: operation, Request: plan, Note: "Policy requires --yes for this operation."}, nil
	}

	timeout := options.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, operation.Method, requestURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return CallResult{}, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := DoWithResilience(req)
	if err != nil {
		return CallResult{}, err
	}
	defer response.Body.Close()
	data, err := DecodeResponse(response.Body, response.Header.Get("Content-Type"))
	if err != nil {
		LogCall(operationID, response.StatusCode, time.Since(startedAt), false, err)
		return CallResult{}, err
	}
	LogCall(operationID, response.StatusCode, time.Since(startedAt), false, nil)
	return CallResult{
		DryRun: false,
		Operation: operation,
		Request: plan,
		Response: &ResponseEnvelope{OK: response.StatusCode >= 200 && response.StatusCode < 300, Status: response.StatusCode, StatusText: response.Status, Data: data},
	}, nil
}

func DecodeResponse(reader io.Reader, contentType string) (any, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, nil
	}
	var value any
	if strings.Contains(contentType, "json") || json.Valid(content) {
		if err := json.Unmarshal(content, &value); err == nil {
			return value, nil
		}
	}
	return string(content), nil
}

func RedactHeaders(headers map[string]string) map[string]string {
	redacted := map[string]string{}
	for key, value := range headers {
		lower := strings.ToLower(key)
		if strings.Contains(lower, "authorization") ||
			strings.Contains(lower, "token") ||
			strings.Contains(lower, "secret") ||
			strings.Contains(lower, "api-key") ||
			strings.Contains(lower, "apikey") ||
			strings.Contains(lower, "subscription-key") ||
			strings.Contains(lower, "subscription_key") ||
			strings.Contains(lower, "ocp-apim-subscription-key") ||
			lower == "key" {
			redacted[key] = "[redacted]"
		} else {
			redacted[key] = value
		}
	}
	return redacted
}

func GraphQLPayload(operation Operation, variables map[string]string) map[string]any {
	if operation.GraphQL == nil {
		return map[string]any{}
	}
	variableDefs := []string{}
	fieldArgs := []string{}
	for _, arg := range operation.GraphQL.Args {
		argType := arg.Type
		if argType == "" {
			argType = "String"
		}
		variableDefs = append(variableDefs, "$"+arg.Name+": "+argType)
		fieldArgs = append(fieldArgs, arg.Name+": $"+arg.Name)
	}
	operationName := operation.GraphQL.Field
	prefix := operation.GraphQL.Kind
	if prefix == "" {
		prefix = "query"
	}
	query := prefix + " " + operationName
	if len(variableDefs) > 0 {
		query += "(" + strings.Join(variableDefs, ", ") + ")"
	}
	query += " { " + operation.GraphQL.Field
	if len(fieldArgs) > 0 {
		query += "(" + strings.Join(fieldArgs, ", ") + ")"
	}
	query += " { __typename } }"
	vars := map[string]any{}
	for key, value := range variables {
		vars[key] = value
	}
	return map[string]any{"query": query, "variables": vars}
}
