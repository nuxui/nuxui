// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

import (
	"syscall"
	"unsafe"
)

type COMPOSITIONFORM struct {
	Style      int32
	CurrentPos POINT
	Area       RECT
}

// bit field for IMC_SETCOMPOSITIONWINDOW, IMC_SETCANDIDATEWINDOW
const (
	CFS_DEFAULT        = 0x0000
	CFS_RECT           = 0x0001
	CFS_POINT          = 0x0002
	CFS_FORCE_POSITION = 0x0020
	CFS_CANDIDATEPOS   = 0x0040
	CFS_EXCLUDE        = 0x0080
)

// parameter of ImmGetCompositionString
const (
	GCS_COMPREADSTR      = 0x0001
	GCS_COMPREADATTR     = 0x0002
	GCS_COMPREADCLAUSE   = 0x0004
	GCS_COMPSTR          = 0x0008
	GCS_COMPATTR         = 0x0010
	GCS_COMPCLAUSE       = 0x0020
	GCS_CURSORPOS        = 0x0080
	GCS_DELTASTART       = 0x0100
	GCS_RESULTREADSTR    = 0x0200
	GCS_RESULTREADCLAUSE = 0x0400
	GCS_RESULTSTR        = 0x0800
	GCS_RESULTCLAUSE     = 0x1000
)

var (
	modimm32                     = syscall.NewLazyDLL("imm32.dll")
	procImmGetContext            = modimm32.NewProc("ImmGetContext")
	procImmSetCompositionWindow  = modimm32.NewProc("ImmSetCompositionWindow")
	procImmReleaseContext        = modimm32.NewProc("ImmReleaseContext")
	procImmGetCompositionStringW = modimm32.NewProc("ImmGetCompositionStringW")
)

func ImmGetContext(hwnd uintptr) (himc uintptr) {
	himc, _, _ = procImmGetContext.Call(hwnd)
	return
}

func ImmSetCompositionWindow(himc uintptr, form *COMPOSITIONFORM) (err error) {
	r0, _, err := procImmSetCompositionWindow.Call(himc, uintptr(unsafe.Pointer(form)))
	if r0 != 0 {
		err = nil
	}
	return
}

func ImmReleaseContext(hwnd uintptr, himc uintptr) (err error) {
	r0, _, err := procImmReleaseContext.Call(hwnd, himc)
	if r0 != 0 {
		err = nil
	}
	return
}
func ImmGetCompositionStringW(himc uintptr, param int32, lpBuf uintptr, dwBufLen int32) (len int32) {
	r0, _, _ := procImmGetCompositionStringW.Call(
		himc,
		uintptr(param),
		lpBuf,
		uintptr(dwBufLen),
	)
	len = int32(r0)
	return
}
