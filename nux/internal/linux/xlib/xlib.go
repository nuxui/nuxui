// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go:build (linux && !android)

package xlib

/*
#cgo pkg-config: x11

#include <X11/X.h>
#include <X11/XKBlib.h>
#include <X11/Xatom.h>
#include <X11/Xlib.h>
#include <X11/Xresource.h>
#include <X11/Xutil.h>
#include <X11/extensions/Xrender.h>
#include <X11/keysymdef.h>

#include <locale.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/syscall.h>
#include <unistd.h>

XIC nux_XCreateIC(XIM xim, Window window);
*/
import "C"
import (
	"os"
	"unsafe"
)

func XCreateWindow(display *Display, parent Window, x, y int32, width, height, borderWidth uint32, depth int32,
	class uint32, visual *Visual, valuemask CW, attrs *XSetWindowAttributes) Window {
	return Window(C.XCreateWindow((*C.Display)(display), C.Window(parent), C.int(x), C.int(y),
		C.uint(width), C.uint(height), C.uint(borderWidth), C.int(depth), C.uint(class), (*C.Visual)(visual),
		C.ulong(valuemask), (*C.XSetWindowAttributes)(unsafe.Pointer(attrs))))
}

func GetDisplayName() string {
	return os.Getenv("DISPLAY")
}

func XOpenDefaultDisplay() *Display {
	return (*Display)(C.XOpenDisplay(nil))
}

func XOpenDisplay(name string) *Display {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Display)(C.XOpenDisplay(cname))
}

func XDisplayWidth(display *Display, screenNumber int32) int32 {
	return int32(C.XDisplayWidth((*C.Display)(display), C.int(screenNumber)))
}

func XDisplayHeight(display *Display, screenNumber int32) int32 {
	return int32(C.XDisplayHeight((*C.Display)(display), C.int(screenNumber)))
}

func XDefaultScreenOfDisplay(display *Display) *Screen {
	return (*Screen)(C.XDefaultScreenOfDisplay((*C.Display)(display)))
}

func XDefaultScreen(display *Display) int32 {
	return int32(C.XDefaultScreen((*C.Display)(display)))
}

func XCloseDisplay(display *Display) {
	C.XCloseDisplay((*C.Display)(display))
}

func XInitThreads() {
	C.XInitThreads()
}

func XrmInitialize() {
	C.XrmInitialize()
}

func XSupportsLocale() bool {
	return C.XSupportsLocale() > 0
}

func XSetLocaleModifiers(modifier_list string) {
	cstr := C.CString(modifier_list)
	defer C.free(unsafe.Pointer(cstr))
	C.XSetLocaleModifiers(cstr)
}

func XGetLocaleModifiers() string {
	return C.GoString(C.XSetLocaleModifiers(nil))
}

func XRootWindow(display *Display, screen_number int32) Window {
	return Window(C.XRootWindow((*C.Display)(display), C.int(screen_number)))
}

func XDefaultDepth(display *Display, screen_number int32) int32 {
	return int32(C.XDefaultDepth((*C.Display)(display), C.int(screen_number)))
}

func XDefaultVisual(display *Display, screen_number int32) *Visual {
	return (*Visual)(C.XDefaultVisual((*C.Display)(display), C.int(screen_number)))
}

func XMapWindow(display *Display, window Window) int32 {
	return int32(C.XMapWindow((*C.Display)(display), C.Window(window)))
}

func XNextEvent(display *Display, event *XEvent) int32 {
	return int32(C.XNextEvent((*C.Display)(display), (*C.XEvent)(event)))
}

func XFilterEvent(event *XEvent, window Window) bool {
	return C.XFilterEvent((*C.XEvent)(event), C.Window(window)) > 0
}

func XGetAtomName(display *Display, atom Atom) string {
	return C.GoString(C.XGetAtomName((*C.Display)(display), C.Atom(atom)))
}

func XInternAtom(display *Display, atom_name string, only_if_exists bool) Atom {
	var exist C.Bool
	if only_if_exists {
		exist = 1
	}
	cstr := C.CString(atom_name)
	defer C.free(unsafe.Pointer(cstr))
	return Atom(C.XInternAtom((*C.Display)(display), cstr, exist))
}

func XSendEvent(display *Display, window Window, propagate bool, eventMask EventMask, event *XEvent) Status {
	var p C.Bool
	if propagate {
		p = 1
	}
	return Status(C.XSendEvent((*C.Display)(display), C.Window(window), p, C.long(eventMask), (*C.XEvent)(event)))
}

func Xutf8LookupString(xic XIC, event *XKeyPressedEvent) (buffer string, keySym KeySym, status Status) {
	var k *C.KeySym = (*C.KeySym)(&keySym)
	var s *C.Status = (*C.Status)(&status)
	var len C.int = 20
	text := make([]C.char, int(len+1))

	len = C.Xutf8LookupString(C.XIC(xic), (*C.XKeyPressedEvent)(unsafe.Pointer(event)), (*C.char)(unsafe.Pointer(&text[0])), len, k, s)
	if status == XBufferOverflow {
		text = make([]C.char, int(len+1))
		C.Xutf8LookupString(C.XIC(xic), (*C.XKeyPressedEvent)(unsafe.Pointer(event)), (*C.char)(unsafe.Pointer(&text[0])), len, k, s)
	}

	buffer = C.GoString((*C.char)(unsafe.Pointer(&text[0])))
	return
}

