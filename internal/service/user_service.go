package service

import (
	"context"
	"fmt"
	"taskmanager/internal/file"
	"taskmanager/internal/item"
	"taskmanager/internal/store"
)

type UserService struct {
	MemoryStore *store.MemoryStore[item.User]
	FileStorage *file.FileStorage[item.User]
}

func NewUserService() *UserService {
	return &UserService{
		MemoryStore: store.NewMemoryStore[item.User](),
		FileStorage: file.NewFileStorage[item.User](file.FILES_DIR + "users.json"),
	}
}

func (s *UserService) LoadFromFile(ctx context.Context) error {
	data, err := s.FileStorage.Load()
	if err != nil {
		return err
	}

	userMap := map[int]item.User{}
	nextId := store.INITIAL_ID
	for _, task := range data {
		if nextId < task.Id {
			nextId = task.Id
		}

		userMap[task.Id] = task
	}

	s.MemoryStore.SetData(userMap)
	s.MemoryStore.SetNextId(nextId)

	return nil
}

func (s *UserService) SaveToFile(ctx context.Context) error {
	return s.FileStorage.Save(s.MemoryStore.List())
}

func (s *UserService) AddUsers(ctx context.Context, users []item.User) error {
	s.MemoryStore.Add(users, func(t *item.User, id int) { t.Id = id })

	return s.SaveToFile(ctx)
}

func (s *UserService) DeleteUser(ctx context.Context, userId int) error {
	s.MemoryStore.Delete(userId)

	return s.SaveToFile(ctx)
}

func (s *UserService) UpdateUsers(ctx context.Context, users []item.User) error {
	s.MemoryStore.Update(users, func(t item.User) int {
		return t.Id
	})

	return s.SaveToFile(ctx)
}

func (s *UserService) ListUsers(ctx context.Context) []item.User {
	return s.MemoryStore.List()
}

func (s *UserService) GetUser(ctx context.Context, userId int) (*item.User, error) {
	user, ok := s.MemoryStore.Get(userId)
	if !ok {
		return nil, fmt.Errorf("couldn't get user by id %d", userId)
	}

	return &user, nil
}
