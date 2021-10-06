// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type KeyCode int

const (
	Mod_CapsLock uint32 = 0x10000 << iota
	Mod_Shift    uint32 = 0x10000 << iota
	Mod_Control  uint32 = 0x10000 << iota
	Mod_Alt      uint32 = 0x10000 << iota
	Mod_Super    uint32 = 0x10000 << iota
	Mod_Mask     uint32 = 0xFFFF0000
)
