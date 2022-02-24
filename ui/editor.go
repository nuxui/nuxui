// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"math"
	"time"
	"unicode/utf8"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

type Editor interface {
	nux.Widget
	nux.Size
	Visual

	Text() string
	SetText(text string)
}

func NewEditor(attrs ...nux.Attr) Editor {
	attr := nux.MergeAttrs(attrs...)
	me := &editor{
		cursorPosition:     0,
		flicker:            false,
		text:               attr.GetString("text", ""),
		textSize:           attr.GetFloat32("textSize", 12),
		textColor:          attr.GetColor("textColor", nux.White),
		textHighlightColor: attr.GetColor("textHighlightColor", nux.Transparent),
		paint:              nux.NewPaint(),
		// paint:              nux.NewPaint(attr.GetAttr("font", nux.Attr{})),
		paintFlicker: nux.NewPaint(),
		// ellipsize: ellipsizeFromName(attr.GetString("ellipsize", "none")),
	}

	me.WidgetBase = nux.NewWidgetBase(attrs...)
	me.WidgetSize = nux.NewWidgetSize(attrs...)
	me.WidgetVisual = NewWidgetVisual(me, attrs...)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)

	nux.OnTapDown(me, me.onTapDown)
	nux.OnPanDown(me, me.onPanStart)
	nux.OnPanUpdate(me, me.OnPanUpdate)

	return me
}

type editor struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	text               string
	textSize           float32
	textColor          nux.Color
	textHighlightColor nux.Color
	editingText        string
	editingLoc         int32
	ellipsize          int
	paint              nux.Paint
	paintFlicker       nux.Paint

	cursorPosition uint32
	flicker        bool
	flickerTimer   nux.Timer
	focus          bool
}

func (me *editor) Mount() {
	log.D("nuxui", "====== Editor Mount")
	nux.OnTapDown(me, me.onTapDown)
	nux.OnPanDown(me, me.onPanStart)
	nux.OnPanUpdate(me, me.OnPanUpdate)
}

func (me *editor) onTapDown(detail nux.GestureDetail) {
	frame := me.Frame()
	me.cursorPosition = me.paint.CharacterIndexForPoint(me.text, float32(frame.Width), float32(frame.Height), detail.X(), detail.Y())
	nux.RequestFocus(me)
}

func (me *editor) onPanStart(detail nux.GestureDetail) {
	me.cursorPosition = uint32(strlen(me.text))
	nux.RequestFocus(me)
}

func (me *editor) OnPanUpdate(detail nux.GestureDetail) {

}

func (me *editor) Focusable() bool {
	return true
}

func (me *editor) HasFocus() bool {
	return me.focus
}

func (me *editor) FocusChanged(focus bool) {
	me.focus = focus
	if focus {
		nux.StartTextInput()
		frame := me.Frame()
		// TODO:: cursor position
		nux.SetTextInputRect(float32(frame.X), float32(frame.Y)+me.paint.TextSize(), 0, 0)
		me.startTick()
	} else {
		nux.StopTextInput()
		me.stopTick()
	}
}

func (me *editor) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *editor) Text() string {
	return me.text
}

func (me *editor) SetText(text string) {
	if me.text != text {
		me.text = text
		nux.RequestLayout(me)
		nux.RequestRedraw(me)
	}
}

func (me *editor) Measure(width, height int32) {
	if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {

		w := width
		h := height

		outW, outH := me.paint.MeasureText(me.text, float32(nux.MeasureSpecValue(w)), float32(nux.MeasureSpecValue(h)))
		// log.V("nuxui", "Editor MeasureText %s %d %d", me.text, outW, outH)
		frame := me.Frame()
		if nux.MeasureSpecMode(width) == nux.Auto {
			frame.Width = nux.MeasureSpec(int32(math.Ceil(float64(outW)))+frame.Padding.Left+frame.Padding.Right, nux.Pixel)
		} else {
			frame.Width = width
		}

		if nux.MeasureSpecMode(height) == nux.Auto {
			frame.Height = nux.MeasureSpec(int32(math.Ceil(float64(outH)))+frame.Padding.Top+frame.Padding.Bottom, nux.Pixel)
		} else {
			frame.Height = height
		}
	}
}

func (me *editor) Layout(x, y, width, height int32) {
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}
}

