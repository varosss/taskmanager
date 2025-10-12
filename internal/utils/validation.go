package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"taskmanager/internal/item"
)

type ListTasksRequest struct {
	UserId int
}

type AddTasksRequest struct {
	UserId int
	Tasks  []item.Task
}

type UpdateTasksRequest struct {
	UserId int
	Tasks  []item.Task
}

type DeleteTaskRequest struct {
	TaskId int
}

type AddUsersRequest struct {
	Users []item.User
}

func ValidateListTasksRequest(r *http.Request) (*ListTasksRequest, error) {
	strUserId := r.URL.Query().Get("user_id")
	if strUserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return nil, fmt.Errorf("user_id must be an integer")
	}

	return &ListTasksRequest{
		UserId: userId,
	}, nil
}

func ValidateAddTasksRequest(r *http.Request) (*AddTasksRequest, error) {
	strUserId := r.URL.Query().Get("user_id")
	if strUserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return nil, fmt.Errorf("user_id must be an integer")
	}

	var tasks []item.Task
	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("at least one task is required")
	}

	return &AddTasksRequest{
		UserId: userId,
		Tasks:  tasks,
	}, nil
}

func ValidateUpdateTasksRequest(r *http.Request) (*UpdateTasksRequest, error) {
	strUserId := r.URL.Query().Get("user_id")
	if strUserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		return nil, fmt.Errorf("user_id must be an integer")
	}

	var tasks []item.Task
	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("at least one task is required")
	}

	return &UpdateTasksRequest{
		UserId: userId,
		Tasks:  tasks,
	}, nil
}

func ValidateDeleteTaskRequest(r *http.Request) (*DeleteTaskRequest, error) {
	strTaskId := r.URL.Query().Get("task_id")
	if strTaskId == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	taskId, err := strconv.Atoi(strTaskId)
	if err != nil {
		return nil, fmt.Errorf("task_id must be an integer")
	}

	return &DeleteTaskRequest{
		TaskId: taskId,
	}, nil
}

func ValidateAddUsersRequest(r *http.Request) (*AddUsersRequest, error) {
	var users []item.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("at least one user is required")
	}

	return &AddUsersRequest{
		Users: users,
	}, nil
}
