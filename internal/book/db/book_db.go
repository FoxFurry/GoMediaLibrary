package db

import (
	"database/sql"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/domain/repository"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	"log"
)

type BookDBRepository struct {
	database *sql.DB
}

func NewBookRepo(db *sql.DB) BookDBRepository {
	return BookDBRepository{database: db}
}

var _ repository.BookRepository = &BookDBRepository{}

const (
	querySaveBook           = `INSERT INTO bookstore (title, author, year, description) VALUES ($1, $2, $3, $4) RETURNING id`
	queryGetBook            = `SELECT * FROM bookstore WHERE id=$1`
	queryGetAll             = `SELECT * FROM bookstore`
	querySearchByAuthorBook = `SELECT * FROM bookstore WHERE author=$1`
	querySearchByTitleBook  = `SELECT * FROM bookstore WHERE title=$1`
	queryUpdateBook         = `UPDATE bookstore SET title=$2, author=$3, year=$4, description=$5 WHERE id=$1`
	queryDeleteBook             = `DELETE FROM bookstore WHERE id=$1`
	queryDeleteAllBooksAndAlter = `DELETE FROM bookstore; ALTER SEQUENCE bookstore_id_seq RESTART WITH 1`
)

func (r *BookDBRepository) SaveBook(book *entity.Book) (*entity.Book, error) {
	if !book.IsValid() {
		log.Printf("Invalid request: %v", book)
		return nil, errors.BookBadRequest{}
	}

	var bookID uint64

	err := r.database.QueryRow(querySaveBook, book.Title, book.Author, book.Year, book.Description).Scan(&bookID)

	if err != nil {
		log.Printf("Unable to save book to db: %v", err)
		return nil, errors.BookBadScanOptions{Msg: err.Error()}
	}

	returnBook := *book
	returnBook.ID = bookID
	return &returnBook, nil
}

func (r *BookDBRepository) GetBook(bookID uint64) (*entity.Book, error) {
	var book entity.Book

	row := r.database.QueryRow(queryGetBook, bookID)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)

	if err == sql.ErrNoRows {
		log.Printf("Book id#%v not found", bookID)
		return nil, errors.BookNotFound{}
	}

	return &book, nil
}

func (r *BookDBRepository) GetAllBooks() ([]entity.Book, error) {
	var books []entity.Book

	rows, err := r.database.Query(queryGetAll)
	if err != nil {
		log.Printf("Unable to get all books: %v", err)
		return nil, errors.BookCouldNotQuery{Msg: err.Error()}
	}

	defer rows.Close()

	for rows.Next() {
		var tempBook entity.Book
		err = rows.Scan(&tempBook.ID, &tempBook.Title, &tempBook.Author, &tempBook.Year, &tempBook.Description)

		if err != nil {
			log.Printf("Unable to scan the user: %v", err)
			continue
		}

		books = append(books, tempBook)
	}

	if len(books) == 0 {
		log.Printf("Could not get all the books\n")
		return nil, errors.BookNotFound{}
	}

	return books, nil
}

func (r *BookDBRepository) SearchByAuthor(author string) ([]entity.Book, error) {
	if author == "" {
		log.Printf("Author field is empty")
		return nil, errors.BookBadRequest{}
	}

	rows, err := r.database.Query(querySearchByAuthorBook, author)

	if err != nil {
		log.Printf("Could not get all books with author %v: %v", author, err)
		return nil, errors.BookCouldNotQuery{Msg: err.Error()}
	}

	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var tempBook entity.Book

		err = rows.Scan(&tempBook.ID, &tempBook.Title, &tempBook.Author, &tempBook.Year, &tempBook.Description)

		if err != nil {
			log.Printf("Could not scan the row: %v", err)
			continue
		}

		books = append(books, tempBook)
	}

	if len(books) == 0 {
		log.Printf("Could not get all the books by author: %v\n", author)
		return books, errors.BookNotFoundByAuthor{Author: author}
	}

	return books, nil
}

func (r *BookDBRepository) SearchByTitle(title string) (*entity.Book, error) {
	var book entity.Book

	if title == "" {
		log.Printf("Title field is empty")
		return nil, errors.BookBadRequest{}
	}

	row := r.database.QueryRow(querySearchByTitleBook, title)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)
	if err == sql.ErrNoRows {
		log.Printf("Book title#%v not found", title)
		return nil, errors.BookNotFoundByTitle{Title: title}
	}

	return &book, nil
}

func (r *BookDBRepository) UpdateBook(bookID uint64, book *entity.Book) (*entity.Book, error) {
	if bookID < 1 {
		log.Printf("Serial is less than 1")
		return nil, errors.BookBadRequest{}
	} else if !book.IsValid() {
		log.Printf("Invalid request: %v", book)
		return book, errors.BookBadRequest{}
	}

	_, err := r.database.Exec(queryUpdateBook, bookID, book.Title, book.Author, book.Year, book.Description)

	returnBook := *book
	returnBook.ID = bookID
	if err != nil {
		log.Printf("Unable to update book: %v", err)
		return nil, errors.BookCouldNotQuery{Msg: err.Error()}
	}

	return &returnBook, nil
}

func (r *BookDBRepository) DeleteBook(bookID uint64) (int64, error) {
	if bookID < 1 {
		log.Printf("Serial is less than 1")
		return 0, errors.BookBadRequest{}
	}
	res, err := r.database.Exec(queryDeleteBook, bookID)

	if err != nil {
		log.Printf("Unable to delete book: %v", err)
		return 0, errors.BookCouldNotQuery{Msg: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Unable to get affected rows book: %v", err)
		return 0, errors.BookCouldNotQuery{Msg: err.Error()}
	}

	if rowsAffected == 0 {
		return 0, errors.BookNotFound{}
	}

	log.Printf("Deleted rows: %v", rowsAffected)

	return rowsAffected, err
}

func (r *BookDBRepository) DeleteAllBooks() (int64, error) {
	res, err := r.database.Exec(queryDeleteAllBooksAndAlter)

	if err != nil {
		log.Printf("Unable to delete book or alter the sequence: %v", err)
		return 0, errors.BookCouldNotQuery{Msg: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Unable to get affected rows book: %v", err)
		return 0, errors.BookCouldNotQuery{Msg: err.Error()}
	}

	if rowsAffected == 0 {
		return 0, errors.BookNotFound{}
	}

	log.Printf("Rows affected: %v", rowsAffected)

	return rowsAffected, err
}
