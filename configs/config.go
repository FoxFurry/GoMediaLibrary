package configs

import (
	"github.com/spf13/viper"
	"log"
)

// LoadConfig loads configuration from ./configs/environment.yaml file
func LoadConfig() {
	viper.SetConfigName("environment")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal reading environment: %+v", err)
	}
}
