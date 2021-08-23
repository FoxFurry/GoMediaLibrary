package controllers

import (
	"database/sql"
	bookDB "github.com/foxfurry/simple-rest/internal/book/db"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	"github.com/foxfurry/simple-rest/internal/book/http/validators"
	"github.com/gin-gonic/gin"
	"io"
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
			errors.HandleBookError(c, errors.BookEmptyBody{})
			return
		}else {
			errors.HandleBookError(c, errors.BookBadBody{Msg: validators.Translate(err)})
			return
		}
	}

	saveBook, err := b.dbRepo.SaveBook(&book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, saveBook)
}

func (b *BookService) GetBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		errors.HandleBookError(c, errors.BookInvalidSerial{})
		return
	}

	getBook, err := b.dbRepo.GetBook(uint64(id))
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, getBook)
}

func (b *BookService) GetAllBooks(c *gin.Context) {
	allBooks, err := b.dbRepo.GetAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, allBooks)
}

func (b *BookService) SearchByAuthor(c *gin.Context) {
	author := c.Param("author")

	bookByAuthor, err := b.dbRepo.SearchByAuthor(author)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, bookByAuthor)
}

func (b *BookService) SearchByTitle(c *gin.Context) {
	title := c.Param("title")

	bookByTitle, err := b.dbRepo.SearchByTitle(title)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, bookByTitle)
}

func (b *BookService) UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		errors.HandleBookError(c, errors.BookInvalidSerial{})
		return
	}

	var book entity.Book

	if err = c.ShouldBindJSON(&book); err != nil {
		if err == io.EOF {
			errors.HandleBookError(c, errors.BookEmptyBody{})
			return
		}else {
			errors.HandleBookError(c, errors.BookBadBody{Msg: validators.Translate(err)})
			return
		}
	}

	updatedRows, err := b.dbRepo.UpdateBook(uint64(id), &book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, updatedRows)
}

func (b *BookService) DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		errors.HandleBookError(c, errors.BookInvalidSerial{})
		return
	}

	_, err = b.dbRepo.DeleteBook(uint64(id))
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.Status(200)
}

func (b *BookService) DeleteAllBooks(c *gin.Context) {
	deletedRows, err := b.dbRepo.DeleteAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, gin.H{"Deleted rows": deletedRows})
}
