package routers

import (
	"github.com/gorilla/mux"
)

func RegisterBookRoutes(router *mux.Router){

	router.HandleFunc("/api/book/{id}", handlers.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/book", handlers.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newuser", handlers.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/book/{id}", handlers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteuser/{id}", handlers.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}
