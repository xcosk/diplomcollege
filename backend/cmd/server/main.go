package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"golearn/internal/db"
	"golearn/internal/handlers"
	"golearn/internal/middleware"
	"golearn/internal/seed"
)

func main() {
	database, err := db.Open("./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}
	if err := seed.Seed(database); err != nil {
		log.Fatal(err)
	}

	secret := os.Getenv("GOLEARN_JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	app := handlers.NewApp(database, []byte(secret))
	mux := http.NewServeMux()
	app.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           middleware.CORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Backend running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
