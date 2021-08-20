package db

import (
	"daily-wallpaper/constant"
	"daily-wallpaper/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
)

var db *gorm.DB

func OpenDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(filepath.Join(constant.AppHome, "app.db")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}
	err = db.AutoMigrate(&models.DownloadHistory{})
	if err != nil {
		panic("failed to migrate db model")
	}
}

func CloseDB() {
	if db == nil {
		return
	}
	s, err := db.DB()
	if err != nil {
		return
	}
	_ = s.Close()
}

func DB() *gorm.DB {
	if db == nil {
		panic("call OpenDB first.")
	}
	return db
}
