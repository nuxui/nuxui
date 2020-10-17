// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

/*
#import <Cocoa/Cocoa.h>

CGEventRef eventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon);

CFMachPortRef newEventTap(CGEventTapLocation tap, CGEventTapPlacement place, CGEventTapOptions options, CGEventMask eventsOfInterest, uint64_t index);
*/
import "C"
import (
	"sync"
)

type EventTapLocation C.uint32_t
type EventTapPlacement C.uint32_t
type EventTapOptions C.uint32_t
type EventType C.uint32_t
type EventMask C.uint64_t

// https://developer.apple.com/documentation/coregraphics/cgeventtaplocation?language=objc
const (
	HIDEventTap              EventTapLocation = C.kCGHIDEventTap
	SessionEventTap          EventTapLocation = C.kCGSessionEventTap
	AnnotatedSessionEventTap EventTapLocation = C.kCGAnnotatedSessionEventTap
)

// https://developer.apple.com/documentation/coregraphics/cgeventtapplacement?language=objc
const (
	HeadInsertEventTap EventTapPlacement = C.kCGHeadInsertEventTap
	TailAppendEventTap EventTapPlacement = C.kCGTailAppendEventTap
)

// https://developer.apple.com/documentation/coregraphics/cgeventtapoptions?language=objc
const (
	EventTapOptionDefault    EventTapOptions = C.kCGEventTapOptionDefault
	EventTapOptionListenOnly EventTapOptions = C.kCGEventTapOptionListenOnly
)

// https://developer.apple.com/documentation/coregraphics/cgeventtype?changes=lat_3_1_4__5&language=objc
const (
	EventNull              EventType = C.kCGEventNull
	EventLeftMouseDown     EventType = C.kCGEventLeftMouseDown
	EventLeftMouseUp       EventType = C.kCGEventLeftMouseUp
	EventRightMouseDown    EventType = C.kCGEventRightMouseDown
	EventRightMouseUp      EventType = C.kCGEventRightMouseUp
	EventMouseMoved        EventType = C.kCGEventMouseMoved
	EventLeftMouseDragged  EventType = C.kCGEventLeftMouseDragged
	EventRightMouseDragged EventType = C.kCGEventRightMouseDragged
	EventKeyDown           EventType = C.kCGEventKeyDown
	EventKeyUp             EventType = C.kCGEventKeyUp
	EventFlagsChanged      EventType = C.kCGEventFlagsChanged
	EventScrollWheel       EventType = C.kCGEventScrollWheel
	EventTabletPointer     EventType = C.kCGEventTabletPointer
	EventTabletProximity   EventType = C.kCGEventTabletProximity
	EventOtherMouseDown    EventType = C.kCGEventOtherMouseDown
	EventOtherMouseUp      EventType = C.kCGEventOtherMouseUp
	EventOtherMouseDragged EventType = C.kCGEventOtherMouseDragged
	// EventTapDisabledByTimeout   EventType = C.kCGEventTapDisabledByTimeout
	// EventTapDisabledByUserInput EventType = C.kCGEventTapDisabledByUserInput
)

// https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.8.sdk/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/CGEventTypes.h
const (
	EventMaskForAllEvents EventMask = 0xFFFFFFFFFFFFFFFF // C.kCGEventMaskForAllEvents
)

// https://developer.apple.com/documentation/coregraphics/1454426-cgeventtapcreate
type EventTapCallBack func(event Event, data interface{})

type MachPort interface{}

type machPort struct {
	ptr C.CFMachPortRef
}

type callbackAndData struct {
	callback EventTapCallBack
	data     interface{}
}

var callbacks map[uint64]callbackAndData = map[uint64]callbackAndData{}
var callbackIndex uint64 = 0
var callbackIndexMutex sync.Mutex

func newIndex() uint64 {
	callbackIndexMutex.Lock()
	callbackIndex++
	callbackIndexMutex.Unlock()
	return callbackIndex
}

func NewEventTap(tap EventTapLocation, place EventTapPlacement, options EventTapOptions, eventsOfInterest EventMask, callback EventTapCallBack, userInfo interface{}) MachPort {
	m := &machPort{}
	// t.ptr = C.CGEventTapCreate(C.CGEventTapLocation(tap), C.CGEventTapPlacement(place), C.CGEventTapOptions(options), C.uint64_t(eventsOfInterest), C.eventCallback, unsafe.Pointer(t))
	index := newIndex()
	callbacks[index] = callbackAndData{
		callback: callback,
		data:     userInfo,
	}
	m.ptr = C.newEventTap(C.CGEventTapLocation(tap), C.CGEventTapPlacement(place), C.CGEventTapOptions(options), C.uint64_t(eventsOfInterest), C.uint64_t(index))
	return m
}

//export eventCallback_go
func eventCallback_go(proxy C.CGEventTapProxy, eventRef C.CGEventRef, index C.uint64_t) C.CGEventRef {
	e := &event{ptr: eventRef}
	cbAndData := callbacks[uint64(index)]
	cbAndData.callback(e, cbAndData.data)
	return eventRef
}

//////////////////////////////////////////////////////////////////////////////////////
type RunLoopMode C.CFStringRef

var RunLoopCommonModes RunLoopMode = RunLoopMode(C.kCFRunLoopCommonModes)
var RunLoopDefaultMode RunLoopMode = RunLoopMode(C.kCFRunLoopDefaultMode)

type RunLoopSource interface{}
type runLoopSource struct {
	ptr C.CFRunLoopSourceRef
}

type RunLoop interface {
	AddSource(source RunLoopSource, mode RunLoopMode)
}
type runLoop struct {
	ptr C.CFRunLoopRef
}

type Allocator C.CFAllocatorRef

// type MachPort C.CFMachPortRef
type Index C.CFIndex

var AllocatorDefault = Allocator(C.kCFAllocatorDefault)

func NewRunLoopSource(allocator Allocator, port MachPort, order Index) RunLoopSource {
	s := &runLoopSource{}
	s.ptr = C.CFMachPortCreateRunLoopSource(C.CFAllocatorRef(allocator), C.CFMachPortRef(port.(*machPort).ptr), C.CFIndex(order))
	return s
}

func CurrentRunLoop() RunLoop {
	r := &runLoop{}
	r.ptr = C.CFRunLoopGetCurrent()
	return r
}

func (me *runLoop) AddSource(source RunLoopSource, mode RunLoopMode) {
	C.CFRunLoopAddSource(me.ptr, source.(*runLoopSource).ptr, C.CFRunLoopMode(mode))
}
