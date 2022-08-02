// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ios

/*
#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>
#import <CoreGraphics/CoreGraphics.h>

uintptr_t nux_UIFont_fontWithName(char* name, CGFloat size){
	return (uintptr_t)[UIFont fontWithName: [NSString stringWithUTF8String:name] size: size];
}

uintptr_t nux_UIFont_systemFontOfSize(CGFloat size, UIFontWeight weight){
	return (uintptr_t)[UIFont systemFontOfSize: size weight: weight];
}

uintptr_t nux_NewNSTextStorage(uintptr_t font, char* text, uint32_t fgColor, uint32_t bgColor){
	CGFloat fr = (CGFloat)((fgColor>>24)&0xff) / 255.0;
	CGFloat fg = (CGFloat)((fgColor>>16)&0xff) / 255.0;
	CGFloat fb = (CGFloat)((fgColor>>8)&0xff)  / 255.0;
	CGFloat fa = (CGFloat)((fgColor)&0xff)     / 255.0;

	CGFloat br = (CGFloat)((bgColor>>24)&0xff) / 255.0;
	CGFloat bg = (CGFloat)((bgColor>>16)&0xff) / 255.0;
	CGFloat bb = (CGFloat)((bgColor>>8)&0xff)  / 255.0;
	CGFloat ba = (CGFloat)((bgColor)&0xff)     / 255.0;

	NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:[NSString stringWithUTF8String:text]
		attributes:@{
			NSFontAttributeName: (UIFont*)font,
			NSForegroundColorAttributeName: [UIColor colorWithRed:fr green:fg blue:fb alpha:fa],
			NSBackgroundColorAttributeName: [UIColor colorWithRed:br green:bg blue:bb alpha:ba],
		}];
	return (uintptr_t)textStorage;
}

uintptr_t nux_NewNSLayoutManager(){
	return (uintptr_t)[[NSLayoutManager alloc] init];
}

uintptr_t nux_NewNSTextContainer(CGFloat width, CGFloat height){
	return (uintptr_t)[[NSTextContainer alloc] initWithSize: CGSizeMake(width, height)];
}

void nux_NSTextContainer_setSize(uintptr_t nsTextContainer, CGFloat width, CGFloat height){
	[(NSTextContainer*)nsTextContainer setSize: CGSizeMake(width, height)];
}

void nux_NSLayoutManager_addTextContainer(uintptr_t nsLayoutManager, uintptr_t nsTextContainer){
	[(NSLayoutManager*)nsLayoutManager addTextContainer:(NSTextContainer*)nsTextContainer];
}

void nux_NSLayoutManager_removeTextContainerAtIndex(uintptr_t nsLayoutManager, NSUInteger index){
	[(NSLayoutManager*)nsLayoutManager removeTextContainerAtIndex:index];
}

void nux_NSLayoutManager_measureText(uintptr_t nsLayoutManager, uintptr_t nsTextContainer,
	uintptr_t font, char* text, CGFloat width, CGFloat height, CGFloat* outWidth, CGFloat* outHeight){

	NSLayoutManager *layoutManager = (NSLayoutManager*)nsLayoutManager;
	NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:[NSString stringWithUTF8String:text]
		attributes:@{
			NSFontAttributeName: (UIFont*)font,
		}];
	[textStorage addLayoutManager:layoutManager];

	NSTextContainer *textContainer = (NSTextContainer*)nsTextContainer;
    [textContainer setLineFragmentPadding:0];
	[textContainer setSize:CGSizeMake(width, height)];
	[layoutManager glyphRangeForTextContainer:textContainer];

	CGRect rect = [layoutManager usedRectForTextContainer: textContainer];
	*outWidth = rect.size.width;
	*outHeight = rect.size.height;

	[textStorage removeLayoutManager:layoutManager];
	[textStorage release];
}

void nux_NSLayoutManager_drawText(uintptr_t nsLayoutManager, uintptr_t nsTextContainer,
	uintptr_t font, char* text, CGFloat width, CGFloat height, uint32_t fgColor, uint32_t bgColor){

	CGFloat fr = (CGFloat)((fgColor>>24)&0xff) / 255.0;
	CGFloat fg = (CGFloat)((fgColor>>16)&0xff) / 255.0;
	CGFloat fb = (CGFloat)((fgColor>>8)&0xff)  / 255.0;
	CGFloat fa = (CGFloat)((fgColor)&0xff)     / 255.0;

	CGFloat br = (CGFloat)((bgColor>>24)&0xff) / 255.0;
	CGFloat bg = (CGFloat)((bgColor>>16)&0xff) / 255.0;
	CGFloat bb = (CGFloat)((bgColor>>8)&0xff)  / 255.0;
	CGFloat ba = (CGFloat)((bgColor)&0xff)     / 255.0;

	NSLayoutManager *layoutManager = (NSLayoutManager*)nsLayoutManager;
	NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:[NSString stringWithUTF8String:text]
		attributes:@{
			NSFontAttributeName: (UIFont*)font,
			NSForegroundColorAttributeName: [UIColor colorWithRed:fr green:fg blue:fb alpha:fa],
			NSBackgroundColorAttributeName: [UIColor colorWithRed:br green:bg blue:bb alpha:ba],
		}];
	[textStorage addLayoutManager:layoutManager];

	NSTextContainer *textContainer = (NSTextContainer*)nsTextContainer;
    [textContainer setLineFragmentPadding:0];
	[textContainer setSize:CGSizeMake(width, height)];
    NSRange glyphRange = [layoutManager glyphRangeForTextContainer:textContainer];
	[layoutManager drawBackgroundForGlyphRange: glyphRange atPoint: CGPointMake(0,0)];
    [layoutManager drawGlyphsForGlyphRange: glyphRange atPoint: CGPointMake(0,0)];

	[textStorage removeLayoutManager:layoutManager];
	[textStorage release];
}

NSUInteger nux_NSLayoutManager_characterIndexForPoint(uintptr_t nsLayoutManager, uintptr_t nsTextContainer,
	uintptr_t font, char* text, CGFloat width, CGFloat height, CGFloat x, CGFloat y, CGFloat* fraction){

		NSLayoutManager *layoutManager = (NSLayoutManager*)nsLayoutManager;
	NSTextStorage *textStorage = [[NSTextStorage alloc]initWithString:[NSString stringWithUTF8String:text]
		attributes:@{
			NSFontAttributeName: (UIFont*)font,
		}];
	[textStorage addLayoutManager:layoutManager];

	NSTextContainer *textContainer = (NSTextContainer*)nsTextContainer;
    [textContainer setLineFragmentPadding:0];
	[textContainer setSize:CGSizeMake(width, height)];
	[layoutManager glyphRangeForTextContainer:textContainer];
	NSUInteger index = [layoutManager characterIndexForPoint: CGPointMake(x,y)
		inTextContainer:textContainer fractionOfDistanceBetweenInsertionPoints:fraction];

	[textStorage removeLayoutManager:layoutManager];
	[textStorage release];
	return index;
}
*/
import "C"
import (
	"unsafe"
)

