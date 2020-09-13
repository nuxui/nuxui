// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

// nux.AddMixin(widget, mixin)
// nux.DelMixin(widget, mixin)
// nux.ClearMixin(widget)

var mixinStorage map[Widget][]Mixin = map[Widget][]Mixin{}

func AddMixins(widget Widget, mixins ...Mixin) {
	if m, ok := mixinStorage[widget]; ok {
		m = append(m, mixins...)
		mixinStorage[widget] = m
	} else {
		mixinStorage[widget] = mixins
	}

	for _, m := range mixinStorage[widget] {
		m.AssignOwner(widget)
	}
}

func DelMixins(widget Widget, mixins ...Mixin) {
	if mixs, ok := mixinStorage[widget]; ok {
		for _, d := range mixins {
			for i, m := range mixs {
				if m == d {
					mixs = append(mixs[:i], mixs[i+1:]...)
					m.AssignOwner(nil)
				}
			}
		}
		mixinStorage[widget] = mixs
	}
}

func ClearMixins(widget Widget) {
	delete(mixinStorage, widget)
}

func Mixins(widget Widget) []Mixin {
	return mixinStorage[widget]
}
