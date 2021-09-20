// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

var scrollWidget map[Widget]GestureCallback = map[Widget]GestureCallback{}

func OnScrollX(widget Widget, callback GestureCallback) {
	scrollWidget[widget] = callback
}

func OnScrollY(widget Widget, callback GestureCallback) {
	scrollWidget[widget] = callback
}

func handleScrollEvent(event ScrollEvent) {
	log.I("nuxui", "ScrollX=%f, ScrollY=%f", event.ScrollX(), event.ScrollY())
	for k, v := range scrollWidget {
		v(scrollEventToDetail(event, k))
	}
}
