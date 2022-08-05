// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

uintptr_t nux_NSSavePanel_shared(){
	return (uintptr_t)([NSSavePanel savePanel]);
}

int nux_NSSavePanel_runModal(uintptr_t nsSavePanel){
	[(NSSavePanel*)(nsSavePanel) setAllowsOtherFileTypes: YES];
	return (int)([(NSSavePanel*)(nsSavePanel) runModal]);
}

void nux_NSSavePanel_setDirectoryURL(uintptr_t nsSavePanel, uintptr_t url){
	[(NSSavePanel*)(nsSavePanel) setDirectoryURL: (NSURL*)url];
}

void nux_NSSavePanel_setNameFieldStringValue(uintptr_t nsSavePanel, char* value){
	[(NSSavePanel*)(nsSavePanel) setNameFieldStringValue: [NSString stringWithUTF8String:value]];
}

uintptr_t nux_NSSavePanel_URL(uintptr_t nsSavePanel){
	return (uintptr_t)([(NSSavePanel*)(nsSavePanel) URL]);
}

uintptr_t nux_NSOpenPanel_shared(){
	return (uintptr_t)([NSOpenPanel openPanel]);
}

void nux_NSOpenPanel_setCanChooseFiles(uintptr_t nsOpenPanel, int can){
	[(NSOpenPanel*)(nsOpenPanel) setCanChooseFiles:can>0];
}

void nux_NSOpenPanel_setCanChooseDirectories(uintptr_t nsOpenPanel, int can){
	[(NSOpenPanel*)(nsOpenPanel) setCanChooseDirectories:can>0];
}

void nux_NSOpenPanel_setAllowsMultipleSelection(uintptr_t nsOpenPanel, int allow){
	[(NSOpenPanel*)(nsOpenPanel) setAllowsMultipleSelection:allow>0];
}

void nux_NSOpenPanel_setCanCreateDirectories(uintptr_t nsOpenPanel, int can){
	[(NSOpenPanel*)(nsOpenPanel) setCanCreateDirectories:can>0];
}

void nux_NSOpenPanel_setAllowedContentTypes(uintptr_t nsOpenPanel, void* types, int count){
	NSArray<UTType*> *arr = [NSArray arrayWithObjects:(UTType**)types count:count];
	[(NSOpenPanel*)(nsOpenPanel) setAllowedContentTypes:arr];
}

uintptr_t nux_NSOpenPanel_URLs(uintptr_t nsOpenPanel){
	return (uintptr_t)[(NSOpenPanel*)(nsOpenPanel) URLs];
}

*/
import "C"
import "unsafe"

const (
	NSModalResponseOK         NSModalResponse = 1
	NSModalResponseCancel     NSModalResponse = 0
	NSModalResponseContinue   NSModalResponse = -1002
	NSModalResponseStop       NSModalResponse = -1000
	NSModalResponseAbort      NSModalResponse = -1001
	NSAlertFirstButtonReturn  NSModalResponse = 1000
	NSAlertSecondButtonReturn NSModalResponse = 1001
	NSAlertThirdButtonReturn  NSModalResponse = 1002
)

// -------------- NSSavePanel -----------------
func SharedNSSavePanel() NSSavePanel {
	return NSSavePanel(C.nux_NSSavePanel_shared())
}

func (me NSSavePanel) RunModal() NSModalResponse {
	return NSModalResponse(C.nux_NSSavePanel_runModal(C.uintptr_t(me)))
}

func (me NSSavePanel) URL() NSURL {
	return NSURL(C.nux_NSSavePanel_URL(C.uintptr_t(me)))
}

func (me NSSavePanel) SetDirectoryURL(url NSURL) {
	C.nux_NSSavePanel_setDirectoryURL(C.uintptr_t(me), C.uintptr_t(url))
}

func (me NSSavePanel) SetNameFieldStringValue(value string) {
	cstr := C.CString(value)
	defer C.free(unsafe.Pointer(cstr))
	C.nux_NSSavePanel_setNameFieldStringValue(C.uintptr_t(me), cstr)
}

// -------------- NSOpenPanel -----------------
func SharedNSOpenPanel() NSOpenPanel {
	return NSOpenPanel(C.nux_NSOpenPanel_shared())
}

func (me NSOpenPanel) SetCanChooseFiles(can bool) {
	b := C.int(0)
	if can {
		b = C.int(1)
	}
	C.nux_NSOpenPanel_setCanChooseFiles(C.uintptr_t(me), b)
}

func (me NSOpenPanel) SetCanChooseDirectories(can bool) {
	b := C.int(0)
	if can {
		b = C.int(1)
	}
	C.nux_NSOpenPanel_setCanChooseDirectories(C.uintptr_t(me), b)
}

func (me NSOpenPanel) SetAllowsMultipleSelection(allow bool) {
	b := C.int(0)
	if allow {
		b = C.int(1)
	}
	C.nux_NSOpenPanel_setAllowsMultipleSelection(C.uintptr_t(me), b)
}

func (me NSOpenPanel) SetAllowedContentTypes(types []UTType) {
	if len(types) == 0 {
		C.nux_NSOpenPanel_setAllowedContentTypes(C.uintptr_t(me), nil, 0)
	} else {
		C.nux_NSOpenPanel_setAllowedContentTypes(C.uintptr_t(me), unsafe.Pointer(&types[0]), C.int(len(types)))
	}
}

func (me NSOpenPanel) SetCanCreateDirectories(can bool) {
	b := C.int(0)
	if can {
		b = C.int(1)
	}
	C.nux_NSOpenPanel_setCanCreateDirectories(C.uintptr_t(me), b)
}

func (me NSOpenPanel) SetDirectoryURL(url NSURL) {
	C.nux_NSSavePanel_setDirectoryURL(C.uintptr_t(me), C.uintptr_t(url))
}

func (me NSOpenPanel) RunModal() NSModalResponse {
	return NSModalResponse(C.nux_NSSavePanel_runModal(C.uintptr_t(me)))
}

func (me NSOpenPanel) URLs() (urls []NSURL) {
	arr := NSArray(C.nux_NSOpenPanel_URLs(C.uintptr_t(me)))
	count := arr.Count()
	for i := 0; i != count; i++ {
		url := NSURL(arr.ObjectAtIndex(i))
		urls = append(urls, url)
	}
	return
}
