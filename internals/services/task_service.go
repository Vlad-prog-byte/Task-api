package services

import (
	"errors"
	"task-api/internals/models"
	"task-api/internals/storages"
)

type TaskService interface {
    Add(title, description string) (models.Task, error)
    GetAll() []models.Task
    Update(id int, title, description string, isDone bool) (models.Task, error)
    Delete(id int) (models.Task, error)
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
	task := service.store.Add(title, description)
	return task, nil
}

func (service *taskService) GetAll() []models.Task {
	return service.store.GetAll()
}

func (service *taskService) Update(id int, title, description string, isDone bool) (models.Task, error) {
	return service.store.Update(id, title, description, isDone)
}

func (service *taskService) Delete(id int) (models.Task, error) {
	return service.store.Delete(id)
}

func (service *taskService) Reset() {
	service.store.Reset()
}