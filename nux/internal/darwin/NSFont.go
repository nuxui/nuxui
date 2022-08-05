// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

uintptr_t nux_NSFontWithName(char* name, CGFloat size){
	return (uintptr_t)([NSFont fontWithName: [NSString stringWithUTF8String:name] size: size]);
}

uintptr_t nux_NSFontSystemFontOfSize(CGFloat size, NSFontWeight weight){
	return (uintptr_t)([NSFont systemFontOfSize:size weight: weight]);
}

uintptr_t nux_SharedNSFontManager(){
	return (uintptr_t)([NSFontManager sharedFontManager]);
}

uintptr_t nux_NSFontManager_fontWithFamily(uintptr_t fontManager, char *family, NSFontTraitMask traits , NSInteger weight , CGFloat size){
	NSFont* font = [(NSFontManager*)fontManager fontWithFamily: [NSString stringWithUTF8String:family]
                    traits:traits
                    weight:weight
                      size:size];
	return (uintptr_t)font;
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
			NSFontAttributeName: (NSFont*)font,
			NSForegroundColorAttributeName: [NSColor colorWithSRGBRed:fr green:fg blue:fb alpha:fa],
			NSBackgroundColorAttributeName: [NSColor colorWithSRGBRed:br green:bg blue:bb alpha:ba],
		}];
	return (uintptr_t)textStorage;
}

uintptr_t nux_NewNSLayoutManager(){
	return (uintptr_t)[[NSLayoutManager alloc] init];
}

uintptr_t nux_NewNSTextContainer(CGFloat width, CGFloat height){
	return (uintptr_t)[[NSTextContainer alloc] initWithSize: NSMakeSize(width, height)];

}

void nux_NSTextContainer_setSize(uintptr_t nsTextContainer, CGFloat width, CGFloat height){
	[(NSTextContainer*)nsTextContainer setSize: NSMakeSize(width, height)];
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
			NSFontAttributeName: (NSFont*)font,
		}];
	[textStorage addLayoutManager:layoutManager];

	NSTextContainer *textContainer = (NSTextContainer*)nsTextContainer;
    [textContainer setLineFragmentPadding:0];
	[textContainer setSize:NSMakeSize(width, height)];
	[layoutManager glyphRangeForTextContainer:textContainer];

	NSRect rect = [layoutManager usedRectForTextContainer: textContainer];
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
			NSFontAttributeName: (NSFont*)font,
			NSForegroundColorAttributeName: [NSColor colorWithSRGBRed:fr green:fg blue:fb alpha:fa],
			NSBackgroundColorAttributeName: [NSColor colorWithSRGBRed:br green:bg blue:bb alpha:ba],
		}];
	[textStorage addLayoutManager:layoutManager];

	NSTextContainer *textContainer = (NSTextContainer*)nsTextContainer;
    [textContainer setLineFragmentPadding:0];
	[textContainer setSize:NSMakeSize(width, height)];
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
			NSFontAttributeName: (NSFont*)font,
		}];
	[textStorage addLayoutManager:layoutManager];

	NSTextContainer *textContainer = (NSTextContainer*)nsTextContainer;
    [textContainer setLineFragmentPadding:0];
	[textContainer setSize:NSMakeSize(width, height)];
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

type NSFontWeight C.NSFontWeight

const (
	NSFontWeightUltraLight NSFontWeight = -0.800000
	NSFontWeightThin       NSFontWeight = -0.600000
	NSFontWeightLight      NSFontWeight = -0.400000
	NSFontWeightRegular    NSFontWeight = 0.000000
	NSFontWeightMedium     NSFontWeight = 0.230000
	NSFontWeightSemibold   NSFontWeight = 0.300000
	NSFontWeightBold       NSFontWeight = 0.400000
	NSFontWeightHeavy      NSFontWeight = 0.560000
	NSFontWeightBlack      NSFontWeight = 0.620000
)

type FontMask C.NSFontTraitMask

const (
	NSBoldFontMask                    FontMask = C.NSBoldFontMask
	NSCompressedFontMask              FontMask = C.NSCompressedFontMask
	NSCondensedFontMask               FontMask = C.NSCondensedFontMask
	NSExpandedFontMask                FontMask = C.NSExpandedFontMask
	NSFixedPitchFontMask              FontMask = C.NSFixedPitchFontMask
	NSItalicFontMask                  FontMask = C.NSItalicFontMask
	NSNarrowFontMask                  FontMask = C.NSNarrowFontMask
	NSNonStandardCharacterSetFontMask FontMask = C.NSNonStandardCharacterSetFontMask
	NSPosterFontMask                  FontMask = C.NSPosterFontMask
	NSSmallCapsFontMask               FontMask = C.NSSmallCapsFontMask
	NSUnboldFontMask                  FontMask = C.NSUnboldFontMask
	NSUnitalicFontMask                FontMask = C.NSUnitalicFontMask
)

func NSFontWithName(name string, size float32) NSFont {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (NSFont)(C.nux_NSFontWithName(cname, C.CGFloat(size)))
}

func NSFontSystemFontOfSize(size float32, weight NSFontWeight) NSFont {
	return (NSFont)(C.nux_NSFontSystemFontOfSize(C.CGFloat(size), C.NSFontWeight(weight)))
}

func SharedNSFontManager() NSFontManager {
	return (NSFontManager)(C.nux_SharedNSFontManager())
}

func (me NSFontManager) FontWithFamily(family string, traits uint32, weight int32, size float32) NSFont {
	cfamily := C.CString(family)
	defer C.free(unsafe.Pointer(cfamily))
	return (NSFont)(C.nux_NSFontManager_fontWithFamily(C.uintptr_t(me),
		cfamily,
		C.NSFontTraitMask(traits),
		C.NSInteger(weight),
		C.CGFloat(size)))
}

func NewNSTextStorage(font NSFont, text string, fgcolor, bgcolor uint32) NSTextStorage {
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

func (me NSLayoutManager) MeasureText(textContainer NSTextContainer, font NSFont, text string, width, height float32) (textWidth, textHeight float32) {
	var w, h C.CGFloat
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.nux_NSLayoutManager_measureText(C.uintptr_t(me), C.uintptr_t(textContainer), C.uintptr_t(font), ctext, C.CGFloat(width), C.CGFloat(height), &w, &h)
	return float32(w), float32(h)
}

func (me NSLayoutManager) DrawText(textContainer NSTextContainer, font NSFont, text string, width, height float32, fgColor, bgColor uint32) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.nux_NSLayoutManager_drawText(C.uintptr_t(me), C.uintptr_t(textContainer), C.uintptr_t(font), ctext, C.CGFloat(width), C.CGFloat(height), C.uint32_t(fgColor), C.uint32_t(bgColor))
}

func (me NSLayoutManager) CharacterIndexForPoint(textContainer NSTextContainer, font NSFont, text string, width, height, x, y float32) (index uint32, fraction float32) {
	f := C.CGFloat(0)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	i := C.nux_NSLayoutManager_characterIndexForPoint(C.uintptr_t(me), C.uintptr_t(textContainer), C.uintptr_t(font), ctext, C.CGFloat(width), C.CGFloat(height), C.CGFloat(x), C.CGFloat(y), &f)
	return uint32(i), float32(f)
}
