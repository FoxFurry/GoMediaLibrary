package controllers

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/foxfurry/simple-rest/internal/book/http/validators"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

const (
	querySaveBook               = `INSERT INTO bookstore (title, author, year, description) VALUES ($1, $2, $3, $4) RETURNING id`
	queryGetBook                = `SELECT * FROM bookstore WHERE id=$1`
	queryGetAll                 = `SELECT * FROM bookstore`
	querySearchByAuthorBook     = `SELECT * FROM bookstore WHERE author=$1`
	querySearchByTitleBook      = `SELECT * FROM bookstore WHERE title=$1`
	queryUpdateBook             = `UPDATE bookstore SET title=$2, author=$3, year=$4, description=$5 WHERE id=$1`
	queryDeleteBook             = `DELETE FROM bookstore WHERE id=$1`
	queryDeleteAllBooksAndAlter = `DELETE FROM bookstore; ALTER SEQUENCE bookstore_id_seq RESTART WITH 1`
)

func init(){
	validators.RegisterBookValidators()
}

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Could not create a new mock: %v", err)
	}

	return db, mock
}

func TestBookService_SaveBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	saveUrl := "/book"
	saveMethod := "POST"

	saveBookServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		requestBody    string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(querySaveBook)).WithArgs("Test 1", "Test 1", 1, "Test 1").WillReturnRows(rows)
			},
			service:        repo,
			requestBody:    "{\"title\":\"Test 1\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"}",
			url:            "/book",
			method:         "POST",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1,\"title\":\"Test 1\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"}",
		},
		{
			testName:       "Test Unsuccessful: Empty title",
			service:        repo,
			requestBody:    "{\"title\":\"\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"}",
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid request body: Title cannot be empty\"}",
		},
		{
			testName:       "Test Unsuccessful: Empty request body",
			service:        repo,
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Expected body, found EOF\"}",
		},
		{
			testName:       "Test Unsuccessful: DB is closed",
			mockFunc: func() {
				db.Close()
			},
			service:        repo,
			requestBody:    "{\"title\":\"Test 1\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"}",
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error:\":\"Could not execute query: sql: database is closed\"}",
		},
	}

	for _, tc := range saveBookServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.requestBody))

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.SaveBook(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestBookService_GetBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	getMethod := "GET"
	getURL := "/book"

	getBookServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		method         string
		url            string
		params         []gin.Param
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).AddRow(1, "Test 1", "Test 1", 1, "Test 1")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBook)).WithArgs(1).WillReturnRows(rows)
			},
			service: repo,
			url:     getURL,
			method:  getMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1,\"title\":\"Test 1\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"}",
		},
		{
			testName: "Test Unsuccessful: Invalid param value",
			service:  repo,
			url:      getURL,
			method:   getMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "0",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid serial. Serial must be more than 1\"}",
		},
		{
			testName:       "Test Unsuccessful: Empty param",
			service:        repo,
			url:            getURL,
			method:         getMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid serial. Serial must be more than 1\"}",
		},
		{
			testName: "Test Unsuccessful: Book not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBook)).WithArgs(666).WillReturnRows(rows)
			},
			service: repo,
			url:     getURL,
			method:  getMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "666",
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"error:\":\"Book(s) not found in db\"}",
		},
		{
			testName: "Test Unsuccessful: DB is closed",
			mockFunc: func() {
				db.Close()
			},
			service: repo,
			url:     getURL,
			method:  getMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "69",
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error:\":\"Could not execute query: sql: database is closed\"}",
		},
	}

	for _, tc := range getBookServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = tc.params

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.GetBook(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestBookService_GetAllBooks(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	getAllMethod := "GET"
	getAllURL := "/book"

	getAllBooksServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "Test 1", "Test 1", 1, "Test 1").
					AddRow(2, "Test 2", "Test 2", 2, "Test 2").
					AddRow(3, "Test 3", "Test 3", 3, "Test 3")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)
			},
			service:        repo,
			url:            getAllURL,
			method:         getAllMethod,
			expectedStatus: http.StatusOK,
			expectedBody:   "[{\"id\":1,\"title\":\"Test 1\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"},{\"id\":2,\"title\":\"Test 2\",\"author\":\"Test 2\",\"year\":2,\"description\":\"Test 2\"},{\"id\":3,\"title\":\"Test 3\",\"author\":\"Test 3\",\"year\":3,\"description\":\"Test 3\"}]",
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAll)).WillReturnRows(rows)
			},
			service:        repo,
			url:            getAllURL,
			method:         getAllMethod,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"error:\":\"Book(s) not found in db\"}",
		},
		{
			testName: "Test Unsuccessful: DB is closed",
			mockFunc: func() {
				db.Close()
			},
			service:        repo,
			url:            getAllURL,
			method:         getAllMethod,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error:\":\"Could not execute query: sql: database is closed\"}",
		},
	}

	for _, tc := range getAllBooksServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.GetAllBooks(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestBookService_SearchByAuthor(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	searchAuthorMethod := "GET"
	searchAuthorURL := "/book/author"

	searchByAuthorServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		method         string
		params         []gin.Param
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "Test 1", "Test", 1, "Test 1").
					AddRow(2, "Test 2", "Test", 2, "Test 2").
					AddRow(3, "Test 3", "Test", 3, "Test 3")
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByAuthorBook)).WithArgs("Test").WillReturnRows(rows)
			},
			service: repo,
			url:     searchAuthorURL,
			method:  searchAuthorMethod,
			params: []gin.Param{
				{
					Key:   "author",
					Value: "Test",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "[{\"id\":1,\"title\":\"Test 1\",\"author\":\"Test\",\"year\":1,\"description\":\"Test 1\"},{\"id\":2,\"title\":\"Test 2\",\"author\":\"Test\",\"year\":2,\"description\":\"Test 2\"},{\"id\":3,\"title\":\"Test 3\",\"author\":\"Test\",\"year\":3,\"description\":\"Test 3\"}]",
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByAuthorBook)).WithArgs("Test").WillReturnRows(rows)
			},
			service: repo,
			url:     searchAuthorURL,
			method:  searchAuthorMethod,
			params: []gin.Param{
				{
					Key:   "author",
					Value: "Test",
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"error:\":\"Book(s) with author Test not found in db\"}",
		},
		{
			testName:       "Test Unsuccessful: Invalid author",
			service:        repo,
			url:            searchAuthorURL,
			method:         searchAuthorMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid request body: Author cannot be empty\"}",
		},
		{
			testName: "Test Unsuccessful: DB is closed",
			service:  repo,
			mockFunc: func() {
				db.Close()
			},
			url:    searchAuthorURL,
			method: searchAuthorMethod,
			params: []gin.Param{
				{
					Key:   "author",
					Value: "Test 1",
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error:\":\"Could not execute query: sql: database is closed\"}",
		},
	}

	for _, tc := range searchByAuthorServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = tc.params
			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.SearchByAuthor(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestBookService_SearchByTitle(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	searchTitleMethod := "GET"
	searchTitleURL := "/book/title"

	searchByTitleServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		method         string
		params         []gin.Param
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "Test 1", "Test 1", 1, "Test 1")
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByTitleBook)).WithArgs("Test 1").WillReturnRows(rows)
			},
			service: repo,
			url:     searchTitleURL,
			method:  searchTitleMethod,
			params: []gin.Param{
				{
					Key:   "title",
					Value: "Test 1",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1,\"title\":\"Test 1\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"}",
		},
		{
			testName: "Test Unsuccessful: Book not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(querySearchByTitleBook)).WithArgs("Test 1").WillReturnRows(rows)
			},
			service: repo,
			url:     searchTitleURL,
			method:  searchTitleMethod,
			params: []gin.Param{
				{
					Key:   "title",
					Value: "Test 1",
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"error:\":\"Book(s) with title Test 1 not found in db\"}",
		},
		{
			testName: "Test Unsuccessful: Invalid title",
			service:  repo,
			url:      searchTitleURL,
			method:   searchTitleMethod,
			params: []gin.Param{
				{
					Key:   "title",
					Value: "",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid request body: Title cannot be empty\"}",
		},
		{
			testName: "Test Unsuccessful: DB is closed",
			service:  repo,
			mockFunc: func() {
				db.Close()
			},
			url:    searchTitleURL,
			method: searchTitleMethod,
			params: []gin.Param{
				{
					Key:   "title",
					Value: "Test 1",
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error:\":\"Could not execute query: sql: database is closed\"}",
		},
	}

	for _, tc := range searchByTitleServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = tc.params
			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.SearchByTitle(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	updateMethod := "PUT"
	updateURL := "/book"

	UpdateBookServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		requestBody    string
		method         string
		params         []gin.Param
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBook)).WithArgs(1, "Test 2", "Test 2", 2, "Test 2").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			requestBody: "{\"title\":\"Test 2\",\"author\":\"Test 2\",\"year\":2,\"description\":\"Test 2\"}",
			service:     repo,
			url:         updateURL,
			method:      updateMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1,\"title\":\"Test 2\",\"author\":\"Test 2\",\"year\":2,\"description\":\"Test 2\"}",
		},
		{
			testName:    "Test Unsuccessful: Invalid serial",
			requestBody: "{\"title\":\"Test 2\",\"author\":\"Test 2\",\"year\":2,\"description\":\"Test 2\"}",
			service:     repo,
			url:         updateURL,
			method:      updateMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "0",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid serial. Serial must be more than 1\"}",
		},
		{
			testName:       "Test Unsuccessful: Empty serial",
			service:        repo,
			url:            updateURL,
			method:         updateMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid serial. Serial must be more than 1\"}",
		},
		{
			testName:    "Test Unsuccessful: Invalid request body",
			requestBody: "{\"title\":\"\",\"author\":\"Test 2\",\"year\":2,\"description\":\"Test 2\"}",
			service:     repo,
			url:         updateURL,
			method:      updateMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid request body: Title cannot be empty\"}",
		},
		{
			testName: "Test Unsuccessful: Empty request body",
			service:  repo,
			url:      updateURL,
			method:   updateMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Expected body, found EOF\"}",
		},
	}

	for _, tc := range UpdateBookServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.requestBody))
			fmt.Println(req)
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = tc.params

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.UpdateBook(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestBookService_DeleteBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	deleteMethod := "DELETE"
	deleteUrl := "/book"

	deleteBookServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		method         string
		url            string
		params         []gin.Param
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteBook)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			service: repo,
			url:     deleteUrl,
			method:  deleteMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteBook)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			service: repo,
			url:     deleteUrl,
			method:  deleteMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"error:\":\"Book(s) not found in db\"}",
		},
		{
			testName: "Test Unsuccessful: Invalid serial",
			service:  repo,
			url:      deleteUrl,
			method:   deleteMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "0",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid serial. Serial must be more than 1\"}",
		},
		{
			testName:       "Test Unsuccessful: Empty serial",
			service:        repo,
			url:            deleteUrl,
			method:         deleteMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"error:\":\"Invalid serial. Serial must be more than 1\"}",
		},
	}

	for _, tc := range deleteBookServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = tc.params

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.DeleteBook(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestBookService_DeleteAllBooks(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookService(db)

	deleteAllMethod := "DELETE"
	deleteAllURL := "/book"

	deleteAllBooksServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteAllBooksAndAlter)).WillReturnResult(sqlmock.NewResult(0, 4))
			},
			service:        repo,
			url:            deleteAllURL,
			method:         deleteAllMethod,
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"Deleted rows\":4}",
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteAllBooksAndAlter)).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			service:        repo,
			url:            deleteAllURL,
			method:         deleteAllMethod,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"error:\":\"Book(s) not found in db\"}",
		},
		{
			testName: "Test Unsuccessful: DB is closed",
			mockFunc: func() {
				db.Close()
			},
			service:        repo,
			url:            deleteAllURL,
			method:         deleteAllMethod,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error:\":\"Could not execute query: sql: database is closed\"}",
		},
	}

	for _, tc := range deleteAllBooksServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.DeleteAllBooks(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

/*
Actual, expected in response
 */