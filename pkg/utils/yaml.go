package utils

import (
	"os"

	"github.com/goccy/go-yaml"
)

func GetYAMLFromFile(path string, target any) error {
	raw, rawErr := os.ReadFile(path)
	if rawErr != nil {
		return rawErr
	}
	yamlErr := yaml.Unmarshal(raw, target)
	return yamlErr
}

func PutYAMLToFile(path string, obj any) error {
	b, bErr := yaml.MarshalWithOptions(
		obj,
		yaml.Indent(4),
		yaml.IndentSequence(true),
	)

	if bErr != nil {
		return bErr
	}

	err := os.WriteFile(path, b, FMode)
	return err
}
