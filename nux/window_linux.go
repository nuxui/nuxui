// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/linux/xlib"
	"runtime"
	"time"
)

const XALL_EVENT = (xlib.KeyPressMask | xlib.KeyReleaseMask | xlib.ButtonPressMask |
	xlib.ButtonReleaseMask | xlib.EnterWindowMask | xlib.LeaveWindowMask | xlib.PointerMotionMask |
	/*xlib.PointerMotionHintMask*/ xlib.Button1MotionMask | xlib.Button2MotionMask |
	xlib.Button3MotionMask | xlib.Button4MotionMask | xlib.Button5MotionMask | xlib.ButtonMotionMask |
	xlib.KeymapStateMask | xlib.ExposureMask | xlib.VisibilityChangeMask | xlib.StructureNotifyMask |
	/*xlib.ResizeRedirectMask |*/ xlib.SubstructureNotifyMask | /*xlib.SubstructureRedirectMask |*/
	xlib.FocusChangeMask | xlib.PropertyChangeMask /*| xlib.ColormapChangeMask | xlib.OwnerGrabButtonMask*/)

const XALL_EVENT2 = (xlib.StructureNotifyMask | xlib.KeyPressMask | xlib.KeyReleaseMask |
	xlib.PointerMotionMask | xlib.ButtonPressMask | xlib.ButtonReleaseMask |
	xlib.ExposureMask | xlib.FocusChangeMask | xlib.VisibilityChangeMask |
	xlib.EnterWindowMask | xlib.LeaveWindowMask | xlib.PropertyChangeMask)
const AllMaskBits = (xlib.CWBorderPixel | xlib.CWColormap | xlib.CWEventMask)

type nativeWindow struct {
	window      xlib.Window
	parent      xlib.Window
	display     *xlib.Display
	visual      *xlib.Visual
	colormap    xlib.Colormap
	screenNum   int32
	depth       int32
	canvas      *canvas
	title       string
	width       int32
	height      int32
	sizeChanged bool
}

func newNativeWindow(attr Attr) *nativeWindow {
	me := &nativeWindow{}
	me.display = theApp.native.display
	me.screenNum = xlib.XDefaultScreen(me.display)
	me.parent = xlib.XRootWindow(me.display, me.screenNum)
	me.depth = xlib.XDefaultDepth(me.display, me.screenNum)
	me.visual = xlib.XDefaultVisual(me.display, me.screenNum)

	var attrs xlib.XSetWindowAttributes
	attrs.EventMask = XALL_EVENT
	attrs.Colormap = xlib.XCreateColormap(me.display, me.parent, me.visual, xlib.AllocNone)
	me.colormap = attrs.Colormap

	if attrs.Colormap == 0 {
		log.E("nuxui", "XCreateColormap failed")
	}

	me.width = 600
	me.height = 400
	me.sizeChanged = true
	me.window = xlib.XCreateWindow(me.display, me.parent, 0, 0, uint32(me.width), uint32(me.height), 0, me.depth, xlib.InputOutput, me.visual, AllMaskBits, &attrs)

	runtime.SetFinalizer(me, freeNativeWindow)
	return me
}

// XResizeWindow
// XGetWindowProperty

func freeNativeWindow(me *nativeWindow) {
	xlib.XFreeColormap(me.display, me.colormap)
}

func (me *nativeWindow) Center() {

}

func (me *nativeWindow) Show() {
	xlib.XMapWindow(me.display, me.window)
}

func (me *nativeWindow) ContentSize() (width, height int32) {
	var attrs xlib.XWindowAttributes
	xlib.XGetWindowAttributes(me.display, me.window, &attrs)
	return int32(attrs.Width), int32(attrs.Height)
}

func (me *nativeWindow) Title() string {
	// hard to get title from xlib
	return me.title
}

func (me *nativeWindow) SetTitle(title string) {
	me.title = title
	utf8Str := xlib.XInternAtom(me.display, "UTF8_STRING", false)
	xlib.XChangeProperty(me.display, me.window, xlib.XA_WM_NAME, utf8Str, 8, xlib.PropModeReplace, title)
}

