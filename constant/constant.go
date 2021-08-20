package constant

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	PackageName                      = "daily-wallpaper"
	ConfigFileName                   = "config.yml"
	DefaultDirectoryCreatePermission = 0755
	DefaultFileCreatePermission      = 0644
)

var AppHome string

func init() {
	home, _ := os.UserHomeDir()
	AppHome = filepath.Join(home, fmt.Sprintf(".%s", PackageName))
}