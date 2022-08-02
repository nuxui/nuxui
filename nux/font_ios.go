// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/ios"
	"runtime"
)

type nativeFont struct {
	ptr ios.UIFont
}

func createNativeFont(family string, traits uint32, weight FontWeight, size int32) *nativeFont {
	font := ios.UIFont_fontWithName(family, float32(size))
	if font.IsNil() {
		font = ios.UIFont_systemFontOfSize(float32(size), fontWeightToNative(weight))
	}
	me := &nativeFont{
		ptr: font,
	}
	runtime.SetFinalizer(me, freeNativeFont)
	return me
}

func freeNativeFont(me *nativeFont) {
	ios.NSObject_release(uintptr(me.ptr))
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
	layout    ios.NSLayoutManager
	container ios.NSTextContainer
}

func newNativeFontLayout() *nativeFontLayout {
	me := &nativeFontLayout{
		layout:    ios.NewNSLayoutManager(),
		container: ios.NewNSTextContainer(0, 0),
	}
	me.layout.AddTextContainer(me.container)
	runtime.SetFinalizer(me, freeNativeFontLayout)
	return me
}

func freeNativeFontLayout(me *nativeFontLayout) {
	me.layout.RemoveTextContainerAtIndex(0)
	ios.NSObject_release(uintptr(me.container))
	ios.NSObject_release(uintptr(me.layout))
}

func (me *nativeFontLayout) MeasureText(font Font, text string, width, height int32) (textWidth, textHeight int32) {
	w, h := me.layout.MeasureText(me.container, font.native().ptr, text, float32(width), float32(height))
	return int32(w + 0.999999999999999), int32(h + 0.999999999999999)
}

func (me *nativeFontLayout) DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32) {
	me.layout.DrawText(me.container, font.native().ptr, text, float32(width), float32(height), uint32(paint.Color()), 0) //TODO:: bgcolor
}

func (me *nativeFontLayout) CharacterIndexForPoint(font Font, text string, width, height int32, x, y float32) uint32 {
	index, fraction := me.layout.CharacterIndexForPoint(me.container, font.native().ptr, text, float32(width), float32(height), x, y)
	if fraction > 0.5 {
		index++
	}
	return index
}

func fontWeightToNative(weight FontWeight) ios.UIFontWeight {
	switch weight {
	case FontWeight_Thin:
		return ios.UIFontWeightThin
	case FontWeight_ExtraLight:
		return ios.UIFontWeightUltraLight
	case FontWeight_Light:
		return ios.UIFontWeightLight
	case FontWeight_Normal:
		return ios.UIFontWeightRegular
	case FontWeight_Medium:
		return ios.UIFontWeightMedium
	case FontWeight_SemiBold:
		return ios.UIFontWeightSemibold
	case FontWeight_Bold:
		return ios.UIFontWeightBold
	case FontWeight_ExtraBold:
		return ios.UIFontWeightHeavy
	case FontWeight_Black:
		return ios.UIFontWeightBlack
	}
	return ios.UIFontWeightRegular
}
