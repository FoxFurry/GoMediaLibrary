package db

import (
	"database/sql"
	"github.com/foxfurry/simple-rest/configs"
	dbpool "github.com/foxfurry/simple-rest/internal/common/database"
	_ "github.com/foxfurry/simple-rest/internal/common/testing"
	"github.com/spf13/viper"
	"testing"
)

var db *sql.DB

func init(){
	configs.LoadConfig()

	db = dbpool.CreateDBPool(
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetInt("database.maxidleconnections"),
		viper.GetInt("database.maxopenconnections"),
		viper.GetDuration("database.maxconnidletime"),
	)
}

func TestBookDBRepository_GetBook(t *testing.T) {

	t.Errorf("Sucking fucking")
}