func (me *editor) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	runes := []rune(me.text)
	front := string(runes[:me.cursorPosition])
	end := string(runes[me.cursorPosition:])
	all := fmt.Sprintf("%s%s%s", front, me.editingText, end)

	frame := me.Frame()
	canvas.Save()
	canvas.Translate(float32(frame.X), float32(frame.Y))
	canvas.Translate(float32(frame.Padding.Left), float32(frame.Padding.Top))

	me.paint.SetColor(me.textColor)
	canvas.DrawText(all, float32(frame.Width-frame.Padding.Left-frame.Padding.Right), float32(frame.Height-frame.Padding.Top-frame.Padding.Bottom), me.paint)

	// draw cursor
	m := fmt.Sprintf("%s%s", front, string([]rune(me.editingText)[:me.editingLoc]))
	outW, outH := me.paint.MeasureText(m, 1000, 1000) // TODO:: 1000
	canvas.Translate(float32(outW), 0)
	if me.flicker {
		me.paintFlicker.SetColor(0xcc000000)
		canvas.DrawRect(0, 1, 1, float32(outH-1), me.paintFlicker)
	}

	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}

}

func (me *editor) startTick() {
	if me.flickerTimer != nil {
		me.flickerTimer.Cancel()
		me.flickerTimer = nil
	}

	if me.flicker != true {
		me.flicker = true // make first is show
		nux.RequestRedraw(me)
	}

	me.flickerTimer = nux.NewInterval(500*time.Millisecond, func() {
		me.flicker = !me.flicker
		nux.RequestRedraw(me)
	})
}

func (me *editor) stopTick() {
	if me.flickerTimer != nil {
		me.flickerTimer.Cancel()
		me.flickerTimer = nil
	}

	me.flicker = false
	nux.RequestRedraw(me)
}

func (me *editor) OnTypeEvent(event nux.TypeEvent) bool {
	log.V("nuxui", "Editor OnTypeEvent == %s", event.Text())

	switch event.Action() {
	case nux.Action_Preedit:
		me.editingText = event.Text()
		me.editingLoc = event.Location()
		nux.RequestLayout(me)
		nux.RequestRedraw(me)
		me.startTick()
	case nux.Action_Input:
		me.editingText = ""
		me.editingLoc = 0
		runes := []rune(me.text)
		front := string(runes[:me.cursorPosition])
		end := string(runes[me.cursorPosition:])
		ret := fmt.Sprintf("%s%s%s", front, event.Text(), end)
		me.cursorPosition += uint32(strlen(event.Text()))
		me.SetText(ret)
		me.startTick()
	}
	// log.V("nuxui", "OnTypeEvent editingText=%s text=%s", me.editingText, me.text)
	return true
}

func (me *editor) OnKeyEvent(event nux.KeyEvent) bool {
	log.V("nuxui", "Editor OnKeyEvent")
	switch event.KeyCode() {
	case nux.Key_BackSpace:
		if event.Action() == nux.Action_Down {
			if utf8.RuneCountInString(me.text) == 0 {
				return true
			}

			if me.cursorPosition == 0 {
				return true
			}
			me.cursorPosition--

			runes := []rune(me.text)
			front := string(runes[:me.cursorPosition])
			end := string(runes[me.cursorPosition+1:])
			ret := fmt.Sprintf("%s%s", front, end)
			me.SetText(ret)
			me.startTick()
			return true
		}
	case nux.Key_Delete:
		if event.Action() == nux.Action_Down {
			if strlen(me.text) == 0 {
				return true
			}

			runes := []rune(me.text)
			if me.cursorPosition+1 > uint32(len(runes)) {
				me.cursorPosition = uint32(len(runes))
				return true
			}

			front := string(runes[:me.cursorPosition])
			end := string(runes[me.cursorPosition+1:])
			ret := fmt.Sprintf("%s%s", front, end)
			me.SetText(ret)
			me.startTick()
			return true
		}
	case nux.Key_Left:
		if event.Action() == nux.Action_Down && me.editingText == "" {
			if me.cursorPosition == 0 {
				return true
			}
			me.cursorPosition--
			nux.RequestRedraw(me)
			me.startTick()
		}
	case nux.Key_Right:
		if event.Action() == nux.Action_Down && me.editingText == "" {
			if me.cursorPosition == uint32(strlen(me.text)) {
				return true
			}
			me.cursorPosition++
			nux.RequestRedraw(me)
			me.startTick()
		}
	case nux.Key_Tab:
		if event.Action() == nux.Action_Down {
			runes := []rune(me.text)
			front := string(runes[:me.cursorPosition])
			end := string(runes[me.cursorPosition:])
			ret := fmt.Sprintf("%s\t%s", front, end)
			me.cursorPosition++
			me.SetText(ret)
			me.startTick()
			return true
		}
	case nux.Key_Up:
	case nux.Key_Down:

	}
	return true
}

func strlen(text string) int {
	return len([]rune(text))
}
