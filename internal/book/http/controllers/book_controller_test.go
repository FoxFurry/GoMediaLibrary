package controllers

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
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
			testName: "Test Unsuccessful: Invalid request body",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(querySaveBook)).WithArgs("Test 1", "Test 1", 1, "Test 1").WillReturnRows(rows)
			},
			service:        repo,
			requestBody:    "{\"title\":\"\",\"author\":\"Test 1\",\"year\":1,\"description\":\"Test 1\"}",
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Title, author and year are mandatory fields",
		},
		{
			testName: "Test Unsuccessful: Empty request body",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(querySaveBook)).WithArgs("Test 1", "Test 1", 1, "Test 1").WillReturnRows(rows)
			},
			service:        repo,
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Title, author and year are mandatory fields",
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
			url:     "/book",
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
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).AddRow(1, "Test 1", "Test 1", 1, "Test 1")
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBook)).WithArgs(1).WillReturnRows(rows)
			},
			service: repo,
			url:     "/book",
			method:  getMethod,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "0",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid serial. Serial must be more than 1",
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
