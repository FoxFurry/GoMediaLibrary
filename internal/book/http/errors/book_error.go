package errors

import (
	"fmt"
	"github.com/foxfurry/simple-rest/internal/common/server"
	"github.com/gin-gonic/gin"
)

type BookNotFoundByTitle struct {
	Title string
}

type BookNotFoundByAuthor struct {
	Author string
}

type BooksNotFound struct{}

type BookBadBody struct{
	Msg string
}

type BookBadScanOptions struct {
	Msg string
}

type BookTitleAlreadyExists struct{}

type BookCouldNotQuery struct {
	Msg string
}

type BookInvalidSerial struct{}

type BookUnexpectedError struct {
	Msg string
}

type BookEmptyBody struct{}

func (b BookNotFoundByTitle) Error() string {
	return fmt.Sprintf("Book(s) with title %v not found in db", b.Title)
}

func (b BookNotFoundByAuthor) Error() string {
	return fmt.Sprintf("Book(s) with author %v not found in db", b.Author)
}

func (b BooksNotFound) Error() string {
	return "Book(s) not found in db"
}

func (b BookBadBody) Error() string {
	return fmt.Sprintf("Invalid request body: %v", b.Msg)
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

func (b BookInvalidSerial) Error() string {
	return "Invalid serial. Serial must be more than 1"
}

func (b BookUnexpectedError) Error() string {
	return fmt.Sprintf("Unexpected error: %v", b.Msg)
}

func (b BookEmptyBody) Error() string {
	return "Expected body, found EOF"
}

func HandleBookError(c *gin.Context, err error) {
	switch err.(type) {
	case BooksNotFound, BookNotFoundByAuthor, BookNotFoundByTitle:
		server.RespondNotFound(c, err.Error())
	case BookBadBody, BookInvalidSerial, BookEmptyBody:
		server.RespondBadRequest(c, err.Error())
	case BookTitleAlreadyExists:
		server.RespondAlreadyExists(c, err.Error())
	case BookUnexpectedError, BookCouldNotQuery:
		server.RespondInternalError(c, err.Error())
	default:
		server.RespondInternalError(c, fmt.Sprintf("Internal Error: %v", err))
	}
}
