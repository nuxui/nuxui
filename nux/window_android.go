// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build android

package nux

/*
#cgo LDFLAGS: -landroid -llog

#include <android/configuration.h>
#include <android/input.h>
#include <android/keycodes.h>
#include <android/looper.h>
#include <android/native_activity.h>
#include <android/native_window.h>
#include <EGL/egl.h>
#include <jni.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"errors"

	"github.com/nuxui/nuxui/log"
)

// merge window and activity ?
type window struct {
	actptr  *C.ANativeActivity
	windptr *C.ANativeWindow
	decor   Parent
}

func newWindow() *window {
	return &window{}
}

func (me *window) Creating(attr Attr) {
	// TODO width ,height background drawable ...
	// create decor widget

	if me.decor == nil {
		me.decor = NewDecor()
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
	if me.windptr == nil {
		return 0
	}

	return int32(C.ANativeWindow_getWidth(me.windptr))
}
func (me *window) Height() int32 {
	if me.windptr == nil {
		return 0
	}

	return int32(C.ANativeWindow_getHeight(me.windptr))
}

func (me *window) LockCanvas() (Canvas, error) {
	var buffer C.ANativeWindow_Buffer
	var rect C.ARect
	if C.ANativeWindow_lock(me.windptr, &buffer, &rect) != 0 {
		return nil, errors.New("Unable to lock android window buffer")
	}

	surface := NewSurfaceFromData(buffer.bits, FORMAT_ARGB32, int(buffer.width), int(buffer.height), int(buffer.stride*4))
	return surface.GetCanvas(), nil
}

func (me *window) UnlockCanvas() error {
	if C.ANativeWindow_unlockAndPost(me.windptr) != 0 {
		errors.New("Unable to unlock android window buffer")
	}
	return nil
}

func (me *window) Decor() Widget {
	return me.decor
}
