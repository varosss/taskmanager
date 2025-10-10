package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"taskmanager/item"
)

type Task = item.Task
type Status = item.Status

type TaskStore struct {
	LastTaskId int
	User       string
	Tasks      map[int]Task
}

func NewTaskStore(user string) *TaskStore {
	return &TaskStore{
		LastTaskId: 0,
		User:       user,
		Tasks:      make(map[int]Task, 0),
	}
}

func (ts *TaskStore) Add(task Task) int {
	ts.LastTaskId++

	task.Id = ts.LastTaskId

	if task.Status == 0 {
		task.Status = item.StatusInQueue
	}

	ts.Tasks[ts.LastTaskId] = task

	return task.Id
}

func (ts *TaskStore) Update(task Task) int {
	ts.Tasks[task.Id] = task

	return task.Id
}

func (ts *TaskStore) Get(id int) Task {
	return ts.Tasks[id]
}

func (ts *TaskStore) List(statuses ...Status) []Task {
	keys := make([]int, 0, len(ts.Tasks))
	for k := range ts.Tasks {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	actualStatuses := make(map[Status]bool, len(statuses))
	for _, status := range statuses {
		actualStatuses[status] = true
	}

	tasksList := []Task{}
	for _, k := range keys {
		if len(statuses) > 0 && actualStatuses[ts.Tasks[k].Status] {
			tasksList = append(tasksList, ts.Tasks[k])
		} else if len(statuses) == 0 {
			tasksList = append(tasksList, ts.Tasks[k])
		}
	}

	return tasksList
}

func (ts *TaskStore) ListByStatus(status Status) []Task {
	tasksList := []Task{}
	for _, task := range ts.List() {
		if task.Status == status {
			tasksList = append(tasksList, task)
		}
	}

	return tasksList
}

// func (ts *TaskStore) List() {
// 	fmt.Println("Список задач:")

// 	for _, task := range ts.All() {
// 		fmt.Printf("%d. -- %s -- %s\n", task.Id, task.Title, task.Status.String())
// 	}

// 	fmt.Print("\n")
// }

func (ts *TaskStore) Delete(id int) {
	delete(ts.Tasks, id)
}

func (ts *TaskStore) ListByTitle(title string) []Task {
	tasksList := []Task{}
	for _, task := range ts.Tasks {
		if task.Title == title {
			tasksList = append(tasksList, task)
		}
	}

	return tasksList
}

func (ts *TaskStore) GroupByCategoryMap() map[string][]Task {
	grouped := make(map[string][]Task)
	for _, task := range ts.List() {
		grouped[task.Category] = append(grouped[task.Category], task)
	}
	return grouped
}

func (ts *TaskStore) LoadFromFile() error {
	filename := fmt.Sprintf("upload/%s.json", ts.User)

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var loaded []Task
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}

	ts.Tasks = make(map[int]Task, len(loaded))
	for _, task := range loaded {
		if ts.LastTaskId < task.Id {
			ts.LastTaskId = task.Id
		}

		ts.Tasks[task.Id] = task
	}

	return nil
}

func (ts *TaskStore) SaveToFile() error {
	filename := fmt.Sprintf("upload/%s.json", ts.User)

	data, err := json.MarshalIndent(ts.List(), "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
