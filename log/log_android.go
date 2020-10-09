// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build android
package log

/*
#cgo LDFLAGS: -llog

#include <stdlib.h>
#include <android/log.h>

void log_print(int level, char* tag, char* msg){
	__android_log_print(level, tag, msg);
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
		level:  VERBOSE,
		logs:   make(chan string, lBufferSize),
	}

	return me
}

func (me *logger) output(depth int, level Level, levelTag string, tag string, format string, msg ...interface{}) {
	ctag := C.CString(tag)
	str := fmt.Sprintf(format, msg...)
	cmsg := C.CString(str)
	C.log_print(C.int(level), ctag, cmsg)
	C.free(unsafe.Pointer(ctag))
	C.free(unsafe.Pointer(cmsg))
}
