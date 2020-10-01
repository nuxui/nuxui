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

var gestureArenaManagerInstance *gestureArenaManager = &gestureArenaManager{
	arenas: map[int64]*gestureArena{},
}

func GestureArenaManager() *gestureArenaManager {
	return gestureArenaManagerInstance
}

type gestureArenaManager struct {
	arenas map[int64]*gestureArena
}

func (me *gestureArenaManager) Add(pointer int64, member GestureArenaMember) {
	if arena, ok := me.arenas[pointer]; ok {
		arena.add(member)
	} else {
		arena := newGestureArena()
		arena.add(member)
		me.arenas[pointer] = arena
	}
}

func (me *gestureArenaManager) Close(pointer int64) {
	if arena, ok := me.arenas[pointer]; ok {
		arena.isOpen = false
		me.tryToResolveArena(pointer, arena)
	}
}

func (me *gestureArenaManager) Sweep(pointer int64) {
	if arena, ok := me.arenas[pointer]; ok {
		if arena.isHeld {
			arena.hasPendingSweep = true
			return
		}

		delete(me.arenas, pointer)

		if len(arena.members) != 0 {
			i := 0
			for _, v := range arena.members {
				if i == 0 {
					log.V("nuxui", "Sweep AccpetGesture %T", v)
					v.AccpetGesture(pointer)
				} else {
					log.V("nuxui", "Sweep RejectGesture %T", v)
					v.RejectGesture(pointer)
				}
				i++
			}
		}
	}
}

func (me *gestureArenaManager) Hold(pointer int64) (success bool) {
	if arena, ok := me.arenas[pointer]; ok {
		arena.isHeld = true
		return true
	}
	return false
}

func (me *gestureArenaManager) Release(pointer int64) {
	if arena, ok := me.arenas[pointer]; ok {
		arena.isHeld = false
		if arena.hasPendingSweep {
			me.Sweep(pointer)
		}
	}
}

func (me *gestureArenaManager) Resolve(pointer int64, member GestureArenaMember, accepted bool) {
	if arena, ok := me.arenas[pointer]; ok {
		if accepted {
			if arena.isOpen {
				if arena.eagerWinner == nil {
					arena.eagerWinner = member
				}
			} else {
				me.resolveInFavorOf(pointer, arena, member)
			}
		} else {
			arena.remove(member)
			log.V("nuxui", "arena remain = %d", len(arena.members))

			member.RejectGesture(pointer)
			if !arena.isOpen {
				me.tryToResolveArena(pointer, arena)
			}
		}
	}
}

func (me *gestureArenaManager) tryToResolveArena(pointer int64, arena *gestureArena) {
	if arena.isOpen {
		log.Fatal("nuxui", "the GestureArena should closed before resolve")
	}

	if len(arena.members) == 1 {
		me.resolveByDefault(pointer, arena)
	} else if len(arena.members) == 0 {
		delete(me.arenas, pointer)
	} else if arena.eagerWinner != nil {
		me.resolveInFavorOf(pointer, arena, arena.eagerWinner)
	}
}

func (me *gestureArenaManager) resolveByDefault(pointer int64, arena *gestureArena) {
	if arena.isOpen {
		log.Fatal("nuxui", "the GestureArena should closed before resolve")

	}

	if arena, ok := me.arenas[pointer]; ok {
		if len(arena.members) != 1 {
			log.Fatal("nuxui", "resolveByDefault should called when members number is 1")
		}

		delete(me.arenas, pointer)

		arena.members[0].AccpetGesture(pointer)
	}
}

func (me *gestureArenaManager) resolveInFavorOf(pointer int64, arena *gestureArena, member GestureArenaMember) {
	if arena.isOpen {
		log.Fatal("nuxui", "the GestureArena should closed before resolve")

	}

	delete(me.arenas, pointer)
	for _, m := range arena.members {
		if m != member {
			m.RejectGesture(pointer)
		}
	}
	member.AccpetGesture(pointer)
}
