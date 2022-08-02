// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package ios

/*
#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>

NSUInteger nux_UIEvent_type(uintptr_t uievent) {
  return [(UIEvent *)uievent type];
}

short nux_UIEvent_subtype(uintptr_t uievent) {
  return [(UIEvent *)uievent subtype];
}

uintptr_t nux_UIEvent_allTouches(uintptr_t uievent){
  return (uintptr_t)[[(UIEvent *)uievent allTouches] allObjects];
}

void nux_UITouch_locationInView(uintptr_t uitouch, uintptr_t uiview, CGFloat* outX, CGFloat* outY){
	CGPoint p = [(UITouch *)uitouch locationInView: (UIView*)uiview];
	if (outX) { *outX = p.x; };
	if (outY) { *outY = p.y; };
}

UITouchPhase nux_UITouch_phase(uintptr_t uitouch){
	return [(UITouch *)uitouch phase];
}
*/
import "C"

const (
	UIEventTypeTouches       UIEventType = C.UIEventTypeTouches
	UIEventTypeMotion        UIEventType = C.UIEventTypeMotion
	UIEventTypeRemoteControl UIEventType = C.UIEventTypeRemoteControl
	UIEventTypePresses       UIEventType = C.UIEventTypePresses
	UIEventTypeHover         UIEventType = C.UIEventTypeHover
	UIEventTypeScroll        UIEventType = C.UIEventTypeScroll
	UIEventTypeTransform     UIEventType = C.UIEventTypeTransform
)

const (
	UIEventSubtypeNone                              UIEventSubtype = C.UIEventSubtypeNone
	UIEventSubtypeMotionShake                       UIEventSubtype = C.UIEventSubtypeMotionShake
	UIEventSubtypeRemoteControlPlay                 UIEventSubtype = C.UIEventSubtypeRemoteControlPlay
	UIEventSubtypeRemoteControlPause                UIEventSubtype = C.UIEventSubtypeRemoteControlPause
	UIEventSubtypeRemoteControlStop                 UIEventSubtype = C.UIEventSubtypeRemoteControlStop
	UIEventSubtypeRemoteControlTogglePlayPause      UIEventSubtype = C.UIEventSubtypeRemoteControlTogglePlayPause
	UIEventSubtypeRemoteControlNextTrack            UIEventSubtype = C.UIEventSubtypeRemoteControlNextTrack
	UIEventSubtypeRemoteControlPreviousTrack        UIEventSubtype = C.UIEventSubtypeRemoteControlPreviousTrack
	UIEventSubtypeRemoteControlBeginSeekingBackward UIEventSubtype = C.UIEventSubtypeRemoteControlBeginSeekingBackward
	UIEventSubtypeRemoteControlEndSeekingBackward   UIEventSubtype = C.UIEventSubtypeRemoteControlEndSeekingBackward
	UIEventSubtypeRemoteControlBeginSeekingForward  UIEventSubtype = C.UIEventSubtypeRemoteControlBeginSeekingForward
	UIEventSubtypeRemoteControlEndSeekingForward    UIEventSubtype = C.UIEventSubtypeRemoteControlEndSeekingForward
)

const (
	UIKeyModifierAlphaShift UIKeyModifierFlags = C.UIKeyModifierAlphaShift
	UIKeyModifierShift      UIKeyModifierFlags = C.UIKeyModifierShift
	UIKeyModifierControl    UIKeyModifierFlags = C.UIKeyModifierControl
	UIKeyModifierAlternate  UIKeyModifierFlags = C.UIKeyModifierAlternate
	UIKeyModifierCommand    UIKeyModifierFlags = C.UIKeyModifierCommand
	UIKeyModifierNumericPad UIKeyModifierFlags = C.UIKeyModifierNumericPad
)

const (
	UITouchTypeDirect   UITouchType = C.UITouchTypeDirect
	UITouchTypeIndirect UITouchType = C.UITouchTypeIndirect
	UITouchTypePencil   UITouchType = C.UITouchTypePencil
)

const (
	UITouchPhaseBegan         UITouchPhase = C.UITouchPhaseBegan
	UITouchPhaseMoved         UITouchPhase = C.UITouchPhaseMoved
	UITouchPhaseStationary    UITouchPhase = C.UITouchPhaseStationary
	UITouchPhaseEnded         UITouchPhase = C.UITouchPhaseEnded
	UITouchPhaseCancelled     UITouchPhase = C.UITouchPhaseCancelled
	UITouchPhaseRegionEntered UITouchPhase = C.UITouchPhaseRegionEntered
	UITouchPhaseRegionMoved   UITouchPhase = C.UITouchPhaseRegionMoved
	UITouchPhaseRegionExited  UITouchPhase = C.UITouchPhaseRegionExited
)

func (me UIEvent) Type() UIEventType {
	return UIEventType(C.nux_UIEvent_type(C.uintptr_t(me)))
}

// https://developer.apple.com/documentation/uikit/uievent/1613824-subtype?language=objc
func (me UIEvent) Subtype() UIEventSubtype {
	return UIEventSubtype(C.nux_UIEvent_subtype(C.uintptr_t(me)))
}

// https://developer.apple.com/documentation/uikit/uievent/1613836-alltouches?language=objc
func (me UIEvent) AllTouches() NSArray {
	return NSArray(C.nux_UIEvent_allTouches(C.uintptr_t(me)))
}

func (me UITouch) LocationInView(view UIView) (x, y float32) {
	var outX, outY C.CGFloat
	C.nux_UITouch_locationInView(C.uintptr_t(me), C.uintptr_t(view), &outX, &outY)
	return float32(outX), float32(outY)
}

func (me UITouch) Phase() UITouchPhase {
	return UITouchPhase(C.nux_UITouch_phase(C.uintptr_t(me)))
}
