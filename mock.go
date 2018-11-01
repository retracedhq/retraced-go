package retraced

import (
	"errors"
	"sync"
)

type MockEventsPager struct {
	currentPage int
	Pages       [][]*EventNode
	sync.Mutex
}

func (p *MockEventsPager) NextPage() error {
	p.Lock()
	defer p.Unlock()
	if !p.HasNextPage() {
		return errors.New("no next page")
	}
	p.currentPage++
	return nil
}

func (p *MockEventsPager) TotalPages() int {
	return len(p.Pages)
}

func (p *MockEventsPager) HasNextPage() bool {
	return p.currentPage < len(p.Pages)-1
}

func (p *MockEventsPager) HasPreviousPage() bool {
	return p.currentPage > 0
}

func (p *MockEventsPager) CurrentPageNumber() int {
	return p.currentPage
}

func (p *MockEventsPager) CurrentResults() []*EventNode {
	return p.Pages[p.currentPage]
}

func (p *MockEventsPager) TotalCount() int {
	var total int
	for _, page := range p.Pages {
		total += len(page)
	}
	return total
}
