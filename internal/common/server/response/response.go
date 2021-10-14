package response

import (
	"github.com/gin-gonic/gin"
)

type genericResponse struct {
	DataField interface{} `json:"data,omitempty"`
	ErrorField interface{} `json:"error,omitempty"`
}

func Respond(c *gin.Context, status int, respData interface{}, respError interface{}) {
	c.JSON(status, genericResponse{
		DataField:  respData,
		ErrorField: respError,
	})
}
