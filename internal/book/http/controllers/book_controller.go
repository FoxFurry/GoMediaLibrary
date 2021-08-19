package controllers

import (
	"database/sql"
	bookDB "github.com/foxfurry/simple-rest/internal/book/db"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	"github.com/gin-gonic/gin"
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

	if err := c.BindJSON(&book); err != nil {
		errors.HandleBookError(c, errors.BookBadRequest{})
		return
	}

	saveBook, err := b.dbRepo.SaveBook(&book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, saveBook)
}

func (b *BookService) GetBook(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.Atoi(params)

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
	books, err := b.dbRepo.GetAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, books)
}

func (b *BookService) SearchByAuthor(c *gin.Context) {
	author := c.Param("author")

	byAuthor, err := b.dbRepo.SearchByAuthor(author)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, byAuthor)
}

func (b *BookService) SearchByTitle(c *gin.Context) {
	title := c.Param("title")

	byTitle, err := b.dbRepo.SearchByTitle(title)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, byTitle)
}

func (b *BookService) UpdateBook(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.Atoi(params)

	if err != nil {
		errors.HandleBookError(c, errors.BookInvalidSerial{})
		return
	}

	var book entity.Book

	if err = c.BindJSON(&book); err != nil {
		errors.HandleBookError(c, errors.BookBadRequest{})
		return
	}

	updatedRows, err := b.dbRepo.UpdateBook(uint64(id), &book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, updatedRows)
}

func (b *BookService) DeleteBook(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.Atoi(params)

	if err != nil {
		errors.HandleBookError(c, errors.BookInvalidSerial{})
		return
	}

	deletedRows, err := b.dbRepo.DeleteBook(uint64(id))
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, deletedRows)
}

func (b *BookService) DeleteAllBooks(c *gin.Context) {
	deletedRows, err := b.dbRepo.DeleteAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, deletedRows)
}
