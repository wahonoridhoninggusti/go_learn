package main

import (
	"fmt"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token != "secret" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}
func secureHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are authorized!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.Handle("/secure", AuthMiddleware(http.HandlerFunc(secureHandler)))

	serve := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	fmt.Println("Server running at http://localhost:8081")

	if err := serve.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
