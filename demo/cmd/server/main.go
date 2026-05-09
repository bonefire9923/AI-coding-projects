package main

import (
	"log"
	"net/http"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/api"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/repository"
	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/service"
)

func main() {
	repo := repository.NewMemoryMessageRepository()
	svc := service.NewMessageService(repo)
	server := api.NewServer(svc)
	mux := http.NewServeMux()
	server.Register(mux)
	log.Println("demo server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
