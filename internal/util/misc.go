package util

import (
	"daily-wallpaper/internal/config"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var platform = runtime.GOOS

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
	_ = os.MkdirAll(path, config.DefaultDirectoryCreatePermission)
}

func GetLastModifyTime(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return stat.ModTime().Unix()
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
