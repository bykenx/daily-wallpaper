package api

import (
	"errors"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit -framework Foundation
#import <AppKit/AppKit.h>
#import <Foundation/Foundation.h>
#include <stdlib.h>

char* setWallpaperForAllScreens(const char* path) {
	@autoreleasepool {
		NSString *pathString = [NSString stringWithUTF8String:path];
		NSURL *imageURL = [NSURL fileURLWithPath:pathString];
		NSArray<NSScreen *> *screens = [NSScreen screens];
		NSWorkspace *workspace = [NSWorkspace sharedWorkspace];

		for (NSScreen *screen in screens) {
			NSError *error = nil;
			BOOL ok = [workspace setDesktopImageURL:imageURL forScreen:screen options:@{} error:&error];
			if (!ok) {
				NSString *message = error ? [error localizedDescription] : @"设置壁纸失败";
				return strdup([message UTF8String]);
			}
		}

		return NULL;
	}
}
*/
import "C"

func SetWallpaper(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	errMsg := C.setWallpaperForAllScreens(cPath)
	if errMsg != nil {
		defer C.free(unsafe.Pointer(errMsg))
		return errors.New(C.GoString(errMsg))
	}
	return nil
}
