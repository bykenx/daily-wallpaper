package util

import (
	"crypto/sha1"
	"daily-wallpaper/internal/source"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
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

func buildImageFileName(url, suffix string) (string, string) {
	sum := sha1.Sum([]byte(url))
	shortHash := hex.EncodeToString(sum[:])[:7]
	return fmt.Sprintf("%s_%s%s", time.Now().Format("200601021504"), shortHash, suffix), shortHash
}

func downloadSource(url string) ([]byte, string, error) {
	res, err := http.Get(source.GetSafeUrl(url))
	if err != nil {
		return nil, "", DownloadError{msg: err.Error()}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, "", DownloadError{msg: "链接请求错误"}
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", DownloadError{msg: err.Error()}
	}
	mediaType := res.Header.Get("content-type")
	suffix := knownMediaSuffix[mediaType]
	if suffix == "" {
		return nil, "", DownloadError{msg: "不支持的资源类型"}
	}
	return bytes, suffix, nil
}
