package response

import "errors"

var (
	ErrBookNotFound     = errors.New("book not found")
	ErrBookAlreadyExist = errors.New("book already exists")
	ErrInvalidBookID    = errors.New("invalid book ID")
	ErrEmptyBookTitle   = errors.New("book title cannot be empty")
	ErrNoBooks          = errors.New("no books available")
)
