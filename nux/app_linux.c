// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && !android)

#include <X11/X.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xresource.h>
#include <X11/keysymdef.h>
#include <X11/extensions/Xrender.h>

#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <cairo/cairo-xlib.h>

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <strings.h>

#include "_cgo_export.h"

static const char *event_names[] = {
" no event 0 ",
" no event 1 ",
" KeyPress ",
" KeyRelease ",
" ButtonPress ",
" ButtonRelease ",
" MotionNotify ",
" EnterNotify ",
" LeaveNotify ",
" FocusIn ",
" FocusOut ",
" KeymapNotify ",
" Expose ",
" GraphicsExpose ",
" NoExpose ",
" VisibilityNotify ",
" CreateNotify ",
" DestroyNotify ",
" UnmapNotify ",
" MapNotify ",
" MapRequest ",
" ReparentNotify ",
" ConfigureNotify ",
" ConfigureRequest ",
" GravityNotify ",
" ResizeRequest ",
" CirculateNotify ",
" CirculateRequest ",
" PropertyNotify ",
" SelectionClear ",
" SelectionRequest ",
" SelectionNotify ",
" ColormapNotify ",
" ClientMessage ",
" MappingNotify ",
};

#define XALL_EVENT KeyPressMask|KeyReleaseMask|ButtonPressMask|ButtonReleaseMask|EnterWindowMask|LeaveWindowMask|PointerMotionMask|/*PointerMotionHintMask*/Button1MotionMask|Button2MotionMask|Button3MotionMask|Button4MotionMask|Button5MotionMask|ButtonMotionMask|KeymapStateMask|ExposureMask|VisibilityChangeMask|StructureNotifyMask|ResizeRedirectMask|SubstructureNotifyMask|SubstructureRedirectMask|FocusChangeMask|PropertyChangeMask|ColormapChangeMask|OwnerGrabButtonMask

#define AllMaskBits (CWBackPixmap|CWBackPixel|CWBorderPixmap|\
		     CWBorderPixel|CWBitGravity|CWWinGravity|\
		     CWBackingStore|CWBackingPlanes|CWBackingPixel|\
		     CWOverrideRedirect|CWSaveUnder|CWEventMask|\
		     CWDontPropagate|CWColormap|CWCursor)

void run(){
    Display* display;
    Visual* visual;
    Window window;
    int screen_num;

    unsigned int display_width, display_height;
    unsigned int width, height;
    char *display_name = getenv("DISPLAY");
    unsigned long valuemask = AllMaskBits;
    XSetWindowAttributes winattr = { 0 };

    display = XOpenDisplay(display_name);
    if (display == NULL) {
        printf("cannot connect to X server\n");
        exit(1);
    }

    screen_num = DefaultScreen(display);
    visual = DefaultVisual(display, screen_num);
    display_width = DisplayWidth(display, screen_num);
    display_height = DisplayHeight(display, screen_num);
    width = display_width/2;
    height = display_height/2;
    printf("display_width=%d, display_height=%d, screen_num=%d\n",display_width, display_height, screen_num);

    winattr.event_mask = XALL_EVENT;
    window = XCreateWindow(
        display, 
        RootWindow(display, screen_num),    // parent window
        0,                                  // x
        0,                                  // y
        width,                              // width
        height,                             // height
        0,                                  // border_width
        DefaultDepth(display, screen_num),  // depth
        InputOutput,                        // class
        visual,                             // visual 
        valuemask,                          // valuemask 
        &winattr                            // attributes
    );
    XMapWindow(display, window);

    /* flush all pending requests to the X server. */
    XFlush(display);

    windowCreated(window, display, visual);

    int done = 0;
    int tag = 0;
    XEvent event;
    while (!done) {

        // printf("XPending %d\n", XPending(display));
        // TODO:: set event 0
        XNextEvent(display, &event);
        switch (event.type) {
        case ConfigureNotify:
            printf("%d, %s\n", event.type, event_names[event.type]);
            windowResized(event.xconfigure.window, event.xconfigure.width, event.xconfigure.height);
        break;
        case MapNotify:
            printf("%d, %s\n", event.type, event_names[event.type]);
        break;
        case Expose:
        {
            printf("%d, %s\n", event.type, event_names[event.type]);
            windowDraw(event.xexpose.window);
            break;
        }
        case MotionNotify:
        // do not print
        break;
        case KeymapNotify:
            printf("%d, %s, %s\n", event.type, event_names[event.type], event.xkeymap.key_vector);
        break;
        case VisibilityNotify:
            printf("%d, %s, state=%d\n", event.type, event_names[event.type], event.xvisibility.state);
            break;
        case KeyPress:
            printf("%d, %s, keycode=%d, XK_q=%d\n", event.type, event_names[event.type], event.xkey.keycode, XK_q);
            if (event.xkey.keycode == 'q'){
                done = 1;
                XDestroyWindow(display, window);
            }

            XEvent event_send;
            event_send.type = Expose;
            event_send.xexpose.send_event = 1;
            event_send.xexpose.display = display;
            event_send.xexpose.window = window;

            XSendEvent(display, window, 0, ExposureMask, &event_send);
            break;
        case KeyRelease:
            printf("%d, %s, keycode=%d, XK_q=%d\n", event.type, event_names[event.type], event.xkey.keycode, XK_q);
            if (tag > 3){
                XUnmapWindow(display, window);  // hide window
            }else{
                XMapWindow(display, window);    // show window
            }
        break;
        default: /* ignore any other event types. */
            if (event.type < 35){
                printf("%d, %s\n", event.type, event_names[event.type]);
            }else{
                printf("Event NO: %d\n", event.type);
            }
            break;
        } /* end switch on event type */
    } /* end while events handling */

    XCloseDisplay(display);
}

void window_getSize(Display* display, Window window, int *width, int *height){
    XWindowAttributes attribs;
    XGetWindowAttributes(display, window, &attribs);

    if (width)
        *width = attribs.width;
    if (height)
        *height = attribs.height;
}