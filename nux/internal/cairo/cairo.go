// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cairo

/*
#cgo pkg-config: cairo
#cgo pkg-config: libjpeg

#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>

#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdio.h>

#include "jpeg.h"
*/
import "C"

import (
	"unsafe"
)

const (
	CAIRO_SVG_UNIT_PX             = C.CAIRO_SVG_UNIT_PX
	CAIRO_FORMAT_INVALID   Format = -1
	CAIRO_FORMAT_ARGB32    Format = 0
	CAIRO_FORMAT_RGB24     Format = 1
	CAIRO_FORMAT_A8        Format = 2
	CAIRO_FORMAT_A1        Format = 3
	CAIRO_FORMAT_RGB16_565 Format = 4
	CAIRO_FORMAT_RGB30     Format = 5
)

type Surface C.cairo_surface_t
type Cairo C.cairo_t
type Format C.cairo_format_t
type SVGUnit C.cairo_svg_unit_t
type Path C.cairo_path_t

func ImageSurfaceCreate(format Format, width, height int32) *Surface {
	return (*Surface)(C.cairo_image_surface_create(C.cairo_format_t(format), C.int(width), C.int(height)))
}

func ImageSurfaceCreateFromPNG(filename string) *Surface {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	return (*Surface)(C.cairo_image_surface_create_from_png(cstr))
}

func ImageSurfaceCreateFromJPEG(filename string) *Surface {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	return (*Surface)(C.cairo_image_surface_create_from_jpeg(cstr))
}

func ImageSurfaceCreateForData(data *uint8, format Format, width, height, stride int32) *Surface {
	return (*Surface)(C.cairo_image_surface_create_for_data((*C.uchar)(data), C.cairo_format_t(format), C.int(width),
		C.int(height), C.int(stride)))
}

func SVGSurfaceCreate(filename string, width, height float32) *Surface {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	return (*Surface)(C.cairo_svg_surface_create(cstr, C.double(width), C.double(height)))
}

func ImageSurfaceGetWidth(surface *Surface) int32 {
	return int32(C.cairo_image_surface_get_width((*C.cairo_surface_t)(surface)))
}

func ImageSurfaceGetHeight(surface *Surface) int32 {
	return int32(C.cairo_image_surface_get_height((*C.cairo_surface_t)(surface)))
}

func ImageSurfaceGetStride(surface *Surface) int32 {
	return int32(C.cairo_image_surface_get_stride((*C.cairo_surface_t)(surface)))
}

func ImageSurfaceGetData(surface *Surface) *byte {
	return (*byte)(unsafe.Pointer(C.cairo_image_surface_get_data((*C.cairo_surface_t)(surface))))
}

func ImageSurfaceGetFormat(surface *Surface) Format {
	return Format(C.cairo_image_surface_get_format((*C.cairo_surface_t)(surface)))
}

func SVGSurfaceSetDocumentUnit(surface *Surface, unit SVGUnit) {
	C.cairo_svg_surface_set_document_unit((*C.cairo_surface_t)(surface), C.cairo_svg_unit_t(unit))
}

func (me *Surface) Flush() {
	C.cairo_surface_flush((*C.cairo_surface_t)(me))
}

func (me *Surface) Finish() {
	C.cairo_surface_finish((*C.cairo_surface_t)(me))
}

func (me *Surface) Destroy() {
	C.cairo_surface_destroy((*C.cairo_surface_t)(me))
}

func Create(surface *Surface) *Cairo {
	return (*Cairo)(C.cairo_create((*C.cairo_surface_t)(surface)))
}

func (me *Cairo) SetSourceSurface(surface *Surface, x, y float32) {
	C.cairo_set_source_surface((*C.cairo_t)(me), (*C.cairo_surface_t)(surface), C.double(x), C.double(y))
}

func (me *Cairo) SetSourceRGBA(r, g, b, a float32) {
	C.cairo_set_source_rgba((*C.cairo_t)(me), C.double(r), C.double(g), C.double(b), C.double(a))
}

func (me *Cairo) SetLineWidth(width float32) {
	C.cairo_set_line_width((*C.cairo_t)(me), C.double(width))
}

func (me *Cairo) Rectangle(x, y, width, height float32) {
	C.cairo_rectangle((*C.cairo_t)(me), C.double(x), C.double(y), C.double(width), C.double(height))
}

func (me *Cairo) Fill() {
	C.cairo_fill((*C.cairo_t)(me))
}

func (me *Cairo) Stroke() {
	C.cairo_stroke((*C.cairo_t)(me))
}

func (me *Cairo) Save() {
	C.cairo_save((*C.cairo_t)(me))
}

func (me *Cairo) Restore() {
	C.cairo_restore((*C.cairo_t)(me))
}

func (me *Cairo) Translate(x, y float32) {
	C.cairo_translate((*C.cairo_t)(me), C.double(x), C.double(y))
}

func (me *Cairo) Scale(x, y float32) {
	C.cairo_scale((*C.cairo_t)(me), C.double(x), C.double(y))
}

func (me *Cairo) Rotate(angle float32) {
	C.cairo_rotate((*C.cairo_t)(me), C.double(angle))
}

func (me *Cairo) Clip() {
	C.cairo_clip((*C.cairo_t)(me))
}

func (me *Cairo) Arc(xc, yc, radius, angle1, angle2 float32) {
	C.cairo_arc((*C.cairo_t)(me), C.double(xc), C.double(yc), C.double(radius), C.double(angle1), C.double(angle2))
}

func (me *Cairo) Paint() {
	C.cairo_paint((*C.cairo_t)(me))
}

func (me *Cairo) Destroy() {
	C.cairo_destroy((*C.cairo_t)(me))
}

func (me *Cairo) NewPath() {
	C.cairo_new_path((*C.cairo_t)(me))
}

func (me *Cairo) NewSubPath() {
	C.cairo_new_sub_path((*C.cairo_t)(me))
}

func (me *Cairo) ClosePath() {
	C.cairo_close_path((*C.cairo_t)(me))
}

func (me *Cairo) CopyPath() *Path {
	return (*Path)(C.cairo_copy_path((*C.cairo_t)(me)))
}

func (me *Cairo) CopyPathFlat() *Path {
	return (*Path)(C.cairo_copy_path_flat((*C.cairo_t)(me)))
}

func (me *Cairo) AppendPath(path *Path) {
	C.cairo_append_path((*C.cairo_t)(me), (*C.cairo_path_t)(path))
}

func (me *Path) Destroy() {
	C.cairo_path_destroy((*C.cairo_path_t)(me))
}

func (me *Cairo) HasCurrentPoint() bool {
	return C.cairo_has_current_point((*C.cairo_t)(me)) > 0
}

func (me *Cairo) MoveTo(x, y float32) {
	C.cairo_move_to((*C.cairo_t)(me), C.double(x), C.double(y))
}

func (me *Cairo) LineTo(x, y float32) {
	C.cairo_line_to((*C.cairo_t)(me), C.double(x), C.double(y))
}

func (me *Cairo) CurveTo(x1, y1, x2, y2, x3, y3 float32) {
	C.cairo_curve_to((*C.cairo_t)(me), C.double(x1), C.double(y1), C.double(x2), C.double(y2), C.double(x3), C.double(y3))
}
