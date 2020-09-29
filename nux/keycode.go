// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

//go:generate stringer -type=KeyCode
type KeyCode int

const (
	Mod_CapsLock uint32 = 0x10000 << iota
	Mod_Shift    uint32 = 0x10000 << iota
	Mod_Control  uint32 = 0x10000 << iota
	Mod_Alt      uint32 = 0x10000 << iota
	Mod_Super    uint32 = 0x10000 << iota
	Mod_Mask     uint32 = 0xFFFF0000
)

//TODO::equalasciicode
//TODO::codenumber
const (
	Key_Unknown KeyCode = iota
	Key_A
	Key_B
	Key_C
	Key_D
	Key_E
	Key_F
	Key_G
	Key_H
	Key_I
	Key_J
	Key_K
	Key_L
	Key_M
	Key_N
	Key_O
	Key_P
	Key_Q
	Key_R
	Key_S
	Key_T
	Key_U
	Key_V
	Key_W
	Key_X
	Key_Y
	Key_Z
	Key_0
	Key_1
	Key_2
	Key_3
	Key_4
	Key_5
	Key_6
	Key_7
	Key_8
	Key_9
	Key_F1
	Key_F2
	Key_F3
	Key_F4
	Key_F5
	Key_F6
	Key_F7
	Key_F8
	Key_F9
	Key_F10
	Key_F11
	Key_F12
	Key_F13
	Key_F14
	Key_F15
	Key_F16
	Key_F17
	Key_F18
	Key_F19
	Key_F20

	Key_Return
	Key_Tab
	Key_Space
	Key_Delete
	Key_Escape
	Key_Super
	Key_CapsLock
	Key_AltLeft
	Key_AltRight
	Key_ShiftLeft
	Key_ShiftRight
	Key_ControlLeft
	Key_ControlRight
	Key_Equal        //=+
	Key_Minus        //-_ TODO:: sub
	Key_BracketLeft  //[{
	Key_BracketRight //]}
	Key_Quote        //'"
	Key_Semicolon    //;:
	Key_Comma        //,<
	Key_Slash        // / ?
	Key_Backslash    // | \
	Key_Period       //.>
	Key_Grave        //`~

	Key_Pad0
	Key_Pad1
	Key_Pad2
	Key_Pad3
	Key_Pad4
	Key_Pad5
	Key_Pad6
	Key_Pad7
	Key_Pad8
	Key_Pad9
	Key_PadAdd      //+
	Key_PadSubtract //-
	Key_PadMultiply //*
	Key_PadDivide   ///
	Key_PadEquals   //=
	Key_PadClear    //num_lockTODO::
	Key_PadDecimal  //.
	Key_PadEnter
	Key_ArrowUp
	Key_ArrowDown
	Key_ArrowLeft
	Key_ArrowRight
	Key_PageUp
	Key_PageDown
	Key_Home
	Key_End

	Key_Help          //TODO::
	Key_ForwardDelete //TODO::

	Key_VolumeUp
	Key_VolumeDown
	Key_Mute
)

