// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"strings"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/util"
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

	ScrollX() float32
	ScrollY() float32
	SetScroll(x, y float32)
	ScrollTo(x, y float32)

	AddSizeObserver(observer func())
	RemoveSizeObserver(observer func())
}

type Rect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

func (me *Rect) Expand(rect *Rect) {
	if me.X > rect.X {
		me.X = rect.X
	}
	if me.Y > rect.Y {
		me.Y = rect.Y
	}
	if me.Width < rect.Width {
		me.Width = rect.Width
	}
	if me.Height < rect.Height {
		me.Height = rect.Height
	}
}

type Side struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type Frame struct {
	X       int32
	Y       int32
	Width   int32
	Height  int32
	Padding Side
	Margin  Side
}

func (me *Frame) String() string {
	return fmt.Sprintf("{x:%d, y:%d, width: %d, height: %d, padding:{left: %d, top: %d, right: %d, bottom: %d}, margin:{left: %d, top: %d, right: %d, bottom: %d}",
		me.X, me.Y, me.Width, me.Height, me.Padding.Left, me.Padding.Top, me.Padding.Right, me.Padding.Bottom, me.Margin.Left, me.Margin.Top, me.Margin.Right, me.Margin.Bottom)
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

/*
padding: {left: 10px, top: 10px, right: 10px, bottom: 10px}
padding: 10px, 10px, 10px, 10px
padding: 10px
*/
func NewPadding(attr any) *Padding {
	if attr == nil {
		return &Padding{}
	}
	switch t := attr.(type) {
	case Attr:
		return &Padding{
			Left:   getPadding(t, "left", "0"),
			Top:    getPadding(t, "top", "0"),
			Right:  getPadding(t, "right", "0"),
			Bottom: getPadding(t, "bottom", "0"),
		}
	case string:
		dimens := strings.Split(t, ",")
		if len(dimens) == 1 {
			d := SDimen(strings.TrimSpace(dimens[0]))
			return &Padding{d, d, d, d}
		} else {
			p := &Padding{}
			for i, d := range dimens {
				d = strings.TrimSpace(d)
				switch i {
				case 0:
					p.Left = SDimen(d)
				case 1:
					p.Top = SDimen(d)
				case 2:
					p.Right = SDimen(d)
				case 3:
					p.Bottom = SDimen(d)
				}
			}
			return p
		}
	}
	return &Padding{}
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

/*
margin: !auto 10px 10dp 1em 1wt 5% !ratio !unlimit
margin: {left: 10px, top: 10px, right: 10px, bottom: 10px}
margin: 10px, 10px, 10px, 10px
margin: 10px
*/
type Margin struct {
	Left   Dimen
	Top    Dimen
	Right  Dimen
	Bottom Dimen
}

func NewMargin(attr any) *Margin {
	if attr == nil {
		return &Margin{}
	}
	switch t := attr.(type) {
	case Attr:
		return &Margin{
			Left:   getMargin(t, "left", "0"),
			Top:    getMargin(t, "top", "0"),
			Right:  getMargin(t, "right", "0"),
			Bottom: getMargin(t, "bottom", "0"),
		}
	case string:
		dimens := strings.Split(t, ",")
		if len(dimens) == 1 {
			d := SDimen(strings.TrimSpace(dimens[0]))
			return &Margin{d, d, d, d}
		} else {
			m := &Margin{}
			for i, d := range dimens {
				d = strings.TrimSpace(d)
				switch i {
				case 0:
					m.Left = SDimen(d)
				case 1:
					m.Top = SDimen(d)
				case 2:
					m.Right = SDimen(d)
				case 3:
					m.Bottom = SDimen(d)
				}
			}
			return m
		}
	}
	return &Margin{}
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
	scrollX       float32
	scrollY       float32
	sizeObservers []func()
}

func NewWidgetSize(attr Attr) *WidgetSize {
	if attr == nil {
		attr = Attr{}
	}

	me := &WidgetSize{
		sizeObservers: []func(){},
	}

	me.width = attr.GetDimen("width", "auto")
	me.height = attr.GetDimen("height", "auto")

	if padding := attr.Get("padding", nil); padding != nil {
		me.padding = NewPadding(padding)
	}

	if margin := attr.Get("margin", nil); margin != nil {
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

// TODO:: why not use value instead ref
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

func (me *WidgetSize) ScrollX() float32 {
	return me.scrollX
}

func (me *WidgetSize) ScrollY() float32 {
	return me.scrollY
}

func (me *WidgetSize) SetScroll(x, y float32) {
	me.scrollX = x
	me.scrollY = y
}

func (me *WidgetSize) ScrollTo(x, y float32) {
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
