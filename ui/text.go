// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

// TODO Text can automatically fine-tune the spacing to ensure that the font occupies the entire line. Basic Text does not do this and uses the new AlignedText

import (
	"math"
	"time"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
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
	textColor          nux.Color
	textHighlightColor nux.Color
	paint              nux.Paint
	ellipsize          int //TODO
	font               nux.Font
	fontLayout         nux.FontLayout

	downTime time.Time
}

func NewText(attr nux.Attr) Text {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &text{
		text:               attr.GetString("text", ""),
		textColor:          attr.GetColor("textColor", nux.Black),
		textHighlightColor: attr.GetColor("textHighlightColor", nux.Transparent),
		paint:              nux.NewPaint(),
		font:               nux.NewFont(attr.GetAttr("font", nil)),
		fontLayout:         nux.NewFontLayout(),
	}

	me.WidgetBase = nux.NewWidgetBase(attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)

	return me
}

func (me *text) OnMount() {
	// TODO:: widget gesture penetrate
	if me.Background() != nil && me.Background().HasState() {
		nux.OnTapDown(me, me.onTapDown)
		nux.OnTapUp(me, me.onTapUp)
		nux.OnTapCancel(me, me.onTapUp)
	}
}

func (me *text) OnEject() {
}

func (me *text) onTapDown(detail nux.GestureDetail) {
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().AddState(nux.State_Pressed)
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().AddState(nux.State_Pressed)
			changed = true
		}
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *text) onTapUp(detail nux.GestureDetail) {
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().DelState(nux.State_Pressed)
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().DelState(nux.State_Pressed)
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

func (me *text) Measure(width, height nux.MeasureDimen) {
	frame := me.Frame()

	outW, outH := me.fontLayout.MeasureText(me.font, me.text, width.Value(), height.Value())
	// log.I("nuxui", "fontLayout.MeasureText outW=%d, outH=%d", outW, outH)

	hPPx, hPPt, vPPx, vPPt, paddingMeasuredFlag := measurePadding(width, height, me.Padding(), frame, float32(outH), 0)
	if hPPt >= 100.0 || vPPt >= 100.0 {
		log.Fatal("nuxui", "padding percent size should at 0% ~ 100%")
	}

	if width.Mode() == nux.Pixel {
		frame.Width = width.Value()
	} else {
		w := (float32(outW) + hPPx) / (1.0 - hPPt/100.0)
		frame.Width = int32(math.Ceil(float64(w)))
		width = nux.MeasureSpec(frame.Width, nux.Pixel)
	}

	if height.Mode() == nux.Pixel {
		frame.Height = height.Value()
	} else {
		h := (float32(outH) + vPPx) / (1.0 - vPPt/100.0)
		frame.Height = int32(math.Ceil(float64(h)))
		height = nux.MeasureSpec(frame.Height, nux.Pixel)
	}

	if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
		measurePadding(width, height, me.Padding(), frame, float32(outH), paddingMeasuredFlag)
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

func (me *text) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	frame := me.Frame()
	canvas.Save()
	canvas.Translate(float32(frame.X), float32(frame.Y))
	canvas.Translate(float32(frame.Padding.Left), float32(frame.Padding.Top))

	if me.text != "" {
		me.paint.SetColor(me.textColor)
		me.fontLayout.DrawText(canvas, me.font, me.paint, me.text, frame.Width, frame.Height)
	}
	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
