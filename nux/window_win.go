// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

/*
#include <windows.h>
#include <windowsx.h>
*/
import "C"

import (
	"syscall"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

// merge window and activity ?
type window struct {
	windptr      C.HWND // TODO:: need by user code, how to use unsafe.Pointer
	decor        Parent
	buffer       unsafe.Pointer
	bufferWidth  int32
	bufferHeight int32
	delegate     WindowDelegate
	focusWidget  Widget

	initEvent PointerEvent
	timer     Timer

	surface        *Surface
	surfaceResized bool
	paintStruct    C.PAINTSTRUCT
}

func newWindow() *window {
	return &window{
		surfaceResized: true,
	}
}

func (me *window) Creating(attr Attr) {
	// TODO width ,height background drawable ...
	// create decor widget

	if me.decor == nil {
		me.CreatDecor(attr)
		GestureBinding().AddGestureHandler(me.decor, &decorGestureHandler{})
	}

	if c, ok := me.decor.(Creating); ok {
		c.Creating(attr)
	}
}

// TODO:: Created(content Widget) content is decor
func (me *window) Created(data interface{}) {
	main := data.(string)
	if main == "" {
		log.Fatal("nuxui", "no main widget found.")
	} else {
		mainWidgetCreator := FindRegistedWidgetCreatorByName(main)
		widgetTree := RenderWidget(mainWidgetCreator())
		me.decor.AddChild(widgetTree)
	}

	if c, ok := me.decor.(Created); ok {
		c.Created(me.decor)
	}
}

func (me *window) CreatDecor(attr Attr) Widget {
	creator := FindRegistedWidgetCreatorByName("github.com/nuxui/nuxui/ui.Layer")
	w := creator()
	if p, ok := w.(Parent); ok {
		me.decor = p
	} else {
		log.Fatal("nuxui", "decor must is a Parent")
	}

	decorWindowList[w] = me

	return me.decor
}

func (me *window) Measure(width, height int32) {
	if me.decor == nil {
		return
	}

	if s, ok := me.decor.(Size); ok {
		if s.MeasuredSize().Width == width && s.MeasuredSize().Height == height {
			// return
		}

		s.MeasuredSize().Width = width
		s.MeasuredSize().Height = height
	}

	me.surfaceResized = true

	if f, ok := me.decor.(Measure); ok {
		f.Measure(width, height)
	}
}

func (me *window) Layout(dx, dy, left, top, right, bottom int32) {
	if me.decor == nil {
		return
	}

	if f, ok := me.decor.(Layout); ok {
		f.Layout(dx, dy, left, top, right, bottom)
	}
}

func (me *window) Draw(canvas Canvas) {
	log.V("nuxui", "window Draw start")
	if me.decor != nil {
		if f, ok := me.decor.(Draw); ok {
			log.V("nuxui", "window Draw canvas save")
			canvas.Save()
			// TODO:: canvas clip
			f.Draw(canvas)
			canvas.Restore()
		}
	}
	log.V("nuxui", "window Draw end")
}

func (me *window) ID() uint64 {
	return 0
}

func (me *window) Width() int32 {
	var rect C.RECT
	if C.GetWindowRect(me.windptr, &rect) > 0 {
		return int32(rect.right - rect.left)
	}
	return 0
}

func (me *window) Height() int32 {
	var rect C.RECT
	if C.GetWindowRect(me.windptr, &rect) > 0 {
		return int32(rect.bottom - rect.top)
	}
	return 0
}

func (me *window) ContentWidth() int32 {
	var rect C.RECT
	if C.GetClientRect(me.windptr, &rect) > 0 {
		return int32(rect.right - rect.left)
	}
	return 0
}

func (me *window) ContentHeight() int32 {
	var rect C.RECT
	if C.GetClientRect(me.windptr, &rect) > 0 {
		return int32(rect.bottom - rect.top)
	}
	return 0
}

func (me *window) LockCanvas() (Canvas, error) {
	hdc := C.BeginPaint(me.windptr, &me.paintStruct)
	me.surface = newSurfaceWin32(hdc)
	return me.surface.GetCanvas(), nil
}

func (me *window) UnlockCanvas() error {
	me.surface.GetCanvas().Destroy()
	me.surface.Flush()
	C.EndPaint(me.windptr, &me.paintStruct)
	return nil
}

func (me *window) Decor() Widget {
	return me.decor
}

// TODO:: use int? set 0.5 and return 127,0.498039
func (me *window) Alpha() float32 {
	var pcrKey C.COLORREF
	var pbAlpha C.BYTE
	var pdwFlags C.DWORD
	if C.GetLayeredWindowAttributes(me.windptr, &pcrKey, &pbAlpha, &pdwFlags) > 0 {
		return float32(float32(pbAlpha) / 255.0)
	}
	return 1.0
}

func (me *window) SetAlpha(alpha float32) {
	C.SetWindowLong(me.windptr, C.GWL_EXSTYLE, C.GetWindowLong(me.windptr, C.GWL_EXSTYLE)|C.WS_EX_LAYERED)
	C.SetLayeredWindowAttributes(me.windptr, 0, C.BYTE(255.0*alpha), C.LWA_ALPHA)
}

func (me *window) Title() string {
	textLen := C.GetWindowTextLengthW(me.windptr) + 1
	buf := make([]uint16, textLen)
	C.GetWindowTextW(me.windptr, (*C.WCHAR)(&buf[0]), textLen)
	return syscall.UTF16ToString(buf)
}

func (me *window) SetTitle(title string) {
	ctitle, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		log.E("nux", "set window title error: %s", err.Error())
	} else {
		C.SetWindowTextW(me.windptr, (*C.WCHAR)(ctitle))
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

	gestureManagerInstance.handlePointerEvent(me.Decor(), e)
}

func (me *window) handleScrollEvent(e ScrollEvent) {
	gestureManagerInstance.handleScrollEvent(me.Decor(), e)
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
				compt = cpt.Component()
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

func (me *window) handleTypingEvent(e TypingEvent) {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(TypingEventHandler); ok {
			f.OnTypingEvent(e)
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
			ms := s.MeasuredSize()
			if x >= float32(ms.Position.X) && x <= float32(ms.Position.X+ms.Width) &&
				y >= float32(ms.Position.Y) && y <= float32(ms.Position.Y+ms.Height) {
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
