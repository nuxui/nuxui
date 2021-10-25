// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

// lifecycle

type OnCreate interface {
	OnCreate()
}

type OnMount interface {
	OnMount()
}

type OnEject interface {
	OnEject()
}

// // Active, it does not mean Focused, it should be Actived and it can also run animation
// type Actived interface {
// 	Actived()
// }

// //Deactived window widget lost focus
// type Deactived interface {
// 	Deactived()
// }

type Measure interface {
	Measure(width, height int32)
}

type Layout interface {
	Layout(dx, dy, left, top, right, bottom int32)
}

type Draw interface {
	Draw(Canvas)
}
