// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

import (
	"runtime"
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/cairo"
	"nuxui.org/nuxui/nux/internal/linux/xlib"
)

const XALL_EVENT = (xlib.KeyPressMask | xlib.KeyReleaseMask | xlib.ButtonPressMask |
	xlib.ButtonReleaseMask | xlib.EnterWindowMask | xlib.LeaveWindowMask | xlib.PointerMotionMask |
	/*xlib.PointerMotionHintMask*/ xlib.Button1MotionMask | xlib.Button2MotionMask |
	xlib.Button3MotionMask | xlib.Button4MotionMask | xlib.Button5MotionMask | xlib.ButtonMotionMask |
	xlib.KeymapStateMask | xlib.ExposureMask | xlib.VisibilityChangeMask | xlib.StructureNotifyMask |
	xlib.ResizeRedirectMask | xlib.SubstructureNotifyMask | xlib.SubstructureRedirectMask |
	xlib.FocusChangeMask | xlib.PropertyChangeMask | xlib.ColormapChangeMask | xlib.OwnerGrabButtonMask)

const AllMaskBits = (xlib.CWBackPixmap | xlib.CWBackPixel | xlib.CWBorderPixmap |
	xlib.CWBorderPixel | xlib.CWBitGravity | xlib.CWWinGravity |
	xlib.CWBackingStore | xlib.CWBackingPlanes | xlib.CWBackingPixel |
	xlib.CWOverrideRedirect | xlib.CWSaveUnder | xlib.CWEventMask |
	xlib.CWDontPropagate | xlib.CWColormap | xlib.CWCursor)

type nativeWindow struct {
	window  xlib.Window
	parent  xlib.Window
	display *xlib.Display
	visual  *xlib.Visual
	colormap xlib.Colormap
	screenNum int32
	depth int32
	surface *cairo.Surface
	title string
}

func newNativeWindow(attr Attr) *nativeWindow {
	me := &nativeWindow{}
	me.display = theApp.nativeApp.display
	me.screenNum = xlib.XDefaultScreen(me.display)
	me.parent = xlib.XRootWindow(me.display, me.screenNum)
	me.depth = xlib.XDefaultDepth(me.display, me.screenNum)
	me.visual = xlib.XDefaultVisual(me.display, me.screenNum)

	var attrs xlib.XSetWindowAttributes
	attrs.EventMask = XALL_EVENT
	attrs.Colormap = xlib.XCreateColormap(me.display, me.parent, me.visual, xlib.AllocNone)
	me.colormap = attrs.Colormap
	
	if attrs.Colormap == 0 {
		log.E("nuxui","XCreateColormap failed")
	}

	me.window = xlib.XCreateWindow(me.display, me.parent, 0, 0, 600, 400, 0, me.depth, xlib.InputOutput, me.visual, AllMaskBits, &attrs)

	runtime.SetFinalizer(me, freeNativeWindow)
	return me
}

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
	// me.surface = C.cairo_xlib_surface_create(me.display, me.windptr, me.visual, C.int(w), C.int(h))
	me.surface = cairo.XlibSurfaceCreate(me.display, xlib.Drawable(me.window), me.visual, w, h)
	return newCanvas(me.surface)
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