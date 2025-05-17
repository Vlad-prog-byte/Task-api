package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-api/internals/models"
	"task-api/internals/storages"
	"testing"
)

func TestGetTasks(t *testing.T) {
	handler := Handler{Store: storages.NewMemoryStore()}
	handler.Store.Reset()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "[]" {
		t.Logf("Empty task list returned — OK")
	} else {
		t.Logf("Response body: %s", body)
	}
}

func TestAddTask(t *testing.T) {
	handler := Handler{Store: storages.NewMemoryStore()}
	handler.Store.Reset()
	var b bytes.Buffer
	task := models.CreateTask{
		Title: "check adding task",
		Description: "using httptest module built-in",
	}
	if err := json.NewEncoder(&b).Encode(task); err != nil {
    	t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tasks", &b)
	handler.AddTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w = httptest.NewRecorder()
	handler.GetTasks(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
	var tasks []models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("count tasks must be 1")
	}
	if tasks[0].Title != "check adding task" {
		t.Errorf("invalid title of the adding task")
	}
	if tasks[0].Description != "using httptest module built-in" {
		t.Errorf("invalid description of the adding task")
	}
}

func TestDeleteTask(t *testing.T) {
	handler := Handler{Store: storages.NewMemoryStore()}
	handler.Store.Reset()

	var b bytes.Buffer
	task := models.CreateTask{
		Title: "check adding task",
		Description: "using httptest module built-in",
	}
	if err := json.NewEncoder(&b).Encode(task); err != nil {
    	t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tasks", &b)
	handler.AddTask(w, req)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/task?id=0", nil)
	handler.DeleteTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	handler.GetTasks(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
	var tasks []models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks after deletion, got %d", len(tasks))
		t.Errorf("%d", tasks[0].ID)
	}
}

func TestPutTask(t *testing.T) {
	handler := Handler{Store: storages.NewMemoryStore()}
	handler.Store.Reset()

	var b bytes.Buffer
	task := models.CreateTask{
		Title: "check adding task",
		Description: "using httptest module built-in",
	}
	if err := json.NewEncoder(&b).Encode(task); err != nil {
    	t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tasks", &b)
	handler.AddTask(w, req)

	task_edit := models.EditTask{
		CreateTask: models.CreateTask{
			Title: "put",
			Description: "put",
		},
	}
	if err := json.NewEncoder(&b).Encode(task_edit); err != nil {
    	t.Fatal(err)
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/task?id=0", &b)
	handler.PutTask(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}


	req = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w = httptest.NewRecorder()
	handler.GetTasks(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
	var tasks []models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("count tasks must be 1")
	}
	if tasks[0].Title != "put" {
		t.Errorf("invalid title of the updating task")
	}
	if tasks[0].Description != "put" {
		t.Errorf("invalid description of the updating task")
	}
}