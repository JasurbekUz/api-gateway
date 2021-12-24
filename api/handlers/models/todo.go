package models

type CreateTodo struct {
	Assignee string `json:"assignee"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Deadline string `json:"deadline"`
	Status   string `json:"status"`
}

type GetTodo struct {
	Id       string `json:"id"`
	Assignee string `json:"assignee"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Deadline string `json:"deadline"`
	Status   string `json:"status"`
}

type GetTodos struct {
	Todos []GetTodo `json:"todos"`
	Count int64     `json:"count"`
}
