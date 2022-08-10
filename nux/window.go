// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/util"
)

var _ Window = (*window)(nil)

type Window interface {
	Component

	Decor() Widget
	Title() string
	SetTitle(string)
	// Size()(width, height int32)
	// SetSize(width, height int32)
	ContentSize() (width, height int32)
	// Alpha() float32
	// SetAlpha(alpha float32)
	// SetDelegate(delegate WindowDelegate)
	// Delegate() WindowDelegate
	Center()
	Show()

	// private methods
	// created()
	createDecor(Attr) Parent
	mountWidget()
	ejectWidget()
	resize()
	draw()
	handlePointerEvent(event PointerEvent) bool
	handleScrollEvent(event ScrollEvent) bool
	handleKeyEvent(event KeyEvent) bool
	handleTypingEvent(event TypingEvent) bool
	requestFocus(widget Widget)
	native() *nativeWindow
}

type window struct {
	*ComponentBase

	nativeWnd   *nativeWindow
	focusWidget Widget
	decor       Widget
	initEvent   PointerEvent
	timer       Timer
}

func newWindow(attr Attr) Window {
	if theApp.mainWindow != nil {
		log.Fatal("nuxui", "nuxui only supported single window")
	}

	if attr == nil {
		attr = Attr{}
	}

	me := &window{}
	theApp.mainWindow = me
	me.ComponentBase = NewComponentBase(me, attr)
	me.decor = me.createDecor(attr)
	OnHoverEnter(me.decor, func(detail GestureDetail) {})
	if content := attr.GetAttr("content", nil); content != nil {
		InflateLayoutAttr(me.decor, content, nil)
	}

	me.nativeWnd = newNativeWindow(attr)
	return me
}

func (me *window) native() *nativeWindow {
	return me.nativeWnd
}

func (me *window) createDecor(attr Attr) Parent {
	decorAttr := Attr{
		"background": "#ffffff",
	}

	creator := FindTypeCreator("nuxui.org/nuxui/ui.Layer")
	w := creator(MergeAttrs(decorAttr, attr))
	if p, ok := w.(Parent); ok {
		p.Info().Self = p
		return p
	} else {
		log.Fatal("nuxui", "decor must be a Parent")
	}

	return nil
}

func (me *window) Decor() Widget {
	return me.decor
}

func (me *window) mountWidget() {
	mountWidget(me.decor, nil)
}

func (me *window) ejectWidget() {
	EjectChild(me.decor)
}

func (me *window) Center() {
	me.nativeWnd.Center()
}

func (me *window) Show() {
	me.nativeWnd.Show()
}

func (me *window) Title() string {
	return me.nativeWnd.Title()
}

func (me *window) SetTitle(title string) {
	me.nativeWnd.SetTitle(title)
}

func (me *window) ContentSize() (width, height int32) {
	return me.nativeWnd.ContentSize()
}

func (me *window) handleKeyEvent(e KeyEvent) bool {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(KeyEventHandler); ok {
			if f.OnKeyEvent(e) {
				return false
			} else {
				goto other
			}
		}
	} else {
		goto other
	}

other:
	if me.decor != nil {
		me.handleOtherWidgetKeyEvent(me.decor.(Parent), e)
	}
	return false
}

