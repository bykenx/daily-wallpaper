//go:build linux

package platform

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/getlantern/systray"
)

const desktopTemplate = `[Desktop Entry]
Name={{.Name}}
GenericName={{.Name}}
Exec={{.Exec}}
Terminal=false
Type=Application
X-GNOME-Autostart-enabled=true
X-GNOME-Autostart-Delay=10`

func SetStartAtLogin(startAtLogin bool) bool {
	autostartConfigPath := filepath.Join(getXDGConfigPath(), "autostart")
	desktopPath := filepath.Join(autostartConfigPath, "daily-wallpaper.desktop")
	if startAtLogin {
		if !pathExists(autostartConfigPath) {
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
		if pathExists(desktopPath) {
			if err := os.Remove(desktopPath); err != nil {
				return false
			}
		}
	}
	return true
}

func SetWallpaper(path string) error {
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	var cmd string
	if strings.Contains(desktop, "xfce") {
		cmd = "xfconf-query -c xfce4-desktop -p /backdrop/screen0/monitorVGA-1/workspace0/last-image -s \"" + path + "\""
	} else if strings.Contains(desktop, "gnome") {
		cmd = "gsettings set org.gnome.desktop.background picture-uri \"file://" + path + "\""
	} else if strings.Contains(desktop, "kde") {
		cmd = "qdbus org.kde.plasmashell /PlasmaShell org.kde.PlasmaShell.evaluateScript 'var allDesktops = desktops();print (allDesktops);for (i=0;i<allDesktops.length;i++) {d = allDesktops[i];d.wallpaperPlugin = \"org.kde.image\";d.currentConfigGroup = Array(\"Wallpaper\", \"org.kde.image\", \"General\");d.writeConfig(\"Image\", \"" + path + "\")}'"
	} else {
		log.Printf("不支持的桌面环境: %s", desktop)
		return nil
	}
	cmd = strings.ReplaceAll(cmd, "\\", "\\\\")
	return exec.Command("/bin/sh", "-c", cmd).Run()
}

func RunTray(onReady, onExit func()) {
	systray.Run(onReady, onExit)
}

func OpenUrl(url string) {
	_ = exec.Command("xdg-open", url).Start()
}

func getXDGConfigPath() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return os.Getenv("XDG_CONFIG_HOME")
	}
	return filepath.Join(os.Getenv("HOME"), ".config")
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
