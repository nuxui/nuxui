// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go:build (linux && !android)

package xlib

/*
#include <X11/X.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xresource.h>
#include <X11/keysymdef.h>
#include <X11/extensions/Xrender.h>
#include <X11/XKBlib.h>
#include <X11/Xatom.h>

#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <cairo/cairo-xlib.h>

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <strings.h>
#include <string.h>
#include <locale.h>
#include <sys/syscall.h>
#include <stdint.h>
*/
import "C"

// X11/x_copy.h

type EventMask C.long

const (
	NoEventMask              EventMask = 0
	KeyPressMask             EventMask = 1 << 0
	KeyReleaseMask           EventMask = 1 << 1
	ButtonPressMask          EventMask = 1 << 2
	ButtonReleaseMask        EventMask = 1 << 3
	EnterWindowMask          EventMask = 1 << 4
	LeaveWindowMask          EventMask = 1 << 5
	PointerMotionMask        EventMask = 1 << 6
	PointerMotionHintMask    EventMask = 1 << 7
	Button1MotionMask        EventMask = 1 << 8
	Button2MotionMask        EventMask = 1 << 9
	Button3MotionMask        EventMask = 1 << 10
	Button4MotionMask        EventMask = 1 << 11
	Button5MotionMask        EventMask = 1 << 12
	ButtonMotionMask         EventMask = 1 << 13
	KeymapStateMask          EventMask = 1 << 14
	ExposureMask             EventMask = 1 << 15
	VisibilityChangeMask     EventMask = 1 << 16
	StructureNotifyMask      EventMask = 1 << 17
	ResizeRedirectMask       EventMask = 1 << 18
	SubstructureNotifyMask   EventMask = 1 << 19
	SubstructureRedirectMask EventMask = 1 << 20
	FocusChangeMask          EventMask = 1 << 21
	PropertyChangeMask       EventMask = 1 << 22
	ColormapChangeMask       EventMask = 1 << 23
	OwnerGrabButtonMask      EventMask = 1 << 24
)

type ErrorCode int

const (
	Success             ErrorCode = 0  /* everything's okay */
	BadRequest          ErrorCode = 1  /* bad request code */
	BadValue            ErrorCode = 2  /* int parameter out of range */
	BadWindow           ErrorCode = 3  /* parameter not a Window */
	BadPixmap           ErrorCode = 4  /* parameter not a Pixmap */
	BadAtom             ErrorCode = 5  /* parameter not an Atom */
	BadCursor           ErrorCode = 6  /* parameter not a Cursor */
	BadFont             ErrorCode = 7  /* parameter not a Font */
	BadMatch            ErrorCode = 8  /* parameter mismatch */
	BadDrawable         ErrorCode = 9  /* parameter not a Pixmap or Window */
	BadAccess           ErrorCode = 10 /* depending on context:*/
	BadAlloc            ErrorCode = 11 /* insufficient resources */
	BadColor            ErrorCode = 12 /* no such colormap */
	BadGC               ErrorCode = 13 /* parameter not a GC */
	BadIDChoice         ErrorCode = 14 /* choice not in range or already used */
	BadName             ErrorCode = 15 /* font or color name doesn't exist */
	BadLength           ErrorCode = 16 /* Request length incorrect */
	BadImplementation   ErrorCode = 17 /* server is defective */
	FirstExtensionError ErrorCode = 128
	LastExtensionError  ErrorCode = 255
)

const (
	InputOutput = 1
	InputOnly   = 2
)

const None = 0

// For CreateColormap
const (
	AllocNone = 0 // create map with no entries
	AllocAll  = 1 // allocate entire map writeable
)

const (
	XBufferOverflow = -1
	XLookupNone     = 1
	XLookupChars    = 2
	XLookupKeySym   = 3
	XLookupBoth     = 4
)

type CW uint64

const (
	CWBackPixmap       CW = (1 << 0)
	CWBackPixel        CW = (1 << 1)
	CWBorderPixmap     CW = (1 << 2)
	CWBorderPixel      CW = (1 << 3)
	CWBitGravity       CW = (1 << 4)
	CWWinGravity       CW = (1 << 5)
	CWBackingStore     CW = (1 << 6)
	CWBackingPlanes    CW = (1 << 7)
	CWBackingPixel     CW = (1 << 8)
	CWOverrideRedirect CW = (1 << 9)
	CWSaveUnder        CW = (1 << 10)
	CWEventMask        CW = (1 << 11)
	CWDontPropagate    CW = (1 << 12)
	CWColormap         CW = (1 << 13)
	CWCursor           CW = (1 << 14)
)

type EventType int

const (
	KeyPress         EventType = 2
	KeyRelease       EventType = 3
	ButtonPress      EventType = 4
	ButtonRelease    EventType = 5
	MotionNotify     EventType = 6
	EnterNotify      EventType = 7
	LeaveNotify      EventType = 8
	FocusIn          EventType = 9
	FocusOut         EventType = 10
	KeymapNotify     EventType = 11
	Expose           EventType = 12
	GraphicsExpose   EventType = 13
	NoExpose         EventType = 14
	VisibilityNotify EventType = 15
	CreateNotify     EventType = 16
	DestroyNotify    EventType = 17
	UnmapNotify      EventType = 18
	MapNotify        EventType = 19
	MapRequest       EventType = 20
	ReparentNotify   EventType = 21
	ConfigureNotify  EventType = 22
	ConfigureRequest EventType = 23
	GravityNotify    EventType = 24
	ResizeRequest    EventType = 25
	CirculateNotify  EventType = 26
	CirculateRequest EventType = 27
	PropertyNotify   EventType = 28
	SelectionClear   EventType = 29
	SelectionRequest EventType = 30
	SelectionNotify  EventType = 31
	ColormapNotify   EventType = 32
	ClientMessage    EventType = 33
	MappingNotify    EventType = 34
	GenericEvent     EventType = 35
	LASTEvent        EventType = 36
)

type PropertyMode int

const (
	PropModeReplace PropertyMode = 0
	PropModePrepend PropertyMode = 1
	PropModeAppend  PropertyMode = 2
)

const (
	XA_WM_NAME = Atom(C.XA_WM_NAME)
)
