// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && cairo

package nux

/*
#cgo LDFLAGS: -limm32
#cgo pkg-config: cairo
#cgo pkg-config: pango
#cgo pkg-config: pangocairo
#cgo pkg-config: gobject-2.0
#cgo pkg-config: libjpeg

#include <cairo/cairo.h>
#include <cairo/cairo-win32.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <pango/pangocairo.h>

#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdio.h>

#include <windows.h>
#include <windowsx.h>

void measureText(cairo_t* cr, char* fontFamily, int fontWeight, int fontSize,
	char* text, int width, int height, int* outWidth, int* outHeight){
	PangoLayout *layout;
	PangoFontDescription *font_description;

	font_description = pango_font_description_new ();
	pango_font_description_set_family (font_description, fontFamily);
	// pango_font_description_set_family_static (font_description, "Apple Color Emoji");
	// pango_font_description_set_weight (font_description, fontWeight);
	// pango_font_description_set_absolute_size (font_description, fontSize * PANGO_SCALE);
	// pango_font_description_set_stretch(font_description, PANGO_STRETCH_NORMAL);
	// pango_font_description_set_style(font_description, PANGO_STYLE_ITALIC);
	// pango_font_description_set_variant(font_description, PANGO_VARIANT_SMALL_CAPS);
	// pango_font_description_set_gravity(font_description, PANGO_GRAVITY_NORTH);

	layout = pango_cairo_create_layout (cr);
	pango_layout_set_font_description (layout, font_description);
	pango_layout_set_width (layout, width * PANGO_SCALE);
	pango_layout_set_height (layout, height * PANGO_SCALE);
	pango_layout_set_text (layout, text, -1);
	pango_layout_set_wrap (layout, PANGO_WRAP_WORD_CHAR);
	// pango_layout_set_justify(layout, TRUE);
	// pango_layout_set_indent(layout, 4);
	// pango_layout_set_markup(layout, "*", 10);
	// pango_layout_set_single_paragraph_mode(layout, TRUE);
	// pango_layout_set_alignment(layout,PANGO_ALIGN_RIGHT);

	pango_layout_get_size(layout, outWidth, outHeight);

	pango_font_description_free (font_description);
	g_object_unref (layout);
}

void drawText(cairo_t* cr, char* fontFamily, int fontWeight, int fontSize,
	char* text, int width, int height){
	PangoLayout *layout;
	PangoFontDescription *font_description;

	font_description = pango_font_description_new ();
	pango_font_description_set_family (font_description, fontFamily);
	// pango_font_description_set_family_static (font_description, "Apple Color Emoji");
	// pango_font_description_set_weight (font_description, fontWeight);
	// pango_font_description_set_absolute_size (font_description, fontSize * PANGO_SCALE);
	// pango_font_description_set_stretch(font_description, PANGO_STRETCH_NORMAL);
	// pango_font_description_set_style(font_description, PANGO_STYLE_ITALIC);
	// pango_font_description_set_variant(font_description, PANGO_VARIANT_SMALL_CAPS);
	// pango_font_description_set_gravity(font_description, PANGO_GRAVITY_NORTH);


	layout = pango_cairo_create_layout (cr);
	pango_layout_set_font_description (layout, font_description);
	pango_layout_set_width (layout, width * PANGO_SCALE);
	pango_layout_set_height (layout, height * PANGO_SCALE);
	pango_layout_set_text (layout, text, -1);
	pango_layout_set_wrap (layout, PANGO_WRAP_WORD_CHAR);
	// pango_layout_set_justify(layout, TRUE);
	// pango_layout_set_indent(layout, 4);
	// pango_layout_set_markup(layout, "*", 10);
	// pango_layout_set_single_paragraph_mode(layout, TRUE);
	// pango_layout_set_alignment(layout,PANGO_ALIGN_RIGHT);

	pango_cairo_show_layout (cr, layout);

	pango_font_description_free (font_description);
	g_object_unref (layout);
}

HDC toHDC(uintptr_t hdc){
	return (HDC)hdc;
}

*/
import "C"

import (
	"unsafe"

	"nuxui.org/nuxui/log"
)

type canvas struct {
	surface *C.cairo_surface_t
	ptr     *C.cairo_t
}

func newCanvas(hdc uintptr) *canvas {
	me := &canvas{}
	me.surface = C.cairo_win32_surface_create(C.toHDC(C.uintptr_t(hdc)))
	me.ptr = C.cairo_create(me.surface)
	return me
}

func (me *canvas) ResetClip() {
}

func (me *canvas) Save() {
	C.cairo_save(me.ptr)
}

func (me *canvas) Restore() {
	C.cairo_restore(me.ptr)
}

func (me *canvas) Translate(x, y float32) {
	C.cairo_translate(me.ptr, C.double(x), C.double(y))
}

func (me *canvas) Scale(x, y float32) {
	C.cairo_scale(me.ptr, C.double(x), C.double(y))
}

func (me *canvas) Rotate(angle float32) {
	C.cairo_rotate(me.ptr, C.double(_RADIAN*angle))
}

func (me *canvas) Skew(x, y float32) {
	// TODO::
}

func (me *canvas) Transform(a, b, c, d, e, f float32) {
	// TODO::
}

func (me *canvas) SetMatrix(matrix Matrix) {
	// TODO::
}

func (me *canvas) GetMatrix() Matrix {
	// TODO::
	return Matrix{}
}

