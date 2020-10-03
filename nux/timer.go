// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"math"
	"sync"
	"time"

	"github.com/nuxui/nuxui/log"
)

type timerLoop struct {
	mux     sync.Mutex
	seq     uint32
	ding    chan *timer
	cancel  map[uint32]struct{}
	running bool
}

var timerLoopInstance = &timerLoop{
	seq:     0,
	ding:    make(chan *timer),
	cancel:  map[uint32]struct{}{},
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
		if _, ok := me.cancel[t.seq]; ok {
			delete(me.cancel, t.seq)
			me.mux.Unlock()
			continue
		}
		me.mux.Unlock()

		if t.ui {
			RunOnUI(t.callback)
		} else {
			t.callback()
		}
	}
}

func (me *timerLoop) newSeq() uint32 {
	me.mux.Lock()
	defer me.mux.Unlock()
	if me.seq > math.MaxUint32-2 {
		me.seq = 0
	} else {
		me.seq++
	}
	return me.seq
}

func (me *timerLoop) cancelTimer(t *timer) {
	me.mux.Lock()
	me.cancel[t.seq] = struct{}{}
	me.mux.Unlock()
}

type Timer interface {
	Cancel()
}

type timer struct {
	ui       bool
	running  bool
	seq      uint32
	callback func()
}

func (me *timer) Cancel() {
	timerLoopInstance.cancelTimer(me)
	me.running = false
}

func NewInterval(duration time.Duration, callback func()) Timer {
	t := &timer{
		ui:       false,
		running:  true,
		seq:      timerLoopInstance.newSeq(),
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
		ui:       true,
		running:  true,
		seq:      timerLoopInstance.newSeq(),
		callback: callback,
	}
	go func(t *timer) {
		time.Sleep(duration)
		timerLoopInstance.ding <- t
	}(t)

	return t
}
