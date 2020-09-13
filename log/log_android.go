// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build android
package log

/*
#cgo LDFLAGS: -llog

#include <stdlib.h>
#include <android/log.h>

void log_print(int level, char* tag, char* fmt, char* msg){
	__android_log_print(level, tag, fmt, msg);
}

*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

func new(out io.Writer, prefix string, flags int, depth int) Logger {
	me := &logger{
		depth:  depth,
		out:    out,
		flags:  flags,
		prefix: prefix,
	}

	return me
}

func (me *logger) output(depth int, level int, levelTag string, tag string, format string, msg ...interface{}) {
	ctag := C.CString(tag)
	cfmt := C.CString("%s")
	str := fmt.Sprintf(format, msg...)
	cmsg := C.CString(str)
	C.log_print(C.int(level), ctag, cfmt, cmsg)
	C.free(unsafe.Pointer(ctag))
	C.free(unsafe.Pointer(cfmt))
	C.free(unsafe.Pointer(cmsg))
}
