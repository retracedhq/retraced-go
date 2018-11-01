package retraced

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamFixPanic(t *testing.T) {
	s := &Stream{
		ec: &MockEventsPager{
			Pages: [][]*EventNode{
				{{}, {}, {}},
				{{}, {}},
			},
		},
	}
	for i := 0; i < s.ec.TotalCount(); i++ {
		_, err := s.Read()
		assert.NoError(t, err)
	}
	_, err := s.Read()
	assert.Equal(t, io.EOF, err)
}
