package service

import (
	"context"
	"taskmanager/internal/file"
	"taskmanager/internal/item"
	"taskmanager/internal/store"
)

type TaskService struct {
	MemoryStore *store.TaskMemoryStore
	FileStorage *file.TaskFileStorage
}

func NewTaskService() *TaskService {
	return &TaskService{
		MemoryStore: store.NewTaskMemoryStore(),
		FileStorage: file.NewTaskFileStorage(),
	}
}

func (s *TaskService) LoadFromFile(ctx context.Context) error {
	data, err := s.FileStorage.Load()
	if err != nil {
		return err
	}

	taskMap := map[int]item.Task{}
	nextId := 0
	for _, task := range data {
		if nextId < task.Id {
			nextId = task.Id
		}

		taskMap[task.Id] = task
	}

	s.MemoryStore.SetData(taskMap)
	s.MemoryStore.SetNextId(nextId)

	return nil
}

func (s *TaskService) SaveToFile(ctx context.Context) error {
	return s.FileStorage.Save(s.MemoryStore.List())
}

func (s *TaskService) AddTasks(ctx context.Context, tasks []item.Task) error {
	s.MemoryStore.Add(tasks)

	return s.SaveToFile(ctx)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskId int) error {
	s.MemoryStore.Delete(taskId)

	return s.SaveToFile(ctx)
}

func (s *TaskService) UpdateTasks(ctx context.Context, tasks []item.Task) error {
	s.MemoryStore.Update(tasks)

	return s.SaveToFile(ctx)
}

func (s *TaskService) ListTasksByUserId(ctx context.Context, userId int) []item.Task {
	allTasks := s.MemoryStore.List()
	filtered := []item.Task{}

	for _, task := range allTasks {
		if task.UserId == userId {
			filtered = append(filtered, task)
		}
	}

	return filtered
}
