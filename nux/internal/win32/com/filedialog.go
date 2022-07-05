// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package com

import (
	"syscall"
	"unsafe"
)

// ShObjIdl_core.h

var (
	CLSID_FileOpenDialog = GUID{0xDC1C5A9C, 0xE88A, 0x4DDE, [8]byte{0xA5, 0xA1, 0x60, 0xF8, 0x2A, 0x20, 0xAE, 0xF7}}
	CLSID_FileSaveDialog = GUID{0xC0B4E2F3, 0xBA21, 0x4773, [8]byte{0x8D, 0xBA, 0x33, 0x5E, 0xC9, 0x46, 0xEB, 0x8B}}
	IID_IFileDialog      = GUID{0x42F85136, 0xDB7E, 0x439C, [8]byte{0x85, 0xF1, 0xE4, 0x07, 0x5D, 0x13, 0x5F, 0xC8}}
	IID_IFileSaveDialog  = GUID{0x84BCCD23, 0x5FDE, 0x4CDB, [8]byte{0xAE, 0xA4, 0xAF, 0x64, 0xB8, 0x3D, 0x78, 0xAB}}
	IID_IFileOpenDialog  = GUID{0xD57C7288, 0xD4AD, 0x4768, [8]byte{0xBE, 0x02, 0x9D, 0x96, 0x95, 0x32, 0xD9, 0x60}}
)

type IModalWindow struct {
	IUnknown
}

type IModalWindowVtbl struct {
	IUnknownVtbl
	Show uintptr
}

func (me *IModalWindow) Show(hwnd uintptr) (HRESULT, error) {
	vtbl := (*IModalWindowVtbl)(unsafe.Pointer(me.Vtbl))
	ret, _, err := syscall.Syscall(vtbl.Show, 2,
		uintptr(unsafe.Pointer(me)),
		hwnd,
		0)
	if err == 0 || ret == ERROR_CANCELLED {
		return HRESULT(ret), nil
	}
	return HRESULT(ret), err
}

type IFileDialog struct {
	IModalWindow
}

type IFileDialogVtbl struct {
	IModalWindowVtbl
	SetFileTypes        uintptr
	SetFileTypeIndex    uintptr
	GetFileTypeIndex    uintptr
	Advise              uintptr
	Unadvise            uintptr
	SetOptions          uintptr
	GetOptions          uintptr
	SetDefaultFolder    uintptr
	SetFolder           uintptr
	GetFolder           uintptr
	GetCurrentSelection uintptr
	SetFileName         uintptr
	GetFileName         uintptr
	SetTitle            uintptr
	SetOkButtonLabel    uintptr
	SetFileNameLabel    uintptr
	GetResult           uintptr
	AddPlace            uintptr
	SetDefaultExtension uintptr
	Close               uintptr
	SetClientGuid       uintptr
	ClearClientData     uintptr
	SetFilter           uintptr
}

func (me *IFileDialog) SetOptions(fos FILEOPENDIALOGOPTIONS) error {
	vtbl := (*IFileDialogVtbl)(unsafe.Pointer(me.Vtbl))
	_, _, err := syscall.Syscall(vtbl.SetOptions, 2,
		uintptr(unsafe.Pointer(me)),
		uintptr(fos),
		0)
	if err == 0 {
		return nil
	}
	return err
}

func (me *IFileDialog) GetOptions() (FILEOPENDIALOGOPTIONS, error) {
	var fos FILEOPENDIALOGOPTIONS
	vtbl := (*IFileDialogVtbl)(unsafe.Pointer(me.Vtbl))
	_, _, err := syscall.Syscall(vtbl.GetOptions, 2,
		uintptr(unsafe.Pointer(me)),
		uintptr(unsafe.Pointer(&fos)),
		0)
	if err == 0 {
		return fos, nil
	}
	return 0, err
}

func (me *IFileDialog) GetResult() (*IShellItem, error) {
	var item *IShellItem
	vtbl := (*IFileDialogVtbl)(unsafe.Pointer(me.Vtbl))
	_, _, err := syscall.Syscall(vtbl.GetResult, 2,
		uintptr(unsafe.Pointer(me)),
		uintptr(unsafe.Pointer(&item)),
		0)
	if err == 0 {
		return item, nil
	}
	return nil, err
}

type IFileOpenDialog struct {
	IFileDialog
}

type IFileOpenDialogVtbl struct {
	IFileDialogVtbl
	GetResults       uintptr
	GetSelectedItems uintptr
}

func (me *IFileOpenDialog) GetResults() (*IShellItemArray, error) {
	var arr *IShellItemArray
	vtbl := (*IFileOpenDialogVtbl)(unsafe.Pointer(me.Vtbl))
	_, _, err := syscall.Syscall(vtbl.GetResults, 2,
		uintptr(unsafe.Pointer(me)),
		uintptr(unsafe.Pointer(&arr)),
		0)
	if err == 0 {
		return arr, nil
	}
	return nil, err
}
