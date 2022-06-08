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

type Window C.Window // XID => unsigned long
type Display C.Display
type Screen C.Screen
type Visual C.Visual
type XEvent C.XEvent
type Atom C.Atom // unsigned long
type XIC C.XIC   // unsigned long
type XIM C.XIM
type Status C.Status
type KeySym C.KeySym
type Colormap C.Colormap
type Cursor C.Cursor
type Pixmap C.Pixmap
type KeyCode C.KeyCode
type Drawable C.Drawable

// X11/xlib.h => struct XAnyEvent
type XAnyEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
}

type XButtonPressedEvent XButtonEvent
type XButtonReleasedEvent XButtonEvent
type XButtonEvent struct {
	Type       EventType
	Serial     C.ulong
	SendEvent  C.Bool
	Display    *Display
	Window     Window
	Root       Window
	SubWindow  Window
	Time       C.Time
	X          C.int
	Y          C.int
	Xroot      C.int
	Yroot      C.int
	State      C.uint
	Button     Button
	SameScreen C.Bool
}

type XKeyPressedEvent XKeyEvent
type XKeyReleasedEvent XKeyEvent
type XKeyEvent struct {
	Type       EventType
	Serial     C.ulong
	SendEvent  C.Bool
	Display    *Display
	Window     Window
	Root       Window
	SubWindow  Window
	Time       C.Time
	X          C.int
	Y          C.int
	Xroot      C.int
	Yroot      C.int
	State      C.uint
	Keycode    C.uint
	SameScreen C.Bool
}

type XPointerMovedEvent XMotionEvent
type XMotionEvent struct {
	Type       EventType
	Serial     C.ulong
	SendEvent  C.Bool
	Display    *Display
	Window     Window
	Root       Window
	SubWindow  Window
	Time       C.Time
	X          C.int
	Y          C.int
	Xroot      C.int
	Yroot      C.int
	State      ButtonMask
	IsHint     C.char
	SameScreen C.Bool
}

type XEnterWindowEvent XCrossingEvent
type XLeaveWindowEvent XCrossingEvent
type XCrossingEvent struct {
	Type       EventType
	Serial     C.ulong
	SendEvent  C.Bool
	Display    *Display
	Window     Window
	Root       Window
	SubWindow  Window
	Time       C.Time
	X          C.int
	Y          C.int
	Xroot      C.int
	Yroot      C.int
	Mode       C.int
	Detail     C.int
	SameScreen C.Bool
	Focus      C.Bool
	State      C.uint
}

type XFocusInEvent XFocusChangeEvent
type XFocusOutEvent XFocusChangeEvent
type XFocusChangeEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
	Mode      C.int
	Detail    C.int
}

type XKeymapEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
	KeyVector []C.char
}

type XExposeEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
	X         C.int
	Y         C.int
	Width     C.int
	Height    C.int
	Count     C.int
}

type XGraphicsExposeEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Drawable  C.Drawable
	X         C.int
	Y         C.int
	Width     C.int
	Height    C.int
	Count     C.int
	MajorCode C.int
	MinorCode C.int
}

type XNoExposeEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Drawable  C.Drawable
	MajorCode C.int
	MinorCode C.int
}

type XVisibilityEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
	State     C.int
}

type XCreateWindowEvent struct {
	Type             EventType
	Serial           C.ulong
	SendEvent        C.Bool
	Display          *Display
	Parent           Window
	Window           Window
	X                C.int
	Y                C.int
	Width            C.int
	Height           C.int
	BorderWidth      C.int
	OverrideRedirect C.Bool
}

type XDestroyWindowEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Event     Window
	Window    Window
}

type XUnmapEvent struct {
	Type          EventType
	Serial        C.ulong
	SendEvent     C.Bool
	Display       *Display
	Event         Window
	Window        Window
	FromConfigure C.Bool
}

type XMapEvent struct {
	Type             EventType
	Serial           C.ulong
	SendEvent        C.Bool
	Display          *Display
	Event            Window
	Window           Window
	OverrideRedirect C.Bool
}

type XMapRequestEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Parent    Window
	Window    Window
}

