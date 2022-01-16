package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DB_FILE = "./tinybeans_photos.db"
)

func GetDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{})
	if err != nil {
		log.Fatal("unable to open database")
	}

	db.AutoMigrate(&Entry{})
	db.AutoMigrate(&Like{})

	return db, nil
}
