package controllers

import (
	"database/sql"
	bookDB "github.com/foxfurry/simple-rest/internal/book/db"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	"github.com/foxfurry/simple-rest/internal/common/server/common_response"
	"github.com/foxfurry/simple-rest/internal/common/server/common_translators"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

type BookService struct {
	dbRepo bookDB.BookDBRepository
}

func NewBookService(db *sql.DB) BookService {
	return BookService{
		dbRepo: bookDB.NewBookRepo(db),
	}
}

func (b *BookService) SaveBook(c *gin.Context) {
	var book entity.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		if err == io.EOF {
			errors.HandleBookError(c, errors.NewBookEmptyBody())
			return
		} else {
			errors.HandleBookError(c, errors.NewBookValidatorError(common_translators.Translate(err)))
			return
		}
	}

	saveBook, err := b.dbRepo.SaveBook(&book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, saveBook, nil)
}

func (b *BookService) GetBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		errors.HandleBookError(c, errors.NewBookInvalidSerial())
		return
	}

	getBook, err := b.dbRepo.GetBook(uint64(id))
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, getBook, nil)
}

func (b *BookService) GetAllBooks(c *gin.Context) {
	allBooks, err := b.dbRepo.GetAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, allBooks, nil)
}

func (b *BookService) SearchByAuthor(c *gin.Context) {
	author := c.Param("author")

	booksByAuthor, err := b.dbRepo.SearchByAuthor(author)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, booksByAuthor, nil)
}

func (b *BookService) SearchByTitle(c *gin.Context) {
	title := c.Param("title")

	bookByTitle, err := b.dbRepo.SearchByTitle(title)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, bookByTitle, nil)
}

func (b *BookService) UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		errors.HandleBookError(c, errors.NewBookInvalidSerial())
		return
	}

	var book entity.Book

	if err = c.ShouldBindJSON(&book); err != nil {
		if err == io.EOF {
			errors.HandleBookError(c, errors.NewBookEmptyBody())
			return
		} else {
			errors.HandleBookError(c, errors.NewBookValidatorError(common_translators.Translate(err)))
			return
		}
	}

	updatedBook, err := b.dbRepo.UpdateBook(uint64(id), &book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, updatedBook, nil)
}

func (b *BookService) DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		errors.HandleBookError(c, errors.NewBookInvalidSerial())
		return
	}

	_, err = b.dbRepo.DeleteBook(uint64(id))
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, nil, nil)
}

func (b *BookService) DeleteAllBooks(c *gin.Context) {
	deletedRows, err := b.dbRepo.DeleteAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	common_response.Respond(c, http.StatusOK, deletedRows, nil)
}
