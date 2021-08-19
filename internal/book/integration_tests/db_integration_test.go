package integration_tests

//
//import (
//	goerrors "errors"
//	"fmt"
//	"github.com/foxfurry/simple-rest/configs"
//	"github.com/foxfurry/simple-rest/internal/book/db"
//	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
//	"github.com/foxfurry/simple-rest/internal/book/http/errors"
//	dbpool "github.com/foxfurry/simple-rest/internal/common/database"
//	"github.com/spf13/viper"
//	"github.com/stretchr/testify/assert"
//	"log"
//	"testing"
//)
//
//type BookDBMock struct {
//	TestName            string
//	Input               entity.Book
//	InputArray          []entity.Book
//	ExpectedResult      entity.Book
//	ExpectedArrayResult []entity.Book
//	Parameter           string
//	ExpectedRows        int64
//	ExpectedError       error
//}
//
//var repo db.BookDBRepository
//
//func init() {
//	configs.LoadConfig()
//	repo = db.NewBookRepo(dbpool.CreateDBPool(
//		viper.GetString("database_test.host"),
//		viper.GetInt("database_test.port"),
//		viper.GetString("database_test.user"),
//		viper.GetString("database_test.password"),
//		viper.GetString("database_test.dbname"),
//		viper.GetInt("database_test.maxidleconnections"),
//		viper.GetInt("database_test.maxopenconnections"),
//		viper.GetDuration("database_test.maxconnidletime"),
//	))
//}
//
//func TestBookDBRepository_SaveBook(t *testing.T) {
//	saveBookMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//			},
//			ExpectedResult: entity.Book{
//				Title:       "Test 1",
//				Author:      "Test 1",
//				Year:        1,
//				Description: "Test 1",
//			},
//			ExpectedError: nil,
//		},
//		{
//			TestName: "Test Unsuccessful Bad Request",
//			InputArray: []entity.Book{
//				{},
//			},
//			ExpectedResult: entity.Book{},
//			ExpectedError:  errors.BookBadRequest{},
//		},
//		{
//			TestName: "Test Unsuccessful Title Already Exists",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//			},
//			ExpectedResult: entity.Book{},
//			ExpectedError:  errors.BookTitleAlreadyExists{},
//		},
//	}
//
//	for _, c := range saveBookMocks {
//		if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//			t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//		}
//		c := c // To isolate test cases and be sure they won't be changed
//		t.Run(c.TestName, func(t *testing.T) {
//			for _, val := range c.InputArray {
//				log.Println(val)
//				book, err := repo.SaveBook(&val)
//				if err != nil {
//					assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//				} else {
//					assert.True(t, val.EqualNoID(*book))
//				}
//			}
//		})
//	}
//}
//
//func TestBookDBRepository_GetBook(t *testing.T) {
//	getBookMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			Input: entity.Book{
//				Title:       "Test 1",
//				Author:      "Test 1",
//				Year:        1,
//				Description: "Test 1",
//			},
//			ExpectedResult: entity.Book{
//				ID:          1,
//				Title:       "Test 1",
//				Author:      "Test 1",
//				Year:        1,
//				Description: "Test 1",
//			},
//			ExpectedError: nil,
//		},
//		{
//			TestName: "Test Unsuccessful",
//			Input: entity.Book{
//				Title:       "Test 2",
//				Author:      "Test 2",
//				Year:        2,
//				Description: "Test 2",
//			},
//			ExpectedResult: entity.Book{
//				ID:          3,
//				Title:       "Test 2",
//				Author:      "Test 2",
//				Year:        2,
//				Description: "Test 2",
//			},
//			ExpectedError: errors.BookNotFound{},
//		},
//	}
//
//	if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//		t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//	}
//
//	for _, c := range getBookMocks {
//		c := c // To isolate test cases and be sure they won't be changed
//		t.Run(c.TestName, func(t *testing.T) {
//			if _, err := repo.SaveBook(&c.Input); err != nil {
//				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
//			}
//			getBook, err := repo.GetBook(c.ExpectedResult.ID)
//
//			if err != nil {
//				assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//			} else {
//				assert.True(t, c.Input.EqualNoID(*getBook))
//			}
//		})
//	}
//}
//
//func TestBookDBRepository_GetAllBooks(t *testing.T) {
//	getAllBooksMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 2",
//					Author:      "Test 2",
//					Year:        2,
//					Description: "Test 2",
//				},
//				{
//					Title:       "Test 3",
//					Author:      "Test 3",
//					Year:        3,
//					Description: "Test 3",
//				},
//			},
//			ExpectedArrayResult: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 2",
//					Author:      "Test 2",
//					Year:        2,
//					Description: "Test 2",
//				},
//				{
//					Title:       "Test 3",
//					Author:      "Test 3",
//					Year:        3,
//					Description: "Test 3",
//				},
//			},
//			ExpectedError: nil,
//		},
//		{
//			TestName:            "Test Unsuccessful",
//			InputArray:          []entity.Book{},
//			ExpectedArrayResult: []entity.Book{},
//			ExpectedError:       errors.BookNotFound{},
//		},
//	}
//
//	for _, c := range getAllBooksMocks {
//		if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//			t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//		}
//		c := c
//		t.Run(c.TestName, func(t *testing.T) {
//			for _, val := range c.InputArray {
//				if _, err := repo.SaveBook(&val); err != nil {
//					t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
//				}
//			}
//			books, err := repo.GetAllBooks()
//
//			if err != nil {
//				assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//			} else {
//				assert.True(t, entity.BookArrayEqualNoID(c.ExpectedArrayResult, books))
//			}
//		})
//	}
//}
//
//func TestBookDBRepository_SearchByAuthor(t *testing.T) {
//	searchByAuthorMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test Author",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 2",
//					Author:      "Test Author",
//					Year:        2,
//					Description: "Test 2",
//				},
//				{
//					Title:       "Test 3",
//					Author:      "Test NonAuthor",
//					Year:        3,
//					Description: "Test 3",
//				},
//				{
//					Title:       "Test 4",
//					Author:      "Test Author",
//					Year:        4,
//					Description: "Test 4",
//				},
//			},
//			ExpectedArrayResult: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test Author",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 2",
//					Author:      "Test Author",
//					Year:        2,
//					Description: "Test 2",
//				},
//				{
//					Title:       "Test 4",
//					Author:      "Test Author",
//					Year:        4,
//					Description: "Test 4",
//				},
//			},
//			ExpectedError: nil,
//			Parameter:     "Test Author",
//		},
//		{
//			TestName:       "Test Unsuccessful",
//			Input:          entity.Book{},
//			ExpectedResult: entity.Book{},
//			ExpectedError:  errors.BookNotFoundByAuthor{Author: "Test Author"},
//			Parameter:      "Test Author",
//		},
//		{
//			TestName:       "Test Unsuccessful",
//			Input:          entity.Book{},
//			ExpectedResult: entity.Book{},
//			ExpectedError:  errors.BookBadRequest{},
//			Parameter:      "",
//		},
//	}
//
//	for _, c := range searchByAuthorMocks {
//		if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//			t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//		}
//		c := c // To isolate test cases and be sure they won't be changed
//		t.Run(c.TestName, func(t *testing.T) {
//			for _, val := range c.InputArray {
//				if _, err := repo.SaveBook(&val); err != nil {
//					t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
//				}
//			}
//
//			books, err := repo.SearchByAuthor(c.Parameter)
//
//			if err != nil {
//				assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//			} else {
//				assert.True(t, entity.BookArrayEqualNoID(c.ExpectedArrayResult, books), "Book arrays are not same:\nExpected: %v\nActual: %v", c.ExpectedArrayResult, books)
//			}
//		})
//	}
//}
//
//func TestBookDBRepository_SearchByTitle(t *testing.T) {
//	searchByTitleMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 2",
//					Author:      "Test 2",
//					Year:        2,
//					Description: "Test 2",
//				},
//			},
//			ExpectedResult: entity.Book{
//				Title:       "Test 2",
//				Author:      "Test 2",
//				Year:        2,
//				Description: "Test 2",
//			},
//			ExpectedError: nil,
//			Parameter:     "Test 2",
//		},
//		{
//			TestName:            "Test Unsuccessful Not Found",
//			InputArray:          []entity.Book{},
//			ExpectedArrayResult: []entity.Book{},
//			ExpectedError:       errors.BookNotFoundByTitle{Title: "Test"},
//			Parameter:           "Test",
//		},
//		{
//			TestName:            "Test Unsuccessful Empty Title",
//			InputArray:          []entity.Book{},
//			ExpectedArrayResult: []entity.Book{},
//			ExpectedError:       errors.BookBadRequest{},
//			Parameter:           "",
//		},
//	}
//
//	for _, c := range searchByTitleMocks {
//		if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//			t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//		}
//		c := c // To isolate test cases and be sure they won't be changed
//		t.Run(c.TestName, func(t *testing.T) {
//			for _, val := range c.InputArray {
//				if _, err := repo.SaveBook(&val); err != nil {
//					t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
//				}
//			}
//			book, err := repo.SearchByTitle(c.Parameter)
//
//			if err != nil {
//				assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//			} else {
//				assert.True(t, c.ExpectedResult.EqualNoID(*book))
//			}
//		})
//	}
//
//}
//
//func TestBookDBRepository_UpdateBook(t *testing.T) {
//	updateBookMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 1 Updated",
//					Author:      "Test 1 Updated",
//					Year:        2,
//					Description: "Test 1 Updated",
//				},
//			},
//			ExpectedResult: entity.Book{
//				ID:          1,
//				Title:       "Test 1 Updated",
//				Author:      "Test 1 Updated",
//				Year:        2,
//				Description: "Test 1 Updated",
//			},
//			ExpectedError: nil,
//		},
//		{
//			TestName: "Test Unsuccessful Bad Request Invalid Serial",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{},
//			},
//			ExpectedResult: entity.Book{},
//			ExpectedError:  errors.BookBadRequest{},
//		},
//		{
//			TestName: "Test Unsuccessful Book Not Found",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 1 Updated",
//					Author:      "Test 1 Updated",
//					Year:        2,
//					Description: "Test 1 Updated",
//				},
//			},
//			ExpectedResult: entity.Book{
//				ID:          2,
//				Title:       "Test 1 Updated",
//				Author:      "Test 1 Updated",
//				Year:        2,
//				Description: "Test 1 Updated",
//			},
//			ExpectedError: errors.BookNotFound{},
//		},
//		{
//			TestName: "Test Unsuccessful Bad Request Invalid Request",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{},
//			},
//			ExpectedResult: entity.Book{
//				ID: 666,
//			},
//			ExpectedError: errors.BookBadRequest{},
//		},
//	}
//
//	for _, c := range updateBookMocks {
//		if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//			t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//		}
//		c := c // To isolate test cases and be sure they won't be changed
//		t.Run(c.TestName, func(t *testing.T) {
//			_, err := repo.SaveBook(&c.InputArray[0])
//			if err != nil {
//				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
//			}
//
//			book, err := repo.UpdateBook(c.ExpectedResult.ID, &c.InputArray[1])
//
//			if err != nil {
//				assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//			} else {
//				assert.True(t, c.ExpectedResult.Equal(*book))
//			}
//		})
//	}
//}
//
//func TestBookDBRepository_DeleteBook(t *testing.T) {
//	deleteBookMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			Input: entity.Book{
//				ID:          1,
//				Title:       "Test 1",
//				Author:      "Test 1",
//				Year:        1,
//				Description: "Test 1",
//			},
//			ExpectedRows:  1,
//			ExpectedError: nil,
//		},
//		{
//			TestName: "Test Unsuccessful Book Not Found",
//			Input: entity.Book{
//				ID:          2,
//				Title:       "Test 1",
//				Author:      "Test 1",
//				Year:        1,
//				Description: "Test 1",
//			},
//			ExpectedRows:  1,
//			ExpectedError: errors.BookNotFound{},
//		},
//		{
//			TestName: "Test Unsuccessful Invalid ID",
//			Input: entity.Book{
//				ID:          0,
//				Title:       "Test 1",
//				Author:      "Test 1",
//				Year:        1,
//				Description: "Test 1",
//			},
//			ExpectedRows:  0,
//			ExpectedError: errors.BookBadRequest{},
//		},
//	}
//
//	for _, c := range deleteBookMocks {
//		if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//			t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//		}
//		c := c // To isolate test cases and be sure they won't be changed
//		t.Run(c.TestName, func(t *testing.T) {
//			_, err := repo.SaveBook(&c.Input)
//			if err != nil {
//				t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
//			}
//			rows, err := repo.DeleteBook(c.Input.ID)
//			if err != nil {
//				assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//			} else {
//				assert.Equal(t, rows, c.ExpectedRows)
//			}
//		})
//	}
//}
//
//func TestBookDBRepository_DeleteAllBooks(t *testing.T) {
//	deleteAllBooksMocks := []BookDBMock{
//		{
//			TestName: "Test Successful",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 2",
//					Author:      "Test 2",
//					Year:        2,
//					Description: "Test 2",
//				},
//			},
//			ExpectedRows:  2,
//			ExpectedError: nil,
//		},
//		{
//			TestName: "Test Successful",
//			InputArray: []entity.Book{
//				{
//					Title:       "Test 1",
//					Author:      "Test 1",
//					Year:        1,
//					Description: "Test 1",
//				},
//				{
//					Title:       "Test 2",
//					Author:      "Test 2",
//					Year:        2,
//					Description: "Test 2",
//				},
//				{
//					Title:       "Test 3",
//					Author:      "Test 3",
//					Year:        3,
//					Description: "Test 3",
//				},
//				{
//					Title:       "Test 4",
//					Author:      "Test 4",
//					Year:        4,
//					Description: "Test 4",
//				},
//				{
//					Title:       "Test 5",
//					Author:      "Test 5",
//					Year:        5,
//					Description: "Test 5",
//				},
//			},
//			ExpectedRows:  5,
//			ExpectedError: nil,
//		},
//	}
//	for _, c := range deleteAllBooksMocks {
//		if _, err := repo.DeleteAllBooks(); err != nil && !goerrors.Is(err, errors.BookNotFound{}) {
//			t.Errorf(fmt.Sprintf("Could not delete all the books: %v", err))
//		}
//		c := c // To isolate test cases and be sure they won't be changed
//		t.Run(c.TestName, func(t *testing.T) {
//			for _, val := range c.InputArray {
//				if _, err := repo.SaveBook(&val); err != nil {
//					t.Errorf(fmt.Sprintf("Could not save the book: %v: ", err))
//				}
//			}
//
//			rows, err := repo.DeleteAllBooks()
//			if err != nil {
//				assert.True(t, goerrors.Is(err, c.ExpectedError), "Errors are not same:\nExpected: %v\nActual: %v", c.ExpectedError, err)
//			} else {
//				assert.Equal(t, rows, c.ExpectedRows)
//			}
//		})
//	}
//}
