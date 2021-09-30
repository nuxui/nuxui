// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin

#include "_cgo_export.h"
#import <Cocoa/Cocoa.h>
#include <cairo/cairo.h>
#include <cairo/cairo-quartz.h>

static const NSRange kEmptyRange = { NSNotFound, 0 };
bool mainWindowCreated = false;


@interface NuxText : NSView <NSTextInputClient> {
    NSString *_markedText;
    NSRange   _markedRange;
    NSRange   _selectedRange;
    NSRect _inputRect;
}
@end

@implementation NuxText
- (NSRect)inputRect
{
    return _inputRect;
}

- (void)setInputRect:(NSRect)rect
{
    _inputRect = rect;
}

- (BOOL)hasMarkedText
{
    return _markedText != nil;
}

- (NSRange)markedRange
{
    return _markedRange;
}

- (NSRange)selectedRange
{
    return _selectedRange;
}

- (void)setMarkedText:(id)aString selectedRange:(NSRange)selectedRange replacementRange:(NSRange)replacementRange
{
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

    go_typingEvent((uintptr_t)[self window], (char*)[aString UTF8String], 1, (int) selectedRange.location, (int) selectedRange.length);
}

- (void)unmarkText
{
    [_markedText release];
    _markedText = nil;

    go_typingEvent((uintptr_t)[self window], "", 1, 0, 0);
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
    if ([aString isKindOfClass:[NSAttributedString class]]) {
        aString = [aString string];
    }

    char* text = (char*)[aString UTF8String];
    NSLog(@"insertText %@", aString);
    
    // Don't post text events for unprintable characters
    if ((unsigned char)*text < ' ' || *text == 127) {
        return ;
    }

    go_typingEvent((uintptr_t)[self window], text, 0, 0, 0);
}

- (NSUInteger)characterIndexForPoint:(NSPoint)point
{
    NSLog(@"NuxText characterIndexForPoint");
    return 0;
}

- (NSRect)firstRectForCharacterRange:(NSRange)range actualRange:(NSRangePointer)actualRange
{
    // convert to screen rect
    return NSMakeRect(
        _inputRect.origin.x + [[self window] frame].origin.x, 
        _inputRect.origin.y + [[self window] frame].origin.y,
        _inputRect.size.width, 
        _inputRect.size.height);
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
@property (nonatomic, strong, readwrite) id nuxtext;
@property (nonatomic, strong, readwrite) id nuxview;
@property (nonatomic, readwrite) cairo_t* cairo;
@property (nonatomic, readwrite) CGContextRef cgContext;
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
        {
            CGFloat x = [theEvent locationInWindow].x;
	        CGFloat y = [[theEvent window] contentView].bounds.size.height - [theEvent locationInWindow].y;
			go_mouseEvent(windptr, etype, x, y, [theEvent buttonNumber], 0, 0);
			break;
        }
		case NSEventTypePressure:
        {
            CGFloat x = [theEvent locationInWindow].x;
	        CGFloat y = [[theEvent window] contentView].bounds.size.height - [theEvent locationInWindow].y;
			go_mouseEvent(windptr, etype, x, y, [theEvent buttonNumber], [theEvent pressure], [theEvent stage]);
			break;
        }
		case NSEventTypeScrollWheel:
        {
            CGFloat x = [theEvent locationInWindow].x;
	        CGFloat y = [[theEvent window] contentView].bounds.size.height - [theEvent locationInWindow].y;
            CGFloat scrollX = [theEvent scrollingDeltaX];
            CGFloat scrollY = [theEvent scrollingDeltaY];

            if ([theEvent hasPreciseScrollingDeltas])
            {
                scrollX *= 0.1;
                scrollY *= 0.1;
            }

            if (fabs(scrollX) > 0.0 || fabs(scrollY) > 0.0){
			    go_scrollEvent(windptr, x, y, scrollX, scrollY);
            }
			break;
        }
        case NSEventTypeKeyDown:
        case NSEventTypeKeyUp:
        // NuxText typing event conflict
            NSLog(@"============= NSEventTypeKeyDown");
			go_keyEvent(windptr, etype, [theEvent keyCode], [theEvent modifierFlags], [theEvent isARepeat],  (char *)[[theEvent characters] UTF8String]);
        	break; // TODO if true return else [super sendEvent:theEvent];
        case NSEventTypeFlagsChanged:
			go_keyEvent(windptr, etype, [theEvent keyCode], [theEvent modifierFlags], 0, "");
        	break;
        break;
        case NSEventTypeAppKitDefined:
        // NSLog(@"###  NSEventTypeAppKitDefined");
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

@interface NuxView : NSView{}
@end

int cg_changed = 0;

@implementation NuxView
- (void)drawRect:(NSRect)dirtyRect
{
    NuxWindow* w = (NuxWindow*)[self window];
    go_drawEvent((uintptr_t)w);
}
@end

@interface NuxWindowController : NSWindowController
@end

@implementation NuxWindowController
- (void)windowDidLoad
{
    //TODO:: not work ?
    NSLog(@"NuxWindowController windowDidLoad ");
}
@end


@interface NuxWindowDelegate : NSObject <NSWindowDelegate>
@end

@implementation NuxWindowDelegate

- (NSSize)windowWillResize:(NSWindow *)sender  toSize:(NSSize)frameSize
{
    NSLog(@"NuxWindowDelegate windowWillResize %@, sender=", sender);
    return frameSize;
}

- (void)windowDidResize:(NSNotification *)notification
{
	NSWindow* window = [notification object];
	NSLog(@"NuxWindowDelegate windowDidResize %@, keyWindow=", window);
    if (!mainWindowCreated){
        mainWindowCreated = true;
        windowCreated((uintptr_t)[notification object]);
    }
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
	// windowDraw((uintptr_t)[notification object]);
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
- (void)terminate:(id)sender;
- (void)sendEvent:(NSEvent *)theEvent;

+ (void)registerUserDefaults;
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
    case NSEventTypeApplicationDefined:
        if([theEvent timestamp] == 0 && 
        [theEvent data1] == 20180101 && 
        [theEvent data2] == 20201120 && 
        [theEvent subtype] == 0 && 
        [theEvent windowNumber] == 0){
            go_backToUI();
        }else{
            [super sendEvent:theEvent];
        }
        break;
    default:
        [super sendEvent:theEvent];
    }
}

