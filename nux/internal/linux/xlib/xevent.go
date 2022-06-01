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

int XEvent_Type(XEvent* event){
	return event->type;
}
*/
import "C"

import "unsafe"

func (me *XEvent) Type() EventType {
	return EventType(C.XEvent_Type((*C.XEvent)(me)))
}

func (me *XEvent) Convert() any {
	switch me.Type() {
	case KeyPress:
		return (*XKeyPressedEvent)(unsafe.Pointer(me))
	case KeyRelease:
		return (*XKeyReleasedEvent)(unsafe.Pointer(me))
	case ButtonPress:
		return (*XButtonPressedEvent)(unsafe.Pointer(me))
	case ButtonRelease:
		return (*XButtonReleasedEvent)(unsafe.Pointer(me))
	case MotionNotify:
		return (*XPointerMovedEvent)(unsafe.Pointer(me))
	case EnterNotify:
		return (*XEnterWindowEvent)(unsafe.Pointer(me))
	case LeaveNotify:
		return (*XLeaveWindowEvent)(unsafe.Pointer(me))
	case FocusIn:
		return (*XFocusInEvent)(unsafe.Pointer(me))
	case FocusOut:
		return (*XFocusOutEvent)(unsafe.Pointer(me))
	case KeymapNotify:
		return (*XKeymapEvent)(unsafe.Pointer(me))
	case Expose:
		return (*XExposeEvent)(unsafe.Pointer(me))
	case GraphicsExpose:
		return (*XGraphicsExposeEvent)(unsafe.Pointer(me))
	case NoExpose:
		return (*XNoExposeEvent)(unsafe.Pointer(me))
	case VisibilityNotify:
		return (*XVisibilityEvent)(unsafe.Pointer(me))
	case CreateNotify:
		return (*XCreateWindowEvent)(unsafe.Pointer(me))
	case DestroyNotify:
		return (*XDestroyWindowEvent)(unsafe.Pointer(me))
	case UnmapNotify:
		return (*XUnmapEvent)(unsafe.Pointer(me))
	case MapNotify:
		return (*XMapEvent)(unsafe.Pointer(me))
	case MapRequest:
		return (*XMapRequestEvent)(unsafe.Pointer(me))
	case ReparentNotify:
		return (*XReparentEvent)(unsafe.Pointer(me))
	case ConfigureNotify:
		return (*XConfigureEvent)(unsafe.Pointer(me))
	case ConfigureRequest:
		return (*XConfigureRequestEvent)(unsafe.Pointer(me))
	case GravityNotify:
		return (*XGravityEvent)(unsafe.Pointer(me))
	case ResizeRequest:
		return (*XResizeRequestEvent)(unsafe.Pointer(me))
	case CirculateNotify:
		return (*XCirculateEvent)(unsafe.Pointer(me))
	case CirculateRequest:
		return (*XCirculateRequestEvent)(unsafe.Pointer(me))
	case PropertyNotify:
		return (*XPropertyEvent)(unsafe.Pointer(me))
	case SelectionClear:
		return (*XSelectionClearEvent)(unsafe.Pointer(me))
	case SelectionRequest:
		return (*XSelectionRequestEvent)(unsafe.Pointer(me))
	case SelectionNotify:
		return (*XSelectionEvent)(unsafe.Pointer(me))
	case ColormapNotify:
		return (*XColormapEvent)(unsafe.Pointer(me))
	case ClientMessage:
		return (*XClientMessageEvent)(unsafe.Pointer(me))
	case MappingNotify:
		return (*XMappingEvent)(unsafe.Pointer(me))
	case GenericEvent:
		return (*XGenericEvent)(unsafe.Pointer(me))
	case LASTEvent:
		return (*XAnyEvent)(unsafe.Pointer(me))
	default:
		return (*XAnyEvent)(unsafe.Pointer(me))
	}
}
