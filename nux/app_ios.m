// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

#include "_cgo_export.h"
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>

static uint64_t mainThreadId;

uint64_t threadID() {
	uint64_t id;
	if (pthread_threadid_np(pthread_self(), &id)) {
		abort();
	}
	return id;
}

int isMainThread(){
    return mainThreadId == threadID();
}

@interface NuxWindow : UIWindow
@end

@interface NuxViewController : UIViewController
@end

@implementation NuxViewController
- (void)viewDidLoad {
  [super viewDidLoad];
  NuxWindow * window = (NuxWindow*)([[[UIApplication sharedApplication] windows] objectAtIndex:0]);
  NSLog(@"nuxui NuxViewController viewDidLoad %d", (uintptr_t)window);
  go_windowCreated((uintptr_t)window);
  go_windowResized((uintptr_t)window);
	CGSize size = [UIScreen mainScreen].bounds.size;

}

- (void)viewWillTransitionToSize:(CGSize)size withTransitionCoordinator:(id<UIViewControllerTransitionCoordinator>)coordinator {
	[super viewWillTransitionToSize:size withTransitionCoordinator:coordinator];
  NSLog(@"nuxui NuxViewController viewWillTransitionToSize");
}

- (void)viewWillAppear:(BOOL)animated
{
	[super viewWillAppear:animated];
  NSLog(@"nuxui NuxViewController viewWillAppear");
}
@end



@implementation NuxWindow
- (void)drawRect:(CGRect)rect;
{
    NSLog(@"nuxui NuxWindow drawRect for ios, %d", (uintptr_t)self);
    CGContextRef context = UIGraphicsGetCurrentContext();
    go_windowDraw((uintptr_t)self);
}
@end

@interface MyDelegate : UIResponder< UIApplicationDelegate >
@end

@implementation MyDelegate
- ( BOOL ) application: ( UIApplication * ) application
           didFinishLaunchingWithOptions: ( NSDictionary * ) launchOptions {

  NuxWindow *window = [[NuxWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
  window.backgroundColor = [ UIColor whiteColor ];

  // UILabel *label = [ [ UILabel alloc ] init ];
  // label.text = @"Hello, World!";
  // label.center = CGPointMake( 100, 100 );
  // [ label sizeToFit ];
  // [ window addSubview: label ];

	NuxViewController *controller = [[NuxViewController alloc] initWithNibName:nil bundle:nil];
	window.rootViewController = controller;
	[window makeKeyAndVisible];
	return YES;
}
@end

void runApp(void)
{
	NSLog(@"gomobile run ios app basic nuxui");
	@autoreleasepool {
		UIApplicationMain(0, nil, nil, NSStringFromClass([MyDelegate class]));
	}
}

void window_getSize(uintptr_t window, int32_t *width, int32_t *height){
  CGSize size = [UIScreen mainScreen].bounds.size;
  if(width){*width = (int32_t)size.width;};
  if(height){*height = (int32_t)size.height;};
}

void window_getContentSize(uintptr_t window, int32_t *width, int32_t *height){
  CGSize size = [UIScreen mainScreen].bounds.size;
  if(width){*width = (int32_t)size.width;};
  if(height){*height = (int32_t)size.height;};
}

CGContextRef window_getCGContext(uintptr_t window){
    return UIGraphicsGetCurrentContext();
}