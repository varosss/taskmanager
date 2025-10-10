package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"taskmanager/item"
	"taskmanager/store"
)

type Task = item.Task

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
		http.Error(w, "User is invalid", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	taskStore := store.NewTaskStore(user)
	taskStore.LoadFromFile()

	statusStr := r.URL.Query().Get("status")

	if statusStr == "" {
		json.NewEncoder(w).Encode(taskStore.List())
		return
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil || status == int(item.StatusUnkown) {
		http.Error(w, "Unknown status", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(taskStore.ListByStatus(item.StatusFromInt(status)))
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
