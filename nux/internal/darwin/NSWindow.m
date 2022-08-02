// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "_cgo_export.h"
#import <Cocoa/Cocoa.h>

@interface NuxView : NSView <NSTextInputClient> {
  NSString *_markedText;
  NSRange _markedRange;
  NSRange _selectedRange;
  NSRect _inputRect;
}
@end // @interface NuxView

@implementation NuxView
- (void)drawRect:(NSRect)dirtyRect {
  CGContextRef cgContext = [NSGraphicsContext currentContext].CGContext;
  [NSGraphicsContext saveGraphicsState];
  // Flip text
  NSGraphicsContext *context =
      [NSGraphicsContext graphicsContextWithCGContext:cgContext flipped:true];
  [NSGraphicsContext setCurrentContext:context];

  go_nux_windowDrawRect((uintptr_t)[self window]);

  [NSGraphicsContext restoreGraphicsState];
}

- (NSRect)inputRect {
  return _inputRect;
}

- (void)setInputRect:(NSRect)rect {
  _inputRect = rect;
}

- (BOOL)hasMarkedText {
  return _markedText != nil;
}

- (NSRange)markedRange {
  return _markedRange;
}

- (NSRange)selectedRange {
  return _selectedRange;
}

- (void)setMarkedText:(id)aString
        selectedRange:(NSRange)selectedRange
     replacementRange:(NSRange)replacementRange {
  if ([aString isKindOfClass:[NSAttributedString class]]) {
    aString = [aString string];
  }

  if ([aString length] == 0) {
    [self unmarkText];
    return;
  }

  if (_markedText != aString) {
    [_markedText release];
    _markedText = [aString retain];
  }

  _selectedRange = selectedRange;
  _markedRange = NSMakeRange(0, [aString length]);

  go_nux_windowTypingEvent((uintptr_t)[self window], (char*)[aString UTF8String], 1,
  (int) selectedRange.location, (int) selectedRange.length);
}

- (void)unmarkText {
  [_markedText release];
  _markedText = nil;

  go_nux_windowTypingEvent((uintptr_t)[self window], "", 1, 0, 0);
}

- (NSArray<NSAttributedStringKey> *)validAttributesForMarkedText {
  return [NSArray array];
}

- (NSAttributedString *)attributedSubstringForProposedRange:(NSRange)range
                                                actualRange:(NSRangePointer)
                                                                actualRange {
  return nil;
}

- (void)insertText:(id)aString replacementRange:(NSRange)replacementRange {
  if ([aString isKindOfClass:[NSAttributedString class]]) {
    aString = [aString string];
  }

  char *text = (char *)[aString UTF8String];
  // NSLog(@"insertText %@", aString);

  // Don't post text events for unprintable characters
  if ((unsigned char)*text < ' ' || *text == 127) {
    return;
  }

  go_nux_windowTypingEvent((uintptr_t)[self window], text, 0, 0, 0);
}

- (NSUInteger)characterIndexForPoint:(NSPoint)point {
  NSLog(@"NuxView characterIndexForPoint");
  return 0;
}

- (NSRect)firstRectForCharacterRange:(NSRange)range
                         actualRange:(NSRangePointer)actualRange {
  // convert to screen rect
  return NSMakeRect(_inputRect.origin.x + [[self window] frame].origin.x,
                    _inputRect.origin.y + [[self window] frame].origin.y,
                    _inputRect.size.width, _inputRect.size.height);
}

- (void)doCommandBySelector:(SEL)selector {
  // NSLog(@"NuxView doCommandBySelector");
  // do nothing otherwise beep
}

- (void)keyDown:(NSEvent *)theEvent {
  // NSLog(@"NuxView keyDown");
  [self interpretKeyEvents:[NSArray arrayWithObject:theEvent]];
}
@end // @implementation NuxView

////////////////////// NuxWindow //////////////////////////
@interface NuxWindow : NSWindow
@end //@interface NuxWindow

@implementation NuxWindow
- (BOOL)canBecomeKeyWindow {
  return YES;
}

- (BOOL)canBecomeMainWindow {
  return YES;
}

- (void)doCommandBySelector:(SEL)aSelector {
  // do nothing otherwise beep
}

- (void)sendEvent:(NSEvent *)theEvent {
  if (go_nux_window_sendEvent((uintptr_t)theEvent) > 0) {
    // handled by user
  }else{
    [super sendEvent:theEvent];
  }
}
@end // @implementation NuxWindow

@interface NuxWindowController : NSWindowController
@end // @interface NuxWindowController

@implementation NuxWindowController
- (void)windowDidLoad {
  // TODO:: not work ?
  NSLog(@"NuxWindowController windowDidLoad ");
}
@end // @implementation NuxWindowController

@interface NuxWindowDelegate : NSObject <NSWindowDelegate>
@end // @interface NuxWindowDelegate

@implementation NuxWindowDelegate

- (NSSize)windowWillResize:(NSWindow *)sender toSize:(NSSize)frameSize {
  // NSLog(@"NuxWindowDelegate windowWillResize sender=%@", sender);
  return frameSize;
}

- (void)windowDidResize:(NSNotification *)notification {
  NSWindow *window = [notification object];
  // NSLog(@"NuxWindowDelegate windowDidResize keyWindow=%@", window);
  // if (!mainWindowCreated){
  //     mainWindowCreated = true;
  //     // go_windowCreated((uintptr_t)[notification object]);
  // }
  go_nux_windowDidResize((uintptr_t)[notification object]);
}

