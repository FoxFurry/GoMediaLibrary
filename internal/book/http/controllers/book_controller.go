package controllers

import (
	"database/sql"
	"encoding/json"
	bookDB "github.com/foxfurry/simple-rest/internal/book/db"
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"github.com/foxfurry/simple-rest/internal/book/http/errors"
	"github.com/foxfurry/simple-rest/internal/common/server"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type BookApp struct {
	repository bookDB.BookDBRepository
}

func NewBookApp(db *sql.DB) BookApp{
	return BookApp{
		repository: bookDB.NewBookRepo(db),
	}
}

func (b *BookApp) SaveBook(w http.ResponseWriter, r *http.Request){
	var book entity.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		server.RespondBadRequest(w, err.Error())
		return
	}

	saveBook, err := b.repository.SaveBook(&book)
	if err != nil {
		server.RespondInternalError(w, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(saveBook); err != nil {
		server.RespondInternalError(w, err.Error())
		return
	}

}

func (b *BookApp) GetBook(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		server.RespondInternalError(w, err.Error())
	}

	getBook, err := b.repository.GetBook(uint64(id))
	if err != nil {
		server.RespondInternalError(w, err.Error())
	}

	if err := json.NewEncoder(w).Encode(getBook); err != nil {
		server.RespondInternalError(w, err.Error())
		return
	}
}

func (b *BookApp) GetAllBooks(w http.ResponseWriter, r *http.Request){
	books, err := b.repository.GetAllBooks()
	if err == sql.ErrNoRows {
		errors.BookNotFound(w)
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		server.RespondInternalError(w, err.Error())
		return
	}

}

func (b *BookApp) SearchByAuthor(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	author := params["author"]

	byAuthor, err := b.repository.SearchByAuthor(author)
	if err != nil {
		server.RespondInternalError(w, err.Error())
	}

	if err := json.NewEncoder(w).Encode(byAuthor); err != nil {
		server.RespondInternalError(w, err.Error())
		return
	}
}

func (b *BookApp) SearchByTitle(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	author := params["title"]

	byTitle, err := b.repository.SearchByTitle(author)
	if err != nil {
		server.RespondInternalError(w, err.Error())
	}

	if err := json.NewEncoder(w).Encode(byTitle); err != nil {
		server.RespondInternalError(w, err.Error())
		return
	}
}

