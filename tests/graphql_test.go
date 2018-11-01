package tests

import (
	"os"
	"testing"
	"time"

	retraced "github.com/retracedhq/retraced-go"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func fakeEvent() *retraced.Event {
	return &retraced.Event{
		Action: "go.test",
		Group: &retraced.Group{
			ID:   uuid.NewV4().String(),
			Name: "domain.xyz",
		},
		Created: time.Now(),
		CRUD:    "c",
		Target: &retraced.Target{
			ID:     uuid.NewV4().String(),
			Name:   "API",
			Type:   "server",
			Href:   "https://api.retraced.io",
			Fields: retraced.Fields{"x": "t1", "y": "t2", "z": "t3"},
		},
		Description: "Fake fake fake",
		SourceIP:    "1.1.1.100",
		Actor: &retraced.Actor{
			ID:     uuid.NewV4().String(),
			Name:   "user1@domain.xyz",
			Href:   "https://domain.xyz/users/user1",
			Fields: retraced.Fields{"x": "a1", "y": "a2", "z": "a3"},
		},
		Fields:      retraced.Fields{"x": "e1", "y": "e2", "z": "e3"},
		IsFailure:   false,
		IsAnonymous: false,
		Component:   "Go SDK",
		Version:     "v1",
	}
}

func TestClientQuery(t *testing.T) {
	projectID := os.Getenv("PROJECT_ID")
	token := os.Getenv("PUBLISHER_API_KEY")
	apiEndpoint := os.Getenv("PUBLISHER_API_ENDPOINT")

	client, err := retraced.NewClient(projectID, token)
	if err != nil {
		t.Fatal(err)
	}
	client.Endpoint = apiEndpoint

	uniqueActorID := uuid.NewV4().String()
	for i := 0; i < 10; i++ {
		e := fakeEvent()
		e.Actor.ID = uniqueActorID
		_, err := client.ReportEvent(e)
		if err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second * 2)

	sq := &retraced.StructuredQuery{
		ActorID: uniqueActorID,
	}
	mask := &retraced.EventNodeMask{
		ID:          true,
		Action:      true,
		GroupID:     true,
		CRUD:        true,
		Fields:      true,
		ActorFields: true,
	}
	eventsConn, err := client.Query(sq, mask, 3)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, eventsConn.CurrentResults(), 3)
	assert.Equal(t, 10, eventsConn.TotalCount())
	assert.True(t, eventsConn.HasNextPage())
	assert.False(t, eventsConn.HasPreviousPage())
	assert.Equal(t, 1, eventsConn.CurrentPageNumber())
	assert.Equal(t, 4, eventsConn.TotalPages())

	for _, node := range eventsConn.CurrentResults() {
		assert.NotEmpty(t, node.ID)
		assert.Equal(t, "go.test", node.Action)
		assert.NotEmpty(t, node.Group.ID)
		assert.Empty(t, node.Group.Name)
		assert.Equal(t, "c", node.CRUD)
		assert.Equal(t, "e1", node.Fields["x"])
		assert.Equal(t, "e2", node.Fields["y"])
		assert.Equal(t, "a3", node.Actor.Fields["z"])

		assert.Nil(t, node.Target)
		assert.Nil(t, node.Display)
		assert.Empty(t, node.Actor.ID)
		assert.Empty(t, node.Actor.Name)
		assert.Empty(t, node.Actor.Href)
		assert.Empty(t, node.Group.Name)
		assert.Empty(t, node.Description)
		assert.Zero(t, node.Created)
		assert.Zero(t, node.Received)
		assert.Zero(t, node.CanonicalTime)
		assert.Empty(t, node.SourceIP)
		assert.Empty(t, node.Component)
		assert.Empty(t, node.Version)
		assert.Empty(t, node.Country)
		assert.Empty(t, node.LocSubdiv1)
		assert.Empty(t, node.LocSubdiv2)
	}

	assert.NoError(t, eventsConn.NextPage())
	assert.NoError(t, eventsConn.NextPage())
	assert.NoError(t, eventsConn.NextPage())

	assert.Equal(t, 10, eventsConn.TotalCount())
	assert.Equal(t, 4, eventsConn.CurrentPageNumber())
	assert.Equal(t, 4, eventsConn.TotalPages())
	assert.True(t, eventsConn.HasPreviousPage())
	assert.False(t, eventsConn.HasNextPage())
}
