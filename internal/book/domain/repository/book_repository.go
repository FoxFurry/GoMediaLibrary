package repository

import "github.com/foxfurry/simple-rest/internal/book/domain/entity"

type BookRepository interface {
	SaveBook(*entity.Book) (*entity.Book, error)
	GetBook(uint64) (*entity.Book, error)
	GetAllBooks() ([]entity.Book, error)
	SearchByAuthor(string) ([]entity.Book, error) // An author can have multiple books
	SearchByTitle(string) (*entity.Book, error)
	UpdateBook(uint64, *entity.Book) (int64, error)
	DeleteBook(uint64) (int64, error)
	DeleteAllBooks() (int64, error)
}

