package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	bookdb "github.com/foxfurry/medialib/internal/book/db"
	"github.com/foxfurry/medialib/internal/book/domain/entity"
	"github.com/foxfurry/medialib/internal/book/http/errors"
	"github.com/foxfurry/medialib/internal/book/http/validators"
	"github.com/foxfurry/medialib/internal/common/server/translator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

type expectedErrors struct {
	Msg    string                  `json:"msg,omitempty"`
	Fields []translator.FieldError `json:"fields,omitempty"`
}

type singleResponse struct {
	Data  *entity.Book   `json:"data"`
	Error expectedErrors `json:"error"`
}

type rowsResponse struct {
	Data  int   `json:"data"`
	Error expectedErrors `json:"error"`
}

type arrayResponse struct {
	Data  []entity.Book  `json:"data"`
	Error expectedErrors `json:"error"`
}

func (e expectedErrors) isEmpty() bool {
	return e.Msg == "" && e.Fields == nil
}

func init() {
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

	repo := NewBookController(db)

	saveUrl := "/book"
	saveMethod := "POST"

	saveBookServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		requestBody    *entity.Book
		method         string
		url            string
		expectedStatus int
		expectedBody   *entity.Book
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QuerySaveBook)).WithArgs("Test 1", "Test 1", 1, "Test 1").WillReturnRows(rows)
			},
			service: repo,
			requestBody: &entity.Book{
				Title:       "Test 1",
				Author:      "Test 1",
				Year:        1,
				Description: "Test 1",
			},
			url:            "/book",
			method:         "POST",
			expectedStatus: http.StatusOK,
			expectedBody: &entity.Book{
				Title:       "Test 1",
				Author:      "Test 1",
				Year:        1,
				Description: "Test 1",
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Empty title",
			service:  repo,
			requestBody: &entity.Book{
				Author:      "Ray Bradbury",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusBadRequest,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldTitleEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Empty request body",
			service:        repo,
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusBadRequest,
			expectedError: expectedErrors{
				Msg: errors.NewBookEmptyBody().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: DB is closed",
			mockFunc: func() {
				db.Close()
			},
			service: repo,
			requestBody: &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			url:            saveUrl,
			method:         saveMethod,
			expectedStatus: http.StatusInternalServerError,
			expectedError: expectedErrors{
				Msg: errors.NewBookCouldNotQuery("sql: database is closed").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range saveBookServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()

			var jsonRequest []byte = nil

			if tc.requestBody != nil {
				jsonRequest, _ = json.Marshal(tc.requestBody)
			}

			req, _ := http.NewRequest(tc.method, tc.url, bytes.NewReader(jsonRequest))

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.SaveBook(c)

			assert.Equal(t, tc.expectedStatus, w.Code)

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := singleResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			if tc.expectedBody != nil {
				assert.True(t, tc.expectedBody.EqualNoID(*resultBody.Data), "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedBody, resultBody.Data)
			} else if resultBody.Data != nil {
				t.Errorf("Expected result body to be nil, found %+v", resultBody.Data)
			}
			if !tc.expectedError.isEmpty() {
				assert.Equal(t, tc.expectedError, resultBody.Error, "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedError, resultBody.Error)
			} else if !resultBody.Error.isEmpty() {
				t.Errorf("Expected error to be nil, found %+v", resultBody.Error)
			}

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestBookService_GetBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookController(db)

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
		expectedBody   *entity.Book
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).AddRow(1, "Test 1", "Test 1", 1, "Test 1")
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QueryGetBook)).WithArgs(1).WillReturnRows(rows)
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
			expectedBody: &entity.Book{
				Title:       "Test 1",
				Author:      "Test 1",
				Year:        1,
				Description: "Test 1",
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg: errors.NewBookInvalidSerial().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Empty param",
			service:        repo,
			url:            getURL,
			method:         getMethod,
			expectedStatus: http.StatusBadRequest,
			expectedError: expectedErrors{
				Msg: errors.NewBookInvalidSerial().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Book not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QueryGetBook)).WithArgs(666).WillReturnRows(rows)
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
			expectedError: expectedErrors{
				Msg: errors.NewBooksNotFound().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg: errors.NewBookCouldNotQuery("sql: database is closed").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := singleResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			if tc.expectedBody != nil {
				assert.True(t, tc.expectedBody.EqualNoID(*resultBody.Data), "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedBody, resultBody.Data)
			} else if resultBody.Data != nil {
				t.Errorf("Expected result body to be nil, found %+v", resultBody.Data)
			}
			if !tc.expectedError.isEmpty() {
				assert.Equal(t, tc.expectedError, resultBody.Error, "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedError, resultBody.Error)
			} else if !resultBody.Error.isEmpty() {
				t.Errorf("Expected error to be nil, found %+v", resultBody.Error)
			}

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestBookService_GetAllBooks(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookController(db)

	getAllMethod := "GET"
	getAllURL := "/book"

	getAllBooksServiceMocks := []struct {
		testName          string
		mockFunc          func()
		service           BookService
		method            string
		url               string
		expectedStatus    int
		expectedBodyArray []entity.Book
		expectedError     expectedErrors
		expectedHeader    map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "Test 1", "Test 1", 1, "Test 1").
					AddRow(2, "Test 2", "Test 2", 2, "Test 2").
					AddRow(3, "Test 3", "Test 3", 3, "Test 3")
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QueryGetAll)).WillReturnRows(rows)
			},
			service:        repo,
			url:            getAllURL,
			method:         getAllMethod,
			expectedStatus: http.StatusOK,
			expectedBodyArray: []entity.Book{
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
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QueryGetAll)).WillReturnRows(rows)
			},
			service:           repo,
			url:               getAllURL,
			method:            getAllMethod,
			expectedStatus:    http.StatusNotFound,
			expectedBodyArray: []entity.Book{},
			expectedError: expectedErrors{
				Msg: errors.NewBooksNotFound().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg: errors.NewBookCouldNotQuery("sql: database is closed").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := &arrayResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			assert.Equal(t, tc.expectedStatus, w.Code)

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}

			if tc.expectedBodyArray != nil {
				assert.True(t, entity.BookArrayEqualNoID(tc.expectedBodyArray, resultBody.Data), "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedBodyArray, resultBody.Data)
			} else if resultBody.Data != nil {
				t.Errorf("Expected result body to be nil, found %+v", resultBody.Data)
			}
			if !tc.expectedError.isEmpty() {
				assert.Equal(t, tc.expectedError, resultBody.Error, "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedError, resultBody.Error)
			} else if !resultBody.Error.isEmpty() {
				t.Errorf("Expected error to be nil, found %+v", resultBody.Error)
			}
		})
	}
}

