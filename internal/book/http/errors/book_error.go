package errors

import (
	"github.com/foxfurry/simple-rest/internal/common/server"
	"net/http"
)

func BookNotFound(w http.ResponseWriter) {
	server.RespondNotFound(w, "Specified book not found")
}


