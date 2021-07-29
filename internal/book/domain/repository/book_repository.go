package repository

import "github.com/foxfurry/simple-rest/internal/book/domain/entity"


type BookRepository interface {
	SaveBook(book *entity.Book) (*entity.Book, error)
	GetBook(uint64) (*entity.Book, error)
	GetAllBooks() ([]entity.Book, error)
	SearchByAuthor(author string) ([]entity.Book, error)
	SearchByTitle(title string) (*entity.Book, error)
}