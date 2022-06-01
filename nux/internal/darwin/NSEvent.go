// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

NSUInteger     nux_NSEvent_Type(uintptr_t nsevent);
short          nux_NSEvent_Subtype(uintptr_t nsevent);
NSPoint        nux_NSEvent_LocationInWindow(uintptr_t nsevent);
uintptr_t      nux_NSEvent_Window(uintptr_t nsevent);
NSInteger      nux_NSEvent_ButtonNumber(uintptr_t nsevent);
BOOL           nux_NSEvent_HasPreciseScrollingDeltas(uintptr_t nsevent);
CGFloat        nux_NSEvent_ScrollingDeltaX(uintptr_t nsevent);
CGFloat        nux_NSEvent_ScrollingDeltaY(uintptr_t nsevent);
unsigned short nux_NSEvent_KeyCode(uintptr_t nsevent);
NSUInteger     nux_NSEvent_ModifierFlags(uintptr_t nsevent);
BOOL           nux_NSEvent_IsARepeat(uintptr_t nsevent);
char*          nux_NSEvent_Characters(uintptr_t nsevent);
float          nux_NSEvent_Pressure(uintptr_t nsevent);
NSInteger      nux_NSEvent_Stage(uintptr_t nsevent);

*/
import "C"

const (
	NSEventTypeLeftMouseDown      NSEventType = C.NSEventTypeLeftMouseDown
	NSEventTypeLeftMouseUp        NSEventType = C.NSEventTypeLeftMouseUp
	NSEventTypeRightMouseDown     NSEventType = C.NSEventTypeRightMouseDown
	NSEventTypeRightMouseUp       NSEventType = C.NSEventTypeRightMouseUp
	NSEventTypeMouseMoved         NSEventType = C.NSEventTypeMouseMoved
	NSEventTypeLeftMouseDragged   NSEventType = C.NSEventTypeLeftMouseDragged
	NSEventTypeRightMouseDragged  NSEventType = C.NSEventTypeRightMouseDragged
	NSEventTypeMouseEntered       NSEventType = C.NSEventTypeMouseEntered
	NSEventTypeMouseExited        NSEventType = C.NSEventTypeMouseExited
	NSEventTypeKeyDown            NSEventType = C.NSEventTypeKeyDown
	NSEventTypeKeyUp              NSEventType = C.NSEventTypeKeyUp
	NSEventTypeFlagsChanged       NSEventType = C.NSEventTypeFlagsChanged
	NSEventTypeAppKitDefined      NSEventType = C.NSEventTypeAppKitDefined
	NSEventTypeSystemDefined      NSEventType = C.NSEventTypeSystemDefined
	NSEventTypeApplicationDefined NSEventType = C.NSEventTypeApplicationDefined
	NSEventTypePeriodic           NSEventType = C.NSEventTypePeriodic
	NSEventTypeCursorUpdate       NSEventType = C.NSEventTypeCursorUpdate
	NSEventTypeScrollWheel        NSEventType = C.NSEventTypeScrollWheel
	NSEventTypeTabletPoint        NSEventType = C.NSEventTypeTabletPoint
	NSEventTypeTabletProximity    NSEventType = C.NSEventTypeTabletProximity
	NSEventTypeOtherMouseDown     NSEventType = C.NSEventTypeOtherMouseDown
	NSEventTypeOtherMouseUp       NSEventType = C.NSEventTypeOtherMouseUp
	NSEventTypeOtherMouseDragged  NSEventType = C.NSEventTypeOtherMouseDragged
	NSEventTypeGesture            NSEventType = C.NSEventTypeGesture
	NSEventTypeMagnify            NSEventType = C.NSEventTypeMagnify
	NSEventTypeSwipe              NSEventType = C.NSEventTypeSwipe
	NSEventTypeRotate             NSEventType = C.NSEventTypeRotate
	NSEventTypeBeginGesture       NSEventType = C.NSEventTypeBeginGesture
	NSEventTypeEndGesture         NSEventType = C.NSEventTypeEndGesture
	NSEventTypeSmartMagnify       NSEventType = C.NSEventTypeSmartMagnify
	NSEventTypeQuickLook          NSEventType = C.NSEventTypeQuickLook
	NSEventTypePressure           NSEventType = C.NSEventTypePressure
	NSEventTypeDirectTouch        NSEventType = C.NSEventTypeDirectTouch
)

