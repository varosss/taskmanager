package handlers

import (
	"context"
	"log"
	"net/http"
	"taskmanager/internal/service"
	"taskmanager/internal/utils"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{Service: service.NewUserService()}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.Service.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	users := h.Service.ListUsers(ctx)
	utils.RespondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) AddUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method is not allowed")

		return
	}

	if err := h.Service.LoadFromFile(context.Background()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't load data from file")

		log.Println(err.Error())

		return
	}

	req, err := utils.ValidateAddUsersRequest(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	err = h.Service.AddUsers(ctx, req.Users)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())

		log.Println(err.Error())

		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "added"})
}
