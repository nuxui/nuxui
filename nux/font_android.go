// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/nux/internal/android"
	"runtime"
)

type nativeFont struct {
}

func createNativeFont(family string, traits uint32, weight FontWeight, size int32) *nativeFont {
	me := &nativeFont{}
	runtime.SetFinalizer(me, freeNativeFont)
	return me
}

func freeNativeFont(me *nativeFont) {
	// darwin.NSObject_release(uintptr(me.ptr))
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
	// layout android.StaticLayout
}

func newNativeFontLayout() *nativeFontLayout {
	me := &nativeFontLayout{
		// layout: 0,
	}
	runtime.SetFinalizer(me, freeNativeFontLayout)
	return me
}

func freeNativeFontLayout(me *nativeFontLayout) {
	// if !me.layout.IsNil() {
	// android.DeleteGlobalRef(android.JObject(me.layout))
	// }
}

func (me *nativeFontLayout) MeasureText(font Font, paint Paint, text string, width, height int32) (textWidth, textHeight int32) {
	// if me.layout.IsNil() {
	paint.native().ref.SetTextSize(float32(font.Size()) * android.GetDisplayMetrics().ScaledDensity)
	layout := android.NewStaticLayout(text, width, paint.native().ref)
	defer android.DeleteGlobalRef(android.JObject(layout))
	// }
	return layout.GetSize()
}

func (me *nativeFontLayout) DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32) {
	// if me.layout.IsNil() {
	paint.native().ref.SetTextSize(float32(font.Size()) * android.GetDisplayMetrics().ScaledDensity)
	layout := android.NewStaticLayout(text, width, paint.native().ref)
	// }
	defer android.DeleteGlobalRef(android.JObject(layout))
	layout.Draw(canvas.native().ref)
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
