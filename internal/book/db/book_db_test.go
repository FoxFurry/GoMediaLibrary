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
	"testing"
)

var repo BookDBRepository

var genericTestCases []entity.Book

func init() {
	configs.LoadConfig()

	repo = NewBookRepo(dbpool.CreateDBPool(
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetInt("database.maxidleconnections"),
		viper.GetInt("database.maxopenconnections"),
		viper.GetDuration("database.maxconnidletime"),
	))

	genericTestCases = []entity.Book{
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
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	for _, c := range genericTestCases {
		c := c // To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			book, err := repo.SaveBook(&c)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}

			assert.True(t, c.EqualNoID(*book))
		})
	}
}

func TestBookDBRepository_GetBook(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	for _, c := range genericTestCases {
		c := c // To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			savedBook, err := repo.SaveBook(&c)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
			getBook, err := repo.GetBook(savedBook.ID)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not get the book: %v: ", err))
			}

			assert.True(t, c.EqualNoID(*getBook))
		})
	}
}

func TestBookDBRepository_GetAllBooks(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	for _, c := range genericTestCases {
		c := c
		t.Run(c.Title, func(t *testing.T) {
			if _, err := repo.SaveBook(&c); err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
		})
	}
	allBooks, err := repo.GetAllBooks()

	if err != nil {
		t.Errorf(fmt.Sprintf("Could not get all books: %v: ", err))
	}

	assert.True(t, entity.BookArrayEqualNoID(allBooks, genericTestCases))
}

func TestBookDBRepository_SearchByAuthor(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	authorDifferentTestCases := []entity.Book{
		{
			Title:       "Test 1",
			Author:      "Test 1",
			Year:        1,
			Description: "Test 1",
		},
		{
			Title:       "Test 2",
			Author:      "Test 2",
			Year:        2,
			Description: "Test 2",
		},
		{
			Title:       "Test 3",
			Author:      "Test 3",
			Year:        3,
			Description: "Test 3",
		},
		{
			Title:       "Test 4",
			Author:      "Test 4",
			Year:        4,
			Description: "Test 4",
		},
	}

	authorSameTestCases := []entity.Book{
		{
			Title:       "Test 1",
			Author:      "Test",
			Year:        1,
			Description: "Test 1",
		},
		{
			Title:       "Test 2",
			Author:      "Test",
			Year:        2,
			Description: "Test",
		},
		{
			Title:       "Test 3",
			Author:      "Test",
			Year:        3,
			Description: "Test 3",
		},
		{
			Title:       "Test 4",
			Author:      "Test",
			Year:        4,
			Description: "Test 4",
		},
	}

	for _, c := range authorDifferentTestCases {
		c := c // To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			if _, err := repo.SaveBook(&c); err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
			getBook, err := repo.SearchByAuthor(c.Author)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not get the book by author: %v: ", err))
			}

			if len(getBook) != 1 {
				t.Errorf(fmt.Sprintf("Expected one book as result, got %v", len(getBook)))
			}
			assert.True(t, c.EqualNoID(getBook[0]))
		})
	}

	for _, c := range authorSameTestCases {
		c := c
		t.Run(c.Title, func(t *testing.T) {
			if _, err := repo.SaveBook(&c); err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
		})
	}

	searchedBooks, err := repo.SearchByAuthor(authorSameTestCases[0].Author)
	if err != nil {
		t.Errorf(fmt.Sprintf("Could not search by author: %v", err))
	}

	assert.True(t, entity.BookArrayEqualNoID(authorSameTestCases, searchedBooks))
}

func TestBookDBRepository_SearchByTitle(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	titleDifferentTestCase := []entity.Book{
		{
			Title:       "Test 1",
			Author:      "Test 1",
			Year:        1,
			Description: "Test 1",
		},
		{
			Title:       "Test 2",
			Author:      "Test 2",
			Year:        2,
			Description: "Test 2",
		},
		{
			Title:       "Test 3",
			Author:      "Test 3",
			Year:        3,
			Description: "Test 3",
		},
		{
			Title:       "Test 4",
			Author:      "Test 4",
			Year:        4,
			Description: "Test 4",
		},
	}

	for _, c := range titleDifferentTestCase {
		c := c // To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			if _, err := repo.SaveBook(&c); err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
			getBook, err := repo.SearchByTitle(c.Title)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not get the book by title: %v: ", err))
			}

			assert.True(t, c.EqualNoID(*getBook))
		})
	}

}

func TestBookDBRepository_UpdateBook(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	bookBeforeUpdate, err := repo.SaveBook(&genericTestCases[0])
	if err != nil {
		t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
	}

	bookAfterUpdate, err := repo.UpdateBook(bookBeforeUpdate.ID, &genericTestCases[1])
	if err != nil {
		t.Errorf(fmt.Sprintf("Could not update the book: %v: ", err))
	}

	assert.True(t, !bookBeforeUpdate.Equal(*bookAfterUpdate))

	bookAfterUpdateGet, err := repo.GetBook(bookBeforeUpdate.ID)
	if err != nil {
		t.Errorf(fmt.Sprintf("Could not get the book: %v: ", err))
	}

	assert.True(t, bookAfterUpdateGet.EqualNoID(genericTestCases[1]))
}

func TestBookDBRepository_DeleteBook(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	for _, c := range genericTestCases {
		c := c // To isolate test cases and be sure they won't be changed
		t.Run(c.Title, func(t *testing.T) {
			book, err := repo.SaveBook(&c)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}

			if _, err = repo.DeleteBook(book.ID); err != nil {
				t.Errorf(fmt.Sprintf("Could not delete the book: %v: ", err))
			}
			_, err = repo.GetBook(book.ID)

			assert.Equal(t, err, errors.BookNotFound{})
		})
	}
}

func TestBookDBRepository_DeleteAllBooks(t *testing.T) {
	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
	}

	_, err := repo.GetBook(0)

	assert.Equal(t, err, errors.BookNotFound{})
}

