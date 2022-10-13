package api

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func SetWallpaper(path string) error {
	// just test on gnome
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
	cmd = strings.Replace(cmd, "\\", "\\\\", -1)
	return exec.Command("/bin/sh", "-c", cmd).Run()
}
