package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	packageName                      = "daily-wallpaper"
	configFileName                   = "config.yml"
	defaultDirectoryCreatePermission = 0755
	defaultFileCreatePermission      = 0644
)

var appHome string

func init() {
	home, _ := os.UserHomeDir()
	appHome = filepath.Join(home, fmt.Sprintf(".%s", packageName))
}