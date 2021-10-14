package errors

import (
	common_response "github.com/foxfurry/medialib/internal/common/server/common_response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonError struct {
	Msg string	`json:"msg"`
}

func (ce CommonError) Error() string {
	return ce.Msg
}

func respondWithError(c *gin.Context, status int, err error) {
	common_response.Respond(c, status, nil, err)
}

func RespondNotFound(c *gin.Context, err error) {
	respondWithError(c, http.StatusNotFound, err)
}

func RespondInternalError(c *gin.Context, err error) {
	respondWithError(c, http.StatusInternalServerError, err)
}

func RespondBadRequest(c *gin.Context, err error) {
	respondWithError(c, http.StatusBadRequest, err)
}

func RespondAlreadyExists(c *gin.Context, err error) {
	respondWithError(c, http.StatusConflict, err)
}
