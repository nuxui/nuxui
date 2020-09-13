// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

// life cycle
type Creating interface {
	// Creating(Template, Attr)
	Creating(attr Attr)
}

type Created interface {
	Created(content Widget)
}

type AnyCreated interface {
	Created(data interface{})
}

// type Mounting interface {
// 	Mounting()
// }

type Mounted interface {
	Mounted()
}

// Active, it does not mean Focused, it should be Actived and it can also run animation
type Actived interface {
	Actived()
}

//Deactived window widget lost focus
type Deactived interface {
	Deactived()
}

// type Unmounting interface {
// 	Unounting()
// }

// widnow hide, go background
// widget removed from parent
type Unmounted interface {
	Unounted()
}

type Destroying interface {
	Destroying()
}

type Destroyed interface {
	Destroyed()
}

type Measure interface {
	Measure(width, height int32)
}

type Layout interface {
	Layout(dx, dy, left, top, right, bottom int32)
}

type Draw interface {
	Draw(Canvas)
}
