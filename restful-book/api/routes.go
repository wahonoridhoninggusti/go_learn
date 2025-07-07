package api

import (
	"net/http"

	"github.com/wahonoridhoninggusti/go_learn/restful-book/api/handlers"
	"github.com/wahonoridhoninggusti/go_learn/restful-book/repository"
	"github.com/wahonoridhoninggusti/go_learn/restful-book/service"
)

func NewRouter() http.Handler {

	bookRepo := repository.NewBookRepository()
	bookService := service.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookHandler.GetAll(w, r)
		case http.MethodPost:
			bookHandler.Create(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookHandler.GetById(w, r)
		case http.MethodPut:
			bookHandler.PutById(w, r)
		case http.MethodDelete:
			bookHandler.DeleteById(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	return mux
}
