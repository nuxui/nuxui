// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package ios

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Foundation -framework CoreGraphics -framework UIKit -framework CoreText -framework GLKit

#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>
#import <CoreGraphics/CoreGraphics.h>

void nux_CGContextDrawImage(CGContextRef ctx, CGRect rect, CGImageRef image){
	CGContextSaveGState(ctx);
	CGContextTranslateCTM(ctx, 0, rect.size.height);
	CGContextScaleCTM(ctx, 1, -1);
	CGContextDrawImage(ctx, rect, image);
	CGContextRestoreGState(ctx);
}

CGContextRef nux_UIGraphicsGetCurrentContext(){
    return UIGraphicsGetCurrentContext();
}
*/
import "C"
import "unsafe"

func UIGraphicsGetCurrentContext() CGContextRef {
	return CGContextRef(C.nux_UIGraphicsGetCurrentContext())
}

func CGContextResetClip(ctx CGContextRef) {
	C.CGContextResetClip(C.CGContextRef(ctx))
}

func CGContextSaveGState(ctx CGContextRef) {
	C.CGContextSaveGState(C.CGContextRef(ctx))

}

func CGContextRestoreGState(ctx CGContextRef) {
	C.CGContextRestoreGState(C.CGContextRef(ctx))
}

func CGContextTranslateCTM(ctx CGContextRef, x, y float32) {
	C.CGContextTranslateCTM(C.CGContextRef(ctx), C.CGFloat(x), C.CGFloat(y))
}

func CGContextScaleCTM(ctx CGContextRef, x, y float32) {
	C.CGContextScaleCTM(C.CGContextRef(ctx), C.CGFloat(x), C.CGFloat(y))
}

func CGContextRotateCTM(ctx CGContextRef, angle float32) {
	C.CGContextRotateCTM(C.CGContextRef(ctx), C.CGFloat(_RADIAN*angle))
}

func CGContextClipToRect(ctx CGContextRef, rect CGRect) {
	C.CGContextClipToRect(C.CGContextRef(ctx), C.CGRect(rect))
}

func CGContextClip(ctx CGContextRef) {
	C.CGContextClip(C.CGContextRef(ctx))
}

func CGContextSetAlpha(ctx CGContextRef, alpha float32) {
	C.CGContextSetAlpha(C.CGContextRef(ctx), C.CGFloat(alpha))
}

func CGContextSetLineWidth(ctx CGContextRef, width float32) {
	C.CGContextSetLineWidth(C.CGContextRef(ctx), C.CGFloat(width))
}

func CGContextFillEllipseInRect(ctx CGContextRef, rect CGRect) {
	C.CGContextFillEllipseInRect(C.CGContextRef(ctx), C.CGRect(rect))
}

func CGContextSetRGBFillColor(ctx CGContextRef, red, green, blue, alpha float32) {
	C.CGContextSetRGBFillColor(C.CGContextRef(ctx), C.CGFloat(red), C.CGFloat(green), C.CGFloat(blue), C.CGFloat(alpha))
}

func CGContextSetRGBStrokeColor(ctx CGContextRef, red, green, blue, alpha float32) {
	C.CGContextSetRGBStrokeColor(C.CGContextRef(ctx), C.CGFloat(red), C.CGFloat(green), C.CGFloat(blue), C.CGFloat(alpha))
}

func CGContextStrokePath(ctx CGContextRef) {
	C.CGContextStrokePath(C.CGContextRef(ctx))
}

func CGContextFillPath(ctx CGContextRef) {
	C.CGContextFillPath(C.CGContextRef(ctx))
}

func CGContextFlush(ctx CGContextRef) {
	C.CGContextFlush(C.CGContextRef(ctx))
}

func CGContextSetShadowWithColor(ctx CGContextRef, offset CGSize, blur float32, color CGColorRef) {
	C.CGContextSetShadowWithColor(C.CGContextRef(ctx), C.CGSize(offset), C.CGFloat(blur), C.CGColorRef(color))
}

func CGContextAddArc(ctx CGContextRef, x, y, radius, startAngle, endAngle float32, clockwise int) {
	C.CGContextAddArc(C.CGContextRef(ctx), C.CGFloat(x), C.CGFloat(y), C.CGFloat(radius), C.CGFloat(startAngle), C.CGFloat(endAngle), C.int(clockwise))
}

func CGContextFillRect(ctx CGContextRef, rect CGRect) {
	C.CGContextFillRect(C.CGContextRef(ctx), C.CGRect(rect))
}

func CGContextStrokeRectWithWidth(ctx CGContextRef, rect CGRect, width float32) {
	C.CGContextStrokeRectWithWidth(C.CGContextRef(ctx), C.CGRect(rect), C.CGFloat(width))
}

func CGContextDrawImage(ctx CGContextRef, rect CGRect, image CGImageRef) {
	// C.CGContextDrawImage(C.CGContextRef(ctx), C.CGRect(rect), C.CGImageRef(image))
	C.nux_CGContextDrawImage(C.CGContextRef(ctx), C.CGRect(rect), C.CGImageRef(image))
}

func CGContextSetLineDash(ctx CGContextRef, phase float32, dash []float32, count int) {
	if len(dash) == 0 || count == 0 {
		C.CGContextSetLineDash(C.CGContextRef(ctx), 0, nil, 0)
		return
	}

	d := make([]C.CGFloat, len(dash))
	for i, v := range dash {
		d[i] = C.CGFloat(v)
	}
	C.CGContextSetLineDash(C.CGContextRef(ctx), C.CGFloat(phase), (*C.CGFloat)(unsafe.Pointer(&d[0])), C.size_t(count))
}
