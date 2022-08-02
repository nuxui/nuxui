// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

#include "_cgo_export.h"
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>
#include <pthread.h>

@interface NuxViewController : UIViewController
@end

@implementation NuxViewController
- (void)viewDidLoad {
    [super viewDidLoad];
    NSLog(@"nuxui NuxViewController viewDidLoad");
    go_nux_windowDidLoad(0);
}

- (void)viewWillTransitionToSize:(CGSize)size withTransitionCoordinator:(id<UIViewControllerTransitionCoordinator>)coordinator {
    [super viewWillTransitionToSize:size withTransitionCoordinator:coordinator];
    NSLog(@"nuxui NuxViewController viewWillTransitionToSize");
}

- (void)viewWillAppear:(BOOL)animated
{
    [super viewWillAppear:animated];
}
@end

@interface NuxWindow : UIWindow
@end

@implementation NuxWindow
- (void)drawRect:(CGRect)rect {
    go_nux_windowDrawRect((uintptr_t)self);
}

- (void)sendEvent:(UIEvent *)event {
    go_nux_window_sendEvent((uintptr_t)event);
}
@end

uintptr_t nux_NewUIWindow(CGFloat width, CGFloat height) {
//   NuxWindow *window = [[NuxWindow alloc] initWithFrame:CGRectMake(0, 0, width, height)];
    NuxWindow *window = [[NuxWindow alloc] initWithFrame:[[UIScreen mainScreen] bounds]];
    window.backgroundColor = [ UIColor whiteColor ];

    NuxViewController *controller = [[NuxViewController alloc] initWithNibName:nil bundle:nil];
    window.rootViewController = controller;
    return (uintptr_t)window;
}

void nux_UIWindow_makeKeyAndVisible(uintptr_t nuxwindow){
    [(NuxWindow*)nuxwindow makeKeyAndVisible];
}

CGRect nux_UIWindow_frame(uintptr_t nuxwindow){
    return [(NuxWindow*)nuxwindow frame];
}

void nux_UIWindow_InvalidateRect_async(uintptr_t nuxwindow, CGFloat x, CGFloat y, CGFloat width, CGFloat height) {
  dispatch_async(dispatch_get_main_queue(), ^{
    UIWindow *w = (UIWindow *)nuxwindow;
    if (w != nil){
        [w setNeedsDisplay];
        // TODO:: setNeedsDisplayInRect
    }
  });
}