package main

import (
	"daily-wallpaper/internal/app"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(app.OnReady, app.OnExit)
}
