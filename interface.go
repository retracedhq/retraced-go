package retraced

type EventsPager interface {
	NextPage() error
	TotalPages() int
	HasNextPage() bool
	HasPreviousPage() bool
	CurrentPageNumber() int
	CurrentResults() []*EventNode
	TotalCount() int
}
