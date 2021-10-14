package integration_tests

import (
	"bytes"
	"encoding/json"
	"github.com/foxfurry/medialib/app"
	"github.com/foxfurry/medialib/configs"
	"github.com/foxfurry/medialib/internal/book/domain/entity"
	"github.com/foxfurry/medialib/internal/book/http/errors"
	"github.com/foxfurry/medialib/internal/book/http/validators"
	"github.com/foxfurry/medialib/internal/common/server/translator"
	_ "github.com/foxfurry/medialib/internal/common/tests"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

var baseURL = "http://localhost:8080/book"

type expectedErrors struct {
	Msg    string                  `json:"msg,omitempty"`
	Fields []translator.FieldError `json:"fields,omitempty"`
}

type singleResponse struct {
	Data  *entity.Book   `json:"data"`
	Error expectedErrors `json:"error"`
}

type arrayResponse struct {
	Data  []entity.Book  `json:"data"`
	Error expectedErrors `json:"error"`
}

func (e expectedErrors) isEmpty() bool {
	return e.Msg == "" && e.Fields == nil
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	configs.LoadConfig()
	go app.NewTestApp().Start()
	return m.Run()
}

func purgeDB() {
	req, err := http.NewRequest(http.MethodDelete, baseURL, strings.NewReader(""))

	if err != nil {
		log.Fatalf("Could not create a request: %v", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Panicf("Could not purge database: %v %v", res, err)
	}
}

func TestSaveBook(t *testing.T) {
	saveURL := baseURL
	saveMethod := http.MethodPost

	saveBookTests := []struct {
		testName       string
		requestBody    *entity.Book
		requestMethod  string
		requestURL     string
		expectedStatus int
		expectedBody   *entity.Book
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful",
			requestBody: &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 200,
			expectedBody: &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Missing title",
			requestBody: &entity.Book{
				Author:      "Ray Bradbury",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
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
			testName: "Test Unsuccessful: Missing author",
			requestBody: &entity.Book{
				Title:       "Fahrenheit 451",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
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
			testName: "Test Unsuccessful: Missing year",
			requestBody: &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldYearEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Missing title and author",
			requestBody: &entity.Book{
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldTitleEmpty,
					validators.FieldAuthorEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Missing title and year",
			requestBody: &entity.Book{
				Author:      "Ray Bradbury",
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
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
			testName: "Test Unsuccessful: Missing author and year",
			requestBody: &entity.Book{
				Title:       "Fahrenheit 451",
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldAuthorEmpty,
					validators.FieldYearEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing title, author and year",
			requestBody:    &entity.Book{},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldTitleEmpty,
					validators.FieldAuthorEmpty,
					validators.FieldYearEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Unsuccessful: Invalid year",
			requestBody: &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Year:        2069,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldYearInvalid,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing body",
			requestBody:    nil,
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedError: expectedErrors{
				Msg: errors.NewBookEmptyBody().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range saveBookTests {
		purgeDB()
		t.Run(tc.testName, func(t *testing.T) {
			client := &http.Client{}
			var jsonRequest []byte = nil

			if tc.requestBody != nil {
				jsonRequest, _ = json.Marshal(tc.requestBody)
			}

			req, _ := http.NewRequest(tc.requestMethod, tc.requestURL, bytes.NewReader(jsonRequest))
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Could not complete the request: %v", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := singleResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			for header, expected := range tc.expectedHeader {
				actual := resp.Header.Get(header)
				assert.Equal(t, expected, actual)
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
		})
	}
}

func TestGetBook(t *testing.T) {
	getURL := baseURL + "/"
	getMethod := http.MethodGet

	getBookTests := []struct {
		testName       string
		requestBody    []entity.Book
		requestMethod  string
		requestURL     string
		requestID      string
		expectedStatus int
		expectedBody   *entity.Book
		expectedError  expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test successful",
			requestBody: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
				},
				{
					Title:       "Brave New World",
					Author:      "Aldous Huxley",
					Year:        1932,
					Description: "Brave New World is a dystopian social science fiction novel by English author Aldous Huxley, written in 1931 and published in 1932.",
				},
			},
			requestMethod:  getMethod,
			requestURL:     getURL,
			requestID:      "2",
			expectedStatus: 200,
			expectedBody: &entity.Book{
				Title:       "Brave New World",
				Author:      "Aldous Huxley",
				Year:        1932,
				Description: "Brave New World is a dystopian social science fiction novel by English author Aldous Huxley, written in 1931 and published in 1932.",
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Data(s) not found",
			requestBody:    []entity.Book{},
			requestMethod:  getMethod,
			requestURL:     getURL,
			requestID:      "1",
			expectedStatus: 404,
			expectedError: expectedErrors{
				Msg: errors.NewBooksNotFound().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Invalid serial",
			requestBody:    []entity.Book{},
			requestMethod:  getMethod,
			requestURL:     getURL,
			requestID:      "0",
			expectedStatus: 400,
			expectedError: expectedErrors{
				Msg: errors.NewBookInvalidSerial().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range getBookTests {
		purgeDB()
		t.Run(tc.testName, func(t *testing.T) {
			client := &http.Client{}

			if tc.requestBody != nil {
				var jsonRequest []byte = nil

				for _, book := range tc.requestBody {
					jsonRequest, _ = json.Marshal(book)

					req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonRequest))

					resp, err := client.Do(req)
					if err != nil || resp.StatusCode != http.StatusOK {
						t.Errorf("Could not complete the request: %v", err)
					}
				}
			}

			url := getURL + tc.requestID

			resp, err := http.Get(url)
			if err != nil {
				t.Errorf("Could not complete the request: %v", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := singleResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			for header, expected := range tc.expectedHeader {
				actual := resp.Header.Get(header)
				assert.Equal(t, expected, actual)
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
		})
	}
}

func TestGetAllBook(t *testing.T) {
	getAllURL := baseURL
	getAllMethod := "GET"

	getAllBookTests := []struct {
		testName          string
		requestBody       []entity.Book
		requestMethod     string
		requestURL        string
		expectedStatus    int
		expectedBodyArray []entity.Book
		expectedError     expectedErrors
		expectedHeader    map[string]string
	}{
		{
			testName: "Test Successful: Multiple books",
			requestBody: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
				{
					Title:       "Brave New World",
					Author:      "Aldous Huxley",
					Year:        1932,
					Description: "Brave New World is a dystopian social science fiction novel by English author Aldous Huxley, written in 1931 and published in 1932.",
				},
				{
					Title:       "Nineteen Eighty-Four",
					Author:      "George Orwell",
					Year:        1949,
					Description: "Nineteen Eighty-Four, often referred to as 1984, is a dystopian social science fiction novel by the English novelist George Orwell.",
				},
			},
			requestMethod:  getAllMethod,
			requestURL:     getAllURL,
			expectedStatus: 200,
			expectedBodyArray: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
				{
					Title:       "Brave New World",
					Author:      "Aldous Huxley",
					Year:        1932,
					Description: "Brave New World is a dystopian social science fiction novel by English author Aldous Huxley, written in 1931 and published in 1932.",
				},
				{
					Title:       "Nineteen Eighty-Four",
					Author:      "George Orwell",
					Year:        1949,
					Description: "Nineteen Eighty-Four, often referred to as 1984, is a dystopian social science fiction novel by the English novelist George Orwell.",
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Successful: Single book",
			requestBody: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
			},
			requestMethod:  getAllMethod,
			requestURL:     getAllURL,
			expectedStatus: 200,
			expectedBodyArray: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Data(s) not found",
			requestBody:    []entity.Book{},
			requestMethod:  getAllMethod,
			requestURL:     getAllURL,
			expectedStatus: 404,
			expectedError: expectedErrors{
				Msg: errors.NewBooksNotFound().Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range getAllBookTests {
		purgeDB()
		t.Run(tc.testName, func(t *testing.T) {
			client := &http.Client{}

			if tc.requestBody != nil {
				var jsonRequest []byte = nil

				for _, book := range tc.requestBody {
					jsonRequest, _ = json.Marshal(book)

					req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonRequest))

					resp, err := client.Do(req)
					if err != nil || resp.StatusCode != http.StatusOK {
						t.Errorf("Could not complete the request: %v %v", resp.Status, err)
					}
				}
			}

			resp, err := http.Get(getAllURL)
			if err != nil {
				t.Errorf("Could not complete the request: %v", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := &arrayResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			for header, expected := range tc.expectedHeader {
				actual := resp.Header.Get(header)
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

func TestSearchByAuthorBook(t *testing.T) {
	searchByAuthorURL := baseURL + "/author/"
	searchByAuthorMethod := http.MethodGet

	searchByAuthorTests := []struct {
		testName       string
		requestBody    []entity.Book
		requestMethod  string
		requestURL     string
		requestAuthor  string
		expectedStatus    int
		expectedBodyArray []entity.Book
		expectedError     expectedErrors
		expectedHeader map[string]string
	}{
		{
			testName: "Test Successful: Multiple books",
			requestBody: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
				{
					Title:       "The Martian Chronicles",
					Author:      "Ray Bradbury",
					Year:        1950,
					Description: "The Martian Chronicles is a science fiction fix-up novel, published in 1950, by American writer Ray Bradbury that chronicles the exploration and settlement of Mars, the home of indigenous Martians, by Americans leaving a troubled Earth that is eventually devastated by nuclear war.",
				},
				{
					Title:       "Dandelion Wine",
					Author:      "Ray Bradbury",
					Year:        1957,
					Description: "Dandelion Wine is a 1957 novel by Ray Bradbury set in the summer of 1928 in the fictional town of Green Town, Illinois, based upon Bradbury's childhood home of Waukegan, Illinois.",
				},
			},
			requestMethod:  searchByAuthorMethod,
			requestURL:     searchByAuthorURL,
			requestAuthor:  "Ray Bradbury",
			expectedStatus: 200,
			expectedBodyArray: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
				{
					Title:       "The Martian Chronicles",
					Author:      "Ray Bradbury",
					Year:        1950,
					Description: "The Martian Chronicles is a science fiction fix-up novel, published in 1950, by American writer Ray Bradbury that chronicles the exploration and settlement of Mars, the home of indigenous Martians, by Americans leaving a troubled Earth that is eventually devastated by nuclear war.",
				},
				{
					Title:       "Dandelion Wine",
					Author:      "Ray Bradbury",
					Year:        1957,
					Description: "Dandelion Wine is a 1957 novel by Ray Bradbury set in the summer of 1928 in the fictional town of Green Town, Illinois, based upon Bradbury's childhood home of Waukegan, Illinois.",
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName: "Test Successful: Single book",
			requestBody: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
				{
					Title:       "Brave New World",
					Author:      "Aldous Huxley",
					Year:        1932,
					Description: "Brave New World is a dystopian social science fiction novel by English author Aldous Huxley, written in 1931 and published in 1932.",
				},
				{
					Title:       "Nineteen Eighty-Four",
					Author:      "George Orwell",
					Year:        1949,
					Description: "Nineteen Eighty-Four, often referred to as 1984, is a dystopian social science fiction novel by the English novelist George Orwell.",
				},
			},
			requestMethod:  searchByAuthorMethod,
			requestURL:     searchByAuthorURL,
			requestAuthor:  "Ray Bradbury",
			expectedStatus: 200,
			expectedBodyArray: []entity.Book{
				{
					Title:       "Fahrenheit 451",
					Author:      "Ray Bradbury",
					Year:        1953,
					Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works.",
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Book(s) not found",
			requestBody:    []entity.Book{},
			requestMethod:  searchByAuthorMethod,
			requestURL:     searchByAuthorURL,
			requestAuthor:  "Ray Bradbury",
			expectedStatus: 404,
			expectedError: expectedErrors{
				Msg: errors.NewBookNotFoundByAuthor("Ray Bradbury").Error(),
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Invalid author",
			requestBody:    []entity.Book{},
			requestMethod:  searchByAuthorMethod,
			requestURL:     searchByAuthorURL,
			requestAuthor:  "",
			expectedStatus: 400,
			expectedError: expectedErrors{
				Fields: []translator.FieldError{
					validators.FieldAuthorEmpty,
				},
			},
			expectedHeader: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range searchByAuthorTests {
		purgeDB()
		t.Run(tc.testName, func(t *testing.T) {
			client := &http.Client{}

			if tc.requestBody != nil {
				var jsonRequest []byte = nil

				for _, book := range tc.requestBody {
					jsonRequest, _ = json.Marshal(book)

					req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonRequest))

					resp, err := client.Do(req)
					if err != nil || resp.StatusCode != http.StatusOK {
						t.Errorf("Could not complete the request: %v", err)
					}
				}
			}

			url := tc.requestURL + tc.requestAuthor
			resp, err := http.Get(url)
			if err != nil {
				t.Errorf("Could not complete the request: %v", err)
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}
			resultBody := &arrayResponse{}

			if err = json.Unmarshal(body, &resultBody); err != nil {
				t.Errorf("Unable to unmarshal the body")
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			for header, expected := range tc.expectedHeader {
				actual := resp.Header.Get(header)
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
