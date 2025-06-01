package main

import (
	"fmt"
	"net/http"
	"task-api/internals/handlers"
	"task-api/internals/logger"
	"task-api/internals/services"
	"github.com/gorilla/mux"
)

func main () {
	service := services.NewTaskService()
	logger := logger.NewTaskLogger()
	handler := handlers.Handler{Service: service, Logger: logger}
	router := mux.NewRouter()
	router.HandleFunc("/task", handler.AddTask).Methods("POST")
	router.HandleFunc("/task", handler.PutTask).Methods("PUT")
	router.HandleFunc("/task", handler.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks", handler.GetTasks).Methods("GET")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
        fmt.Println("Server error:", err)
    }
	logger.Close()
}

