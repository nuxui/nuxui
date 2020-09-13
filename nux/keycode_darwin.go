// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package nux

// https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.6.sdk/System/Library/Frameworks/Carbon.framework/Versions/A/Frameworks/HIToolbox.framework/Versions/A/Headers/Events.h

const (
	KeyCode_A              KeyCode = 0x00 //  kVK_ANSI_A                    = 0x00,
	KeyCode_S              KeyCode = 0x01 //  kVK_ANSI_S                    = 0x01,
	KeyCode_D              KeyCode = 0x02 //  kVK_ANSI_D                    = 0x02,
	KeyCode_F              KeyCode = 0x03 //  kVK_ANSI_F                    = 0x03,
	KeyCode_H              KeyCode = 0x04 //  kVK_ANSI_H                    = 0x04,
	KeyCode_G              KeyCode = 0x05 //  kVK_ANSI_G                    = 0x05,
	KeyCode_Z              KeyCode = 0x06 //  kVK_ANSI_Z                    = 0x06,
	KeyCode_X              KeyCode = 0x07 //  kVK_ANSI_X                    = 0x07,
	KeyCode_C              KeyCode = 0x08 //  kVK_ANSI_C                    = 0x08,
	KeyCode_V              KeyCode = 0x09 //  kVK_ANSI_V                    = 0x09,
	KeyCode_B              KeyCode = 0x0B //  kVK_ANSI_B                    = 0x0B,
	KeyCode_Q              KeyCode = 0x0C //  kVK_ANSI_Q                    = 0x0C,
	KeyCode_W              KeyCode = 0x0D //  kVK_ANSI_W                    = 0x0D,
	KeyCode_E              KeyCode = 0x0E //  kVK_ANSI_E                    = 0x0E,
	KeyCode_R              KeyCode = 0x0F //  kVK_ANSI_R                    = 0x0F,
	KeyCode_Y              KeyCode = 0x10 //  kVK_ANSI_Y                    = 0x10,
	KeyCode_T              KeyCode = 0x11 //  kVK_ANSI_T                    = 0x11,
	KeyCode_1              KeyCode = 0x12 //  kVK_ANSI_1                    = 0x12,
	KeyCode_2              KeyCode = 0x13 //  kVK_ANSI_2                    = 0x13,
	KeyCode_3              KeyCode = 0x14 //  kVK_ANSI_3                    = 0x14,
	KeyCode_4              KeyCode = 0x15 //  kVK_ANSI_4                    = 0x15,
	KeyCode_6              KeyCode = 0x16 //  kVK_ANSI_6                    = 0x16,
	KeyCode_5              KeyCode = 0x17 //  kVK_ANSI_5                    = 0x17,
	KeyCode_Equal          KeyCode = 0x18 //  kVK_ANSI_Equal                = 0x18,
	KeyCode_9              KeyCode = 0x19 //  kVK_ANSI_9                    = 0x19,
	KeyCode_7              KeyCode = 0x1A //  kVK_ANSI_7                    = 0x1A,
	KeyCode_Minus          KeyCode = 0x1B //  kVK_ANSI_Minus                = 0x1B,
	KeyCode_8              KeyCode = 0x1C //  kVK_ANSI_8                    = 0x1C,
	KeyCode_0              KeyCode = 0x1D //  kVK_ANSI_0                    = 0x1D,
	KeyCode_RightBracket   KeyCode = 0x1E //  kVK_ANSI_RightBracket         = 0x1E,
	KeyCode_O              KeyCode = 0x1F //  kVK_ANSI_O                    = 0x1F,
	KeyCode_U              KeyCode = 0x20 //  kVK_ANSI_U                    = 0x20,
	KeyCode_LeftBracket    KeyCode = 0x21 //  kVK_ANSI_LeftBracket          = 0x21,
	KeyCode_I              KeyCode = 0x22 //  kVK_ANSI_I                    = 0x22,
	KeyCode_P              KeyCode = 0x23 //  kVK_ANSI_P                    = 0x23,
	KeyCode_L              KeyCode = 0x25 //  kVK_ANSI_L                    = 0x25,
	KeyCode_J              KeyCode = 0x26 //  kVK_ANSI_J                    = 0x26,
	KeyCode_Quote          KeyCode = 0x27 //  kVK_ANSI_Quote                = 0x27,
	KeyCode_K              KeyCode = 0x28 //  kVK_ANSI_K                    = 0x28,
	KeyCode_Semicolon      KeyCode = 0x29 //  kVK_ANSI_Semicolon            = 0x29,
	KeyCode_Backslash      KeyCode = 0x2A //  kVK_ANSI_Backslash            = 0x2A,
	KeyCode_Comma          KeyCode = 0x2B //  kVK_ANSI_Comma                = 0x2B,
	KeyCode_Slash          KeyCode = 0x2C //  kVK_ANSI_Slash                = 0x2C,
	KeyCode_N              KeyCode = 0x2D //  kVK_ANSI_N                    = 0x2D,
	KeyCode_M              KeyCode = 0x2E //  kVK_ANSI_M                    = 0x2E,
	KeyCode_Period         KeyCode = 0x2F //  kVK_ANSI_Period               = 0x2F,
	KeyCode_Grave          KeyCode = 0x32 //  kVK_ANSI_Grave                = 0x32,
	KeyCode_NumpadDecimal  KeyCode = 0x41 //  kVK_ANSI_KeypadDecimal        = 0x41,
	KeyCode_NumpadMultiply KeyCode = 0x43 //  kVK_ANSI_KeypadMultiply       = 0x43,
	KeyCode_NumpadAdd      KeyCode = 0x45 //  kVK_ANSI_KeypadPlus           = 0x45,
	KeyCode_NumpadClear    KeyCode = 0x47 //  kVK_ANSI_KeypadClear          = 0x47,
	KeyCode_NumpadDivide   KeyCode = 0x4B //  kVK_ANSI_KeypadDivide         = 0x4B,
	KeyCode_NumpadEnter    KeyCode = 0x4C //  kVK_ANSI_KeypadEnter          = 0x4C,
	KeyCode_NumpadSubtract KeyCode = 0x4E //  kVK_ANSI_KeypadMinus          = 0x4E,
	KeyCode_NumpadEquals   KeyCode = 0x51 //  kVK_ANSI_KeypadEquals         = 0x51,
	KeyCode_Numpad0        KeyCode = 0x52 //  kVK_ANSI_Keypad0              = 0x52,
	KeyCode_Numpad1        KeyCode = 0x53 //  kVK_ANSI_Keypad1              = 0x53,
	KeyCode_Numpad2        KeyCode = 0x54 //  kVK_ANSI_Keypad2              = 0x54,
	KeyCode_Numpad3        KeyCode = 0x55 //  kVK_ANSI_Keypad3              = 0x55,
	KeyCode_Numpad4        KeyCode = 0x56 //  kVK_ANSI_Keypad4              = 0x56,
	KeyCode_Numpad5        KeyCode = 0x57 //  kVK_ANSI_Keypad5              = 0x57,
	KeyCode_Numpad6        KeyCode = 0x58 //  kVK_ANSI_Keypad6              = 0x58,
	KeyCode_Numpad7        KeyCode = 0x59 //  kVK_ANSI_Keypad7              = 0x59,
	KeyCode_Numpad8        KeyCode = 0x5B //  kVK_ANSI_Keypad8              = 0x5B,
	KeyCode_Numpad9        KeyCode = 0x5C //  kVK_ANSI_Keypad9              = 0x5C
)

