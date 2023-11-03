package database

import "database/sql"

// User represents a user in the system.
type User struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Provider string `db:"provider" json:"provider"`
}

func (db *Database) GetUserByEmail(email string) (User, error) {
	var user User
	err := db.conn.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	return user, err
}

func (db *Database) UserExists(email string) (bool, error) {
	var user User
	err := db.conn.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			// The error is a "no rows in result set" error,
			// which means the user does not exist
			return false, nil
		}
		// An error occurred during the query execution
		return false, err
	}
	// No error occurred and a user was found
	return true, nil
}

func (db *Database) CreateUser(user User) error {
	_, err := db.conn.Exec("INSERT INTO users (id, name, email, provider) VALUES ($1, $2, $3, $4)", user.ID, user.Name, user.Email, user.Provider)
	return err
}
