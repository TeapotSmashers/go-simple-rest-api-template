package database

import (
	"time"
)

type Todo struct {
	ID          int       `db:"id" json:"id"`                   // Matches SERIAL type
	UserID      string    `db:"user_id" json:"-"`               // Matches VARCHAR(255) and references users(id)
	Title       string    `db:"title" json:"title"`             // Matches VARCHAR(255)
	Description string    `db:"description" json:"description"` // Matches TEXT and can be NULL
	Done        bool      `db:"done" json:"done"`               // Matches BOOLEAN
	CreatedAt   time.Time `db:"created_at" json:"created_at"`   // Matches TIMESTAMPTZ
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`   // Matches TIMESTAMPTZ
}

type CreateTodo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (db *Database) GetTodos(userID string) ([]Todo, error) {
	var todos []Todo
	err := db.conn.Select(&todos, "SELECT * FROM todos WHERE user_id=$1", userID)
	return todos, err
}

func (db *Database) GetTodoByID(id int) (Todo, error) {
	var todo Todo
	err := db.conn.Get(&todo, "SELECT * FROM todos WHERE id=$1", id)
	return todo, err
}

// CreateTodoForUser inserts a new TODO for a user and returns the new ID
func (db *Database) CreateTodoForUser(userID string, todo CreateTodo) (Todo, error) {
	var newTodo Todo

	// Use QueryRow to execute the query and RETURNING clause to get all the fields
	err := db.conn.QueryRowx(`
		INSERT INTO todos (user_id, title, description) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, title, description, done, created_at, updated_at`,
		userID, todo.Title, todo.Description,
	).StructScan(&newTodo)

	if err != nil {
		// Handle the error; the insertion failed
		return Todo{}, err
	}

	// Return the new Todo struct
	return newTodo, nil
}

func (db *Database) SetTodoCompletedForUser(userID string, id int, done bool) error {
	_, err := db.conn.Exec("UPDATE todos SET done=$1 WHERE id=$2 AND user_id=$3", done, id, userID)
	return err
}

func (db *Database) UpdateTodoForUser(userID string, id int, todo UpdateTodo) (Todo, error) {
	var updatedTodo Todo

	// Use QueryRowx to execute the query, RETURNING clause to get all fields of the updated row
	err := db.conn.QueryRowx(`
		UPDATE todos 
		SET title=$1, description=$2, done=$3, updated_at=NOW() 
		WHERE id=$4 AND user_id=$5
		RETURNING id, user_id, title, description, done, created_at, updated_at`,
		todo.Title, todo.Description, todo.Done, id, userID,
	).StructScan(&updatedTodo)

	if err != nil {
		// Handle the error; the update failed
		return Todo{}, err
	}

	// Return the updated Todo struct
	return updatedTodo, nil
}

func (db *Database) DeleteTodoForUser(userID string, id int) error {
	_, err := db.conn.Exec("DELETE FROM todos WHERE id=$1 AND user_id=$2", id, userID)
	return err
}
