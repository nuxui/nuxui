// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package nux

// https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.6.sdk/System/Library/Frameworks/Carbon.framework/Versions/A/Frameworks/HIToolbox.framework/Versions/A/Headers/Events.h

const (
	Key_A             KeyCode = 0x00 // kVK_ANSI_A
	Key_B             KeyCode = 0x0B // kVK_ANSI_B
	Key_C             KeyCode = 0x08 // kVK_ANSI_C
	Key_D             KeyCode = 0x02 // kVK_ANSI_D
	Key_E             KeyCode = 0x0E // kVK_ANSI_E
	Key_F             KeyCode = 0x03 // kVK_ANSI_F
	Key_G             KeyCode = 0x05 // kVK_ANSI_G
	Key_H             KeyCode = 0x04 // kVK_ANSI_H
	Key_I             KeyCode = 0x22 // kVK_ANSI_I
	Key_J             KeyCode = 0x26 // kVK_ANSI_J
	Key_K             KeyCode = 0x28 // kVK_ANSI_K
	Key_L             KeyCode = 0x25 // kVK_ANSI_L
	Key_M             KeyCode = 0x2E // kVK_ANSI_M
	Key_N             KeyCode = 0x2D // kVK_ANSI_N
	Key_O             KeyCode = 0x1F // kVK_ANSI_O
	Key_P             KeyCode = 0x23 // kVK_ANSI_P
	Key_Q             KeyCode = 0x0C // kVK_ANSI_Q
	Key_R             KeyCode = 0x0F // kVK_ANSI_R
	Key_S             KeyCode = 0x01 // kVK_ANSI_S
	Key_T             KeyCode = 0x11 // kVK_ANSI_T
	Key_U             KeyCode = 0x20 // kVK_ANSI_U
	Key_V             KeyCode = 0x09 // kVK_ANSI_V
	Key_W             KeyCode = 0x0D // kVK_ANSI_W
	Key_X             KeyCode = 0x07 // kVK_ANSI_X
	Key_Y             KeyCode = 0x10 // kVK_ANSI_Y
	Key_Z             KeyCode = 0x06 // kVK_ANSI_Z
	Key_0             KeyCode = 0x1D // kVK_ANSI_0
	Key_1             KeyCode = 0x12 // kVK_ANSI_1
	Key_2             KeyCode = 0x13 // kVK_ANSI_2
	Key_3             KeyCode = 0x14 // kVK_ANSI_3
	Key_4             KeyCode = 0x15 // kVK_ANSI_4
	Key_5             KeyCode = 0x17 // kVK_ANSI_5
	Key_6             KeyCode = 0x16 // kVK_ANSI_6
	Key_7             KeyCode = 0x1A // kVK_ANSI_7
	Key_8             KeyCode = 0x1C // kVK_ANSI_8
	Key_9             KeyCode = 0x19 // kVK_ANSI_9
	Key_F1            KeyCode = 0x7A // kVK_F1
	Key_F2            KeyCode = 0x78 // kVK_F2
	Key_F3            KeyCode = 0x63 // kVK_F3
	Key_F4            KeyCode = 0x76 // kVK_F4
	Key_F5            KeyCode = 0x60 // kVK_F5
	Key_F6            KeyCode = 0x61 // kVK_F6
	Key_F7            KeyCode = 0x62 // kVK_F7
	Key_F8            KeyCode = 0x64 // kVK_F8
	Key_F9            KeyCode = 0x65 // kVK_F9
	Key_F10           KeyCode = 0x6D // kVK_F10
	Key_F11           KeyCode = 0x67 // kVK_F11
	Key_F12           KeyCode = 0x6F // kVK_F12
	Key_F13           KeyCode = 0x69 // kVK_F13  PrintScreen
	Key_F14           KeyCode = 0x6B // kVK_F14
	Key_F15           KeyCode = 0x71 // kVK_F15
	Key_F16           KeyCode = 0x6A // kVK_F16
	Key_F17           KeyCode = 0x40 // kVK_F17
	Key_F18           KeyCode = 0x4F // kVK_F18
	Key_F19           KeyCode = 0x50 // kVK_F19
	Key_F20           KeyCode = 0x5A // kVK_F20
	Key_Enter         KeyCode = 0x24 // kVK_Return
	Key_Tab           KeyCode = 0x30 // kVK_Tab
	Key_Space         KeyCode = 0x31 // kVK_Space
	Key_Delete        KeyCode = 0x33 // kVK_Delete
	Key_Escape        KeyCode = 0x35 // kVK_Escape
	Key_CapsLock      KeyCode = 0x39 // kVK_CapsLock
	Key_Alt           KeyCode = 0x3A // kVK_Option
	Key_RightAlt      KeyCode = 0x3D // kVK_RightOption
	Key_Shift         KeyCode = 0x38 // kVK_Shift
	Key_RightShift    KeyCode = 0x3C // kVK_RightShift
	Key_Control       KeyCode = 0x3B // kVK_Control
	Key_RightControl  KeyCode = 0x3E // kVK_RightControl
	Key_Command       KeyCode = 0x37 // kVK_Command 0x36
	Key_Equal         KeyCode = 0x18 // kVK_ANSI_Equal =+
	Key_Minus         KeyCode = 0x1B // kVK_ANSI_Minus -_
	Key_LeftBracket   KeyCode = 0x21 // kVK_ANSI_LeftBracket [{
	Key_RightBracket  KeyCode = 0x1E // kVK_ANSI_RightBracket ]}
	Key_Quote         KeyCode = 0x27 // kVK_ANSI_Quote '"
	Key_Semicolon     KeyCode = 0x29 // kVK_ANSI_Semicolon ;:
	Key_Comma         KeyCode = 0x2B // kVK_ANSI_Comma ,<
	Key_Period        KeyCode = 0x2F // kVK_ANSI_Period .>
	Key_Slash         KeyCode = 0x2C // kVK_ANSI_Slash /?
	Key_Backslash     KeyCode = 0x2A // kVK_ANSI_Backslash \|
	Key_Grave         KeyCode = 0x32 // kVK_ANSI_Grave `~
	Key_Menu          KeyCode = 0x6E // Window Menu
	Key_Function      KeyCode = 0x3F // kVK_Function FN
	Key_LeftArrow     KeyCode = 0x7B // kVK_LeftArrow
	Key_RightArrow    KeyCode = 0x7C // kVK_RightArrow
	Key_DownArrow     KeyCode = 0x7D // kVK_DownArrow
	Key_UpArrow       KeyCode = 0x7E // kVK_UpArrow
	Key_Mute          KeyCode = 0x4A // kVK_Mute
	Key_VolumeUp      KeyCode = 0x48 // kVK_VolumeUp
	Key_VolumeDown    KeyCode = 0x49 // kVK_VolumeDown
	Key_Home          KeyCode = 0x73 // kVK_Home
	Key_End           KeyCode = 0x77 // kVK_End
	Key_PageUp        KeyCode = 0x74 // kVK_PageUp
	Key_PageDown      KeyCode = 0x79 // kVK_PageDown
	Key_ForwardDelete KeyCode = 0x75 // kVK_ForwardDelete
	Key_Insert        KeyCode = 0x72 // kVK_Help
	Key_Pad0          KeyCode = 0x52 // kVK_ANSI_Keypad0
	Key_Pad1          KeyCode = 0x53 // kVK_ANSI_Keypad1
	Key_Pad2          KeyCode = 0x54 // kVK_ANSI_Keypad2
	Key_Pad3          KeyCode = 0x55 // kVK_ANSI_Keypad3
	Key_Pad4          KeyCode = 0x56 // kVK_ANSI_Keypad4
	Key_Pad5          KeyCode = 0x57 // kVK_ANSI_Keypad5
	Key_Pad6          KeyCode = 0x58 // kVK_ANSI_Keypad6
	Key_Pad7          KeyCode = 0x59 // kVK_ANSI_Keypad7
	Key_Pad8          KeyCode = 0x5B // kVK_ANSI_Keypad8
	Key_Pad9          KeyCode = 0x5C // kVK_ANSI_Keypad9
	Key_PadPlus       KeyCode = 0x45 // kVK_ANSI_KeypadPlus +
	Key_PadMinus      KeyCode = 0x4E // kVK_ANSI_KeypadMinus -
	Key_PadMultiply   KeyCode = 0x43 // kVK_ANSI_KeypadMultiply *
	Key_PadDivide     KeyCode = 0x4B // kVK_ANSI_KeypadDivide /
	Key_PadDecimal    KeyCode = 0x41 // kVK_ANSI_KeypadDecimal .
	Key_PadEnter      KeyCode = 0x4C // kVK_ANSI_KeypadEnter
	Key_PadNumLock    KeyCode = 0x47 // kVK_ANSI_KeypadClear

	Key_PadEquals KeyCode = 0x51 // kVK_ANSI_KeypadEquals
)
