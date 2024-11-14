package shared

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJSONFile[T any](sourcePath string) (obj T, err error) {
	bytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return obj, fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
	}
	if err = json.Unmarshal(bytes, &obj); err != nil {
		return obj, fmt.Errorf("failed to parse source file %s: %w", sourcePath, err)
	}
	return obj, nil
}

func WriteJSONFile(path string, object map[string]any) error {
	bytes, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON file %s: %w", path, err)
	}
	bytes = append(bytes, '\n')
	if err = os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file %s: %w", path, err)
	}
	return nil
}

func OmitEmpty(v any) (isEmpty bool) {
	switch v := v.(type) {
	case string:
		return v == ""
	case bool:
		return !v
	case float64:
		return v == 0
	case []any:
		for i := range v {
			OmitEmpty(v[i])
		}
		return len(v) == 0
	case map[string]any:
		for key, value := range v {
			if isKeyEmpty := OmitEmpty(value); isKeyEmpty {
				delete(v, key)
			}
		}
		return len(v) == 0
	default:
		return v == nil
	}
}
