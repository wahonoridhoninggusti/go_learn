package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wahonoridhoninggusti/go_learn/restful-book/domain/models"
	"github.com/wahonoridhoninggusti/go_learn/restful-book/domain/response"
	"github.com/wahonoridhoninggusti/go_learn/restful-book/service"
)

// "github.com/wahonoridhoninggusti/go_learn/restful-book/domain/models"
type BookHandler struct {
	Service service.BookService
}

func NewBookHandler(s service.BookService) *BookHandler {
	return &BookHandler{Service: s}
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	if err := h.Service.CreateBook(&book); err != nil {
		if err == response.ErrBookAlreadyExist {
			response.Error(w, err, http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response.JSON(w, book, "Success", http.StatusCreated)
}

func (h *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	title := query.Get("title")
	author := query.Get("author")

	if title != "" {
		book, err := h.Service.SearchBooksByTitle(title)
		if err != nil {
			response.Error(w, err, http.StatusNotFound)
			return
		}
		response.JSON(w, book, "Success", http.StatusOK)
		return
	}

	if author != "" {
		book, err := h.Service.SearchBooksByAuthor(author)
		if err != nil {
			response.Error(w, err, http.StatusNotFound)
			return
		}
		response.JSON(w, book, "Success", http.StatusOK)
		return
	}

	books, err := h.Service.GetAllBooks()
	if err != nil {
		if err == response.ErrNoBooks {
			response.Error(w, err, http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	response.JSON(w, books, "Success", http.StatusOK)
}

func (h *BookHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	response.JSON(w, book, "Success", http.StatusOK)
}

func (h *BookHandler) PutById(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	book.ID = id

	if err := h.Service.UpdateBook(id, book); err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	response.JSON(w, book, "Updated success", http.StatusOK)
}

func (h *BookHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if err := h.Service.DeleteBook(id); err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	response.JSON(w, nil, "data deleted!", http.StatusOK)
}

func (h *BookHandler) GetByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	book, err := h.Service.SearchBooksByTitle(title)
	if err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	response.JSON(w, book, "Success", http.StatusOK)
}

func (h *BookHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	book, err := h.Service.SearchBooksByAuthor(author)
	if err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	response.JSON(w, book, "Success", http.StatusOK)
}

func extractID(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
