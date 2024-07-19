package utils

import (
	"map-reduce/config"
	"os"
	"path/filepath"
)

func ReadInput(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteOutput(filePath string, data string) error {
	dir := filepath.Dir(filePath)
	err := os.Mkdir(dir, os.ModePerm) // 0777 for directories
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(data), 0644) // 0644 for file
}

func SplitInput(content string) []string {
	chunkSize := config.CHUNK_SIZE
	var chunks []string
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		end = min(len(content), end)
		chunks = append(chunks, content[i:end])
	}
	return chunks
}