type UIFontWeight C.UIFontWeight

const (
	UIFontWeightUltraLight UIFontWeight = -0.800000
	UIFontWeightThin       UIFontWeight = -0.600000
	UIFontWeightLight      UIFontWeight = -0.400000
	UIFontWeightRegular    UIFontWeight = 0.000000
	UIFontWeightMedium     UIFontWeight = 0.230000
	UIFontWeightSemibold   UIFontWeight = 0.300000
	UIFontWeightBold       UIFontWeight = 0.400000
	UIFontWeightHeavy      UIFontWeight = 0.560000
	UIFontWeightBlack      UIFontWeight = 0.620000
)

func UIFont_fontWithName(name string, size float32) UIFont {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	return UIFont(C.nux_UIFont_fontWithName(cstr, C.CGFloat(size)))
}

func UIFont_systemFontOfSize(size float32, weight UIFontWeight) UIFont {
	return UIFont(C.nux_UIFont_systemFontOfSize(C.CGFloat(size), C.UIFontWeight(weight)))
}

func NewNSTextStorage(font UIFont, text string, fgcolor, bgcolor uint32) NSTextStorage {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return NSTextStorage(C.nux_NewNSTextStorage((C.uintptr_t)(font), ctext, C.uint32_t(fgcolor), C.uint32_t(bgcolor)))
}

func NewNSTextContainer(width, height float32) NSTextContainer {
	return NSTextContainer(C.nux_NewNSTextContainer(C.CGFloat(width), C.CGFloat(height)))
}

func (me NSTextContainer) SetSize(width, height float32) {
	C.nux_NSTextContainer_setSize(C.uintptr_t(me), C.CGFloat(width), C.CGFloat(height))
}

func NewNSLayoutManager() NSLayoutManager {
	return NSLayoutManager(C.nux_NewNSLayoutManager())
}

func (me NSLayoutManager) AddTextContainer(textContainer NSTextContainer) {
	C.nux_NSLayoutManager_addTextContainer(C.uintptr_t(me), C.uintptr_t(textContainer))
}

func (me NSLayoutManager) RemoveTextContainerAtIndex(index uint32) {
	C.nux_NSLayoutManager_removeTextContainerAtIndex(C.uintptr_t(me), C.NSUInteger(index))
}

func (me NSLayoutManager) MeasureText(textContainer NSTextContainer, font UIFont, text string, width, height float32) (textWidth, textHeight float32) {
	var w, h C.CGFloat
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.nux_NSLayoutManager_measureText(C.uintptr_t(me), C.uintptr_t(textContainer), C.uintptr_t(font), ctext, C.CGFloat(width), C.CGFloat(height), &w, &h)
	return float32(w), float32(h)
}

func (me NSLayoutManager) DrawText(textContainer NSTextContainer, font UIFont, text string, width, height float32, fgColor, bgColor uint32) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.nux_NSLayoutManager_drawText(C.uintptr_t(me), C.uintptr_t(textContainer), C.uintptr_t(font), ctext, C.CGFloat(width), C.CGFloat(height), C.uint32_t(fgColor), C.uint32_t(bgColor))
}

func (me NSLayoutManager) CharacterIndexForPoint(textContainer NSTextContainer, font UIFont, text string, width, height, x, y float32) (index uint32, fraction float32) {
	var f C.CGFloat
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	i := C.nux_NSLayoutManager_characterIndexForPoint(C.uintptr_t(me), C.uintptr_t(textContainer), C.uintptr_t(font), ctext, C.CGFloat(width), C.CGFloat(height), C.CGFloat(x), C.CGFloat(y), &f)
	return uint32(i), float32(f)
}
