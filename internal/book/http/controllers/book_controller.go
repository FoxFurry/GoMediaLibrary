package controllers

import (
	"database/sql"
	"encoding/json"
	bookDB "github.com/foxfurry/simple-rest/internal/book/db"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BookApp struct {
	dbRepo bookDB.BookDBRepository
}

func NewBookApp(db *sql.DB) BookApp{
	return BookApp{
		dbRepo: bookDB.NewBookRepo(db),
	}
}

func (b *BookApp) SaveBook(w http.ResponseWriter, r *http.Request){
	var book entity.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		errors.HandleBookError(w, err)
		return
	}

	saveBook, err := b.dbRepo.SaveBook(&book)
	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(saveBook); err != nil {
		errors.HandleBookError(w, err)
		return
	}

}

func (b *BookApp) GetBook(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	getBook, err := b.dbRepo.GetBook(uint64(id))
	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(getBook); err != nil {
		errors.HandleBookError(w, err)
		return
	}
}

func (b *BookApp) GetAllBooks(w http.ResponseWriter, r *http.Request){
	books, err := b.dbRepo.GetAllBooks()
	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(books); err != nil {
		errors.HandleBookError(w, err)
		return
	}

}

func (b *BookApp) SearchByAuthor(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	author := params["author"]

	byAuthor, err := b.dbRepo.SearchByAuthor(author)
	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(byAuthor); err != nil {
		errors.HandleBookError(w, err)
		return
	}
}

func (b *BookApp) SearchByTitle(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	author := params["title"]

	byTitle, err := b.dbRepo.SearchByTitle(author)
	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(byTitle); err != nil {
		errors.HandleBookError(w, err)
		return
	}
}

func (b *BookApp) UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	var book *entity.Book

	if err = json.NewDecoder(r.Body).Decode(book); err != nil {
		errors.HandleBookError(w, err)
		return
	}

	updatedRows, err := b.dbRepo.UpdateBook(uint64(id), book)
	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(updatedRows); err != nil {
		errors.HandleBookError(w, err)
		return
	}
}

func (b *BookApp) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	deletedRows, err := b.dbRepo.DeleteBook(uint64(id))
	if err != nil {
		errors.HandleBookError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(deletedRows); err != nil {
		errors.HandleBookError(w, err)
		return
	}
}