func XClearArea(display *Display, window Window, x, y, width, height int32, exposures bool) int {
	var exp C.Bool
	if exposures {
		exp = 1
	}

	return int(C.XClearArea((*C.Display)(display), C.Window(window), C.int(x), C.int(y), C.uint(width), C.uint(height), exp))
}

func XClearWindow(display *Display, window Window) int {
	return int(C.XClearWindow((*C.Display)(display), C.Window(window)))
}

func XOpenIM(display *Display) XIM {
	return XIM(C.XOpenIM((*C.Display)(display), nil, nil, nil))
}

func XCloseIM(xim XIM) Status {
	return Status(C.XCloseIM(C.XIM(xim)))
}

func XCreateIC(xim XIM, window Window) XIC {
	return XIC(C.nux_XCreateIC(C.XIM(xim), C.Window(window)))
}

func XSetICFocus(xic XIC) {
	C.XSetICFocus(C.XIC(xic))
}

func XCreateColormap(display *Display, window Window, visual *Visual, alloc int) Colormap {
	return Colormap(C.XCreateColormap((*C.Display)(display), C.Window(window), (*C.Visual)(visual), C.int(alloc)))
}

func XkbKeycodeToKeysym(display *Display, keyCode KeyCode, group, level int) KeySym {
	return KeySym(C.XkbKeycodeToKeysym((*C.Display)(display), C.KeyCode(keyCode), C.int(group), C.int(level)))
}

func XGetKeyboardControl(display *Display) (state XKeyboardState) {
	C.XGetKeyboardControl((*C.Display)(display), (*C.XKeyboardState)(unsafe.Pointer(&state)))
	return
}

func XGetWindowAttributes(display *Display, window Window, attrs *XWindowAttributes) Status {
	return Status(C.XGetWindowAttributes((*C.Display)(display), C.Window(window), (*C.XWindowAttributes)(unsafe.Pointer(attrs))))
}

func XSelectInput(display *Display, window Window, eventMask EventMask) int32 {
	return int32(C.XSelectInput((*C.Display)(display), C.Window(window), C.long(eventMask)))
}

func XChangeProperty(display *Display, window Window, property, typee Atom, format int, mode PropertyMode, data string /*, nelements int => strlen(data)*/) int {
	cstr := C.CString(data)
	defer C.free(unsafe.Pointer(cstr))
	return int(C.XChangeProperty((*C.Display)(display), C.Window(window), C.Atom(property), C.Atom(typee), C.int(format), C.int(mode), (*C.uchar)(unsafe.Pointer(cstr)), C.int(C.strlen(cstr))))
}

func XCreateFontCursor(display *Display, shape CursorShape) Cursor {
	return Cursor(C.XCreateFontCursor((*C.Display)(display), C.uint(shape)))
}

func XDefineCursor(display *Display, window Window, cursor Cursor) ErrorCode {
	return ErrorCode(C.XDefineCursor((*C.Display)(display), C.Window(window), C.Cursor(cursor)))
}

func XUndefineCursor(display *Display, window Window) ErrorCode {
	return ErrorCode(C.XUndefineCursor((*C.Display)(display), C.Window(window)))
}

func XFlush(display *Display) {
	C.XFlush((*C.Display)(display))
}

func XSync(display *Display, discard bool) {
	var d C.Bool
	if discard {
		d = 1
	}
	C.XSync((*C.Display)(display), d)
}

func XFree(data unsafe.Pointer) ErrorCode {
	return ErrorCode(C.XFree(data))
}

func XFreeColormap(display *Display, colormap Colormap) ErrorCode {
	return ErrorCode(C.XFreeColormap((*C.Display)(display), C.Colormap(colormap)))
}

// func XFreeColors(display *Display, colormap Colormap, pixels *uint64, npixels int, planes uint64)ErrorCode{
// 	return ErrorCode(C.XFreeColors((*C.Display)(display), C.Colormap(colormap), ))
// }

// extern int XFreeColors(
//     Display*		/* display */,
//     Colormap		/* colormap */,
//     unsigned long*	/* pixels */,
//     int			/* npixels */,
//     unsigned long	/* planes */
// );

// extern int XFreeCursor(
//     Display*		/* display */,
//     Cursor		/* cursor */
// );

// extern int XFreeExtensionList(
//     char**		/* list */
// );

// extern int XFreeFont(
//     Display*		/* display */,
//     XFontStruct*	/* font_struct */
// );

// extern int XFreeFontInfo(
//     char**		/* names */,
//     XFontStruct*	/* free_info */,
//     int			/* actual_count */
// );

// extern int XFreeFontNames(
//     char**		/* list */
// );

// extern int XFreeFontPath(
//     char**		/* list */
// );

// extern int XFreeGC(
//     Display*		/* display */,
//     GC			/* gc */
// );

// extern int XFreeModifiermap(
//     XModifierKeymap*	/* modmap */
// );

// extern int XFreePixmap(
//     Display*		/* display */,
//     Pixmap		/* pixmap */
// );
