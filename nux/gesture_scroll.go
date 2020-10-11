// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

var scrollWidget map[Widget]GestureCallback = map[Widget]GestureCallback{}

func OnScrollX(widget Widget, callback GestureCallback) {
	scrollWidget[widget] = callback
}

func OnScrollY(widget Widget, callback GestureCallback) {
	scrollWidget[widget] = callback
}

func handleScrollEvent(event Event) {
	for k, v := range scrollWidget {
		v(eventToDetail(event, k))
	}
}
