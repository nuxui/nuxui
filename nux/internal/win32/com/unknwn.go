// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package com

import (
	"syscall"
	"unsafe"
)

var (
	IID_IUnknown = GUID{0x00000000, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
)

type IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

type IUnknown struct {
	Vtbl *IUnknownVtbl
}

func (self *IUnknown) QueryInterface(iid *GUID, object *uintptr) uint32 {
	r1, _, _ := syscall.Syscall(
		self.Vtbl.QueryInterface,
		3,
		uintptr(unsafe.Pointer(self)),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(object)))
	return uint32(r1)
}

func (self *IUnknown) AddRef() uint32 {
	r1, _, _ := syscall.Syscall(self.Vtbl.AddRef, 1, uintptr(unsafe.Pointer(self)), 0, 0)
	return uint32(r1)
}

func (self *IUnknown) Release() uint32 {
	r1, _, _ := syscall.Syscall(self.Vtbl.Release, 1, uintptr(unsafe.Pointer(self)), 0, 0)
	return uint32(r1)
}
