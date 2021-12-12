package utils

import (
	"os/exec"
)

func OpenUrl(url string) {
	_ = exec.Command(`open`, url).Start()
}
