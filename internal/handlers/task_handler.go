package handlers

import (
	"context"
	"log"
	"net/http"
	"taskmanager/internal/service"
	"taskmanager/internal/utils"
)

type TaskHandler struct {
	TaskService *service.TaskService
	UserService *service.UserService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		TaskService: service.NewTaskService(),
		UserService: service.NewUserService(),
	}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.TaskService.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	if err := h.UserService.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := utils.ValidateListTasksRequest(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	_, err = h.UserService.GetUser(ctx, req.UserId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	tasks := h.TaskService.ListTasksByUserId(ctx, req.UserId)
	utils.RespondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) AddTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.TaskService.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	if err := h.UserService.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := utils.ValidateAddTasksRequest(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	_, err = h.UserService.GetUser(ctx, req.UserId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	for i := 0; i < len(req.Tasks); i++ {
		req.Tasks[i].UserId = req.UserId
	}

	err = h.TaskService.AddTasks(ctx, req.Tasks)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		log.Println(err.Error())

		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "added"})
}

func (h *TaskHandler) UpdateTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPatch {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.TaskService.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	if err := h.UserService.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := utils.ValidateUpdateTasksRequest(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	_, err = h.UserService.GetUser(ctx, req.UserId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	for i := 0; i < len(req.Tasks); i++ {
		req.Tasks[i].UserId = req.UserId
	}

	err = h.TaskService.UpdateTasks(ctx, req.Tasks)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "added"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.TaskService.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := utils.ValidateDeleteTaskRequest(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	err = h.TaskService.DeleteTask(ctx, req.TaskId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
