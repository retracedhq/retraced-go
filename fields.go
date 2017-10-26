package retraced

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Fields map[string]string

type fieldsList []struct{ Key, Value string }

// UnmarshalJSON handles [{key: "", value: ""},...] as returned by GraphQL.
func (fields *Fields) UnmarshalJSON(data []byte) error {
	list := make(fieldsList, 0)
	err := json.Unmarshal(data, &list)
	if err != nil {
		return fmt.Errorf("Fields.UnmarshalJSON: %v", err)
	}
	if len(list) > 0 && *fields == nil {
		*fields = make(Fields, len(list))
	}
	for _, f := range list {
		(*fields)[f.Key] = f.Value
	}

	return nil
}

// json without chance of error
func (fields Fields) String() string {
	var s []string
	for k, v := range fields {
		s = append(s, fmt.Sprintf("%q:%q", k, v))
	}

	return fmt.Sprintf("{%s}", strings.Join(s, ","))
}