func TestBookService_SearchByAuthor(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookController(db)

	searchAuthorMethod := "GET"
	searchAuthorURL := "/book/author"

	searchByAuthorServiceMocks := []struct {
		testName          string
		mockFunc          func()
		service           BookService
		method            string
		params            []gin.Param
		url               string
		expectedStatus    int
		expectedBodyArray []entity.Book
		expectedError     expectedErrors
		expectedHeader    map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "Test 1", "Test", 1, "Test 1").
					AddRow(2, "Test 2", "Test", 2, "Test 2").
					AddRow(3, "Test 3", "Test", 3, "Test 3")
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QuerySearchByAuthorBook)).WithArgs("Test").WillReturnRows(rows)
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
			expectedBodyArray: []entity.Book{
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
					Description: "Test 2",
				},
				{
					Title:       "Test 3",
					Author:      "Test",
					Year:        3,
					Description: "Test 3",
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QuerySearchByAuthorBook)).WithArgs("Test").WillReturnRows(rows)
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
			expectedError: expectedErrors{
				Msg: errors.NewBookNotFoundByAuthor("Test").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Invalid author",
			service:        repo,
			url:            searchAuthorURL,
			method:         searchAuthorMethod,
			expectedStatus: http.StatusBadRequest,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldAuthorEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg: errors.NewBookCouldNotQuery("sql: database is closed").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := &arrayResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			assert.Equal(t, tc.expectedStatus, w.Code)

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}

			if tc.expectedBodyArray != nil {
				assert.True(t, entity.BookArrayEqualNoID(tc.expectedBodyArray, resultBody.Data), "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedBodyArray, resultBody.Data)
			} else if resultBody.Data != nil {
				t.Errorf("Expected result body to be nil, found %+v", resultBody.Data)
			}
			if !tc.expectedError.isEmpty() {
				assert.Equal(t, tc.expectedError, resultBody.Error, "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedError, resultBody.Error)
			} else if !resultBody.Error.isEmpty() {
				t.Errorf("Expected error to be nil, found %+v", resultBody.Error)
			}
		})
	}
}

