// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

// TODO Text can automatically fine-tune the spacing to ensure that the font occupies the entire line. Basic Text does not do this and uses the new AlignedText

import (
	"math"
	"time"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Text interface {
	nux.Widget
	nux.Size
	Visual

	Text() string
	SetText(text string)
}

type text struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	text               string
	textSize           float32
	textColor          nux.Color
	textHighlightColor nux.Color
	paint              nux.Paint
	ellipsize          int

	downTime time.Time
}

func NewText(attrs ...nux.Attr) Text {
	attr := nux.MergeAttrs(attrs...)
	me := &text{
		text:               attr.GetString("text", ""),
		textSize:           attr.GetFloat32("textSize", 12),
		textColor:          attr.GetColor("textColor", nux.White),
		textHighlightColor: attr.GetColor("textHighlightColor", nux.Transparent),
		paint:              nux.NewPaint(),
		// paint:              nux.NewPaint(attr.GetAttr("font", nux.Attr{})),
		// ellipsize: ellipsizeFromName(attr.GetString("ellipsize", "none")),
	}

	me.WidgetBase = nux.NewWidgetBase(attrs...)
	me.WidgetSize = nux.NewWidgetSize(attrs...)
	me.WidgetVisual = NewWidgetVisual(me, attrs...)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)

	return me
}

func (me *text) Mount() {
	log.I("nuxui", "text Mount")
	nux.OnTapDown(me, me.onTapDown)
	nux.OnTapUp(me, me.onTapUp)
	nux.OnTapCancel(me, me.onTapUp)
}

func (me *text) Eject() {
	log.I("nuxui", "text Eject")
}

func (me *text) onTapDown(detail nux.GestureDetail) {
	log.V("nuxui", "text onTapDown")
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().SetState(nux.Attr{"state": "pressed"})
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().SetState(nux.Attr{"state": "pressed"})
			changed = true
		}
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *text) onTapUp(detail nux.GestureDetail) {
	log.V("nuxui", "text onTapUp")
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().SetState(nux.Attr{"state": "normal"})
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().SetState(nux.Attr{"state": "normal"})
			changed = true
		}
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *text) onSizeChanged() {
	nux.RequestLayout(me)
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
func (me *text) Layout(x, y, width, height int32) {
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}
}

func (me *text) Measure(width, height int32) {
	// measuredDuration := log.Time()
	// defer log.TimeEnd(measuredDuration, "nuxui", "ui.Text Measure ")

	var vPPt float32 // horizontal padding percent
	var vPPx float32 // horizontal padding pixel
	var hPPt float32
	var hPPx float32

	frame := me.Frame()

	// 1. Calculate its own padding size
	if me.Padding() != nil {
		switch me.Padding().Left.Mode() {
		case nux.Pixel:
			l := me.Padding().Left.Value()
			frame.Padding.Left = util.Roundi32(l)
			hPPx += l
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				l := me.Padding().Left.Value() / 100 * float32(nux.MeasureSpecValue(width))
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
			case nux.Auto:
				hPPt += me.Padding().Left.Value()
			}
		}

		switch me.Padding().Right.Mode() {
		case nux.Pixel:
			r := me.Padding().Right.Value()
			frame.Padding.Right = util.Roundi32(r)
			hPPx += r
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				r := me.Padding().Right.Value() / 100 * float32(nux.MeasureSpecValue(width))
				frame.Padding.Right = util.Roundi32(r)
				hPPx += r
			case nux.Auto:
				hPPt += me.Padding().Right.Value()
			}
		}

		switch me.Padding().Top.Mode() {
		case nux.Pixel:
			t := me.Padding().Top.Value()
			frame.Padding.Top = util.Roundi32(t)
			vPPx += t
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				t := me.Padding().Top.Value() / 100 * float32(nux.MeasureSpecValue(height))
				frame.Padding.Top = util.Roundi32(t)
				vPPx += t
			case nux.Auto:
				vPPt += me.Padding().Top.Value()
			}
		}

		switch me.Padding().Bottom.Mode() {
		case nux.Pixel:
			b := me.Padding().Bottom.Value()
			frame.Padding.Bottom = util.Roundi32(b)
			vPPx += b
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				b := me.Padding().Bottom.Value() / 100 * float32(nux.MeasureSpecValue(height))
				frame.Padding.Bottom = util.Roundi32(b)
				vPPx += b
			case nux.Auto:
				vPPt += me.Padding().Bottom.Value()
			}
		}
	}

	if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {
		w := int32(width)
		h := int32(height)

		me.paint.SetTextSize(me.textSize)
		outW, outH := me.paint.MeasureText(me.text, float32(nux.MeasureSpecValue(w)), float32(nux.MeasureSpecValue(h)))

		frame := me.Frame()
		if nux.MeasureSpecMode(width) == nux.Auto {
			w := (float32(outW) + hPPx) / (1.0 - hPPt/100.0)
			frame.Width = nux.MeasureSpec(int32(math.Ceil(float64(w))), nux.Pixel)
		} else {
			frame.Width = width
		}

		if nux.MeasureSpecMode(height) == nux.Auto {
			h := (float32(outH) + vPPx) / (1.0 - vPPt/100.0)
			frame.Height = nux.MeasureSpec(int32(math.Ceil(float64(h))), nux.Pixel)
		} else {
			frame.Height = height
		}
	}
}

func (me *text) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		log.V("nuxui", "text color drawable")
		me.Background().Draw(canvas)
	}

	frame := me.Frame()
	canvas.Save()
	canvas.Translate(float32(frame.X), float32(frame.Y))
	canvas.Translate(float32(frame.Padding.Left), float32(frame.Padding.Top))

	if me.text != "" {
		me.paint.SetTextSize(me.textSize)
		me.paint.SetColor(me.textColor)
		canvas.DrawText(me.text, float32(frame.Width), float32(frame.Height), me.paint)
	}
	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}

}
