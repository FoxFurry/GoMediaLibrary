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

type BookNotFound struct{}

type BookBadRequest struct{}

type BookBadScanOptions struct{
	Msg string
}

type BookTitleAlreadyExists struct{}

type BookCouldNotQuery struct{
	Msg string
}

func (b BookNotFoundByTitle) Error() string {
	return fmt.Sprintf("Book(s) with title %v not found in db", b.Title)
}

func (b BookNotFoundByAuthor) Error() string {
	return fmt.Sprintf("Book(s) with author %v not found in db", b.Author)
}

func (b BookNotFound) Error() string {
	return "Book(s) not found in db"
}

func (b BookBadRequest) Error() string {
	return "Title, author and year are mandatory fields"
}

func (b BookTitleAlreadyExists) Error() string {
	return "Requested title already exists"
}

func (b BookBadScanOptions) Error() string {
	return fmt.Sprintf("Bad SQL scan options: %v", b.Msg)
}

func (b BookCouldNotQuery) Error() string {
	return fmt.Sprintf("Could not execute query: %v", b.Msg)
}

func HandleBookError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case BookNotFound:
		server.RespondNotFound(w, err.Error())
	case BookNotFoundByAuthor:
		server.RespondNotFound(w, err.Error())
	case BookNotFoundByTitle:
		server.RespondNotFound(w, err.Error())
	case BookBadRequest:
		server.RespondBadRequest(w, err.Error())
	case BookTitleAlreadyExists:
		server.RespondAlreadyExists(w, err.Error())
	default:
		server.RespondInternalError(w, fmt.Sprintf("Internal Error: %v", err))
	}
}
