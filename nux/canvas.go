// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

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
	// DrawText(text string, width, height float32, paint Paint)

	ResetClip()

	Flush()
	Destroy()

	native() *canvas
}
