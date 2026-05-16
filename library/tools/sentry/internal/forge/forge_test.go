package forge

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestOperations(t *testing.T) {
	manifest := LoadManifest()
	if manifest.Slug != "sentry" {
		t.Fatalf("slug = %s", manifest.Slug)
	}
	if len(manifest.Operations) == 0 {
		t.Fatal("expected operations")
	}
}

func TestGoldenSnapshotsParse(t *testing.T) {
	dir := "testdata/golden"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("no golden snapshots yet")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Errorf("%s: %v", e.Name(), err)
			continue
		}
		var value any
		if err := json.Unmarshal(data, &value); err != nil {
			t.Errorf("%s: invalid JSON: %v", e.Name(), err)
		}
	}
}

func TestWriteOperationsDryRun(t *testing.T) {
	t.Setenv("SENTRY_BASE_URL", "https://example.com")
	t.Setenv("GUTENBERG_AUDIT_FILE", filepath.Join(t.TempDir(), "audit.jsonl"))
	for _, operation := range Operations() {
		if operation.Risk == "read" {
			continue
		}
		pathParams := map[string]string{}
		queryParams := map[string]string{}
		for _, parameter := range operation.Parameters {
			if parameter.In == "path" {
				pathParams[parameter.Name] = "test"
			} else if parameter.In == "query" && parameter.Required {
				queryParams[parameter.Name] = "test"
			}
		}
		result, err := CallOperation(context.Background(), operation.ID, CallOptions{
			Body:        map[string]any{"example": true},
			PathParams:  pathParams,
			QueryParams: queryParams,
		})
		if err != nil {
			t.Fatalf("%s: %v", operation.ID, err)
		}
		if !result.DryRun {
			t.Fatalf("%s: expected dry-run", operation.ID)
		}
		return
	}
}
