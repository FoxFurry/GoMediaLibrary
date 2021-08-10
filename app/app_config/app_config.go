package app_config

import (
	"fmt"
	"github.com/foxfurry/simple-rest/app"
	"github.com/foxfurry/simple-rest/app/env_setup"
	"github.com/foxfurry/simple-rest/internal/book/http/routers"
	dbpool "github.com/foxfurry/simple-rest/internal/common/database"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
)


func NewApp() *app.App {
	env_setup.Config()
	newApp := &app.App{
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

	newApp.Router.Use(setGlobalHeaders)
	routers.RegisterBookRoutes(newApp)

	showRoutes(newApp.Router)
	return newApp
}





func setGlobalHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
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