func (me *canvas) ClipRect(x, y, width, height float32) {
	if width < 0 || height < 0 {
		log.Fatal("nuxui", "invalid rect for clip")
	}
	C.cairo_rectangle(me.ptr, C.double(x), C.double(y), C.double(width), C.double(height))
	C.cairo_clip(me.ptr)
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	// TODO::
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}
func (me *canvas) SetAlpha(alpha float32) {
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	if width < 0 || height < 0 {
		return
	}

	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		// C.cairo_identity_matrix(me.ptr)
		me.Save()
		me.Translate(0.5, 0.5)

	}

	C.cairo_rectangle(me.ptr, C.double(x), C.double(y), C.double(width), C.double(height))
	me.drawPaint(paint)

	if fix {
		me.Restore()
	}
}

func (me *canvas) DrawRoundRect(x, y, width, height float32, rLT, rRT, rRB, rLB float32, paint Paint) {
	if width < 0 || height < 0 {
		return
	}

	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		// C.cairo_identity_matrix(me.ptr)
		me.Save()
		me.Translate(0.5, 0.5)
	}

	C.cairo_new_sub_path(me.ptr)
	C.cairo_arc(me.ptr, C.double(x+width-rRT), C.double(y+rRT), C.double(rRT), -90*_RADIAN, 0)
	C.cairo_arc(me.ptr, C.double(x+width-rRB), C.double(y+height-rRB), C.double(rRB), 0, 90*_RADIAN)
	C.cairo_arc(me.ptr, C.double(x+rLB), C.double(y+height-rLB), C.double(rLB), 90*_RADIAN, 180*_RADIAN)
	C.cairo_arc(me.ptr, C.double(x+rLT), C.double(y+rLT), C.double(rLT), 180*_RADIAN, 270*_RADIAN)
	C.cairo_close_path(me.ptr)
	me.drawPaint(paint)

	if fix {
		me.Restore()
	}
}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {
	if useCenter {
		// TODO
		C.cairo_arc(me.ptr, C.double(x), C.double(y), C.double(radius), C.double(startAngle*_RADIAN), C.double(endAngle*_RADIAN))
		me.drawPaint(paint)
	} else {
		C.cairo_arc(me.ptr, C.double(x), C.double(y), C.double(radius), C.double(startAngle*_RADIAN), C.double(endAngle*_RADIAN))
		me.drawPaint(paint)
	}
}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
	me.Save()
	var centerX, centerY, scaleX, scaleY, radius float32
	if width > height {
		centerX = x + width/2.0
		centerY = y + width/2.0
		scaleX = 1.0
		scaleY = height / width
		radius = width / 2.0
	} else {
		centerX = x + height/2.0
		centerY = y + height/2.0
		scaleX = width / height
		scaleY = 1.0
		radius = height / 2.0
	}

	C.cairo_scale(me.ptr, C.double(scaleX), C.double(scaleY))
	C.cairo_arc(me.ptr, C.double(centerX), C.double(centerY), C.double(radius), C.double(0), C.double(_PI2))
	me.drawPaint(paint)
	me.Restore()
}

func (me *canvas) DrawPath(path Path) {
	// TODO::
}

func (me *canvas) drawPaint(paint Paint) {
	a, r, g, b := paint.Color().ARGBf()
	// C.cairo_fill_preserve(me.ptr)
	C.cairo_set_source_rgba(me.ptr, C.double(r), C.double(g), C.double(b), C.double(a))
	C.cairo_set_line_width(me.ptr, C.double(paint.Width()))
	switch paint.Style() {
	case PaintStyle_Stroke:
		C.cairo_stroke(me.ptr)
	case PaintStyle_Fill:
		C.cairo_fill(me.ptr)
	case PaintStyle_Both:
		C.cairo_stroke(me.ptr)
		C.cairo_fill(me.ptr)
	}
}

func (me *canvas) DrawColor(color Color) {
	a, r, g, b := color.ARGBf()
	C.cairo_set_source_rgba(me.ptr, C.double(r), C.double(g), C.double(b), C.double(a))
	C.cairo_paint(me.ptr)
}

func (me *canvas) DrawImage(img Image) {
	C.cairo_set_source_surface(me.ptr, img.(*nativeImage).ptr, 0, 0)
	C.cairo_paint(me.ptr)
}

func (me *canvas) DrawText(text string, width, height float32, paint Paint) {
	cfamily := C.CString("")
	ctext := C.CString(text)
	a, r, g, b := paint.Color().ARGBf()
	C.cairo_set_source_rgba(me.ptr, C.double(r), C.double(g), C.double(b), C.double(a))
	C.drawText(me.ptr, cfamily, C.int(1), C.int(paint.TextSize()), ctext, C.int(width), C.int(height))
	me.drawPaint(paint)
	C.free(unsafe.Pointer(cfamily))
	C.free(unsafe.Pointer(ctext))
}

func (me *canvas) Flush() {
	C.cairo_surface_flush(me.surface)
}

func (me *canvas) Destroy() {
	C.cairo_destroy(me.ptr)
}

var canvas4measure *C.cairo_t = C.cairo_create(nil)

func (me *paint) MeasureText(text string, width, height float32) (outWidth float32, outHeight float32) {
	if text == "" {
		return 0, 0
	}

	cfamily := C.CString("")
	ctext := C.CString(text)
	var w, h C.int
	C.measureText(canvas4measure, cfamily, C.int(1), C.int(me.TextSize()), ctext, C.int(width), C.int(height), &w, &h)
	outWidth = float32(int32(float64(w)/float64(C.PANGO_SCALE) + 0.99999))
	outHeight = float32(int32(float64(h)/float64(C.PANGO_SCALE) + 0.99999))
	C.free(unsafe.Pointer(cfamily))
	C.free(unsafe.Pointer(ctext))
	return
	return
}

func (me *paint) CharacterIndexForPoint(text string, width, height float32, x, y float32) uint32 {
	return 0
}
