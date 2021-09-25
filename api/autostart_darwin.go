//+build darwin

package api

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include <Cocoa/Cocoa.h>

BOOL isStartAtLogin()
{
    BOOL  isHaveAdd = NO;
    NSString* appPath = [[NSBundle mainBundle] bundlePath];
    LSSharedFileListRef loginItems = LSSharedFileListCreate(NULL, kLSSharedFileListSessionLoginItems, NULL);
    UInt32 seedValue = 0;
    NSArray* loginItemsArray = (NSArray*)LSSharedFileListCopySnapshot(loginItems, &seedValue);
    CFURLRef tempUrl = (CFURLRef)[NSURL fileURLWithPath:appPath];
    for(NSInteger i = 0 ; i< [loginItemsArray count]; i ++ )
    {
        LSSharedFileListItemRef itemRef = (LSSharedFileListItemRef)[loginItemsArray objectAtIndex:i];
        if (LSSharedFileListItemResolve(itemRef, 0,&tempUrl, NULL) == noErr)
        {
            NSString * urlPath = [(NSURL*)tempUrl path];
            if ([urlPath compare:appPath] == NSOrderedSame)
            {
                isHaveAdd = YES;
                break;
            }
        }
    }
    [loginItemsArray release];
    CFRelease(loginItems);
    return isHaveAdd;
}

BOOL setStartAtLogin(BOOL startAtLogin)
{
    NSString* appPath = [[NSBundle mainBundle] bundlePath];
    BOOL result = NO;
    if (startAtLogin)
    {
        if (!isStartAtLogin())
        {
            CFURLRef url = (CFURLRef)[NSURL fileURLWithPath:appPath];
            LSSharedFileListRef newloginItems = LSSharedFileListCreate(NULL, kLSSharedFileListSessionLoginItems, NULL);
            LSSharedFileListItemRef item = LSSharedFileListInsertItemURL(newloginItems, kLSSharedFileListItemLast, NULL, NULL, url, NULL, NULL);
            if (item)
            {
                result = YES;
                CFRelease(item);
            }
            if (newloginItems)
            {
                CFRelease(newloginItems);
            }
        }
    }
    else
    {
        LSSharedFileListRef loginItems = LSSharedFileListCreate(NULL,kLSSharedFileListSessionLoginItems, NULL);
        UInt32 seedValue = 0;
        NSArray* loginItemsArray = (NSArray*)LSSharedFileListCopySnapshot(loginItems, &seedValue);
        CFURLRef tempUrl = (CFURLRef)[NSURL fileURLWithPath:appPath];
        for(NSInteger i = 0 ; i < [loginItemsArray count]; i ++ )
        {
            LSSharedFileListItemRef itemRef = (LSSharedFileListItemRef)[loginItemsArray objectAtIndex:i];
            if (LSSharedFileListItemResolve(itemRef, 0,&tempUrl, NULL) == noErr)
            {
                NSString * urlPath = [(NSURL*)tempUrl path];
                if ([urlPath compare:appPath] == NSOrderedSame)
                {
                    OSStatus status = LSSharedFileListItemRemove(loginItems,itemRef);
                    result = (status == noErr);
                }
            }
        }
        [loginItemsArray release];
        CFRelease(loginItems);
    }
    return result;
}
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
