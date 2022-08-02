// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "_cgo_export.h"
#import <Cocoa/Cocoa.h>

@interface NuxApplication : NSApplication
- (void)terminate:(id)sender;
- (void)sendEvent:(NSEvent *)theEvent;

+ (void)registerUserDefaults;
@end // @interface NuxApplication

@implementation NuxApplication
- (void)terminate:(id)sender {
  // NSLog(@"NuxApplication terminate");
  [super terminate:sender];
}

- (void)sendEvent:(NSEvent *)theEvent {
  if (go_nux_app_sendEvent((uintptr_t)theEvent) > 0) {
    // handled by user
  }else{
    [super sendEvent:theEvent];
  }
}

+ (void)registerUserDefaults {
  NSDictionary *appDefaults = [[NSDictionary alloc]
      initWithObjectsAndKeys:[NSNumber numberWithBool:NO],
                             @"AppleMomentumScrollSupported",
                             [NSNumber numberWithBool:NO],
                             @"ApplePressAndHoldEnabled",
                             [NSNumber numberWithBool:YES],
                             @"ApplePersistenceIgnoreState", nil];
  [[NSUserDefaults standardUserDefaults] registerDefaults:appDefaults];
  [appDefaults release];
}

@end // @implementation NuxApplication

@interface NuxApplicationDelegate : NSObject <NSApplicationDelegate>
@end // @interface NuxApplicationDelegate

@implementation NuxApplicationDelegate
// Launching Applications
- (void)applicationWillFinishLaunching:(NSNotification *)notification {
  go_nux_app_delegate(1, (uintptr_t)notification);
}

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
  if (true /*TODO:: background*/) {
    [NSApp activateIgnoringOtherApps:YES];
  }

  [NuxApplication registerUserDefaults];
  // [NSApp activateWithOptions:(NSApplicationActivateAllWindows |
  // NSApplicationActivateIgnoringOtherApps)];
  // [NSApp activateIgnoringOtherApps:YES];
  // [[NSRunningApplication currentApplication]
  // activateWithOptions:(NSApplicationActivateAllWindows |
  // NSApplicationActivateIgnoringOtherApps)]; [window orderFrontRegardless];
  // [window makeKeyAndOrderFront:window];
  // go_nativeLoopPrepared();
  go_nux_app_delegate(2, (uintptr_t)notification);
}

// Managing Active Status
- (void)applicationWillBecomeActive:(NSNotification *)notification {
  go_nux_app_delegate(3, (uintptr_t)notification);
}

- (void)applicationDidBecomeActive:(NSNotification *)notification {
  go_nux_app_delegate(4, (uintptr_t)notification);

  // NSLog(@"window = %@", [[NSApplication sharedApplication] mainWindow]);
  // NSLog(@"window2 = %@", [[[[NSApplication sharedApplication] windows]
  // objectAtIndex:0] title]); NSLog(@"main window title = %@", [[[NSApplication
  // sharedApplication] mainWindow] title]); NSLog(@"nsapp main window title =
  // %@", [[NSApp mainWindow] title]); NSLog(@"key window title = %@",
  // [[[NSApplication sharedApplication] keyWindow] title]);

  // [[NSApp mainWindow] orderFrontRegardless];
  // windowCreated((uintptr_t)[NSApp mainWindow]);
}

- (void)applicationWillResignActive:(NSNotification *)notification {
  go_nux_app_delegate(5, (uintptr_t)notification);
}

- (void)applicationDidResignActive:(NSNotification *)notification {
  go_nux_app_delegate(6, (uintptr_t)notification);
}

// Terminating Applications
- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication*)sender
{
  return go_nux_app_delegate(7, (uintptr_t)sender);
}

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
  return go_nux_app_delegate(8, (uintptr_t)sender) > 0;
}

- (void)applicationWillTerminate:(NSNotification *)notification {
  go_nux_app_delegate(9, (uintptr_t)notification);
}

// Hiding Applications
- (void)applicationWillHide:(NSNotification *)notification {
  go_nux_app_delegate(10, (uintptr_t)notification);
}

- (void)applicationDidHide:(NSNotification *)notification {
  go_nux_app_delegate(11, (uintptr_t)notification);
}

- (void)applicationWillUnhide:(NSNotification *)notification {
  go_nux_app_delegate(12, (uintptr_t)notification);
}

- (void)applicationDidUnhide:(NSNotification *)notification {
  go_nux_app_delegate(13, (uintptr_t)notification);
}

// Managing Windows
- (void)applicationWillUpdate:(NSNotification *)notification {
  // when has input
}

- (void)applicationDidUpdate:(NSNotification *)notification {
  // when has input
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  return NO;
}

@end // @implementation NuxApplicationDelegate

uintptr_t nux_NSApp_SharedApplication() {
  NSApplication *app = [NuxApplication sharedApplication];
  // show icon at docker
  [app setActivationPolicy:NSApplicationActivationPolicyRegular];
  [app setDelegate:[[NuxApplicationDelegate alloc] init]];
  return (uintptr_t)app;
}

uintptr_t nux_NSApp() { return (uintptr_t)NSApp; }

void nux_NSApp_Run(uintptr_t app) { [(NSApplication *)app run]; }

uintptr_t nux_NSApp_KeyWindow(uintptr_t app) {
  return (uintptr_t)[(NSApplication *)app keyWindow];
}

void nux_NSApp_Terminate(uintptr_t app) {
  dispatch_async(dispatch_get_main_queue(), ^{
    [(NSApplication *)app terminate:nil];
  });
}

void nux_BackToUI() {
  dispatch_async(dispatch_get_main_queue(), ^{
    go_nux_backToUI();
  });
}



