package api

import (
	"daily-wallpaper/utils"
	"os"
	"path/filepath"
	"text/template"
)

const desktopTemplate = `[Desktop Entry]
Name={{.Name}}
GenericName={{.Name}}
Exec={{.Exec}}
Terminal=false
Type=Application
X-GNOME-Autostart-enabled=true
X-GNOME-Autostart-Delay=10`

func getXDGConfigPath() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return os.Getenv("XDG_CONFIG_HOME")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".config")
	}
}

func SetStartAtLogin(startAtLogin bool) bool {
	// 获取当前程序路径
	autostartConfigPath := filepath.Join(getXDGConfigPath(), "autostart")
	desktopPath := filepath.Join(autostartConfigPath, "daily-wallpaper.desktop")
	if startAtLogin {
		if !utils.IsDirExists(autostartConfigPath) {
			if err := os.MkdirAll(autostartConfigPath, 0755); err != nil {
				return false
			}
		}
		f, err := os.Create(desktopPath)
		if err != nil {
			return false
		}
		defer f.Close()
		exePath, err := os.Executable()
		if err != nil {
			return false
		}
		t := template.Must(template.New("desktop").Parse(desktopTemplate))
		if err := t.Execute(f, map[string]interface{}{
			"Name": "Daily Wallpaper",
			"Exec": exePath,
		}); err != nil {
			return false
		}
	} else {
		if utils.IsDirExists(desktopPath) {
			if err := os.Remove(desktopPath); err != nil {
				return false
			}
		}
	}
	return true
}
