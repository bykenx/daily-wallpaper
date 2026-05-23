//go:build !darwin && !linux

package platform

import "github.com/getlantern/systray"

func SetTrayIcon(iconData []byte) {
	systray.SetIcon(iconData)
}
