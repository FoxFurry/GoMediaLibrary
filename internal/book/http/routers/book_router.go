package routers

import (
	"github.com/foxfurry/simple-rest/app"
	"github.com/foxfurry/simple-rest/internal/book/http/controllers"
)

func RegisterBookRoutes(mainApp *app.App){
	bookRepo := controllers.NewBookApp(mainApp.Database)

	mainApp.Router.HandleFunc("/api/book/author={author}", bookRepo.SearchByAuthor).Methods("GET", "OPTIONS")
	mainApp.Router.HandleFunc("/api/book/title={title}", bookRepo.SearchByTitle).Methods("GET", "OPTIONS")
	mainApp.Router.HandleFunc("/api/book/{id}", bookRepo.GetBook).Methods("GET", "OPTIONS")
	mainApp.Router.HandleFunc("/api/book", bookRepo.GetAllBooks).Methods("GET", "OPTIONS")

	mainApp.Router.HandleFunc("/api/newuser", bookRepo.SaveBook).Methods("POST", "OPTIONS")

}
