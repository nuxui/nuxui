// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

/*
#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>

NSUInteger characterIndexForPoint(CGFloat x, CGFloat y, char* text, CGFloat size, CGFloat width, CGFloat height){
	@autoreleasepool{
		NSString *textstr = [NSString stringWithUTF8String:text];
		UIFont* font = [UIFont systemFontOfSize:size];
		NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:textstr
			attributes:@{
				NSFontAttributeName: font,
			}];
		NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init];
		NSTextContainer *textContainer = [[NSTextContainer alloc] initWithSize: CGSizeMake(width, height)];
		[textContainer setLineFragmentPadding:0];
		[layoutManager addTextContainer:textContainer];
		[textStorage addLayoutManager:layoutManager];
		[layoutManager glyphRangeForTextContainer:textContainer];
		CGFloat fraction = 0;
    	NSUInteger index =  [layoutManager characterIndexForPoint: CGPointMake(x,y) inTextContainer:textContainer fractionOfDistanceBetweenInsertionPoints:&fraction];
		if (fraction > 0.5){
			index++;
		}
		return index;
	}
}

void measureText(char* text, CGFloat size, CGFloat width, CGFloat height, CGFloat* outWidth, CGFloat* outHeight){
	@autoreleasepool{
		NSString *textstr = [NSString stringWithUTF8String:text];
		UIFont* font = [UIFont systemFontOfSize:size];
		NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:textstr
			attributes:@{
				NSFontAttributeName: font,
			}];
		NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init];
		NSTextContainer *textContainer = [[NSTextContainer alloc] initWithSize: CGSizeMake(width, height)];
		[textContainer setLineFragmentPadding:0];
		[layoutManager addTextContainer:textContainer];
		[textStorage addLayoutManager:layoutManager];
		[layoutManager glyphRangeForTextContainer:textContainer];

		CGRect rect = [layoutManager usedRectForTextContainer: textContainer];
		*outWidth = rect.size.width;
		*outHeight = rect.size.height;
	}
}

void drawText(char *text, CGFloat width, CGFloat height, uint32_t color, uint32_t bgColor, CGFloat size){
	@autoreleasepool{
		NSString *textstr = [NSString stringWithUTF8String:text];
		CGFloat a = (CGFloat)((color>>24)&0xff) / 255.0;
		CGFloat r = (CGFloat)((color>>16)&0xff) / 255.0;
		CGFloat g = (CGFloat)((color>>8)&0xff) / 255.0;
		CGFloat b = (CGFloat)((color)&0xff) / 255.0;

		CGFloat a0 = (CGFloat)((bgColor>>24)&0xff) / 255.0;
		CGFloat r0 = (CGFloat)((bgColor>>16)&0xff) / 255.0;
		CGFloat g0 = (CGFloat)((bgColor>>8)&0xff) / 255.0;
		CGFloat b0 = (CGFloat)((bgColor)&0xff) / 255.0;

		UIFont* font = [UIFont systemFontOfSize:size];
		NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:textstr
			attributes:@{
				NSFontAttributeName: font,
				NSBackgroundColorAttributeName: [UIColor colorWithRed:r0 green:g0 blue:b0 alpha:a0],
				NSForegroundColorAttributeName: [UIColor colorWithRed:r green:g blue:b alpha:a],
			}];
		NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init];
		NSTextContainer *textContainer = [[NSTextContainer alloc] initWithSize: CGSizeMake(width, height)];
		[textContainer setLineFragmentPadding:0];
		[layoutManager addTextContainer:textContainer];
		[textStorage addLayoutManager:layoutManager];
		NSRange glyphRange = [layoutManager glyphRangeForTextContainer:textContainer];
		[layoutManager drawGlyphsForGlyphRange: glyphRange atPoint: CGPointMake(0,0)];
	}
}

void drawImage(CGContextRef ctx, CGFloat x, CGFloat y, CGFloat width, CGFloat height, CGImageRef image){
	CGContextSaveGState(ctx);
	CGContextTranslateCTM(ctx, 0, height);
	CGContextScaleCTM(ctx, 1, -1);
	CGContextDrawImage(ctx, CGRectMake(x, y, width, height), image);
	CGContextRestoreGState(ctx);
}

void setShadow(CGContextRef ctx, CGFloat x, CGFloat y, CGFloat blur, CGFloat a, CGFloat r, CGFloat g, CGFloat b){
	CGColorRef color = [[UIColor colorWithRed:r green:g blue:b alpha:a] CGColor];
	CGContextSetShadowWithColor(ctx, CGSizeMake(x, -y), blur, color);
}
*/
import "C"
import (
	"unsafe"
)

type canvas struct {
	ptr  C.CGContextRef
}

func newCanvas(ref C.CGContextRef) *canvas {
	return &canvas{
		ptr:  ref,
	}
}

func (me *canvas) ResetClip() {
	C.CGContextResetClip(me.ptr)
}

func (me *canvas) Save() {
	C.CGContextSaveGState(me.ptr)
}

func (me *canvas) Restore() {
	C.CGContextRestoreGState(me.ptr)
}

func (me *canvas) Translate(x, y float32) {
	C.CGContextTranslateCTM(me.ptr, C.CGFloat(x), C.CGFloat(y))
}

func (me *canvas) Scale(x, y float32) {
	C.CGContextScaleCTM(me.ptr, C.CGFloat(x), C.CGFloat(y))
}

func (me *canvas) Rotate(angle float32) {
	C.CGContextRotateCTM(me.ptr, C.CGFloat(angle))
}

