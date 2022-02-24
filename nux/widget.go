// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

type WidgetCreator func(...Attr) Widget

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
	Self    Widget
	ID      string
	Parent  Parent
	Mounted bool
	Mixins  []any
}

type WidgetBase struct {
	info *WidgetInfo
}

func NewWidgetBase(attrs ...Attr) *WidgetBase {
	attr := MergeAttrs(attrs...)
	return &WidgetBase{
		info: &WidgetInfo{
			ID:      attr.GetString("id", ""),
			Mounted: false,
		},
	}
}

func (me *WidgetBase) Info() *WidgetInfo {
	return me.info
}

func AddMixins(widget Widget, mixin any) {
	for _, m := range widget.Info().Mixins {
		if m == mixin {
			log.E("nuxui", "the mixin %T is already existed.", mixin)
			return
		}
	}
	widget.Info().Mixins = append(widget.Info().Mixins, mixin)
}

type viewfuncs interface {
	Measure
	Layout
	Draw
}

func isView(widget Widget) bool {
	_, ret := widget.(viewfuncs)

	return ret
}