func TestBookService_SearchByTitle(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookController(db)

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
		expectedBody   *entity.Book
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"}).
					AddRow(1, "Test 1", "Test 1", 1, "Test 1")
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QuerySearchByTitleBook)).WithArgs("Test 1").WillReturnRows(rows)
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
			expectedBody:   &entity.Book{
				Title:       "Test 1",
				Author:      "Test 1",
				Year:        1,
				Description: "Test 1",
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Book not found",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "description"})
				mock.ExpectQuery(regexp.QuoteMeta(bookdb.QuerySearchByTitleBook)).WithArgs("Test 1").WillReturnRows(rows)
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
			expectedError: expectedErrors{
				Msg:    errors.NewBookNotFoundByTitle("Test 1").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldTitleEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg:    errors.NewBookCouldNotQuery("sql: database is closed").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := singleResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			if tc.expectedBody != nil {
				assert.True(t, tc.expectedBody.EqualNoID(*resultBody.Data), "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedBody, resultBody.Data)
			} else if resultBody.Data != nil {
				t.Errorf("Expected result body to be nil, found %+v", resultBody.Data)
			}
			if !tc.expectedError.isEmpty() {
				assert.Equal(t, tc.expectedError, resultBody.Error, "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedError, resultBody.Error)
			} else if !resultBody.Error.isEmpty() {
				t.Errorf("Expected error to be nil, found %+v", resultBody.Error)
			}

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookController(db)

	updateMethod := "PUT"
	updateURL := "/book"

	UpdateBookServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		requestBody    *entity.Book
		method         string
		params         []gin.Param
		url            string
		expectedStatus int
		expectedBody   *entity.Book
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(bookdb.QueryUpdateBook)).WithArgs(1, "Test 2", "Test 2", 2, "Test 2").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			requestBody: &entity.Book{
				Title:       "Test 2",
				Author:      "Test 2",
				Year:        2,
				Description: "Test 2",
			},
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
			expectedBody:   &entity.Book{
				Title:       "Test 2",
				Author:      "Test 2",
				Year:        2,
				Description: "Test 2",
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:    "Test Unsuccessful: Invalid serial",
			requestBody: &entity.Book{
				Title:       "Test 2",
				Author:      "Test 2",
				Year:        2,
				Description: "Test 2",
			},
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
			expectedError: expectedErrors{
				Msg:    errors.NewBookInvalidSerial().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Empty serial",
			service:        repo,
			url:            updateURL,
			method:         updateMethod,
			expectedStatus: http.StatusBadRequest,
			expectedError: expectedErrors{
				Msg:   errors.NewBookInvalidSerial().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:    "Test Unsuccessful: Invalid request body",
			requestBody: &entity.Book{
				Author:      "Test 2",
				Description: "Test 2",
			},
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
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldTitleEmpty,
					validators.FieldYearEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg:   errors.NewBookEmptyBody().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range UpdateBookServiceMocks {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()

			var jsonRequest []byte = nil

			if tc.requestBody != nil {
				jsonRequest, _ = json.Marshal(tc.requestBody)
			}

			req, _ := http.NewRequest(tc.method, tc.url, bytes.NewReader(jsonRequest))
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = tc.params

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			tc.service.UpdateBook(c)

			assert.Equal(t, tc.expectedStatus, w.Code)

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := singleResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			if tc.expectedBody != nil {
				assert.True(t, tc.expectedBody.EqualNoID(*resultBody.Data), "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedBody, resultBody.Data)
			} else if resultBody.Data != nil {
				t.Errorf("Expected result body to be nil, found %+v", resultBody.Data)
			}
			if !tc.expectedError.isEmpty() {
				assert.Equal(t, tc.expectedError, resultBody.Error, "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedError, resultBody.Error)
			} else if !resultBody.Error.isEmpty() {
				t.Errorf("Expected error to be nil, found %+v", resultBody.Error)
			}

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestBookService_DeleteBook(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookController(db)

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
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(bookdb.QueryDeleteBook)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
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
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(bookdb.QueryDeleteBook)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
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
			expectedError: expectedErrors{
				Msg:    errors.NewBooksNotFound().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg:    errors.NewBookInvalidSerial().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Empty serial",
			service:        repo,
			url:            deleteUrl,
			method:         deleteMethod,
			expectedStatus: http.StatusBadRequest,
			expectedError: expectedErrors{
				Msg:    errors.NewBookInvalidSerial().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestBookService_DeleteAllBooks(t *testing.T) {
	db, mock := newMock()
	defer db.Close()

	repo := NewBookController(db)

	deleteAllMethod := "DELETE"
	deleteAllURL := "/book"

	deleteAllBooksServiceMocks := []struct {
		testName       string
		mockFunc       func()
		service        BookService
		method         string
		url            string
		expectedStatus int
		expectedBody   int
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(bookdb.QueryDeleteAllBooksAndAlter)).WillReturnResult(sqlmock.NewResult(0, 4))
			},
			service:        repo,
			url:            deleteAllURL,
			method:         deleteAllMethod,
			expectedStatus: http.StatusOK,
			expectedBody:   4,
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Book(s) not found",
			mockFunc: func() {
				mock.ExpectExec(regexp.QuoteMeta(bookdb.QueryDeleteAllBooksAndAlter)).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			service:        repo,
			url:            deleteAllURL,
			method:         deleteAllMethod,
			expectedStatus: http.StatusNotFound,
			expectedError: expectedErrors{
				Msg: errors.NewBooksNotFound().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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
			expectedError: expectedErrors{
				Msg: errors.NewBookCouldNotQuery("sql: database is closed").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
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

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := &rowsResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			assert.Equal(t, tc.expectedStatus, w.Code)

			assert.Equal(t, tc.expectedBody, resultBody.Data)

			if !tc.expectedError.isEmpty() {
				assert.Equal(t, tc.expectedError, resultBody.Error, "Values are not equal:\nExpected: %+v\nActual: %+v", tc.expectedError, resultBody.Error)
			} else if !resultBody.Error.isEmpty() {
				t.Errorf("Expected error to be nil, found %+v", resultBody.Error)
			}

			for header, expected := range tc.expectedHeader {
				actual := w.Header().Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}
