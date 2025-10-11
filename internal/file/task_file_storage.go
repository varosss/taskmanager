package file

import (
	"encoding/json"
	"os"
	"taskmanager/internal/item"
)

const FILES_DIR = "data"

type TaskFileStorage struct {
	Path string
}

func NewTaskFileStorage() *TaskFileStorage {
	return &TaskFileStorage{
		Path: "data/tasks.json",
	}
}

func (f *TaskFileStorage) Load() ([]item.Task, error) {
	_, err := os.Stat(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(FILES_DIR, 0755); err != nil {
				return nil, err
			}

			file, err := os.Create(f.Path)
			if err != nil {
				return nil, err
			}

			f.Save([]item.Task{})

			defer file.Close()

		} else {
			return nil, err
		}
	}

	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}

	var result []item.Task
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (f *TaskFileStorage) Save(data []item.Task) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.Path, bytes, 0644)
}
