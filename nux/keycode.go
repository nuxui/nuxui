// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type KeyCode uint32

const (
	Key_Unknow KeyCode = 0xffffffff
)

const (
	Mod_CapsLock uint32 = 0x10000 << iota
	Mod_Shift    uint32 = 0x10000 << iota
	Mod_Control  uint32 = 0x10000 << iota
	Mod_Alt      uint32 = 0x10000 << iota
	Mod_Super    uint32 = 0x10000 << iota
	Mod_Mask     uint32 = 0xFFFF0000
)

func (keyCode KeyCode) String() string {
	if keyCode == Key_A {
		return "Key_A"
	} else if keyCode == Key_B {
		return "Key_B"
	} else if keyCode == Key_C {
		return "Key_C"
	} else if keyCode == Key_D {
		return "Key_D"
	} else if keyCode == Key_E {
		return "Key_E"
	} else if keyCode == Key_F {
		return "Key_F"
	} else if keyCode == Key_G {
		return "Key_G"
	} else if keyCode == Key_H {
		return "Key_H"
	} else if keyCode == Key_I {
		return "Key_I"
	} else if keyCode == Key_J {
		return "Key_J"
	} else if keyCode == Key_K {
		return "Key_K"
	} else if keyCode == Key_L {
		return "Key_L"
	} else if keyCode == Key_M {
		return "Key_M"
	} else if keyCode == Key_N {
		return "Key_N"
	} else if keyCode == Key_O {
		return "Key_O"
	} else if keyCode == Key_P {
		return "Key_P"
	} else if keyCode == Key_Q {
		return "Key_Q"
	} else if keyCode == Key_R {
		return "Key_R"
	} else if keyCode == Key_S {
		return "Key_S"
	} else if keyCode == Key_T {
		return "Key_T"
	} else if keyCode == Key_U {
		return "Key_U"
	} else if keyCode == Key_V {
		return "Key_V"
	} else if keyCode == Key_W {
		return "Key_W"
	} else if keyCode == Key_X {
		return "Key_X"
	} else if keyCode == Key_Y {
		return "Key_Y"
	} else if keyCode == Key_Z {
		return "Key_Z"
	} else if keyCode == Key_0 {
		return "Key_0"
	} else if keyCode == Key_1 {
		return "Key_1"
	} else if keyCode == Key_2 {
		return "Key_2"
	} else if keyCode == Key_3 {
		return "Key_3"
	} else if keyCode == Key_4 {
		return "Key_4"
	} else if keyCode == Key_5 {
		return "Key_5"
	} else if keyCode == Key_6 {
		return "Key_6"
	} else if keyCode == Key_7 {
		return "Key_7"
	} else if keyCode == Key_8 {
		return "Key_8"
	} else if keyCode == Key_9 {
		return "Key_9"
	} else if keyCode == Key_F1 {
		return "Key_F1"
	} else if keyCode == Key_F2 {
		return "Key_F2"
	} else if keyCode == Key_F3 {
		return "Key_F3"
	} else if keyCode == Key_F4 {
		return "Key_F4"
	} else if keyCode == Key_F5 {
		return "Key_F5"
	} else if keyCode == Key_F6 {
		return "Key_F6"
	} else if keyCode == Key_F7 {
		return "Key_F7"
	} else if keyCode == Key_F8 {
		return "Key_F8"
	} else if keyCode == Key_F9 {
		return "Key_F9"
	} else if keyCode == Key_F10 {
		return "Key_F10"
	} else if keyCode == Key_F11 {
		return "Key_F11"
	} else if keyCode == Key_F12 {
		return "Key_F12"
	} else if keyCode == Key_F13 {
		return "Key_F13"
	} else if keyCode == Key_F14 {
		return "Key_F14"
	} else if keyCode == Key_F15 {
		return "Key_F15"
	} else if keyCode == Key_F16 {
		return "Key_F16"
	} else if keyCode == Key_F17 {
		return "Key_F17"
	} else if keyCode == Key_F18 {
		return "Key_F18"
	} else if keyCode == Key_F19 {
		return "Key_F19"
	} else if keyCode == Key_F20 {
		return "Key_F20"
	} else if keyCode == Key_Return {
		return "Key_Return"
	} else if keyCode == Key_Tab {
		return "Key_Tab"
	} else if keyCode == Key_Space {
		return "Key_Space"
	} else if keyCode == Key_BackSpace {
		return "Key_BackSpace"
	} else if keyCode == Key_Escape {
		return "Key_Escape"
	} else if keyCode == Key_CapsLock {
		return "Key_CapsLock"
	} else if keyCode == Key_Alt {
		return "Key_Alt"
	} else if keyCode == Key_RightAlt {
		return "Key_RightAlt"
	} else if keyCode == Key_Shift {
		return "Key_Shift"
	} else if keyCode == Key_RightShift {
		return "Key_RightShift"
	} else if keyCode == Key_Control {
		return "Key_Control"
	} else if keyCode == Key_RightControl {
		return "Key_RightControl"
	} else if keyCode == Key_Command {
		return "Key_Command"
	} else if keyCode == Key_Equal {
		return "Key_Equal"
	} else if keyCode == Key_Minus {
		return "Key_Minus"
	} else if keyCode == Key_LeftBracket {
		return "Key_LeftBracket"
	} else if keyCode == Key_RightBracket {
		return "Key_RightBracket"
	} else if keyCode == Key_Quote {
		return "Key_Quote"
	} else if keyCode == Key_Semicolon {
		return "Key_Semicolon"
	} else if keyCode == Key_Comma {
		return "Key_Comma"
	} else if keyCode == Key_Period {
		return "Key_Period"
	} else if keyCode == Key_Slash {
		return "Key_Slash"
	} else if keyCode == Key_Backslash {
		return "Key_Backslash"
	} else if keyCode == Key_Grave {
		return "Key_Grave"
	} else if keyCode == Key_Menu {
		return "Key_Menu"
	} else if keyCode == Key_Function {
		return "Key_Function"
	} else if keyCode == Key_Left {
		return "Key_Left"
	} else if keyCode == Key_Right {
		return "Key_Right"
	} else if keyCode == Key_Down {
		return "Key_Down"
	} else if keyCode == Key_Up {
		return "Key_Up"
	} else if keyCode == Key_Mute {
		return "Key_Mute"
	} else if keyCode == Key_VolumeUp {
		return "Key_VolumeUp"
	} else if keyCode == Key_VolumeDown {
		return "Key_VolumeDown"
	} else if keyCode == Key_Home {
		return "Key_Home"
	} else if keyCode == Key_End {
		return "Key_End"
	} else if keyCode == Key_PageUp {
		return "Key_PageUp"
	} else if keyCode == Key_PageDown {
		return "Key_PageDown"
	} else if keyCode == Key_Delete {
		return "Key_Delete"
	} else if keyCode == Key_Insert {
		return "Key_Insert"
	} else if keyCode == Key_Help {
		return "Key_Help"
	} else if keyCode == Key_Pad0 {
		return "Key_Pad0"
	} else if keyCode == Key_Pad1 {
		return "Key_Pad1"
	} else if keyCode == Key_Pad2 {
		return "Key_Pad2"
	} else if keyCode == Key_Pad3 {
		return "Key_Pad3"
	} else if keyCode == Key_Pad4 {
		return "Key_Pad4"
	} else if keyCode == Key_Pad5 {
		return "Key_Pad5"
	} else if keyCode == Key_Pad6 {
		return "Key_Pad6"
	} else if keyCode == Key_Pad7 {
		return "Key_Pad7"
	} else if keyCode == Key_Pad8 {
		return "Key_Pad8"
	} else if keyCode == Key_Pad9 {
		return "Key_Pad9"
	} else if keyCode == Key_PadDecimal {
		return "Key_PadDecimal"
	} else if keyCode == Key_PadPlus {
		return "Key_PadPlus"
	} else if keyCode == Key_PadMinus {
		return "Key_PadMinus"
	} else if keyCode == Key_PadMultiply {
		return "Key_PadMultiply"
	} else if keyCode == Key_PadDivide {
		return "Key_PadDivide"
	} else if keyCode == Key_PadEnter {
		return "Key_PadEnter"
	} else if keyCode == Key_PadNumLock {
		return "Key_PadNumLock"
	} else if keyCode == Key_PadBegin {
		return "Key_PadBegin"
	} else if keyCode == Key_Clear {
		return "Key_Clear"
	} else if keyCode == Key_ScrollLock {
		return "Key_ScrollLock"
	} else if keyCode == Key_Pause {
		return "Key_Pause"
	} else if keyCode == Key_Snapshot {
		return "Key_Snapshot"
	} else if keyCode == Key_PadEquals {
		return "Key_PadEquals"
	}

	return "! Unknow KeyCode"
}
