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

uintptr_t nux_NSCursor_ArrowCursor(){
	return (uintptr_t)[NSCursor arrowCursor];
}

uintptr_t nux_NSCursor_IBeamCursor(){
	return (uintptr_t)[NSCursor IBeamCursor];
}
*/
import "C"

func NSCursor_IBeamCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_IBeamCursor())
}

func NSCursor_ArrowCursor() NSCursor {
	return NSCursor(C.nux_NSCursor_ArrowCursor())
}

func (me NSCursor) Set() {
	C.nux_NSCursor_set(C.uintptr_t(me))
}
