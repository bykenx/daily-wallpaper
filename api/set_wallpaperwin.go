//+build windows

package api

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

var (
	user32DLL           = windows.NewLazyDLL("user32.dll")
	procSystemParamInfo = user32DLL.NewProc("SystemParametersInfoW")
)

func SetWallpaper(path string) error {
	imagePath, _ := windows.UTF16PtrFromString(path)
	_, _, err := procSystemParamInfo.Call(20, 0, uintptr(unsafe.Pointer(imagePath)), 0x001a)
	if err != nil && err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
