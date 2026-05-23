//go:build windows

package platform

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Matches bundle id style from scripts/build.py; must stay stable across versions.
const scheduledTaskName = "DailyWallpaper_com.bykenx.daily-wallpaper"

func schtasksExe() string {
	root := os.Getenv("SystemRoot")
	if root == "" {
		root = os.Getenv("windir")
	}
	if root != "" {
		p := filepath.Join(root, "System32", "schtasks.exe")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "schtasks.exe"
}

func resolvedExecutable() (string, bool) {
	p, err := os.Executable()
	if err != nil {
		return "", false
	}
	p = filepath.Clean(p)
	if r, err := filepath.EvalSymlinks(p); err == nil {
		p = filepath.Clean(r)
	}
	return p, true
}

func scheduledTaskExists() bool {
	cmd := exec.Command(schtasksExe(), "/Query", "/TN", scheduledTaskName, "/FO", "LIST")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run() == nil
}

func SetStartAtLogin(startAtLogin bool) bool {
	if startAtLogin {
		exe, ok := resolvedExecutable()
		if !ok {
			return false
		}
		args := []string{
			"/Create",
			"/SC", "ONLOGON",
			"/TN", scheduledTaskName,
			"/TR", exe,
			"/RL", "LIMITED",
			"/F",
		}
		cmd := exec.Command(schtasksExe(), args...)
		return cmd.Run() == nil
	}
	if !scheduledTaskExists() {
		return true
	}
	cmd := exec.Command(schtasksExe(), "/Delete", "/TN", scheduledTaskName, "/F")
	return cmd.Run() == nil
}
