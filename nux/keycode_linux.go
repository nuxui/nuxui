// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

// https://code.woboq.org/qt5/include/X11/keysymdef.h.html

const (
	Key_A            KeyCode = 0x61       // XK_a
	Key_B            KeyCode = 0x62       // XK_b
	Key_C            KeyCode = 0x63       // XK_c
	Key_D            KeyCode = 0x64       // XK_d
	Key_E            KeyCode = 0x65       // XK_e
	Key_F            KeyCode = 0x66       // XK_f
	Key_G            KeyCode = 0x67       // XK_g
	Key_H            KeyCode = 0x68       // XK_h
	Key_I            KeyCode = 0x69       // XK_i
	Key_J            KeyCode = 0x6a       // XK_j
	Key_K            KeyCode = 0x6b       // XK_k
	Key_L            KeyCode = 0x6c       // XK_l
	Key_M            KeyCode = 0x6d       // XK_m
	Key_N            KeyCode = 0x6e       // XK_n
	Key_O            KeyCode = 0x6f       // XK_o
	Key_P            KeyCode = 0x70       // XK_p
	Key_Q            KeyCode = 0x71       // XK_q
	Key_R            KeyCode = 0x72       // XK_r
	Key_S            KeyCode = 0x73       // XK_s
	Key_T            KeyCode = 0x74       // XK_t
	Key_U            KeyCode = 0x75       // XK_u
	Key_V            KeyCode = 0x76       // XK_v
	Key_W            KeyCode = 0x77       // XK_w
	Key_X            KeyCode = 0x78       // XK_x
	Key_Y            KeyCode = 0x79       // XK_y
	Key_Z            KeyCode = 0x7a       // XK_z
	Key_0            KeyCode = 0x30       // XK_0
	Key_1            KeyCode = 0x31       // XK_1
	Key_2            KeyCode = 0x32       // XK_2
	Key_3            KeyCode = 0x33       // XK_3
	Key_4            KeyCode = 0x34       // XK_4
	Key_5            KeyCode = 0x35       // XK_5
	Key_6            KeyCode = 0x36       // XK_6
	Key_7            KeyCode = 0x37       // XK_7
	Key_8            KeyCode = 0x38       // XK_8
	Key_9            KeyCode = 0x39       // XK_9
	Key_F1           KeyCode = 0xffbe     // XK_F1
	Key_F2           KeyCode = 0xffbf     // XK_F2
	Key_F3           KeyCode = 0xffc0     // XK_F3
	Key_F4           KeyCode = 0xffc1     // XK_F4
	Key_F5           KeyCode = 0xffc2     // XK_F5
	Key_F6           KeyCode = 0xffc3     // XK_F6
	Key_F7           KeyCode = 0xffc4     // XK_F7
	Key_F8           KeyCode = 0xffc5     // XK_F8
	Key_F9           KeyCode = 0xffc6     // XK_F9
	Key_F10          KeyCode = 0xffc7     // XK_F10
	Key_F11          KeyCode = 0xffc8     // XK_F11
	Key_F12          KeyCode = 0xffc9     // XK_F12
	Key_F13          KeyCode = 0xffca     // XK_F13
	Key_F14          KeyCode = 0xffcb     // XK_F14
	Key_F15          KeyCode = 0xffcc     // XK_F15
	Key_F16          KeyCode = 0xffcd     // XK_F16
	Key_F17          KeyCode = 0xffce     // XK_F17
	Key_F18          KeyCode = 0xffcf     // XK_F18
	Key_F19          KeyCode = 0xffd0     // XK_F19
	Key_F20          KeyCode = 0xffd1     // XK_F20
	Key_Return       KeyCode = 0xff0d     // XK_Return
	Key_Tab          KeyCode = 0xff09     // XK_Tab
	Key_Space        KeyCode = 0x0020     // XK_space
	Key_BackSpace    KeyCode = 0xff08     // XK_BackSpace
	Key_Escape       KeyCode = 0xff1b     // XK_Escape
	Key_CapsLock     KeyCode = 0xffe5     // XK_Caps_Lock
	Key_Alt          KeyCode = 0xffe9     // XK_Alt_L
	Key_RightAlt     KeyCode = 0xffea     // XK_Alt_R
	Key_Shift        KeyCode = 0xffe1     // XK_Shift_L
	Key_RightShift   KeyCode = 0xffe2     // XK_Shift_R
	Key_Control      KeyCode = 0xffe3     // XK_Control_L
	Key_RightControl KeyCode = 0xffe4     // XK_Control_R
	Key_Command      KeyCode = Key_Unknow // unknow
	Key_Equal        KeyCode = 0x003d     // XK_equal =+
	Key_Minus        KeyCode = 0x002d     // XK_minus -_
	Key_LeftBracket  KeyCode = 0x005b     // XK_bracketleft [{
	Key_RightBracket KeyCode = 0x005d     // XK_bracketright ]}
	Key_Quote        KeyCode = 0x0027     // XK_apostrophe '"
	Key_Semicolon    KeyCode = 0x003b     // XK_semicolon ;:
	Key_Comma        KeyCode = 0x002c     // XK_comma ,<
	Key_Period       KeyCode = 0x002e     // XK_period .>
	Key_Slash        KeyCode = 0x002f     // XK_slash /?
	Key_Backslash    KeyCode = 0x005c     // XK_backslash \|
	Key_Grave        KeyCode = 0x0060     // XK_grave `~
	Key_Menu         KeyCode = 0xff67     // XK_Menu
	Key_Function     KeyCode = Key_Unknow // Unknow
	Key_Left         KeyCode = 0xff51     // XK_Left
	Key_Right        KeyCode = 0xff53     // XK_Right
	Key_Up           KeyCode = 0xff52     // XK_Up
	Key_Down         KeyCode = 0xff54     // XK_Down
	Key_Mute         KeyCode = Key_Unknow // Unknow
	Key_VolumeUp     KeyCode = Key_Unknow // Unknow
	Key_VolumeDown   KeyCode = Key_Unknow // Unknow
	Key_Home         KeyCode = 0xff50     // XK_Home
	Key_End          KeyCode = 0xff57     // XK_End
	Key_PageUp       KeyCode = 0xff55     // XK_Page_Up
	Key_PageDown     KeyCode = 0xff56     // XK_Page_Down
	Key_Delete       KeyCode = 0xffff     // XK_Delete
	Key_Insert       KeyCode = 0xff63     // XK_Insert
	Key_Help         KeyCode = Key_Unknow // Unknow
	Key_Pad0         KeyCode = 0xffb0     // XK_KP_0
	Key_Pad1         KeyCode = 0xffb1     // XK_KP_1
	Key_Pad2         KeyCode = 0xffb2     // XK_KP_2
	Key_Pad3         KeyCode = 0xffb3     // XK_KP_3
	Key_Pad4         KeyCode = 0xffb4     // XK_KP_4
	Key_Pad5         KeyCode = 0xffb5     // XK_KP_5
	Key_Pad6         KeyCode = 0xffb6     // XK_KP_6
	Key_Pad7         KeyCode = 0xffb7     // XK_KP_7
	Key_Pad8         KeyCode = 0xffb8     // XK_KP_8
	Key_Pad9         KeyCode = 0xffb9     // XK_KP_9
	Key_PadDecimal   KeyCode = 0xffae     // XK_KP_Decimal
	Key_PadPlus      KeyCode = 0xffab     // XK_KP_Add +
	Key_PadMinus     KeyCode = 0xffad     // XK_KP_Subtract -
	Key_PadMultiply  KeyCode = 0xffaa     // XK_KP_Multiply *
	Key_PadDivide    KeyCode = 0xffaf     // XK_KP_Divide /
	Key_PadEnter     KeyCode = 0xff8d     // XK_KP_Enter
	Key_PadNumLock   KeyCode = 0xff7f     // XK_Num_Lock
	Key_PadBegin     KeyCode = 0xff9d     // XK_KP_Begin
	Key_Clear        KeyCode = 0xff0b     // XK_Clear
	Key_ScrollLock   KeyCode = 0xff14     // XK_Scroll_Lock
	Key_Pause        KeyCode = 0xff13     // XK_Pause
	Key_Snapshot     KeyCode = 0xff15     // XK_Sys_Req
	Key_PadEquals    KeyCode = Key_Unknow // Unknow
	Key_Back         KeyCode = Key_Unknow // AKEYCODE_BACK

)
