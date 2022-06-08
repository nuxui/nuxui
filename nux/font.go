// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "nuxui.org/nuxui/log"

type FontWeight int

const (
	FontWeight_Thin       FontWeight = 100
	FontWeight_ExtraLight FontWeight = 200
	FontWeight_Light      FontWeight = 300
	FontWeight_Normal     FontWeight = 400
	FontWeight_Medium     FontWeight = 500
	FontWeight_SemiBold   FontWeight = 600
	FontWeight_Bold       FontWeight = 700
	FontWeight_ExtraBold  FontWeight = 800
	FontWeight_Black      FontWeight = 900
)

type Font interface {
	Family() string
	// SetFamily(family string)
	Size() int32
	SetSize(size int32)

	native() *nativeFont
}

/*
family:
size:
weight: string or int
if font not exist rollback to default system font
*/
func NewFont(attr Attr) Font {
	me := &font{
		family: attr.GetString("family", "sans-serif"),
		weight: FontWeightFromName(attr.GetString("weight", "normal")),
		size:   attr.GetInt32("size", 14),
	}
	me.nativeFont = createNativeFont_(me.family, me.traits, me.weight, me.size)
	return me
}

func createNativeFont_(family string, traits uint32, weight FontWeight, size int32) *nativeFont {
	return createNativeFont(family, traits, weight, size)
}

type font struct {
	nativeFont *nativeFont
	family     string
	size       int32
	traits     uint32
	weight     FontWeight
}

func (me *font) native() *nativeFont {
	return me.nativeFont
}

func (me *font) Family() string {
	return me.family
}

// func (me *font) SetFamily(family string) {
// 	me.family = family
// 	me.nativeFont.SetFamily(family)
// }

func (me *font) Size() int32 {
	return me.size
}

func (me *font) SetSize(size int32) {
	me.size = size
}

func FontWeightFromName(name string) FontWeight {
	switch name {
	case "thin":
		return FontWeight_Thin
	case "extraLight":
		return FontWeight_ExtraLight
	case "light":
		return FontWeight_Light
	case "normal":
		return FontWeight_Normal
	case "medium":
		return FontWeight_Medium
	case "semiBold":
		return FontWeight_SemiBold
	case "bold":
		return FontWeight_Bold
	case "extraBold":
		return FontWeight_ExtraBold
	case "black":
		return FontWeight_Black
	default:
		log.E("nuxui", "unknow font weight name %s, go back to normal", name)
		return FontWeight_Normal
	}
}

type FontLayout interface {
	MeasureText(font Font, text string, width, height int32) (textWidth, textHeight int32)
	DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32)
	CharacterIndexForPoint(font Font, text string, width, height int32, x, y float32) uint32

	native() *nativeFontLayout
}

func NewFontLayout() FontLayout {
	me := &fontLayout{
		naiveLayout: newNativeFontLayout(),
	}
	return me
}

type fontLayout struct {
	naiveLayout *nativeFontLayout
}

func (me *fontLayout) MeasureText(font Font, text string, width, height int32) (textWidth, textHeight int32) {
	return me.naiveLayout.MeasureText(font, text, width, height)
}

func (me *fontLayout) DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32) {
	me.naiveLayout.DrawText(canvas, font, paint, text, width, height)
}

func (me *fontLayout) CharacterIndexForPoint(font Font, text string, width, height int32, x, y float32) uint32 {
	return me.naiveLayout.CharacterIndexForPoint(font, text, width, height, x, y)
}

func (me *fontLayout) native() *nativeFontLayout {
	return me.naiveLayout
}
