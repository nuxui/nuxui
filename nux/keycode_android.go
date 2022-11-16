// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android
// +build android

package nux

// https://cs.android.com/android/platform/superproject/+/master:frameworks/native/include/android/keycodes.h

const (
	Key_A            KeyCode = 29         // AKEYCODE_A
	Key_B            KeyCode = 30         // AKEYCODE_B
	Key_C            KeyCode = 31         // AKEYCODE_C
	Key_D            KeyCode = 32         // AKEYCODE_D
	Key_E            KeyCode = 33         // AKEYCODE_E
	Key_F            KeyCode = 34         // AKEYCODE_F
	Key_G            KeyCode = 35         // AKEYCODE_G
	Key_H            KeyCode = 36         // AKEYCODE_H
	Key_I            KeyCode = 37         // AKEYCODE_I
	Key_J            KeyCode = 38         // AKEYCODE_J
	Key_K            KeyCode = 39         // AKEYCODE_K
	Key_L            KeyCode = 40         // AKEYCODE_L
	Key_M            KeyCode = 41         // AKEYCODE_M
	Key_N            KeyCode = 42         // AKEYCODE_N
	Key_O            KeyCode = 43         // AKEYCODE_O
	Key_P            KeyCode = 44         // AKEYCODE_P
	Key_Q            KeyCode = 45         // AKEYCODE_Q
	Key_R            KeyCode = 46         // AKEYCODE_R
	Key_S            KeyCode = 47         // AKEYCODE_S
	Key_T            KeyCode = 48         // AKEYCODE_T
	Key_U            KeyCode = 49         // AKEYCODE_U
	Key_V            KeyCode = 50         // AKEYCODE_V
	Key_W            KeyCode = 51         // AKEYCODE_W
	Key_X            KeyCode = 52         // AKEYCODE_X
	Key_Y            KeyCode = 53         // AKEYCODE_Y
	Key_Z            KeyCode = 54         // AKEYCODE_Z
	Key_0            KeyCode = 7          // AKEYCODE_0
	Key_1            KeyCode = 8          // AKEYCODE_1
	Key_2            KeyCode = 9          // AKEYCODE_2
	Key_3            KeyCode = 10         // AKEYCODE_3
	Key_4            KeyCode = 11         // AKEYCODE_4
	Key_5            KeyCode = 12         // AKEYCODE_5
	Key_6            KeyCode = 13         // AKEYCODE_6
	Key_7            KeyCode = 14         // AKEYCODE_7
	Key_8            KeyCode = 15         // AKEYCODE_8
	Key_9            KeyCode = 16         // AKEYCODE_9
	Key_F1           KeyCode = 131        // AKEYCODE_F1
	Key_F2           KeyCode = 132        // AKEYCODE_F2
	Key_F3           KeyCode = 133        // AKEYCODE_F3
	Key_F4           KeyCode = 134        // AKEYCODE_F4
	Key_F5           KeyCode = 135        // AKEYCODE_F5
	Key_F6           KeyCode = 136        // AKEYCODE_F6
	Key_F7           KeyCode = 137        // AKEYCODE_F7
	Key_F8           KeyCode = 138        // AKEYCODE_F8
	Key_F9           KeyCode = 139        // AKEYCODE_F9
	Key_F10          KeyCode = 140        // AKEYCODE_F10
	Key_F11          KeyCode = 141        // AKEYCODE_F11
	Key_F12          KeyCode = 142        // AKEYCODE_F12
	Key_F13          KeyCode = Key_Unknow // VK_F13
	Key_F14          KeyCode = Key_Unknow // VK_F14
	Key_F15          KeyCode = Key_Unknow // VK_F15
	Key_F16          KeyCode = Key_Unknow // VK_F16
	Key_F17          KeyCode = Key_Unknow // VK_F17
	Key_F18          KeyCode = Key_Unknow // VK_F18
	Key_F19          KeyCode = Key_Unknow // VK_F19
	Key_F20          KeyCode = Key_Unknow // VK_F20
	Key_Return       KeyCode = 66         // AKEYCODE_ENTER
	Key_Tab          KeyCode = 61         // AKEYCODE_TAB
	Key_Space        KeyCode = 62         // AKEYCODE_SPACE
	Key_BackSpace    KeyCode = 67         // AKEYCODE_DEL
	Key_Escape       KeyCode = 111        // AKEYCODE_ESCAPE
	Key_CapsLock     KeyCode = 115        // AKEYCODE_CAPS_LOCK
	Key_Alt          KeyCode = 57         // AKEYCODE_ALT_LEFT ? AKEYCODE_META_LEFT
	Key_RightAlt     KeyCode = 58         // AKEYCODE_ALT_RIGHT
	Key_Shift        KeyCode = 59         // AKEYCODE_SHIFT_LEFT
	Key_RightShift   KeyCode = 60         // AKEYCODE_SHIFT_RIGHT
	Key_Control      KeyCode = 113        // AKEYCODE_CTRL_LEFT
	Key_RightControl KeyCode = 114        // AKEYCODE_CTRL_RIGHT
	Key_Command      KeyCode = Key_Unknow // VK_LWIN VK_RWIN 0x5C
	Key_Equal        KeyCode = 70         // AKEYCODE_EQUALS =+
	Key_Minus        KeyCode = 69         // AKEYCODE_MINUS -_
	Key_LeftBracket  KeyCode = 71         // AKEYCODE_LEFT_BRACKET [{
	Key_RightBracket KeyCode = 72         // AKEYCODE_RIGHT_BRACKET ]}
	Key_Quote        KeyCode = 75         // AKEYCODE_APOSTROPHE '"
	Key_Semicolon    KeyCode = 74         // AKEYCODE_SEMICOLON ;:
	Key_Comma        KeyCode = 55         // AKEYCODE_COMMA ,<
	Key_Period       KeyCode = 56         // AKEYCODE_PERIOD .>
	Key_Slash        KeyCode = 76         // AKEYCODE_SLASH /?
	Key_Backslash    KeyCode = 73         // AKEYCODE_BACKSLASH \|
	Key_Grave        KeyCode = 68         // AKEYCODE_GRAVE `~
	Key_Menu         KeyCode = 82         // AKEYCODE_MENU
	Key_Function     KeyCode = 119        // AKEYCODE_FUNCTION
	Key_Left         KeyCode = 21         // AKEYCODE_DPAD_LEFT
	Key_Right        KeyCode = 22         // AKEYCODE_DPAD_RIGHT
	Key_Up           KeyCode = 19         // AKEYCODE_DPAD_UP
	Key_Down         KeyCode = 20         // AKEYCODE_DPAD_DOWN
	Key_Mute         KeyCode = 91         // AKEYCODE_MUTE ? AKEYCODE_VOLUME_MUTE 164
	Key_VolumeUp     KeyCode = 24         // AKEYCODE_VOLUME_UP
	Key_VolumeDown   KeyCode = 25         // AKEYCODE_VOLUME_DOWN
	Key_Home         KeyCode = 122        // AKEYCODE_MOVE_HOME ? Phone Home
	Key_End          KeyCode = 123        // AKEYCODE_MOVE_END
	Key_PageUp       KeyCode = 92         // AKEYCODE_PAGE_UP
	Key_PageDown     KeyCode = 93         // AKEYCODE_PAGE_DOWN
	Key_Delete       KeyCode = 112        // AKEYCODE_FORWARD_DEL
	Key_Insert       KeyCode = 124        // AKEYCODE_INSERT
	Key_Help         KeyCode = 259        // AKEYCODE_HELP
	Key_Pad0         KeyCode = 144        // AKEYCODE_NUMPAD_0
	Key_Pad1         KeyCode = 145        // AKEYCODE_NUMPAD_1
	Key_Pad2         KeyCode = 146        // AKEYCODE_NUMPAD_2
	Key_Pad3         KeyCode = 147        // AKEYCODE_NUMPAD_3
	Key_Pad4         KeyCode = 148        // AKEYCODE_NUMPAD_4
	Key_Pad5         KeyCode = 149        // AKEYCODE_NUMPAD_5
	Key_Pad6         KeyCode = 150        // AKEYCODE_NUMPAD_6
	Key_Pad7         KeyCode = 151        // AKEYCODE_NUMPAD_7
	Key_Pad8         KeyCode = 152        // AKEYCODE_NUMPAD_8
	Key_Pad9         KeyCode = 153        // AKEYCODE_NUMPAD_9
	Key_PadDecimal   KeyCode = 158        // AKEYCODE_NUMPAD_DOT .
	Key_PadPlus      KeyCode = 157        // AKEYCODE_NUMPAD_ADD +
	Key_PadMinus     KeyCode = 156        // AKEYCODE_NUMPAD_SUBTRACT -
	Key_PadMultiply  KeyCode = 155        // AKEYCODE_NUMPAD_MULTIPLY *
	Key_PadDivide    KeyCode = 154        // AKEYCODE_NUMPAD_DIVIDE /
	Key_PadEnter     KeyCode = 160        // AKEYCODE_NUMPAD_ENTER
	Key_PadNumLock   KeyCode = 143        // AKEYCODE_NUM_LOCK
	Key_PadBegin     KeyCode = Key_Unknow // Unknow
	Key_Clear        KeyCode = 28         // AKEYCODE_CLEAR
	Key_ScrollLock   KeyCode = 116        // AKEYCODE_SCROLL_LOCK
	Key_Pause        KeyCode = 121        // AKEYCODE_BREAK
	Key_Snapshot     KeyCode = 120        // AKEYCODE_SYSRQ
	Key_PadEquals    KeyCode = 161        // AKEYCODE_NUMPAD_EQUALS
	Key_Back         KeyCode = 4          // AKEYCODE_BACK
)
