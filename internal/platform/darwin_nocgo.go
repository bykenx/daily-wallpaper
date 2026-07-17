//go:build darwin && !cgo

package platform

import (
	"errors"
	"log"
	"os/exec"
)

func SetStartAtLogin(startAtLogin bool) bool {
	log.Println("[autostart] CGO is disabled; login item is unavailable on macOS")
	return false
}

func SetWallpaper(path string) error {
	return errors.New("CGO is disabled; cannot set wallpaper on macOS")
}

func SetTrayIcon(iconData []byte) {
}

func RunTray(onReady, onExit func()) {
}

func OpenUrl(url string) {
	_ = exec.Command("open", url).Start()
}
