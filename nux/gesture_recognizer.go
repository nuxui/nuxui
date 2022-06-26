// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/util"
)

type GestureCallback func(detail GestureDetail)

type GestureRecognizer interface {
	GestureArenaMember
	PointerAllowed(event PointerEvent) bool
	HandleAllowedPointer(event PointerEvent)
	Clear(widget Widget) // clear callbacks of widget when ejected
}

// ------------- gestureRecognizerManager -------------------------

func GestureManager() *gestureRecognizerManager {
	return gestureManager
}

var gestureManager *gestureRecognizerManager = &gestureRecognizerManager{
	gestureRecognizers: map[Widget][]GestureRecognizer{},
}

type gestureRecognizerManager struct {
	gestureRecognizers map[Widget][]GestureRecognizer
}

func (me *gestureRecognizerManager) getGestureRecognizers(widget Widget) []GestureRecognizer {
	return me.gestureRecognizers[widget]
}

func (me *gestureRecognizerManager) AddGestureRecognizer(widget Widget, recognizer GestureRecognizer) {
	if recognizer == nil || widget == nil {
		return
	}

	if rs, ok := me.gestureRecognizers[widget]; ok {
		if debug_gesture {
			name := util.TypeName(recognizer)
			for _, r := range rs {
				rname := util.TypeName(r)
				if name == rname {
					log.Fatal("nuxui", "The gesture recognizer '%s' is existed.", name)
					break
				}
			}
		}
		me.gestureRecognizers[widget] = append(rs, recognizer)
	} else {
		me.gestureRecognizers[widget] = []GestureRecognizer{recognizer}
	}
}

func (me *gestureRecognizerManager) ClearGestureRecognizer(widget Widget) {
	if rs, ok := me.gestureRecognizers[widget]; ok {
		for _, r := range rs {
			r.Clear(widget)
		}

		delete(me.gestureRecognizers, widget)
	}
}

func (me *gestureRecognizerManager) FindGestureRecognizer(widget Widget, recognizerType GestureRecognizer) GestureRecognizer {
	if rs, ok := me.gestureRecognizers[widget]; ok {
		name := util.TypeName(recognizerType)
		for _, r := range rs {
			if util.TypeName(r) == name {
				return r
			}
		}
	}

	return nil
}
