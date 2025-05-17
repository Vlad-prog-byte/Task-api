package models


type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    IsDone      bool      `json:"is_done"`
}

type CreateTask struct {
	Title	string `json:"title"`
	Description string `json:"description"`
}

type EditTask struct {
    CreateTask
    IsDone bool `json:"is_done"`
}
