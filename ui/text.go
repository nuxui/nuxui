// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

// TODO Text can automatically fine-tune the spacing to ensure that the font occupies the entire line. Basic Text does not do this and uses the new AlignedText

import (
	"time"

	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Text interface {
	nux.Widget
	nux.Size
	nux.Creating
	nux.Layout
	nux.Measure
	nux.Draw
	Visual

	Text() string
	SetText(text string)
}

type text struct {
	// TODO use TextMeasure combination, me.MeasureText()
	// TextMeasure
	nux.WidgetBase
	nux.WidgetSize
	WidgetVisual

	text      string
	Font      nux.Font
	ellipsize int

	downTime time.Time
}

func NewText() Text {
	me := &text{}
	me.WidgetSize.Owner = me
	me.WidgetVisual.Owner = me
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)

	me.Creating(nux.Attr{})
	return me
}

func (me *text) Creating(attr nux.Attr) {
	if attr == nil {
		attr = nux.Attr{}
	}

	me.WidgetBase.Creating(attr)
	me.WidgetSize.Creating(attr)
	me.WidgetVisual.Creating(attr)

	me.Font.Creating(attr.GetAttr("font", nux.Attr{}))

	me.text = attr.GetString("text", "")

	ellipsize := attr.GetString("ellipsize", "none")
	switch ellipsize {
	case "none":
		// text.ellipsize = C.PANGO_ELLIPSIZE_NONE

	}
}

func (me *text) Created(content nux.Widget) {
	nux.OnTapDown(me, me.onTapDown)
	nux.OnTapUp(me, me.onTapUp)
	nux.OnTapCancel(me, me.onTapUp)
}

func (me *text) onTapDown(widget nux.Widget) {
	me.SetBackgroundColor(0xFF938276)
	me.downTime = time.Now()
	nux.NewTimerBackToUI(nux.GESTURE_DOWN2UP_DELAY*time.Millisecond, func() {

	})
}

func (me *text) onTapUp(widget nux.Widget) {
	if sub := time.Now().Sub(me.downTime); sub < nux.GESTURE_DOWN2UP_DELAY*time.Millisecond {
		nux.NewTimerBackToUI(nux.GESTURE_DOWN2UP_DELAY*time.Millisecond-sub, func() {
			me.doTapUp(widget)
		})
	} else {
		me.doTapUp(widget)
	}
}

func (me *text) doTapUp(widget nux.Widget) {
	me.SetBackgroundColor(0xFFFFFFFF)
}

func (me *text) onSizeChanged(widget nux.Widget) {

}
func (me *text) onVisualChanged(widget nux.Widget) {
	nux.RequestRedraw(me)
}

func (me *text) Text() string {
	return me.text
}

func (me *text) SetText(text string) {
	if me.text != text {
		me.text = text
		nux.RequestLayout(me)
	}
}

// Responsible for determining the position of the widget align, margin...
func (me *text) Layout(dx, dy, left, top, right, bottom int32) {
	// log.V("nuxui", "text layout %d, %d, %d, %d, %d, %d", dx, dy, left, top, right, bottom)
}

func (me *text) Measure(width, height int32) {
	// measuredDuration := log.Time()
	// defer log.TimeEnd("nuxui", "ui.Text Measure", measuredDuration)

	var vPPt float32 // horizontal padding percent
	var vPPx float32 // horizontal padding pixel
	var hPPt float32
	var hPPx float32

	ms := me.MeasuredSize()

	// 1. Calculate its own padding size
	if me.HasPadding() {
		switch me.PaddingLeft().Mode() {
		case nux.Pixel:
			l := me.PaddingLeft().Value()
			ms.Padding.Left = util.Roundi32(l)
			hPPx += l
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				l := me.PaddingLeft().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Left = util.Roundi32(l)
				hPPx += l
			case nux.Auto:
				hPPt += me.PaddingLeft().Value()
			}
		}

		switch me.PaddingRight().Mode() {
		case nux.Pixel:
			r := me.PaddingRight().Value()
			ms.Padding.Right = util.Roundi32(r)
			hPPx += r
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				r := me.PaddingRight().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Right = util.Roundi32(r)
				hPPx += r
			case nux.Auto:
				hPPt += me.PaddingRight().Value()
			}
		}

		switch me.PaddingTop().Mode() {
		case nux.Pixel:
			t := me.PaddingTop().Value()
			ms.Padding.Top = util.Roundi32(t)
			vPPx += t
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				t := me.PaddingTop().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Top = util.Roundi32(t)
				vPPx += t
			case nux.Auto:
				vPPt += me.PaddingTop().Value()
			}
		}

		switch me.PaddingBottom().Mode() {
		case nux.Pixel:
			b := me.PaddingBottom().Value()
			ms.Padding.Bottom = util.Roundi32(b)
			vPPx += b
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				b := me.PaddingBottom().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Bottom = util.Roundi32(b)
				vPPx += b
			case nux.Auto:
				vPPt += me.PaddingBottom().Value()
			}
		}
	}

	if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {
		w := int32(width)
		h := int32(height)

		outW, outH := nux.MeasureText(me.text, &me.Font, w, h)
		// log.V("nuxui", "Text '%s' Measure: %d, %d\n", me.text, outW, outH)
		ms := me.MeasuredSize()
		if nux.MeasureSpecMode(width) == nux.Auto {
			w := (float32(outW) + hPPx) / (1.0 - hPPt/100.0)
			ms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
		} else {
			ms.Width = width
		}

		if nux.MeasureSpecMode(height) == nux.Auto {
			h := (float32(outH) + vPPx) / (1.0 - vPPt/100.0)
			ms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
		} else {
			ms.Height = height
		}
	}
}

func (me *text) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	ms := me.MeasuredSize()
	canvas.Save()
	canvas.Translate(ms.Padding.Left, ms.Padding.Top)
	canvas.ClipRect(0, 0,
		ms.Width-ms.Padding.Left-ms.Padding.Right,
		ms.Height-ms.Padding.Top-ms.Padding.Bottom)

	if me.text != "" {
		canvas.DrawText(me.text, &me.Font, int32(ms.Width), int32(ms.Height), &nux.Paint{me.Font.Color, nux.FILL, 2})
	}

	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
