package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"x-clone/server/constants"
	"x-clone/server/startup"
)

func main() {
	startServer()
}

func startServer() {
	db := startup.StartDatabase()

	if db == nil {
		fmt.Println("Error starting the database")
		return
	}

	s, serverError := db.DB()
	if serverError != nil {
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(s)

	if startRoutesErr := startup.StartRoutes(db); startRoutesErr != nil {
		return
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == constants.Empty {
		log.Panic("serverPort environment variable is not set")
	}

	fmt.Printf("Server running on port %s", serverPort)
	serverError = http.ListenAndServe(":"+serverPort, nil)
	if serverError != nil {
		return
	}
}
