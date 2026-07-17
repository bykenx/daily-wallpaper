//go:build linux && cgo

package platform

import "github.com/getlantern/systray"

func SetTrayIcon(iconData []byte) {
	systray.SetIcon(iconData)
}
