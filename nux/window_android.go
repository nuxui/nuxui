// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

/*
#cgo LDFLAGS: -llog

#include <jni.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

jobject surfaceHolder_lockCanvas(jobject surfaceHolder);
void surfaceHolder_unlockCanvas(jobject surfaceHolder, jobject canvas);

*/
import "C"
import (
	"nuxui.org/nuxui/log"
)

type window struct {
	jnienv        *C.JNIEnv
	activity      C.jobject
	surfaceHolder C.jobject
	decor         Parent
	width         int32
	height        int32
	delegate      WindowDelegate
	focusWidget   Widget

	initEvent PointerEvent
	timer     Timer
}

func newWindow(attr Attr) *window {
	me := &window{}

	me.CreateDecor(me.attr)
	return me
}

func (me *window) CreateDecor(attr Attr) Widget {
	creator := FindRegistedWidgetCreator("nuxui.org/nuxui/ui.Layer")
	attr.Set("padding", Attr{"top": "75px"})
	w := creator(attr)
	if p, ok := w.(Parent); ok {
		me.decor = p
	} else {
		log.Fatal("nuxui", "decor must is a Parent")
	}

	decorWindowList[w] = me

	return me.decor
}

func (me *window) Draw(canvas Canvas) {
	if me.decor != nil {
		if f, ok := me.decor.(Draw); ok {
			canvas.Save()
			f.Draw(canvas)
			canvas.Restore()
		}
	}
}

func (me *window) ID() uint64 {
	return 0
}

func (me *window) Size() (width, height int32) {
	if me.jnienv == nil {
		return 0, 0
	}
	return me.width, me.height
}

func (me *window) ContentSize() (width, height int32) {
	return me.Size()
}

func (me *window) LockCanvas() Canvas {
	canvas := C.surfaceHolder_lockCanvas(me.surfaceHolder)
	if canvas == 0 {
		return nil
	}
	return newCanvas(canvas)
}

func (me *window) UnlockCanvas(c Canvas) {
	C.surfaceHolder_unlockCanvas(me.surfaceHolder, c.(*canvas).ptr)
}

func (me *window) Decor() Widget {
	return me.decor
}

func (me *window) Alpha() float32 {
	return 0
}

func (me *window) SetAlpha(alpha float32) {
}

func (me *window) Title() string {
	return ""
}

func (me *window) SetTitle(title string) {
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
