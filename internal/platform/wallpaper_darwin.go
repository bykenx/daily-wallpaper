package platform

import (
	"errors"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit -framework Foundation -framework ApplicationServices
#import <AppKit/AppKit.h>
#import <ApplicationServices/ApplicationServices.h>
#import <Foundation/Foundation.h>
#include <stdlib.h>

static BOOL screenIsLocked(void) {
	CFDictionaryRef session = CGSessionCopyCurrentDictionary();
	if (session == NULL) {
		return NO;
	}
	CFBooleanRef locked = (CFBooleanRef)CFDictionaryGetValue(session, CFSTR("CGSSessionScreenIsLocked"));
	BOOL isLocked = (locked == kCFBooleanTrue);
	CFRelease(session);
	return isLocked;
}

char* setWallpaperForAllScreens(const char* path) {
	@autoreleasepool {
		if (screenIsLocked()) {
			return strdup("屏幕已锁定，稍后重试");
		}

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
