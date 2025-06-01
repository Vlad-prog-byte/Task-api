package services

import (
	"errors"
	"task-api/internals/models"
	"task-api/internals/storages"
)

type TaskService interface {
    Add(title, description string) (models.Task, error)
    GetAll() []models.Task
    Update(id int, title, description string, isDone bool) error
    Delete(id int) error
	Reset()
}

type taskService struct {
	store storages.TaskStore
}

func NewTaskService() TaskService {
	return &taskService{store: storages.NewMemoryStore()}
}

func (service *taskService) Add(title, description string) (models.Task, error){
	if title == "" {
		return models.Task{}, errors.New("title is required")
	}
	return service.store.Add(title, description), nil
}

func (service *taskService) GetAll() []models.Task {
	return service.store.GetAll()
}

func (service *taskService) Update(id int, title, description string, isDone bool) error {
	return service.store.Update(id, title, description, isDone)
}

func (service *taskService) Delete(id int) error {
	return service.store.Delete(id)
}

func (service *taskService) Reset() {
	service.store.Reset()
}