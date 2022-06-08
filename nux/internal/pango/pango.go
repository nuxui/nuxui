// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pango

/*
#cgo pkg-config: pango
#cgo pkg-config: pangocairo

#include <pango/pangocairo.h>

*/
import "C"

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/cairo"
	"unsafe"
)

const Scale = C.PANGO_SCALE

type FontDescription C.PangoFontDescription

func FontDescriptionNew() *FontDescription {
	return (*FontDescription)(C.pango_font_description_new())
}

func (me *FontDescription) Free() {
	C.pango_font_description_free((*C.PangoFontDescription)(me))
}

func (me *FontDescription) SetSize(size int32) {
	C.pango_font_description_set_size((*C.PangoFontDescription)(me), C.gint(size))
}

func (me *FontDescription) Size() int32 {
	return int32(C.pango_font_description_get_size((*C.PangoFontDescription)(me)))
}

func (me *FontDescription) SetFamily(family string) {
	cstr := C.CString(family)
	defer C.free(unsafe.Pointer(cstr))
	C.pango_font_description_set_family((*C.PangoFontDescription)(me), cstr)
}

func (me *FontDescription) Family() string {
	return C.GoString(C.pango_font_description_get_family((*C.PangoFontDescription)(me)))
}

func (me *FontDescription) SetWeight(weight Weight) {
	C.pango_font_description_set_weight((*C.PangoFontDescription)(me), C.PangoWeight(weight))
}

func (me *FontDescription) Weight() Weight {
	return Weight(C.pango_font_description_get_weight((*C.PangoFontDescription)(me)))
}

type Layout C.PangoLayout

func CairoCreateLayout(cr *cairo.Cairo) *Layout {
	return (*Layout)(C.pango_cairo_create_layout((*C.cairo_t)(cr)))
}

func (me *Layout) Free() {
	C.g_object_unref(C.gpointer((*C.PangoLayout)(me)))
}

func (me *Layout) CairoUpdateLayout(cr *cairo.Cairo) {
	C.pango_cairo_update_layout((*C.cairo_t)(cr), (*C.PangoLayout)(me))
}

func (me *Layout) SetFontDescription(fd *FontDescription) {
	C.pango_layout_set_font_description((*C.PangoLayout)(me), (*C.PangoFontDescription)(fd))
}

func (me *Layout) SetWidth(width int32) {
	C.pango_layout_set_width((*C.PangoLayout)(me), C.int(width))
}

func (me *Layout) SetHeight(height int32) {
	C.pango_layout_set_height((*C.PangoLayout)(me), C.int(height))
}

func (me *Layout) SetText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.pango_layout_set_text((*C.PangoLayout)(me), cstr, C.int(C.strlen(cstr)))
}

func (me *Layout) SetWrap(mode WrapMode) {
	C.pango_layout_set_wrap((*C.PangoLayout)(me), C.PangoWrapMode(mode))
}

func (me *Layout) GetSize() (width, height int32) {
	var w, h C.int
	C.pango_layout_get_size((*C.PangoLayout)(me), &w, &h)
	return int32(w), int32(h)
}

func (me *Layout) GetPixelSize() (width, height int32) {
	var w, h C.int
	C.pango_layout_get_pixel_size((*C.PangoLayout)(me), &w, &h)
	return int32(w), int32(h)
}

func (me *Layout) CairoShowLayout(cr *cairo.Cairo) {
	C.pango_cairo_show_layout((*C.cairo_t)(cr), (*C.PangoLayout)(me))
}

// https://docs.gtk.org/Pango/method.Layout.xy_to_index.html
func (me *Layout) XYtoIndex(x, y int32) (index, trailing int32, hit bool) {
	var i, t C.int
	hit = C.pango_layout_xy_to_index((*C.PangoLayout)(me), C.int(x), C.int(y), &i, &t) > 0
	index = int32(i)
	trailing = int32(t)
	return
}

func TestLayout() {
	return
	text := "hello"
	family := "sans-serif"
	cstr := C.CString(text)
	cfamily := C.CString(family)
	// defer C.free(unsafe.Pointer(cstr))

	cr := C.cairo_create(nil)
	fd := C.pango_font_description_new()
	C.pango_font_description_set_family(fd, cfamily)
	C.pango_font_description_set_size(fd, 14*C.PANGO_SCALE)
	// C.pango_font_description_set_size(fd, 14)
	// C.pango_font_description_set_absolute_size(fd, 14)

	layout := C.pango_cairo_create_layout(cr)
	C.pango_layout_set_font_description(layout, fd)
	C.pango_layout_set_wrap(layout, 1)
	C.pango_layout_set_width(layout, 30)
	C.pango_layout_set_height(layout, 30)
	// C.pango_layout_set_width(layout, 30*C.PANGO_SCALE)
	// C.pango_layout_set_height(layout, 30*C.PANGO_SCALE)
	C.pango_layout_set_text(layout, cstr, -1)
	var w, h C.int
	C.pango_layout_get_size(layout, &w, &h)
	w2 := float32(w) / float32(C.PANGO_SCALE)
	h2 := float32(h) / float32(C.PANGO_SCALE)
	log.E("nuxui", "TestLayout w:%d, h:%d, w2:%f h2:%f", w, h, w2, h2)
	var w3, h3 C.int
	C.pango_layout_get_pixel_size(layout, &w3, &h3)
	log.E("nuxui", "TestLayout pixel w3:%d h3:%d, scale=%d", w3, h3, C.PANGO_SCALE)
}
