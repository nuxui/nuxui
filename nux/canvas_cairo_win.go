// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package nux

/*
#cgo pkg-config: cairo
#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <cairo/cairo-win32.h>

#cgo pkg-config: pango
#cgo pkg-config: pangocairo
#cgo pkg-config: gobject-2.0
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

*/
import "C"

import (
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

const (
	FORMAT_INVALID   = -1
	FORMAT_ARGB32    = 0
	FORMAT_RGB24     = 1
	FORMAT_A8        = 2
	FORMAT_A1        = 3
	FORMAT_RGB16_565 = 4
	FORMAT_RGB30     = 5
)

const (
	ANTIALIAS_DEFAULT  = C.CAIRO_ANTIALIAS_DEFAULT
	ANTIALIAS_NONE     = C.CAIRO_ANTIALIAS_NONE
	ANTIALIAS_GRAY     = C.CAIRO_ANTIALIAS_GRAY
	ANTIALIAS_SUBPIXEL = C.CAIRO_ANTIALIAS_SUBPIXEL
	ANTIALIAS_FAST     = C.CAIRO_ANTIALIAS_FAST
	ANTIALIAS_GOOD     = C.CAIRO_ANTIALIAS_GOOD
	ANTIALIAS_BEST     = C.CAIRO_ANTIALIAS_BEST
)

const (
	PI     = 3.1415926535897932384626433832795028841971
	PI2    = PI * 2
	DEGREE = PI / 180.0
)

/////////////////////////////////////////////////////////////////////////////////////////
/////////////                        Surface                          ///////////////////
/////////////////////////////////////////////////////////////////////////////////////////

// Surface c
type Surface struct {
	surface *C.cairo_surface_t
	canvas  *canvas
}

func NewSurfaceFromData(data unsafe.Pointer, format, width, height, stride int) *Surface {
	s := C.cairo_image_surface_create_for_data((*C.uchar)(data), C.cairo_format_t(format),
		C.int(width), C.int(height), C.int(stride))

	cairo := C.cairo_create(s)
	return &Surface{surface: s, canvas: &canvas{ptr: cairo}}
}

func newSurfaceWin32(hdc C.HDC) *Surface {
	s := C.cairo_win32_surface_create(hdc)
	cairo := C.cairo_create(s)
	return &Surface{surface: s, canvas: &canvas{ptr: cairo}}
}

// func newSurfaceQuartzWithCGContext(context uintptr, width, height int32) *Surface {
// 	s := C.nux_cairo_quartz_surface_create_for_cg_context(C.uintptr_t(context), C.uint(width), C.uint(height))
// 	cairo := C.cairo_create(s)
// 	return &Surface{surface: s, canvas: &canvas{ptr: cairo}}
// }

// func newSurfaceQuartz(width, height int32) *Surface {
// 	s := C.cairo_quartz_surface_create(C.CAIRO_FORMAT_ARGB32, C.uint(width), C.uint(height))
// 	cairo := C.cairo_create(s)
// 	return &Surface{surface: s, canvas: &canvas{ptr: cairo}}
// }

func (me *Surface) WriteToPng(fileName string) {
	name := C.CString(fileName)
	defer C.free(unsafe.Pointer(name))
	C.cairo_surface_write_to_png(me.surface, name)
}

func (me *Surface) GetCanvas() Canvas {
	if me.canvas == nil {
		log.Fatal("nuxui", "create surface by call NewSurface...")
	}
	return me.canvas
}

func (me *Surface) Flush() {
	C.cairo_surface_flush(me.surface)
}

func (me *Surface) Destory() {
	me.Flush()
	C.cairo_surface_finish(me.surface)
	C.cairo_surface_destroy(me.surface)
	me.surface = nil
	me.canvas = nil
}

/////////////////////////////////////////////////////////////////////////////////////////
/////////////                        Canvas                           ///////////////////
/////////////////////////////////////////////////////////////////////////////////////////

func NewCanvas() Canvas {
	return &canvas{}
}

// Canvas c
type canvas struct {
	ptr *C.cairo_t
}

func (me *canvas) Save() {
	C.cairo_save(me.ptr)
}

func (me *canvas) Restore() {
	C.cairo_restore(me.ptr)
}

func (me *canvas) Translate(x, y int32) {
	C.cairo_translate(me.ptr, C.double(x), C.double(y))
}

func (me *canvas) TranslateF(x, y float32) {
	C.cairo_translate(me.ptr, C.double(x), C.double(y))
}

func (me *canvas) Scale(x, y int32) {
	C.cairo_scale(me.ptr, C.double(x), C.double(y))
}

func (me *canvas) ScaleF(x, y float32) {
	C.cairo_scale(me.ptr, C.double(x), C.double(y))
}

func (me *canvas) Rotate(angle int32) {
	C.cairo_rotate(me.ptr, C.double(angle))
}

func (me *canvas) RotateF(angle float32) {
	C.cairo_rotate(me.ptr, C.double(angle))
}

// TODO https://cairographics.org/manual/cairo-Transformations.html#cairo-transform
func (me *canvas) Transform() {
}

// TODO https://cairographics.org/manual/cairo-Transformations.html#cairo-set-matrix
func (me *canvas) SetMatrix() {
}

// TODO https://cairographics.org/manual/cairo-Transformations.html#cairo-get-matrix
func (me *canvas) GetMatrix() {
}

// TODO https://developer.android.com/reference/android/graphics/Canvas#skew(float,%20float)
func (me *canvas) Skew() {
}

func (me *canvas) ClipRect(left, top, right, bottom int32) {
	me.ClipRectF(float32(left), float32(top), float32(right), float32(bottom))
}

func (me *canvas) ClipRectF(left, top, right, bottom float32) {
	if right < left || bottom < top {
		log.Fatal("nuxui", "invalid rect for clip")
	}
	C.cairo_rectangle(me.ptr, C.double(left), C.double(top), C.double(right-left), C.double(bottom-top))
	C.cairo_clip(me.ptr)
}

func (me *canvas) ClipPath() {
	// TODO::
}

func (me *canvas) DrawRect(left, top, right, bottom int32, paint *Paint) {
	me.DrawRectF(float32(left), float32(top), float32(right), float32(bottom), paint)
}

func (me *canvas) DrawRectF(left, top, right, bottom float32, paint *Paint) {
	if right <= left || bottom <= top {
		return
	}

	fix := paint.Style == STROKE && int32(paint.Width)%2 != 0
	if fix {
		// C.cairo_identity_matrix(me.ptr)
		me.Save()
		me.TranslateF(-0.5, -0.5)

	}

	C.cairo_rectangle(me.ptr, C.double(left), C.double(top), C.double(right-left), C.double(bottom-top))
	me.drawPaint(paint)

	if fix {
		me.Restore()
	}
}

func (me *canvas) DrawArc(x, y, radius, angle1, angle2 float32, useCenter bool, paint *Paint) {
	if useCenter {
		// TODO
		C.cairo_arc(me.ptr, C.double(x), C.double(y), C.double(radius), C.double(angle1*DEGREE), C.double(angle2*DEGREE))
		me.drawPaint(paint)
	} else {
		C.cairo_arc(me.ptr, C.double(x), C.double(y), C.double(radius), C.double(angle1*DEGREE), C.double(angle2*DEGREE))
		me.drawPaint(paint)
	}
}

func (me *canvas) DrawOval(left, top, right, bottom int32, paint *Paint) {
	me.DrawOvalF(float32(left), float32(top), float32(right), float32(bottom), paint)
}

func (me *canvas) DrawOvalF(left, top, right, bottom float32, paint *Paint) {
	if left > right || top > bottom {
		return
	}

	me.Save()
	width := right - left
	height := bottom - top
	var centerX, centerY, scaleX, scaleY, radius float32
	if width > height {
		centerX = left + width/2.0
		centerY = top + width/2.0
		scaleX = 1.0
		scaleY = height / width
		radius = width / 2.0
	} else {
		centerX = left + height/2.0
		centerY = top + height/2.0
		scaleX = width / height
		scaleY = 1.0
		radius = height / 2.0
	}

	C.cairo_scale(me.ptr, C.double(scaleX), C.double(scaleY))
	C.cairo_arc(me.ptr, C.double(centerX), C.double(centerY), C.double(radius), C.double(0), C.double(PI2))
	me.drawPaint(paint)
	me.Restore()
}

func (me *canvas) DrawRoundRect(left, top, right, bottom, radius int32, paint *Paint) {
	me.DrawRoundRectF(float32(left), float32(top), float32(right), float32(bottom), float32(radius), paint)
}

func (me *canvas) DrawRoundRectF(left, top, right, bottom, radius float32, paint *Paint) {
	C.cairo_new_sub_path(me.ptr)
	C.cairo_arc(me.ptr, C.double(right-radius), C.double(top+radius), C.double(radius), -90*DEGREE, 0)
	C.cairo_arc(me.ptr, C.double(right-radius), C.double(bottom-radius), C.double(radius), 0, 90*DEGREE)
	C.cairo_arc(me.ptr, C.double(left+radius), C.double(bottom-radius), C.double(radius), 90*DEGREE, 180*DEGREE)
	C.cairo_arc(me.ptr, C.double(left+radius), C.double(top+radius), C.double(radius), 180*DEGREE, 270*DEGREE)
	C.cairo_close_path(me.ptr)
	me.drawPaint(paint)
}

func (me *canvas) drawPaint(paint *Paint) {
	a := float32((paint.Color>>24)&0xff) / 255
	r := float32((paint.Color>>16)&0xff) / 255
	g := float32((paint.Color>>8)&0xff) / 255
	b := float32((paint.Color)&0xff) / 255
	// C.cairo_fill_preserve(me.ptr)
	C.cairo_set_source_rgba(me.ptr, C.double(r), C.double(g), C.double(b), C.double(a))
	C.cairo_set_line_width(me.ptr, C.double(paint.Width))
	switch paint.Style {
	case STROKE:
		C.cairo_stroke(me.ptr)
	case FILL:
		C.cairo_fill(me.ptr)
	}
}

func (me *canvas) DrawColor(color Color) {
	a := float32((color>>24)&0xff) / 255
	r := float32((color>>16)&0xff) / 255
	g := float32((color>>8)&0xff) / 255
	b := float32((color)&0xff) / 255
	C.cairo_set_source_rgba(me.ptr, C.double(r), C.double(g), C.double(b), C.double(a))
	// t1 := time.Now()
	C.cairo_paint(me.ptr)
	// log.V("nuxui", "cairo_paint used time %d", time.Now().Sub(t1).Milliseconds())
}

func (me *canvas) DrawAlpha(alpha float32) {
	C.cairo_paint_with_alpha(me.ptr, C.double(alpha))
}

func (me *canvas) SetColor(color Color) {
	a := float32((color>>24)&0xff) / 255
	r := float32((color>>16)&0xff) / 255
	g := float32((color>>8)&0xff) / 255
	b := float32((color)&0xff) / 255
	C.cairo_set_source_rgba(me.ptr, C.double(r), C.double(g), C.double(b), C.double(a))
}

/*
typedef struct {
    double x_bearing;
    double y_bearing;
    double width;
    double height;
    double x_advance;
    double y_advance;
} cairo_text_extents_t;
*/
func (me *canvas) GetTextRect(text string, fontFamily string, fontSize float32) C.cairo_text_extents_t {
	str := C.CString(text)
	font := C.CString(fontFamily)
	C.cairo_select_font_face(me.ptr, font, C.CAIRO_FONT_SLANT_NORMAL, C.CAIRO_FONT_WEIGHT_NORMAL)

	C.cairo_set_font_size(me.ptr, C.double(fontSize))

	var extents C.cairo_text_extents_t
	C.cairo_text_extents(me.ptr, str, (*C.cairo_text_extents_t)(unsafe.Pointer(&extents)))
	C.free(unsafe.Pointer(str))
	C.free(unsafe.Pointer(font))
	return extents
}

func (me *canvas) DrawText(text string, font *Font, width, height int32, paint *Paint) {
	fontFamily := C.CString(font.Family)
	ctext := C.CString(text)
	me.SetColor(paint.Color)
	C.drawText(me.ptr, fontFamily, C.int(font.Weight), C.int(font.Size), ctext, C.int(width), C.int(height))
	me.drawPaint(paint)
	C.free(unsafe.Pointer(fontFamily))
	C.free(unsafe.Pointer(ctext))
}

func (me *canvas) MeasureText(text string, font *Font, width, height int32) (outWidth, outHeight int32) {
	fontFamily := C.CString(font.Family)
	ctext := C.CString(text)
	var w, h C.int
	C.measureText(me.ptr, fontFamily, C.int(font.Weight), C.int(font.Size), ctext, C.int(width), C.int(height), &w, &h)
	outWidth = int32(float64(w)/float64(C.PANGO_SCALE) + 0.99999)
	outHeight = int32(float64(h)/float64(C.PANGO_SCALE) + 0.99999)
	C.free(unsafe.Pointer(fontFamily))
	C.free(unsafe.Pointer(ctext))
	return
}

func (me *canvas) DrawImage(img Image) {
	C.cairo_set_source_rgba(me.ptr, C.double(0), C.double(1.0), C.double(0), C.double(1))
	C.cairo_set_source_surface(me.ptr, img.Buffer(), 0, 0)
	C.cairo_paint(me.ptr)
}

func (me *canvas) deviceToUser(x, y float32) {
	dx := C.double(x)
	dy := C.double(y)
	C.cairo_device_to_user(me.ptr, &dx, &dy)
}

func (me *canvas) UserToDevice(x, y float32) {
	me.userToDevice(x, y)
}

func (me *canvas) userToDevice(x, y float32) {
	dx := C.double(x)
	dy := C.double(y)
	C.cairo_user_to_device(me.ptr, &dx, &dy)
	log.V("nuxui", "userToDevice x=%f, y=%f, dx=%f, dy=%f", x, y, dx, dy)
	C.cairo_device_to_user(me.ptr, &dx, &dy)
	log.V("nuxui", "deviceToUser x=%f, y=%f, dx=%f, dy=%f", x, y, dx, dy)
}

func (me *canvas) SetAntialias(a int) {
	C.cairo_set_antialias(me.ptr, C.cairo_antialias_t(a))
}

func (me *canvas) GetAntialias() int {
	return int(C.cairo_get_antialias(me.ptr))
}

/////////////////////////////////////////////////////////////////////////////////////////
/////////////                         Text                            ///////////////////
/////////////////////////////////////////////////////////////////////////////////////////
// measureText canvas use an empty bitmap
var canvas4measure Canvas = NewSurfaceFromData(unsafe.Pointer(&[]C.uchar{}), FORMAT_ARGB32, 0, 0, 0).GetCanvas()

func MeasureText(text string, font *Font, width, height int32) (outWidth, outHeight int32) {
	// md := log.Time()
	// defer log.TimeEnd("nuxui", "canvas MeasureText", md)
	return canvas4measure.MeasureText(text, font, width, height)
}

/////////////////////////////////////////////////////////////////////////////////////////
/////////////                         Paint                           ///////////////////
/////////////////////////////////////////////////////////////////////////////////////////

const (
	STROKE        = 0
	FILL          = 1
	FILLANDSTROKE = 2
)

type Paint struct {
	Color Color
	Style int32
	Width float32
}
