package routers

import (
	"github.com/foxfurry/simple-rest/app"
	"github.com/foxfurry/simple-rest/internal/book/http/controllers"
)

func RegisterBookRoutes(mainApp *app.App) {
	bookRepo := controllers.NewBookApp(mainApp.Database)

	mainApp.Router.HandleFunc("/book/author={author}", bookRepo.SearchByAuthor).Methods("GET", "OPTIONS")
	mainApp.Router.HandleFunc("/book/title={title}", bookRepo.SearchByTitle).Methods("GET", "OPTIONS")
	mainApp.Router.HandleFunc("/book/{id}", bookRepo.GetBook).Methods("GET", "OPTIONS")
	mainApp.Router.HandleFunc("/book", bookRepo.GetAllBooks).Methods("GET", "OPTIONS")

	mainApp.Router.HandleFunc("/book", bookRepo.SaveBook).Methods("POST", "OPTIONS")
	mainApp.Router.HandleFunc("/book", bookRepo.UpdateBook).Methods("PUT", "OPTIONS")
	mainApp.Router.HandleFunc("/book", bookRepo.DeleteBook).Methods("DELETE", "OPTIONS")
}
