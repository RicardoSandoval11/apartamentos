package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabase(dbString string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}

	client, err := db.DB()

	if err != nil {
		panic(fmt.Sprintf("failed to create db client: %s", err.Error()))
	}

	client.SetMaxIdleConns(5)
	client.SetMaxOpenConns(10)
	client.SetConnMaxLifetime(time.Hour)

	return db
}
