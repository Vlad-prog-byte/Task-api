package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"task-api/internals/logger"
	"task-api/internals/models"
	"task-api/internals/services"
)

type Handler struct {
	Service services.TaskService
	Logger logger.TaskLogger
}

func SendSimpleResponse(w http.ResponseWriter, response map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler)AddTask(w http.ResponseWriter, req *http.Request) {
	var data models.CreateTask
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
        return
	}
	task, err := h.Service.Add(data.Title, data.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
        return
	}
	go h.Logger.Log(logger.LogEvent{Type: logger.LogCreate, Task: task})
	res, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
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
	task, err := h.Service.Update(id, data.Title, data.Description, data.IsDone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	go h.Logger.Log(logger.LogEvent{Type: logger.LogUpdate, Task: task})
	res, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func (h *Handler)DeleteTask(w http.ResponseWriter, req *http.Request) {
	id, err := parseIDFromQuery(req)
	if err != nil {
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	task, err := h.Service.Delete(id)
	if err != nil {
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	go h.Logger.Log(logger.LogEvent{Type: logger.LogDelete, Task: task})
	res, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h *Handler)GetTasks(w http.ResponseWriter, req *http.Request) {
	res, err := json.Marshal(h.Service.GetAll())
	if err != nil {
		SendSimpleResponse(w, map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
