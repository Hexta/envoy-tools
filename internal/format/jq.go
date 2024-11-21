package format

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

const indentation = 2

func JQ(data interface{}, opts ...any) (string, error) {
	if len(opts) == 0 {
		return "", &InvalidFormatError{Format: "empty JQ query"}
	}
	query := opts[0].(string)

	jsonStr, err := JSON(data)
	if err != nil {
		return "", err
	}

	dataUnmarshalled := make(map[string]any)
	err = json.Unmarshal([]byte(jsonStr), &dataUnmarshalled)
	if err != nil {
		return "", &InvalidFormatError{Format: "invalid JSON"}
	}

	q, err := gojq.Parse(query)
	if err != nil {
		return "", err
	}

	iter := q.Run(dataUnmarshalled)

	sb := strings.Builder{}

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			return "", err
		}
		s, err := json.MarshalIndent(
			v,
			"",
			strings.Repeat(" ", indentation))

		if err != nil {
			return "", &InvalidFormatError{Format: "invalid JSON"}
		}

		sb.WriteString(fmt.Sprintf("%s\n", s))
	}

	return sb.String(), nil
}
