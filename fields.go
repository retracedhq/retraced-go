package retraced

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Fields map[string]string

type fieldsList []struct{ Key, Value string }

// UnmarshalJSON handles
// [{key: "", value: ""},...] as returned by GraphQL.
// {"key": "value", ...} as returned by json.Marshal when data is marshalled in Go.
func (fields *Fields) UnmarshalJSON(data []byte) error {
	list := make(fieldsList, 0)
	err := json.Unmarshal(data, &list)
	if err != nil {
		f := make(map[string]interface{})
		if err := json.Unmarshal(data, &f); err != nil {
			return fmt.Errorf("Fields.UnmarshalJSON: %v", err)
		}

		*fields = make(Fields, len(list))
		for key, val := range f {
			(*fields)[key] = fmt.Sprintf("%v", val)
		}
		return nil
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
