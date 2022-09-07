package main

import (
	"database/sql"
	"github.com/ChrisCrawford1/Command/internal/handlers"
	"github.com/ChrisCrawford1/Command/internal/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	loadEnvironment()
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal(err)
		return
	}

	requestHandler := &handlers.RequestHandler{
		Commands: models.CommandModel{DB: db},
		Users:    models.UserModel{DB: db},
	}

	log.Println("Starting server on port 8000...")
	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: Routes(requestHandler),
	}

	err = httpServer.ListenAndServe()
	log.Fatal(err)
}

func loadEnvironment() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("An error occurred when loading the environment variables: %s", err.Error())
	}
}
