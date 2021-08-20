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

func DownloadFileAndSave(url string) (string, error) {
	history := models.DownloadHistory{}
	db.DB().Where(&models.DownloadHistory{Url: url}).First(&history)
	if history.Dir != "" && utils.IsFile(history.Dir) {
		return history.Dir, nil
	}
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
	history.Url = url
	history.Dir = filePath
	db.DB().Save(&history)
	return filePath, nil
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
