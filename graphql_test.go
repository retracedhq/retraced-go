package retraced

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructuredQueryString(t *testing.T) {
	sq := &StructuredQuery{
		ActorID:  "actor1*",
		Location: "Los Angeles",
	}

	assert.Equal(t, "actor.id:actor1* location:\"Los Angeles\"", sq.String())
}
