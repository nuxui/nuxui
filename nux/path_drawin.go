// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

*/
import "C"

type path struct {
	ptr C.CGPathRef
}

func NewPath() Path {
	return &path{}
}

func (me *path) AddRoundRect(left, top, right, bottom, rx, ry float32) {

}

type gradient struct {
	ptr C.CGGradientRef
}

func NewLinearGradient(startX, startY, endX, endY float32) Gradient {
	me := &gradient{}
	return me
}
