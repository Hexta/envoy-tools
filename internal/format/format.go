package format

import (
	"fmt"
	"strings"
)

type UnknownFormatError struct {
	Format string
}

func (e *UnknownFormatError) Error() string {
	return fmt.Sprintf("unknown format: %s", e.Format)
}

type InvalidFormatError struct {
	Format string
}

func (e *InvalidFormatError) Error() string {
	return fmt.Sprintf("invalid format: %s", e.Format)
}

type UnknownInputTypeError struct {
	Type string
}

func (e *UnknownInputTypeError) Error() string {
	return fmt.Sprintf("unknown input type: %s", e.Type)
}

type Options struct {
	Indent    int
	StatsOnly bool
}

type Format struct {
	Name       string
	Expression string
}

const (
	JSONFormat = "json"
	YAMLFormat = "yaml"
	TextFormat = "text"
	JQFormat   = "jq"
)

func (f *Format) String() string {
	return f.Name
}

func (f *Format) Set(s string) error {
	tokens := strings.SplitN(s, "=", 2)
	if len(tokens) == 0 {
		return &InvalidFormatError{Format: s}
	}

	name := tokens[0]
	expression := ""

	if len(tokens) == 2 {
		expression = tokens[1]
	}

	switch name {
	case JSONFormat, YAMLFormat, TextFormat, JQFormat:
		*f = Format{Name: name, Expression: expression}
		return nil
	default:
		return &UnknownFormatError{Format: name}
	}
}

func (f *Format) Type() string {
	return "Format"
}

func Apply(format Format, input any, opts Options) (string, error) {
	switch format.Name {
	case JSONFormat:
		return JSON(input)
	case YAMLFormat:
		return YAML(input)
	case TextFormat:
		return Text(input, opts)
	case JQFormat:
		return JQ(input, format.Expression)
	default:
		return "", &UnknownFormatError{Format: format.Name}
	}
}
