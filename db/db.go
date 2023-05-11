package db

import (
	"crypto/sha1"
	"daily-wallpaper/constant"
	"daily-wallpaper/models"
	"daily-wallpaper/utils"
	"encoding/hex"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(filepath.Join(constant.AppHome, "app.db")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	if db.Migrator().HasTable("download_histories") {
		db.Migrator().RenameTable("download_histories", &models.ImageItem{})
	} else {
		err = db.AutoMigrate(&models.ImageItem{})
	}

	{
		// 迁移数据
		var list []models.ImageItem

		db.Where("sha1 is null OR trim(sha1) = '' AND valid = 1").Find(&list)

		for _, item := range list {
			println(item.Dir)
			if utils.IsFile(item.Dir) {
				data, _ := os.ReadFile(item.Dir)
				sum := sha1.Sum(data)
				item.Sha1 = hex.EncodeToString(sum[:])
			} else {
				item.Valid = false
			}
			db.Updates(item)
		}
	}

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
