package services

import (
	"task-api/internals/models"
	"task-api/internals/storages"
)

type TaskService interface {
    Add(title, description string)
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

func (service *taskService) Add(title, description string) {
	service.store.Add(title, description)
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