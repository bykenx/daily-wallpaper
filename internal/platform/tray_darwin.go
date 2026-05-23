//go:build darwin

package platform

import "github.com/getlantern/systray"

func SetTrayIcon(iconData []byte) {
	systray.SetTemplateIcon(iconData, iconData)
}
