// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

char* nux_NSURL_path(uintptr_t nsurl){
	return (char*) [[(NSURL*)nsurl path] UTF8String];
}

uintptr_t nux_NSURL_fileURLWithPath(char* path){
	return (uintptr_t)[NSURL fileURLWithPath: [NSString stringWithUTF8String:path]];
}
*/
import "C"
import "unsafe"

func NSURLfileURLWithPath(path string) NSURL {
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	return NSURL(C.nux_NSURL_fileURLWithPath(cstr))
}

func (me NSURL) Path() string {
	return C.GoString(C.nux_NSURL_path(C.uintptr_t(me)))
}
