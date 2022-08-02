// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa -framework UniformTypeIdentifiers

CGColorRef nux_CGColorMake(CGFloat red, CGFloat green, CGFloat blue, CGFloat alpha){
	return [[NSColor colorWithSRGBRed:red green:green blue:blue alpha:alpha] CGColor];
}

void nux_NSRect_size(NSRect rect, CGFloat *outWidth, CGFloat* outHeight){
	*outWidth = rect.size.width;
	*outHeight = rect.size.height;
}

int nux_NSArray_count(uintptr_t nsarray){
	return (int)((NSArray*)nsarray).count;
}

uintptr_t nux_NSArray_objectAtIndex(uintptr_t nsarray, NSUInteger index){
	return (uintptr_t)[(NSArray*)nsarray objectAtIndex:index];
}

uintptr_t nux_UTType_typeWithFilenameExtension(char* ext){
	return (uintptr_t)[UTType typeWithFilenameExtension: [NSString stringWithUTF8String:ext]];
}
*/
import "C"
import "unsafe"

const (
	_PI     = 3.1415926535897932384626433832795028841971
	_PI2    = _PI * 2
	_RADIAN = _PI / 180.0
)

type CGPoint C.CGPoint
type CGSize C.CGSize
type CGRect C.CGRect
type CGAffineTransform C.CGAffineTransform

// typedef struct CF_BRIDGED_TYPE(id) CGPath *CGMutablePathRef;
// typedef const struct CF_BRIDGED_TYPE(id) CGPath *CGPathRef;
// /Library/Developer/CommandLineTools/SDKs/MacOSX11.3.sdk/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/CGPath.h
type CGPathRef C.CGPathRef
type CGMutablePathRef C.CGMutablePathRef
type CGContextRef C.CGContextRef
type CGImageRef C.CGImageRef
type CGColorRef C.CGColorRef

type NSPoint C.NSPoint
type NSRect C.NSRect
type NSView C.uintptr_t
type NSApplication C.uintptr_t
type NSEvent C.uintptr_t
type NSWindow C.uintptr_t
type NSFont C.uintptr_t
type NSFontManager C.uintptr_t
type NSLayoutManager C.uintptr_t
type NSTextContainer C.uintptr_t
type NSTextStorage C.uintptr_t
type NSCursor C.uintptr_t
type NSSavePanel C.uintptr_t
type NSOpenPanel C.uintptr_t
type NSWorkspace C.uintptr_t
type NSNotification C.uintptr_t
type NSArray C.uintptr_t
type NSURL C.uintptr_t
type UTType C.uintptr_t

type NSWindowStyleMask uint32
type NSEventType uint32
type NSEventSubtype int32
type NSEventModifierFlags uint32
type NSModalResponse int32

type WindowEvent struct {
	Window NSWindow
	Type   int
}

type TypingEvent struct {
	Window   NSWindow
	Text     string
	Action   int32 // 0 = Action_Input, 1 = Action_Preedit
	Location int32
	Length   int32
}

func (me NSFont) IsNil() bool        { return me == 0 }
func (me NSView) IsNil() bool        { return me == 0 }
func (me NSWindow) IsNil() bool      { return me == 0 }
func (me NSApplication) IsNil() bool { return me == 0 }
func (me UTType) IsNil() bool        { return me == 0 }
func (me NSURL) IsNil() bool         { return me == 0 }

func CGRectMake(x, y, width, height float32) CGRect {
	return CGRect(C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height)))
}

func CGSizeMake(width, height float32) CGSize {
	return CGSize(C.CGSizeMake(C.CGFloat(width), C.CGFloat(height)))
}

func CGColorMake(red, green, blue, alpha float32) CGColorRef {
	return CGColorRef(C.nux_CGColorMake(C.CGFloat(red), C.CGFloat(green), C.CGFloat(blue), C.CGFloat(alpha)))
}

func NSMakeRect(x, y, width, height float32) NSRect {
	return NSRect(C.NSMakeRect(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height)))
}

func (me NSRect) Size() (width, height float32) {
	var w, h C.CGFloat
	C.nux_NSRect_size(C.NSRect(me), &w, &h)
	return float32(w), float32(h)
}

func (me NSArray) Count() int {
	return int(C.nux_NSArray_count(C.uintptr_t(me)))
}

func (me NSArray) ObjectAtIndex(index int) uintptr {
	return uintptr(C.nux_NSArray_objectAtIndex(C.uintptr_t(me), C.NSUInteger(index)))
}

func UTTypeWithFilenameExtension(ext string) UTType {
	cstr := C.CString(ext)
	defer C.free(unsafe.Pointer(cstr))
	return UTType(C.nux_UTType_typeWithFilenameExtension(cstr))
}
