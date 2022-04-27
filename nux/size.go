// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/util"
)

type Size interface {
	Width() Dimen
	SetWidth(width Dimen)
	Height() Dimen
	SetHeight(height Dimen)
	Margin() *Margin     // can be nil
	SetMargin(*Margin)   // can be nil
	Padding() *Padding   // can be nil
	SetPadding(*Padding) // can be nil

	Frame() *Frame // not nil

	ScrollX() int32
	ScrollY() int32
	SetScroll(x, y int32)
	ScrollTo(x, y int32)

	AddSizeObserver(observer func())
	RemoveSizeObserver(observer func())
}

type Frame struct {
	X       int32
	Y       int32
	Width   int32
	Height  int32
	Padding Rect
	Margin  Rect
}

func (me *Frame) String() string {
	return fmt.Sprintf("{width: %d, height: %d, padding:{left: %d, top: %d, right: %d, bottom: %d}, margin:{left: %d, top: %d, right: %d, bottom: %d}", me.Width,
		me.Height, me.Padding.Left, me.Padding.Top, me.Padding.Right, me.Padding.Bottom, me.Margin.Left, me.Margin.Top, me.Margin.Right, me.Margin.Bottom)
}

func (me *Frame) Clear() {
	me.Width = 0
	me.Height = 0
	me.Margin.Left = 0
	me.Margin.Top = 0
	me.Margin.Right = 0
	me.Margin.Bottom = 0
	me.Padding.Left = 0
	me.Padding.Top = 0
	me.Padding.Right = 0
	me.Padding.Bottom = 0
}

/* padding:  only supported definite size and Percent size
10px 10dp 1em 5% !auto !wt !ratio !unlimit
padding Percent size is associate with parent size, eg: 5% = 0.05*parent.Width
*/
type Padding struct {
	Left   Dimen
	Top    Dimen
	Right  Dimen
	Bottom Dimen
}

func NewPadding(attr Attr) *Padding {
	if attr == nil {
		return &Padding{}
	}
	return &Padding{
		Left:   getPadding(attr, "left", "0"),
		Top:    getPadding(attr, "top", "0"),
		Right:  getPadding(attr, "right", "0"),
		Bottom: getPadding(attr, "bottom", "0"),
	}
}

func getPadding(attr Attr, key string, defaultValue string) Dimen {
	d := attr.GetDimen(key, defaultValue)
	switch d.Mode() {
	case Auto, Weight, Unlimit, Ratio:
		log.Fatal("nuxui", "Unsupported padding dimension mode %s at %s: %s, only supported definite size and Percent size", d.Mode(), key, d)
	default:
		return d
	}
	return 0
}

func (me *Padding) Equal(value *Padding) bool {
	if value != nil {
		return me.Left == value.Left &&
			me.Top == value.Top &&
			me.Right == value.Right &&
			me.Bottom == value.Top
	}
	return false
}

// margin: !auto 10px 10dp 1em 1wt 5% !ratio !unlimit
type Margin struct {
	Left   Dimen
	Top    Dimen
	Right  Dimen
	Bottom Dimen
}

func NewMargin(attr Attr) *Margin {
	if attr == nil {
		return &Margin{}
	}

	return &Margin{
		Left:   getMargin(attr, "left", "0"),
		Top:    getMargin(attr, "top", "0"),
		Right:  getMargin(attr, "right", "0"),
		Bottom: getMargin(attr, "bottom", "0"),
	}
}

func getMargin(attr Attr, key string, defaultValue string) Dimen {
	d := attr.GetDimen(key, defaultValue)
	switch d.Mode() {
	case Auto, Unlimit, Ratio:
		log.Fatal("nuxui", "Unsupported padding dimension mode %s at %s: %s, only supported Pixel, dp, Weight, Percent.", d.Mode(), key, d)
	default:
		return d
	}
	return 0
}

func (me *Margin) Equal(value *Margin) bool {
	if value != nil {
		return me.Left == value.Left &&
			me.Top == value.Top &&
			me.Right == value.Right &&
			me.Bottom == value.Top
	}
	return false
}

type WidgetSize struct {
	width         Dimen
	height        Dimen
	padding       *Padding
	margin        *Margin
	frame         Frame
	scrollX       int32
	scrollY       int32
	sizeObservers []func()
}

func NewWidgetSize(attrs ...Attr) *WidgetSize {
	me := &WidgetSize{
		sizeObservers: []func(){},
	}

	attr := MergeAttrs(attrs...)

	me.width = attr.GetDimen("width", "auto")
	me.height = attr.GetDimen("height", "auto")

	if padding := attr.GetAttr("padding", nil); padding != nil {
		me.padding = NewPadding(padding)
	}

	if margin := attr.GetAttr("margin", nil); margin != nil {
		me.margin = NewMargin(margin)
	}

	return me
}

func (me *WidgetSize) Width() Dimen {
	return me.width
}

func (me *WidgetSize) SetWidth(width Dimen) {
	if me.width != width {
		me.width = width
		me.doSizeChanged()
	}
}

func (me *WidgetSize) Height() Dimen {
	return me.height
}

func (me *WidgetSize) SetHeight(height Dimen) {
	if me.height != height {
		me.height = height
		me.doSizeChanged()
	}
}

func (me *WidgetSize) Padding() *Padding {
	return me.padding
}

func (me *WidgetSize) SetPadding(padding *Padding) {
	if (me.padding == nil && padding != nil) || (me.padding != nil && !me.padding.Equal(padding)) {
		me.padding = padding
		me.doSizeChanged()
	}
}

func (me *WidgetSize) Margin() *Margin {
	return me.margin
}

func (me *WidgetSize) SetMargin(margin *Margin) {
	if (me.margin == nil && margin != nil) || (me.margin != nil && !me.margin.Equal(margin)) {
		me.margin = margin
		me.doSizeChanged()
	}
}

func (me *WidgetSize) ScrollX() int32 {
	return me.scrollX
}

func (me *WidgetSize) ScrollY() int32 {
	return me.scrollY
}

func (me *WidgetSize) SetScroll(x, y int32) {
	me.scrollX = x
	me.scrollY = y
}

func (me *WidgetSize) ScrollTo(x, y int32) {
	me.scrollX += x
	me.scrollY += y
}

func (me *WidgetSize) Frame() *Frame {
	return &me.frame
}

func (me *WidgetSize) AddSizeObserver(observer func()) {
	if observer == nil {
		return
	}

	if debug_size {
		for _, cb := range me.sizeObservers {
			if util.SameFunc(cb, observer) {
				log.Fatal("nuxui", "The OnSizeChanged callback is existed.")
			}
		}
	}

	me.sizeObservers = append(me.sizeObservers, observer)
}

func (me *WidgetSize) RemoveSizeObserver(observer func()) {
	if observer != nil {
		for i, cb := range me.sizeObservers {
			if util.SameFunc(cb, observer) {
				me.sizeObservers = append(me.sizeObservers[:i], me.sizeObservers[i+1:]...)
			}
		}
	}
}

func (me *WidgetSize) doSizeChanged() {
	for _, observer := range me.sizeObservers {
		observer()
	}
}
