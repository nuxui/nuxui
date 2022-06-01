// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#import <Cocoa/Cocoa.h>

NSUInteger nux_NSEvent_Type(uintptr_t nsevent) {
  return [(NSEvent *)nsevent type];
}

short nux_NSEvent_Subtype(uintptr_t nsevent) {
  return [(NSEvent *)nsevent subtype];
}

NSPoint nux_NSEvent_LocationInWindow(uintptr_t nsevent){
  NSPoint p = [(NSEvent *)nsevent locationInWindow];
	p.y = [[(NSEvent *)nsevent window] contentView].bounds.size.height - p.y;
  return p;
}

uintptr_t nux_NSEvent_Window(uintptr_t nsevent){
  return (uintptr_t)[(NSEvent *)nsevent window];
}

NSInteger nux_NSEvent_ButtonNumber(uintptr_t nsevent){
  return [(NSEvent *)nsevent buttonNumber];
}

BOOL nux_NSEvent_HasPreciseScrollingDeltas(uintptr_t nsevent){
  return [(NSEvent *)nsevent hasPreciseScrollingDeltas];
}

CGFloat nux_NSEvent_ScrollingDeltaX(uintptr_t nsevent){
  return [(NSEvent *)nsevent scrollingDeltaX];
}

CGFloat nux_NSEvent_ScrollingDeltaY(uintptr_t nsevent){
  return [(NSEvent *)nsevent scrollingDeltaY];
}

unsigned short nux_NSEvent_KeyCode(uintptr_t nsevent){
  return [(NSEvent *)nsevent keyCode];
}

NSUInteger nux_NSEvent_ModifierFlags(uintptr_t nsevent){
  return [(NSEvent *)nsevent modifierFlags];
}

BOOL nux_NSEvent_IsARepeat(uintptr_t nsevent){
  return [(NSEvent *)nsevent isARepeat];
}

char* nux_NSEvent_Characters(uintptr_t nsevent){
  return (char*)[[(NSEvent *)nsevent characters] UTF8String];
}

float nux_NSEvent_Pressure(uintptr_t nsevent){
  return [(NSEvent *)nsevent pressure];
}

NSInteger nux_NSEvent_Stage(uintptr_t nsevent){
  return [(NSEvent *)nsevent stage];
}
