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

type ErrorCode C.uchar

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

type EventType C.int

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

type ButtonMask C.uint

const (
	Button1Mask = (1 << 8)
	Button2Mask = (1 << 9)
	Button3Mask = (1 << 10)
	Button4Mask = (1 << 11)
	Button5Mask = (1 << 12)
	AnyModifier = (1 << 15)
)

type Button C.uint

const (
	Button1 Button = 1
	Button2 Button = 2
	Button3 Button = 3
	Button4 Button = 4
	Button5 Button = 5
)

const (
	XA_WM_NAME = Atom(C.XA_WM_NAME)
)

type CursorShape uint32

const (
	XC_num_glyphs          CursorShape = 154
	XC_X_cursor            CursorShape = 0
	XC_arrow               CursorShape = 2
	XC_based_arrow_down    CursorShape = 4
	XC_based_arrow_up      CursorShape = 6
	XC_boat                CursorShape = 8
	XC_bogosity            CursorShape = 10
	XC_bottom_left_corner  CursorShape = 12
	XC_bottom_right_corner CursorShape = 14
	XC_bottom_side         CursorShape = 16
	XC_bottom_tee          CursorShape = 18
	XC_box_spiral          CursorShape = 20
	XC_center_ptr          CursorShape = 22
	XC_circle              CursorShape = 24
	XC_clock               CursorShape = 26
	XC_coffee_mug          CursorShape = 28
	XC_cross               CursorShape = 30
	XC_cross_reverse       CursorShape = 32
	XC_crosshair           CursorShape = 34
	XC_diamond_cross       CursorShape = 36
	XC_dot                 CursorShape = 38
	XC_dotbox              CursorShape = 40
	XC_double_arrow        CursorShape = 42
	XC_draft_large         CursorShape = 44
	XC_draft_small         CursorShape = 46
	XC_draped_box          CursorShape = 48
	XC_exchange            CursorShape = 50
	XC_fleur               CursorShape = 52
	XC_gobbler             CursorShape = 54
	XC_gumby               CursorShape = 56
	XC_hand1               CursorShape = 58
	XC_hand2               CursorShape = 60
	XC_heart               CursorShape = 62
	XC_icon                CursorShape = 64
	XC_iron_cross          CursorShape = 66
	XC_left_ptr            CursorShape = 68
	XC_left_side           CursorShape = 70
	XC_left_tee            CursorShape = 72
	XC_leftbutton          CursorShape = 74
	XC_ll_angle            CursorShape = 76
	XC_lr_angle            CursorShape = 78
	XC_man                 CursorShape = 80
	XC_middlebutton        CursorShape = 82
	XC_mouse               CursorShape = 84
	XC_pencil              CursorShape = 86
	XC_pirate              CursorShape = 88
	XC_plus                CursorShape = 90
	XC_question_arrow      CursorShape = 92
	XC_right_ptr           CursorShape = 94
	XC_right_side          CursorShape = 96
	XC_right_tee           CursorShape = 98
	XC_rightbutton         CursorShape = 100
	XC_rtl_logo            CursorShape = 102
	XC_sailboat            CursorShape = 104
	XC_sb_down_arrow       CursorShape = 106
	XC_sb_h_double_arrow   CursorShape = 108
	XC_sb_left_arrow       CursorShape = 110
	XC_sb_right_arrow      CursorShape = 112
	XC_sb_up_arrow         CursorShape = 114
	XC_sb_v_double_arrow   CursorShape = 116
	XC_shuttle             CursorShape = 118
	XC_sizing              CursorShape = 120
	XC_spider              CursorShape = 122
	XC_spraycan            CursorShape = 124
	XC_star                CursorShape = 126
	XC_target              CursorShape = 128
	XC_tcross              CursorShape = 130
	XC_top_left_arrow      CursorShape = 132
	XC_top_left_corner     CursorShape = 134
	XC_top_right_corner    CursorShape = 136
	XC_top_side            CursorShape = 138
	XC_top_tee             CursorShape = 140
	XC_trek                CursorShape = 142
	XC_ul_angle            CursorShape = 144
	XC_umbrella            CursorShape = 146
	XC_ur_angle            CursorShape = 148
	XC_watch               CursorShape = 150
	XC_xterm               CursorShape = 152
)

