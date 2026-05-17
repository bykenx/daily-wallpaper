//go:build darwin

package api

import "github.com/getlantern/systray"

func SetTrayIcon(iconData []byte) {
	systray.SetTemplateIcon(iconData, iconData)
}
