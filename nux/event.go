// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "time"

type Event interface {
	ID() uint64
	Time() time.Time
	Type() EventType
	Action() EventAction
}

type event struct {
	id     uint64
	time   time.Time
	etype  EventType
	action EventAction
}

func (me *event) ID() uint64 {
	return me.id
}

func (me *event) Time() time.Time {
	return me.time
}

func (me *event) Type() EventType {
	return me.etype
}

func (me *event) Action() EventAction {
	return me.action
}
