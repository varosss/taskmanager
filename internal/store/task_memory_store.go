package store

import (
	"sort"
	"taskmanager/internal/item"
)

type TaskMemoryStore struct {
	data   map[int]item.Task // map[UserID]map[ID]item.Task
	nextId int
}

func NewTaskMemoryStore() *TaskMemoryStore {
	return &TaskMemoryStore{
		data:   make(map[int]item.Task),
		nextId: 1,
	}
}

func (ms *TaskMemoryStore) SetData(data map[int]item.Task) {
	ms.data = data
}

func (ms *TaskMemoryStore) SetNextId(nextId int) {
	ms.nextId = nextId
}

func (ms *TaskMemoryStore) Add(tasks []item.Task) {
	for _, task := range tasks {
		task.Id = ms.nextId
		ms.data[ms.nextId] = task

		ms.nextId++
	}
}

func (ms *TaskMemoryStore) Delete(id int) {
	delete(ms.data, id)
}

func (ms *TaskMemoryStore) Update(tasks []item.Task) {
	for _, task := range tasks {
		ms.data[task.Id] = task
	}
}

func (ms *TaskMemoryStore) Get(id int) item.Task {
	return ms.data[id]
}

func (ms *TaskMemoryStore) Data() map[int]item.Task {
	return ms.data
}

func (ms *TaskMemoryStore) List() []item.Task {
	keys := []int{}
	for key := range ms.data {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	taskList := []item.Task{}
	for _, key := range keys {
		taskList = append(taskList, ms.data[key])
	}

	return taskList
}
