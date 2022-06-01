// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

void nux_NSDrawText(char *text, CGFloat width, CGFloat height, uint32_t fgColor, uint32_t bgColor, CGFloat size){
	@autoreleasepool{
		NSString *textstr = [NSString stringWithUTF8String:text];
		CGFloat fr = (CGFloat)((fgColor>>16)&0xff) / 255.0;
		CGFloat fg = (CGFloat)((fgColor>>8)&0xff)  / 255.0;
		CGFloat fb = (CGFloat)((fgColor)&0xff)     / 255.0;
		CGFloat fa = (CGFloat)((fgColor>>24)&0xff) / 255.0;

		CGFloat br = (CGFloat)((bgColor>>16)&0xff) / 255.0;
		CGFloat bg = (CGFloat)((bgColor>>8)&0xff)  / 255.0;
		CGFloat bb = (CGFloat)((bgColor)&0xff)     / 255.0;
		CGFloat ba = (CGFloat)((bgColor>>24)&0xff) / 255.0;

		NSFont* font = [NSFont systemFontOfSize:size];
		NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:textstr
			attributes:@{
				NSFontAttributeName: font,
				NSForegroundColorAttributeName: [NSColor colorWithSRGBRed:fr green:fg blue:fb alpha:fa],
				NSBackgroundColorAttributeName: [NSColor colorWithSRGBRed:br green:bg blue:bb alpha:ba],
			}];
		NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init];
		NSTextContainer *textContainer = [[NSTextContainer alloc] initWithSize: NSMakeSize(width, height)];
		[textContainer setLineFragmentPadding:0];
		[layoutManager addTextContainer:textContainer];
		[textStorage addLayoutManager:layoutManager];
		NSRange glyphRange = [layoutManager glyphRangeForTextContainer:textContainer];
		[layoutManager drawGlyphsForGlyphRange: glyphRange atPoint: CGPointMake(0,0)];
	}
}

NSUInteger nux_NSCharacterIndexForPoint(CGFloat x, CGFloat y, char* text, CGFloat size, CGFloat width, CGFloat height){
	@autoreleasepool{
		NSString *textstr = [NSString stringWithUTF8String:text];
		NSFont* font = [NSFont systemFontOfSize:size];
		NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:textstr
			attributes:@{
				NSFontAttributeName: font,
			}];
		NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init];
		NSTextContainer *textContainer = [[NSTextContainer alloc] initWithSize: NSMakeSize(width, height)];
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

void nux_NSMeasureText(char* text, CGFloat size, CGFloat width, CGFloat height, CGFloat* outWidth, CGFloat* outHeight){
	@autoreleasepool{
		NSString *textstr = [NSString stringWithUTF8String:text];
		NSFont* font = [NSFont systemFontOfSize:size];
		NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:textstr
			attributes:@{
				NSFontAttributeName: font,
			}];
		NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init];
		NSTextContainer *textContainer = [[NSTextContainer alloc] initWithSize: NSMakeSize(width, height)];
		[textContainer setLineFragmentPadding:0];
		[layoutManager addTextContainer:textContainer];
		[textStorage addLayoutManager:layoutManager];
		[layoutManager glyphRangeForTextContainer:textContainer];

		NSRect rect = [layoutManager usedRectForTextContainer: textContainer];
		*outWidth = rect.size.width;
		*outHeight = rect.size.height;
	}
}
*/
import "C"

import (
	"unsafe"
)

func NSDrawText(text string, width, height float32, fgColor, bgColor uint32, txtSize float32) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.nux_NSDrawText(ctext, C.CGFloat(width), C.CGFloat(height), C.uint32_t(fgColor), C.uint32_t(bgColor), C.CGFloat(txtSize))
}

func NSMeasureText(text string, width, height float32, fontSize float32) (textWidth, textHeight float32) {
	if text == "" {
		return 0, 0
	}
	var w, h C.CGFloat
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.nux_NSMeasureText(ctext, C.CGFloat(fontSize), C.CGFloat(width), C.CGFloat(height), &w, &h)
	return float32(w), float32(h)
}

func NSCharacterIndexForPoint(text string, width, height float32, x, y float32, fontSize float32) uint32 {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	index := C.nux_NSCharacterIndexForPoint(C.CGFloat(x), C.CGFloat(y), ctext, C.CGFloat(fontSize), C.CGFloat(width), C.CGFloat(height))
	return uint32(index)
}
