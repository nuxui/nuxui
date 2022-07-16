// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go:build (linux && !android)

package linux

/*
#include <stdint.h>
#include <sys/syscall.h>
#include <unistd.h>
#include <locale.h>
#include <stdlib.h>
#include <stdio.h>

uint64_t currentThreadID(){
	return (uint64_t)syscall(SYS_gettid);
}

*/
import "C"
import "unsafe"

func CurrentThreadID() uint64 {
	return uint64(C.currentThreadID())
}

func SetLocale(category Category, locale string) {
	clocale := C.CString(locale)
	defer C.free(unsafe.Pointer(clocale))
	C.setlocale(C.int(category), clocale)
}

func GetLocale(category Category) string {
	loc := C.setlocale(C.int(category), nil)
	return C.GoString(loc)
}

func System(cmd string) {
	cstr := C.CString(cmd)
	defer C.free(unsafe.Pointer(cstr))
	C.system(cstr)
}
