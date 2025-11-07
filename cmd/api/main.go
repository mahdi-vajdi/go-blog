package main

import (
	"github.com/mahdi-vajdi/go-blog/internal/server"
	"github.com/mahdi-vajdi/go-blog/internal/store"
	"log"
)

func main() {
	memoryStore := store.NewMemoryStore()

	apiServer := server.NewAPIServer(":8080", memoryStore)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
