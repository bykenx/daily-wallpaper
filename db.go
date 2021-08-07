package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
)

var db *gorm.DB

func openDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(filepath.Join(appHome, "app.db")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}
	err = db.AutoMigrate(&DownloadHistory{})
	if err != nil {
		panic("failed to migrate db model")
	}
}

func closeDB() {
	s, err := db.DB()
	if err != nil {
		return
	}
	s.Close()
}
