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

char* window_title(uintptr_t window);
void window_setTitle(uintptr_t window, char* title);
float window_alpha(uintptr_t window);
void window_setAlpha(uintptr_t window, float alpha);
int32_t window_getWidth(uintptr_t window);
int32_t window_getHeight(uintptr_t window);
int32_t window_getContentWidth(uintptr_t window);
int32_t window_getContentHeight(uintptr_t window);
uint8_t* renew_buffer(int32_t width, int32_t height);
size_t get_pitch(int32_t width);
void flush_buffer(uintptr_t window, void* buffer0);
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
}

func newWindow() *window {
	return &window{}
}

func (me *window) Creating(attr Attr) {
	// TODO width ,height background drawable ...
	// create decor widget

	if me.decor == nil {
		me.decor = NewDecor()
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
		s.MeasuredSize().Width = width
		s.MeasuredSize().Height = height
	}

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
	// buffer size changed
	w := me.ContentWidth()
	h := me.ContentHeight()
	if me.bufferWidth != w || me.bufferHeight != h {
		log.V("nuxui", "buffer size changed xxxxx")
		if me.buffer != nil {
			C.free(me.buffer)
		}

		me.buffer = unsafe.Pointer(C.renew_buffer(C.int32_t(w), C.int32_t(h)))
		me.bufferWidth = w
		me.bufferHeight = h
	}

	surface := NewSurfaceFromData(unsafe.Pointer(me.buffer), FORMAT_ARGB32, int(w), int(h), int(C.get_pitch(C.int32_t(w))))
	return surface.GetCanvas(), nil
}

func (me *window) UnlockCanvas() error {
	C.flush_buffer(me.windptr, me.buffer)
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
	if me.delegate != nil {
		if f, ok := me.delegate.(windowDelegate_HandlePointerEvent); ok {
			f.HandlePointerEvent(e)
			return
		}
	}

	gestureManagerInstance.handlePointerEvent(me.Decor(), e)
}