/* keycodes for keys that are independent of keyboard layout*/
const (
	KeyCode_Return        KeyCode = 0x24
	KeyCode_Tab           KeyCode = 0x30
	KeyCode_Space         KeyCode = 0x31
	KeyCode_Delete        KeyCode = 0x33
	KeyCode_Escape        KeyCode = 0x35
	KeyCode_Command       KeyCode = 0x37
	KeyCode_Shift         KeyCode = 0x38
	KeyCode_CapsLock      KeyCode = 0x39
	KeyCode_Option        KeyCode = 0x3A
	KeyCode_Control       KeyCode = 0x3B
	KeyCode_RightShift    KeyCode = 0x3C
	KeyCode_RightOption   KeyCode = 0x3D
	KeyCode_RightControl  KeyCode = 0x3E
	KeyCode_Function      KeyCode = 0x3F
	KeyCode_F17           KeyCode = 0x40
	KeyCode_VolumeUp      KeyCode = 0x48
	KeyCode_VolumeDown    KeyCode = 0x49
	KeyCode_Mute          KeyCode = 0x4A
	KeyCode_F18           KeyCode = 0x4F
	KeyCode_F19           KeyCode = 0x50
	KeyCode_F20           KeyCode = 0x5A
	KeyCode_F5            KeyCode = 0x60
	KeyCode_F6            KeyCode = 0x61
	KeyCode_F7            KeyCode = 0x62
	KeyCode_F3            KeyCode = 0x63
	KeyCode_F8            KeyCode = 0x64
	KeyCode_F9            KeyCode = 0x65
	KeyCode_F11           KeyCode = 0x67
	KeyCode_F13           KeyCode = 0x69
	KeyCode_F16           KeyCode = 0x6A
	KeyCode_F14           KeyCode = 0x6B
	KeyCode_F10           KeyCode = 0x6D
	KeyCode_F12           KeyCode = 0x6F
	KeyCode_F15           KeyCode = 0x71
	KeyCode_Help          KeyCode = 0x72
	KeyCode_Home          KeyCode = 0x73
	KeyCode_PageUp        KeyCode = 0x74
	KeyCode_ForwardDelete KeyCode = 0x75
	KeyCode_F4            KeyCode = 0x76
	KeyCode_End           KeyCode = 0x77
	KeyCode_F2            KeyCode = 0x78
	KeyCode_PageDown      KeyCode = 0x79
	KeyCode_F1            KeyCode = 0x7A
	KeyCode_LeftArrow     KeyCode = 0x7B
	KeyCode_RightArrow    KeyCode = 0x7C
	KeyCode_DownArrow     KeyCode = 0x7D
	KeyCode_UpArrow       KeyCode = 0x7E
)
