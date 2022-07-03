// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

// https://github.com/tpn/winsdk-10/blob/master/Include/10.0.10240.0/um/WinUser.h

const (
	Key_A            KeyCode = 0x41       // VK_A
	Key_B            KeyCode = 0x42       // VK_B
	Key_C            KeyCode = 0x43       // VK_C
	Key_D            KeyCode = 0x44       // VK_D
	Key_E            KeyCode = 0x45       // VK_E
	Key_F            KeyCode = 0x46       // VK_F
	Key_G            KeyCode = 0x47       // VK_G
	Key_H            KeyCode = 0x48       // VK_H
	Key_I            KeyCode = 0x49       // VK_I
	Key_J            KeyCode = 0x4A       // VK_J
	Key_K            KeyCode = 0x4B       // VK_K
	Key_L            KeyCode = 0x4C       // VK_L
	Key_M            KeyCode = 0x4D       // VK_M
	Key_N            KeyCode = 0x4E       // VK_N
	Key_O            KeyCode = 0x4F       // VK_O
	Key_P            KeyCode = 0x50       // VK_P
	Key_Q            KeyCode = 0x51       // VK_Q
	Key_R            KeyCode = 0x52       // VK_R
	Key_S            KeyCode = 0x53       // VK_S
	Key_T            KeyCode = 0x54       // VK_T
	Key_U            KeyCode = 0x55       // VK_U
	Key_V            KeyCode = 0x56       // VK_V
	Key_W            KeyCode = 0x57       // VK_W
	Key_X            KeyCode = 0x58       // VK_X
	Key_Y            KeyCode = 0x59       // VK_Y
	Key_Z            KeyCode = 0x5A       // VK_Z
	Key_0            KeyCode = 0x30       // VK_0
	Key_1            KeyCode = 0x31       // VK_1
	Key_2            KeyCode = 0x32       // VK_2
	Key_3            KeyCode = 0x33       // VK_3
	Key_4            KeyCode = 0x34       // VK_4
	Key_5            KeyCode = 0x35       // VK_5
	Key_6            KeyCode = 0x36       // VK_6
	Key_7            KeyCode = 0x37       // VK_7
	Key_8            KeyCode = 0x38       // VK_8
	Key_9            KeyCode = 0x39       // VK_9
	Key_F1           KeyCode = 0x70       // VK_F1
	Key_F2           KeyCode = 0x71       // VK_F2
	Key_F3           KeyCode = 0x72       // VK_F3
	Key_F4           KeyCode = 0x73       // VK_F4
	Key_F5           KeyCode = 0x74       // VK_F5
	Key_F6           KeyCode = 0x75       // VK_F6
	Key_F7           KeyCode = 0x76       // VK_F7
	Key_F8           KeyCode = 0x77       // VK_F8
	Key_F9           KeyCode = 0x78       // VK_F9
	Key_F10          KeyCode = 0x79       // VK_F10
	Key_F11          KeyCode = 0x7A       // VK_F11
	Key_F12          KeyCode = 0x7B       // VK_F12
	Key_F13          KeyCode = 0x7C       // VK_F13
	Key_F14          KeyCode = 0x7D       // VK_F14
	Key_F15          KeyCode = 0x7E       // VK_F15
	Key_F16          KeyCode = 0x7F       // VK_F16
	Key_F17          KeyCode = 0x80       // VK_F17
	Key_F18          KeyCode = 0x81       // VK_F18
	Key_F19          KeyCode = 0x82       // VK_F19
	Key_F20          KeyCode = 0x83       // VK_F20
	Key_Return       KeyCode = 0x0D       // VK_RETURN
	Key_Tab          KeyCode = 0x09       // VK_TAB
	Key_Space        KeyCode = 0x20       // VK_SPACE
	Key_BackSpace    KeyCode = 0x08       // VK_BACK
	Key_Escape       KeyCode = 0x1B       // VK_ESCAPE
	Key_CapsLock     KeyCode = 0x14       // VK_CAPITAL
	Key_Alt          KeyCode = 0x12       // VK_MENU
	Key_RightAlt     KeyCode = 0x12       // VK_MENU
	Key_Shift        KeyCode = 0x10       // VK_SHIFT
	Key_RightShift   KeyCode = 0x10       // VK_SHIFT
	Key_Control      KeyCode = 0x11       // VK_CONTROL
	Key_RightControl KeyCode = 0x11       // VK_CONTROL
	Key_Command      KeyCode = 0x5B       // VK_LWIN VK_RWIN 0x5C
	Key_Equal        KeyCode = 0xBB       // VK_OEM_PLUS =+
	Key_Minus        KeyCode = 0xBD       // VK_OEM_MINUS -_
	Key_LeftBracket  KeyCode = 0xDB       // VK_OEM_4 [{
	Key_RightBracket KeyCode = 0xDD       // VK_OEM_6 ]}
	Key_Quote        KeyCode = 0xDE       // VK_OEM_7 '"
	Key_Semicolon    KeyCode = 0xBA       // VK_OEM_1 ;:
	Key_Comma        KeyCode = 0xBC       // VK_OEM_COMMA ,<
	Key_Period       KeyCode = 0xBE       // VK_OEM_PERIOD .>
	Key_Slash        KeyCode = 0xBF       // VK_OEM_2 /?
	Key_Backslash    KeyCode = 0xDC       // VK_OEM_5 \|
	Key_Grave        KeyCode = 0xC0       // VK_OEM_3 `~
	Key_Menu         KeyCode = 0x5D       // VK_APPS
	Key_Function     KeyCode = 0x3F       // kVK_Function FN TODO::
	Key_Left         KeyCode = 0x25       // VK_LEFT
	Key_Right        KeyCode = 0x27       // VK_RIGHT
	Key_Down         KeyCode = 0x28       // VK_DOWN
	Key_Up           KeyCode = 0x26       // VK_UP
	Key_Mute         KeyCode = 0xAD       // VK_VOLUME_MUTE
	Key_VolumeUp     KeyCode = 0xAF       // VK_VOLUME_UP
	Key_VolumeDown   KeyCode = 0xAE       // VK_VOLUME_DOWN
	Key_Home         KeyCode = 0x24       // VK_HOME
	Key_End          KeyCode = 0x23       // VK_END
	Key_PageUp       KeyCode = 0x21       // VK_PRIOR
	Key_PageDown     KeyCode = 0x22       // VK_NEXT
	Key_Delete       KeyCode = 0x2E       // VK_DELETE
	Key_Insert       KeyCode = 0x2D       // VK_INSERT
	Key_Help         KeyCode = 0x2F       // VK_HELP
	Key_Pad0         KeyCode = 0x60       // VK_NUMPAD0
	Key_Pad1         KeyCode = 0x61       // VK_NUMPAD1
	Key_Pad2         KeyCode = 0x62       // VK_NUMPAD2
	Key_Pad3         KeyCode = 0x63       // VK_NUMPAD3
	Key_Pad4         KeyCode = 0x64       // VK_NUMPAD4
	Key_Pad5         KeyCode = 0x65       // VK_NUMPAD5
	Key_Pad6         KeyCode = 0x66       // VK_NUMPAD6
	Key_Pad7         KeyCode = 0x67       // VK_NUMPAD7
	Key_Pad8         KeyCode = 0x68       // VK_NUMPAD8
	Key_Pad9         KeyCode = 0x69       // VK_NUMPAD9
	Key_PadDecimal   KeyCode = 0x6E       // VK_DECIMAL .
	Key_PadPlus      KeyCode = 0x6B       // VK_ADD +
	Key_PadMinus     KeyCode = 0x6D       // VK_SUBTRACT -
	Key_PadMultiply  KeyCode = 0x6A       // VK_MULTIPLY *
	Key_PadDivide    KeyCode = 0x6F       // VK_DIVIDE /
	Key_PadEnter     KeyCode = 0x0D       // VK_RETURN
	Key_PadNumLock   KeyCode = 0x90       // VK_NUMLOCK
	Key_PadBegin     KeyCode = Key_Unknow // Unknow
	Key_Clear        KeyCode = 0x0C       // VK_CLEAR
	Key_ScrollLock   KeyCode = 0x91       // VK_SCROLL
	Key_Pause        KeyCode = 0x13       // VK_PAUSE
	Key_Snapshot     KeyCode = 0x2C       // VK_SNAPSHOT
	Key_PadEquals    KeyCode = Key_Unknow // Unknow
	Key_Back         KeyCode = Key_Unknow // AKEYCODE_BACK
)
