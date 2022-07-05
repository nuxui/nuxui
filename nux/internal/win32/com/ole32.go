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
	modole32             = syscall.NewLazyDLL("ole32.dll")
	procCoInitialize     = modole32.NewProc("CoInitialize")
	procCoInitializeEx   = modole32.NewProc("CoInitializeEx")
	procCoUninitialize   = modole32.NewProc("CoUninitialize")
	procCoCreateInstance = modole32.NewProc("CoCreateInstance")
	procCoTaskMemFree    = modole32.NewProc("CoTaskMemFree")
)

func CoInitialize(reserved *byte) (HRESULT, error) {
	ret, _, err := procCoInitialize.Call(uintptr(unsafe.Pointer(reserved)))
	if ret == S_OK {
		return HRESULT(ret), nil
	}
	return HRESULT(ret), err
}

// https://docs.microsoft.com/en-us/windows/win32/api/combaseapi/nf-combaseapi-coinitializeex
func CoInitializeEx(reserved *byte, flags COINIT) (HRESULT, error) {
	ret, _, err := procCoInitializeEx.Call(uintptr(unsafe.Pointer(reserved)), uintptr(flags))
	if ret == S_OK {
		return HRESULT(ret), nil
	}
	return HRESULT(ret), err
}

func CoUninitialize() {
	procCoUninitialize.Call()
}

// https://docs.microsoft.com/en-us/windows/win32/api/combaseapi/nf-combaseapi-cocreateinstance
// https://github.com/tpn/winsdk-10/blob/master/Include/10.0.10240.0/shared/guiddef.h
func CoCreateInstance(clsid *GUID, outer *IUnknown, clsContext CLSCTX, iid *GUID, object unsafe.Pointer) (HRESULT, error) {
	ret, _, err := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(clsid)),
		uintptr(unsafe.Pointer(outer)),
		uintptr(clsContext),
		uintptr(unsafe.Pointer(iid)),
		uintptr(object),
	)

	if ret == S_OK {
		return HRESULT(ret), nil
	}
	return HRESULT(ret), err
}

func CoTaskMemFree(object unsafe.Pointer) {
	procCoTaskMemFree.Call(uintptr(object))
}
