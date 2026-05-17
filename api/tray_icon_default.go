//go:build !darwin && !linux

package api

import "github.com/getlantern/systray"

func SetTrayIcon(iconData []byte) {
	systray.SetIcon(iconData)
}
