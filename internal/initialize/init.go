package initialize

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sankalpmukim/todos-backend/internal/auth"
	"github.com/sankalpmukim/todos-backend/internal/database"
	"github.com/sankalpmukim/todos-backend/pkg/logs"
)

func InitAll() error {
	if err := initializeEnv(); err != nil {
		fmt.Println("Error initializing env")
		return err
	}
	if err := logs.Initialize(); err != nil {
		fmt.Println("Error initializing logs")
		return err
	}
	if err := database.Initialize(); err != nil {
		fmt.Println("Error initializing database")
		return err
	}
	auth.Initialize()
	return nil
}

func initializeEnv() error {
	err := godotenv.Load()
	return err
}
