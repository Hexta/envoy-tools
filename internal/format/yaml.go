package format

import "github.com/goccy/go-yaml"

func YAML(data interface{}) (string, error) {
	bytes, err := yaml.MarshalWithOptions(data, yaml.UseLiteralStyleIfMultiline(true))
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
