package service

import (
	"context"
	"taskmanager/internal/file"
	"taskmanager/internal/item"
	"taskmanager/internal/store"
)

type TaskService struct {
	MemoryStore *store.MemoryStore[item.Task]
	FileStorage *file.FileStorage[item.Task]
}

func NewTaskService() *TaskService {
	return &TaskService{
		MemoryStore: store.NewMemoryStore[item.Task](),
		FileStorage: file.NewFileStorage[item.Task](file.FILES_DIR + "tasks.json"),
	}
}

func (s *TaskService) LoadFromFile(ctx context.Context) error {
	data, err := s.FileStorage.Load()
	if err != nil {
		return err
	}

	taskMap := map[int]item.Task{}
	nextId := store.INITIAL_ID
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
	s.MemoryStore.Add(tasks, func(t *item.Task, id int) { t.Id = id })

	return s.SaveToFile(ctx)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskId int) error {
	s.MemoryStore.Delete(taskId)

	return s.SaveToFile(ctx)
}

func (s *TaskService) UpdateTasks(ctx context.Context, tasks []item.Task) error {
	s.MemoryStore.Update(tasks, func(t item.Task) int {
		return t.Id
	})

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
