package errors

import (
	"fmt"
	"github.com/foxfurry/simple-rest/internal/common/server"
	"net/http"
)

type BookNotFoundByTitle struct {
	Title string
}

type BookNotFoundByAuthor struct {
	Author string
}

type BookNotFound struct {}

func (b BookNotFoundByTitle) Error() string {
	return fmt.Sprintf("Book(s) with title %v not found in db", b.Title)
}

func (b BookNotFoundByAuthor) Error() string {
	return fmt.Sprintf("Book(s) with author %v not found in db", b.Author)
}

func (b BookNotFound) Error() string {
	return "Book(s) not found in db"
}

func HandleBookError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case BookNotFound:
		server.RespondNotFound(w, err.Error())
	case BookNotFoundByAuthor:
		server.RespondNotFound(w, err.Error())
	case BookNotFoundByTitle:
		server.RespondNotFound(w, err.Error())
	default:
		server.RespondInternalError(w, fmt.Sprintf("Internal Error: %v", err))
	}
}
