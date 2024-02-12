package main

import (
	"brenonaraujo/rinhabackend-q12024/api"
	"brenonaraujo/rinhabackend-q12024/infra/database"
	"log"
	"os"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if err := database.ConnectDB(databaseUrl); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.GetDBPool().Close()
	api.Run()
}
