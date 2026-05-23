package util

import "os/exec"

func OpenUrl(url string) {
	_ = exec.Command(`xdg-open`, url).Start()
}
