package main

import (
	"fmt"
	"net/http"

	"taskmanager/handlers"
)

func main() {
	port := ":8080"

	fmt.Printf("Hello! Serving task manager at http://localhost%s !\n", port)

	http.HandleFunc("/", handlers.HealthHandler)
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListTasksHandler(w, r)
		case http.MethodPost:
			handlers.CreateTasksHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/task/", handlers.TaskByIdHandler)

	http.ListenAndServe(port, nil)
}
