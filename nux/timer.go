// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"sync"
	"time"

	"github.com/nuxui/nuxui/log"
)

type timerLoop struct {
	mux     sync.Mutex
	n       int32
	ding    chan *timer
	cancel  map[int32]struct{}
	running bool
}

var timerLoopInstance = &timerLoop{
	n:       0,
	ding:    make(chan *timer),
	cancel:  map[int32]struct{}{},
	running: false,
}

func (me *timerLoop) init() {
	me.mux.Lock()
	if !me.running {
		go me.loop()
	}
	me.mux.Unlock()
}

func (me *timerLoop) loop() {
	me.running = true
	for {
		t := <-me.ding

		me.mux.Lock()
		if _, ok := me.cancel[t.id]; ok {
			delete(me.cancel, t.id)
			me.mux.Unlock()
			continue
		}
		me.mux.Unlock()

		//RunOnUI(t.callback)
	}
}

func (me *timerLoop) newId() int32 {
	me.mux.Lock()
	defer me.mux.Unlock()
	me.n++
	if me.n == 0 {
		me.n++
	}
	return me.n
}

func (me *timerLoop) cancelTimer(t *timer) {
	me.mux.Lock()
	me.cancel[t.id] = struct{}{}
	me.mux.Unlock()
}

type Timer interface {
	Stop()
}

type timer struct {
	id       int32
	callback func()
}

func (me *timer) Stop() {
	timerLoopInstance.cancelTimer(me)
}

func NewTimerBackToUI(duration time.Duration, callback func()) Timer {
	timerLoopInstance.init()
	log.V("nux", fmt.Sprintf("callback = %p", &callback))
	t := &timer{
		id:       timerLoopInstance.newId(),
		callback: callback,
	}
	go func(t *timer) {
		time.Sleep(duration)
		timerLoopInstance.ding <- t
	}(t)

	return t
}
