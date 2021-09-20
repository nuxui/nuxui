// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

// TODO:: float32 or double?
type Canvas interface {
	Save()
	Restore()

	Translate(x, y int32)
	TranslateF(x, y float32)
	Scale(x, y int32)
	ScaleF(x, y float32)
	Rotate(angle int32)
	RotateF(angle float32)
	// Transform()
	// SetMatrix()
	// GetMatrix()
	// Skew()
	ClipRect(left, top, right, bottom int32)
	ClipRectF(left, top, right, bottom float32)
	// ClipPath()

	///////////// draw
	DrawRect(left, top, right, bottom int32, paint *Paint)
	DrawRectF(left, top, right, bottom float32, paint *Paint)
	DrawArc(x, y, radius, angle1, angle2 float32, useCenter bool, paint *Paint)
	DrawOval(left, top, right, bottom int32, paint *Paint)
	DrawOvalF(left, top, right, bottom float32, paint *Paint)
	DrawRoundRect(left, top, right, bottom int32, radius int32, paint *Paint)
	DrawRoundRectF(left, top, right, bottom float32, radius float32, paint *Paint)
	// drawPaint(paint *Paint)
	DrawColor(color Color)
	DrawAlpha(alpha float32)
	SetColor(color Color)
	// GetTextRect(text string, fontFamily string, fontSize float32) C.cairo_text_extents_t
	DrawText(text string, font *Font, width, height int32, paint *Paint)

	MeasureText(text string, font *Font, width, height int32) (outWidth, outHeight int32)

	// DrawImage(src string, width, height int32)
	// DrawPNG(*cPNGImage)
	DrawImage(img Image)

	SetAntialias(a int)
	GetAntialias() int

	UserToDevice(x, y float32)
}
