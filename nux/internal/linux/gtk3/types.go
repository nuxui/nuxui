// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gtk3

/*
#cgo pkg-config: gtk+-3.0
#cgo LDFLAGS: -lX11

#include <X11/Xlib.h>
#include <gtk/gtk.h>
#include <stdlib.h>


*/
import "C"
import "unsafe"

type GPointer uintptr
type GtkWidget GPointer
type GtkDialog GPointer
type GtkWindow GPointer
type GtkFileChooser GPointer
type GSList GPointer

func (me GSList) Free() {
	C.g_slist_free((*C.GSList)(unsafe.Pointer(me)))
}

func (me GSList) Length() uint32 {
	return uint32(C.g_slist_length((*C.GSList)(unsafe.Pointer(me))))
}

func (me GSList) DataAt(index uint32) GPointer {
	return GPointer(C.g_slist_nth_data((*C.GSList)(unsafe.Pointer(me)), C.guint(index)))
}

func (me GPointer) Free() {
	C.g_free(C.gpointer(me))
}