func (me *canvas) Skew(x, y float32) {
	// TODO::
}

func (me *canvas) Transform(a, b, c, d, e, f float32) {
	// TODO::
}

func (me *canvas) SetMatrix(matrix Matrix) {
	// TODO::
}

func (me *canvas) GetMatrix() Matrix {
	// TODO::
	return Matrix{}
}

func (me *canvas) ClipRect(x, y, width, height float32) {
	C.CGContextClipToRect(me.ptr, C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height)))
}

func (me *canvas) ClipRoundRect(x, y, width, height, radius float32) {
	rect := C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))
	path := C.CGPathCreateWithRoundedRect(rect, C.CGFloat(radius), C.CGFloat(radius), nil)
	C.CGContextAddPath(me.ptr, path)
	C.CGContextClip(me.ptr)
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}
func (me *canvas) SetAlpha(alpha float32) {
	C.CGContextSetAlpha(me.ptr, C.CGFloat(alpha))
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	a, r, g, b := paint.Color().ARGBf()
	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		x += 1
		y += 1
	}

	if fix {
		me.Save()
		me.Translate(-0.5, -0.5)
	}

	rect := C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))

	switch paint.Style() {
	case PaintStyle_Fill:
		C.CGContextSetRGBFillColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextFillRect(me.ptr, rect)
	case PaintStyle_Stroke:
		C.CGContextSetRGBStrokeColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextStrokeRectWithWidth(me.ptr, rect, C.CGFloat(paint.Width()))
	case PaintStyle_Both:
		C.CGContextSetRGBFillColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextFillRect(me.ptr, rect)
		C.CGContextSetRGBStrokeColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextStrokeRectWithWidth(me.ptr, rect, C.CGFloat(paint.Width()))
	}

	if fix {
		me.Restore()
	}
}

func (me *canvas) DrawRoundRect(x, y, width, height float32, radius float32, paint Paint) {
	a, r, g, b := paint.Color().ARGBf()
	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		x += 1
		y += 1
	}

	if fix {
		me.Save()
		me.Translate(-0.5, -0.5)
	}

	rect := C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))

	path := C.CGPathCreateWithRoundedRect(rect, C.CGFloat(radius), C.CGFloat(radius), nil)
	C.CGContextAddPath(me.ptr, path)
	C.CGContextSetLineWidth(me.ptr, C.CGFloat(paint.Width()))

	hasShadow := false
	if sc, sx, sy, sb := paint.Shadow(); sc != 0 && sb > 0 {
		hasShadow = true
		me.Save()
		a0, r0, g0, b0 := sc.ARGBf()
		C.setShadow(me.ptr, C.CGFloat(sx), C.CGFloat(sy), C.CGFloat(sb), C.CGFloat(a0), C.CGFloat(r0), C.CGFloat(g0), C.CGFloat(b0))
	}

	switch paint.Style() {
	case PaintStyle_Fill:
		C.CGContextSetRGBFillColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextFillPath(me.ptr)
	case PaintStyle_Stroke:
		C.CGContextSetRGBStrokeColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextStrokePath(me.ptr)
	case PaintStyle_Both:
		C.CGContextSetRGBFillColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextFillPath(me.ptr)
		C.CGContextSetRGBStrokeColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
		C.CGContextStrokePath(me.ptr)
	}
	C.CGPathRelease(path)

	if hasShadow {
		me.Restore()
	}

	if fix {
		me.Restore()
	}
}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {
	// TODO:: useCenter
	a, r, g, b := paint.Color().ARGBf()
	C.CGContextSetRGBFillColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))
	C.CGContextSetRGBStrokeColor(me.ptr, C.CGFloat(r), C.CGFloat(g), C.CGFloat(b), C.CGFloat(a))

	var clockwise C.int = 0
	if useCenter {
		clockwise = 1
	}
	C.CGContextAddArc(me.ptr, C.CGFloat(x), C.CGFloat(y), C.CGFloat(radius), C.CGFloat(startAngle), C.CGFloat(endAngle), clockwise)
}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
	C.CGContextFillEllipseInRect(me.ptr, C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height)))
}

func (me *canvas) DrawPath(path Path) {
	// TODO::
}

func (me *canvas) DrawImage(img Image) {
	w, h := img.Size()
	C.drawImage(me.ptr, 0, 0, C.double(w), C.double(h), img.(*nativeImage).ptr)
}

func (me *canvas) DrawText(text string, width, height float32, paint Paint) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.drawText(ctext, C.CGFloat(width), C.CGFloat(height), C.uint32_t(paint.Color()), 0, C.CGFloat(paint.TextSize()))
}

func (me *canvas) Flush() {
	C.CGContextFlush(me.ptr)
}

func (me *canvas) Destroy() {
}

func (me *paint) MeasureText(text string, width, height float32) (outWidth float32, outHeight float32) {
	if text == "" {
		return 0, 0
	}
	
	var w, h C.CGFloat
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.measureText(ctext, C.CGFloat(me.textSize), C.CGFloat(width), C.CGFloat(height), &w, &h)
	return float32(w), float32(h)
}

func (me *paint) CharacterIndexForPoint(text string, width, height float32, x, y float32) uint32 {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	index := C.characterIndexForPoint(C.CGFloat(x), C.CGFloat(y), ctext, C.CGFloat(me.textSize), C.CGFloat(width), C.CGFloat(height))
	return uint32(index)
}
