// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin
// +build !ios

#include "_cgo_export.h"
#import <Cocoa/Cocoa.h>

void terminate()
{
	[NSApp terminate:nil];
};

static const NSRange kEmptyRange = { NSNotFound, 0 };
@interface NuxText : NSView <NSTextInputClient> {}
@end

@implementation NuxText
- (BOOL)hasMarkedText
{
    NSLog(@"NuxText hasMarkedText");
    return NO;  // TODO:: what is 
}

- (NSRange)markedRange
{
    NSLog(@"NuxText markedRange");
    return kEmptyRange;
}

- (NSRange)selectedRange
{
    NSLog(@"NuxText selectedRange");
    return kEmptyRange;
}

- (void)setMarkedText:(id)string selectedRange:(NSRange)selectedRange replacementRange:(NSRange)replacementRange
{
    NSLog(@"NuxText setMarkedText");
    // do nothing
}

- (void)unmarkText
{
    NSLog(@"NuxText unmarkText");
    // do nothing
}

- (NSArray<NSAttributedStringKey> *)validAttributesForMarkedText
{
    return [NSArray array];
}

- (NSAttributedString *)attributedSubstringForProposedRange:(NSRange)range actualRange:(NSRangePointer)actualRange
{
    return nil;
}

- (void)insertText:(id)aString replacementRange:(NSRange)replacementRange
{
    const char *str;
    if ([aString isKindOfClass: [NSAttributedString class]]) {
        str = [[aString string] UTF8String];
    } else {
        str = [aString UTF8String];
    }
    NSLog(@"NuxText insertText %@", aString);
}

- (NSUInteger)characterIndexForPoint:(NSPoint)point
{
    NSLog(@"NuxText characterIndexForPoint");
    return 0;
}

- (NSRect)firstRectForCharacterRange:(NSRange)range actualRange:(NSRangePointer)actualRange
{
    NSLog(@"NuxText firstRectForCharacterRange");
    return NSZeroRect;
}

- (void)doCommandBySelector:(SEL)selector
{
    NSLog(@"NuxText doCommandBySelector");
    // do nothing otherwise beep
}

- (void) keyDown: (NSEvent *) theEvent
{
    NSLog(@"NuxText keyDown");
    [self interpretKeyEvents: [NSArray arrayWithObject: theEvent]];
}

@end // NuxText


@interface NuxWindow : NSWindow
@end

@implementation NuxWindow
- (BOOL)canBecomeKeyWindow
{
    return YES;
}

- (BOOL)canBecomeMainWindow
{
    return YES;
}

- (void)doCommandBySelector:(SEL)aSelector
{
    // do nothing otherwise beep
}

