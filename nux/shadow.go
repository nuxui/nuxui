// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "math"

type OutsideShadow interface {
	SetOutsideShadow(offsetX, offsetY, radius int32, color Color)
	OutsideShadow() (offsetX, offsetY, radius int32, color Color)
}

type WidgetOutsideShadow struct {
	offsetX     int32
	offsetY     int32
	radius      int32
	color       Color
	dataChanged func()
}

func ceil2i32(v float32) int32 {
	return int32(math.Ceil(float64(v)))
}

func NewWidgetOutsideShadow(dataChanged func(), attr Attr) *WidgetOutsideShadow {
	me := &WidgetOutsideShadow{
		offsetX:     ceil2i32(attr.GetDimen("offsetX", "0").Value()),
		offsetY:     ceil2i32(attr.GetDimen("offsetY", "0").Value()),
		radius:      ceil2i32(attr.GetDimen("radius", "0").Value()),
		color:       attr.GetColor("color", Transparent),
		dataChanged: dataChanged,
	}
	return me
}

func (me *WidgetOutsideShadow) SetOutsideShadow(offsetX, offsetY, radius int32, color Color) {
	if me.offsetX != offsetX ||
		me.offsetY != offsetY ||
		me.radius != radius ||
		me.color != color {
		me.offsetX = offsetX
		me.offsetY = offsetY
		me.radius = radius
		me.color = color
		me.dataChanged()
	}
}

func (me *WidgetOutsideShadow) OutsideShadow() (offsetX, offsetY, radius int32, color Color) {
	return me.offsetX, me.offsetY, me.radius, me.color
}
