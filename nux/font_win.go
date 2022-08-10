// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && !cairo

package nux

import (
	"nuxui.org/nuxui/nux/internal/win32"
	"runtime"
)

// http://etutorials.org/Programming/visual-c-sharp/Part+III+Programming+Windows+Forms/Chapter+14+GDI/Using+Fonts/
type nativeFont struct {
	font   *win32.GpFont
	family *win32.GpFontFamily
}

func createNativeFont(family string, traits uint32, weight FontWeight, size int32) *nativeFont {
	me := &nativeFont{
		font:   &win32.GpFont{},
		family: &win32.GpFontFamily{},
	}

	win32.GdipGetGenericFontFamilySansSerif(&me.family)
	win32.GdipCreateFont(me.family, float32(size), 0, 0, &me.font)
	runtime.SetFinalizer(me, freeNativeFont)
	return me
}

func freeNativeFont(me *nativeFont) {

}

func (me *nativeFont) SetFamily(family string) {
	// me.fd.SetFamily(family)
}

func (me *nativeFont) SetSize(size int32) {
	// me.fd.SetSize(size * pango.Scale)
}

func (me *nativeFont) SetWeight(weight int32) {
	// me.fd.SetWeight(pango.WEIGHT_NORMAL) //TODO::
}

type nativeFontLayout struct {
	// layout    darwin.NSLayoutManager
	// container darwin.NSTextContainer
}

func newNativeFontLayout() *nativeFontLayout {
	me := &nativeFontLayout{
		// layout:    darwin.NewNSLayoutManager(),
		// container: darwin.NewNSTextContainer(0, 0),
	}
	// me.layout.AddTextContainer(me.container)
	// runtime.SetFinalizer(me, freeNativeFontLayout)
	return me
}

func freeNativeFontLayout(me *nativeFontLayout) {
	// me.layout.RemoveTextContainerAtIndex(0)
	// darwin.NSObject_release(uintptr(me.container))
	// darwin.NSObject_release(uintptr(me.layout))
}

var globalDC uintptr
var globalGP *win32.GpGraphics

func init() {
	globalDC, _ = win32.GetDC(0)
	globalGP = &win32.GpGraphics{}
	win32.GdipCreateFromHDC(globalDC, &globalGP)
}

func (me *nativeFontLayout) MeasureText(font Font, text string, width, height int32) (textWidth, textHeight int32) {
	if text == "" {
		return 0, 0
	}
	rect := &win32.RectF{0, 0, float32(width), float32(height)}
	size := &win32.RectF{}
	win32.GdipMeasureString(globalGP, text, font.native().font, rect, nil, size, nil, nil)
	return int32(math.Ceil(float64(size.Width))), int32(math.Ceil(float64(size.Height)))
}

func (me *nativeFontLayout) DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32) {
	rect := &win32.RectF{0, 0, float32(width), float32(height)}
	brush := &win32.GpBrush{}
	win32.GdipCreateSolidFill(win32.ARGB(paint.Color().ARGB()), &brush)
	win32.GdipDrawString(canvas.native().ptr, text, font.native().font, rect, nil, brush)
}

func (me *nativeFontLayout) CharacterIndexForPoint(font Font, text string, width, height int32, x, y float32) uint32 {
	// index, fraction := me.layout.CharacterIndexForPoint(me.container, font.native().ptr, text, float32(width), float32(height), x, y)
	// if fraction > 0.5 {
	// 	index++
	// }
	// return index
	return 0
}

// func fontWeightToNative(weight FontWeight) darwin.NSFontWeight {
// 	switch weight {
// 	case FontWeight_Thin:
// 		return darwin.NSFontWeightThin
// 	case FontWeight_ExtraLight:
// 		return darwin.NSFontWeightUltraLight
// 	case FontWeight_Light:
// 		return darwin.NSFontWeightLight
// 	case FontWeight_Normal:
// 		return darwin.NSFontWeightRegular
// 	case FontWeight_Medium:
// 		return darwin.NSFontWeightMedium
// 	case FontWeight_SemiBold:
// 		return darwin.NSFontWeightSemibold
// 	case FontWeight_Bold:
// 		return darwin.NSFontWeightBold
// 	case FontWeight_ExtraBold:
// 		return darwin.NSFontWeightHeavy
// 	case FontWeight_Black:
// 		return darwin.NSFontWeightBlack
// 	}
// 	return darwin.NSFontWeightRegular
// }
