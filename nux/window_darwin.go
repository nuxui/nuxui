// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>

char*   window_title(uintptr_t window);
void    window_setTitle(uintptr_t window, char* title);
float   window_alpha(uintptr_t window);
void    window_setAlpha(uintptr_t window, float alpha);
void    window_getSize(uintptr_t window, int32_t *width, int32_t *height);
void    window_getContentSize(uintptr_t window, int32_t *width, int32_t *height);
CGContextRef window_getCGContext(uintptr_t window);
*/
import "C"

import (
	"unsafe"

	"nuxui.org/nuxui/log"
)

type window struct {
	windptr      C.uintptr_t
	decor        Parent
	buffer       unsafe.Pointer
	bufferWidth  int32
	bufferHeight int32
	delegate     WindowDelegate
	focusWidget  Widget

	initEvent PointerEvent
	timer     Timer

	canvas Canvas
}

func newWindow(attr Attr) *window {
	me := &window{}

	me.CreateDecor(attr)
	return me
}

func (me *window) CreateDecor(attr Attr) Widget {
	creator := FindRegistedWidgetCreator("nuxui.org/nuxui/ui.Layer")
	w := creator(attr)
	if p, ok := w.(Parent); ok {
		me.decor = p
	} else {
		log.Fatal("nuxui", "decor must is a Parent")
	}

	decorWindowList[w] = me

	return me.decor
}

var TestDraw func(Canvas)

func (me *window) Draw(canvas Canvas) {
	// log.V("nuxui", "window Draw start")

	if me.decor != nil {
		if f, ok := me.decor.(Draw); ok {
			_, h := me.ContentSize()
			canvas.Save()
			canvas.Translate(0, float32(h))
			canvas.Scale(1, -1)
			f.Draw(canvas)
			canvas.Restore()
		}
	}

	if TestDraw != nil {
		TestDraw(canvas)
		return
	}
}

func (me *window) ID() uint64 {
	return 0
}

func (me *window) Size() (width, height int32) {
	var w, h C.int32_t
	C.window_getSize(me.windptr, &w, &h)
	return int32(w), int32(h)
}

func (me *window) ContentSize() (width, height int32) {
	var w, h C.int32_t
	C.window_getContentSize(me.windptr, &w, &h)
	return int32(w), int32(h)
}

func (me *window) LockCanvas() Canvas {
	me.canvas = newCanvas(C.window_getCGContext(me.windptr))
	return me.canvas
}

func (me *window) UnlockCanvas(c Canvas) {
	me.canvas.Flush()
}

func (me *window) Decor() Widget {
	return me.decor
}

func (me *window) Alpha() float32 {
	return float32(C.window_alpha(me.windptr))
}

func (me *window) SetAlpha(alpha float32) {
	C.window_setAlpha(me.windptr, C.float(alpha))
}

func (me *window) Title() string {
	return C.GoString(C.window_title(me.windptr))
}

func (me *window) SetTitle(title string) {
	ctitle := C.CString(title)
	C.window_setTitle(me.windptr, ctitle)
	C.free(unsafe.Pointer(ctitle))
}

func (me *window) SetDelegate(delegate WindowDelegate) {
	me.delegate = delegate
}

func (me *window) Delegate() WindowDelegate {
	return me.delegate
}

func (me *window) handlePointerEvent(e PointerEvent) {
	me.switchFocusIfPossible(e)

	// log.I("nuxui", "handlePointerEvent x=%f, y=%f", e.X(), e.Y())

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
			f := s.Frame()
			if x >= float32(f.X) && x <= float32(f.X+f.Width) &&
				y >= float32(f.Y) && y <= float32(f.Y+f.Height) {
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
