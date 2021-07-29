package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	bookDB "github.com/foxfurry/simple-rest/internal/book/db"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type BookApp struct {
	repository bookDB.BookRepo
}

func NewBookApp(db *sql.DB) BookApp{
	return BookApp{
		repository: bookDB.NewBookRepo(db),
	}
}

func (b *BookApp) SaveBook(book *entity.Book) (*entity.Book, error){

}

func (b *BookApp) GetBook(uint64) (*entity.Book, error){

}

func (b *BookApp) GetAllBooks() ([]entity.Book, error){

}

func (b *BookApp) SearchByAuthor(author string) ([]entity.Book, error){

}

func (b *BookApp) SearchByTitle(title string) (*entity.Book, error){

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var book entity.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body: %v\n", err)
	}

	insertID :=

	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to get book id: %v\n", err)
	}

	user, err := usermodel.GetUser(int64(id))

	if err != nil {
		log.Fatalf("Unable to get book: %v\n", err)
	}

	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users, err := usermodel.GetAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all users: %v\n", err)
	}

	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	requestParams := mux.Vars(r)

	id, err := strconv.Atoi(requestParams["id"])

	if err != nil {
		log.Fatalf("Unable to get book id: %v\n", err)
	}

	var user usermodel.User

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body: %v\n", err)
	}

	updatedRows := usermodel.UpdateUser(int64(id), user)

	res := response{
		ID:      int64(id),
		Message: fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows),
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	requestParams := mux.Vars(r)

	id, err := strconv.Atoi(requestParams["id"])

	if err != nil {
		log.Fatalf("Unable to get book id: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
	}

	deletedRows := usermodel.DeleteUser(int64(id))

	res := response{
		ID:      int64(id),
		Message: fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows),
	}

	json.NewEncoder(w).Encode(res)
}