func (me *nativeWindow) lockCanvas() Canvas {
	w, h := me.ContentSize()
	// log.I("nuxui", "lockCanvas w=%d, h=%d", w, h)
	// me.surface = cairo.XlibSurfaceCreate(me.display, xlib.Drawable(me.window), me.visual, w, h)
	if me.canvas == nil || me.sizeChanged {
		me.sizeChanged = false
		me.canvas = canvasFromWindow(me.display, xlib.Drawable(me.window), me.visual, w, h)
	}
	return me.canvas
}

func (me *nativeWindow) unlockCanvas(canvas Canvas) {
	canvas.Flush()
	// xlib.XFlush(me.display)
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

func (me *nativeWindow) handleNativeEvent(event any) bool {
	switch e := event.(type) {
	case *xlib.XPointerMovedEvent, *xlib.XButtonPressedEvent, *xlib.XButtonReleasedEvent:
		return me.handlePointerEvent(e)
	case *xlib.XKeyPressedEvent, *xlib.XKeyReleasedEvent:
	case *xlib.XEnterWindowEvent, *xlib.XLeaveWindowEvent:
	case *xlib.XFocusInEvent, *xlib.XFocusOutEvent:
	case *xlib.XPropertyEvent:
		// log.I("nuxui", "XPropertyEvent: %s", xlib.XGetAtomName(me.display, e.Atom))
	case *xlib.XConfigureEvent:
		// w, h := me.ContentSize()
		if me.width != int32(e.Width) || me.height != int32(e.Height) {
			me.width = int32(e.Width)
			me.height = int32(e.Height)
			me.sizeChanged = true
			theApp.window.resize()
		}
		// log.I("nuxui", "XConfigureEvent: width=%d, height=%d, w=%d, h=%d", e.Width, e.Height, w, h)
	case *xlib.XClientMessageEvent:
		if xlib.XGetAtomName(me.display, e.MessageType) == "nux_user_backToUI" {
			backToUI()
		}
	case *xlib.XExposeEvent:
		theApp.window.draw()
	case *xlib.XMapEvent: // show window
		theApp.window.resize()
	case *xlib.XResizeRequestEvent:
		// log.I("nuxui", "XResizeRequestEvent w=%d, h=%d", e.Width, e.Height)
		// theApp.window.measure()
		// theApp.window.layout()
		// InvalidateRect(0,0,0,0)
	case *xlib.XUnmapEvent: // hide window
	case *xlib.XVisibilityEvent:
	case *xlib.XDestroyWindowEvent:
	case *xlib.XErrorEvent:
		log.E("nuxui", "XErrorEvent %d", e.ErrorCode)
	}
	return false
}

var lastMouseEvent map[MouseButton]PointerEvent = map[MouseButton]PointerEvent{}

func (me *nativeWindow) handlePointerEvent(xevent any) bool {
	e := &pointerEvent{
		event: event{
			window: theApp.MainWindow(),
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer: 0,
		button:  MB_None,
		kind:    Kind_Mouse,
	}

	switch xe := xevent.(type) {
	case *xlib.XPointerMovedEvent:
		e.x = float32(xe.X)
		e.y = float32(xe.Y)
		e.action = Action_Hover
		e.button = MB_None
		e.pointer = 0
	case *xlib.XButtonPressedEvent:
		e.x = float32(xe.X)
		e.y = float32(xe.Y)
		e.action = Action_Down
		e.button = xbutton2nuxbutton(xe.Button)
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case *xlib.XButtonReleasedEvent:
		e.x = float32(xe.X)
		e.y = float32(xe.Y)
		e.action = Action_Up
		e.button = xbutton2nuxbutton(xe.Button)
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	}

	return App().MainWindow().handlePointerEvent(e)
}

func xbutton2nuxbutton(button xlib.Button) MouseButton {
	switch button {
	case xlib.Button1:
		return MB_Left
	case xlib.Button2:
		return MB_Middle
	case xlib.Button3:
		return MB_Right
	case 8:
		return MB_X2
	default:
		log.E("nuxui", "not handle button number %d", button)
	}
	return MB_None
}