+ (void)registerUserDefaults
{
    NSDictionary *appDefaults = [[NSDictionary alloc] initWithObjectsAndKeys:
                                 [NSNumber numberWithBool:NO], @"AppleMomentumScrollSupported",
                                 [NSNumber numberWithBool:NO], @"ApplePressAndHoldEnabled",
                                 [NSNumber numberWithBool:YES], @"ApplePersistenceIgnoreState",
                                 nil];
    [[NSUserDefaults standardUserDefaults] registerDefaults:appDefaults];
    [appDefaults release];
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
    if(true/*TODO:: background*/){
        [NSApp activateIgnoringOtherApps:YES];
    }

    [NuxApplication registerUserDefaults];
	// [NSApp activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
    // [NSApp activateIgnoringOtherApps:YES];
	// [[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
	// [window orderFrontRegardless];
	// [window makeKeyAndOrderFront:window];
}

// Managing Active Status
- (void)applicationWillBecomeActive:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationWillBecomeActive");
}

- (void)applicationDidBecomeActive:(NSNotification *)notification
{
	NSLog(@"NuxApplicationDelegate applicationDidBecomeActive");

	// NSLog(@"window = %@", [[NSApplication sharedApplication] mainWindow]);
	// NSLog(@"window2 = %@", [[[[NSApplication sharedApplication] windows] objectAtIndex:0] title]);
	// NSLog(@"main window title = %@", [[[NSApplication sharedApplication] mainWindow] title]);
	// NSLog(@"nsapp main window title = %@", [[NSApp mainWindow] title]);
	// NSLog(@"key window title = %@", [[[NSApplication sharedApplication] keyWindow] title]);

	// [[NSApp mainWindow] orderFrontRegardless];
	// windowCreated((uintptr_t)[NSApp mainWindow]);
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
		
		[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular]; // show icon at docker
		[NSApp setDelegate:[[NuxApplicationDelegate alloc] init]];

		NSUInteger windowStyle =  NSWindowStyleMaskTitled | NSWindowStyleMaskBorderless | NSWindowStyleMaskClosable |NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable | NSWindowStyleMaskUnifiedTitleAndToolbar | NSWindowStyleMaskDocModalWindow;
		NSRect windowRect = NSMakeRect(0, 0, 800, 600);
		NuxWindow * window = [[NuxWindow alloc] initWithContentRect:windowRect
											styleMask:windowStyle
											backing:NSBackingStoreBuffered
											defer:NO];

        // NuxWindowController * windowController = [[NuxWindowController alloc] initWithWindow:window];

        NuxView* nuxview = (NuxView*)[[[NuxView alloc] initWithFrame:NSMakeRect(0, 0, 0, 0)] autorelease]; 
        [window setContentView:nuxview];
        window.nuxview = nuxview;

		window.title = @"nuxui title";
		[window setDelegate:[[NuxWindowDelegate alloc] init]];
		[window setAcceptsMouseMovedEvents:YES];
        [window center];
        [window makeKeyAndOrderFront:nil];

	    [NSApp run];

	}
}
// ######################  begin  #####################

void loop_event(){
    NSEvent* event = [NSApp nextEventMatchingMask:NSEventMaskAny untilDate:nil inMode:NSDefaultRunLoopMode dequeue:YES];
    if (event != NULL){
        [NSApp sendEvent: event];
    }
}

void terminate()
{
    dispatch_async(dispatch_get_main_queue(), ^{
        @autoreleasepool{
	        [NSApp terminate:nil];
        }
    });
}

