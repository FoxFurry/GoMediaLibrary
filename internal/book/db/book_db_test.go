package db

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	_ "github.com/foxfurry/simple-rest/internal/common/tests"
	"log"
	"regexp"
	"testing"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Could not create a new mock: %v", err)
	}

	return db, mock
}

func TestBookDBRepository_SaveBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	saveBookMocks := []struct {
		testName       string
		input          entity.Book
		expectedOutput entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
	}{
		{
			testName: "Test Successful",
			input: entity.Book{
				Title:       "test title",
				Author:      "test author",
				Year:        1,
				Description: "test description",
			},
			expectedOutput: entity.Book{
				ID:          1,
				Title:       "test title",
				Author:      "test author",
				Year:        1,
				Description: "test description",
			},
			expectedError: nil,
			mockFunc: func() {
				rows := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(querySaveBook)).WithArgs("test title", "test author", 1, "test description").WillReturnRows(rows)
			},
			mockRepo: repo,
		},
		{
			testName: "Test Unsuccessful: Invalid Book Request",
			input: entity.Book{
				Author:      "test author",
				Year:        1,
				Description: "test description",
			},
			expectedOutput: entity.Book{},
			expectedError:  errors.BookBadRequest{},
			mockRepo:       repo,
		},
		{
			testName: "Test Unsuccessful: Invalid Repository",
			input: entity.Book{
				Title:       "test title",
				Author:      "test author",
				Year:        1,
				Description: "test description",
			},
			expectedOutput: entity.Book{},
			expectedError:  errors.BookBadScanOptions{Msg: "sql: no rows in result set"},
			mockFunc: func() {
				rows := mock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(querySaveBook)).WithArgs("test title", "test author", 1, "test description").WillReturnRows(rows)
			},
			mockRepo: repo,
		},
	}

	for _, test := range saveBookMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.SaveBook(&test.input)
			if (err != nil) && (err != test.expectedError) {
				t.Errorf("Unexpected error:\nExpected: %v\nActual: %v", test.expectedError, err)
				return
			} else if (err == nil) && !test.expectedOutput.Equal(*res) {
				t.Errorf("Unexpected result:\nExpected: %v\nActual: %v", test.expectedOutput, res)
			}
		})
	}
}

func TestBookDBRepository_GetBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	getBookMocks := []struct {
		testName       string
		expectedOutput entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
		getID          uint64
	}{
		{
			testName: "Test Successful",
			expectedOutput: entity.Book{
				ID:          1,
				Title:       "test title",
				Author:      "test author",
				Year:        1,
				Description: "test description",
			},
			expectedError: nil,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).AddRow(1, "test title", "test author", 1, "test description")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBook)).WithArgs(1).WillReturnRows(rows)
			},
			mockRepo: repo,
			getID:    1,
		},
		{
			testName:       "Test Unsuccessful: Book not found",
			expectedOutput: entity.Book{},
			expectedError:  errors.BookNotFound{},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBook)).WithArgs(2).WillReturnRows(rows)
			},
			mockRepo: repo,
			getID:    2,
		},
	}

	for _, test := range getBookMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.GetBook(test.getID)
			if (err != nil) && (err != test.expectedError) {
				t.Errorf("Unexpected error:\nExpected: %v\nActual: %v", test.expectedError, err)
				return
			} else if (err == nil) && !test.expectedOutput.Equal(*res) {
				t.Errorf("Unexpected result:\nExpected: %v\nActual: %v", test.expectedOutput, res)
			}
		})
	}
}

