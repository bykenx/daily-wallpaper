package api

import (
	"daily-wallpaper/constant"
	"daily-wallpaper/db"
	"daily-wallpaper/models"
	"daily-wallpaper/sources"
	"daily-wallpaper/utils"
	"fmt"
	"io/ioutil"
	"net/http"
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

type DownloadError struct {
	msg string
	error
}

func (e DownloadError) Error() string {
	return e.msg
}

func GetOrDownload(url string) (string, error) {
	var model models.ImageItem
	result := db.DB().Where(&models.ImageItem{Url: url, Valid: true}).First(&model)
	if result.Error == nil {
		if utils.IsFile(model.Dir) {
			return model.Dir, nil
		} else {
			model.Valid = false
			db.DB().Updates(model)
		}
	}
	var newModel models.ImageItem
	bytes, suffix, err := downloadSource(url)
	if err != nil {
		return "", DownloadError{msg: err.Error()}
	}
	category := knownSuffixCategory[suffix]
	if category == "" {
		return "", DownloadError{msg: "不支持保存的文件类型"}
	}
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), suffix)
	categoryPath := filepath.Join(constant.AppHome, category)
	utils.MkdirIfNotExists(categoryPath)
	filePath := filepath.Join(categoryPath, fileName)
	_ = ioutil.WriteFile(filePath, bytes, constant.DefaultFileCreatePermission)
	newModel.Url = url
	newModel.Dir = filePath
	db.DB().Save(&newModel)
	return filePath, nil
}

func GetImageListPagination(start, limit int) []string {
	if limit <= 0 {
		limit = 10
	}
	offset := start * limit
	var list []models.ImageItem
	var resultList []string
	db.DB().Limit(limit).Offset(offset).Find(&list)
	for _, item := range list {
		resultList = append(resultList, item.Url)
	}
	return resultList
}

func downloadSource(url string) ([]byte, string, error) {
	res, err := http.Get(sources.GetSafeUrl(url))
	if err != nil {
		return nil, "", DownloadError{msg: err.Error()}
	}
	if res.StatusCode != 200 {
		return nil, "", DownloadError{msg: "链接请求错误"}
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", DownloadError{msg: err.Error()}
	}
	metaTye := res.Header.Get("content-type")
	suffix := knownMediaSuffix[metaTye]
	if suffix == "" {
		return nil, "", DownloadError{msg: "不支持的资源类型"}
	}
	return bytes, suffix, nil
}