func (me KeyCode) String() string {
	switch me {
	case Key_Unknown:
		return "Key_Unknown "
	case Key_A:
		return "Key_A"
	case Key_B:
		return "Key_B"
	case Key_C:
		return "Key_C"
	case Key_D:
		return "Key_D"
	case Key_E:
		return "Key_E"
	case Key_F:
		return "Key_F"
	case Key_G:
		return "Key_G"
	case Key_H:
		return "Key_H"
	case Key_I:
		return "Key_I"
	case Key_J:
		return "Key_J"
	case Key_K:
		return "Key_K"
	case Key_L:
		return "Key_L"
	case Key_M:
		return "Key_M"
	case Key_N:
		return "Key_N"
	case Key_O:
		return "Key_O"
	case Key_P:
		return "Key_P"
	case Key_Q:
		return "Key_Q"
	case Key_R:
		return "Key_R"
	case Key_S:
		return "Key_S"
	case Key_T:
		return "Key_T"
	case Key_U:
		return "Key_U"
	case Key_V:
		return "Key_V"
	case Key_W:
		return "Key_W"
	case Key_X:
		return "Key_X"
	case Key_Y:
		return "Key_Y"
	case Key_Z:
		return "Key_Z"
	case Key_0:
		return "Key_0"
	case Key_1:
		return "Key_1"
	case Key_2:
		return "Key_2"
	case Key_3:
		return "Key_3"
	case Key_4:
		return "Key_4"
	case Key_5:
		return "Key_5"
	case Key_6:
		return "Key_6"
	case Key_7:
		return "Key_7"
	case Key_8:
		return "Key_8"
	case Key_9:
		return "Key_9"
	case Key_F1:
		return "Key_F1"
	case Key_F2:
		return "Key_F2"
	case Key_F3:
		return "Key_F3"
	case Key_F4:
		return "Key_F4"
	case Key_F5:
		return "Key_F5"
	case Key_F6:
		return "Key_F6"
	case Key_F7:
		return "Key_F7"
	case Key_F8:
		return "Key_F8"
	case Key_F9:
		return "Key_F9"
	case Key_F10:
		return "Key_F10"
	case Key_F11:
		return "Key_F11"
	case Key_F12:
		return "Key_F12"
	case Key_F13:
		return "Key_F13"
	case Key_F14:
		return "Key_F14"
	case Key_F15:
		return "Key_F15"
	case Key_F16:
		return "Key_F16"
	case Key_F17:
		return "Key_F17"
	case Key_F18:
		return "Key_F18"
	case Key_F19:
		return "Key_F19"
	case Key_F20:
		return "Key_F20"

	case Key_Return:
		return "Key_Return"
	case Key_Tab:
		return "Key_Tab"
	case Key_Space:
		return "Key_Space"
	case Key_Delete:
		return "Key_Delete"
	case Key_Escape:
		return "Key_Escape"
	case Key_Super:
		return "Key_Super"
	case Key_CapsLock:
		return "Key_CapsLock"
	case Key_AltLeft:
		return "Key_AltLeft"
	case Key_AltRight:
		return "Key_AltRight"
	case Key_ShiftLeft:
		return "Key_ShiftLeft"
	case Key_ShiftRight:
		return "Key_ShiftRight"
	case Key_ControlLeft:
		return "Key_ControlLeft"
	case Key_ControlRight:
		return "Key_ControlRight"
	case Key_Equal:
		return "Key_Equal"
	case Key_Minus:
		return "Key_Minus"
	case Key_BracketLeft:
		return "Key_BracketLeft"
	case Key_BracketRight:
		return "Key_BracketRight"
	case Key_Quote:
		return "Key_Quote"
	case Key_Semicolon:
		return "Key_Semicolon"
	case Key_Comma:
		return "Key_Comma"
	case Key_Slash:
		return "Key_Slash"
	case Key_Backslash:
		return "Key_Backslash"
	case Key_Period:
		return "Key_Period"
	case Key_Grave:
		return "Key_Grave"

	case Key_Pad0:
		return "Key_Pad0"
	case Key_Pad1:
		return "Key_Pad1"
	case Key_Pad2:
		return "Key_Pad2"
	case Key_Pad3:
		return "Key_Pad3"
	case Key_Pad4:
		return "Key_Pad4"
	case Key_Pad5:
		return "Key_Pad5"
	case Key_Pad6:
		return "Key_Pad6"
	case Key_Pad7:
		return "Key_Pad7"
	case Key_Pad8:
		return "Key_Pad8"
	case Key_Pad9:
		return "Key_Pad9"
	case Key_PadAdd:
		return "Key_PadAdd"
	case Key_PadSubtract:
		return "Key_PadSubtract"
	case Key_PadMultiply:
		return "Key_PadMultiply"
	case Key_PadDivide:
		return "Key_PadDivide"
	case Key_PadEquals:
		return "Key_PadEquals"
	case Key_PadClear:
		return "Key_PadClear"
	case Key_PadDecimal:
		return "Key_PadDecimal"
	case Key_PadEnter:
		return "Key_PadEnter"
	case Key_ArrowUp:
		return "Key_ArrowUp"
	case Key_ArrowDown:
		return "Key_ArrowDown"
	case Key_ArrowLeft:
		return "Key_ArrowLeft"
	case Key_ArrowRight:
		return "Key_ArrowRight"
	case Key_PageUp:
		return "Key_PageUp"
	case Key_PageDown:
		return "Key_PageDown"
	case Key_Home:
		return "Key_Home"
	case Key_End:
		return "Key_End"

	case Key_Help:
		return "Key_Help"
	case Key_ForwardDelete:
		return "Key_ForwardDelete"

	case Key_VolumeUp:
		return "Key_VolumeUp"
	case Key_VolumeDown:
		return "Key_VolumeDown"
	case Key_Mute:
		return "Key_Mute"
	}
	log.Fatal("nuxui", "KeyCode %d not handled in switch", me)
	return ""
}
