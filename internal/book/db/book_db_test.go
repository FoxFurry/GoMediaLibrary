package db

import (
	"database/sql"
	goerrors "errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/foxfurry/medialib/internal/book/domain/entity"
	"github.com/foxfurry/medialib/internal/book/http/errors"
	"github.com/foxfurry/medialib/internal/book/http/validators"
	ct "github.com/foxfurry/medialib/internal/common/server/common_translators"
	_ "github.com/foxfurry/medialib/internal/common/tests"
	"github.com/stretchr/testify/assert"
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
				mock.ExpectQuery(regexp.QuoteMeta(QuerySaveBook)).WithArgs("test title", "test author", 1, "test description").WillReturnRows(rows)
			},
			mockRepo: repo,
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
			expectedError:  errors.NewBookCouldNotQuery("sql: no rows in result set"),
			mockFunc: func() {
				rows := mock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(QuerySaveBook)).WithArgs("test title", "test author", 1, "test description").WillReturnRows(rows)
			},
			mockRepo: repo,
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
			mockFunc: func() {
				db.Close()
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
				mock.ExpectQuery(regexp.QuoteMeta(QueryGetBook)).WithArgs(1).WillReturnRows(rows)
			},
			mockRepo: repo,
			getID:    1,
		},
		{
			testName:       "Test Unsuccessful: Book not found",
			expectedOutput: entity.Book{},
			expectedError:  errors.NewBooksNotFound(),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(QueryGetBook)).WithArgs(2).WillReturnRows(rows)
			},
			mockRepo: repo,
			getID:    2,
		},
		{
			testName:       "Test Unsuccessful: Invalid serial",
			expectedOutput: entity.Book{},
			expectedError:  errors.NewBookInvalidSerial(),
			mockRepo: repo,
			getID:    0,
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
			mockFunc: func() {
				db.Close()
			},
			mockRepo: repo,
			getID:    1,
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
				mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnRows(rows)
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
				mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnRows(rows)
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
			expectedError: errors.NewBooksNotFound(),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(QueryGetAll)).WillReturnRows(rows)
			},
			mockRepo: repo,
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
			mockFunc: func() {
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
				mock.ExpectQuery(regexp.QuoteMeta(QuerySearchByAuthorBook)).WithArgs("test author").WillReturnRows(rows)
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
				mock.ExpectQuery(regexp.QuoteMeta(QuerySearchByAuthorBook)).WithArgs("test author").WillReturnRows(rows)
			},
			mockRepo: repo,
			author:   "test author",
		},
		{
			testName:      "Test Unsuccessful: Invalid author",
			expectedError: errors.NewBookValidatorError([]ct.FieldError{validators.FieldAuthorEmpty}),
			author:        "",
		},
		{
			testName:      "Test Unsuccessful: Books not found",
			expectedError: errors.NewBookNotFoundByAuthor("test author"),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(QuerySearchByAuthorBook)).WithArgs("test author").WillReturnRows(rows)
			},
			mockRepo: repo,
			author:   "test author",
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
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

			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err, "Values are not equal:\nExpected: %+v\nActual: %+v", test.expectedError, err)
			}
			if test.expectedOutput != nil {
				assert.True(t, entity.BookArrayEqualNoID(test.expectedOutput, res))
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
		expectedOutput *entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
		title          string
	}{
		{
			testName: "Test Successful",
			expectedOutput: &entity.Book{
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
				mock.ExpectQuery(regexp.QuoteMeta(QuerySearchByTitleBook)).WithArgs("test title").WillReturnRows(rows)
			},
			mockRepo: repo,
			title:    "test title",
		},
		{
			testName:      "Test Unsuccessful: Invalid title",
			expectedError: errors.NewBookValidatorError([]ct.FieldError{validators.FieldTitleEmpty}),
			title:         "",
		},
		{
			testName:      "Test Unsuccessful: No books found",
			expectedError: errors.NewBookNotFoundByTitle("test title"),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(QuerySearchByTitleBook)).WithArgs("test title").WillReturnRows(rows)
			},
			mockRepo: repo,
			title:    "test title",
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
			mockFunc: func() {
				db.Close()
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

			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err, "Values are not equal:\nExpected: %+v\nActual: %+v", test.expectedError, err)
			}
			if test.expectedOutput != nil {
				assert.True(t, test.expectedOutput.Equal(*res))
			}
		})
	}
}

