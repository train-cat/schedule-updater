package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Load MySQL driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var db *gorm.DB

// InitDatabase connection
func InitDatabase() {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	hostname := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	name := viper.GetString("database.name")

	d, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, name))

	if err != nil {
		log.Fatal(err.Error())
	}

	db = d

	db.SingularTable(true)
}
