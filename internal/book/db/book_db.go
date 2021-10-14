package db

import (
	"database/sql"
	"github.com/foxfurry/medialib/internal/book/domain/entity"
	"github.com/foxfurry/medialib/internal/book/domain/repository"
	"github.com/foxfurry/medialib/internal/book/http/errors"
	"github.com/foxfurry/medialib/internal/book/http/validators"
	ct "github.com/foxfurry/medialib/internal/common/server/common_translators"
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
	QuerySaveBook = `INSERT INTO bookstore (title, author, year, description) VALUES ($1, $2, $3, $4) RETURNING id`
	QueryGetBook            = `SELECT * FROM bookstore WHERE id=$1`
	QueryGetAll             = `SELECT * FROM bookstore`
	QuerySearchByAuthorBook = `SELECT * FROM bookstore WHERE author=$1`
	QuerySearchByTitleBook = `SELECT * FROM bookstore WHERE title=$1`
	QueryUpdateBook             = `UPDATE bookstore SET title=$2, author=$3, year=$4, description=$5 WHERE id=$1`
	QueryDeleteBook             = `DELETE FROM bookstore WHERE id=$1`
	QueryDeleteAllBooksAndAlter = `DELETE FROM bookstore; ALTER SEQUENCE bookstore_id_seq RESTART WITH 1`
)

func (r *BookDBRepository) SaveBook(book *entity.Book) (*entity.Book, error) {
	var bookID uint64

	err := r.database.QueryRow(QuerySaveBook, book.Title, book.Author, book.Year, book.Description).Scan(&bookID)

	if err != nil {
		log.Printf("Unable to save book to db: %v", err)
		return nil, errors.NewBookCouldNotQuery(err.Error())
	}

	returnBook := *book
	returnBook.ID = bookID
	return &returnBook, nil
}

func (r *BookDBRepository) GetBook(bookID uint64) (*entity.Book, error) {
	if bookID < 1 {
		log.Printf("Serial is less than 1")
		return nil, errors.NewBookInvalidSerial()
	}
	var book entity.Book

	row := r.database.QueryRow(QueryGetBook, bookID)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)

	if err == sql.ErrNoRows {
		log.Printf("Book id#%v not found", bookID)
		return nil, errors.NewBooksNotFound()
	} else if err != nil {
		log.Printf("Could not execute the query: %v", err)
		return nil, errors.NewBookCouldNotQuery(err.Error())
	}

	return &book, nil
}

func (r *BookDBRepository) GetAllBooks() ([]entity.Book, error) {
	var books []entity.Book

	rows, err := r.database.Query(QueryGetAll)
	if err != nil {
		log.Printf("Unable to get all books: %v", err)
		return nil, errors.NewBookCouldNotQuery(err.Error())
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
		return nil, errors.NewBooksNotFound()
	}

	return books, nil
}

func (r *BookDBRepository) SearchByAuthor(author string) ([]entity.Book, error) {
	if author == "" {
		log.Printf("Author field is empty")
		return nil, errors.NewBookValidatorError([]ct.FieldError{validators.FieldAuthorEmpty})
	}
	rows, err := r.database.Query(QuerySearchByAuthorBook, author)

	if err != nil {
		log.Printf("Could not get all books with author %v: %v", author, err)
		return nil, errors.NewBookCouldNotQuery(err.Error())
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
		return books, errors.NewBookNotFoundByAuthor(author)
	}

	return books, nil
}

func (r *BookDBRepository) SearchByTitle(title string) (*entity.Book, error) {
	var book entity.Book

	if title == "" {
		log.Printf("Title field is empty")
		return nil, errors.NewBookValidatorError([]ct.FieldError{validators.FieldTitleEmpty})
	}

	row := r.database.QueryRow(QuerySearchByTitleBook, title)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)
	if err == sql.ErrNoRows {
		log.Printf("Book title#%v not found", title)
		return nil, errors.NewBookNotFoundByTitle(title)
	} else if err != nil {
		log.Printf("Could not execute the query: %v", err)
		return nil, errors.NewBookCouldNotQuery(err.Error())
	}

	return &book, nil
}

func (r *BookDBRepository) UpdateBook(bookID uint64, book *entity.Book) (*entity.Book, error) {
	if bookID < 1 {
		log.Printf("Serial is less than 1")
		return nil, errors.NewBookInvalidSerial()
	}
	_, err := r.database.Exec(QueryUpdateBook, bookID, book.Title, book.Author, book.Year, book.Description)

	returnBook := *book
	returnBook.ID = bookID
	if err != nil {
		log.Printf("Unable to update book: %v", err)
		return nil, errors.NewBookCouldNotQuery(err.Error())
	}

	return &returnBook, nil
}

func (r *BookDBRepository) DeleteBook(bookID uint64) (int64, error) {
	if bookID < 1 {
		log.Printf("Serial is less than 1")
		return 0, errors.NewBookInvalidSerial()
	}

	res, err := r.database.Exec(QueryDeleteBook, bookID)

	if err != nil {
		log.Printf("Unable to delete book: %v", err)
		return 0, errors.NewBookCouldNotQuery(err.Error())
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Unable to get affected rows book: %v", err)
		return 0, errors.NewBookCouldNotQuery(err.Error())
	}

	if rowsAffected == 0 {
		return 0, errors.NewBooksNotFound()
	}

	log.Printf("Deleted rows: %v", rowsAffected)

	return rowsAffected, err
}

func (r *BookDBRepository) DeleteAllBooks() (int64, error) {
	res, err := r.database.Exec(QueryDeleteAllBooksAndAlter)

	if err != nil {
		log.Printf("Unable to delete book or alter the sequence: %v", err)
		return 0, errors.NewBookCouldNotQuery(err.Error())
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Unable to get affected rows book: %v", err)
		return 0, errors.NewBookCouldNotQuery(err.Error())
	}

	if rowsAffected == 0 {
		return 0, errors.NewBooksNotFound()
	}

	log.Printf("Rows affected: %v", rowsAffected)

	return rowsAffected, err
}
