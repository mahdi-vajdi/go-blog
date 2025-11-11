package main

import (
	"github.com/mahdi-vajdi/go-blog/internal/server"
	"github.com/mahdi-vajdi/go-blog/internal/store"
	"log"
)

func main() {
	postgresStore, err := store.NewPostgresStore("postgresql://mahdi:mahdi@localhost:5432/go_blog")
	if err != nil {
		log.Fatalf("Failed to connect to the postgres database: %v", err)
	}

	if err := postgresStore.Init(); err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	log.Println("Database connection successful and tables initialized.")

	apiServer := server.NewAPIServer(":8080", postgresStore)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
