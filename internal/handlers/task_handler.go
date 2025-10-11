package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"taskmanager/internal/service"
)

type TaskHandler struct {
	Service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{Service: s}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.Service.LoadFromFile(context.Background()); err != nil {
		respondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := validateListTasksRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())

		return
	}

	tasks := h.Service.ListTasksByUserId(ctx, req.UserId)
	respondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) AddTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.Service.LoadFromFile(context.Background()); err != nil {
		respondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := validateAddTasksRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())

		return
	}

	for i := 0; i < len(req.Tasks); i++ {
		req.Tasks[i].UserId = req.UserId
	}

	err = h.Service.AddTasks(ctx, req.Tasks)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())

		log.Println(err.Error())

		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "added"})
}

func (h *TaskHandler) UpdateTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.Service.LoadFromFile(context.Background()); err != nil {
		respondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := validateUpdateTasksRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())

		return
	}

	for _, task := range req.Tasks {
		task.UserId = req.UserId
	}

	err = h.Service.UpdateTasks(ctx, req.Tasks)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())

		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "added"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.Service.LoadFromFile(context.Background()); err != nil {
		respondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := validateDeleteTaskRequest(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())

		return
	}

	err = h.Service.DeleteTask(ctx, req.TaskId)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())

		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// Вспомогательные функции
func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
