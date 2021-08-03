package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type App struct {
	*http.Server
	Router   *mux.Router
	Database *sql.DB
}

func (a *App) Start() {
	log.Fatal(http.ListenAndServe(viper.GetString("server.port"), a.Router))
}
