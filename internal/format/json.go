package format

import (
	"encoding/json"
	"strings"
)

func JSON(data any) (string, error) {
	bytes, err := json.MarshalIndent(data, "", strings.Repeat(" ", 2))
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
