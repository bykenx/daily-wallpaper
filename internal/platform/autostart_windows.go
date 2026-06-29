//go:build windows

package platform

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/sys/windows"
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

func currentUserID() (string, bool) {
	u, err := user.Current()
	if err != nil {
		return "", false
	}
	return u.Username, true
}

func decodeConsoleOutput(b []byte) string {
	s := strings.TrimSpace(string(b))
	if s == "" {
		return ""
	}
	if utf8.Valid(b) {
		return s
	}
	if decoded, ok := decodeCodePage(b, windows.GetACP()); ok {
		return decoded
	}
	return s
}

func decodeCodePage(b []byte, cp uint32) (string, bool) {
	if len(b) == 0 || cp == 0 {
		return "", false
	}
	n, err := windows.MultiByteToWideChar(cp, 0, &b[0], int32(len(b)), nil, 0)
	if err != nil || n == 0 {
		return "", false
	}
	buf := make([]uint16, n)
	if _, err := windows.MultiByteToWideChar(cp, 0, &b[0], int32(len(b)), &buf[0], n); err != nil {
		return "", false
	}
	return strings.TrimSpace(windows.UTF16ToString(buf)), true
}

func xmlEscape(s string) string {
	var buf strings.Builder
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

func scheduledTaskXML(exe, userID string) string {
	exe = xmlEscape(exe)
	userID = xmlEscape(userID)
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Description>Daily Wallpaper autostart</Description>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
      <UserId>%s</UserId>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <UserId>%s</UserId>
      <LogonType>InteractiveToken</LogonType>
      <RunLevel>LeastPrivilege</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>%s</Command>
    </Exec>
  </Actions>
</Task>`, userID, userID, exe)
}

func writeUTF16XML(path, content string) error {
	u16 := utf16.Encode([]rune(content))
	b := make([]byte, 2+len(u16)*2)
	b[0], b[1] = 0xFF, 0xFE
	for i, r := range u16 {
		b[2+i*2] = byte(r)
		b[2+i*2+1] = byte(r >> 8)
	}
	return os.WriteFile(path, b, 0600)
}

func runSchtasks(quiet bool, args ...string) error {
	cmd := exec.Command(schtasksExe(), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := decodeConsoleOutput(stdout.Bytes())
	errOut := decodeConsoleOutput(stderr.Bytes())
	if err != nil {
		msg := errOut
		if msg == "" {
			msg = out
		}
		if msg != "" {
			log.Printf("[autostart] schtasks failed: %s", msg)
		} else {
			log.Printf("[autostart] schtasks failed: %v", err)
		}
		return err
	}
	if !quiet && out != "" {
		log.Printf("[autostart] schtasks: %s", out)
	}
	return nil
}

func scheduledTaskExists() bool {
	return runSchtasks(true, "/Query", "/TN", scheduledTaskName, "/FO", "LIST") == nil
}

func createScheduledTask(exe string) error {
	userID, ok := currentUserID()
	if !ok {
		return fmt.Errorf("could not resolve current user")
	}
	xmlPath := filepath.Join(os.TempDir(), scheduledTaskName+".xml")
	defer os.Remove(xmlPath)
	if err := writeUTF16XML(xmlPath, scheduledTaskXML(exe, userID)); err != nil {
		return err
	}
	return runSchtasks(false, "/Create", "/TN", scheduledTaskName, "/XML", xmlPath, "/F")
}

func deleteScheduledTask() error {
	if !scheduledTaskExists() {
		return nil
	}
	return runSchtasks(false, "/Delete", "/TN", scheduledTaskName, "/F")
}

func SetStartAtLogin(startAtLogin bool) bool {
	if startAtLogin {
		exe, ok := resolvedExecutable()
		if !ok {
			log.Printf("[autostart] enable failed: could not resolve executable")
			return false
		}
		if err := createScheduledTask(exe); err != nil {
			log.Printf("[autostart] enable failed")
			return false
		}
		log.Printf("[autostart] enable succeeded")
		return true
	}
	if !scheduledTaskExists() {
		return true
	}
	if err := deleteScheduledTask(); err != nil {
		log.Printf("[autostart] disable failed")
		return false
	}
	log.Printf("[autostart] disable succeeded")
	return true
}
