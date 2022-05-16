// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "nuxui.org/nuxui/log"

type WidgetCreator func(Attr) Widget

type Widget interface {
	Info() *WidgetInfo
}

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

func NewWidgetBase(attr Attr) *WidgetBase {
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

func FindChild(widget Widget, id string) Widget {
	if id == "" {
		log.Fatal("nuxui", "the widget %T id must be specified", widget)
		return nil
	}

	w := findChild(widget, id)
	if w == nil {
		log.Fatal("nuxui", "the id '%s' was not found in widget %T\n", id, widget)
		return nil
	}

	return w
}

func findChild(widget Widget, id string) Widget {
	if id == widget.Info().ID {
		return widget
	}

	if c, ok := widget.(Component); ok {
		if ret := findChild(c.Content(), id); ret != nil {
			return ret
		}
	}

	if p, ok := widget.(Parent); ok {
		for _, child := range p.Children() {
			if _, ok := child.(Component); ok {
				continue
			}

			if ret := findChild(child, id); ret != nil {
				return ret
			}
		}
	}

	return nil
}