- (void)windowWillStartLiveResize:(NSNotification *)notification {
  // NSLog(@"NuxWindowDelegate windowWillStartLiveResize");
}

- (void)windowDidEndLiveResize:(NSNotification *)notification {
  // NSLog(@"NuxWindowDelegate windowDidEndLiveResize");
}

// Updating Windows

- (void)windowDidUpdate:(NSNotification *)notification {
  // when has input
  // NSLog(@"NuxWindowDelegate windowDidUpdate");
  // windowDraw((uintptr_t)[notification object]);
}

// Exposing Windows

- (void)windowDidExpose:(NSNotification *)notification {
  // NSLog(@"NuxWindowDelegate windowDidExpose");
}

// Managing Key Status
- (void)windowDidBecomeKey:(NSNotification *)notification {
  // NSLog(@"NuxWindowDelegate windowDidBecomeKey");
}

- (void)windowDidResignKey:(NSNotification *)notification {
  // NSLog(@"NuxWindowDelegate windowDidResignKey");
}

// Managing Occlusion State

- (void)windowDidChangeOcclusionState:(NSNotification *)notification {
  // NSLog(@"NuxWindowDelegate windowDidChangeOcclusionState");
}

@end // @implementation NuxWindowDelegate

uintptr_t nux_NewNSWindow(CGFloat width, CGFloat height) {
  NSRect rect = NSMakeRect(0, 0, width, height);
  id window = [[NuxWindow alloc]
      initWithContentRect:rect
                styleMask:NSWindowStyleMaskTitled |
                          NSWindowStyleMaskBorderless |
                          NSWindowStyleMaskClosable |
                          NSWindowStyleMaskMiniaturizable |
                          NSWindowStyleMaskResizable |
                          NSWindowStyleMaskUnifiedTitleAndToolbar |
                          NSWindowStyleMaskDocModalWindow
                  backing:NSBackingStoreBuffered
                    defer:NO];

  id root = [[NuxView alloc] initWithFrame:NSMakeRect(0, 0, 0, 0)];
  [window setContentView:root];
  [window setDelegate:[[NuxWindowDelegate alloc] init]];
  [root setInputRect:NSMakeRect(0, 0, 100, 100)];
  // [window makeFirstResponder:root];
  [window setAcceptsMouseMovedEvents:YES];
  return (uintptr_t)window;
}

char *nux_NSWindow_Title(uintptr_t window) {
  return (char *)[[(NSWindow *)window title] UTF8String];
}

void nux_NSWindow_SetTitle(uintptr_t window, char *title) {
  ((NSWindow *)window).title = [NSString stringWithUTF8String:title];
}

float nux_NSWindow_Alpha(uintptr_t window) {
  return (float)([(NSWindow *)window alphaValue]);
}

void nux_NSWindow_SetAlpha(uintptr_t window, float alpha) {
  ((NSWindow *)window).alphaValue = alpha;
}

void nux_NSWindow_Size(uintptr_t window, int32_t *width, int32_t *height) {
  CGSize size = [(NSWindow *)window frame].size;
  if (width) {
    *width = (int32_t)size.width;
  };
  if (height) {
    *height = (int32_t)size.height;
  };
}

void nux_NSWindow_ContentSize(uintptr_t window, int32_t *width,
                              int32_t *height) {
  CGSize size = [[(NSWindow *)window contentView] bounds].size;
  if (width) {
    *width = (int32_t)size.width;
  };
  if (height) {
    *height = (int32_t)size.height;
  };
}

void nux_NSWindow_Center(uintptr_t window) { [(NSWindow *)window center]; }

void nux_NSWindow_MakeKeyAndOrderFront(uintptr_t window) {
  [(NSWindow *)window makeKeyAndOrderFront:nil];
}

void nux_NSWindow_SetContentView(uintptr_t window, uintptr_t view) {
  [(NSWindow *)window setContentView:(NSView *)view];
}

void nux_NSWindow_StartTextInput_async(uintptr_t window) {
  dispatch_async(dispatch_get_main_queue(), ^{
    [(NSWindow *)window makeFirstResponder:[(NSWindow *)window contentView]];
  });
}

void nux_NSWindow_StopTextInput_async(uintptr_t window) {
  dispatch_async(dispatch_get_main_queue(), ^{
    [(NSWindow *)window makeFirstResponder:nil];
  });
}

void nux_NSWindow_SetTextInputRect_async(uintptr_t window, CGFloat x, CGFloat y,
                                   CGFloat width, CGFloat height) {
  dispatch_async(dispatch_get_main_queue(), ^{
    NSWindow *w = (NSWindow *)window;
    NuxView *v = (NuxView *)[w contentView];
    if (w != nil && v != nil) {
      // fix y coordinate to orignal 
      CGFloat fixy = [v bounds].size.height - y - height;
      [v setInputRect:NSMakeRect(x, fixy, width, height)];
    }
  });
}

void nux_NSWindow_InvalidateRect_async(uintptr_t window, CGFloat x, CGFloat y,
                                 CGFloat width, CGFloat height) {
  dispatch_async(dispatch_get_main_queue(), ^{
    NSWindow *w = (NSWindow *)window;
    if (w != nil && [w contentView] != nil) {
      if (width == 0 || height == 0) {
        [[w contentView] setNeedsDisplay:YES];
      } else {
        // fix y coordinate to orignal 
        CGFloat fixy = [w contentView].bounds.size.height - y - height;
        [[w contentView] setNeedsDisplayInRect:NSMakeRect(x, fixy, width, height)];
      }
    }
  });
}
