package env_setup

import (
	"flag"
	"github.com/spf13/viper"
	"log"
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

func Config() {
	flag.StringVar(&configFilePath, configFileFlagName, configDefaultPath, configFilePathUsage)
	flag.StringVar(&env, envFlagName, envDefault, envUsage)
	flag.Parse()

	generalConfig(configFilePath, env)
}

func TestConfig() {
	env = "test"
	configFilePath = "./configs"

	viper.SetConfigName("conf_" + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configFilePath)
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Fatal reading Config: %+v", err)
	}
}

func generalConfig(path string, env string) {
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
		log.Fatalf("Fatal reading Config: %+v", err)
	}
}
