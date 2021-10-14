package router

import (
	"database/sql"
	"github.com/foxfurry/medialib/internal/book/http/controllers"
	"github.com/foxfurry/medialib/internal/book/http/validators"
	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(router *gin.Engine, db *sql.DB) {
	bookRepo := controllers.NewBookService(db)

	book := router.Group("/book")
	{
		book.GET("/:id", bookRepo.GetBook)

		book.GET("/title/:title", bookRepo.SearchByTitle)
		book.GET("/title/", bookRepo.SearchByTitle)

		book.GET("/author/:author", bookRepo.SearchByAuthor)
		book.GET("/author/", bookRepo.SearchByAuthor)

		book.GET("/", bookRepo.GetAllBooks)

		book.POST("/", bookRepo.SaveBook)

		book.PUT("/", bookRepo.UpdateBook)

		book.DELETE("/:id", bookRepo.DeleteBook)
		book.DELETE("/", bookRepo.DeleteAllBooks)
	}

	validators.RegisterBookValidators()
}
