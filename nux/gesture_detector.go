// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/util"

type GestureMixin interface {
	Mixin
	ComptHelper
}

func NewGestureMixin() GestureMixin {
	return &gestureMixin{}
}

type gestureMixin struct {
	WidgetMixin
	ComptHelperPart

	tapGestureRecognizer TapGestureRecognizer
}

func (me *gestureMixin) Creating(attr Attr) {
	if onTap := attr.GetString("onTap", ""); onTap != "" {
		// log.V("nuxui", "xxxxxxxxx onTap = %s", onTap)
		OnTap(me.Owner(), func(detail GestureDetail) {
			// log.V("nuxui", "xxxxxxxxx onTap 2 = %s", onTap)
			util.ReflectCall(me.Component(), onTap, detail.Target())
		})
	}
	if onTapDown := attr.GetString("onTapDown", ""); onTapDown != "" {
		OnTapDown(me.Owner(), func(detail GestureDetail) {
			util.ReflectCall(me.Component(), onTapDown, detail.Target())
		})
	}
	if onTapUp := attr.GetString("onTapUp", ""); onTapUp != "" {
		OnTapUp(me.Owner(), func(detail GestureDetail) {
			util.ReflectCall(me.Component(), onTapUp, detail.Target())
		})
	}
	if onTapCancel := attr.GetString("onTapCancel", ""); onTapCancel != "" {
		OnTapCancel(me.Owner(), func(detail GestureDetail) {
			util.ReflectCall(me.Component(), onTapCancel, detail.Target())
		})
	}
}
