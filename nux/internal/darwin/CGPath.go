// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

#define _RADIAN (3.1415926535897932384626433832795028841971/180.0)

void CGPathAddRoundRectPath_(CGMutablePathRef path, CGFloat x, CGFloat y, CGFloat width, CGFloat height,
	CGFloat rLT, CGFloat rRT, CGFloat rRB, CGFloat rLB){
	CGPathAddArc(path, nil, x+width-rRT, y+rRT, rRT, -90*_RADIAN, 0, false);
	CGPathAddArc(path, nil, x+width-rRB, y+height-rRB, rRB, 0, 90*_RADIAN, false);
	CGPathAddArc(path, nil, x+rLB, y+height-rLB, rLB, 90*_RADIAN, 180*_RADIAN, false);
	CGPathAddArc(path, nil, x+rLT, y+rLT, rLT, 180*_RADIAN, 270*_RADIAN, false);
}
*/
import "C"

func CGPathCreateMutable() CGMutablePath {
	return CGMutablePath(C.CGPathCreateMutable())
}

func CGContextAddPath(ctx CGContext, path CGPath) {
	C.CGContextAddPath(C.CGContextRef(ctx), C.CGPathRef(path))
}

// https://developer.apple.com/documentation/coregraphics/1411218-cgpathcreatewithroundedrect/
func CGPathCreateWithRoundedRect(rect CGRect, cornerWidth, cornerHeight float32, transform *CGAffineTransform) CGPath {
	return CGPath(C.CGPathCreateWithRoundedRect(C.CGRect(rect), C.CGFloat(cornerWidth), C.CGFloat(cornerHeight), (*C.CGAffineTransform)(transform)))
}

func CGPathAddRoundRectPath(path CGMutablePath, x, y, width, height, rLT, rRT, rRB, rLB float32) {
	C.CGPathAddRoundRectPath_(C.CGMutablePathRef(path), C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height), C.CGFloat(rLT), C.CGFloat(rRT), C.CGFloat(rRB), C.CGFloat(rLB))
}

func CGPathCloseSubpath(path CGMutablePath) {
	C.CGPathCloseSubpath(C.CGMutablePathRef(path))
}

func CGPathRelease(path CGPath) {
	C.CGPathRelease(C.CGPathRef(path))
}
