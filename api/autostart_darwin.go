//+build darwin

package api

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "autostart.h"
*/
import "C"

func SetStartAtLogin(startAtLogin bool) bool {
	var startAtLoginVal C.BOOL
	if startAtLogin {
		startAtLoginVal = 1
	} else {
		startAtLoginVal = 0
	}
	return C.setStartAtLogin(C.BOOL(startAtLoginVal)) == 1
}
