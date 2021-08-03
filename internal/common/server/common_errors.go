package server

import "net/http"

type ErrorResponse struct {
	errMsg  string
	errType int
}

func respondWithError(w http.ResponseWriter, response ErrorResponse) {
	http.Error(w, response.errMsg, response.errType)
}

func RespondNotFound(w http.ResponseWriter, error string) {
	respondWithError(w, ErrorResponse{errMsg: error, errType: 404})
}

func RespondInternalError(w http.ResponseWriter, error string) {
	respondWithError(w, ErrorResponse{errMsg: error, errType: 500})
}

func RespondBadRequest(w http.ResponseWriter, error string) {
	respondWithError(w, ErrorResponse{errMsg: error, errType: 400})
}
