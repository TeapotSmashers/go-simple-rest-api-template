package routers

import (
	"github.com/go-chi/chi"
	"github.com/sankalpmukim/todos-backend/internal/handlers"
)

var Todos *chi.Mux

// GET /todos
// PUT /todos/id
// POST /todos
// POST/todos/id
// DELETE /todos/id

func init() {
	Todos = chi.NewRouter()
	Todos.Get("/", handlers.GetTodos)
	Todos.Put("/{id}", handlers.UpdateTodoDone)
	Todos.Post("/", handlers.CreateTodo)
	Todos.Post("/{id}", handlers.UpdateTodo)
	Todos.Delete("/{id}", handlers.DeleteTodo)
}
