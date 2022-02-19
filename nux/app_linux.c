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
#include <locale.h>

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

#define Button6            6
#define Button7            7

#define XALL_EVENT KeyPressMask|KeyReleaseMask|ButtonPressMask|ButtonReleaseMask|EnterWindowMask|LeaveWindowMask|PointerMotionMask|/*PointerMotionHintMask*/Button1MotionMask|Button2MotionMask|Button3MotionMask|Button4MotionMask|Button5MotionMask|ButtonMotionMask|KeymapStateMask|ExposureMask|VisibilityChangeMask|StructureNotifyMask|ResizeRedirectMask|SubstructureNotifyMask|SubstructureRedirectMask|FocusChangeMask|PropertyChangeMask|ColormapChangeMask|OwnerGrabButtonMask

#define AllMaskBits (CWBackPixmap|CWBackPixel|CWBorderPixmap|\
		     CWBorderPixel|CWBitGravity|CWWinGravity|\
		     CWBackingStore|CWBackingPlanes|CWBackingPixel|\
		     CWOverrideRedirect|CWSaveUnder|CWEventMask|\
		     CWDontPropagate|CWColormap|CWCursor)

int _go_nativeLoopPrepared = 0;
static pid_t mainThreadId;
static short g_ime_pos_x ,g_ime_pos_y;

int isMainThread(){
    return mainThreadId == gettid();
}

