package file

import (
	"encoding/json"
	"os"
)

const FILES_DIR = "data/"

type FileStorage[T any] struct {
	Path string
}

func NewFileStorage[T any](path string) *FileStorage[T] {
	return &FileStorage[T]{Path: path}
}

func (f *FileStorage[T]) Load() ([]T, error) {
	_, err := os.Stat(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(FILES_DIR, 0755); err != nil {
				return nil, err
			}

			empty := []T{}
			if err := f.Save(empty); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}

	var result []T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (f *FileStorage[T]) Save(data []T) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.Path, bytes, 0644)
}
