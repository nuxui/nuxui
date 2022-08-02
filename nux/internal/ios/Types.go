// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package ios

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Foundation -framework CoreGraphics -framework UIKit -framework CoreText -framework GLKit

#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION -DGLES_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework GLKit -framework OpenGLES -framework QuartzCore

#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>
#import <CoreGraphics/CoreGraphics.h>

CGColorRef nux_CGColorMake(CGFloat red, CGFloat green, CGFloat blue, CGFloat alpha){
	return [[UIColor colorWithRed:red green:green blue:blue alpha:alpha] CGColor];
}

CGFloat nux_CGSize_width(CGSize size){
	return size.width;
}

CGFloat nux_CGSize_height(CGSize size){
	return size.height;
}

void nux_CGRect_origin(CGRect rect, CGFloat* outX, CGFloat* outY){
	if (outX) { *outX = rect.origin.x; };
	if (outY) { *outY = rect.origin.y; };
}

void nux_CGRect_size(CGRect rect, CGFloat* outW, CGFloat* outH){
	if (outW) { *outW = rect.size.width; };
	if (outH) { *outH = rect.size.height; };
}

int nux_NSArray_count(uintptr_t nsarray){
	return (int)((NSArray*)nsarray).count;
}

uintptr_t nux_NSArray_objectAtIndex(uintptr_t nsarray, NSUInteger index){
	return (uintptr_t)[(NSArray*)nsarray objectAtIndex:index];
}
*/
import "C"

const (
	_PI     = 3.1415926535897932384626433832795028841971
	_PI2    = _PI * 2
	_RADIAN = _PI / 180.0
)

type UIApplication C.uintptr_t
type UIWindow C.uintptr_t
type UIView C.uintptr_t
type UIImage C.uintptr_t
type UIFont C.uintptr_t
type UIEvent C.uintptr_t
type UITouch C.uintptr_t

type NSLayoutManager C.uintptr_t
type NSTextContainer C.uintptr_t
type NSTextStorage C.uintptr_t
type NSArray C.uintptr_t

type CGPoint C.CGPoint
type CGSize C.CGSize
type CGRect C.CGRect
type CGAffineTransform C.CGAffineTransform
type CGContextRef C.CGContextRef
type CGPathRef C.CGPathRef
type CGMutablePathRef C.CGMutablePathRef
type CGImageRef C.CGImageRef
type CGColorRef C.CGColorRef

type UIEventType int
type UIEventSubtype int
type UIKeyModifierFlags int
type UITouchType int
type UITouchPhase int

type WindowEvent struct {
	Window UIWindow
	Type   int
}

func (me UIFont) IsNil() bool   { return me == 0 }
func (me UIWindow) IsNil() bool { return me == 0 }

func CGRectMake(x, y, width, height float32) CGRect {
	return CGRect(C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height)))
}

func CGSizeMake(width, height float32) CGSize {
	return CGSize(C.CGSizeMake(C.CGFloat(width), C.CGFloat(height)))
}

func CGColorMake(red, green, blue, alpha float32) CGColorRef {
	return CGColorRef(C.nux_CGColorMake(C.CGFloat(red), C.CGFloat(green), C.CGFloat(blue), C.CGFloat(alpha)))
}

func (me CGRect) Origin() (x, y float32) {
	var x0, y0 C.CGFloat
	C.nux_CGRect_origin(C.CGRect(me), &x0, &y0)
	return float32(x0), float32(y0)
}

func (me CGRect) Size() (width, height float32) {
	var w, h C.CGFloat
	C.nux_CGRect_size(C.CGRect(me), &w, &h)
	return float32(w), float32(h)
}

func (me CGSize) Width() float32 {
	return float32(C.nux_CGSize_width(C.CGSize(me)))
}

func (me CGSize) Height() float32 {
	return float32(C.nux_CGSize_height(C.CGSize(me)))
}

func (me NSArray) Count() int {
	return int(C.nux_NSArray_count(C.uintptr_t(me)))
}

func (me NSArray) ObjectAtIndex(index int) uintptr {
	return uintptr(C.nux_NSArray_objectAtIndex(C.uintptr_t(me), C.NSUInteger(index)))
}
