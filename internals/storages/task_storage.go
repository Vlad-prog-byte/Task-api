package storages

import (
	"errors"
	"sync"
	"task-api/internals/models"
)

type TaskStore interface {
    Add(title, desc string)
    GetAll() []models.Task
    Update(id int, title, description string, isDone bool) error
    Delete(id int) error
	Reset()
}

func NewMemoryStore() TaskStore {
	return &taskStore{
		tasks: make([]models.Task, 0),
	}
}

type taskStore struct {
	tasks []models.Task
	id int
	mu sync.Mutex
}

func (store *taskStore) Reset() {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.tasks = make([]models.Task, 0)
	store.id = 0
}

func (store *taskStore)newTask(title, desc string) models.Task {
	task := models.Task{ID: Store.id, Title: title, Description: desc}
	store.id++
	return task
}

var Store = &taskStore{tasks: make([]models.Task, 0)}

func (store *taskStore) Add(title, desc string) {
	task := store.newTask(title, desc)
	store.mu.Lock()
	defer store.mu.Unlock()
	store.tasks = append(store.tasks, task)
}

func (store *taskStore) GetAll() []models.Task {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.tasks
}

func (store *taskStore) Update(id int, title, description string, isDone bool) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	for i, _ := range store.tasks {
		if store.tasks[i].ID == id {
			store.tasks[i].Title = title
			store.tasks[i].Description = description
			store.tasks[i].IsDone = isDone
			return nil
		}
	}
	return errors.New("task with given id not found")
}

func (store *taskStore) Delete(id int) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	for i, task := range store.tasks {
		if task.ID == id {
			store.tasks = append(store.tasks[:i], store.tasks[i + 1:]...)
			return nil
		}
	}
	return errors.New("task with given id not found")
}