// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

type Watcher func()

type OnSizeChanged func(Widget)

// TODO size changed
type Size interface {
	Width() Dimen
	SetWidth(width Dimen)
	Height() Dimen
	SetHeight(height Dimen)
	HasMargin() bool
	MarginLeft() Dimen
	MarginTop() Dimen
	MarginRight() Dimen
	MarginBottom() Dimen
	SetMargin(left, top, right, bottom Dimen)
	SetMarginLeft(left Dimen)
	SetMarginTop(top Dimen)
	SetMarginRight(right Dimen)
	SetMarginBottom(bottom Dimen)
	HasPadding() bool
	PaddingLeft() Dimen
	PaddingTop() Dimen
	PaddingRight() Dimen
	PaddingBottom() Dimen
	SetPadding(left, top, right, bottom Dimen)
	SetPaddingLeft(left Dimen)
	SetPaddingTop(top Dimen)
	SetPaddingRight(right Dimen)
	SetPaddingBottom(bottom Dimen)
	MeasuredSize() *MeasuredSize // not nil
	AddOnSizeChanged(callback OnSizeChanged)
	RemoveOnSizeChanged(callback OnSizeChanged)
}

// type Offset struct {
// 	X int32
// 	Y int32
// }

type MeasuredSize struct {
	Width    int32
	Height   int32
	Padding  Rect
	Margin   Rect
	Position RectXY
}

func (me *MeasuredSize) String() string {
	return fmt.Sprintf("{width: %d, height: %d, padding:{left: %d, top: %d, right: %d, bottom: %d}, margin:{left: %d, top: %d, right: %d, bottom: %d}", me.Width,
		me.Height, me.Padding.Left, me.Padding.Top, me.Padding.Right, me.Padding.Bottom, me.Margin.Left, me.Margin.Top, me.Margin.Right, me.Margin.Bottom)
}

// TODO Margin padding width height gravity{vertical: @parent.left} align{vertical: @parent.left} position: {left: @parent}
// width: @me.height height: @me. $parent.right $boxRoot.left
// padding: !auto 10px 10dp 5% !wt !ratio !unlimit
type Padding struct {
	Left   Dimen
	Top    Dimen
	Right  Dimen
	Bottom Dimen
}

func NewPadding() *Padding {
	return &Padding{}
}

func (me *Padding) Creating(attr Attr) {
	if attr == nil {
		attr = Attr{}
	}

	me.Left = getPadding(attr, "left", "0")
	me.Top = getPadding(attr, "top", "0")
	me.Right = getPadding(attr, "right", "0")
	me.Bottom = getPadding(attr, "bottom", "0")
}

func getPadding(attr Attr, key string, defaultValue string) Dimen {
	d := attr.GetDimen(key, defaultValue)
	switch d.Mode() {
	case Auto, Weight, Unspec, Ratio:
		log.Fatal("nuxui", "Unsupported padding dimension mode %s at %s: %s, only supported Pixel, dp, Percent.", d.Mode(), key, d)
	default:
		return d
	}
	return 0
}

// margin: !auto 10px 10dp 1wt 5% !ratio !unlimit
type Margin struct {
	Left   Dimen
	Top    Dimen
	Right  Dimen
	Bottom Dimen
}

func NewMargin() *Margin {
	return &Margin{}
}

func (me *Margin) Creating(attr Attr) {
	if attr == nil {
		attr = Attr{}
	}

	me.Left = getMargin(attr, "left", "0")
	me.Top = getMargin(attr, "top", "0")
	me.Right = getMargin(attr, "right", "0")
	me.Bottom = getMargin(attr, "bottom", "0")
}

func getMargin(attr Attr, key string, defaultValue string) Dimen {
	d := attr.GetDimen(key, defaultValue)
	switch d.Mode() {
	case Auto, Unspec, Ratio:
		log.Fatal("nuxui", "Unsupported padding dimension mode %s at %s: %s, only supported Pixel, dp, Percent.", d.Mode(), key, d)
	default:
		return d
	}
	return 0
}

type WidgetSize struct {
	Owner                  Widget
	width                  Dimen
	height                 Dimen
	padding                *Padding
	margin                 *Margin
	position               interface{} // layout params
	measured               MeasuredSize
	onSizeChangedCallbacks []unsafe.Pointer
}

