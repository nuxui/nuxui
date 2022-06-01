// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

type OnVisualChanged func(nux.Widget)

type Visual interface {
	Background() nux.Drawable
	SetBackground(nux.Drawable)
	SetBackgroundColor(nux.Color)
	Foreground() nux.Drawable
	SetForeground(nux.Drawable)
	SetForegroundColor(nux.Color)
	Visible() Visible
	SetVisible(visible Visible)
	Disable() bool
	SetDisable(bool)
	Translucent() bool // TODO:: Can the event penetrate ?
	SetTranslucent(bool)
}

type WidgetVisual struct {
	owner       nux.Widget
	background  nux.Drawable
	foreground  nux.Drawable
	visible     Visible
	disable     bool
	translucent bool
}

func NewWidgetVisual(owner nux.Widget, attr nux.Attr) *WidgetVisual {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &WidgetVisual{
		owner:       owner,
		visible:     visibleFromString(attr.GetString("visible", "show")),
		translucent: attr.GetBool("translucent", false),
		background:  nux.InflateDrawable(attr.Get("background", nil)),
		foreground:  nux.InflateDrawable(attr.Get("foreground", nil)),
		disable:     attr.GetBool("disable", false),
	}

	return me
}

func (me *WidgetVisual) Background() nux.Drawable {
	return me.background
}

func (me *WidgetVisual) SetBackground(background nux.Drawable) {
	if me.background != background && (me.background == nil || !me.background.Equal(background)) {
		me.background = background
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) SetBackgroundColor(background nux.Color) {
	if me.background != nil {
		if c, ok := me.background.(ColorDrawable); ok {
			c.SetColor(background)
			me.doVisualChanged()
			return
		}
	}

	b := NewColorDrawableWithColor(background)
	if me.background == nil || !me.background.Equal(b) {
		me.background = b
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) Foreground() nux.Drawable {
	return me.foreground
}

func (me *WidgetVisual) SetForeground(foreground nux.Drawable) {
	if me.foreground != foreground && (me.foreground == nil || !me.foreground.Equal(foreground)) {
		me.foreground = foreground
		me.doVisualChanged()
	}
}

func (me *WidgetVisual) SetForegroundColor(foreground nux.Color) {
	if me.foreground != nil {
		if c, ok := me.foreground.(ColorDrawable); ok {
			c.SetColor(foreground)
			me.doVisualChanged()
			return
		}
	}

	f := NewColorDrawableWithColor(foreground)
	if me.foreground == nil || !me.foreground.Equal(f) {
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

func (me *WidgetVisual) Translucent() bool {
	return me.translucent
}

func (me *WidgetVisual) SetTranslucent(translucent bool) {
	me.translucent = translucent
}

func (me *WidgetVisual) Disable() bool {
	return me.disable
}

func (me *WidgetVisual) SetDisable(disable bool) {
	me.disable = disable
}

func (me *WidgetVisual) doVisualChanged() {
	nux.RequestRedraw(me.owner)
}

func visibleFromString(visible string) Visible {
	switch visible {
	case "show":
		return Show
	case "hide":
		return Hide
	case "gone":
		return Gone
	default:
		log.Fatal("nuxui", "visible should be 'show', 'hide' or 'gone'")
	}
	return Show
}
