package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"taskmanager/item"
	"taskmanager/store"
)

type Task = item.Task
type Status = item.Status

var taskByIdRoute = regexp.MustCompile(`^/task/(\d+)$`)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse := map[string]string{"status": "alive"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

func TaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	matches := taskByIdRoute.FindStringSubmatch(r.URL.Path)
	if len(matches) != 2 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "User is required", http.StatusBadRequest)
		return
	}

	taskStore := store.NewTaskStore(user)
	taskStore.LoadFromFile()

	switch r.Method {
	case http.MethodPut:
		var task Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		taskStore.Update(task)
		taskStore.SaveToFile()

		json.NewEncoder(w).Encode(map[string]string{"result": "success"})

	case http.MethodDelete:
		taskStore.Delete(id)
		taskStore.SaveToFile()

		json.NewEncoder(w).Encode(map[string]string{"result": "success"})

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "User is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	taskStore := store.NewTaskStore(user)
	taskStore.LoadFromFile()

	statusQueries := r.URL.Query()["status"]
	statuses := make(map[Status]bool, len(statusQueries))
	for _, s := range statusQueries {
		if s == "" {
			continue
		}
		sInt, err := strconv.Atoi(s)
		if err != nil {
			http.Error(w, "Invalid status value", http.StatusBadRequest)
			return
		}
		statuses[item.StatusFromInt(sInt)] = true
	}

	title := r.URL.Query().Get("title")

	var filtered []item.Task
	allTasks := taskStore.List()

	for _, task := range allTasks {
		if len(statuses) > 0 && !statuses[task.Status] {
			continue
		}
		if title != "" && !strings.Contains(task.Title, title) {
			continue
		}
		filtered = append(filtered, task)
	}

	json.NewEncoder(w).Encode(filtered)
}

func CreateTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	user := r.URL.Query().Get("user")

	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	taskStore := store.NewTaskStore(user)
	taskStore.LoadFromFile()

	for _, task := range tasks {
		taskStore.Add(task)
	}

	taskStore.SaveToFile()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}
