package app

import (
	"database/sql"
	"fmt"
	"github.com/foxfurry/simple-rest/internal/book/http/routers"
	dbpool "github.com/foxfurry/simple-rest/internal/common/database"
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

func NewApp() *App {
	loadConfig()

	newApp := &App{
		Router: mux.NewRouter(),
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

	newApp.Router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		})
	})

	routers.RegisterBookRoutes(newApp.Router, newApp.Database)

	showRoutes(newApp.Router)

	return newApp
}

func loadConfig() {
	viper.SetConfigName("environment")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal reading environment: %+v", err)
	}
}

func showRoutes(r *mux.Router) {
	log.Println("Registered routes:")

	walkFunc := func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, errPath := route.GetPathTemplate()
		method, errMethod := route.GetMethods()
		if errPath != nil && errMethod != nil {
			return fmt.Errorf("error reading path or methods: %v %v", errPath, errMethod)
		} else {
			log.Printf("%s %+v", path, method)
		}
		return nil
	}

	if err := r.Walk(walkFunc); err != nil {
		log.Printf("Logging error: %v", err)
	}
}