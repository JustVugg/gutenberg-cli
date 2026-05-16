package aggr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// CallSource invokes a single source's chosen operation with the given query parameters.
// It returns the parsed JSON response (or raw string) wrapped in an Item.
func CallSource(ctx context.Context, source SourceDescriptor, params map[string]string) (Item, error) {
	mf, err := source.parsedManifest()
	if err != nil {
		return Item{}, err
	}
	op := mf.findOperation(source.Operation)
	if op == nil {
		return Item{}, fmt.Errorf("source %s: unknown operation %q", source.Slug, source.Operation)
	}
	baseURL := os.Getenv(envKey(source.Slug, "BASE_URL"))
	if baseURL == "" && len(mf.BaseURLs) > 0 {
		baseURL = mf.BaseURLs[0]
	}
	if baseURL == "" {
		return Item{}, fmt.Errorf("source %s: missing base URL", source.Slug)
	}

	requestURL, err := buildURL(baseURL, op, params)
	if err != nil {
		return Item{}, err
	}

	req, err := http.NewRequestWithContext(ctx, op.Method, requestURL, nil)
	if err != nil {
		return Item{}, err
	}
	req.Header.Set("Accept", "application/json")
	apiKey := os.Getenv(envKey(source.Slug, "API_KEY"))
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		payload = string(body)
	}
	return Item{
		Source:     source.Slug,
		Operation:  source.Operation,
		Status:     resp.StatusCode,
		URL:        requestURL,
		FetchedAt:  time.Now().UTC(),
		Payload:    payload,
	}, nil
}

// FanOut calls every source in parallel with shared parameters and returns merged results.
func FanOut(ctx context.Context, params map[string]string, parallel int) []Item {
	if parallel <= 0 {
		parallel = 4
	}
	sem := make(chan struct{}, parallel)
	wg := sync.WaitGroup{}
	results := make([]Item, len(Sources))
	for index, source := range Sources {
		index := index
		source := source
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			item, err := CallSource(ctx, source, params)
			if err != nil {
				results[index] = Item{Source: source.Slug, Operation: source.Operation, Error: err.Error()}
				return
			}
			results[index] = item
		}()
	}
	wg.Wait()
	return Rank(Merge(results))
}

func envKey(slug, suffix string) string {
	return strings.ToUpper(strings.ReplaceAll(slug, "-", "_")) + "_" + suffix
}

func buildURL(baseURL string, op *Operation, params map[string]string) (string, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/") + op.Path)
	if err != nil {
		return "", err
	}
	query := u.Query()
	for _, parameter := range op.Parameters {
		if parameter.In != "query" {
			continue
		}
		if value, ok := params[parameter.Name]; ok {
			query.Set(parameter.Name, value)
		} else if parameter.Required {
			return "", fmt.Errorf("missing required query parameter %q", parameter.Name)
		}
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}
