// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

type HitTestResult interface {
	Add(widget Widget)
	Remove(widget Widget)
	Contain(widget Widget) bool
	Results() []Widget
}

type hitTestResult struct {
	widgets []Widget
}

func NewHitTestResult() HitTestResult {
	h := &hitTestResult{
		widgets: []Widget{},
	}
	return h
}

func (me *hitTestResult) Add(widget Widget) {
	if true /*TODO debug*/ {
		for _, w := range me.widgets {
			if w == widget {
				log.Fatal("hitTest", "The entry is already added.")
			}
		}
	}

	me.widgets = append(me.widgets, widget)
}

func (me *hitTestResult) Remove(widget Widget) {
	for i, w := range me.widgets {
		if w == widget {
			me.widgets = append(me.widgets[:i], me.widgets[i+1:]...)
		}
	}

	if false /*TODO debug*/ {
		log.Fatal("hitTest", "The entry is not exist.")
	}
}

func (me *hitTestResult) Contain(widget Widget) bool {
	for _, w := range me.widgets {
		if w == widget {
			return true
		}
	}
	return false
}

func (me *hitTestResult) Results() []Widget {
	return me.widgets
}
