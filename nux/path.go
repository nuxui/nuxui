// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Path interface {
	AddRoundRect(left, top, right, bottom, rx, ry float32)
	// AddArc(left, top, right, bottom, startAngle, sweepAngle float32)
	// MoveTo(x, y float32)
	// LineTo(x, y float32)
	// Close()
}

type Gradient interface {
}
