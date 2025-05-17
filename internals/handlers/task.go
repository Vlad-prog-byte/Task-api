package handlers

import (
	"task-api/internals/models"
	"task-api/internals/storages"
	"net/http"
	"encoding/json"
	"errors"
	"strconv"
)

type Handler struct {
	Store storages.TaskStore
}
// var (
// 	tasks = make([]models.Task, 0)
// 	nextId = 1
// 	mu sync.Mutex
// )

// func NewTask(title, desc string) models.Task {
// 	task := models.Task{ID: nextId, Title: title, Description: desc}
// 	nextId++
// 	return task
// }

func SendSimpleResponse(w http.ResponseWriter, response map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func AddTask(w http.ResponseWriter, req *http.Request) {
// 	var data models.CreateTask
// 	err := json.NewDecoder(req.Body).Decode(&data)
// 	if err != nil {
// 		SendSimpleResponse(w, map[string]string{"error": err.Error()})
//         return
// 	}
// 	task := NewTask(data.Title, data.Description)
// 	mu.Lock()
// 	defer mu.Unlock()
// 	tasks = append(tasks, task)
// 	SendSimpleResponse(w, map[string]string{"status": "Task added"})
// }

func (h *Handler)AddTask(w http.ResponseWriter, req *http.Request) {
	var data models.CreateTask
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
        return
	}
	h.Store.Add(data.Title, data.Description)
	SendSimpleResponse(w, map[string]string{"status": "Task added"})
}


func parseIDFromQuery(req *http.Request) (int, error) {
	idStr, ok := req.URL.Query()["id"]
	if !ok {
		return 0, errors.New("id parameter in " + req.Method + " method is required")
	}
	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return 0, err
	}
	return id, nil
}

// func PutTask(w http.ResponseWriter, req *http.Request) {
// 	id, err := parseIDFromQuery(req)
// 	if err != nil {
// 		SendSimpleResponse(w, map[string]string{"error": err.Error()})
// 		return
// 	}

// 	var data models.EditTask
// 	err = json.NewDecoder(req.Body).Decode(&data)
// 	if err != nil {
// 		SendSimpleResponse(w, map[string]string{"error": err.Error()})
//         return
// 	}
// 	mu.Lock()
// 	defer mu.Unlock()
// 	for ind, task := range tasks {
// 		if task.ID != id {
// 			continue
// 		}
// 		tasks[ind].Title = data.Title
// 		tasks[ind].Description = data.Description
// 		tasks[ind].IsDone = data.IsDone
// 		SendSimpleResponse(w, map[string]string{"status": "Task edited"})
// 		return
// 	}
// 	SendSimpleResponse(w, map[string]string{"error": "Task not found with given id"})
// }


func (h *Handler)PutTask(w http.ResponseWriter, req *http.Request) {
	id, err := parseIDFromQuery(req)
	if err != nil {
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}

	var data models.EditTask
	if err = json.NewDecoder(req.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
        return
	}
	if err = h.Store.Update(id, data.Title, data.Description, data.IsDone); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	SendSimpleResponse(w, map[string]string{"status": "Task edited"})

}


// func DeleteTask(w http.ResponseWriter, req *http.Request) {
// 	id, err := parseIDFromQuery(req)
// 	if err != nil {
// 		SendSimpleResponse(w, map[string]string{"error": err.Error()})
// 		return
// 	}

// 	mu.Lock()
// 	defer mu.Unlock()
// 	for ind, task := range tasks {
// 		if task.ID != id {
// 			continue
// 		}
// 		tasks = append(tasks[:ind], tasks[ind + 1:]...)
// 		SendSimpleResponse(w, map[string]string{"status": "Task deleted"})
// 		return
// 	}
// 	SendSimpleResponse(w, map[string]string{"error": "Task with given id does not exist"})
// }


func (h *Handler)DeleteTask(w http.ResponseWriter, req *http.Request) {
	id, err := parseIDFromQuery(req)
	if err != nil {
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	if err = h.Store.Delete(id); err != nil {
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	SendSimpleResponse(w, map[string]string{"status": "Task deleted"})
}


// func GetTasks(w http.ResponseWriter, req *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	res, err := json.Marshal(tasks)
// 	if err != nil {
// 		SendSimpleResponse(w, map[string]string{"error": err.Error()})
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(res)
// }

func (h *Handler)GetTasks(w http.ResponseWriter, req *http.Request) {
	res, err := json.Marshal(h.Store.GetAll())
	if err != nil {
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
