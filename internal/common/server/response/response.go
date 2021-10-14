package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type responseTemplate struct {
	DataField interface{} `json:"data,omitempty"`
	ErrorField interface{} `json:"error,omitempty"`
}

func respond(c *gin.Context, status int, respData interface{}, respError interface{}) {
	c.JSON(status, responseTemplate{
		DataField:  respData,
		ErrorField: respError,
	})
}

func OK(c *gin.Context, respData interface{}) {
	respond(c, http.StatusOK, respData, nil)
}

func NotFound(c *gin.Context, err error) {
	respond(c, http.StatusNotFound, nil, err)
}

func InternalError(c *gin.Context, err error) {
	respond(c, http.StatusInternalServerError, nil, err)
}

func BadRequest(c *gin.Context, err error) {
	respond(c, http.StatusBadRequest, nil, err)
}

func AlreadyExists(c *gin.Context, err error) {
	respond(c, http.StatusConflict, nil, err)
}