package utils

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

func OpenUrl(url string) {
	verbPtr, _ := syscall.UTF16PtrFromString("open")
	filePtr, _ := syscall.UTF16PtrFromString("cmd")
	argsPtr, _ := syscall.UTF16PtrFromString(fmt.Sprintf("/c start %s", url))
	windows.ShellExecute(0, verbPtr, filePtr, argsPtr, nil, 0) // hide window
}