const (
	NSEventSubtypeWindowExposed          NSEventSubtype = C.NSEventSubtypeWindowExposed
	NSEventSubtypeApplicationActivated   NSEventSubtype = C.NSEventSubtypeApplicationActivated
	NSEventSubtypeApplicationDeactivated NSEventSubtype = C.NSEventSubtypeApplicationDeactivated
	NSEventSubtypeWindowMoved            NSEventSubtype = C.NSEventSubtypeWindowMoved
	NSEventSubtypeScreenChanged          NSEventSubtype = C.NSEventSubtypeScreenChanged
	NSEventSubtypePowerOff               NSEventSubtype = C.NSEventSubtypePowerOff
	NSEventSubtypeMouseEvent             NSEventSubtype = C.NSEventSubtypeMouseEvent
	NSEventSubtypeTabletPoint            NSEventSubtype = C.NSEventSubtypeTabletPoint
	NSEventSubtypeTabletProximity        NSEventSubtype = C.NSEventSubtypeTabletProximity
	NSEventSubtypeTouch                  NSEventSubtype = C.NSEventSubtypeTouch
)

const (
	NSEventModifierFlagCapsLock   NSEventModifierFlags = C.NSEventModifierFlagCapsLock
	NSEventModifierFlagShift      NSEventModifierFlags = C.NSEventModifierFlagShift
	NSEventModifierFlagControl    NSEventModifierFlags = C.NSEventModifierFlagControl
	NSEventModifierFlagOption     NSEventModifierFlags = C.NSEventModifierFlagOption
	NSEventModifierFlagCommand    NSEventModifierFlags = C.NSEventModifierFlagCommand
	NSEventModifierFlagNumericPad NSEventModifierFlags = C.NSEventModifierFlagNumericPad
	NSEventModifierFlagHelp       NSEventModifierFlags = C.NSEventModifierFlagHelp
	NSEventModifierFlagFunction   NSEventModifierFlags = C.NSEventModifierFlagFunction
)

func (me NSEvent) Type() NSEventType {
	return NSEventType(C.nux_NSEvent_Type(C.uintptr_t(me)))
}

// https://developer.apple.com/documentation/appkit/nsevent/1527726-subtype?language=objc
func (me NSEvent) Subtype() NSEventSubtype {
	return NSEventSubtype(C.nux_NSEvent_Subtype(C.uintptr_t(me)))
}

func (me NSEvent) LocationInWindow() (x, y float32) {
	point := C.nux_NSEvent_LocationInWindow(C.uintptr_t(me))
	return float32(point.x), float32(point.y)
}

func (me NSEvent) Window() NSWindow {
	return NSWindow(C.nux_NSEvent_Window(C.uintptr_t(me)))
}

func (me NSEvent) ButtonNumber() int32 {
	return int32(C.nux_NSEvent_ButtonNumber(C.uintptr_t(me)))
}

func (me NSEvent) HasPreciseScrollingDeltas() bool {
	return C.nux_NSEvent_HasPreciseScrollingDeltas(C.uintptr_t(me)) > 0
}

func (me NSEvent) ScrollingDeltaX() float32 {
	return float32(C.nux_NSEvent_ScrollingDeltaX(C.uintptr_t(me)))
}

func (me NSEvent) ScrollingDeltaY() float32 {
	return float32(C.nux_NSEvent_ScrollingDeltaY(C.uintptr_t(me)))
}

func (me NSEvent) KeyCode() uint32 {
	return uint32(C.nux_NSEvent_KeyCode(C.uintptr_t(me)))
}

func (me NSEvent) ModifierFlags() NSEventModifierFlags {
	return NSEventModifierFlags(C.nux_NSEvent_ModifierFlags(C.uintptr_t(me)))
}

func (me NSEvent) IsARepeat() bool {
	return C.nux_NSEvent_IsARepeat(C.uintptr_t(me)) > 0
}

func (me NSEvent) Characters() string {
	return C.GoString(C.nux_NSEvent_Characters(C.uintptr_t(me)))
}

func (me NSEvent) Pressure() float32 {
	return float32(C.nux_NSEvent_Pressure(C.uintptr_t(me)))
}

func (me NSEvent) Stage() int32 {
	return int32(C.nux_NSEvent_Stage(C.uintptr_t(me)))
}
