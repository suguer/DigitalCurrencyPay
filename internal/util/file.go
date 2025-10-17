package util

import "os"

func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