func TestBookDBRepository_GetAllBooks(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	getAllBooksMocks := []struct {
		testName       string
		expectedOutput []entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
	}{
		{
			testName: "Test Successful",
			expectedOutput: []entity.Book{
				{
					ID:          1,
					Title:       "test title 1",
					Author:      "test author 1",
					Year:        1,
					Description: "test description 1",
				},
				{
					ID:          2,
					Title:       "test title 2",
					Author:      "test author 2",
					Year:        2,
					Description: "test description 2",
				},
				{
					ID:          3,
					Title:       "test title 3",
					Author:      "test author 3",
					Year:        3,
					Description: "test description 3",
				},
			},
			expectedError: nil,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "test title 1", "test author 1", 1, "test description 1").
					AddRow(2, "test title 2", "test author 2", 2, "test description 2").
					AddRow(3, "test title 3", "test author 3", 3, "test description 3")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)
			},
			mockRepo: repo,
		},
		{
			testName: "Test Successful: Row 2 unparsed",
			expectedOutput: []entity.Book{
				{
					ID:          1,
					Title:       "test title 1",
					Author:      "test author 1",
					Year:        1,
					Description: "test description 1",
				},
				{
					ID:          3,
					Title:       "test title 3",
					Author:      "test author 3",
					Year:        3,
					Description: "test description 3",
				},
			},
			expectedError: nil,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "test title 1", "test author 1", 1, "test description 1").
					AddRow(2, "test title 2", "test author 2", "error", "test description 2").
					AddRow(3, "test title 3", "test author 3", 3, "test description 3")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)
			},
			mockRepo: repo,
		},
		{
			testName: "Test Unsuccessful: Books not found",
			expectedOutput: []entity.Book{
				{
					ID:          1,
					Title:       "test title 1",
					Author:      "test author 1",
					Year:        1,
					Description: "test description 1",
				},
				{
					ID:          2,
					Title:       "test title 2",
					Author:      "test author 2",
					Year:        2,
					Description: "test description 2",
				},
				{
					ID:          3,
					Title:       "test title 3",
					Author:      "test author 3",
					Year:        3,
					Description: "test description 3",
				},
			},
			expectedError: errors.BookNotFound{},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)
			},
			mockRepo: repo,
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName: "Test Unsuccessful: DB is closed",
			expectedOutput: []entity.Book{
				{
					ID:          1,
					Title:       "test title 1",
					Author:      "test author 1",
					Year:        1,
					Description: "test description 1",
				},
				{
					ID:          2,
					Title:       "test title 2",
					Author:      "test author 2",
					Year:        2,
					Description: "test description 2",
				},
				{
					ID:          3,
					Title:       "test title 3",
					Author:      "test author 3",
					Year:        3,
					Description: "test description 3",
				},
			},
			expectedError: errors.BookCouldNotQuery{Msg: "sql: database is closed"},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)
				db.Close()
			},
			mockRepo: repo,
		},
	}

	for _, test := range getAllBooksMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.GetAllBooks()
			if (err != nil) && (err != test.expectedError) {
				t.Errorf("Unexpected error:\nExpected: %v\nActual: %v", test.expectedError, err)
				return
			} else if (err == nil) && !entity.BookArrayEqualNoID(test.expectedOutput, res) {
				t.Errorf("Unexpected result:\nExpected: %v\nActual: %v", test.expectedOutput, res)
			}
		})
	}
}

func TestBookDBRepository_SearchByAuthor(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	searchByAuthorMocks := []struct {
		testName       string
		expectedOutput []entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
		author         string
	}{
		{
			testName: "Test Successful",
			expectedOutput: []entity.Book{
				{
					ID:          1,
					Title:       "test title 1",
					Author:      "test author",
					Year:        1,
					Description: "test description 1",
				},
				{
					ID:          2,
					Title:       "test title 2",
					Author:      "test author",
					Year:        2,
					Description: "test description 2",
				},
				{
					ID:          3,
					Title:       "test title 3",
					Author:      "test author",
					Year:        3,
					Description: "test description 3",
				},
				{
					ID:          5,
					Title:       "test title 5",
					Author:      "test author",
					Year:        5,
					Description: "test description 5",
				},
			},
			expectedError: nil,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "test title 1", "test author", 1, "test description 1").
					AddRow(2, "test title 2", "test author", 2, "test description 2").
					AddRow(3, "test title 3", "test author", 3, "test description 3").
					AddRow(5, "test title 5", "test author", 5, "test description 5")
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByAuthorBook)).WithArgs("test author").WillReturnRows(rows)
			},
			mockRepo: repo,
			author:   "test author",
		},
		{
			testName: "Test Successful: Row 4 unparsed",
			expectedOutput: []entity.Book{
				{
					ID:          1,
					Title:       "test title 1",
					Author:      "test author",
					Year:        1,
					Description: "test description 1",
				},
				{
					ID:          2,
					Title:       "test title 2",
					Author:      "test author",
					Year:        2,
					Description: "test description 2",
				},
				{
					ID:          3,
					Title:       "test title 3",
					Author:      "test author",
					Year:        3,
					Description: "test description 3",
				},
				{
					ID:          5,
					Title:       "test title 5",
					Author:      "test author",
					Year:        5,
					Description: "test description 5",
				},
			},
			expectedError: nil,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "test title 1", "test author", 1, "test description 1").
					AddRow(2, "test title 2", "test author", 2, "test description 2").
					AddRow(3, "test title 3", "test author", 3, "test description 3").
					AddRow(4, "test title 5", "test author", "error", "test description 4").
					AddRow(5, "test title 5", "test author", 5, "test description 5")
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByAuthorBook)).WithArgs("test author").WillReturnRows(rows)
			},
			mockRepo: repo,
			author:   "test author",
		},
		{
			testName:      "Test Unsuccessful: Invalid author",
			expectedError: errors.BookBadRequest{},
			author:        "",
		},
		{
			testName:      "Test Unsuccessful: Books not found",
			expectedError: errors.BookNotFoundByAuthor{Author: "test author"},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByAuthorBook)).WithArgs("test author").WillReturnRows(rows)
			},
			mockRepo: repo,
			author:   "test author",
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.BookCouldNotQuery{Msg: "sql: database is closed"},
			mockFunc: func() {
				db.Close()
			},
			mockRepo: repo,
			author:   "test author",
		},
	}

	for _, test := range searchByAuthorMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.SearchByAuthor(test.author)
			if (err != nil) && (err != test.expectedError) {
				t.Errorf("Unexpected error:\nExpected: %v\nActual: %v", test.expectedError, err)
				return
			} else if (err == nil) && !entity.BookArrayEqualNoID(test.expectedOutput, res) {
				t.Errorf("Unexpected result:\nExpected: %v\nActual: %v", test.expectedOutput, res)
			}
		})
	}
}

