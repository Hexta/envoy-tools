package format

import yaml "github.com/goccy/go-yaml"

func YAML(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
