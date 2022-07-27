// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Path interface {
	Rect(x, y, width, height float32)
	RoundRect(x, y, width, height, rLT, rRT, rRB, rLB float32)
	Ellipse(cx, cy, rx, ry float32)
	MoveTo(x, y float32)
	LineTo(x, y float32)
	CurveTo(x1, y1, x2, y2, x3, y3 float32)
	CurveToV(x2, y2, x3, y3 float32)
	Close()

	native() *path
}

type Gradient interface {
}

func NewPath() Path {
	return newPath()
}
