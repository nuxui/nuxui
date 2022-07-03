// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

import (
	"syscall"
)

var (
	modkernel32            = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentThreadId = modkernel32.NewProc("GetCurrentThreadId")
)

// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentthreadid
func GetCurrentThreadId() uint32 {
	ret, _, _ := procGetCurrentThreadId.Call()
	return uint32(ret)
}
