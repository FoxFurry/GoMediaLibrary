package server

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	errMsg  string
	errType int
}

func respondWithError(c *gin.Context, response ErrorResponse) {
	c.JSON(response.errType, gin.H{"error:":response.errMsg})
}

func RespondNotFound(c *gin.Context, error string) {
	respondWithError(c, ErrorResponse{errMsg: error, errType: 404})
}

func RespondInternalError(c *gin.Context, error string) {
	respondWithError(c, ErrorResponse{errMsg: error, errType: 500})
}

func RespondBadRequest(c *gin.Context, error string) {
	respondWithError(c, ErrorResponse{errMsg: error, errType: 400})
}

func RespondAlreadyExists(c *gin.Context, error string) {
	respondWithError(c, ErrorResponse{errMsg: error, errType: 403})
}
