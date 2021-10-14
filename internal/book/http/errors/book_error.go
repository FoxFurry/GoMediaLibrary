package errors

import (
	"encoding/json"
	"fmt"
	"github.com/foxfurry/medialib/internal/common/server/errors"
	"github.com/foxfurry/medialib/internal/common/server/response"
	validator "github.com/foxfurry/medialib/internal/common/server/translator"
	"github.com/gin-gonic/gin"
	"log"
)

type bookNotFoundByTitle struct {
	errors.Common
}

type bookNotFoundByAuthor struct {
	errors.Common
}

type booksNotFound struct{
	errors.Common
}

type bookBadBody struct{
	errors.Common
}

type bookBadScanOptions struct {
	errors.Common
}

type bookTitleAlreadyExists struct{
	errors.Common
}

type bookCouldNotQuery struct {
	errors.Common
}

type bookInvalidSerial struct{
	errors.Common
}

type bookUnexpectedError struct {
	errors.Common
}

type bookEmptyBody struct{
	errors.Common
}

type bookValidatorError struct {
	Fields []validator.FieldError	`json:"fields"`
}

func NewBookNotFoundByTitle(title string) bookNotFoundByTitle {
	return bookNotFoundByTitle{
		errors.Common{Msg: fmt.Sprintf("Book(s) with title %v not found in db", title)},
	}
}

func NewBookNotFoundByAuthor(author string) bookNotFoundByAuthor {
	return bookNotFoundByAuthor{
		errors.Common{Msg: fmt.Sprintf("Book(s) with author %v not found in db", author)},
	}
}

func NewBooksNotFound() booksNotFound {
	return booksNotFound{
		errors.Common{Msg: "Book(s) not found in db"},
	}
}

func NewBookTitleAlreadyExists() bookTitleAlreadyExists {
	return bookTitleAlreadyExists{
		errors.Common{Msg: "Requested title already exists"},
	}
}

func NewBookBadScanOptions(msg string) bookBadScanOptions {
	return bookBadScanOptions{
		errors.Common{Msg: fmt.Sprintf("Bad SQL scan options: %v", msg)},
	}
}

func NewBookCouldNotQuery(msg string) bookCouldNotQuery {
	return bookCouldNotQuery{
		errors.Common{Msg: fmt.Sprintf("Could not execute query: %v", msg)},
	}
}

func NewBookInvalidSerial() bookInvalidSerial {
	return bookInvalidSerial{
		errors.Common{Msg: "Invalid serial. Serial must be more than 1"},
	}
}

func NewBookUnexpectedError(msg string) bookUnexpectedError {
	return bookUnexpectedError{
		errors.Common{Msg: fmt.Sprintf("Unexpected error: %v", msg)},
	}
}

func NewBookEmptyBody() bookEmptyBody {
	return bookEmptyBody{
		errors.Common{Msg: "Expected body, found EOF"},
	}
}

func NewBookValidatorError(fields []validator.FieldError) bookValidatorError {
	return bookValidatorError{Fields: fields}
}

func (b bookValidatorError) Error() string {
	var res = ""
	for _, f := range b.Fields {
		tmp, err := json.Marshal(f)
		if err != nil {
			log.Fatalf("Could not marshal field error: %v", err)
		}
		res += fmt.Sprintf("%s", tmp)
	}
	return res
}

func HandleBookError(c *gin.Context, err error) {
	switch err.(type) {
	case booksNotFound, bookNotFoundByAuthor, bookNotFoundByTitle:
		response.NotFound(c, err)
	case bookValidatorError, bookInvalidSerial, bookEmptyBody:
		response.BadRequest(c, err)
	case bookTitleAlreadyExists:
		response.AlreadyExists(c, err)
	case bookUnexpectedError, bookCouldNotQuery:
		response.InternalError(c, err)
	default:
		response.InternalError(c, err)
	}
}