func TestBookDBRepository_SearchByTitle(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	searchByAuthorMocks := []struct {
		testName       string
		expectedOutput entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
		title          string
	}{
		{
			testName: "Test Successful",
			expectedOutput: entity.Book{
				ID:          1,
				Title:       "test title",
				Author:      "test author",
				Year:        1,
				Description: "test description 1",
			},
			expectedError: nil,
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "test title", "test author", 1, "test description 1")
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByTitleBook)).WithArgs("test title").WillReturnRows(rows)
			},
			mockRepo: repo,
			title:    "test title",
		},
		{
			testName:      "Test Unsuccessful: Invalid title",
			expectedError: errors.BookBadRequest{},
			title:         "",
		},
		{
			testName:      "Test Unsuccessful: No books found",
			expectedError: errors.BookNotFoundByTitle{Title: "test title"},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByTitleBook)).WithArgs("test title").WillReturnRows(rows)
			},
			mockRepo: repo,
			title:    "test title",
		},
	}

	for _, test := range searchByAuthorMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.SearchByTitle(test.title)
			if (err != nil) && (err != test.expectedError) {
				t.Errorf("Unexpected error:\nExpected: %v\nActual: %v", test.expectedError, err)
				return
			} else if (err == nil) && !test.expectedOutput.Equal(*res) {
				t.Errorf("Unexpected result:\nExpected: %v\nActual: %v", test.expectedOutput, res)
			}
		})
	}
}

func TestBookDBRepository_UpdateBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	searchByAuthorMocks := []struct {
		testName       string
		input          entity.Book
		expectedOutput entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
		id             uint64
	}{
		{
			testName: "Test Successful",
			input: entity.Book{
				Title:       "test title",
				Author:      "test author",
				Year:        1,
				Description: "test description 1",
			},
			expectedOutput: entity.Book{
				ID:          3,
				Title:       "test title",
				Author:      "test author",
				Year:        1,
				Description: "test description 1",
			},
			expectedError: nil,
			mockFunc: func() {
				rowsFirstQuery := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(3, "test title", "test author", 1, "test description 1")
				rowsSecondQuery := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(3, "test title", "test author", 1, "test description 1")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBook)).WithArgs(3).WillReturnRows(rowsFirstQuery)
				mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBook)).WithArgs(3, "test title", "test author", 1, "test description 1").WillReturnRows(rowsSecondQuery)
			},
			mockRepo: repo,
			id:       3,
		},
	}

	for _, test := range searchByAuthorMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.UpdateBook(test.id, &test.input)
			if (err != nil) && (err != test.expectedError) {
				t.Errorf("Unexpected error:\nExpected: %v\nActual: %v", test.expectedError, err)
				return
			} else if (err == nil) && !test.expectedOutput.Equal(*res) {
				t.Errorf("Unexpected result:\nExpected: %v\nActual: %v", test.expectedOutput, res)
			}
		})
	}
}

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
