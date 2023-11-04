package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sankalpmukim/todos-backend/internal/database"
	"github.com/sankalpmukim/todos-backend/pkg/logs"
)

func getUserId(r *http.Request) string {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		logs.Error("Error getting claims from context", err)
		return ""
	}
	return claims["sub"].(string)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID := getUserId(r)
	if userID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todos from database
	todos, err := database.DB.GetTodos(userID)

	if err != nil {
		logs.Error("Error getting todos from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return todos
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		logs.Error("Error writing JSON response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateTodoDone(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID := getUserId(r)
	if userID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todo id from url
	todoIDString := chi.URLParam(r, "id")
	if todoIDString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert todoID to int
	todoID, err := strconv.Atoi(todoIDString)
	if err != nil {
		logs.Error("Error converting todoID to int", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todo from database
	todo, err := database.DB.GetTodoByID(todoID)
	if err != nil {
		logs.Error("Error getting todo from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if todo belongs to user
	if todo.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// update todo
	toBeSet := !todo.Done
	logs.Info("Setting todo done to", toBeSet)
	if err := database.DB.SetTodoCompletedForUser(userID, todoID, toBeSet); err != nil {
		logs.Error("Error updating todo", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todo.Done = toBeSet
	todo.UpdatedAt = time.Now()
	// return updated todo
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		logs.Error("Error writing JSON response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID := getUserId(r)
	if userID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// decode request body
	var createTodo database.CreateTodo
	if err := json.NewDecoder(r.Body).Decode(&createTodo); err != nil {
		logs.Error("Error decoding request body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create todo
	todoId, err := database.DB.CreateTodoForUser(userID, createTodo)
	if err != nil {
		logs.Error("Error creating todo", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todo from database
	todo, err := database.DB.GetTodoByID(todoId)

	if err != nil {
		logs.Error("Error getting todo from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return created todo
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		logs.Error("Error writing JSON response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID := getUserId(r)
	if userID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todo id from url
	todoIDString := chi.URLParam(r, "id")
	if todoIDString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert todoID to int
	todoID, err := strconv.Atoi(todoIDString)
	if err != nil {
		logs.Error("Error converting todoID to int", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todo from database
	todo, err := database.DB.GetTodoByID(todoID)
	if err != nil {
		logs.Error("Error getting todo from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if todo belongs to user
	if todo.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// decode request body
	var todoUpdate database.Todo
	if err := json.NewDecoder(r.Body).Decode(&todoUpdate); err != nil {
		logs.Error("Error decoding request body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update todo
	if err := database.DB.UpdateTodoForUser(userID, todoID, todoUpdate); err != nil {
		logs.Error("Error updating todo", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return updated todo
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoUpdate); err != nil {
		logs.Error("Error writing JSON response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID := getUserId(r)
	if userID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todo id from url
	todoIDString := chi.URLParam(r, "id")
	if todoIDString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert todoID to int
	todoID, err := strconv.Atoi(todoIDString)
	if err != nil {
		logs.Error("Error converting todoID to int", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get todo from database
	todo, err := database.DB.GetTodoByID(todoID)
	if err != nil {
		logs.Error("Error getting todo from database", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check if todo belongs to user
	if todo.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// delete todo
	if err := database.DB.DeleteTodoForUser(userID, todoID); err != nil {
		logs.Error("Error deleting todo", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return deleted todo
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		logs.Error("Error writing JSON response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