func TestBookDBRepository_UpdateBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	updateBookMocks := []struct {
		testName       string
		input          *entity.Book
		expectedOutput *entity.Book
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
		id             uint64
	}{
		{
			testName: "Test Successful",
			input: &entity.Book{
				Title:       "test title 2",
				Author:      "test author 2",
				Year:        2,
				Description: "test description 2",
			},
			expectedOutput: &entity.Book{
				ID:          3,
				Title:       "test title 2",
				Author:      "test author 2",
				Year:        2,
				Description: "test description 2",
			},
			expectedError: nil,
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(QueryUpdateBook)).WithArgs(3, "test title 2", "test author 2", 2, "test description 2").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			mockRepo: repo,
			id:       3,
		},
		{
			testName:       "Test Unsuccessful: Invalid serial",
			expectedError:  errors.NewBookInvalidSerial(),
			mockRepo: repo,
			id:    0,
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			input: &entity.Book{},
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
			mockFunc: func() {
				db.Close()
			},
			mockRepo: repo,
			id:       1,
		},
	}

	for _, test := range updateBookMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.UpdateBook(test.id, test.input)

			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err, "Values are not equal:\nExpected: %+v\nActual: %+v", test.expectedError, err)
			}
			if test.expectedOutput != nil {
				assert.True(t, test.expectedOutput.Equal(*res))
			}
		})
	}
}

func TestBookDBRepository_DeleteBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	deleteBookMocks := []struct {
		testName       string
		expectedOutput int64
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
		id             uint64
	}{
		{
			testName:       "Test Successful",
			expectedOutput: 1,
			expectedError:  nil,
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(QueryDeleteBook)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			mockRepo: repo,
			id:       1,
		},
		{
			testName:      "Test Unsuccessful: Book not found",
			expectedError: errors.NewBooksNotFound(),
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(QueryDeleteBook)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			mockRepo: repo,
			id:       1,
		},
		{
			testName:      "Test Unsuccessful: Invalid rows affected",
			expectedError: errors.NewBookCouldNotQuery("no RowsAffected available after DDL statement"),
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(QueryDeleteBook)).WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(goerrors.New("no RowsAffected available after DDL statement")))
			},
			mockRepo: repo,
			id:       1,
		},
		{
			testName:       "Test Unsuccessful: Invalid serial",
			expectedError:  errors.NewBookInvalidSerial(),
			mockRepo: repo,
			id:    0,
		},
		{ // DO NOT ADD ANY TC AFTER THIS ONE. IN THIS TC DB IS BEING CLOSED AND NOT REOPENED FOR REST OF TEST
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
			mockFunc: func() {
				db.Close()
			},
			mockRepo: repo,
			id:       1,
		},
	}

	for _, test := range deleteBookMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.DeleteBook(test.id)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err, "Values are not equal:\nExpected: %+v\nActual: %+v", test.expectedError, err)
			}
			assert.Equal(t, test.expectedOutput, res)
		})
	}
}

func TestBookDBRepository_DeleteAllBooks(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookRepo(db)

	updateBookMocks := []struct {
		testName       string
		expectedOutput int64
		expectedError  error
		mockFunc       func()
		mockRepo       BookDBRepository
	}{
		{
			testName:       "Test Successful",
			expectedOutput: 4,
			expectedError:  nil,
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(QueryDeleteAllBooksAndAlter)).WillReturnResult(sqlmock.NewResult(0, 4))
			},
			mockRepo: repo,
		},
		{
			testName:       "Test Unsuccessful: Invalid rows affected",
			expectedOutput: 0,
			expectedError:  errors.NewBookCouldNotQuery("no RowsAffected available after DDL statement"),
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(QueryDeleteAllBooksAndAlter)).WillReturnResult(sqlmock.NewErrorResult(goerrors.New("no RowsAffected available after DDL statement")))
			},
			mockRepo: repo,
		},
		{
			testName:       "Test Unsuccessful: Books not found",
			expectedOutput: 0,
			expectedError:  errors.NewBooksNotFound(),
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(QueryDeleteAllBooksAndAlter)).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			mockRepo: repo,
		},
		{
			testName:      "Test Unsuccessful: DB is closed",
			expectedError: errors.NewBookCouldNotQuery("sql: database is closed"),
			mockFunc: func() {
				db.Close()
			},
			mockRepo: repo,
		},
	}

	for _, test := range updateBookMocks {
		t.Run(test.testName, func(t *testing.T) {
			if test.mockFunc != nil {
				test.mockFunc()
			}

			res, err := test.mockRepo.DeleteAllBooks()
			if (err != nil) && (err != test.expectedError) {
				t.Errorf("Unexpected error:\nExpected: %v\nActual: %v", test.expectedError, err)
				return
			} else if (err == nil) && res != test.expectedOutput {
				t.Errorf("Unexpected result:\nExpected: %v\nActual: %v", test.expectedOutput, res)
			}
		})
	}
}
