package retraced

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Event represents a single audited event.
// Required fields: Action, Group
type Event struct {
	// Action is a short, readable word to describe the action
	Action string `json:"action"`

	// Group is the team that will be able to see this event in the audit log
	Group *Group `json:"group,omitempty"`

	// Created is a timestamp representing when the event took place
	Created time.Time `json:"created"`

	// CRUD is a list of the most basic verbs that describe the type of action
	CRUD string `json:"crud"`

	// Target represents the item that had an action performed on it
	Target *Target `json:"target,omitempty"`

	// Description is a string with the full details of the event
	Description string `json:"description,omitempty"`

	// SourceIP is the client ip address where the event was performed
	SourceIP string `json:"source_ip,omitempty"`

	// Actor represents the entity that performed the action
	Actor *Actor `json:"actor,omitempty"`

	// Fields are any additional data to store with the event
	Fields Fields `json:"fields,omitempty"`

	// IsFailure is an optional flag that, when set, indicates that this audited event is a failed use of privileges
	IsFailure bool `json:"is_failure,omitempty"`

	// IsAnonymous is an optional flag that, when set, indicates that this is an anonymous event
	IsAnonymous bool `json:"is_anonymous,omitempty"`

	// Component is an identifier for a specific component of a vendor app platform
	// Component can be set on each Event, or on an instance of retraced.Client
	Component string `json:"component,omitempty"`

	// Version is an identifier for the specific version of this component, usually a git SHA
	// Version can be set on each Event, or on an instance of retraced.Client
	Version string `json:"version,omitempty"`

	// apiVersion is set here to allow updates to this model without breaking the API server
	apiVersion int

	ExternalID string `json:"external_id,omitempty"`

	Metadata Fields `json:"metadata,omitempty"`
}

// VerifyHash computes a hash of the sent event, and verifies
// that it matches the hash we got back from Retraced API
func (event *Event) VerifyHash(newEvent *NewEventRecord) error {
	// Basic sanity check
	if event.Action == "" {
		return fmt.Errorf("missing required field for hash verification: Action")
	}

	hashTarget := event.BuildHashTarget(newEvent)

	hashBytes := sha256.Sum256(hashTarget)
	result := hex.EncodeToString(hashBytes[:])

	if result != newEvent.Hash {
		return fmt.Errorf("hash mismatch: local[%s] != remote[%s]", result, newEvent.Hash)
	}

	return nil
}

// BuildHashTarget builds a string that will be used to
// compute a hash of the event
func (event *Event) BuildHashTarget(newEvent *NewEventRecord) []byte {

	concat := &bytes.Buffer{}
	fmt.Fprintf(concat, "%s:", encodePassOne(newEvent.ID))
	fmt.Fprintf(concat, "%s:", encodePassOne(event.Action))

	targetId := ""
	if event.Target != nil {
		targetId = event.Target.ID
	}
	fmt.Fprintf(concat, "%s:", encodePassOne(targetId))

	actorId := ""
	if event.Actor != nil {
		actorId = event.Actor.ID
	}
	fmt.Fprintf(concat, "%s:", encodePassOne(actorId))

	groupId := ""
	if event.Group != nil {
		groupId = event.Group.ID
	}
	fmt.Fprintf(concat, "%s:", encodePassOne(groupId))

	fmt.Fprintf(concat, "%s:", encodePassOne(event.SourceIP))

	if event.IsFailure {
		fmt.Fprint(concat, "1:")
	} else {
		fmt.Fprint(concat, "0:")
	}
	if event.IsAnonymous {
		fmt.Fprint(concat, "1:")
	} else {
		fmt.Fprint(concat, "0:")
	}

	if len(event.Fields) == 0 {
		fmt.Fprintf(concat, ":")
	} else {
		allKeys := []string{}
		for k := range event.Fields {
			allKeys = append(allKeys, k)
		}
		sort.Strings(allKeys)
		for i := 0; i < len(allKeys); i++ {
			k := allKeys[i]
			v := event.Fields[k]

			encodedKey := encodePassTwo(encodePassOne(k))
			encodedValue := encodePassTwo(encodePassOne(v))
			fmt.Fprintf(concat, "%s=%s;", encodedKey, encodedValue)
		}
	}

	if event.ExternalID != "" {
		fmt.Fprintf(concat, ":%s", encodePassOne(event.ExternalID))
	}

	if len(event.Metadata) > 0 {
		fmt.Fprintf(concat, ":")
		allKeys := []string{}
		for k := range event.Metadata {
			allKeys = append(allKeys, k)
		}
		sort.Strings(allKeys)
		for i := 0; i < len(allKeys); i++ {
			k := allKeys[i]
			v := event.Metadata[k]

			encodedKey := encodePassTwo(encodePassOne(k))
			encodedValue := encodePassTwo(encodePassOne(v))
			fmt.Fprintf(concat, "%s=%s;", encodedKey, encodedValue)
		}
	}

	return concat.Bytes()
}

func encodePassOne(in string) string {
	s := strings.Replace(in, "%", "%25", -1)
	return strings.Replace(s, ":", "%3A", -1)
}

func encodePassTwo(in string) string {
	s := strings.Replace(in, "=", "%3D", -1)
	return strings.Replace(s, ";", "%3B", -1)
}
