// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package nux

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa -framework OpenGL
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>

char*   window_title(uintptr_t window);
void    window_setTitle(uintptr_t window, char* title);
float   window_alpha(uintptr_t window);
void    window_setAlpha(uintptr_t window, float alpha);
int32_t window_getWidth(uintptr_t window);
int32_t window_getHeight(uintptr_t window);
int32_t window_getContentWidth(uintptr_t window);
int32_t window_getContentHeight(uintptr_t window);
void*   window_getCGContext(uintptr_t window);
*/
import "C"

import (
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

// merge window and activity ?
type window struct {
	windptr      C.uintptr_t
	decor        Parent
	buffer       unsafe.Pointer
	bufferWidth  int32
	bufferHeight int32
	delegate     WindowDelegate
	focusWidget  Widget

	initEvent Event
	timer     Timer

	surface        *Surface
	surfaceResized bool
	cgContext      uintptr
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
		me.decor = newDecor(me)
		GestureBinding().AddGestureHandler(me.decor, &decorGestureHandler{})
	}

	if c, ok := me.decor.(Creating); ok {
		c.Creating(attr)
	}
}

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

// this canvas is nil
func (me *window) Draw(canvas Canvas) {
	if me.decor != nil {
		if f, ok := me.decor.(Draw); ok {
			canvas.Save()
			// TODO:: canvas clip
			canvas.Translate(0, me.ContentHeight())
			canvas.Scale(1, -1)
			f.Draw(canvas)
			canvas.Restore()
		}
	}
}

func (me *window) ID() uint64 {
	return 0
}

func (me *window) Width() int32 {
	return int32(C.window_getWidth(me.windptr))
}

func (me *window) Height() int32 {
	return int32(C.window_getHeight(me.windptr))
}

func (me *window) ContentWidth() int32 {
	return int32(C.window_getContentWidth(me.windptr))
}

func (me *window) ContentHeight() int32 {
	return int32(C.window_getContentHeight(me.windptr))
}

func (me *window) LockCanvas() (Canvas, error) {
	w := me.ContentWidth()
	h := me.ContentHeight()
	if me.surface == nil {
		me.cgContext = uintptr(C.window_getCGContext(me.windptr))
		me.surface = newSurfaceQuartzWithCGContext(me.cgContext, w, h)
	} else {
		// TODO:: did cairo_surface has method to resize instead of recreate
		if me.surfaceResized {
			me.surface.Destory()
			me.cgContext = uintptr(C.window_getCGContext(me.windptr))
			me.surface = newSurfaceQuartzWithCGContext(me.cgContext, w, h)
			me.surfaceResized = false
		}
	}
	return me.surface.GetCanvas(), nil
}

func (me *window) UnlockCanvas() error {
	// me.surface.WriteToPng("./a.png")
	me.surface.Flush()
	return nil
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
	C.window_setTitle(me.windptr, C.CString(title))
}

func (me *window) SetDelegate(delegate WindowDelegate) {
	me.delegate = delegate
}

func (me *window) Delegate() WindowDelegate {
	return me.delegate
}

func (me *window) handlePointerEvent(e Event) {
	me.switchFocusIfPossible(e)

	if me.delegate != nil {
		if f, ok := me.delegate.(windowDelegate_HandlePointerEvent); ok {
			f.HandlePointerEvent(e)
			return
		}
	}

	gestureManagerInstance.handlePointerEvent(me.Decor(), e)
}

func (me *window) handleKeyEvent(e Event) {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(KeyEventHandler); ok {
			if f.OnKeyEvent(e) {
				return
			}
		}
	}

	// TODO:: handle key event
}

func (me *window) handleTypingEvent(e Event) {
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

func (me *window) switchFocusIfPossible(event Event) {
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
