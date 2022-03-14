// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Check interface {
	nux.Widget
	nux.Size
	// nux.Layout
	// nux.Measure
	// nux.Draw
	Visual
}

func NewCheck(attrs ...nux.Attr) Check {
	attr := nux.MergeAttrs(attrs...)
	me := &check{
		paint: nux.NewPaint(),
	}
	me.WidgetBase = nux.NewWidgetBase(attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

type check struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual
	paint nux.Paint
}

func (me *check) Mount() {
}

func (me *check) Measure(width, height int32) {
	frame := me.Frame()

	if p := me.Padding(); p != nil {
		switch p.Left.Mode() {
		case nux.Pixel:
			frame.Padding.Left = util.Roundi32(p.Left.Value())
		}

		switch p.Top.Mode() {
		case nux.Pixel:
			frame.Padding.Top = util.Roundi32(p.Top.Value())
		}

		switch p.Right.Mode() {
		case nux.Pixel:
			frame.Padding.Right = util.Roundi32(p.Right.Value())
		}

		switch p.Bottom.Mode() {
		case nux.Pixel:
			frame.Padding.Bottom = util.Roundi32(p.Bottom.Value())
		}
	}

	// if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {
	// 	var dw, dh int32
	// 	if me.srcDrawable != nil {
	// 		dw, dh = me.srcDrawable.Size()
	// 	}
	// 	if nux.MeasureSpecMode(width) == nux.Auto {
	// 		frame.Width = dw
	// 	} else {
	// 		frame.Width = nux.MeasureSpecValue(width)
	// 	}
	// 	if nux.MeasureSpecMode(height) == nux.Auto {
	// 		frame.Height = dh
	// 	} else {
	// 		frame.Height = nux.MeasureSpecValue(height)
	// 	}
	// }

	// frame.Width += frame.Padding.Left + frame.Padding.Right
	// frame.Height += frame.Padding.Top + frame.Padding.Bottom
}

func (me *check) Layout(x, y, width, height int32) {
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}

	// frame := me.Frame()

}

func (me *check) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *check) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	me.paint.SetColor(0xffff0000)
	frame := me.Frame()
	canvas.DrawRoundRect(
		float32(frame.Padding.Left),
		float32(frame.Padding.Top),
		float32(frame.Width-frame.Padding.Left-frame.Padding.Right),
		float32(frame.Height-frame.Padding.Top-frame.Padding.Bottom),
		3,
		me.paint)

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
