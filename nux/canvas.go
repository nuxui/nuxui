// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "nuxui.org/nuxui/log"

type Matrix struct {
	A float32
	B float32
	C float32
	D float32
	E float32
	F float32
}

type Canvas interface {
	Save()
	Restore()

	Translate(x, y float32)
	Scale(x, y float32)
	Rotate(angle float32)
	Skew(x, y float32)                  // https://developer.android.com/reference/android/graphics/Canvas#skew(float,%20float)
	Transform(a, b, c, d, e, f float32) //
	SetMatrix(matrix Matrix)            // https://cairographics.org/manual/cairo-Transformations.html#cairo-set-matrix
	GetMatrix() Matrix                  // https://cairographics.org/manual/cairo-Transformations.html#cairo-get-matrix

	ClipRect(x, y, width, height float32)
	ClipRoundRect(x, y, width, height, cornerX, cornerY float32)
	ClipPath(path Path)

	///////////// draw
	SetAlpha(alpha float32)
	DrawRect(x, y, width, height float32, paint Paint)
	DrawRoundRect(x, y, width, height float32, rLT, rRT, rRB, rLB float32, paint Paint)
	DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint)
	DrawOval(x, y, width, height float32, paint Paint)
	DrawPath(path Path)
	// DrawColor(color Color)
	DrawImage(img Image)
	DrawText(text string, width, height float32, paint Paint)

	ResetClip()

	Flush()
	Destroy()
}

type PaintStyle int

const (
	PaintStyle_Fill   PaintStyle = 1
	PaintStyle_Stroke PaintStyle = 2
	PaintStyle_Both   PaintStyle = 1 | 2
)

type Paint interface {
	Color() Color
	SetColor(color Color)
	AntiAlias() bool
	SetAntiAlias(antialias bool)
	Width() float32
	SetWidth(width float32)
	Style() PaintStyle
	SetStyle(style PaintStyle)
	TextSize() float32
	SetTextSize(size float32)
	SetShadow(color Color, offsetX, offsetY, blur float32)
	Shadow() (color Color, offsetX, offsetY, blur float32)
	HasShadow() bool
	MeasureText(text string, width, height float32) (outWidth float32, outHeight float32)
	CharacterIndexForPoint(text string, width, height float32, x, y float32) uint32
}

// func NewPaint(attrs ...Attr) Paint {
// 	attr := Attr{}
// 	attr.Merge(attrs...)

// 	me := &paint{
// 		style:      PaintStyle_Fill,
// 		textSize:   attr.GetFloat32("textSize", 14),
// 		color:      attr.GetColor("color", 0xff000000),
// 		antialias:  attr.GetBool("antialias", false),
// 		fontFamily: attr.GetString("family", ""),
// 		fontWeight: FontWeightFromName(attr.GetString("weight", "normal")),
// 		fontItalic: attr.GetString("style", "normal") == "italic",
// 	}
// 	return me
// }

func NewPaint() Paint {

	me := &paint{
		width:      1.0,
		style:      PaintStyle_Fill,
		textSize:   16,
		color:      0xff000000,
		antialias:  false,
		fontFamily: "",
		fontWeight: FontWeight_Normal,
		fontItalic: false,
	}
	return me
}

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
		log.E("nuxui", "unknow font weight name %s", name)
		return FontWeight_Normal
	}
}

type paint struct {
	color       Color
	style       PaintStyle
	width       float32
	textSize    float32
	antialias   bool
	fontFamily  string
	fontWeight  FontWeight
	fontItalic  bool
	shadowColor Color
	shadowX     float32
	shadowY     float32
	shadowBlur  float32
}

func (me *paint) Color() Color {
	return me.color
}

func (me *paint) SetColor(color Color) {
	me.color = color
}

func (me *paint) AntiAlias() bool {
	return me.antialias
}

func (me *paint) SetAntiAlias(antialias bool) {
	me.antialias = antialias
}

func (me *paint) Width() float32 {
	return me.width
}

func (me *paint) SetWidth(width float32) {
	me.width = width
}

func (me *paint) Style() PaintStyle {
	return me.style
}

func (me *paint) SetStyle(style PaintStyle) {
	me.style = style
}

func (me *paint) TextSize() float32 {
	return me.textSize
}

func (me *paint) SetTextSize(size float32) {
	me.textSize = size
}

func (me *paint) SetShadow(color Color, offsetX, offsetY, blur float32) {
	me.shadowColor = color
	me.shadowX = offsetX
	me.shadowY = offsetY
	me.shadowBlur = blur
}

func (me *paint) Shadow() (color Color, offsetX, offsetY, blur float32) {
	return me.shadowColor, me.shadowX, me.shadowY, me.shadowBlur
}

func (me *paint) HasShadow() bool {
	return me.shadowColor != 0 && me.shadowBlur > 0
}

// func (me *paint) MeasureText(text string, start, end int32) (width float32, height float32) {
// 	return
// }
