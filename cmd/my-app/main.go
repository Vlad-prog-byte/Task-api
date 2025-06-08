package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-api/internals/handlers"
	"task-api/internals/logger"
	"task-api/internals/services"
	"time"
	"github.com/gorilla/mux"
)

func startArchiver(ctx context.Context, service services.TaskService, archiverLogger logger.TaskLogger) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
				case <- ctx.Done():
					defer ticker.Stop()
					// Я хз что тут делать
					return
				case <- ticker.C:
					tasks := service.GetAll()
					for i := range tasks {
						if tasks[i].IsDone {
							archiverLogger.Log(logger.LogEvent{Type: logger.LogArchiver, Task: tasks[i]})
						}
					}

			}
		}
	}()
}

func main () {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	service := services.NewTaskService()
	logger := logger.NewTaskLogger()
	defer logger.Close()
	startArchiver(ctx, service, logger)

	handler := handlers.Handler{Service: service, Logger: logger}
	router := mux.NewRouter()
	router.HandleFunc("/task", handler.AddTask).Methods("POST")
	router.HandleFunc("/task", handler.PutTask).Methods("PUT")
	router.HandleFunc("/task", handler.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks", handler.GetTasks).Methods("GET")
	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe error: %v\n", err)
			stop()
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	} else {
		fmt.Println("Server exited properly")
	}
}

