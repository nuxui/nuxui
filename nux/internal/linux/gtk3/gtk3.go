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
#include <libintl.h>

gpointer nux_gtk_file_chooser_dialog_new(char* title, GtkWindow* parent, GtkFileChooserAction action){
	GtkWidget *dialog = gtk_file_chooser_dialog_new((gchar*)title, parent, action,
		dgettext("gtk30", "_Cancel"), GTK_RESPONSE_CANCEL,
		dgettext("gtk30", "_Open"), GTK_RESPONSE_ACCEPT,
		NULL);
	return (gpointer)dialog;
}

*/
import "C"
import "unsafe"

func InitCheck() bool {
	return C.gtk_init_check(nil, nil) > 0
}

// https://docs.gtk.org/gtk3/iface.FileChooser.html
func FileChooserDialogNew(title string, parent GtkWindow, action GtkFileChooserAction) GtkFileChooser {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	return GtkFileChooser(C.nux_gtk_file_chooser_dialog_new(cstr, (*C.GtkWindow)(unsafe.Pointer(parent)), C.GtkFileChooserAction(action)))
}

func (me GtkFileChooser) Run() GtkResponseType {
	return GtkResponseType(C.gtk_dialog_run((*C.GtkDialog)(unsafe.Pointer(me))))
}

func (me GtkFileChooser) GetFilename() string {
	cstr := C.gtk_file_chooser_get_filename((*C.GtkFileChooser)(unsafe.Pointer(me)))
	defer C.g_free(C.gpointer(cstr))
	return C.GoString(cstr)
}

func (me GtkFileChooser) GetFilenames() (fileNames []string) {
	filelist := GSList(unsafe.Pointer(C.gtk_file_chooser_get_filenames((*C.GtkFileChooser)(unsafe.Pointer(me)))))
	defer filelist.Free()

	len := filelist.Length()
	for i := uint32(0); i != len; i++ {
		name := filelist.DataAt(i)
		str := C.GoString((*C.char)(unsafe.Pointer(name)))
		if str != "" {
			fileNames = append(fileNames, str)
		}
		name.Free()
	}
	return
}

func (me GtkFileChooser) SetCurrentFolder(folder string) {
	cstr := C.CString(folder)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_file_chooser_set_current_folder((*C.GtkFileChooser)(unsafe.Pointer(me)), cstr)
}

func (me GtkFileChooser) SetCurrentName(name string) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_file_chooser_set_current_name((*C.GtkFileChooser)(unsafe.Pointer(me)), cstr)
}

func (me GtkFileChooser) SetSelectMultiple(enable bool) {
	var b C.gboolean
	if enable {
		b = 1
	}
	C.gtk_file_chooser_set_select_multiple((*C.GtkFileChooser)(unsafe.Pointer(me)), b)
}

func (me GtkFileChooser) SetCreateFolders(enable bool) {
	var b C.gboolean
	if enable {
		b = 1
	}
	C.gtk_file_chooser_set_create_folders((*C.GtkFileChooser)(unsafe.Pointer(me)), b)
}

func (me GtkFileChooser) SetDoOverwriteConfirmation(enable bool) {
	var b C.gboolean
	if enable {
		b = 1
	}
	C.gtk_file_chooser_set_do_overwrite_confirmation((*C.GtkFileChooser)(unsafe.Pointer(me)), b)
}

func (me GtkFileChooser) SetFilter(filter GtkFileFilter) {
	C.gtk_file_chooser_set_filter((*C.GtkFileChooser)(unsafe.Pointer(me)), (*C.GtkFileFilter)(unsafe.Pointer(filter)))
}

func (me GtkFileChooser) AddFilter(filter GtkFileFilter) {
	C.gtk_file_chooser_add_filter((*C.GtkFileChooser)(unsafe.Pointer(me)), (*C.GtkFileFilter)(unsafe.Pointer(filter)))
}

func (me GtkFileChooser) Close() {
	GtkWidget(me).Destroy()
	/* The Destroy call itself isn't enough to remove the dialog from the screen; apparently
	** that happens once the GTK main loop processes some further events. But if we're
	** in a non-GTK app the main loop isn't running, so we empty the event queue before
	** returning from the dialog functions.
	** Not sure how this interacts with an actual GTK app... */
	for C.gtk_events_pending() != 0 {
		C.gtk_main_iteration()
	}
}

func (me GtkWidget) Destroy() {
	C.gtk_widget_destroy((*C.GtkWidget)(unsafe.Pointer(me)))
}

func FileFilterNew() GtkFileFilter {
	return GtkFileFilter(C.gpointer(C.gtk_file_filter_new()))
}

func (me GtkFileFilter) SetName(name string) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_file_filter_set_name((*C.GtkFileFilter)(unsafe.Pointer(me)), cstr)
}

func (me GtkFileFilter) AddPattern(pattern string) {
	cstr := C.CString(pattern)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_file_filter_add_pattern((*C.GtkFileFilter)(unsafe.Pointer(me)), cstr)
}
