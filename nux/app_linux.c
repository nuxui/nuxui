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

int _go_nativeLoopPrepared = 0;

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
    XkbDescPtr keyboard_map = XkbGetMap(display, XkbAllClientInfoMask, XkbUseCoreKbd);  
    screen_num = DefaultScreen(display);
    visual = DefaultVisual(display, screen_num);
    display_width = DisplayWidth(display, screen_num);
    display_height = DisplayHeight(display, screen_num);
    width = 800;
    height = 600;
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

    go_windowCreated(window, display, visual);

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
            go_windowResized(event.xconfigure.window, event.xconfigure.width, event.xconfigure.height);
        break;
        case MapNotify:
            printf("%d, %s\n", event.type, event_names[event.type]);
        break;
        case Expose:
        {
            if (_go_nativeLoopPrepared == 0){
                _go_nativeLoopPrepared = 1;
                go_nativeLoopPrepared();
            }
            printf("%d, %s\n", event.type, event_names[event.type]);
            go_windowDraw(event.xexpose.window);
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
        case KeyRelease:
        {
            KeySym keysym = XkbKeycodeToKeysym(display, event.xkey.keycode, 0, 0);
            if (keysym >= XK_KP_Decimal || keysym <= XK_KP_9){
                XKeyboardState keyState;
                XGetKeyboardControl(display, &keyState);
                if( (keyState.led_mask & 2) == 2){ // NumLock On
                    keysym = XkbKeycodeToKeysym(display, event.xkey.keycode, 0, 1);
                }else{
                    switch(keysym){
                        case XK_KP_Insert    : keysym = XK_Insert    ;break;
                        case XK_KP_Up        : keysym = XK_Up        ;break;
                        case XK_KP_Down      : keysym = XK_Down      ;break;
                        case XK_KP_Left      : keysym = XK_Left      ;break;
                        case XK_KP_Right     : keysym = XK_Right     ;break;
                        case XK_KP_Home      : keysym = XK_Home      ;break;
                        case XK_KP_End       : keysym = XK_End       ;break;
                        case XK_KP_Page_Up   : keysym = XK_Page_Up   ;break;
                        case XK_KP_Page_Down : keysym = XK_Page_Down ;break;
                        case XK_KP_Delete    : keysym = XK_Delete    ;break;
                        case XK_KP_Begin     : keysym = XK_KP_Begin  ;break;
                        default:break;
                    }
                }
            }

            go_keyEvent(event.xkey.window, event.type, keysym, 0, 0, NULL);
            break;
        }
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

void invalidate(Display *display, Window window){
    XEvent event;
    event.type = Expose;
    event.xexpose.serial = 0;
    event.xexpose.send_event = True;
    event.xexpose.display = display;
    event.xexpose.window = window;
    event.xexpose.window = window;
    event.xexpose.x = 0;
    event.xexpose.y = 0;
    event.xexpose.width = 0;
    event.xexpose.height = 0;
    event.xexpose.count = 0;

    if ( XSendEvent(display, window, False, ExposureMask, &event) == 0 ){
        printf("XSendEvent faild !\n");
    }
}

void window_getSize(Display* display, Window window, int *width, int *height){
    XWindowAttributes attribs;
    XGetWindowAttributes(display, window, &attribs);

    if (width)
        *width = attribs.width;
    if (height)
        *height = attribs.height;
}

 void window_setText(Display* display, Window window, char *name)
 {
    Atom utf8Str = XInternAtom(display, "UTF8_STRING", 0);
    XChangeProperty(display, window, XA_WM_NAME, utf8Str, 8, PropModeReplace, (unsigned char *)name, (int)strlen(name));
 }