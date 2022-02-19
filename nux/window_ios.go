// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ios

package nux

/*
#include <stdint.h>
#import <UIKit/UIKit.h>

// char*   window_title(uintptr_t window);
// void    window_setTitle(uintptr_t window, char* title);
// float   window_alpha(uintptr_t window);
// void    window_setAlpha(uintptr_t window, float alpha);
void window_getSize(uintptr_t window, int32_t *width, int32_t *height);
void window_getContentSize(uintptr_t window, int32_t *width, int32_t *height);
CGContextRef window_getCGContext(uintptr_t window);
*/
import "C"

import (
	"github.com/nuxui/nuxui/log"
	"unsafe"
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

	cgContext uintptr
}

func newWindow(attr Attr) *window {
	me := &window{
	}

	me.CreateDecor(attr)
	GestureManager().AddGestureHandler(me.decor, &decorGestureHandler{})
	log.I("nuxui", "====== mount decor")
	mountWidget(me.decor, nil)
	log.I("nuxui", "====== mount decor end ")
	return me
}

func (me *window) CreateDecor(attr Attr) Widget {
	creator := FindRegistedWidgetCreator("github.com/nuxui/nuxui/ui.Layer")
	w := creator(ctx, attr)
	if p, ok := w.(Parent); ok {
		me.decor = p
	} else {
		log.Fatal("nuxui", "decor must is a Parent")
	}

	decorWindowList[w] = me

	return me.decor
}

func (me *window) Draw(canvas Canvas) {
	log.V("nuxui", "window Draw start")
	if me.decor != nil {
		if f, ok := me.decor.(Draw); ok {
			log.V("nuxui", "window Draw canvas save")
			canvas.Save()
			// TODO:: canvas clip
			// canvas.Translate(0, me.ContentHeight())
			// canvas.Scale(1, -1)
			log.V("nuxui", "window Draw canvas scale then draw")
			f.Draw(canvas)
			canvas.Restore()
		}
	}
	log.V("nuxui", "window Draw end")
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
	return newCanvas(C.window_getCGContext(me.windptr))
}

func (me *window) UnlockCanvas(c Canvas) {
}

func (me *window) Decor() Widget {
	return me.decor
}

func (me *window) Alpha() float32 {
	return 0 //float32(C.window_alpha(me.windptr))
}

func (me *window) SetAlpha(alpha float32) {
	// C.window_setAlpha(me.windptr, C.float(alpha))
}

func (me *window) Title() string {
	return "" //C.GoString(C.window_title(me.windptr))
}

func (me *window) SetTitle(title string) {
	// ctitle := C.CString(title)
	// C.window_setTitle(me.windptr, ctitle)
	// C.free(unsafe.Pointer(ctitle))
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
				y >= float32(frame.Position.Y) && y <= float32(frame.Position.Y+frame.Height) {
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
