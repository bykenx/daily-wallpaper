package main

import (
	"daily-wallpaper/internal/app"
	"daily-wallpaper/internal/platform"
)

func main() {
	platform.RunTray(app.OnReady, app.OnExit)
}