- (void)sendEvent:(NSEvent *)theEvent
{
	CGFloat x = [theEvent locationInWindow].x;
	CGFloat y = [[[theEvent window] contentView] bounds].size.height - [theEvent locationInWindow].y;
	CGFloat screenX = [NSEvent mouseLocation].x;
	CGFloat screenY = [[NSScreen mainScreen] frame].size.height-[NSEvent mouseLocation].y;
	NSEventType etype = [theEvent type];
	uintptr_t windptr = (uintptr_t)[theEvent window];

  	switch (etype) {
        case NSEventTypeMouseEntered:
        case NSEventTypeMouseExited: 
        // ignore
        break;
        case NSEventTypeLeftMouseDown:
        case NSEventTypeLeftMouseUp:
        case NSEventTypeRightMouseDown:
        case NSEventTypeRightMouseUp:
        case NSEventTypeMouseMoved:
        case NSEventTypeLeftMouseDragged:
        case NSEventTypeRightMouseDragged:
		case NSEventTypeOtherMouseDown:
        case NSEventTypeOtherMouseUp:
        case NSEventTypeOtherMouseDragged:
		NSLog(@"###  NuxWindow NSEventTypeAppKitDefined");
			go_mouseEvent(windptr, etype, x, y, screenX, screenY, 0, 0, 0, 0);
			break;
		case NSEventTypeScrollWheel:
			go_mouseEvent(windptr, etype, x, y, screenX, screenY, [theEvent scrollingDeltaX], [theEvent scrollingDeltaY], 0, 0);
			break;
		case NSEventTypePressure:
			go_mouseEvent(windptr, etype, x, y, screenX, screenY, 0, 0, [theEvent pressure], [theEvent stage]);
        	break;
        case NSEventTypeKeyDown:
        case NSEventTypeKeyUp:
			go_keyEvent(windptr, etype, [theEvent keyCode], [theEvent modifierFlags], [theEvent isARepeat],  (char *)[[theEvent characters] UTF8String]);
        	break;
        case NSEventTypeFlagsChanged:
		// NSLog(@"Shift=%d, control=%d opt=%d cmd=%d lock=%d pad=%d help=%d func=%d mask=%d", NSEventModifierFlagShift, NSEventModifierFlagControl, NSEventModifierFlagOption, NSEventModifierFlagCommand, NSEventModifierFlagCapsLock, NSEventModifierFlagNumericPad, NSEventModifierFlagHelp,NSEventModifierFlagFunction, NSEventModifierFlagDeviceIndependentFlagsMask);
			go_keyEvent(windptr, etype, [theEvent keyCode], [theEvent modifierFlags], 0, "");
        	break;
        break;
        case NSEventTypeAppKitDefined:
        NSLog(@"###  NSEventTypeAppKitDefined");
        break;
        case NSEventTypeSystemDefined:
        NSLog(@"###  NSEventTypeSystemDefined");
        break;
        case NSEventTypeApplicationDefined:
        NSLog(@"###  NSEventTypeApplicationDefined");
        break;
        case NSEventTypePeriodic:
        NSLog(@"###  NSEventTypePeriodic");
        break;
        case NSEventTypeCursorUpdate:
        NSLog(@"###  NSEventTypeCursorUpdate");
        break;
        case NSEventTypeTabletPoint:
        NSLog(@"###  NSEventTypeTabletPoint");
        break;
        case NSEventTypeTabletProximity:
        NSLog(@"###  NSEventTypeTabletProximity");
        break;

        case NSEventTypeGesture:
        NSLog(@"###  ###  NSEventTypeGesture");
        break;
        case NSEventTypeMagnify:
        NSLog(@"###  ###  NSEventTypeMagnify");
        break;
        case NSEventTypeSwipe:
        NSLog(@"###  NSEventTypeSwipe");
        break;
        case NSEventTypeRotate:
        NSLog(@"###  NSEventTypeRotate");
        break;
        case NSEventTypeBeginGesture:
        NSLog(@"###  NSEventTypeBeginGesture");
        break;
        case NSEventTypeEndGesture:
        NSLog(@"###  NSEventTypeEndGesture");
        break;
        case NSEventTypeSmartMagnify:
        NSLog(@"###  NSEventTypeSmartMagnify");
        break;

        case NSEventTypeDirectTouch:
        NSLog(@"###  NSEventTypeDirectTouch");
        break;
        case NSEventTypeQuickLook:
        NSLog(@"###  NSEventTypeQuickLook");
        break;
        // case NSEventTypeChangeMode:
        // NSLog(@"###  NSEventTypeChangeMode");
        // break;
        default:
        NSLog(@"unknow darwin event.");
        break;
    }

    [super sendEvent:theEvent];
}
@end // NuxWindow


@interface NuxWindowDelegate : NSObject <NSWindowDelegate>
@end

@implementation NuxWindowDelegate

// Sizing Windows
// - (NSSize)windowWillResize:(NSWindow *)sender  toSize:(NSSize)frameSize
// {

// }

- (void)windowDidResize:(NSNotification *)notification
{
	NSWindow* window = [notification object];
	NSLog(@"NuxWindowDelegate windowDidResize %@, keyWindow=", window);
	windowResized((uintptr_t)[notification object]);
}

- (void)windowWillStartLiveResize:(NSNotification *)notification
{
	NSLog(@"NuxWindowDelegate windowWillStartLiveResize");
}

- (void)windowDidEndLiveResize:(NSNotification *)notification
{
	NSLog(@"NuxWindowDelegate windowDidEndLiveResize");
}

// Updating Windows

- (void)windowDidUpdate:(NSNotification *)notification
{
	// when has input
	// NSLog(@"NuxWindowDelegate windowDidUpdate");
	windowDraw((uintptr_t)[notification object]);
}

// Exposing Windows

- (void)windowDidExpose:(NSNotification *)notification
{
	NSLog(@"NuxWindowDelegate windowDidExpose");
}

// Managing Key Status
- (void)windowDidBecomeKey:(NSNotification *)notification
{
	NSLog(@"NuxWindowDelegate windowDidBecomeKey");
}

- (void)windowDidResignKey:(NSNotification *)notification
{
	NSLog(@"NuxWindowDelegate windowDidResignKey");
}

// Managing Occlusion State

- (void)windowDidChangeOcclusionState:(NSNotification *)notification
{
	NSLog(@"NuxWindowDelegate windowDidChangeOcclusionState");
}

@end // NuxWindowDelegate

@interface NuxApplication : NSApplication
@end

@implementation NuxApplication

- (void)terminate:(id)sender
{
	NSLog(@"NuxApplication terminate");
	[super terminate:sender];
}

