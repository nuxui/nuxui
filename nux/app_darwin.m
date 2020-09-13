// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin
// +build !ios

#include "_cgo_export.h"
// #include <pthread.h>
// #include <stdio.h>

// #include <Carbon/Carbon.h>
#import <Cocoa/Cocoa.h>
// #import <Foundation/Foundation.h>
// #import <OpenGL/gl3.h>

void stopApp(void);

// void makeCurrentContext(GLintptr context) {
// 	NSOpenGLContext* ctx = (NSOpenGLContext*)context;
// 	[ctx makeCurrentContext];
// }

// uint64 threadID() {
// 	uint64 id;
// 	if (pthread_threadid_np(pthread_self(), &id)) {
// 		abort();
// 	}
// 	return id;
// }

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
  switch ([theEvent type]) {
        case NSEventTypeLeftMouseDown:
        case NSEventTypeOtherMouseDown:
        case NSEventTypeRightMouseDown:
        case NSEventTypeLeftMouseUp:
        case NSEventTypeOtherMouseUp:
        case NSEventTypeRightMouseUp:
        case NSEventTypeLeftMouseDragged:
        case NSEventTypeRightMouseDragged:
        case NSEventTypeOtherMouseDragged: /* usually middle mouse dragged */
        case NSEventTypeMouseMoved:
        case NSEventTypeScrollWheel:
            NSLog(@"NuxWindow Mouse Event");
            break;
        case NSEventTypeKeyDown:
            NSLog(@"NuxWindow Key Down Event");
            break;
        case NSEventTypeKeyUp:
            NSLog(@"NuxWindow Key Up Event");
            break;
        case NSEventTypeFlagsChanged:
            NSLog(@"NuxWindow Key Flags Event");
            break;
        default:
            NSLog(@"NuxWindow Other Event");
            break;
    }

    [super sendEvent:theEvent];
}

// TODO:: draw event
// TODO:: 
@end // NuxWindow

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
        case NSEventTypeLeftMouseDown:
        case NSEventTypeOtherMouseDown:
        case NSEventTypeRightMouseDown:
        case NSEventTypeLeftMouseUp:
        case NSEventTypeOtherMouseUp:
        case NSEventTypeRightMouseUp:
        case NSEventTypeLeftMouseDragged:
        case NSEventTypeRightMouseDragged:
        case NSEventTypeOtherMouseDragged: /* usually middle mouse dragged */
        case NSEventTypeMouseMoved:
        case NSEventTypeScrollWheel:
            NSLog(@"NuxApplication Mouse Event");
            break;
        case NSEventTypeKeyDown:
            NSLog(@"NuxApplication Key Down Event");
            break;
        case NSEventTypeKeyUp:
            NSLog(@"NuxApplication Key Up Event");
            break;
        case NSEventTypeFlagsChanged:
            NSLog(@"NuxApplication Key Flags Event");
            break;
        default:
            NSLog(@"NuxApplication Other Event");
            break;
    }

    [super sendEvent:theEvent];
}

@end // NuxApplication 


// ################################## begin NuxView ############################################
@interface NuxView : NSOpenGLView<NSApplicationDelegate, NSWindowDelegate, NSTextInputClient>
{
	
}
@end

@implementation NuxView
- (void)prepareOpenGL {
	[self setWantsBestResolutionOpenGLSurface:YES];
	GLint swapInt = 1;

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
	[[self openGLContext] setValues:&swapInt forParameter:NSOpenGLCPSwapInterval];
#pragma clang diagnostic pop

	// Using attribute arrays in OpenGL 3.3 requires the use of a VBA.
	// But VBAs don't exist in ES 2. So we bind a default one.
	GLuint vba;
	glGenVertexArrays(1, &vba);
	glBindVertexArray(vba);

	// startloop((GLintptr)[self openGLContext]);
}

