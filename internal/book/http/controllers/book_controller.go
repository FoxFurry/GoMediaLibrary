package controllers

import (
	"database/sql"
	bookDB "github.com/foxfurry/simple-rest/internal/book/db"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type BookApp struct {
	dbRepo bookDB.BookDBRepository
}

func NewBookApp(db *sql.DB) BookApp {
	return BookApp{
		dbRepo: bookDB.NewBookRepo(db),
	}
}

func (b *BookApp) SaveBook(c *gin.Context) {
	var book entity.Book

	if err := c.BindJSON(&book); err != nil {
		errors.HandleBookError(c, err)
		return
	}

	saveBook, err := b.dbRepo.SaveBook(&book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, saveBook)
}

func (b *BookApp) GetBook(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.Atoi(params)

	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	getBook, err := b.dbRepo.GetBook(uint64(id))
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, getBook)
}

func (b *BookApp) GetAllBooks(c *gin.Context) {
	books, err := b.dbRepo.GetAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, books)
}

func (b *BookApp) SearchByAuthor(c *gin.Context) {
	author := c.Param("author")

	byAuthor, err := b.dbRepo.SearchByAuthor(author)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, byAuthor)
}

func (b *BookApp) SearchByTitle(c *gin.Context) {
	title := c.Param("title")

	byTitle, err := b.dbRepo.SearchByTitle(title)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, byTitle)
}

func (b *BookApp) UpdateBook(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.Atoi(params)

	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	var book *entity.Book

	if err := c.BindJSON(book); err != nil {
		errors.HandleBookError(c, err)
		return
	}

	updatedRows, err := b.dbRepo.UpdateBook(uint64(id), book)
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, updatedRows)
}

func (b *BookApp) DeleteBook(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.Atoi(params)

	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	deletedRows, err := b.dbRepo.DeleteBook(uint64(id))
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, deletedRows)
}

func (b *BookApp) DeleteAllBooks(c *gin.Context) {
	deletedRows, err := b.dbRepo.DeleteAllBooks()
	if err != nil {
		errors.HandleBookError(c, err)
		return
	}

	c.JSON(200, deletedRows)
}
