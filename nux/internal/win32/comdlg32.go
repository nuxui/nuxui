// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"syscall"
	"unsafe"
)

var (
	modcomdlg32 = syscall.NewLazyDLL("comdlg32.dll")

	procGetSaveFileName      = modcomdlg32.NewProc("GetSaveFileNameW")
	procGetOpenFileName      = modcomdlg32.NewProc("GetOpenFileNameW")
	procCommDlgExtendedError = modcomdlg32.NewProc("CommDlgExtendedError")
)

func GetSaveFileName(ofn *OPENFILENAME) bool {
	ret, _, _ := procGetSaveFileName.Call(uintptr(unsafe.Pointer(ofn)))
	return ret != 0
}

func GetOpenFileName(ofn *OPENFILENAME) bool {
	ret, _, _ := procGetOpenFileName.Call(uintptr(unsafe.Pointer(ofn)))
	return ret != 0
}

func CommDlgExtendedError() uint {
	ret, _, _ := procCommDlgExtendedError.Call()
	return uint(ret)
}