// https://www.x.org/releases/current/doc/libX11/libX11/libX11.html#:~:text=Preedit%20State-,Callbacks,-When%20the%20input
int nux_PreeditStartCallback(XIC xic, XPointer client_data, XPointer call_data){
    printf("nux_PreeditStartCallback x=%d, y=%d\n", g_ime_pos_x, g_ime_pos_y);

    return -1;
}
void nux_PreeditDoneCallback(XIC xic, XPointer client_data, XPointer call_data){
    printf("nux_PreeditDoneCallback\n");

}
void nux_PreeditDrawCallback(XIC xic, XPointer client_data, XIMPreeditDrawCallbackStruct *call_data){
    if(call_data->text != NULL){
        go_typeEvent(1, call_data->text->string.multi_byte, call_data->text->length, call_data->caret);
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

void run(){
    Display* display;
    Visual* visual;
    Window window;
    int screen;
    unsigned int display_width, display_height;
    unsigned int width, height;
    unsigned long valuemask = AllMaskBits;
    XSetWindowAttributes winattr = { 0 };
    XIM xim;
    XIC xic;

    mainThreadId = gettid();

    // HACK: If the application has left the locale as "C" then both wide
    //       character text input and explicit UTF-8 input via XIM will break
    //       This sets the CTYPE part of the current locale from the environment
    //       in the hope that it is set to something more sane than "C"
    if ( strcmp(setlocale(LC_CTYPE, NULL), "C") == 0 ){
        setlocale(LC_CTYPE, "");
    }

    XInitThreads();
    XrmInitialize();

    display = XOpenDisplay(NULL);
    if (!display)
    {
        const char* display_name = getenv("DISPLAY");
        if (display_name){
            printf("X11: Failed to open display %s\n", display_name);
        } else {
            printf("X11: The DISPLAY environment variable is missing\n");
        }
        return;
    }

    // XkbDescPtr keyboard_map = XkbGetMap(display, XkbAllClientInfoMask, XkbUseCoreKbd);  
    screen = DefaultScreen(display);
    visual = DefaultVisual(display, screen);
    display_width = DisplayWidth(display, screen);
    display_height = DisplayHeight(display, screen);
    width = 800;
    height = 600;
    printf("display_width=%d, display_height=%d, screen=%d\n",display_width, display_height, screen);

    if (XSupportsLocale())
    {
        XSetLocaleModifiers("");
    }

    winattr.event_mask = XALL_EVENT;
    window = XCreateWindow(
        display, 
        RootWindow(display, screen),    // parent window
        0,                                  // x
        0,                                  // y
        width,                              // width
        height,                             // height
        0,                                  // border_width
        DefaultDepth(display, screen),  // depth
        InputOutput,                        // class
        visual,                             // visual 
        valuemask,                          // valuemask 
        &winattr                            // attributes
    );
    XMapWindow(display, window);

    xim = XOpenIM(display, 0, NULL, NULL);

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
    // TODO:: free status_attributes
    xic = XCreateIC(xim,
                    /* the following are in attr, val format, terminated by NULL */
                    XNInputStyle, XIMPreeditCallbacks | XIMStatusCallbacks,
                    XNClientWindow, window, 
                    XNPreeditAttributes, preedit_attributes,
                    XNStatusAttributes, status_attributes,
                    NULL);
    /* focus on the only IC */
    XSetICFocus(xic);


    /* flush all pending requests to the X server. */
    XFlush(display);

    go_windowCreated(window, display, visual);

    int done = 0;
    int tag = 0;
    XEvent event;
    Bool filtered;
    while (!done) {
        // printf("XPending %d\n", XPending(display));
        // TODO:: set event 0
        XNextEvent(display, &event);
        filtered = XFilterEvent(&event, None); // filter by input method
        if (event.type < 35){
            printf("%d, %s, filtered=%d, serial: %ld\n", event.type, event_names[event.type], filtered, event.xany.serial);
        }else{
            printf("Event NO: %d, serial: %ld\n", event.type, event.xany.serial);
        }
        if(filtered){
            continue;
        }

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
            go_mouseEvent(event.xmotion.window, event.xmotion.type, (float)event.xmotion.x, (float)event.xmotion.y, event.xmotion.is_hint);
        break;
        case ButtonPress:
        {
            printf("ButtonPress x=%d, y=%d, button=%d, state=%d\n", event.xbutton.x, event.xbutton.y, event.xbutton.button, event.xbutton.state);
            switch(event.xbutton.button){
            case Button4:
                go_scrollEvent(event.xbutton.window, (float)event.xbutton.x, (float)event.xbutton.y, 0.0, 1.0);break;
            case Button5:
                go_scrollEvent(event.xbutton.window, (float)event.xbutton.x, (float)event.xbutton.y, 0.0, -1.0);break;
            case Button6:
                go_scrollEvent(event.xbutton.window, (float)event.xbutton.x, (float)event.xbutton.y, 1.0, 0.0);break;
            case Button7:
                go_scrollEvent(event.xbutton.window, (float)event.xbutton.x, (float)event.xbutton.y, -1.0, 0.0);break;
            default:
                go_mouseEvent(event.xbutton.window, event.xbutton.type, event.xbutton.x, event.xbutton.y, event.xbutton.button);
                break;
            }
            break;
        }
        case ButtonRelease:
        {
            printf("ButtonRelease x=%d, y=%d, button=%d, state=%d\n", event.xbutton.x, event.xbutton.y, event.xbutton.button, event.xbutton.state);

            switch(event.xbutton.button){
            case Button4:break;
            case Button5:break;
            case Button6:break;
            case Button7:break;
            default:
                go_mouseEvent(event.xbutton.window, event.xbutton.type, event.xbutton.x, event.xbutton.y, event.xbutton.button);
                break;
            }
            break;
        }
        case KeymapNotify:
            printf("%d, %s, %s\n", event.type, event_names[event.type], event.xkeymap.key_vector);
        break;
        case VisibilityNotify:
            printf("%d, %s, state=%d\n", event.type, event_names[event.type], event.xvisibility.state);
            break;
        case KeyPress:
        case KeyRelease:
        {
            {
                XPoint spot;
                spot.x = 500;//g_ime_pos_x;
                spot.y = 500; //g_ime_pos_y;
                XVaNestedList preedit_attr;
                preedit_attr = XVaCreateNestedList(0, XNSpotLocation, &spot, NULL);
                XSetICValues(xic, XNPreeditAttributes, preedit_attr, NULL);
                XFree(preedit_attr);
            }
            // 1. handle key event
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
            }

            // 2. handle input event
            if(event.type == KeyPress){
                KeySym keysym;
                Status status;
                int len = Xutf8LookupString(xic, &event.xkey, NULL, 0, &keysym, &status);
                char* text = (char*)calloc(len+1, sizeof(char));
                len = Xutf8LookupString(xic, &event.xkey, text, len, &keysym, &status);
                if (status == XLookupChars){
                    go_typeEvent(0, text, len, 0);
                }else if (status == XLookupBoth){
                    if (keysym >= 0x20 && keysym <= 0x7E){
                        go_typeEvent(0, text, len, 0);
                    }
                }
                free(text);
            }
            break;
        }
        case ClientMessage:
        {
            printf("ClientMessage atom: %s\n", XGetAtomName(display, event.xclient.message_type));
            if (strcmp("_user_ev",XGetAtomName(display, event.xclient.message_type)) == 0){
                go_backToUI();
            }
        }
        default: /* ignore any other event types. */
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
    event.xexpose.x = 0;
    event.xexpose.y = 0;
    event.xexpose.width = 0;
    event.xexpose.height = 0;
    event.xexpose.count = 0;

    if ( XSendEvent(display, window, False, ExposureMask, &event) == 0 ){
        printf("XSendEvent faild !\n");
    }
}

void runOnUI(Display *display, Window window){
    XEvent event;
    event.xclient.type = ClientMessage;
    event.xclient.serial = 0;
    event.xclient.send_event = True;
    event.xclient.display = display;
    event.xclient.window = window;
    event.xclient.message_type = XInternAtom(display, "_user_ev", False);
    event.xclient.format = 32;
    event.xclient.data.l[0] = XInternAtom(display, "_user_ev_runOnUI", False);
    XSendEvent(display, window, False, NoEventMask, &event);
}

void window_getSize(Display* display, Window window, int *width, int *height){
    XWindowAttributes attribs;
    XGetWindowAttributes(display, window, &attribs);
    if (width) {*width = attribs.width;};
    if (height) {*height = attribs.height;};
}

// TODO:: content size
void window_getContentSize(Display* display, Window window, int *width, int *height){
    XWindowAttributes attribs;
    XGetWindowAttributes(display, window, &attribs);
    if (width) {*width = attribs.width;};
    if (height) {*height = attribs.height;};
}

void window_setTitle(Display* display, Window window, char *name){
    Atom utf8Str = XInternAtom(display, "UTF8_STRING", 0);
    XChangeProperty(display, window, XA_WM_NAME, utf8Str, 8, PropModeReplace, (unsigned char *)name, (int)strlen(name));
}

void setTextInputRect(short x, short y){
    g_ime_pos_x = x;
    g_ime_pos_y = y;
}