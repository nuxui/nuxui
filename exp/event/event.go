// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa -framework OpenGL
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>


*/
import "C"

type EventSourceRef interface{}
type eventSourceRef struct {
	ptr C.uintptr_t
}

type Event interface {
	Type() EventType
	SetType(etype EventType)
	SetIntegerValue(key int, value int64)
}

type event struct {
	ptr C.CGEventRef
}

func NewEvent(source EventSourceRef) Event {
	if source == nil {
		source = &eventSourceRef{ptr: C.uintptr_t(0)}
	}
	e := &event{}
	// e.ptr = C.CGEventCreate(((C.CGEventSourceRef)(0)))
	return e
}

func NewMouseEvent(point Point) Event {
	e := &event{}
	// e.ptr = C.CGEventCreateMouseEvent(nil, C.kCGEventRightMouseDown, point, C.kCGMouseButtonRight)
	return e
}

func (e *event) Type() EventType {
	return EventType(C.CGEventGetType(e.ptr))
}

func (e *event) SetType(etype EventType) {
	C.CGEventSetType(e.ptr, C.CGEventType(etype))
}

func (e *event) SetIntegerValue(key int, value int64) {
	C.CGEventSetIntegerValueField(C.CGEventRef(e.ptr), C.CGEventField(C.int32_t(key)), C.int64_t(value))
}
