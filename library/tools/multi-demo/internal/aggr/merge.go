package aggr

// Merge combines fan-out results from N sources.
// Default: concatenate results, preserve source order.
// Override this function for domain-specific merging (e.g. dedupe by id).
func Merge(items []Item) []Item {
	out := make([]Item, 0, len(items))
	for _, item := range items {
		if item.Error != "" {
			out = append(out, item)
			continue
		}
		out = append(out, item)
	}
	return out
}
