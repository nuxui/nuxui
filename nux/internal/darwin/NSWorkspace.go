// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

uintptr_t nux_NSWorkspace_shared(){
	return (uintptr_t)([NSWorkspace sharedWorkspace]);
}

void nux_NSWorkspace_openURL(uintptr_t nsWorkspace, char* url){
	[(NSWorkspace*)(nsWorkspace) openURL: [NSURL fileURLWithPath: [NSString stringWithUTF8String:url]]];
}

void nux_NSWorkspace_activateFileViewerSelectingURLs(uintptr_t nsWorkspace, void* urls, int count){
	NSArray<NSURL*> *arr = [NSArray arrayWithObjects:(NSURL**)urls count:count];
	[(NSWorkspace*)(nsWorkspace) activateFileViewerSelectingURLs:arr];
}
*/
import "C"
import "unsafe"

func SharedNSWorkspace() NSWorkspace {
	return NSWorkspace(C.nux_NSWorkspace_shared())
}

func (me NSWorkspace) OpenURL(url string) {
	cstr := C.CString(url)
	defer C.free(unsafe.Pointer(cstr))
	C.nux_NSWorkspace_openURL(C.uintptr_t(me), cstr)
}

func (me NSWorkspace) ActivateFileViewerSelectingURLs(urls []NSURL) {
	if len(urls) == 0 {
		return
	}
	C.nux_NSWorkspace_activateFileViewerSelectingURLs(C.uintptr_t(me), unsafe.Pointer(&urls[0]), C.int(len(urls)))
}