type XReparentEvent struct {
	Type             EventType
	Serial           C.ulong
	SendEvent        C.Bool
	Display          *Display
	Event            Window
	Window           Window
	Parent           Window
	X                C.int
	Y                C.int
	OverrideRedirect C.Bool
}

type XConfigureEvent struct {
	Type             EventType
	Serial           C.ulong
	SendEvent        C.Bool
	Display          *Display
	Event            Window
	Window           Window
	X                C.int
	Y                C.int
	Width            C.int
	Height           C.int
	BorderWidth      C.int
	Above            Window
	OverrideRedirect C.Bool
}

type XGravityEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Event     Window
	Window    Window
	X         C.int
	Y         C.int
}

type XResizeRequestEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Event     Window
	Window    Window
	Width     C.int
	Height    C.int
}

type XConfigureRequestEvent struct {
	Type        EventType
	Serial      C.ulong
	SendEvent   C.Bool
	Display     *Display
	Parent      Window
	Window      Window
	X           C.int
	Y           C.int
	Width       C.int
	Height      C.int
	BorderWidth C.int
	Above       Window
	Detail      C.int
	ValueMask   C.ulong
}

type XCirculateEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Event     Window
	Window    Window
	Place     C.int
}

type XCirculateRequestEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Parent    Window
	Window    Window
	Place     C.int
}

type XPropertyEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
	Atom      Atom
	Time      C.Time
	State     C.int
}

type XSelectionClearEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
	Selection Atom
	Time      C.Time
}

type XSelectionRequestEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Owner     Window
	Requestor Window
	Selection Atom
	Target    Atom
	Property  Atom
	Time      C.Time
}

type XSelectionEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Requestor Window
	Selection Atom
	Target    Atom
	Property  Atom
	Time      C.Time
}

type XColormapEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Window    Window
	Colormap  C.Colormap
	New       C.Bool
	state     C.int
}

type XClientMessageEvent struct {
	Type        EventType
	Serial      C.ulong
	SendEvent   C.Bool
	Display     *Display
	Window      Window
	MessageType Atom
	Format      C.int
	Data        [20]byte // save Atom data
}

type XMappingEvent struct {
	Type         EventType
	Serial       C.ulong
	SendEvent    C.Bool
	Display      *Display
	Window       Window
	Request      C.int
	FirstKeycode C.int
	Count        C.int
}

type XErrorEvent struct {
	Type        EventType
	Display     *Display
	ResourceID  C.XID
	Serial      C.ulong
	ErrorCode   C.uchar
	RequestCode C.uchar
	MinorCode   C.uchar
}

type XGenericEvent struct {
	Type      EventType
	Serial    C.ulong
	SendEvent C.Bool
	Display   *Display
	Extension C.int
	EvType    EventType
}

type XSetWindowAttributes struct {
	BackgroundPixmap   Pixmap
	BackgroundPixel    C.ulong
	borderPixmap       Pixmap
	BorderPixel        C.ulong
	BitGravity         C.int
	WinGravity         C.int
	BackingStore       C.int
	BackingPlanes      C.ulong
	BackingPixel       C.ulong
	SaveUnder          C.Bool
	EventMask          EventMask
	DoNotPropagateMask C.long
	OverrideRedirect   C.Bool
	Colormap           Colormap
	Cursor             Cursor
}

type XWindowAttributes struct {
	X                  C.int
	Y                  C.int
	Width              C.int
	Height             C.int
	BorderWidth        C.int
	Depth              C.int
	visual             *Visual
	Root               Window
	Class              C.int
	BitGravity         C.int
	WinGravity         C.int
	BackingStore       C.int
	BackingPlanes      C.ulong
	BackingPixel       C.ulong
	SaveUnder          C.Bool
	colormap           Colormap
	MapInstalled       C.Bool
	MapState           C.int
	AllEventMasks      C.long
	YourEventMask      C.long
	DoNotPropagateMask C.long
	OverrideRedirect   C.Bool
	Screen             *Screen
}

type XKeyboardState struct {
	KeyClickPercent  C.int
	BellPercent      C.int
	BellPitch        C.uint
	BellDuration     C.uint
	LedMask          C.ulong
	GlobalAutoRepeat C.int
	AutoRepeats      [32]C.char
}
