// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa

CGColorRef nux_CGColorMake(CGFloat red, CGFloat green, CGFloat blue, CGFloat alpha){
	return [[NSColor colorWithSRGBRed:red green:green blue:blue alpha:alpha] CGColor];
}

void nux_NSRect_size(NSRect rect, CGFloat *outWidth, CGFloat* outHeight){
	*outWidth = rect.size.width;
	*outHeight = rect.size.height;
}
*/
import "C"

const (
	_PI     = 3.1415926535897932384626433832795028841971
	_PI2    = _PI * 2
	_RADIAN = _PI / 180.0
)

type CGPoint C.CGPoint
type CGSize C.CGSize
type CGRect C.CGRect
type CGPath C.CGPathRef
type CGMutablePath C.CGMutablePathRef
type CGContext C.CGContextRef
type CGAffineTransform C.CGAffineTransform
type CGImage C.CGImageRef
type CGColor C.CGColorRef

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

type NSWindowStyleMask uint32
type NSEventType uint32
type NSEventSubtype int32
type NSEventModifierFlags uint32

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

func CGRectMake(x, y, width, height float32) CGRect {
	return CGRect(C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height)))
}

func CGSizeMake(width, height float32) CGSize {
	return CGSize(C.CGSizeMake(C.CGFloat(width), C.CGFloat(height)))
}

func CGColorMake(red, green, blue, alpha float32) CGColor {
	return CGColor(C.nux_CGColorMake(C.CGFloat(red), C.CGFloat(green), C.CGFloat(blue), C.CGFloat(alpha)))
}

func NSMakeRect(x, y, width, height float32) NSRect {
	return NSRect(C.NSMakeRect(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height)))
}

func (me NSRect) Size() (width, height float32) {
	var w, h C.CGFloat
	C.nux_NSRect_size(C.NSRect(me), &w, &h)
	return float32(w), float32(h)
}
