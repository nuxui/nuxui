// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

type GestureArenaMember interface {
	RejectGesture(pointer int64)
	AccpetGesture(pointer int64)
}

type gestureArena struct {
	members         []GestureArenaMember
	isOpen          bool
	isHeld          bool
	hasPendingSweep bool
	eagerWinner     GestureArenaMember
}

func newGestureArena() *gestureArena {
	return &gestureArena{
		members:         []GestureArenaMember{},
		isOpen:          true,
		isHeld:          false,
		hasPendingSweep: false,
		eagerWinner:     nil,
	}
}

func (me *gestureArena) add(member GestureArenaMember) {
	if member == nil {
		log.Fatal("nuxui", "Invalid GestureArenaMember, it can not be nil")
	}

	if !me.isOpen {
		log.Fatal("nuxui", "GestureArena is already closed, can not add any GestureArenaMember to here")
	}

	me.members = append(me.members, member)
}

func (me *gestureArena) remove(member GestureArenaMember) {
	for i, m := range me.members {
		if member == m {
			me.members = append(me.members[:i], me.members[i+1:]...)
			break
		}
	}
}
