package databases

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"distrise/internal/configs"
	"distrise/internal/models"
)

// set up a singleton instance of the database
var dbInstance *gorm.DB

func GetDB() (instance *gorm.DB, err error) {
	if dbInstance == nil {
		instance, err = newDB()
		if err != nil {
			return nil, err
		}
		dbInstance = instance
	}
	return dbInstance, nil
}

func newDB() (*gorm.DB, error) {
	dbConfig := configs.GetConfig().Database
	db, err := gorm.Open(
		postgres.Open(dbConfig.URL))
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&models.CoreEvent{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
