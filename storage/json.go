package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Storage interface {
	Save(key string, data interface{}) error
	Load(key string, target interface{}) error
	Exists(key string) bool
}

type FileStorage struct {
	BasePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &FileStorage{
		BasePath: basePath,
	}, nil
}

func (fs *FileStorage) Save(key string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	filePath := filepath.Join(fs.BasePath, key+".json")
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (fs *FileStorage) Load(key string, target interface{}) error {
	filePath := filepath.Join(fs.BasePath, key+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

func (fs *FileStorage) Exists(key string) bool {
	filePath := filepath.Join(fs.BasePath, key+".json")
	_, err := os.Stat(filePath)
	return err == nil
}

func DefaultStorage(appName string) (*FileStorage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	storagePath := filepath.Join(homeDir, ".config", appName)

	return NewFileStorage(storagePath)
}
