package utils

import (
	"encoding/json"
	"os"
)

func GetJSONFromFile(path string, target any) error {
	data, dataErr := os.ReadFile(path)
	if dataErr != nil {
		return dataErr
	}

	return json.Unmarshal(data, target)
}

func PutJSONToFile(path string, obj any) error {
	data, dataErr := json.MarshalIndent(obj, "", "    ")
	if dataErr != nil {
		return dataErr
	}

	return os.WriteFile(path, data, FMode)
}
