package app

import (
	"database/sql"
	"github.com/foxfurry/simple-rest/internal/book/http/routers"
	dbpool "github.com/foxfurry/simple-rest/internal/common/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// app structure is the core of the project.
// It embeds http server and provides router and database instances
type app struct {
	*http.Server
	Router   *gin.Engine
	Database *sql.DB
}

// Start allows app to serve a http server on port from environment
func (a *app) Start() {
	log.Fatal(http.ListenAndServe(viper.GetString("server.port"), a.Router))
}

// NewApp returns an instance of app with configured router and database.
// Configuration is loaded from viper environment
func NewApp() *app {
	newApp := &app{
		Router: gin.New(),
		Database: dbpool.CreateDBPool(
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

	routers.RegisterBookRoutes(newApp.Router, newApp.Database)

	return newApp
}