- (void)reshape {
	[super reshape];

	// Calculate screen PPI.
	//
	// Note that the backingScaleFactor converts from logical
	// pixels to actual pixels, but both of these units vary
	// independently from real world size. E.g.
	//
	// 13" Retina Macbook Pro, 2560x1600, 227ppi, backingScaleFactor=2, scale=3.15
	// 15" Retina Macbook Pro, 2880x1800, 220ppi, backingScaleFactor=2, scale=3.06
	// 27" iMac,               2560x1440, 109ppi, backingScaleFactor=1, scale=1.51
	// 27" Retina iMac,        5120x2880, 218ppi, backingScaleFactor=2, scale=3.03
	NSScreen *screen = [NSScreen mainScreen];
	double screenPixW = [screen frame].size.width * [screen backingScaleFactor];

	CGDirectDisplayID display = (CGDirectDisplayID)[[[screen deviceDescription] valueForKey:@"NSScreenNumber"] intValue];
	CGSize screenSizeMM = CGDisplayScreenSize(display); // in millimeters
	float ppi = 25.4 * screenPixW / screenSizeMM.width;
	float pixelsPerPt = ppi/72.0;

	// The width and height reported to the geom package are the
	// bounds of the OpenGL view. Several steps are necessary.
	// First, [self bounds] gives us the number of logical pixels
	// in the view. Multiplying this by the backingScaleFactor
	// gives us the number of actual pixels.
	NSRect r = [self bounds];
	int w = r.size.width * [screen backingScaleFactor];
	int h = r.size.height * [screen backingScaleFactor];

	// setGeom(pixelsPerPt, w, h);
    // NSLog(@"reshape width = %f , height = %f", r.size.width, r.size.height);
    // NSLog(@"reshape w = %f , h = %f", w, h);
}

- (void)drawRect:(NSRect)theRect {
	// Called during resize. This gets rid of flicker when resizing.
	// drawgl();
    NSLog(@"drawRect");
    windowDraw((uintptr_t)self.window);
}

- (void)sendEvent:(NSEvent *)theEvent
{
	NSLog(@"sendEvent 0000000000");
    [super sendEvent:theEvent];
}

- (void)insertText:(id)string {
	NSLog(@"insertText 0000000000");
    [super insertText:string];  // have superclass insert it
}

//https://developer.apple.com/documentation/appkit/nsevent
- (void)mouseDown:(NSEvent *)theEvent {
	double scale = [[NSScreen mainScreen] backingScaleFactor];
	NSPoint p = [theEvent locationInWindow];
    windowDraw((uintptr_t)[theEvent window]);
}

- (void)mouseUp:(NSEvent *)theEvent {
	double scale = [[NSScreen mainScreen] backingScaleFactor];
	NSPoint p = [theEvent locationInWindow];
	// eventMouseEnd(p.x * scale, p.y * scale);
	// id contentHeight = [self contentRectForFrameRect: self.frame];
	// CGFloat titleBarHeight = self.frame.size.height - contentHeight;
		NSRect r = [self bounds];


	NSRect r2 = [[[theEvent window] contentView] bounds];
    NSLog(@"mouseUp %f, %f, %f, %f, %f", p.x, p.y, r.size.height, r2.size.height, [[theEvent window] frame].size.height);
    // windowDraw((uintptr_t)[theEvent window]);

	
	onMouseDown((uintptr_t)[theEvent window], (float)p.x, r2.size.height- (float)p.y);
}

- (void)mouseDragged:(NSEvent *)theEvent {
	double scale = [[NSScreen mainScreen] backingScaleFactor];
	NSPoint p = [theEvent locationInWindow];
	// eventMouseDragged(p.x * scale, p.y * scale);
    NSLog(@"mouseDragged");
    windowDraw((uintptr_t)[theEvent window]);
}

- (void)windowDidBecomeKey:(NSNotification *)notification {
	// lifecycleFocused();
}

- (void)windowDidResignKey:(NSNotification *)notification {
	if (![NSApp isHidden]) {
		// lifecycleVisible();
	}
}



- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    NSLog(@"applicationDidFinishLaunching");
	[[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
	[self.window makeKeyAndOrderFront:self];

    windowCreated((uintptr_t)self.window);
    windowResized((uintptr_t)self.window);
    windowDraw((uintptr_t)self.window);
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
	// lifecycleDead();
    NSLog(@"applicationWillTerminate");
}

- (void)applicationDidHide:(NSNotification *)aNotification {
	// lifecycleAlive();
    NSLog(@"applicationDidHide");
}

- (void)applicationWillUnhide:(NSNotification *)notification {
	// lifecycleVisible();
    NSLog(@"applicationWillUnhide");
}

- (void)windowDidResize:(NSNotification *)notification {
	NSLog(@"windowDidResize");
    windowResized((uintptr_t)self.window);

	[[NSOpenGLContext currentContext] update];
}

- (void)windowDidExpose:(NSNotification *)notification {
	NSLog(@"windowDidExpose");
}

- (void)windowWillClose:(NSNotification *)notification {
	// lifecycleAlive();
    NSLog(@"windowWillClose");
    stopApp();
}

- (void)windowDidUpdate:(NSNotification *)notification {
	// lifecycleAlive();
    // NSLog(@"windowDidUpdate");
}

@end
// ################################## end NuxView ############################################

void window_setTitle(uintptr_t window, char* title){
    NSLog(@"window_setTitle");
    NSWindow* w = (NSWindow*)window;
    w.title = [NSString stringWithUTF8String:title];
}

void
runApp(void) {
	[NSAutoreleasePool new];
	[NuxApplication sharedApplication];
	[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

	id menuBar = [[NSMenu new] autorelease];
	id menuItem = [[NSMenuItem new] autorelease];
	id menuItem2 = [[NSMenuItem new] autorelease];
	[menuBar addItem:menuItem];
	[menuBar addItem:menuItem2];
	[NSApp setMainMenu:menuBar];

	id menu = [[NSMenu new] autorelease];
	id menu2 = [[NSMenu new] autorelease];
	id name = [[NSProcessInfo processInfo] processName];

	id hideMenuItem = [[[NSMenuItem alloc] initWithTitle:@"Hide"
		action:@selector(hide:) keyEquivalent:@"h"]
		autorelease];
	[menu addItem:hideMenuItem];

	id hideMenuItem2 = [[[NSMenuItem alloc] initWithTitle:@"Hide2"
		action:@selector(hide:) keyEquivalent:@"h"]
		autorelease];
	[menu2 addItem:hideMenuItem2];

	id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:@"Quit"
		action:@selector(terminate:) keyEquivalent:@"q"]
		autorelease];
	[menu addItem:quitMenuItem];
	[menuItem setSubmenu:menu];
	[menuItem2 setSubmenu:menu2];

	NSRect rect = NSMakeRect(0, 0, 800, 600);

	NuxWindow* window = [[[NuxWindow alloc] initWithContentRect:rect
			styleMask:NSWindowStyleMaskTitled
			backing:NSBackingStoreBuffered
			defer:NO]
		autorelease];
	window.styleMask |= NSWindowStyleMaskResizable;
	window.styleMask |= NSWindowStyleMaskMiniaturizable;
	window.styleMask |= NSWindowStyleMaskClosable;
	window.title = name;
	[window cascadeTopLeftFromPoint:NSMakePoint(20,20)];

	NSOpenGLPixelFormatAttribute attr[] = {
		NSOpenGLPFAOpenGLProfile, NSOpenGLProfileVersion3_2Core,
		NSOpenGLPFAColorSize,     24,
		NSOpenGLPFAAlphaSize,     8,
		NSOpenGLPFADepthSize,     16,
		NSOpenGLPFAAccelerated,
		NSOpenGLPFADoubleBuffer,
		NSOpenGLPFAAllowOfflineRenderers,
		0
	};
	id pixFormat = [[NSOpenGLPixelFormat alloc] initWithAttributes:attr];
	NuxView* view = [[NuxView alloc] initWithFrame:rect pixelFormat:pixFormat];
	id text = [[NuxText alloc] initWithFrame: NSMakeRect(0, 0, 0, 0)];
	[window setContentView:view];
	[window setDelegate:view];
	[NSApp setDelegate:view];

	[[window contentView] addSubview:text];
    [window makeFirstResponder: text];
    nativeLoopPrepared();
	[NSApp run];
}

void stopApp(void) {
	[NSApp terminate:nil];
}
