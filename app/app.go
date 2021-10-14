package app

import (
	"database/sql"
	"github.com/foxfurry/medialib/internal/book/http/router"
	dbpool "github.com/foxfurry/medialib/internal/common/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type IApp interface {
	Start()
}

// mediaApp structure is the core of the project.
// It embeds http server and provides router and database instances
type mediaApp struct {
	server *http.Server
	Router   *gin.Engine
	Database *sql.DB
}

// NewApp returns an instance of mediaApp with configured router and database.
// Configuration is loaded from viper environment
func NewApp() IApp {
	newApp := &mediaApp{
		Router: gin.New(),
		Database: dbpool.Create(
			viper.GetString("Database.host"),
			viper.GetInt("Database.port"),
			viper.GetString("Database.user"),
			viper.GetString("Database.password"),
			viper.GetString("Database.dbname"),
			viper.GetInt("Database.maxidleconnections"),
			viper.GetInt("Database.maxopenconnections"),
			viper.GetDuration("Database.maxconnidletime"),
		),
	}

	router.RegisterBook(newApp.Router, newApp.Database)

	return newApp
}

func NewTestApp() IApp {
	newApp := &mediaApp{
		Router: gin.New(),
		Database: dbpool.Create(
			viper.GetString("database_test.host"),
			viper.GetInt("database_test.port"),
			viper.GetString("database_test.user"),
			viper.GetString("database_test.password"),
			viper.GetString("database_test.dbname"),
			viper.GetInt("database_test.maxidleconnections"),
			viper.GetInt("database_test.maxopenconnections"),
			viper.GetDuration("database_test.maxconnidletime"),
		),
	}

	router.RegisterBook(newApp.Router, newApp.Database)

	return newApp
}

// Start allows mediaApp to serve a http server on port from environment
func (a *mediaApp) Start() {
	log.Fatal(http.ListenAndServe(viper.GetString("server.port"), a.Router))
}