func (me *window) handleOtherWidgetKeyEvent(p Parent, e KeyEvent) bool {
	if p.ChildrenCount() > 0 {
		var compt Widget
		for _, c := range p.Children() {
			compt = nil
			if cpt, ok := c.(Component); ok {
				c = cpt.Content() // TODO:: conent still si compent
				compt = cpt
			}
			if cp, ok := c.(Parent); ok {
				if me.handleOtherWidgetKeyEvent(cp, e) {
					return true
				}
			} else if f, ok := c.(KeyEventHandler); ok {
				if f.OnKeyEvent(e) {
					return true
				}
			}

			if compt != nil {
				if f, ok := compt.(KeyEventHandler); ok {
					if f.OnKeyEvent(e) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (me *window) handlePointerEvent(e PointerEvent) bool {
	me.switchFocusIfPossible(e)

	// if me.delegate != nil {
	// 	if f, ok := me.delegate.(windowDelegate_HandlePointerEvent); ok {
	// 		f.HandlePointerEvent(e)
	// 		return
	// 	}
	// }
	hitTestResultManagerInstance.handlePointerEvent(me.Decor(), e)
	return false
}

func (me *window) handleTypingEvent(e TypingEvent) bool {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(TypingEventHandler); ok {
			f.OnTypingEvent(e)
			return false
		}
	}

	log.E("nuxui", "none widget handle typing event")
	return false
}

func (me *window) handleScrollEvent(e ScrollEvent) bool {
	hitTestResultManagerInstance.handleScrollEvent(me.Decor(), e)
	return false
}

func (me *window) requestFocus(widget Widget) {
	if me.focusWidget == widget {
		return
	}

	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(Focus); ok && f.HasFocus() {
			f.FocusChanged(false)
		}
	}

	if f, ok := widget.(Focus); ok {
		if f.Focusable() {
			me.focusWidget = widget
			if !f.HasFocus() {
				f.FocusChanged(true)
			}
		}
	}
}

func (me *window) switchFocusIfPossible(event PointerEvent) {
	if event.Type() != Type_PointerEvent || !event.IsPrimary() {
		return
	}

	switch event.Action() {
	case Action_Down:
		me.initEvent = event
		if event.Kind() == Kind_Mouse {
			me.switchFocusAtPoint(event.X(), event.Y())
		}
	case Action_Drag:
		if event.Kind() == Kind_Touch {
			if me.timer != nil {
				if me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
					me.switchFocusAtPoint(me.initEvent.X(), me.initEvent.Y())
				}
			}
		}
	case Action_Up:
		if event.Kind() == Kind_Touch {
			if event.Time().UnixNano()-me.initEvent.Time().UnixNano() < GESTURE_LONG_PRESS_TIMEOUT.Nanoseconds() {
				me.switchFocusAtPoint(event.X(), event.Y())
			}
		}
	}

}

func (me *window) switchFocusAtPoint(x, y float32) {
	if me.focusWidget != nil {
		if s, ok := me.focusWidget.(Size); ok {
			f := s.Frame()
			if x >= float32(f.X) && x <= float32(f.X+f.Width) &&
				y >= float32(f.Y) && y <= float32(f.Y+f.Height) {
				// point is in current focus widget, do not need change focus
				return
			}
		}

		if f, ok := me.focusWidget.(Focus); ok && f.HasFocus() {
			me.focusWidget = nil
			f.FocusChanged(false)
		}
	}
}

func (me *window) resize() {
	if me.decor == nil {
		return
	}

	width, height := me.ContentSize()

	if s, ok := me.decor.(Size); ok {
		f := s.Frame()
		f.Width = width
		f.Height = height
	}

	if f, ok := me.decor.(Measure); ok {
		f.Measure(MeasureSpec(width, Pixel), MeasureSpec(height, Pixel))
	}

	if f, ok := me.decor.(Layout); ok {
		f.Layout(0, 0, width, height)
	}
}

func (me *window) lockCanvas() Canvas {
	return me.nativeWnd.lockCanvas()
}

func (me *window) unlockCanvas(canvas Canvas) {
	me.nativeWnd.unlockCanvas(canvas)
}

var TestDraw func(Canvas)

func (me *window) draw() {
	canvas := me.lockCanvas()
	me.nativeWnd.draw(canvas, me.decor)
	me.unlockCanvas(canvas)
}

func measureWindowSize(w, h Dimen) (width, height int32) {
	sw, sh := ScreenSize() // TODO taskbar size
	switch w.Mode() {
	case Auto:
		width = 400
	case Pixel:
		width = util.Roundi32(w.Value())
	case Percent:
		width = util.Roundi32(float32(sw) * w.Value() / 100.0)
	case Weight:
		width = sw
	}

	switch h.Mode() {
	case Auto:
		height = 300
	case Pixel:
		height = util.Roundi32(h.Value())
	case Percent:
		height = util.Roundi32(float32(sh) * h.Value() / 100.0)
	case Weight:
		height = sh
	}

	if w.Mode() == Ratio {
		if h.Mode() != Ratio {
			width = util.Roundi32(float32(height) * w.Value())
		} else {
			log.Fatal("nuxui", "width and height size mode can not both Ratio")
		}
	}

	if h.Mode() == Ratio {
		height = util.Roundi32(float32(width) / h.Value())
	}

	return
}
