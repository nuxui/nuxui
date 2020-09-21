// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"unsafe"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

type OnVisualChanged func(nux.Widget)

type Visual interface {
	Border()
	SetBorder()
	Background() Drawable
	SetBackground(Drawable)
	SetBackgroundColor(nux.Color)
	Foreground() Drawable
	SetForeground(Drawable)
	SetForegroundColor(nux.Color)
	Visible() Visible
	SetVisible(visible Visible)
	AddOnVisualChanged(callback OnVisualChanged)
	RemoveOnVisualChanged(callback OnVisualChanged)
	Translucent() bool // TODO:: Can the event penetrate ?
	SetTranslucent(bool)
}

type WidgetVisual struct {
	Owner                    nux.Widget
	background               Drawable // TODO replace type to *Drawable
	foreground               Drawable
	visible                  Visible
	onVisualChangedCallbacks []unsafe.Pointer
	translucent              bool
}

func (me *WidgetVisual) Creating(attr nux.Attr) {
	bg := NewColorDrawable()
	bg.SetColor(attr.GetColor("background", nux.Transparent))
	me.background = bg

	foreground := NewColorDrawable()
	foreground.SetColor(attr.GetColor("foreground", nux.Transparent))
	me.foreground = foreground

	me.translucent = attr.GetBool("translucent", false)

	visible := attr.GetString("visible", "show")
	switch visible {
	case "show":
		me.visible = Show
	case "hide":
		me.visible = Hide
	case "gone":
		me.visible = Gone
	default:
		log.Fatal("nuxui", "visible should be 'show', 'hide' or 'gone'")
	}
}

func (me *WidgetVisual) Border() {}

func (me *WidgetVisual) SetBorder() {}

func (me *WidgetVisual) Background() Drawable {
	return me.background
}

func (me *WidgetVisual) SetBackground(background Drawable) {
	if me.background != background && (me.background == nil || !me.background.Equal(background)) {
		me.background = background
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) SetBackgroundColor(background nux.Color) {
	b := NewColorDrawable()
	b.SetColor(background)

	if me.background != b && (me.background == nil || !me.background.Equal(b)) {
		me.background = b
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) Foreground() Drawable {
	return me.foreground
}

func (me *WidgetVisual) SetForeground(foreground Drawable) {
	if me.foreground != foreground && (me.foreground == nil || !me.foreground.Equal(foreground)) {
		me.foreground = foreground
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) SetForegroundColor(foreground nux.Color) {
	f := NewColorDrawable()
	f.SetColor(foreground)

	if me.foreground != f && (me.foreground == nil || !me.foreground.Equal(f)) {
		me.foreground = f
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) Visible() Visible {
	return me.visible
}
func (me *WidgetVisual) SetVisible(visible Visible) {
	if visible != Show && visible != Hide && visible != Gone {
		log.Fatal("nuxui", "visible should be 'Show', 'Hide' or 'Gone'")
	}

	if me.visible != visible {
		me.visible = visible
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) AddOnVisualChanged(callback OnVisualChanged) {
	if callback == nil {
		return
	}

	if me.onVisualChangedCallbacks == nil {
		me.onVisualChangedCallbacks = []unsafe.Pointer{}
	}

	p := unsafe.Pointer(&callback)
	for _, o := range me.onVisualChangedCallbacks {
		if o == p {
			log.Fatal("nuxui", "The OnVisualChanged callback is existed.")
		}
	}

	me.onVisualChangedCallbacks = append(me.onVisualChangedCallbacks, unsafe.Pointer(&callback))
}

func (me *WidgetVisual) RemoveOnVisualChanged(callback OnVisualChanged) {
	if me.onVisualChangedCallbacks != nil && callback != nil {
		p := unsafe.Pointer(&callback)
		for i, o := range me.onVisualChangedCallbacks {
			if o == p {
				me.onVisualChangedCallbacks = append(me.onVisualChangedCallbacks[:i], me.onVisualChangedCallbacks[i+1:]...)
			}
		}
	}
}

func (me *WidgetVisual) doVisualChanged() {
	if me.Owner == nil {
		log.Fatal("nuxui", "set target for WidgetVisual first.")
	}

	for _, c := range me.onVisualChangedCallbacks {
		(*(*OnVisualChanged)(c))(me.Owner)
	}
}

func (me *WidgetVisual) Translucent() bool {
	return me.translucent
}

func (me *WidgetVisual) SetTranslucent(translucent bool) {
	me.translucent = translucent
}
