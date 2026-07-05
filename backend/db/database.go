package db

import (
	"time"

	"github.com/RicardoSandoval11/apartamentos/backend/constants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(constants.DATABASE_CONN_STRING), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}

	client, _ := db.DB()

	client.SetMaxIdleConns(5)
	client.SetMaxOpenConns(10)
	client.SetConnMaxLifetime(time.Hour)

	return db
}
