package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/util"
)

func OnHover(widget Widget, callback GestureCallback) {
	addHoverCallback(widget, actionHover, callback)
}

func OnHoverEnter(widget Widget, callback GestureCallback) {
	addHoverCallback(widget, actionHoverEnter, callback)
}

func OnHoverExit(widget Widget, callback GestureCallback) {
	addHoverCallback(widget, actionHoverExit, callback)
}

func addHoverCallback(widget Widget, which int, callback GestureCallback) {
	c, ok := hoverGestureManager_.c[widget]
	if !ok {
		c = [][]GestureCallback{{}, {}, {}}
		hoverGestureManager_.c[widget] = c
	}
	if c[which] == nil {
		c[which] = []GestureCallback{callback}
	} else {
		for _, cb := range c[which] {
			if util.SameFunc(cb, callback) {
				log.Fatal("nuxui", "The %s callback is already existed.", []string{"OnHover", "OnHoverEnter", "OnHoverExit"}[which])
			}
		}
		c[which] = append(c[which], callback)
	}
}

const (
	actionHover = iota
	actionHoverEnter
	actionHoverExit
)

func HoverGestureManager() *hoverGestureManager {
	return hoverGestureManager_
}

var hoverGestureManager_ *hoverGestureManager = &hoverGestureManager{c: map[Widget][][]GestureCallback{}}

type hoverGestureManager struct {
	hoverWidget_ Widget

	c map[Widget][][]GestureCallback
}

func (me *hoverGestureManager) existHoverCallback(widget Widget) (exist bool) {
	_, exist = me.c[widget]
	return
}

func (me *hoverGestureManager) invokeHoverEvent(widget Widget, event PointerEvent) {
	if me.hoverWidget_ == widget {
		if cc, ok := me.c[widget]; ok {
			for _, c := range cc[actionHover] {
				c(pointerEventToDetail(event, widget))
			}
		}
	} else {
		if me.hoverWidget_ != nil {
			if cc, ok := me.c[me.hoverWidget_]; ok {
				for _, c := range cc[actionHoverExit] {
					c(pointerEventToDetail(event, me.hoverWidget_))
				}
			}
		}

		if widget != nil {
			if cc, ok := me.c[widget]; ok {
				for _, c := range cc[actionHoverEnter] {
					c(pointerEventToDetail(event, widget))
				}
			}
		}

		me.hoverWidget_ = widget
	}
}
