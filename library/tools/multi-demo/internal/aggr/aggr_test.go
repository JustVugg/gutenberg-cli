package aggr

import "testing"

func TestSourcesNotEmpty(t *testing.T) {
	if len(Sources) == 0 {
		t.Fatal("aggregator has no sources")
	}
}

func TestEachSourceParses(t *testing.T) {
	for _, source := range Sources {
		if _, err := source.parsedManifest(); err != nil {
			t.Fatalf("source %s: %v", source.Slug, err)
		}
	}
}

func TestMergeIsIdempotentForErrors(t *testing.T) {
	items := []Item{{Source: "a", Error: "boom"}, {Source: "b", Status: 200}}
	merged := Merge(items)
	if len(merged) != 2 {
		t.Fatalf("expected 2 items, got %d", len(merged))
	}
}
