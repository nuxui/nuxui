// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package nux

import ()

type nativeFont struct {
}

func createNativeFont(family string, traits uint32, weight FontWeight, size int32) *nativeFont {
	me := &nativeFont{}
	return me
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
}

func newNativeFontLayout() *nativeFontLayout {
	me := &nativeFontLayout{}
	return me
}

func (me *nativeFontLayout) MeasureText(font Font, paint Paint, text string, width, height int32) (textWidth, textHeight int32) {
	return 800, 800
}

func (me *nativeFontLayout) DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32) {
}

func (me *nativeFontLayout) CharacterIndexForPoint(font Font, text string, width, height int32, x, y float32) uint32 {
	return 0
}
