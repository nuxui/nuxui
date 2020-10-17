// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#import <Cocoa/Cocoa.h>

CGEventRef eventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon)
{
	// printf("####### myCGEventCallback...  %d\n", i);
    return eventCallback_go(proxy, event, (uint64_t)refcon);
	// return event;
}

CFMachPortRef newEventTap(CGEventTapLocation tap, CGEventTapPlacement place, CGEventTapOptions options, CGEventMask eventsOfInterest, uint64_t index){
	return CGEventTapCreate(tap, place, options, eventsOfInterest, eventCallback, (void*)index);
}