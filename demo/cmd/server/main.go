package main

import (
	"log"
	"net/http"
	"os"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/api"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/repository"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/service"
)

func main() {
	var repo repository.MessageRepository
	dataPath := os.Getenv("DEMO_DATA_PATH")
	if dataPath != "" {
		fileRepo, err := repository.NewFileMessageRepository(dataPath)
		if err != nil {
			log.Fatalf("failed to open file repository: %v", err)
		}
		repo = fileRepo
		log.Printf("using file repository: %s", dataPath)
	} else {
		repo = repository.NewMemoryMessageRepository()
	}
	svc := service.NewMessageService(repo)
	server := api.NewServer(svc)
	mux := http.NewServeMux()
	server.Register(mux)
	log.Println("demo server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
