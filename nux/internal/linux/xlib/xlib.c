// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go:build (linux && !android)

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

static const char *event_names[] = {
    " no event 0 ",       " no event 1 ",      " KeyPress ",
    " KeyRelease ",       " ButtonPress ",     " ButtonRelease ",
    " MotionNotify ",     " EnterNotify ",     " LeaveNotify ",
    " FocusIn ",          " FocusOut ",        " KeymapNotify ",
    " Expose ",           " GraphicsExpose ",  " NoExpose ",
    " VisibilityNotify ", " CreateNotify ",    " DestroyNotify ",
    " UnmapNotify ",      " MapNotify ",       " MapRequest ",
    " ReparentNotify ",   " ConfigureNotify ", " ConfigureRequest ",
    " GravityNotify ",    " ResizeRequest ",   " CirculateNotify ",
    " CirculateRequest ", " PropertyNotify ",  " SelectionClear ",
    " SelectionRequest ", " SelectionNotify ", " ColormapNotify ",
    " ClientMessage ",    " MappingNotify ",
};

// https://www.x.org/releases/current/doc/libX11/libX11/libX11.html#:~:text=Preedit%20State-,Callbacks,-When%20the%20input
int nux_PreeditStartCallback(XIC xic, XPointer client_data, XPointer call_data){
    printf("nux_PreeditStartCallback \n" );

    return -1;
}
void nux_PreeditDoneCallback(XIC xic, XPointer client_data, XPointer call_data){
    printf("nux_PreeditDoneCallback\n");

}
void nux_PreeditDrawCallback(XIC xic, XPointer client_data, XIMPreeditDrawCallbackStruct *call_data){
    if(call_data->text != NULL){
        // go_typingEvent(1, call_data->text->string.multi_byte, call_data->text->length, call_data->caret);
       printf("nux_PreeditDrawCallback\n");
    }
}
void nux_PreeditCaretCallback(XIC xic, XPointer client_data, XIMPreeditCaretCallbackStruct *call_data){
    printf("nux_PreeditCaretCallback\n");

}

void nux_PreeditStateNotifyCallback(XIC xic, XPointer client_data, XIMPreeditStateNotifyCallbackStruct *call_data){
    printf("nux_PreeditStateNotifyCallback\n");

}

void nux_StatusStartCallback(XIC ic, XPointer client_data, XPointer call_data){
    printf("nux_StatusStartCallback\n");

}

void nux_StatusDoneCallback(XIC ic, XPointer client_data, XPointer call_data){
    printf("nux_StatusDoneCallback\n");

}

void nux_StatusDrawCallback(XIC ic, XPointer client_data, XIMStatusDrawCallbackStruct *call_data){
    printf("nux_StatusDrawCallback\n");

}

XIC nux_XCreateIC(XIM xim, Window window){
    XIMCallback start_callback;
    start_callback.client_data = NULL;
    start_callback.callback = (XIMProc)nux_PreeditStartCallback;
    XIMCallback done_callback;
    done_callback.client_data = NULL;
    done_callback.callback = (XIMProc)nux_PreeditDoneCallback;
    XIMCallback draw_callback;
    draw_callback.client_data = NULL;
    draw_callback.callback = (XIMProc)nux_PreeditDrawCallback;
    XIMCallback caret_callback;
    caret_callback.client_data = NULL;
    caret_callback.callback = (XIMProc)nux_PreeditCaretCallback;
    XIMCallback state_notify_callback;
    state_notify_callback.client_data = NULL;
    state_notify_callback.callback = (XIMProc)nux_PreeditStateNotifyCallback;

    XVaNestedList preedit_attributes = XVaCreateNestedList(
        0,
        XNPreeditStartCallback, &start_callback,
        XNPreeditDoneCallback, &done_callback,
        XNPreeditDrawCallback, &draw_callback,
        XNPreeditCaretCallback, &caret_callback,
        XNPreeditStateNotifyCallback, &state_notify_callback,
        NULL);

    // TODO:: free preedit_attributes

    XVaNestedList status_attributes = XVaCreateNestedList(
        0,
        XNStatusStartCallback, &nux_StatusStartCallback,
        XNStatusDoneCallback, &nux_StatusDoneCallback,
        XNStatusDrawCallback, &nux_StatusDrawCallback,
        NULL);

  return XCreateIC(xim,
                    /* the following are in attr, val format, terminated by NULL */
                    XNInputStyle, XIMPreeditCallbacks | XIMStatusCallbacks,
                    XNClientWindow, window, 
                    XNPreeditAttributes, preedit_attributes,
                    XNStatusAttributes, status_attributes,
                    NULL);
}