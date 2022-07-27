// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>
#include <pthread.h>

void nux_NSCursor_set(uintptr_t nscursor){
	[(NSCursor*)nscursor set];
}

uintptr_t nux_NSCursor_arrowCursor(){
	return (uintptr_t)[NSCursor arrowCursor];
}

uintptr_t nux_NSCursor_IBeamCursor(){
	return (uintptr_t)[NSCursor IBeamCursor];
}

uintptr_t nux_NSCursor_pointingHandCursor(){
	return (uintptr_t)[NSCursor pointingHandCursor];
}

uintptr_t nux_NSCursor_crosshairCursor(){
	return (uintptr_t)[NSCursor crosshairCursor];
}

uintptr_t nux_NSCursor_openHandCursor(){
	return (uintptr_t)[NSCursor openHandCursor];
}

uintptr_t nux_NSCursor_closedHandCursor(){
	return (uintptr_t)[NSCursor closedHandCursor];
}
*/
import "C"

func NSCursor_IBeamCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_IBeamCursor())
}

func NSCursor_ArrowCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_arrowCursor())
}

func NSCursor_PointingHandCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_pointingHandCursor())
}

func NSCursor_CrosshairCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_crosshairCursor())
}

func NSCursor_OpenHandCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_openHandCursor())
}

func NSCursor_ClosedHandCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_closedHandCursor())
}

func (me NSCursor) Set() {
	C.nux_NSCursor_set(C.uintptr_t(me))
}
