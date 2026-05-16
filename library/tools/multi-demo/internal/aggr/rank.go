package aggr

// Rank orders the merged items.
// Default strategy: by-source-order.
// Override to sort by price, distance, score, etc.
func Rank(items []Item) []Item {
	return items
}
