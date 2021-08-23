package integration_tests

import (
	"bytes"
	"encoding/json"
	"github.com/foxfurry/simple-rest/app"
	"github.com/foxfurry/simple-rest/configs"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	_ "github.com/foxfurry/simple-rest/internal/common/tests"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)



func init() {
	configs.LoadConfig()
	go app.NewTestApp().Start()
}

func purgeDB(){
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/book", nil)

	res, err := (&http.Client{}).Do(req)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		log.Printf("PURGE RESULT: %+v %v", res, err)
		log.Panicf("Could not purge database: %v %v", res.StatusCode, err)
	}
}

func TestSaveBook(t *testing.T) {
	saveURL := "http://localhost:8080/book"
	saveMethod := "POST"

	saveBookTests := []struct{
		testName string
		requestBody *entity.Book
		requestMethod string
		requestURL string
		expectedStatus int
		expectedBody string
		expectedHeader map[string]string
	}{
		{
			testName:       "Test Successful",
			requestBody:    &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 200,
			expectedBody:   "{\"id\":1,\"title\":\"Fahrenheit 451\",\"author\":\"Ray Bradbury\",\"year\":1953,\"description\":\"Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing title",
			requestBody:    &entity.Book{
				Author:      "Ray Bradbury",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Title cannot be empty\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing author",
			requestBody:    &entity.Book{
				Title:      "Fahrenheit 451",
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Author cannot be empty\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing year",
			requestBody:    &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Year cannot be empty\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing title and author",
			requestBody:    &entity.Book{
				Year:        1953,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Title cannot be empty, Author cannot be empty\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing title and year",
			requestBody:    &entity.Book{
				Author:      "Ray Bradbury",
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Title cannot be empty, Year cannot be empty\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing author and year",
			requestBody:    &entity.Book{
				Title:       "Fahrenheit 451",
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Author cannot be empty, Year cannot be empty\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing title, author and year",
			requestBody:    &entity.Book{
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Title cannot be empty, Author cannot be empty, Year cannot be empty\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Invalid year",
			requestBody:    &entity.Book{
				Title:       "Fahrenheit 451",
				Author:      "Ray Bradbury",
				Year:        2069,
				Description: "Fahrenheit 451 is a 1953 dystopian novel by American writer Ray Bradbury. Often regarded as one of his best works",
			},
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Invalid request body: Year should be between -868 and 2021\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
		{
			testName:       "Test Unsuccessful: Missing body",
			requestBody:    nil,
			requestMethod:  saveMethod,
			requestURL:     saveURL,
			expectedStatus: 400,
			expectedBody:   "{\"error:\":\"Expected body, found EOF\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
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

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			assert.Equal(t, tc.expectedBody, string(body))

			for header, expected := range tc.expectedHeader {
				actual := resp.Header.Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	getURL := "http://localhost:8080/book"
	getMethod := "GET"

	getBookTests := []struct{
		testName      string
		requestSave   []entity.Book
		requestMethod string
		requestURL string
		requestID string
		expectedStatus int
		expectedBody string
		expectedHeader map[string]string
	}{
		{
			testName:       "Test successful",
			requestSave:    []entity.Book{
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
			expectedBody:   "{\"id\":2,\"title\":\"Brave New World\",\"author\":\"Aldous Huxley\",\"year\":1932,\"description\":\"Brave New World is a dystopian social science fiction novel by English author Aldous Huxley, written in 1931 and published in 1932.\"}",
			expectedHeader: map[string]string{
				"Content-Type":"application/json; charset=utf-8",
			},
		},
	}

	for _, tc := range getBookTests {
		purgeDB()
		t.Run(tc.testName, func(t *testing.T) {
			client := &http.Client{}

			if tc.requestSave != nil {
				var jsonRequest []byte = nil

				for _, book := range tc.requestSave {
					jsonRequest, _ = json.Marshal(book)

					req, _ := http.NewRequest("POST", tc.requestURL, bytes.NewReader(jsonRequest))

					resp, err := client.Do(req)
					if err != nil || resp.StatusCode != http.StatusOK {
						t.Errorf("Could not complete the request: %v", err)
					}
				}
			}

			url := getURL + "/" + tc.requestID

			resp, err := http.Get(url)
			if err != nil || resp.StatusCode != http.StatusOK {
				t.Errorf("Could not complete the request: %v", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Could not read the body: %v", err)
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			assert.Equal(t, tc.expectedBody, string(body))

			for header, expected := range tc.expectedHeader {
				actual := resp.Header.Get(header)
				assert.Equal(t, expected, actual)
			}
		})
	}
}