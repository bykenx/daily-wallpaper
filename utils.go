package main

import (
	"daily-wallpaper/sources"
	"daily-wallpaper/sources/bing_source"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
)

func OpenUrl(url string) {
	platform := runtime.GOOS
	switch platform {
	case "windows":
		_ = exec.Command(`cmd`, `/c`, `start`, url).Start()
	case "darwin":
		_ = exec.Command(`open`, url).Start()
	case "linux":
		_ = exec.Command(`xdg-open`, url).Start()
	default:
		log.Fatal("Unsupported platform.")
	}
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func IsDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func MkdirIfNotExists(path string) {
	if IsDirExists(path) {
		return
	}
	_ = os.MkdirAll(path, defaultDirectoryCreatePermission)
}

func ginJsonError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": 400,
		"msg":  msg,
		"data": nil,
	})
}

func ginJsonResult(c *gin.Context, obj interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": obj,
	})
}

func GetSource(name string) sources.Source {
	switch name {
	case "bing":
		return bing_source.BingSource{}
	default:
		return nil
	}
}

func GetDescriptions() []sources.Description {
	return []sources.Description{
		{Name: "bing壁纸", Description: "bing壁纸", Enabled: true},
	}
}
