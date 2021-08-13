package db

import (
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

var testCases []struct {
	mockBook entity.Book
}

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

	testCases = []struct{ mockBook entity.Book }{
		{mockBook: entity.Book{
			Title:       "Test 1",
			Author:      "Test Author 1",
			Year:        1,
			Description: "Test Description 1",
		}},
		{mockBook: entity.Book{
			Title:       "Test 2",
			Author:      "",
			Year:        2,
			Description: "",
		}},
		{mockBook: entity.Book{
			Title:       "Test 3",
			Author:      "",
			Year:        3,
			Description: "Test Description 3",
		}},
		{mockBook: entity.Book{
			Title:       "Test 4",
			Author:      "Test Author 4",
			Description: "",
		}},
	}
}

func TestBookDBRepository_SaveBook(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		c := c
		t.Run(c.mockBook.Title, func(t *testing.T) {
			t.Parallel()
			book, err := repo.SaveBook(&c.mockBook)
			if err != nil {
				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
			}
			c.mockBook.ID = book.ID // We just copy ID to expected mock cuz it's generated
			assert.Equal(t, &c.mockBook, book)
		})
	}
}

func TestBookDBRepository_DeleteBook(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		c := c
		t.Run(c.mockBook.Title, func(t *testing.T) {
			t.Parallel()
			book, err := repo.SaveBook(&c.mockBook)
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
