package routers

import (
	"database/sql"
	"github.com/foxfurry/simple-rest/internal/book/http/controllers"
	"github.com/gorilla/mux"
)

func RegisterBookRoutes(router *mux.Router, db *sql.DB) {
	bookRepo := controllers.NewBookApp(db)

	router.HandleFunc("/book/author=\"{author}\"", bookRepo.SearchByAuthor).Methods("GET", "OPTIONS")
	router.HandleFunc("/book/title=\"{title}\"", bookRepo.SearchByTitle).Methods("GET", "OPTIONS")
	router.HandleFunc("/book/{id}", bookRepo.GetBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/book", bookRepo.GetAllBooks).Methods("GET", "OPTIONS")

	router.HandleFunc("/book", bookRepo.SaveBook).Methods("POST", "OPTIONS")
	router.HandleFunc("/book", bookRepo.UpdateBook).Methods("PUT", "OPTIONS")
	router.HandleFunc("/book", bookRepo.DeleteAllBooks).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/book/{id}", bookRepo.DeleteBook).Methods("DELETE", "OPTIONS")
}
