package utils

import (
	"crypto/sha1"
	"daily-wallpaper/internal/config"
	"daily-wallpaper/internal/db"
	"daily-wallpaper/internal/source"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	knownMediaSuffix = map[string]string{
		"image/png":  ".png",
		"image/jpg":  ".jpg",
		"image/jpeg": ".jpg",
	}
	knownSuffixCategory = map[string]string{
		".png": "images",
		".jpg": "images",
	}
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
		return "", err
	}
	category := knownSuffixCategory[suffix]
	if category == "" {
		return "", errors.New("不支持保存的文件类型")
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

func buildImageFileName(url, suffix string) (string, string) {
	sum := sha1.Sum([]byte(url))
	shortHash := hex.EncodeToString(sum[:])[:7]
	return fmt.Sprintf("%s_%s%s", time.Now().Format("200601021504"), shortHash, suffix), shortHash
}

func downloadSource(url string) ([]byte, string, error) {
	res, err := http.Get(source.GetSafeUrl(url))
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, "", errors.New("链接请求错误")
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	mediaType := res.Header.Get("content-type")
	suffix := knownMediaSuffix[mediaType]
	if suffix == "" {
		return nil, "", errors.New("不支持的资源类型")
	}
	return bytes, suffix, nil
}
