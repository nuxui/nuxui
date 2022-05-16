// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/win32"
)

type window struct {
	hwnd        uintptr
	decor       Parent
	delegate    WindowDelegate
	focusWidget Widget

	initEvent PointerEvent
	timer     Timer

	preHdc      uintptr
	hdcBuffer   uintptr
	hBitMap     uintptr
	canvas      *canvas
	paintStruct win32.PAINTSTRUCT
}

func newWindow(attr Attr) *window {
	me := &window{
		hwnd:      0,
		preHdc:    0,
		hdcBuffer: 0,
	}

	me.CreateDecor(attr)
	return me
}

func (me *window) CreateDecor(attr Attr) Widget {
	creator := FindTypeCreator("nuxui.org/nuxui/ui.Layer")
	w := creator(attr)
	if p, ok := w.(Parent); ok {
		me.decor = p
	} else {
		log.Fatal("nuxui", "decor must is a Parent")
	}

	return me.decor
}

var TestDraw func(Canvas)

func (me *window) Draw(canvas Canvas) {
	if me.decor != nil {
		if f, ok := me.decor.(Draw); ok {
			var rectClient win32.RECT
			win32.GetClientRect(me.hwnd, &rectClient)

			win32.PatBlt(me.hdcBuffer, 0, 0, rectClient.Right-rectClient.Left, rectClient.Bottom-rectClient.Top, win32.WHITENESS)

			canvas.Save()
			canvas.ClipRect(0, 0, float32(rectClient.Right-rectClient.Left), float32(rectClient.Bottom-rectClient.Top))
			if TestDraw != nil {
				TestDraw(canvas)
			} else {
				f.Draw(canvas)
			}
			canvas.Restore()
			canvas.Flush()

			win32.BitBlt(me.preHdc, 0, 0, rectClient.Right-rectClient.Left, rectClient.Bottom-rectClient.Top, me.hdcBuffer, 0, 0, win32.SRCCOPY)
		}
	}
}

func (me *window) ID() uint64 {
	return 0
}

func (me *window) Size() (width, height int32) {
	var rect win32.RECT
	if err := win32.GetWindowRect(me.hwnd, &rect); err == nil {
		return rect.Right - rect.Left, rect.Bottom - rect.Top
	}
	return 0, 0
}

func (me *window) ContentSize() (width, height int32) {
	var rect win32.RECT
	if err := win32.GetClientRect(me.hwnd, &rect); err == nil {
		return rect.Right - rect.Left, rect.Bottom - rect.Top
	}
	return 0, 0
}

func (me *window) LockCanvas() Canvas {
	hdc, _ := win32.BeginPaint(me.hwnd, &me.paintStruct)
	// most time preHdc == hdc
	if me.preHdc != hdc {
		if me.preHdc != 0 {
			me.canvas.Destroy()
			win32.DeleteObject(me.hBitMap)
			win32.DeleteDC(me.preHdc)
			win32.DeleteDC(me.hdcBuffer)
		}
		me.preHdc = hdc
		var rectClient win32.RECT
		win32.GetClientRect(me.hwnd, &rectClient)
		me.hdcBuffer, _ = win32.CreateCompatibleDC(hdc)
		me.hBitMap, _ = win32.CreateCompatibleBitmap(hdc, rectClient.Right-rectClient.Left, rectClient.Bottom-rectClient.Top)
		win32.SelectObject(me.hdcBuffer, me.hBitMap)

		me.canvas = newCanvas(me.hdcBuffer)
	}
	return me.canvas
}

func (me *window) UnlockCanvas(c Canvas) {
	win32.EndPaint(me.hwnd, &me.paintStruct)
}

func (me *window) Decor() Widget {
	return me.decor
}

// TODO:: use int? set 0.5 and return 127,0.498039
func (me *window) Alpha() float32 {
	var pbAlpha byte
	extstyle, _ := win32.GetWindowLong(me.hwnd, win32.GWL_EXSTYLE)
	win32.SetWindowLong(me.hwnd, win32.GWL_EXSTYLE, extstyle|win32.WS_EX_LAYERED)
	if err := win32.GetLayeredWindowAttributes(me.hwnd, nil, &pbAlpha, nil); err == nil {
		return float32(float32(pbAlpha) / 255.0)
	}
	return 1.0
}

