// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/nux/internal/android"
)

type canvas struct {
	ref android.Canvas
}

func newCanvas(ref android.Canvas) *canvas {
	return &canvas{
		ref: ref,
	}
}

func (me *canvas) native() *canvas {
	return me
}

func (me *canvas) ResetClip() {
}

func (me *canvas) Save() {
	me.ref.Save()
}

func (me *canvas) Restore() {
	me.ref.Restore()
}

func (me *canvas) Translate(x, y float32) {
	me.ref.Translate(x, y)
}

func (me *canvas) Scale(x, y float32) {
	me.ref.Scale(x, y)
}

// TODO:: radian, angle, degrees
func (me *canvas) Rotate(degrees float32) {
	me.ref.Rotate(degrees)
}

func (me *canvas) Skew(x, y float32) {
	// TODO::
}

func (me *canvas) Transform(a, b, c, d, e, f float32) {
	// TODO::
}

func (me *canvas) SetMatrix(matrix Matrix) {
	// TODO::
}

func (me *canvas) GetMatrix() Matrix {
	// TODO::
	return Matrix{}
}

func (me *canvas) ClipRect(x, y, width, height float32) {
	me.ref.ClipRect(x, y, width, height)
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	me.ref.ClipRoundRect(x, y, width, height, cornerX, cornerY)
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}

func (me *canvas) SetAlpha(alpha float32) {
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	me.ref.DrawRect(x, y, width, height, paint.native().ref)
}

func (me *canvas) DrawRoundRect(x, y, width, height, rLT, rRT, rRB, rLB float32, paint Paint) {
	if width < 0 || height < 0 {
		return
	}
	me.ref.DrawRoundRect(x, y, width, height, rLT, rRT, rRB, rLB, paint.native().ref)
}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {

}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
}

func (me *canvas) DrawPath(path Path, paint Paint) {

}

func (me *canvas) DrawImage(img Image) {
	img.Draw(me)
}

func (me *canvas) Flush() {
}

func (me *canvas) Destroy() {
}
