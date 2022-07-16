// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

#define _RADIAN (3.1415926535897932384626433832795028841971/180.0)

void nux_CGPathAddRoundRectPath(CGMutablePathRef path, CGFloat x, CGFloat y, CGFloat width, CGFloat height,
	CGFloat rLT, CGFloat rRT, CGFloat rRB, CGFloat rLB){
	CGPathAddArc(path, nil, x+width-rRT, y+rRT, rRT, -90*_RADIAN, 0, false);
	CGPathAddArc(path, nil, x+width-rRB, y+height-rRB, rRB, 0, 90*_RADIAN, false);
	CGPathAddArc(path, nil, x+rLB, y+height-rLB, rLB, 90*_RADIAN, 180*_RADIAN, false);
	CGPathAddArc(path, nil, x+rLT, y+rLT, rLT, 180*_RADIAN, 270*_RADIAN, false);
}
*/
import "C"

func CGPathCreateMutable() CGMutablePathRef {
	return CGMutablePathRef(C.CGPathCreateMutable())
}

func CGContextAddPath(ctx CGContextRef, path CGPathRef) {
	C.CGContextAddPath(C.CGContextRef(ctx), C.CGPathRef(path))
}

// https://developer.apple.com/documentation/coregraphics/1411218-cgpathcreatewithroundedrect/
func CGPathCreateWithRoundedRect(rect CGRect, cornerWidth, cornerHeight float32, transform *CGAffineTransform) CGPathRef {
	return CGPathRef(C.CGPathCreateWithRoundedRect(C.CGRect(rect), C.CGFloat(cornerWidth), C.CGFloat(cornerHeight), (*C.CGAffineTransform)(transform)))
}

func CGPathAddRoundRectPath(path CGMutablePathRef, x, y, width, height, rLT, rRT, rRB, rLB float32) {
	C.nux_CGPathAddRoundRectPath(C.CGMutablePathRef(path), C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height), C.CGFloat(rLT), C.CGFloat(rRT), C.CGFloat(rRB), C.CGFloat(rLB))
}

func CGPathAddRoundedRect(path CGMutablePathRef, transform *CGAffineTransform, rect CGRect, cornerWidth, cornerHeight float32) {
	C.CGPathAddRoundedRect(C.CGMutablePathRef(path), (*C.CGAffineTransform)(transform), C.CGRect(rect), C.CGFloat(cornerWidth), C.CGFloat(cornerHeight))
}

func CGPathCloseSubpath(path CGMutablePathRef) {
	C.CGPathCloseSubpath(C.CGMutablePathRef(path))
}

func CGPathRelease(path CGPathRef) {
	C.CGPathRelease(C.CGPathRef(path))
}

func CGPathAddRect(path CGMutablePathRef, transform *CGAffineTransform, rect CGRect) {
	C.CGPathAddRect(C.CGMutablePathRef(path), (*C.CGAffineTransform)(transform), C.CGRect(rect))
}

func CGPathAddEllipseInRect(path CGMutablePathRef, transform *CGAffineTransform, rect CGRect) {
	C.CGPathAddEllipseInRect(C.CGMutablePathRef(path), (*C.CGAffineTransform)(transform), C.CGRect(rect))
}

func CGPathMoveToPoint(path CGMutablePathRef, transform *CGAffineTransform, x, y float32) {
	C.CGPathMoveToPoint(C.CGMutablePathRef(path), (*C.CGAffineTransform)(transform), C.CGFloat(x), C.CGFloat(y))
}

func CGPathAddLineToPoint(path CGMutablePathRef, transform *CGAffineTransform, x, y float32) {
	C.CGPathAddLineToPoint(C.CGMutablePathRef(path), (*C.CGAffineTransform)(transform), C.CGFloat(x), C.CGFloat(y))
}

func CGPathAddCurveToPoint(path CGMutablePathRef, transform *CGAffineTransform, x1, y1, x2, y2, x3, y3 float32) {
	C.CGPathAddCurveToPoint(C.CGMutablePathRef(path), (*C.CGAffineTransform)(transform), C.CGFloat(x1), C.CGFloat(y1), C.CGFloat(x2), C.CGFloat(y2), C.CGFloat(x3), C.CGFloat(y3))
}

func CGPathAddQuadCurveToPoint(path CGMutablePathRef, transform *CGAffineTransform, x2, y2, x3, y3 float32) {
	C.CGPathAddQuadCurveToPoint(C.CGMutablePathRef(path), (*C.CGAffineTransform)(transform), C.CGFloat(x2), C.CGFloat(y2), C.CGFloat(x3), C.CGFloat(y3))
}
