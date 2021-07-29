package db

import (
	"database/sql"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/domain/repository"
)

type BookRepo struct {
	db *sql.DB:
}

func NewBookRepository(db *sql.DB) *BookRepo {
	return &BookRepo{db: db}
}

var _ repository.BookRepository = &BookRepo{}

func (r *BookRepo) SaveBook(book *entity.Book) *entity.Book {

}

func (r *BookRepo) GetBook(uint64) (*entity.Book, error) {

}

func (r *BookRepo) GetAllBooks() ([]entity.Book, error) {

}

func (r *BookRepo) SearchByAuthor(author string) ([]entity.Book, error) {

}

func (r *BookRepo) SearchByTitle(title string) (*entity.Book, error) {

}
