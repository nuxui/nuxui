// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>

#include <X11/Xlib.h>

void window_getSize(Display* display, Window window, int32_t *width, int32_t *height);
void window_setText(Display* display, Window window, char *name);
*/
import "C"

import (
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

// merge window and activity ?
type window struct {
	windptr      C.uintptr_t // TODO:: need by user code, how to use unsafe.Pointer
	decor        Parent
	buffer       unsafe.Pointer
	bufferWidth  int32
	bufferHeight int32
	delegate     WindowDelegate
	focusWidget  Widget

	initEvent PointerEvent
	timer     Timer

	surface        *Surface
	display        *C.Display
	visual         *C.Visual
	surfaceResized bool

	width  int32
	height int32
	title  string
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

	me.width = width
	me.height = height

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
			// canvas.DrawColor(0xffff00ff)
			canvas.Restore()
		}
	}
	log.V("nuxui", "window Draw end")
}

func (me *window) ID() uint64 {
	return 0
}

func (me *window) Width() int32 {
	var width, height C.int
	C.window_getSize(me.display, me.windptr, &width, &height)
	return int32(width)
}

func (me *window) Height() int32 {
	var width, height C.int
	C.window_getSize(me.display, me.windptr, &width, &height)
	return int32(height)
}

func (me *window) ContentWidth() int32 {
	// TODO::
	return me.Width()
}

func (me *window) ContentHeight() int32 {
	// TODO::
	return me.Height()
}

func (me *window) LockCanvas() (Canvas, error) {
	w := me.ContentWidth()
	h := me.ContentHeight()
	log.V("nux", "LockCanvas ###### %d, %d", w, h)
	me.surface = newSurfaceXlib(me.display, me.windptr, me.visual, w, h)
	return me.surface.GetCanvas(), nil
}

func (me *window) UnlockCanvas() error {
	me.surface.GetCanvas().Destroy()
	me.surface.Flush()
	me.surface.Destroy()
	C.XFlush(me.display)
	return nil
}

func (me *window) Decor() Widget {
	return me.decor
}

// TODO:: use int? set 0.5 and return 127,0.498039
func (me *window) Alpha() float32 {
	// TODO::
	return 1.0
}

func (me *window) SetAlpha(alpha float32) {
	// TODO::
}

func (me *window) Title() string {
	return me.title
}

func (me *window) SetTitle(title string) {
	me.title = title
	ctitle := C.CString(title)
	C.window_setText(me.display, me.windptr, ctitle)
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
