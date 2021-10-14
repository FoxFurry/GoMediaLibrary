package controllers

import (
	"database/sql"
	bookDB "github.com/foxfurry/medialib/internal/book/db"
	"github.com/foxfurry/medialib/internal/book/domain/entity"
	"github.com/foxfurry/medialib/internal/book/http/errors"
	"github.com/foxfurry/medialib/internal/common/server/response"
	"github.com/foxfurry/medialib/internal/common/server/translator"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
)

type IService interface {

}

type BookService struct {
	dbRepo bookDB.BookDBRepository
}

func NewBookController(db *sql.DB) BookService {
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
			errors.HandleBookError(c, errors.NewBookValidatorError(translator.Translate(err)))
			return
		}
	}

	saveBook, err := b.dbRepo.SaveBook(&book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	response.OK(c, saveBook)
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

	response.OK(c, getBook)
}

func (b *BookService) GetAllBooks(c *gin.Context) {
	allBooks, err := b.dbRepo.GetAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	response.OK(c, allBooks)
}

func (b *BookService) SearchByAuthor(c *gin.Context) {
	author := c.Param("author")

	booksByAuthor, err := b.dbRepo.SearchByAuthor(author)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	response.OK(c, booksByAuthor)
}

func (b *BookService) SearchByTitle(c *gin.Context) {
	title := c.Param("title")

	bookByTitle, err := b.dbRepo.SearchByTitle(title)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	response.OK(c, bookByTitle)
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
			errors.HandleBookError(c, errors.NewBookValidatorError(translator.Translate(err)))
			return
		}
	}

	updatedBook, err := b.dbRepo.UpdateBook(uint64(id), &book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	response.OK(c, updatedBook)
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

	response.OK(c, nil)
}

func (b *BookService) DeleteAllBooks(c *gin.Context) {
	deletedRows, err := b.dbRepo.DeleteAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	response.OK(c, deletedRows)
}
