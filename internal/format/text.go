package format

import (
	"reflect"

	"github.com/Hexta/envoy-tools/internal/diff"
)

func Text(input any, opts ...any) (string, error) {
	if input == nil {
		return "", nil
	}

	switch inputTyped := input.(type) {
	case *diff.Changes:
		return ChangesAsText(inputTyped, opts[0].(Options)), nil

	default:
		i := reflect.TypeOf(input)
		return "", &UnknownInputTypeError{Type: i.String()}
	}
}
