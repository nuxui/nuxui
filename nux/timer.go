// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"time"

	"github.com/nuxui/nuxui/log"
)

type timerLoop struct {
	ding    chan *timer
	running bool
}

var timerLoopInstance = &timerLoop{
	ding:    make(chan *timer),
	running: false,
}

func (me *timerLoop) init() {
	if !me.running {
		go me.loop()
	}
}

func (me *timerLoop) loop() {
	me.running = true
	for {
		t := <-me.ding
		log.V("nuxui", "t := <-me.ding 0")

		if !t.running {
			continue
		}

		log.V("nuxui", "t := <-me.ding 1")
		if t.backui {
			RunOnUI(t.callback)
		} else {
			go t.callback()
		}
	}
}

type Timer interface {
	Cancel()
	Running() bool
}

type timer struct {
	backui   bool
	running  bool
	callback func()
}

func (me *timer) Cancel() {
	me.running = false
}

func (me *timer) Running() bool {
	return me.running
}

func NewInterval(duration time.Duration, callback func()) Timer {
	t := &timer{
		backui:   false,
		running:  true,
		callback: callback,
	}

	go func(t *timer) {
		for t.running {
			time.Sleep(duration)
			timerLoopInstance.ding <- t
		}
	}(t)

	return t
}

func NewTimerBackToUI(duration time.Duration, callback func()) Timer {
	log.V("nuxui", "callback = %p", &callback)
	t := &timer{
		backui:   true,
		running:  true,
		callback: callback,
	}
	go func(t *timer) {
		time.Sleep(duration)
		log.V("nuxui", "timerLoopInstance.ding <- t 0")
		timerLoopInstance.ding <- t
		log.V("nuxui", "timerLoopInstance.ding <- t 1")

	}(t)

	return t
}
