package main

import (
	// "encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wahonoridhoninggusti/go_learn/restful-book/api"
	// "github.com/google/uuid"
)

func main() {
	serve := &http.Server{
		Addr:         ":8081",
		Handler:      api.NewRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	log.Println("Server starting on :8081")
	if err := serve.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
