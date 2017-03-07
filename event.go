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
	Group *Group `json:"group"`

	// Created is a timestamp representing when the event took place
	Created time.Time `json:"created"`

	// CRUD is a list of the most basic verbs that describe the type of action
	CRUD string `json:"crud"`

	// Target represents the item that had an action performed on it
	Target *Target `json:"target,omitempty"`

	// Description is a string with the full details of the event
	Description string `json:"description"`

	// SourceIP is the client ip address where the event was performed
	SourceIP string `json:"source_ip"`

	// Actor represents the entity that performed the action
	Actor *Actor `json:"actor"`

	// Fields are any additional data to store with the event
	Fields map[string]string `json:"fields,omitempty"`

	// IsFailure is an optional flag that, when set, indicates that this audited event is a failed use of privileges
	IsFailure bool `json:"is_failure"`

	// IsAnonymous is an optional flag that, when set, indicates that this is an anonymous event
	IsAnonymous bool `json:"is_anonymous"`

	// apiVersion is set here to allow updates to this model without breaking the API server
	apiVersion int
}

func (event *Event) VerifyHash(newEvent *NewEventRecord) error {
	// Basic sanity check
	if event.Action == "" {
		return fmt.Errorf("Missing required field for hash verification: Action")
	}
	if event.Group == nil || event.Group.ID == "" {
		return fmt.Errorf("Missing required field for hash verification: Group.Id")
	}

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

	fmt.Fprintf(concat, "%s:", encodePassOne(event.Group.ID))
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

	hashBytes := sha256.Sum256(concat.Bytes())
	result := hex.EncodeToString(hashBytes[:])

	if result != newEvent.Hash {
		return fmt.Errorf("Hash mismatch: local[%s] != remote[%s]", result, newEvent.Hash)
	}

	return nil
}

func encodePassOne(in string) string {
	s := strings.Replace(in, "%", "%25", -1)
	return strings.Replace(s, ":", "%3A", -1)
}

func encodePassTwo(in string) string {
	s := strings.Replace(in, "=", "%3D", -1)
	return strings.Replace(s, ";", "%3B", -1)
}
