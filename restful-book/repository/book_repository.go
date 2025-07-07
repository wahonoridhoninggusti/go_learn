package repository

import (
	"fmt"
	"sync"

	"github.com/wahonoridhoninggusti/go_learn/restful-book/domain/models"
	"github.com/wahonoridhoninggusti/go_learn/restful-book/domain/response"
)

type BookRepository interface {
	GetAll() ([]*models.Book, error)
	GetByID(id string) (*models.Book, error)
	Create(book *models.Book) error
	Update(id string, book models.Book) error
	Delete(id string) error
	SearchByAuthor(author string) ([]*models.Book, error)
	SearchByTitle(title string) ([]*models.Book, error)
}

func NewBookRepository() BookRepository {
	return &BookRepo{
		books: []models.Book{},
	}
}

type BookRepo struct {
	books []models.Book
	mu    sync.RWMutex
}

// Create implements BookRepository.
func (b *BookRepo) Create(book *models.Book) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, s := range b.books {
		if s.Title == book.Title && s.Author == book.Author {
			return response.ErrBookAlreadyExist
		}
	}

	b.books = append(b.books, *book)
	return nil
}

// Delete implements BookRepository.
func (b *BookRepo) Delete(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, books := range b.books {
		if books.ID == id {
			b.books = append(b.books[:i], b.books[i+1:]...)
			return nil
		}
	}
	return response.ErrBookNotFound
}

// GetAll implements BookRepository.
func (b *BookRepo) GetAll() ([]*models.Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.books) == 0 {
		return nil, response.ErrNoBooks
	}
	booksData := make([]*models.Book, 0, len(b.books))
	for _, b := range b.books {
		book := b
		booksData = append(booksData, &book)
	}
	return booksData, nil
}

// GetByID implements BookRepository.
func (b *BookRepo) GetByID(id string) (*models.Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for i := range b.books {
		if b.books[i].ID == id {
			return &b.books[i], nil
		}
	}
	return nil, response.ErrBookNotFound
}

// SearchByAuthor implements BookRepository.
func (b *BookRepo) SearchByAuthor(author string) ([]*models.Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(b.books) == 0 {
		return nil, response.ErrNoBooks
	}
	booksData := make([]*models.Book, 0, len(b.books))

	for i, books := range b.books {
		if b.books[i].Author == author {
			data := books
			booksData = append(booksData, &data)
		}
	}

	return booksData, nil
}

// SearchByTitle implements BookRepository.
func (b *BookRepo) SearchByTitle(title string) ([]*models.Book, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(b.books) == 0 {
		return nil, response.ErrNoBooks
	}
	booksData := make([]*models.Book, 0, len(b.books))

	for i, books := range b.books {
		if b.books[i].Title == title {
			data := books
			booksData = append(booksData, &data)
		}
	}

	return booksData, nil
}

// Update implements BookRepository.
func (b *BookRepo) Update(id string, book models.Book) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, books := range b.books {
		if books.ID == id {
			b.books[i] = book

			fmt.Println(b.books[i])
			return nil
		}
	}
	return response.ErrBookNotFound
}