func (me *window) SetAlpha(alpha float32) {
	extstyle, _ := win32.GetWindowLong(me.hwnd, win32.GWL_EXSTYLE)
	win32.SetWindowLong(me.hwnd, win32.GWL_EXSTYLE, extstyle|win32.WS_EX_LAYERED)
	win32.SetLayeredWindowAttributes(me.hwnd, 0, byte(255.0*alpha), win32.LWA_ALPHA)
}

func (me *window) Title() (title string) {
	title, err := win32.GetWindowText(me.hwnd)
	if err != nil {
		log.E("nux", "get window title error: %s", err.Error())
		return
	}
	return title
}

func (me *window) SetTitle(title string) {
	if err := win32.SetWindowText(me.hwnd, title); err != nil {
		log.E("nux", "set window title error: %s", err.Error())
	}
}

func (me *window) SetDelegate(delegate WindowDelegate) {
	me.delegate = delegate
}

func (me *window) Delegate() WindowDelegate {
	return me.delegate
}

func (me *window) handlePointerEvent(e PointerEvent) {
	me.switchFocusIfPossible(e)

	if me.delegate != nil {
		if f, ok := me.delegate.(windowDelegate_HandlePointerEvent); ok {
			f.HandlePointerEvent(e)
			return
		}
	}

	hitTestResultManagerInstance.handlePointerEvent(me.Decor(), e)
}

func (me *window) handleScrollEvent(e ScrollEvent) {
	hitTestResultManagerInstance.handleScrollEvent(me.Decor(), e)
}

func (me *window) handleKeyEvent(e KeyEvent) {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(KeyEventHandler); ok {
			if f.OnKeyEvent(e) {
				return
			} else {
				goto other
			}
		}
	} else {
		goto other
	}

other:
	if me.decor != nil {
		me.handleOtherWidgetKeyEvent(me.decor, e)
	}
}

func (me *window) handleOtherWidgetKeyEvent(p Parent, e KeyEvent) bool {
	if p.ChildrenCount() > 0 {
		var compt Widget
		for _, c := range p.Children() {
			compt = nil
			if cpt, ok := c.(Component); ok {
				c = cpt.Content()
				compt = cpt
			}
			if cp, ok := c.(Parent); ok {
				if me.handleOtherWidgetKeyEvent(cp, e) {
					return true
				}
			} else if f, ok := c.(KeyEventHandler); ok {
				if f.OnKeyEvent(e) {
					return true
				}
			}

			if compt != nil {
				if f, ok := compt.(KeyEventHandler); ok {
					if f.OnKeyEvent(e) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (me *window) handleTypeEvent(e TypeEvent) {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(TypeEventHandler); ok {
			f.OnTypeEvent(e)
			return
		}
	}

	log.E("nuxui", "none widget handle typing event")
}

func (me *window) requestFocus(widget Widget) {
	if me.focusWidget == widget {
		return
	}

	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(Focus); ok && f.HasFocus() {
			f.FocusChanged(false)
		}
	}

	if f, ok := widget.(Focus); ok {
		if f.Focusable() {
			me.focusWidget = widget
			if !f.HasFocus() {
				f.FocusChanged(true)
			}
		}
	}
}

func (me *window) switchFocusIfPossible(event PointerEvent) {
	if event.Type() != Type_PointerEvent || !event.IsPrimary() {
		return
	}

	switch event.Action() {
	case Action_Down:
		me.initEvent = event
		if event.Kind() == Kind_Mouse {
			me.switchFocusAtPoint(event.X(), event.Y())
		}
	case Action_Move:
		if event.Kind() == Kind_Touch {
			if me.timer != nil {
				if me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
					me.switchFocusAtPoint(me.initEvent.X(), me.initEvent.Y())
				}
			}
		}
	case Action_Up:
		if event.Kind() == Kind_Touch {
			if event.Time().UnixNano()-me.initEvent.Time().UnixNano() < GESTURE_LONG_PRESS_TIMEOUT.Nanoseconds() {
				me.switchFocusAtPoint(event.X(), event.Y())
			}
		}
	}

}

func (me *window) switchFocusAtPoint(x, y float32) {
	if me.focusWidget != nil {
		if s, ok := me.focusWidget.(Size); ok {
			frame := s.Frame()
			if x >= float32(frame.X) && x <= float32(frame.X+frame.Width) &&
				y >= float32(frame.Y) && y <= float32(frame.Y+frame.Height) {
				// point is in current focus widget, do not need change focus
				return
			}
		}

		if f, ok := me.focusWidget.(Focus); ok && f.HasFocus() {
			me.focusWidget = nil
			f.FocusChanged(false)
		}
	}

}
