// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

#include "_cgo_export.h"
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>

@interface NuxApplicationDelegate : UIResponder<UIApplicationDelegate>
@end

@implementation NuxApplicationDelegate
- (BOOL)application:(UIApplication *)application willFinishLaunchingWithOptions:(NSDictionary<UIApplicationLaunchOptionsKey, id> *)launchOptions {
  go_nux_app_delegate(1, (uintptr_t)launchOptions);
  return YES;
}

- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary<UIApplicationLaunchOptionsKey, id> *)launchOptions {
  go_nux_app_delegate(2, (uintptr_t)launchOptions);
	return YES;
}
@end

void nux_UIApplication_Run()
{	
  char* argv[] = {};
  [UIApplication sharedApplication];
	@autoreleasepool {
		UIApplicationMain(0, argv, nil, NSStringFromClass([NuxApplicationDelegate class]));
	}
}

uintptr_t nux_UIApplication_sharedApplication() {
  UIApplication *app = [UIApplication sharedApplication];
  return (uintptr_t)app;
}

void nux_BackToUI() {
  dispatch_async(dispatch_get_main_queue(), ^{
    go_nux_backToUI();
  });
}