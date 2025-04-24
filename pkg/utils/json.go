package utils

import (
	"encoding/json"
	"os"
)

func ToJsonFile(path string, v any) error {
	out, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, out, FMode)
	if err != nil {
		return err
	}
	return nil
}

func FromJsonFile(path string, v any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}
