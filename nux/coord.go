// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"math"

	"github.com/nuxui/nuxui/util"
)

type Coord interface {
	X() float32
	Y() float32
	Distance() float64
	DistanceSquared() float64
	Direction() float64
	Translate(x, y float32)
	Scale(x, y float32)
	Sub(Coord) Coord
	Equal(Coord) bool
	String() string
}

type coord struct {
	X int32
	Y int32
}

func (me *coord) Distance() float64 {
	return math.Sqrt(float64(me.X)*float64(me.X) + float64(me.Y)*float64(me.Y))
}

func (me *coord) DistanceSquared() float64 {
	return float64(me.X)*float64(me.X) + float64(me.Y)*float64(me.Y)
}

func (me *coord) Direction() float64 {
	return math.Atan2(float64(me.Y), float64(me.X))
}

func (me *coord) Translate(x, y int32) {
	me.X += x
	me.Y += y
}

func (me *coord) Scale(x, y float32) {
	me.X = util.Roundi32(float32(me.X) * x)
	me.Y = util.Roundi32(float32(me.X) * x)
}

func (me *coord) Sub(offset coord) coord {
	return coord{X: me.X - offset.X, Y: me.Y - offset.Y}
}

func (me *coord) Equal(offset coord) bool {
	return me.X == offset.X && me.Y == offset.Y
}

func (me *coord) String() string {
	return fmt.Sprintf("coord{X: %d, Y: %d}", me.X, me.Y)
}
