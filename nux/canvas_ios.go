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
*/
import "C"
import (
	"unsafe"
)

type canvas struct {
	ptr  C.CGContextRef
	clip *RectF
}

func newCanvas(ref C.CGContextRef) *canvas {
	return &canvas{
		ptr:  ref,
		clip: &RectF{0, 0, 99999, 99999},
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

func (me *canvas) ClipRect(left, top, right, bottom float32) {
	C.CGContextClipToRect(me.ptr, C.CGRectMake(C.CGFloat(left), C.CGFloat(top), C.CGFloat(right), C.CGFloat(bottom)))
	me.clip.Left = left
	me.clip.Top = top
	me.clip.Right = right
	me.clip.Bottom = bottom
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}
func (me *canvas) SetAlpha(alpha float32) {
	C.CGContextSetAlpha(me.ptr, C.CGFloat(alpha))
}

func (me *canvas) DrawRect(left, top, right, bottom float32, paint Paint) {
	a, r, g, b := paint.Color().ARGBf()
	rect := C.CGRectMake(C.CGFloat(left), C.CGFloat(top), C.CGFloat(right), C.CGFloat(bottom))
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
}

func (me *canvas) DrawRoundRect(left, top, right, bottom float32, radius float32, paint Paint) {
	// TODO::
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

func (me *canvas) DrawOval(left, top, right, bottom float32, paint Paint) {
	C.CGContextFillEllipseInRect(me.ptr, C.CGRectMake(C.CGFloat(left), C.CGFloat(top), C.CGFloat(right), C.CGFloat(bottom)))
}

func (me *canvas) DrawPath(path Path) {
	// TODO::
}
func (me *canvas) DrawColor(color Color) {
	paint := NewPaint()
	paint.SetColor(color)
	paint.SetStyle(PaintStyle_Fill)
	me.DrawRect(me.clip.Left, me.clip.Top, me.clip.Right, me.clip.Bottom, paint)
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
