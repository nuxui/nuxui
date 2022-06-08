// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

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
	// TextSize() float32 // TODO:: no use
	// SetTextSize(size float32)
	SetShadow(color Color, offsetX, offsetY, blur float32)
	Shadow() (color Color, offsetX, offsetY, blur float32)
	HasShadow() bool

	// TODO::
	// MeasureText(text string, width, height float32) (outWidth float32, outHeight float32)
	// CharacterIndexForPoint(text string, width, height float32, x, y float32) uint32
}

// TODO:: Attr
func NewPaint() Paint {
	// if attr == nil {
	// 	attr = Attr{}
	// }

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
