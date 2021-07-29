package db

import (
	"database/sql"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/domain/repository"
	"log"
)

type BookRepo struct{
	database *sql.DB
}

func NewBookRepo(db *sql.DB) BookRepo{
	return BookRepo{database: db}
}

var _ repository.BookRepository = &BookRepo{}

func (r *BookRepo) SaveBook(book *entity.Book) (*entity.Book, error) {
	query := `INSERT INTO bookstore (title, author, year, description) VALUES ($1, $2, $3, $4) RETURNING id`

	var bookID uint64

	err := r.database.QueryRow(query, book.Title, book.Author, book.Year, book.Description).Scan(&bookID)

	if err != nil {
		log.Printf("Unable to save book to db: %v", err)
		return nil, err
	}

	book.ID = bookID

	return book, nil
}

func (r *BookRepo) GetBook(bookID uint64) (*entity.Book, error) {
	var book entity.Book

	query := `SELECT * FROM bookstore WHERE id=$1`

	row := r.database.QueryRow(query, bookID)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)

	switch err {
	case sql.ErrNoRows:
		log.Printf("Book id#%v not found", bookID)
		return &book, err
	case nil:
		return &book, nil
	default:
		log.Fatalf("Unable to scan the row: %v", err)
		return nil,nil
	}
}

func (r *BookRepo) GetAllBooks() ([]entity.Book, error) {
	var books []entity.Book

	query := `SELECT * FROM bookstore`

	rows, err := r.database.Query(query)

	if err != nil {
		log.Printf("Unable to get all books: %v", err)
		return books, err
	}

	defer rows.Close()

	for rows.Next(){
		var tempBook entity.Book

		err := rows.Scan(&tempBook.ID,&tempBook.Title, &tempBook.Author, &tempBook.Year, &tempBook.Description)

		if err != nil {
			log.Printf("Unable to scan the user: %v", err)
		}

		books = append(books, tempBook)
	}
	return books, nil
}

func (r *BookRepo) SearchByAuthor(author string) ([]entity.Book, error) {
	var books []entity.Book

	query := `SELECT * FROM bookstore WHERE author=$1`

	rows, err := r.database.Query(query, author)

	if err != nil {
		log.Printf("Could not get all books with author %v: %v", author, err)
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		var tempBook entity.Book

		err = rows.Scan(&tempBook.ID, &tempBook.Title, &tempBook.Author, &tempBook.Year, &tempBook.Description)

		if err != nil {
			log.Printf("Could not scan the row: %v", err)
			continue
		}

		books = append(books, tempBook)
	}

	return books, nil
}

func (r *BookRepo) SearchByTitle(title string) (*entity.Book, error) {
	var book entity.Book

	query := `SELECT * FROM bookstore WHERE title=$1`

	row := r.database.QueryRow(query, title)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Description)

	switch err {
	case sql.ErrNoRows:
		log.Printf("Book title#%v not found", title)
		return &book, err
	case nil:
		return &book, nil
	default:
		log.Fatalf("Unable to scan the row: %v", err)
		return nil,nil
	}
}
