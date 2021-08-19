package routers

import (
	"database/sql"
	"github.com/foxfurry/simple-rest/internal/book/http/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(router *gin.Engine, db *sql.DB) {
	bookRepo := controllers.NewBookApp(db)

	book := router.Group("/book")
	{
		book.GET("/:id", bookRepo.GetBook)
		book.GET("/title/:title", bookRepo.SearchByTitle)
		book.GET("/author/:author", bookRepo.SearchByAuthor)
		book.GET("/", bookRepo.GetAllBooks)

		book.POST("/", bookRepo.SaveBook)

		book.PUT("/", bookRepo.UpdateBook)

		book.DELETE("/:id", bookRepo.DeleteBook)
		book.DELETE("/", bookRepo.DeleteAllBooks)
	}
}
