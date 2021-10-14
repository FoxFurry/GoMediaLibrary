package errors

import (
	"encoding/json"
	"fmt"
	"github.com/foxfurry/medialib/internal/common/server/common_errors"
	validator "github.com/foxfurry/medialib/internal/common/server/common_translators"
	"github.com/gin-gonic/gin"
	"log"
)

type bookNotFoundByTitle struct {
	common_errors.CommonError
}

type bookNotFoundByAuthor struct {
	common_errors.CommonError
}

type booksNotFound struct{
	common_errors.CommonError
}

type bookBadBody struct{
	common_errors.CommonError
}

type bookBadScanOptions struct {
	common_errors.CommonError
}

type bookTitleAlreadyExists struct{
	common_errors.CommonError
}

type bookCouldNotQuery struct {
	common_errors.CommonError
}

type bookInvalidSerial struct{
	common_errors.CommonError
}

type bookUnexpectedError struct {
	common_errors.CommonError
}

type bookEmptyBody struct{
	common_errors.CommonError
}

type bookValidatorError struct {
	Fields []validator.FieldError	`json:"fields"`
}

func NewBookNotFoundByTitle(title string) bookNotFoundByTitle {
	return bookNotFoundByTitle{
		common_errors.CommonError{Msg: fmt.Sprintf("Book(s) with title %v not found in db", title)},
	}
}

func NewBookNotFoundByAuthor(author string) bookNotFoundByAuthor {
	return bookNotFoundByAuthor{
		common_errors.CommonError{Msg: fmt.Sprintf("Book(s) with author %v not found in db", author)},
	}
}

func NewBooksNotFound() booksNotFound {
	return booksNotFound{
		common_errors.CommonError{Msg: "Book(s) not found in db"},
	}
}

func NewBookTitleAlreadyExists() bookTitleAlreadyExists {
	return bookTitleAlreadyExists{
		common_errors.CommonError{Msg: "Requested title already exists"},
	}
}

func NewBookBadScanOptions(msg string) bookBadScanOptions {
	return bookBadScanOptions{
		common_errors.CommonError{Msg: fmt.Sprintf("Bad SQL scan options: %v", msg)},
	}
}

func NewBookCouldNotQuery(msg string) bookCouldNotQuery {
	return bookCouldNotQuery{
		common_errors.CommonError{Msg: fmt.Sprintf("Could not execute query: %v", msg)},
	}
}

func NewBookInvalidSerial() bookInvalidSerial {
	return bookInvalidSerial{
		common_errors.CommonError{Msg: "Invalid serial. Serial must be more than 1"},
	}
}

func NewBookUnexpectedError(msg string) bookUnexpectedError {
	return bookUnexpectedError{
		common_errors.CommonError{Msg: fmt.Sprintf("Unexpected error: %v", msg)},
	}
}

func NewBookEmptyBody() bookEmptyBody {
	return bookEmptyBody{
		common_errors.CommonError{Msg: "Expected body, found EOF"},
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
		common_errors.RespondNotFound(c, err)
	case bookValidatorError, bookInvalidSerial, bookEmptyBody:
		common_errors.RespondBadRequest(c, err)
	case bookTitleAlreadyExists:
		common_errors.RespondAlreadyExists(c, err)
	case bookUnexpectedError, bookCouldNotQuery:
		common_errors.RespondInternalError(c, err)
	default:
		common_errors.RespondInternalError(c, err)
	}
}
