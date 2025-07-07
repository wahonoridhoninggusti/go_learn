package service

import (
	"github.com/google/uuid"
	"github.com/wahonoridhoninggusti/go_learn/restful-book/domain/models"
	"github.com/wahonoridhoninggusti/go_learn/restful-book/repository"
)

type BookService interface {
	GetAllBooks() ([]*models.Book, error)
	GetBookByID(id string) (*models.Book, error)
	CreateBook(book *models.Book) error
	UpdateBook(id string, book models.Book) error
	DeleteBook(id string) error
	SearchBooksByAuthor(author string) ([]*models.Book, error)
	SearchBooksByTitle(title string) ([]*models.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

// CreateBook implements BookService.
func (b *bookService) CreateBook(book *models.Book) error {
	book.ID = uuid.New().String()
	return b.repo.Create(book)
}

// DeleteBook implements BookService.
func (b *bookService) DeleteBook(id string) error {
	return b.repo.Delete(id)
}

// GetAllBooks implements BookService.
func (b *bookService) GetAllBooks() ([]*models.Book, error) {
	return b.repo.GetAll()
}

// GetBookByID implements BookService.
func (b *bookService) GetBookByID(id string) (*models.Book, error) {
	return b.repo.GetByID(id)
}

// SearchBooksByAuthor implements BookService.
func (b *bookService) SearchBooksByAuthor(author string) ([]*models.Book, error) {
	return b.repo.SearchByAuthor(author)
}

// SearchBooksByTitle implements BookService.
func (b *bookService) SearchBooksByTitle(title string) ([]*models.Book, error) {
	return b.repo.SearchByTitle(title)
}

// UpdateBook implements BookService.
func (b *bookService) UpdateBook(id string, book models.Book) error {
	return b.repo.Update(id, book)
}

func NewBookService(r repository.BookRepository) BookService {
	return &bookService{repo: r}
}
