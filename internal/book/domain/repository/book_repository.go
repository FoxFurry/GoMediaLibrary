package repository

import (
	entity2 "github.com/foxfurry/simple-rest/internal/book/domain/entity"
)

type BookRepository interface {
	SaveBook(book *entity2.Book) (*entity2.Book)
	GetBook(uint64) (*entity2.Book, error)
	GetAllBooks() ([]entity2.Book, error)
	SearchByAuthor(author string) ([]entity2.Book, error)
	SearchByTitle(title string) (*entity2.Book, error)

}