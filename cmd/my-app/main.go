package main

import (
	"fmt"
	"net/http"
	"task-api/internals/handlers"
	"task-api/internals/storages"

	"github.com/gorilla/mux"
)

func main () {
	store := storages.NewMemoryStore()
	handler := handlers.Handler{Store: store}
	router := mux.NewRouter()
	fmt.Println("Server is starting...") // ← добавь это для проверки
	router.HandleFunc("/task", handler.AddTask).Methods("POST")
	router.HandleFunc("/task", handler.PutTask).Methods("PUT")
	router.HandleFunc("/task", handler.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks", handler.GetTasks).Methods("GET")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
        fmt.Println("Server error:", err)
    }
}

