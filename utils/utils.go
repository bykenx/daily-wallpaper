package utils

import (
	"daily-wallpaper/constant"
	"daily-wallpaper/sources"
	"daily-wallpaper/sources/bing_source"
	"daily-wallpaper/sources/unsplash_source"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/sys/windows"
)

var platform = runtime.GOOS

func OpenUrl(url string) {
	switch platform {
	case "windows":
		verbPtr, _ := syscall.UTF16PtrFromString("open")
		filePtr, _ := syscall.UTF16PtrFromString("cmd")
		argsPtr, _ := syscall.UTF16PtrFromString(fmt.Sprintf("/c start %s", url))
		windows.ShellExecute(0, verbPtr, filePtr, argsPtr, nil, 0) // hide window
	case "darwin":
		_ = exec.Command(`open`, url).Start()
	case "linux":
		_ = exec.Command(`xdg-open`, url).Start()
	default:
		log.Println("Unsupported platform.")
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
	return err == nil
}

func MkdirIfNotExists(path string) {
	if IsDirExists(path) {
		return
	}
	_ = os.MkdirAll(path, constant.DefaultDirectoryCreatePermission)
}

func GinJsonError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": 400,
		"msg":  msg,
		"data": nil,
	})
}

func GinJsonResult(c *gin.Context, obj interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": obj,
	})
}

func GetSource(name string) sources.Source {
	switch name {
	case "bing":
		return bing_source.BingSource{}
	case "unsplash":
		return unsplash_source.UnsplashSource{}
	default:
		return nil
	}
}

func GetLastModifyTime(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return stat.ModTime().Unix()
}

func GetDescriptions() []sources.Description {
	return []sources.Description{
		{Name: "bing", Description: "bing壁纸", Enabled: true},
		{Name: "unsplash", Description: "Unsplash", Enabled: true},
	}
}

func GetCwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
	}
	return dir
}

func GetStaticPath() string {
	cwd := GetCwd()
	switch platform {
	case "windows":
		return filepath.Join(cwd, "static")
	case "darwin":
		return filepath.Join(filepath.Dir(cwd), "Resources", "static")
	default:
		return filepath.Join(cwd, "static")
	}
}
