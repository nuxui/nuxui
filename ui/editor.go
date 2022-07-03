// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"math"
	"time"
	"unicode/utf8"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

type Editor interface {
	nux.Widget
	nux.Size
	Visual

	Text() string
	SetText(text string)
}

func NewEditor(attr nux.Attr) Editor {
	if attr == nil {
		attr = nux.Attr{}
	}

	if !attr.Has("theme") {
		attr = nux.MergeAttrs(editor_theme(nux.ThemeLight), attr)
	}

	me := &editor{
		cursorPosition:     0,
		flicker:            false,
		text:               attr.GetString("text", ""),
		textColor:          attr.GetColor("textColor", nux.White),
		textHighlightColor: attr.GetColor("textHighlightColor", nux.Transparent),
		font:               nux.NewFont(attr.GetAttr("font", nil)),
		fontLayout:         nux.NewFontLayout(),
		paint:              nux.NewPaint(),
		// paint:              nux.NewPaint(attr.GetAttr("font", nux.Attr{})),
		paintFlicker: nux.NewPaint(),
		// ellipsize: ellipsizeFromName(attr.GetString("ellipsize", "none")),
	}

	me.WidgetBase = nux.NewWidgetBase(attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)

	return me
}

type editor struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	font               nux.Font
	fontLayout         nux.FontLayout
	text               string
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

func (me *editor) OnMount() {
	// log.I("nuxui", "editor mount")
	nux.OnTapDown(me.Info().Self, me.onTapDown)
	nux.OnTapUp(me.Info().Self, me.onTapUp)
	nux.OnTapCancel(me.Info().Self, me.onTapCancel)
	nux.OnPanDown(me.Info().Self, me.onPanStart)

	nux.OnHoverEnter(me.Info().Self, func(detail nux.GestureDetail) {
		nux.LoadNativeCursor(nux.CursorIBeam).Set()
	})
	nux.OnHoverExit(me.Info().Self, func(detail nux.GestureDetail) {
		nux.LoadNativeCursor(nux.CursorArrow).Set()
	})
}

func (me *editor) onTapDown(detail nux.GestureDetail) {
	// log.I("nuxui", "editor onTapDown")
	frame := me.Frame()
	me.cursorPosition = me.fontLayout.CharacterIndexForPoint(me.font, me.text, frame.Width, frame.Height, detail.X(), detail.Y())
	nux.RequestFocus(me)
}

func (me *editor) onTapUp(detail nux.GestureDetail) {
	// log.I("nuxui", "editor onTapUp")
}

func (me *editor) onTapCancel(detail nux.GestureDetail) {
	// log.I("nuxui", "editor onTapCancel")
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
		nux.SetTextInputRect(float32(frame.X), float32(frame.Y+me.font.Size()), 0, 0)
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

func (me *editor) Measure(width, height nux.MeasureDimen) {
	frame := me.Frame()

	outW, outH := me.fontLayout.MeasureText(me.font, me.text, width.Value(), height.Value())

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
	// log.I("nuxui", "all: %s", all)

	frame := me.Frame()
	canvas.Save()
	canvas.Translate(float32(frame.X), float32(frame.Y))
	canvas.Translate(float32(frame.Padding.Left), float32(frame.Padding.Top))

	me.paint.SetColor(me.textColor)
	// me.fontLayout.DrawText(canvas, me.font, me.paint, all, frame.Width-frame.Padding.Left-frame.Padding.Right, frame.Height-frame.Padding.Top-frame.Padding.Bottom)
	me.fontLayout.DrawText(canvas, me.font, me.paint, all, 10000, 100000)

	if me.flicker {
		m := fmt.Sprintf("%s%s", front, string([]rune(me.editingText)[:me.editingLoc]))
		outW, outH := me.fontLayout.MeasureText(me.font, m, 1000, 1000) // TODO:: 1000
		canvas.Translate(float32(outW), 0)
		me.paintFlicker.SetColor(0x000000cc)
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

	if !me.flicker {
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

func (me *editor) OnTypingEvent(event nux.TypingEvent) bool {
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
	return true
}

func (me *editor) OnKeyEvent(event nux.KeyEvent) bool {
	// log.I("nuxui", "editor OnKeyEvent %s", event)

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
		// log.I("nuxui", "editor Key_Left")
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