// TODO:: windptr
void startTextInput()
{
    dispatch_async(dispatch_get_main_queue(), ^{
        @autoreleasepool{
            NuxWindow* nuxWindow = (NuxWindow*)[NSApp keyWindow];
            if(nuxWindow.nuxtext != nil){
                [nuxWindow.nuxtext removeFromSuperview];
            }
            nuxWindow.nuxtext = [[[NuxText alloc] initWithFrame:NSMakeRect(0, 0, 0, 0)] autorelease];
            [[nuxWindow contentView] addSubview:nuxWindow.nuxtext];
            [nuxWindow makeFirstResponder:nuxWindow.nuxtext];
        }
    });
}

void stopTextInput()
{
    dispatch_async(dispatch_get_main_queue(), ^{
        @autoreleasepool{
            NuxWindow* nuxWindow = (NuxWindow*)[NSApp keyWindow];
            if(nuxWindow.nuxtext != nil){
                [nuxWindow.nuxtext removeFromSuperview];
            }
        }
    });
}

void setTextInputRect(float x, float y, float w, float h)
{
    dispatch_async(dispatch_get_main_queue(), ^{
        @autoreleasepool{
            NuxWindow* nuxWindow = (NuxWindow*)[NSApp keyWindow];
            if(nuxWindow.nuxtext != nil){
                float fixy = [[nuxWindow contentView] bounds].size.height - y;
                [nuxWindow.nuxtext setInputRect:NSMakeRect(x,fixy,w,h)];
            }
        }
    });
}

void invalidate(){
    dispatch_async(dispatch_get_main_queue(), ^{
        @autoreleasepool{
            NSLog(@"------ post invalidate -------");
            NuxWindow* nuxWindow = (NuxWindow*)[NSApp keyWindow];
            if(nuxWindow.nuxview != nil){
                [nuxWindow.nuxview setNeedsDisplay:YES];
            }
        }
    });
}

void backToUI(){
    dispatch_async(dispatch_get_main_queue(), ^{
        @autoreleasepool{
            NSEvent* event = [NSEvent otherEventWithType:NSEventTypeApplicationDefined
                                        location:NSMakePoint(0, 0)
                                   modifierFlags:0
                                       timestamp:0
                                    windowNumber:0
                                         context:nil
                                         subtype:0
                                           data1:20180101
                                           data2:20201120];
            [NSApp sendEvent:event];
        }
    });
}


// ######################  end  #####################




// ######################### nuxwindow ###########################

char* window_title(uintptr_t window){
    return (char *)[[(NSWindow*)window title] UTF8String];
}

void window_setTitle(uintptr_t window, char* title){
    ((NSWindow*)window).title = [NSString stringWithUTF8String:title];
}

float window_alpha(uintptr_t window){
    return (float)([(NSWindow*)window alphaValue]);
}

void window_setAlpha(uintptr_t window, float alpha){
    ((NSWindow*)window).alphaValue = alpha;
}

int32_t window_getWidth(uintptr_t window){
    return (int32_t)([(NSWindow*)window frame].size.width);
}

int32_t window_getHeight(uintptr_t window){
    return (int32_t)([(NSWindow*)window frame].size.height);
}

int32_t window_getContentWidth(uintptr_t window){
    return (int32_t)([[(NSWindow*)window contentView] bounds].size.width);
}

int32_t window_getContentHeight(uintptr_t window){
    return (int32_t)([[(NSWindow*)window contentView] bounds].size.height);
}

uintptr_t window_getCGContext(uintptr_t window){
    // NSWindow* w = (NSWindow*)window;
    // [w graphicsContext] Deprecated macOS 10.0~10.14
    // NSLog(@"1 gc=%@, ref=%@", [w graphicsContext], [w graphicsContext].CGContext);
    // NSLog(@"2 gc=%@, ref=%@", [NSGraphicsContext currentContext], [NSGraphicsContext currentContext].CGContext);
    return (uintptr_t)[NSGraphicsContext currentContext];
}
// ############################### nuxwindow end #################################


// ############################### cursor begin #################################
void cursor_getScreenPosition(CGFloat* outX, CGFloat* outY){
    *outX = [NSEvent mouseLocation].x;
	*outY = [NSScreen mainScreen].frame.size.height-[NSEvent mouseLocation].y;
}

void cursor_positionWindowToScreen(uintptr_t window, CGFloat x, CGFloat y, CGFloat *outX, CGFloat *outY){
    NSWindow* w = (NSWindow*)window;
    *outX = w.frame.origin.x + x;
    *outY = [NSScreen mainScreen].frame.size.height - ([w contentView].bounds.size.height - y + w.frame.origin.y);
}

void cursor_positionScreenToWindow(uintptr_t window, CGFloat x, CGFloat y, CGFloat *outX, CGFloat *outY){
    NSWindow* w = (NSWindow*)window;
    *outX = x - w.frame.origin.x;
    *outY = [w contentView].bounds.size.height - ( ([NSScreen mainScreen].frame.size.height - y) - w.frame.origin.y );
}
// ############################### cursor end   #################################