// https://code.woboq.org/qt5/include/X11/keysymdef.h.html
const (
	XK_VoidSymbol        KeySym = 0xffffff /* Void symbol */
	XK_BackSpace         KeySym = 0xff08   /* Back space, back char */
	XK_Tab               KeySym = 0xff09
	XK_Linefeed          KeySym = 0xff0a /* Linefeed, LF */
	XK_Clear             KeySym = 0xff0b
	XK_Return            KeySym = 0xff0d /* Return, enter */
	XK_Pause             KeySym = 0xff13 /* Pause, hold */
	XK_Scroll_Lock       KeySym = 0xff14
	XK_Sys_Req           KeySym = 0xff15
	XK_Escape            KeySym = 0xff1b
	XK_Delete            KeySym = 0xffff /* Delete, rubout */
	XK_Multi_key         KeySym = 0xff20 /* Multi-key character compose */
	XK_Codeinput         KeySym = 0xff37
	XK_SingleCandidate   KeySym = 0xff3c
	XK_MultipleCandidate KeySym = 0xff3d
	XK_PreviousCandidate KeySym = 0xff3e
	XK_Kanji             KeySym = 0xff21 /* Kanji, Kanji convert */
	XK_Muhenkan          KeySym = 0xff22 /* Cancel Conversion */
	XK_Henkan_Mode       KeySym = 0xff23 /* Start/Stop Conversion */
	XK_Henkan            KeySym = 0xff23 /* Alias for Henkan_Mode */
	XK_Romaji            KeySym = 0xff24 /* to Romaji */
	XK_Hiragana          KeySym = 0xff25 /* to Hiragana */
	XK_Katakana          KeySym = 0xff26 /* to Katakana */
	XK_Hiragana_Katakana KeySym = 0xff27 /* Hiragana/Katakana toggle */
	XK_Zenkaku           KeySym = 0xff28 /* to Zenkaku */
	XK_Hankaku           KeySym = 0xff29 /* to Hankaku */
	XK_Zenkaku_Hankaku   KeySym = 0xff2a /* Zenkaku/Hankaku toggle */
	XK_Touroku           KeySym = 0xff2b /* Add to Dictionary */
	XK_Massyo            KeySym = 0xff2c /* Delete from Dictionary */
	XK_Kana_Lock         KeySym = 0xff2d /* Kana Lock */
	XK_Kana_Shift        KeySym = 0xff2e /* Kana Shift */
	XK_Eisu_Shift        KeySym = 0xff2f /* Alphanumeric Shift */
	XK_Eisu_toggle       KeySym = 0xff30 /* Alphanumeric toggle */
	XK_Kanji_Bangou      KeySym = 0xff37 /* Codeinput */
	XK_Zen_Koho          KeySym = 0xff3d /* Multiple/All Candidate(s) */
	XK_Mae_Koho          KeySym = 0xff3e /* Previous Candidate */
	XK_Home              KeySym = 0xff50
	XK_Left              KeySym = 0xff51 /* Move left, left arrow */
	XK_Up                KeySym = 0xff52 /* Move up, up arrow */
	XK_Right             KeySym = 0xff53 /* Move right, right arrow */
	XK_Down              KeySym = 0xff54 /* Move down, down arrow */
	XK_Prior             KeySym = 0xff55 /* Prior, previous */
	XK_Page_Up           KeySym = 0xff55
	XK_Next              KeySym = 0xff56 /* Next */
	XK_Page_Down         KeySym = 0xff56
	XK_End               KeySym = 0xff57 /* EOL */
	XK_Begin             KeySym = 0xff58 /* BOL */
	XK_Select            KeySym = 0xff60 /* Select, mark */
	XK_Print             KeySym = 0xff61
	XK_Execute           KeySym = 0xff62 /* Execute, run, do */
	XK_Insert            KeySym = 0xff63 /* Insert, insert here */
	XK_Undo              KeySym = 0xff65
	XK_Redo              KeySym = 0xff66 /* Redo, again */
	XK_Menu              KeySym = 0xff67
	XK_Find              KeySym = 0xff68 /* Find, search */
	XK_Cancel            KeySym = 0xff69 /* Cancel, stop, abort, exit */
	XK_Help              KeySym = 0xff6a /* Help */
	XK_Break             KeySym = 0xff6b
	XK_Mode_switch       KeySym = 0xff7e /* Character set switch */
	XK_script_switch     KeySym = 0xff7e /* Alias for mode_switch */
	XK_Num_Lock          KeySym = 0xff7f
	XK_KP_Space          KeySym = 0xff80 /* Space */
	XK_KP_Tab            KeySym = 0xff89
	XK_KP_Enter          KeySym = 0xff8d /* Enter */
	XK_KP_F1             KeySym = 0xff91 /* PF1, KP_A, ... */
	XK_KP_F2             KeySym = 0xff92
	XK_KP_F3             KeySym = 0xff93
	XK_KP_F4             KeySym = 0xff94
	XK_KP_Home           KeySym = 0xff95
	XK_KP_Left           KeySym = 0xff96
	XK_KP_Up             KeySym = 0xff97
	XK_KP_Right          KeySym = 0xff98
	XK_KP_Down           KeySym = 0xff99
	XK_KP_Prior          KeySym = 0xff9a
	XK_KP_Page_Up        KeySym = 0xff9a
	XK_KP_Next           KeySym = 0xff9b
	XK_KP_Page_Down      KeySym = 0xff9b
	XK_KP_End            KeySym = 0xff9c
	XK_KP_Begin          KeySym = 0xff9d
	XK_KP_Insert         KeySym = 0xff9e
	XK_KP_Delete         KeySym = 0xff9f
	XK_KP_Equal          KeySym = 0xffbd /* Equals */
	XK_KP_Multiply       KeySym = 0xffaa
	XK_KP_Add            KeySym = 0xffab
	XK_KP_Separator      KeySym = 0xffac /* Separator, often comma */
	XK_KP_Subtract       KeySym = 0xffad
	XK_KP_Decimal        KeySym = 0xffae
	XK_KP_Divide         KeySym = 0xffaf
	XK_KP_0              KeySym = 0xffb0
	XK_KP_1              KeySym = 0xffb1
	XK_KP_2              KeySym = 0xffb2
	XK_KP_3              KeySym = 0xffb3
	XK_KP_4              KeySym = 0xffb4
	XK_KP_5              KeySym = 0xffb5
	XK_KP_6              KeySym = 0xffb6
	XK_KP_7              KeySym = 0xffb7
	XK_KP_8              KeySym = 0xffb8
	XK_KP_9              KeySym = 0xffb9
	XK_F1                KeySym = 0xffbe
	XK_F2                KeySym = 0xffbf
	XK_F3                KeySym = 0xffc0
	XK_F4                KeySym = 0xffc1
	XK_F5                KeySym = 0xffc2
	XK_F6                KeySym = 0xffc3
	XK_F7                KeySym = 0xffc4
	XK_F8                KeySym = 0xffc5
	XK_F9                KeySym = 0xffc6
	XK_F10               KeySym = 0xffc7
	XK_F11               KeySym = 0xffc8
	XK_L1                KeySym = 0xffc8
	XK_F12               KeySym = 0xffc9
	XK_L2                KeySym = 0xffc9
	XK_F13               KeySym = 0xffca
	XK_L3                KeySym = 0xffca
	XK_F14               KeySym = 0xffcb
	XK_L4                KeySym = 0xffcb
	XK_F15               KeySym = 0xffcc
	XK_L5                KeySym = 0xffcc
	XK_F16               KeySym = 0xffcd
	XK_L6                KeySym = 0xffcd
	XK_F17               KeySym = 0xffce
	XK_L7                KeySym = 0xffce
	XK_F18               KeySym = 0xffcf
	XK_L8                KeySym = 0xffcf
	XK_F19               KeySym = 0xffd0
	XK_L9                KeySym = 0xffd0
	XK_F20               KeySym = 0xffd1
	XK_L10               KeySym = 0xffd1
	XK_F21               KeySym = 0xffd2
	XK_R1                KeySym = 0xffd2
	XK_F22               KeySym = 0xffd3
	XK_R2                KeySym = 0xffd3
	XK_F23               KeySym = 0xffd4
	XK_R3                KeySym = 0xffd4
	XK_F24               KeySym = 0xffd5
	XK_R4                KeySym = 0xffd5
	XK_F25               KeySym = 0xffd6
	XK_R5                KeySym = 0xffd6
	XK_F26               KeySym = 0xffd7
	XK_R6                KeySym = 0xffd7
	XK_F27               KeySym = 0xffd8
	XK_R7                KeySym = 0xffd8
	XK_F28               KeySym = 0xffd9
	XK_R8                KeySym = 0xffd9
	XK_F29               KeySym = 0xffda
	XK_R9                KeySym = 0xffda
	XK_F30               KeySym = 0xffdb
	XK_R10               KeySym = 0xffdb
	XK_F31               KeySym = 0xffdc
	XK_R11               KeySym = 0xffdc
	XK_F32               KeySym = 0xffdd
	XK_R12               KeySym = 0xffdd
	XK_F33               KeySym = 0xffde
	XK_R13               KeySym = 0xffde
	XK_F34               KeySym = 0xffdf
	XK_R14               KeySym = 0xffdf
	XK_F35               KeySym = 0xffe0
	XK_R15               KeySym = 0xffe0
	XK_Shift_L           KeySym = 0xffe1 /* Left shift */
	XK_Shift_R           KeySym = 0xffe2 /* Right shift */
	XK_Control_L         KeySym = 0xffe3 /* Left control */
	XK_Control_R         KeySym = 0xffe4 /* Right control */
	XK_Caps_Lock         KeySym = 0xffe5 /* Caps lock */
	XK_Shift_Lock        KeySym = 0xffe6 /* Shift lock */
	XK_Meta_L            KeySym = 0xffe7 /* Left meta */
	XK_Meta_R            KeySym = 0xffe8 /* Right meta */
	XK_Alt_L             KeySym = 0xffe9 /* Left alt */
	XK_Alt_R             KeySym = 0xffea /* Right alt */
	XK_Super_L           KeySym = 0xffeb /* Left super */
	XK_Super_R           KeySym = 0xffec /* Right super */
	XK_Hyper_L           KeySym = 0xffed /* Left hyper */
	XK_Hyper_R           KeySym = 0xffee /* Right hyper */
)
