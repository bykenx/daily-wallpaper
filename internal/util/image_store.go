package util

import (
	"daily-wallpaper/internal/config"
	"daily-wallpaper/internal/db"
	"os"
	"path/filepath"
)

func GetOrDownload(url string) (string, error) {
	var item db.ImageItem
	result := db.DB().Where(&db.ImageItem{Url: url, Valid: true}).First(&item)
	if result.Error == nil {
		if IsFile(item.Dir) {
			return item.Dir, nil
		}
		item.Valid = false
		db.DB().Updates(item)
	}

	bytes, suffix, err := downloadSource(url)
	if err != nil {
		return "", DownloadError{msg: err.Error()}
	}
	category := knownSuffixCategory[suffix]
	if category == "" {
		return "", DownloadError{msg: "不支持保存的文件类型"}
	}
	fileName, shortHash := buildImageFileName(url, suffix)
	categoryPath := filepath.Join(config.AppHome, category)
	MkdirIfNotExists(categoryPath)
	filePath := filepath.Join(categoryPath, fileName)
	_ = os.WriteFile(filePath, bytes, config.DefaultFileCreatePermission)
	db.DB().Save(&db.ImageItem{
		Url:  url,
		Dir:  filePath,
		Sha1: shortHash,
	})
	return filePath, nil
}

func GetImageListPagination(start, limit int) []string {
	if limit <= 0 {
		limit = 10
	}
	offset := start * limit
	var list []db.ImageItem
	resultList := make([]string, 0)
	db.DB().Where(&db.ImageItem{Valid: true}).Order("create_time DESC").Limit(limit).Offset(offset).Find(&list)
	for _, item := range list {
		resultList = append(resultList, item.Url)
	}
	return resultList
}
