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
	SetDash(dash []float32)
	Dash() []float32

	SetShadow(color Color, offsetX, offsetY, blur float32)
	Shadow() (color Color, offsetX, offsetY, blur float32)
	HasShadow() bool

	native() *paint
}
