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

func (r *BookDBRepository) SaveBook(book *entity.Book) (*entity.Book, error) {
	query := `INSERT INTO bookstore (title, author, year, description) VALUES ($1, $2, $3, $4) RETURNING id`

	var bookID uint64

	err := r.database.QueryRow(query, book.Title, book.Author, book.Year, book.Description).Scan(&bookID)

	if err != nil {
		log.Printf("Unable to save book to db: %v", err)
		return book, err
	}

	returnBook := *book
	returnBook.ID = bookID
	return &returnBook, nil
}

func (r *BookDBRepository) GetBook(bookID uint64) (*entity.Book, error) {
	var book entity.Book

	query := `SELECT * FROM bookstore WHERE id=$1`

	row := r.database.QueryRow(query, bookID)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)

	switch err {
	case sql.ErrNoRows:
		log.Printf("Book id#%v not found", bookID)
		return &book, errors.BookNotFound{}
	case nil:
		return &book, nil
	default:
		log.Printf("Unable to scan the row: %v", err)
		return nil, err
	}
}

func (r *BookDBRepository) GetAllBooks() ([]entity.Book, error) {
	var books []entity.Book

	query := `SELECT * FROM bookstore`

	rows, err := r.database.Query(query)

	if err != nil {
		log.Printf("Unable to get all books: %v", err)
		return books, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Panicf("Could not close the db rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		var tempBook entity.Book

		err = rows.Scan(&tempBook.ID, &tempBook.Title, &tempBook.Author, &tempBook.Year, &tempBook.Description)

		if err != nil {
			log.Printf("Unable to scan the user: %v", err)
		}

		books = append(books, tempBook)
	}

	if len(books) == 0 {
		log.Printf("Could not get all the books\n")
		return books, errors.BookNotFound{}
	}

	return books, nil
}

func (r *BookDBRepository) SearchByAuthor(author string) ([]entity.Book, error) {
	var books []entity.Book

	query := `SELECT * FROM bookstore WHERE author=$1`

	rows, err := r.database.Query(query, author)

	if err != nil {
		log.Printf("Could not get all books with author %v: %v", author, err)
		return books, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Panicf("Could not close the db rows: %v", err)
		}
	}(rows)

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

	query := `SELECT * FROM bookstore WHERE title=$1`

	row := r.database.QueryRow(query, title)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)

	switch err {
	case sql.ErrNoRows:
		log.Printf("Book title#%v not found", title)
		return &book, errors.BookNotFoundByTitle{Title: title}
	case nil:
		return &book, nil
	default:
		log.Fatalf("Unable to scan the row: %v", err)
		return &book, err
	}
}

func (r *BookDBRepository) UpdateBook(bookID uint64, book *entity.Book) (*entity.Book, error) {
	query := `UPDATE bookstore SET title=$2, author=$3, year=$4, description=$5 WHERE id=$1`

	_, err := r.database.Exec(query, bookID, book.Title, book.Author, book.Year, book.Description)

	returnBook := *book
	returnBook.ID = bookID
	if err != nil {
		log.Printf("Unable to update book: %v", err)
		return &returnBook, err
	}

	return &returnBook, err
}

func (r *BookDBRepository) DeleteBook(bookID uint64) (int64, error) {
	query := `DELETE FROM bookstore WHERE id=$1`

	res, err := r.database.Exec(query, bookID)

	if err != nil {
		log.Printf("Unable to delete book: %v", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Unable to get affected rows book: %v", err)
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, errors.BookNotFound{}
	}

	log.Printf("Rows affected: %v", rowsAffected)

	return rowsAffected, err
}

func (r *BookDBRepository) DeleteAllBooks() (int64, error) {
	query := `TRUNCATE bookstore RESTART IDENTITY`

	res, err := r.database.Exec(query)

	if err != nil {
		log.Printf("Unable to delete book: %v", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("Unable to get affected rows book: %v", err)
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, errors.BookNotFound{}
	}

	log.Printf("Rows affected: %v", rowsAffected)

	return rowsAffected, err
}
