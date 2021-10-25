// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"time"

	"github.com/nuxui/nuxui/nux"
)

type Editor interface {
	nux.Widget
	nux.Size
	nux.Layout
	nux.Measure
	nux.Draw
	Visual

	Text() string
	SetText(text string)
}

func NewEditor(context nux.Context, attrs ...nux.Attr) Editor {
	attr := getAttr(attrs...)
	me := &editor{
		cursorPosition: 0,
		flicker:        false,
		text:           attr.GetString("text", ""),
		font:           nux.NewFont(attr),
	}

	ellipsize := attr.GetString("ellipsize", "none")
	switch ellipsize {
	case "none":
		// text.ellipsize = C.PANGO_ELLIPSIZE_NONE
	}

	me.WidgetBase = nux.NewWidgetBase(context, me, attrs...)
	me.WidgetSize = nux.NewWidgetSize(context, me, attrs...)
	me.WidgetVisual = NewWidgetVisual(context, me, attrs...)
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)

	return me
}

type editor struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	text        string
	editingText string
	editingLoc  int32
	font        *nux.Font
	ellipsize   int

	cursorPosition int
	flicker        bool
	flickerTimer   nux.Timer
	focus          bool
}

func (me *editor) OnCreate(content nux.Widget) {

	nux.OnTapDown(me, me.onTapDown)
	nux.OnPanStart(me, me.onPanStart)
	nux.OnPanUpdate(me, me.OnPanUpdate)
}

func (me *editor) onTapDown(detail nux.GestureDetail) {
	// TODO:: x, y --> font position
	me.cursorPosition = len([]rune(me.text))
	nux.RequestFocus(me)
}

func (me *editor) onPanStart(detail nux.GestureDetail) {
	me.cursorPosition = len([]rune(me.text))
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
		ms := me.MeasuredSize()
		// TODO:: cursor position
		nux.SetTextInputRect(float32(ms.Position.X), float32(ms.Position.Y)+me.font.Size, 0, 0)
		me.startTick()
	} else {
		nux.StopTextInput()
		me.stopTick()
	}
}

func (me *editor) onSizeChanged(widget nux.Widget) {
	nux.RequestLayout(me)
}

func (me *editor) onVisualChanged(widget nux.Widget) {
	nux.RequestRedraw(me)
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

func (me *editor) Layout(dx, dy, left, top, right, bottom int32) {
	ms := me.MeasuredSize()

	// setFrame
	ms.Position.Left = left
	ms.Position.Top = top
	ms.Position.Right = right
	ms.Position.Bottom = bottom
	ms.Position.X = dx
	ms.Position.Y = dy
}

func (me *editor) Measure(width, height int32) {
	if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {

		w := width
		h := height

		outW, outH := nux.MeasureText(me.text, me.font, w, h)
		// log.V("nuxui", "Editor MeasureText %s %d %d", me.text, outW, outH)
		ms := me.MeasuredSize()
		if nux.MeasureSpecMode(width) == nux.Auto {
			ms.Width = nux.MeasureSpec(outW+ms.Padding.Left+ms.Padding.Right, nux.Pixel)
		} else {
			ms.Width = width
		}

		if nux.MeasureSpecMode(height) == nux.Auto {
			ms.Height = nux.MeasureSpec(outH+ms.Padding.Top+ms.Padding.Bottom, nux.Pixel)
		} else {
			ms.Height = height
		}
	}
}

func (me *editor) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	runes := []rune(me.text)
	front := string(runes[:me.cursorPosition])
	end := string(runes[me.cursorPosition:len(runes)])
	all := fmt.Sprintf("%s%s%s", front, me.editingText, end)
	m := fmt.Sprintf("%s%s", front, string([]rune(me.editingText)[:me.editingLoc]))
	outW, outH := nux.MeasureText(m, me.font, 1000, 1000)

	// TODO has padding?
	ms := me.MeasuredSize()
	canvas.Save()
	canvas.Translate(ms.Padding.Left, ms.Padding.Top)
	canvas.ClipRect(0, 0,
		ms.Width-ms.Padding.Left-ms.Padding.Right,
		ms.Height-ms.Padding.Top-ms.Padding.Bottom)

	canvas.DrawText(all, me.font, int32(ms.Width), int32(ms.Height), &nux.Paint{me.font.Color, nux.FILL, 2})
	canvas.Translate(outW, 0)
	p := &nux.Paint{Color: me.font.Color, Style: nux.FILL, Width: 1}

	if me.flicker {
		canvas.DrawRect(0, 1, 1, outH-1, p)
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
		end := string(runes[me.cursorPosition:len(runes)])
		ret := fmt.Sprintf("%s%s%s", front, event.Text(), end)
		me.cursorPosition += len([]rune(event.Text()))
		me.SetText(ret)
		me.startTick()
	}
	// log.V("nuxui", "OnTypeEvent editingText=%s text=%s", me.editingText, me.text)
	return true
}

func (me *editor) OnKeyEvent(event nux.KeyEvent) bool {
	switch event.KeyCode() {
	case nux.Key_BackSpace:
		if event.Action() == nux.Action_Down {
			if len(me.text) == 0 {
				return false
			}

			me.cursorPosition--
			if me.cursorPosition < 0 {
				me.cursorPosition = 0
				return true
			}

			runes := []rune(me.text)
			front := string(runes[:me.cursorPosition])
			end := string(runes[me.cursorPosition+1 : len(runes)])
			ret := fmt.Sprintf("%s%s", front, end)
			me.SetText(ret)
			me.startTick()
			return true
		}
	case nux.Key_Delete:
		if event.Action() == nux.Action_Down {
			if len(me.text) == 0 {
				return false
			}

			runes := []rune(me.text)
			if me.cursorPosition+1 > len(runes) {
				me.cursorPosition = len(runes)
				return true
			}

			front := string(runes[:me.cursorPosition])
			end := string(runes[me.cursorPosition+1 : len(runes)])
			ret := fmt.Sprintf("%s%s", front, end)
			me.SetText(ret)
			me.startTick()
			return true
		}
	case nux.Key_Left:
		if event.Action() == nux.Action_Down && me.editingText == "" {
			me.cursorPosition--
			if me.cursorPosition < 0 {
				me.cursorPosition = 0
				return true
			}
			nux.RequestRedraw(me)
			me.startTick()
		}
	case nux.Key_Right:
		if event.Action() == nux.Action_Down && me.editingText == "" {
			me.cursorPosition++
			maxlen := len([]rune(me.text))
			if me.cursorPosition > maxlen {
				me.cursorPosition = maxlen
				return true
			}
			nux.RequestRedraw(me)
			me.startTick()
		}
	case nux.Key_Up:
	case nux.Key_Down:
	case nux.Key_V:
	}
	return true
}
