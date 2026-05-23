//go:build darwin

package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework ServiceManagement -framework CoreServices
#include <Cocoa/Cocoa.h>
#include <CoreServices/CoreServices.h>
#include <ServiceManagement/ServiceManagement.h>
#include <unistd.h>

static NSString *const kLaunchAgentLabel = @"com.bykenx.daily-wallpaper";

static NSString *launchAgentPlistPath(void)
{
	return [NSString stringWithFormat:@"%@/Library/LaunchAgents/%@.plist", NSHomeDirectory(), kLaunchAgentLabel];
}

// Removes legacy session login items created by older versions (LSSharedFileList), only used on macOS 13+ before SMAppService.
static void removeLegacySessionLoginItems(void)
{
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
	NSString *appPath = [[NSBundle mainBundle] bundlePath];
	if (!appPath)
		return;
	LSSharedFileListRef loginItems = LSSharedFileListCreate(NULL, kLSSharedFileListSessionLoginItems, NULL);
	if (!loginItems)
		return;
	UInt32 seedValue = 0;
	NSArray *loginItemsArray = (NSArray *)LSSharedFileListCopySnapshot(loginItems, &seedValue);
	CFURLRef tempUrl = (CFURLRef)[NSURL fileURLWithPath:appPath];
	for (NSInteger i = 0; i < [loginItemsArray count]; i++)
	{
		LSSharedFileListItemRef itemRef = (LSSharedFileListItemRef)[loginItemsArray objectAtIndex:i];
		if (LSSharedFileListItemResolve(itemRef, 0, &tempUrl, NULL) == noErr)
		{
			NSString *urlPath = [(NSURL *)tempUrl path];
			if ([urlPath compare:appPath] == NSOrderedSame)
				LSSharedFileListItemRemove(loginItems, itemRef);
		}
	}
	[loginItemsArray release];
	CFRelease(loginItems);
#pragma clang diagnostic pop
}

static BOOL runLaunchctl(NSArray *arguments)
{
	NSTask *task = [[NSTask alloc] init];
	task.launchPath = @"/bin/launchctl";
	task.arguments = arguments;
	NSPipe *devNull = [NSPipe pipe];
	task.standardOutput = devNull;
	task.standardError = devNull;
	BOOL ok = NO;
	@try
	{
		[task launch];
		[task waitUntilExit];
		ok = ([task terminationStatus] == 0);
	}
	@catch (NSException *e)
	{
		ok = NO;
	}
	[task release];
	return ok;
}

static BOOL launchctlLoadPlist(NSString *plistPath)
{
	NSOperatingSystemVersion ver = [[NSProcessInfo processInfo] operatingSystemVersion];
	if (ver.majorVersion >= 11)
	{
		NSString *domain = [NSString stringWithFormat:@"gui/%u", getuid()];
		return runLaunchctl(@[ @"bootstrap", domain, plistPath ]);
	}
	return runLaunchctl(@[ @"load", @"-w", plistPath ]);
}

static void launchctlUnloadPlistIgnoringErrors(NSString *plistPath)
{
	NSOperatingSystemVersion ver = [[NSProcessInfo processInfo] operatingSystemVersion];
	if (ver.majorVersion >= 11)
	{
		NSString *domain = [NSString stringWithFormat:@"gui/%u", getuid()];
		runLaunchctl(@[ @"bootout", domain, plistPath ]);
		return;
	}
	runLaunchctl(@[ @"unload", @"-w", plistPath ]);
}

static void removeLaunchAgentFiles(void)
{
	NSString *plistPath = launchAgentPlistPath();
	NSFileManager *fm = [NSFileManager defaultManager];
	if (![fm fileExistsAtPath:plistPath])
		return;
	launchctlUnloadPlistIgnoringErrors(plistPath);
	[fm removeItemAtPath:plistPath error:NULL];
}

static BOOL isMainBundleAppBundle(void)
{
	NSString *path = [[NSBundle mainBundle] bundlePath];
	return path != nil && [path hasSuffix:@".app"];
}

static NSArray *launchAgentProgramArguments(void)
{
	if (isMainBundleAppBundle())
		return @[ @"/usr/bin/open", @"-gj", [[NSBundle mainBundle] bundlePath] ];
	NSString *exe = [[NSBundle mainBundle] executablePath];
	if (!exe)
		return nil;
	return @[ exe ];
}

static BOOL installLaunchAgent(void)
{
	NSArray *args = launchAgentProgramArguments();
	if (!args)
		return NO;
	NSString *plistPath = launchAgentPlistPath();
	NSString *agentsDir = [NSString stringWithFormat:@"%@/Library/LaunchAgents", NSHomeDirectory()];
	NSFileManager *fm = [NSFileManager defaultManager];
	if ([fm fileExistsAtPath:plistPath])
		return YES;
	[fm createDirectoryAtPath:agentsDir withIntermediateDirectories:YES attributes:nil error:NULL];
	NSDictionary *plist =
	    [NSDictionary dictionaryWithObjectsAndKeys:kLaunchAgentLabel, @"Label", args, @"ProgramArguments",
						   [NSNumber numberWithBool:YES], @"RunAtLoad", @"Aqua", @"LimitLoadToSessionType", nil];
	if (![plist writeToFile:plistPath atomically:YES])
		return NO;
	return launchctlLoadPlist(plistPath);
}

static BOOL registerLoginItemWithSMAppService(void)
{
	NSError *err = nil;
	SMAppService *service = [SMAppService mainAppService];
	SMAppServiceStatus st = service.status;
	if (st == SMAppServiceStatusEnabled || st == SMAppServiceStatusRequiresApproval)
		return YES;
	if ([service registerAndReturnError:&err])
		return YES;
	if (err)
		NSLog(@"SMAppService register failed: %@", err);
	return NO;
}

BOOL setStartAtLogin(BOOL startAtLogin)
{
	if (@available(macOS 13.0, *))
	{
		NSError *err = nil;
		SMAppService *service = [SMAppService mainAppService];
		if (startAtLogin)
		{
			removeLegacySessionLoginItems();
			removeLaunchAgentFiles();
			if (registerLoginItemWithSMAppService())
				return YES;
			// Unsigned or incomplete .app bundles cannot use SMAppService; fall back to LaunchAgent.
			return installLaunchAgent();
		}
		BOOL ok = [service unregisterAndReturnError:&err];
		if (err)
			NSLog(@"SMAppService unregister failed: %@", err);
		removeLaunchAgentFiles();
		return ok;
	}

	if (startAtLogin)
		return installLaunchAgent();
	removeLaunchAgentFiles();
	return YES;
}
*/
import "C"

func SetStartAtLogin(startAtLogin bool) bool {
	var startAtLoginVal C.BOOL
	if startAtLogin {
		startAtLoginVal = CTypeTrue
	} else {
		startAtLoginVal = CTypeFalse
	}
	return C.setStartAtLogin(C.BOOL(startAtLoginVal)) == CTypeTrue
}
