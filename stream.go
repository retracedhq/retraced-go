package retraced

import (
	"io"
	"sync"
)

// Stream returns a single event on every Read. It wraps an EventsConnection and
// fetches the next page as needed to fullfill Reads.
type Stream struct {
	ec  EventsPager
	i   int
	mtx sync.Mutex
}

func (c *Client) NewStream(sq *StructuredQuery, mask *EventNodeMask) (*Stream, error) {
	conn, err := c.Query(sq, mask, 1000)
	if err != nil {
		return nil, err
	}
	return &Stream{
		ec: conn,
	}, nil
}

// Read returns the next unread Event or io.EOF if there are no more.
// It is safe for concurrent access.
func (s *Stream) Read() (*EventNode, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.i == len(s.ec.CurrentResults()) {
		if s.ec.HasNextPage() {
			if err := s.ec.NextPage(); err != nil {
				return nil, err
			}
			s.i = 0
		}
		if s.i == len(s.ec.CurrentResults()) {
			return nil, io.EOF
		}
	}
	event := s.ec.CurrentResults()[s.i]
	s.i++
	return event, nil
}
