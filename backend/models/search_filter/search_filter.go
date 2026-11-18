package search_filter

type SearchFilter struct {
	Skip    uint64
	Limit   uint64
	Sort    string
	Order   int
	Deleted string
}
