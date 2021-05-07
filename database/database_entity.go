package database

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/yjagdale/siem-data-producer/models/health_models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var databaseConnection *gorm.DB

func GetDBConnection() (*gorm.DB, error) {
	return connectDB()
}

func connectDB() (*gorm.DB, error) {
	var err error
	var dbPath string
	log.Infoln("Connecting to database")
	dbPath = os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database.db"
	} else {
		log.Infoln("DB path provided in env, Using", dbPath)
		dbPath = dbPath + "/database.db"
	}
	databaseConnection, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	log.Debugln("Connected ", err == nil)
	return databaseConnection, err
}

func ValidateConnection() bool {
	if databaseConnection == nil {
		_, err := connectDB()
		if err != nil {
			log.Errorln("Error while connecting to database", err)
			return false
		} else {
			log.Debugln("Connection created successfully")
		}
	}
	health := health_models.Health{}
	err := databaseConnection.Model(health_models.Health{}).Select(&health).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorln("Error while checking health", err)
		return false
	}
	log.Debugln("DB health validated successfully")
	return true
}
