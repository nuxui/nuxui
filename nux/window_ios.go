// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

import (
	"runtime"
	"time"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/ios"
)

type nativeWindow struct {
	ptr ios.UIWindow

	pendingAttr Attr
}

func newNativeWindow(attr Attr) (me *nativeWindow) {
	ios.SetWindowEventHandler(nativeWindowEventHandler)

	width, height := measureWindowSize(attr.GetDimen("width", "100%"), attr.GetDimen("height", "100%"))
	me = &nativeWindow{
		ptr: ios.NewUIWindow(width, height),
	}
	// me.SetTitle(attr.GetString("title", ""))

	runtime.SetFinalizer(me, freeWindow)
	return me
}

func freeWindow(me *nativeWindow) {
}

func (me *nativeWindow) Center() {
}

func (me *nativeWindow) Show() {
	me.ptr.MakeKeyAndVisible()
}

func (me *nativeWindow) ContentSize() (width, height int32) {
	w, h := me.ptr.Frame().Size()
	return int32(w), int32(h)
}

func (me *nativeWindow) Title() string {
	// return me.ptr.Title()
	return ""
}

func (me *nativeWindow) SetTitle(title string) {
	// me.ptr.SetTitle(title)
}

func (me *nativeWindow) lockCanvas() Canvas {
	return newCanvas(ios.UIGraphicsGetCurrentContext())
}

func (me *nativeWindow) unlockCanvas(canvas Canvas) {
	canvas.Flush()
}

func (me *nativeWindow) draw(canvas Canvas, decor Widget) {
	if decor != nil {
		if f, ok := decor.(Draw); ok {
			canvas.Save()
			if TestDraw != nil {
				TestDraw(canvas)
			} else {
				f.Draw(canvas)
			}
			canvas.Restore()
		}
	}
}

func nativeWindowEventHandler(event any) bool {
	switch e := event.(type) {
	case ios.UIEvent:
		return handleUIEvent(e)
	// case *ios.TypingEvent:
	// 	return handleTypingEvent(e)
	case *ios.WindowEvent:
		switch e.Type {
		case ios.Event_WindowDidResize:
			theApp.mainWindow.resize()
		case ios.Event_WindowDrawRect:
			theApp.mainWindow.draw()
		}
	}

	return false
}

func handleUIEvent(event ios.UIEvent) bool {
	switch event.Type() {
	case ios.UIEventTypeTouches:
		if event.AllTouches().Count() == 1 {
			return handlePointerEvent(event)
		}
	case ios.UIEventTypeMotion:
		log.I("nuxui", "handleUIEvent === UIEventTypeMotion")
	case ios.UIEventTypeRemoteControl:
		log.I("nuxui", "handleUIEvent === UIEventTypeRemoteControl")
	case ios.UIEventTypePresses:
		log.I("nuxui", "handleUIEvent === UIEventTypePresses")
	case ios.UIEventTypeHover:
		log.I("nuxui", "handleUIEvent === UIEventTypeHover")
	case ios.UIEventTypeScroll:
		log.I("nuxui", "handleUIEvent === UIEventTypeScroll")
	case ios.UIEventTypeTransform:
		log.I("nuxui", "handleUIEvent === UIEventTypeTransform")
	}
	return false
}

var lastMouseEvent map[PointerButton]PointerEvent = map[PointerButton]PointerEvent{}

func handlePointerEvent(uievent ios.UIEvent) bool {
	touch := ios.UITouch(uievent.AllTouches().ObjectAtIndex(0))
	x, y := touch.LocationInView(0) // Pass nil to get the touch location in the window's coordinates

	e := &pointerEvent{
		event: event{
			window: theApp.MainWindow(),
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer: 0,
		button:  ButtonPrimary,
		kind:    Kind_Touch,
		x:       x,
		y:       y,
	}

	switch touch.Phase() {
	case ios.UITouchPhaseBegan:
		e.action = Action_Down
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case ios.UITouchPhaseMoved:
		e.action = Action_Drag
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case ios.UITouchPhaseEnded:
		e.action = Action_Up
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	case ios.UITouchPhaseCancelled:
		e.action = Action_Up
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	}

	return App().MainWindow().handlePointerEvent(e)
}
