package app

import (
	"database/sql"
	"flag"
	"fmt"
	dbpool "github.com/foxfurry/simple-rest/internal/infrastructure/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const (
	configFilePathUsage = "Config file directory. Config must be of name 'conv_{env}.yaml'."
	configDefaultPath   = "./configs"
	configFileFlagName  = "configFilePath"

	envUsage    = "Environment for dev, prod or test."
	envDefault  = "dev"
	envFlagName = "env"
)

type app struct {
	*http.Server
	router *mux.Router
	database *sql.DB
}

var configFilePath string
var env string

func NewApp() *app {
	config()

	newApp := &app{
		router: mux.NewRouter(),
	}

	return newApp
}

func (a *app) Start() {
	log.Fatal(http.ListenAndServe(viper.GetString("server.port"), a.router))
}

func config() {
	dbpool.ConfigDBPool(
		viper.GetString("database.URL"),
		viper.GetInt("database.maxIdleConnections"),
		viper.GetInt("database.maxOpenConnections"),
		viper.GetDuration("database.maxConnIdleTime"),
	)

	flag.StringVar(&configFilePath, configFileFlagName, configDefaultPath, configFilePathUsage)
	flag.StringVar(&env, envFlagName, envDefault, envUsage)
	yamlConfig(configFilePath, env)
}

func yamlConfig(path string, env string) {
	if flag.Lookup("test.v") != nil {
		env = "test"
		path = "./../../configs"
	}

	log.Println("Environment settings: " + env + "\nConfig path " + path)

	viper.SetConfigFile("conv_" + env)
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Fatal reading config: %+v", err)
	}
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