func (me *WidgetSize) Creating(attr Attr) {
	me.width = attr.GetDimen("width", "auto")
	me.height = attr.GetDimen("height", "auto")

	if padding := attr.Get("padding", nil); padding != nil {
		me.padding = NewPadding()
		me.padding.Creating(padding.(Attr))
	}

	if margin := attr.Get("margin", nil); margin != nil {
		me.margin = NewMargin()
		me.margin.Creating(margin.(Attr))
	}
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

func (me *WidgetSize) HasMargin() bool {
	return me.margin != nil
}

func (me *WidgetSize) MarginLeft() Dimen {
	if me.margin == nil {
		return 0
	}

	return me.margin.Left
}

func (me *WidgetSize) MarginTop() Dimen {
	if me.margin == nil {
		return 0
	}

	return me.margin.Top
}

func (me *WidgetSize) MarginRight() Dimen {
	if me.margin == nil {
		return 0
	}

	return me.margin.Right
}

func (me *WidgetSize) MarginBottom() Dimen {
	if me.margin == nil {
		return 0
	}

	return me.margin.Bottom
}

func (me *WidgetSize) SetMargin(left, top, right, bottom Dimen) {
	if me.margin == nil {
		me.margin = &Margin{
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		}
		me.doSizeChanged()
	} else {
		if me.margin.Left != left || me.margin.Top != top || me.margin.Right != right || me.margin.Bottom != bottom {
			me.margin.Left = left
			me.margin.Top = top
			me.margin.Right = right
			me.margin.Bottom = bottom
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetMarginLeft(left Dimen) {
	if me.margin == nil {
		me.margin = &Margin{
			Left:   left,
			Top:    0,
			Right:  0,
			Bottom: 0,
		}
		me.doSizeChanged()
	} else {
		if me.margin.Left != left {
			me.margin.Left = left
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetMarginTop(top Dimen) {
	if me.margin == nil {
		me.margin = &Margin{
			Left:   0,
			Top:    top,
			Right:  0,
			Bottom: 0,
		}
		me.doSizeChanged()
	} else {
		if me.margin.Top != top {
			me.margin.Top = top
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetMarginRight(right Dimen) {
	if me.margin == nil {
		me.margin = &Margin{
			Left:   0,
			Top:    0,
			Right:  right,
			Bottom: 0,
		}
		me.doSizeChanged()
	} else {
		if me.margin.Right != right {
			me.margin.Right = right
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetMarginBottom(bottom Dimen) {
	if me.margin == nil {
		me.margin = &Margin{
			Left:   0,
			Top:    0,
			Right:  0,
			Bottom: bottom,
		}
		me.doSizeChanged()
	} else {
		if me.margin.Bottom != bottom {
			me.margin.Bottom = bottom
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) HasPadding() bool {
	return me.padding != nil
}

func (me *WidgetSize) PaddingLeft() Dimen {
	if me.padding == nil {
		return 0
	}

	return me.padding.Left
}

func (me *WidgetSize) PaddingTop() Dimen {
	if me.padding == nil {
		return 0
	}

	return me.padding.Top
}

func (me *WidgetSize) PaddingRight() Dimen {
	if me.padding == nil {
		return 0
	}

	return me.padding.Right
}

func (me *WidgetSize) PaddingBottom() Dimen {
	if me.padding == nil {
		return 0
	}

	return me.padding.Bottom
}

func (me *WidgetSize) SetPadding(left, top, right, bottom Dimen) {
	if me.padding == nil {
		me.padding = &Padding{
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		}
		me.doSizeChanged()
	} else {
		if me.padding.Left != left || me.padding.Top != top || me.padding.Right != right || me.padding.Bottom != bottom {
			me.padding.Left = left
			me.padding.Top = top
			me.padding.Right = right
			me.padding.Bottom = bottom
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetPaddingLeft(left Dimen) {
	if me.padding == nil {
		me.padding = &Padding{
			Left:   left,
			Top:    0,
			Right:  0,
			Bottom: 0,
		}
		me.doSizeChanged()
	} else {
		if me.padding.Left != left {
			me.padding.Left = left
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetPaddingTop(top Dimen) {
	if me.padding == nil {
		me.padding = &Padding{
			Left:   0,
			Top:    top,
			Right:  0,
			Bottom: 0,
		}
		me.doSizeChanged()
	} else {
		if me.padding.Top != top {
			me.padding.Top = top
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetPaddingRight(right Dimen) {
	if me.padding == nil {
		me.padding = &Padding{
			Left:   0,
			Top:    0,
			Right:  right,
			Bottom: 0,
		}
		me.doSizeChanged()
	} else {
		if me.padding.Right != right {
			me.padding.Right = right
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) SetPaddingBottom(bottom Dimen) {
	if me.padding == nil {
		me.padding = &Padding{
			Left:   0,
			Top:    0,
			Right:  0,
			Bottom: bottom,
		}
		me.doSizeChanged()
	} else {
		if me.padding.Bottom != bottom {
			me.padding.Bottom = bottom
			me.doSizeChanged()
		}
	}
}

func (me *WidgetSize) MeasuredSize() *MeasuredSize {
	return &me.measured
}

func (me *WidgetSize) AddOnSizeChanged(callback OnSizeChanged) {
	if callback == nil {
		return
	}

	if me.onSizeChangedCallbacks == nil {
		me.onSizeChangedCallbacks = []unsafe.Pointer{}
	}

	p := unsafe.Pointer(&callback)
	for _, o := range me.onSizeChangedCallbacks {
		if o == p {
			log.Fatal("nuxui", "The OnSizeChanged callback is existed.")
		}
	}

	me.onSizeChangedCallbacks = append(me.onSizeChangedCallbacks, unsafe.Pointer(&callback))
}

func (me *WidgetSize) RemoveOnSizeChanged(callback OnSizeChanged) {
	if me.onSizeChangedCallbacks != nil && callback != nil {
		p := unsafe.Pointer(&callback)
		for i, o := range me.onSizeChangedCallbacks {
			if o == p {
				me.onSizeChangedCallbacks = append(me.onSizeChangedCallbacks[:i], me.onSizeChangedCallbacks[i+1:]...)
			}
		}
	}
}

func (me *WidgetSize) doSizeChanged() {
	if me.Owner == nil {
		log.Fatal("nuxui", "set target to WidgetSize first.")
	}

	for _, c := range me.onSizeChangedCallbacks {
		(*(*OnSizeChanged)(c))(me.Owner)
	}
}
