// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Creator func(Context, ...Attr) Widget

type Widget interface {
	// ID() string
	// SetID(string)
	// Parent() Parent
	// AssignParent(parent Parent)
	// IsMounted() bool
	Info() *WidgetInfo
}

// type WidgetInfo interface {
// }

// type State int

type WidgetInfo struct {
	ID      string
	Parent  Parent
	Mounted bool
}

type WidgetBase struct {
	info  *WidgetInfo
	owner Widget
}

func NewWidgetBase(ctx Context, owner Widget, attrs ...Attr) *WidgetBase {
	attr := Attr{}
	if len(attrs) > 0 {
		attr = attrs[0]
	}
	return &WidgetBase{
		owner: owner,
		info: &WidgetInfo{
			ID:      attr.GetString("id", ""),
			Mounted: false,
		},
	}
}

func (me *WidgetBase) Info() *WidgetInfo {
	return me.info
}

// type Template interface {
// 	Template() string
// }

// type Render interface {
// 	Render() Widget
// }

type viewfuncs interface {
	Measure
	Layout
	Draw
}
