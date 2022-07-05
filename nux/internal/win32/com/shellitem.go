// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package com

import (
	"syscall"
	"unsafe"
)

var (
	IID_IShellItem = GUID{0x43826D1E, 0xE718, 0x42EE, [8]byte{0xBC, 0x55, 0xA1, 0xE2, 0x61, 0xC3, 0x7B, 0xFE}}
)

// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitem
type IShellItem struct {
	IUnknown
}

type IShellItemVtbl struct {
	IUnknownVtbl
	BindToHandler  uintptr
	GetParent      uintptr
	GetDisplayName uintptr
	GetAttributes  uintptr
	Compare        uintptr
}

func (me *IShellItem) GetDisplayName(sigdnName SIGDN) (string, error) {
	var szName *uint16
	vtbl := (*IShellItemVtbl)(unsafe.Pointer(me.Vtbl))
	_, _, err := syscall.Syscall(vtbl.GetDisplayName, 3,
		uintptr(unsafe.Pointer(me)),
		uintptr(sigdnName),
		uintptr(unsafe.Pointer(&szName)))
	if err == 0 {
		defer CoTaskMemFree(unsafe.Pointer(szName))
		return uint16PtrToString(szName), nil
	}
	return "", err
}

func uint16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}

	var s []uint16
	end := unsafe.Pointer(p)
	for n := 0; *(*uint16)(end) != 0; n++ {
		s = append(s, *(*uint16)(end))
		end = unsafe.Pointer(uintptr(end) + 2)
	}

	return syscall.UTF16ToString(s)
}

// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitemarray
type IShellItemArray struct {
	IUnknown
}

type IShellItemArrayVtbl struct {
	IUnknownVtbl
	BindToHandler              uintptr
	GetPropertyStore           uintptr
	GetPropertyDescriptionList uintptr
	GetAttributes              uintptr
	GetCount                   uintptr
	GetItemAt                  uintptr
	EnumItems                  uintptr
}

func (me *IShellItemArray) GetCount() (uint32, error) {
	var count uint32
	vtbl := (*IShellItemArrayVtbl)(unsafe.Pointer(me.Vtbl))
	_, _, err := syscall.Syscall(vtbl.GetCount, 2,
		uintptr(unsafe.Pointer(me)),
		uintptr(unsafe.Pointer(&count)),
		0)
	if err == 0 {
		return count, nil
	}
	return 0, err
}

func (me *IShellItemArray) GetItemAt(index uint32) (*IShellItem, error) {
	var item *IShellItem
	vtbl := (*IShellItemArrayVtbl)(unsafe.Pointer(me.Vtbl))
	_, _, err := syscall.Syscall(vtbl.GetItemAt, 3,
		uintptr(unsafe.Pointer(me)),
		uintptr(index),
		uintptr(unsafe.Pointer(&item)))
	if err == 0 {
		return item, nil
	}
	return nil, err
}
