package db

import (
	goerrors "errors"
	"fmt"
	"github.com/foxfurry/simple-rest/configs"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	dbpool "github.com/foxfurry/simple-rest/internal/common/database"
	_ "github.com/foxfurry/simple-rest/internal/common/testing"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var repo BookDBRepository

var testCases []entity.Book

func init() {
	configs.LoadConfig()

	db := dbpool.CreateDBPool(
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetInt("database.maxidleconnections"),
		viper.GetInt("database.maxopenconnections"),
		viper.GetDuration("database.maxconnidletime"),
	)

	repo = BookDBRepository{database: db}

	testCases = []entity.Book{
		{
			Title:       "Test 1",
			Author:      "Test Author 1",
			Year:        1,
			Description: "Test Description 1",
				},
		{
			Title:       "Test 2",
			Author:      "",
			Year:        2,
			Description: "",
				},
		{
			Title:       "Test 3",
			Author:      "",
			Year:        3,
			Description: "Test Description 3",
				},
		{
			Title:       "Test 4",
			Author:      "Test Author 4",
			Description: "",
		},
	}
}


func TestBookDBRepository_SaveBook(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		c := c			// To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			t.Parallel()
			book, err := repo.SaveBook(&c)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
			//c.Author = "Bullshit" // We just copy ID to expected mock cuz it's generated
			log.Printf("---------------\n%v\n%v\n---------------------", c, book)

			assert.True(t, c.EqualNoID(*book))
		})
	}
}

func TestBookDBRepository_DeleteBook(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		c := c			// To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			t.Parallel()
			book, err := repo.SaveBook(&c)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
			_, err = repo.DeleteBook(book.ID)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not delete the book: %v: ", err))
			}
			book, err = repo.GetBook(book.ID)

			assert.Equal(t, err, errors.BookNotFound{})
		})
	}
}

// This test is not parallel, because order of queries matter
func TestBookDBRepository_GetAllBooks(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	for _, c := range testCases {
		c := c			// To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {

			if _, err := repo.SaveBook(&c); err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
		})
	}
	books, err := repo.GetAllBooks()

	if err != nil {
		t.Errorf(fmt.Sprintf("Could not get all books: %v: ", err))
	}
	for v, book := range books {
		book.ID = testCases[v].ID
		assert.Equal(t, testCases[v], book)
	}
}

func TestBookDBRepository_GetBook(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		c := c			// To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			t.Parallel()
			savedBook, err := repo.SaveBook(&c)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
			getBook, err := repo.GetBook(savedBook.ID)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not get the book: %v: ", err))
			}

			assert.Equal(t, &c, getBook)
		})
	}
}