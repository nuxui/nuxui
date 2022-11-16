// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package nux

import ()

type nativeWindow struct {
}

func newNativeWindow(attr Attr) *nativeWindow {
	me := &nativeWindow{}
	return me
}

func freeWindow(me *nativeWindow) {
}

func (me *nativeWindow) Center() {
}

func (me *nativeWindow) Show() {
}

func (me *nativeWindow) ContentSize() (width, height int32) {
	return 800, 800
}

func (me *nativeWindow) Title() string {
	return ""
}

func (me *nativeWindow) SetTitle(title string) {

}

func (me *nativeWindow) lockCanvas() Canvas {
	return newCanvas()
}

func (me *nativeWindow) unlockCanvas(canvas Canvas) {
}

func (me *nativeWindow) draw(canvas Canvas, decor Widget) {
	if decor != nil {
		if f, ok := decor.(Draw); ok {
			canvas.Save()
			if TestDraw != nil {
				TestDraw(canvas)
			} else {
				f.Draw(canvas)
			}
			canvas.Restore()
		}
	}
}
