package app_config

import (
	"flag"
	"fmt"
	"github.com/foxfurry/simple-rest/app"
	"github.com/foxfurry/simple-rest/internal/book/http/routers"
	dbpool "github.com/foxfurry/simple-rest/internal/common/database"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

const (
	configFilePathUsage = "Config file directory. Config must be of name 'conv_{env}.yaml'."
	configDefaultPath   = "./configs"
	configFileFlagName  = "configFilePath"

	envUsage    = "Environment for dev, prod or test."
	envDefault  = "dev"
	envFlagName = "env"
)

var configFilePath string
var env string

func NewApp() *app.App {
	config()

	newApp := &app.App{
		Router:   mux.NewRouter(),
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

func config() {
	flag.StringVar(&configFilePath, configFileFlagName, configDefaultPath, configFilePathUsage)
	flag.StringVar(&env, envFlagName, envDefault, envUsage)
	flag.Parse()

	yamlConfig(configFilePath, env)
}

func yamlConfig(path string, env string) {
	if flag.Lookup("test.v") != nil {
		env = "test"
		path = "./configs"
	}

	log.Println("Environment settings: \"" + env + "\" | Config path \"" + path + "\"")

	viper.SetConfigName("conf_" + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Fatal reading config: %+v", err)
	}
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
			return fmt.Errorf("Error reading path or methods: %v %v", errPath, errMethod)
		} else {
			log.Printf("%s %+v", path, method)
		}
		return nil
	}

	if err := r.Walk(walkFunc); err != nil {
		log.Printf("Logging error: %v", err)
	}
}
