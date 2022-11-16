// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package nux

import ()

type canvas struct {
}

func newCanvas() *canvas {
	return &canvas{}
}

func (me *canvas) native() *canvas {
	return me
}

func (me *canvas) ResetClip() {
}

func (me *canvas) Save() {
}

func (me *canvas) Restore() {
}

func (me *canvas) Translate(x, y float32) {
}

func (me *canvas) Scale(x, y float32) {
}

func (me *canvas) Rotate(angle float32) {
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
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {

}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}

func (me *canvas) SetAlpha(alpha float32) {
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {

}

func (me *canvas) DrawRoundRect(x, y, width, height, rLT, rRT, rRB, rLB float32, paint Paint) {

}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {

}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
}

func (me *canvas) DrawPath(path Path, paint Paint) {

}

func (me *canvas) DrawImage(img Image) {
}

func (me *canvas) Flush() {
}

func (me *canvas) Destroy() {
}
