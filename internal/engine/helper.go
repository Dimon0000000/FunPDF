package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ConcatenatingStrings(str string) string {
	parts := strings.Split(str, "\n")
	return strings.Join(parts, "")
}

func loadTranslatorConfig(name string) (json.RawMessage, error) {
	data, err := os.ReadFile(fmt.Sprintf("conf/%s.json", name))
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return data, nil
}