- (void)sendEvent:(NSEvent *)theEvent
{
    switch ([theEvent type]) {
		case NSEventTypeRightMouseUp:
			NSLog(@"window = %@", [[NSApplication sharedApplication] mainWindow]);
	NSLog(@"window2 = %@", [[[[NSApplication sharedApplication] windows] objectAtIndex:0] title]);
	NSLog(@"main window title = %@", [[[NSApplication sharedApplication] mainWindow] title]);
	NSLog(@"key window title = %@", [[[NSApplication sharedApplication] keyWindow] title]);
		break;
    }

    [super sendEvent:theEvent];
}

@end // NuxApplication 


@interface NuxApplicationDelegate : NSObject <NSApplicationDelegate>
@end

@implementation NuxApplicationDelegate
// Launching Applications
- (void)applicationWillFinishLaunching:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationWillFinishLaunching");
}

- (void)applicationDidFinishLaunching:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationDidFinishLaunching");
	// [NSApp activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
	// [[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
	// [window orderFrontRegardless];
	// [window makeKeyAndOrderFront:window];

    // windowCreated((uintptr_t)window);
    // windowResized((uintptr_t)window);
    // windowDraw((uintptr_t)window);
}

// Managing Active Status
- (void)applicationWillBecomeActive:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationWillBecomeActive");
}

- (void)applicationDidBecomeActive:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationDidBecomeActive");

	NSLog(@"window = %@", [[NSApplication sharedApplication] mainWindow]);
	NSLog(@"window2 = %@", [[[[NSApplication sharedApplication] windows] objectAtIndex:0] title]);
	NSLog(@"main window title = %@", [[[NSApplication sharedApplication] mainWindow] title]);
	NSLog(@"nsapp main window title = %@", [[NSApp mainWindow] title]);
	NSLog(@"key window title = %@", [[[NSApplication sharedApplication] keyWindow] title]);

	// [[NSApp mainWindow] orderFrontRegardless];
	// windowCreated((uintptr_t)[NSApp mainWindow]);
	NSLog(@"### windowCreated");
}

- (void)applicationWillResignActive:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationWillResignActive");
}

- (void)applicationDidResignActive:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationDidResignActive");
}

// // Terminating Applications
// // - (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender
// // {
// // 	NSLog(@"NuxApplicationDelegate applicationShouldTerminate");
// // }

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender
{
	NSLog(@"NuxApplicationDelegate applicationShouldTerminateAfterLastWindowClosed");
	return YES;
}

- (void)applicationWillTerminate:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationWillTerminate");
}


// Hiding Applications
- (void)applicationWillHide:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationWillHide");
}

- (void)applicationDidHide:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationDidHide");
}

- (void)applicationWillUnhide:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationWillUnhide");
}

- (void)applicationDidUnhide:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationDidUnhide");
}

// Managing Windows
- (void)applicationWillUpdate:(NSNotification *)notification
{
	// when has input
	// NSLog(@"NuxApplicationDelegate applicationWillUpdate");
}

- (void)applicationDidUpdate:(NSNotification *)notification
{
	// when has input
	// NSLog(@"NuxApplicationDelegate applicationDidUpdate");
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender hasVisibleWindows:(BOOL)flag
{
	NSLog(@"NuxApplicationDelegate applicationShouldHandleReopen");
	return NO;
}

@end //NuxApplicationDelegate



void runApp()
{
	@autoreleasepool{
		[NuxApplication sharedApplication];
		NSLog(@"go_nativeLoopPrepared");
		go_nativeLoopPrepared();
		
		[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular]; // show icon at docker
		[NSApp activateIgnoringOtherApps:YES];
		[NSApp setDelegate:[[NuxApplicationDelegate alloc] init]];

		NSUInteger windowStyle =  NSWindowStyleMaskTitled | NSWindowStyleMaskBorderless | NSWindowStyleMaskClosable |NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable | NSWindowStyleMaskUnifiedTitleAndToolbar | NSWindowStyleMaskDocModalWindow;
		NSRect windowRect = NSMakeRect(0, 0, 800, 600);
		NSLog(@"NSWindow alloc");
		NSWindow * window = [[NuxWindow alloc] initWithContentRect:windowRect
											styleMask:windowStyle
											backing:NSBackingStoreBuffered
											defer:NO];
		NSLog(@"NSWindow windowCreated %@", window);
		windowCreated((uintptr_t)window);
		window.title = @"nuxui title";
		[window setDelegate:[[NuxWindowDelegate alloc] init]];
		[window orderFrontRegardless];
		[window setAcceptsMouseMovedEvents:YES];

		[NSApp run];
	}
}
