// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/util"
)

type GestureHandler interface {
	AddGestureRecoginer(recognizer GestureRecognizer)
	RemoveGestureRecoginer(recognizer GestureRecognizer)
	HandlePointerEvent(event Event)
	FindGestureRecognizer(recognizerType GestureRecognizer) GestureRecognizer
}

func NewGestureHandler() GestureHandler {
	return &gestureHandler{}
}

type gestureHandler struct {
	gestureRecognizers []GestureRecognizer
}

func (me *gestureHandler) AddGestureRecoginer(recognizer GestureRecognizer) {
	if recognizer == nil {
		return
	}

	if me.gestureRecognizers == nil {
		me.gestureRecognizers = []GestureRecognizer{}
	}

	if true /*TODO debug*/ {
		name := util.GetTypeName(recognizer)
		for _, r := range me.gestureRecognizers {
			n := util.GetTypeName(r)
			if name == n {
				log.Fatal("nuxui", "The gesture recognizer '%s' is already registed.", name)
			}
		}
	}

	me.gestureRecognizers = append(me.gestureRecognizers, recognizer)
}

func (me *gestureHandler) RemoveGestureRecoginer(recognizer GestureRecognizer) {
	if recognizer == nil || me.gestureRecognizers == nil {
		return
	}

	for i, r := range me.gestureRecognizers {
		if r == recognizer {
			me.gestureRecognizers = append(me.gestureRecognizers[:i], me.gestureRecognizers[i+1:]...)
		}
	}
}

func (me *gestureHandler) HandlePointerEvent(event Event) {
	if me.gestureRecognizers != nil {
		for _, r := range me.gestureRecognizers {
			if r.PointerAllowed(event) {
				r.HandlePointerEvent(event)
			}
		}
	}
}

func (me *gestureHandler) FindGestureRecognizer(recognizerType GestureRecognizer) GestureRecognizer {
	name := util.GetTypeName(recognizerType)
	for _, r := range me.gestureRecognizers {
		if util.GetTypeName(r) == name  {
			return r
		}
	}

	return nil
}
