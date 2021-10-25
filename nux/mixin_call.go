// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

// func executeCreating(widget Widget, layout Attr) {
// 	if creating, ok := widget.(Creating); ok {
// 		creating.Creating(layout)
// 	}

// 	for _, mixin := range Mixins(widget) {
// 		if creating, ok := mixin.(Creating); ok {
// 			creating.Creating(layout)
// 		}
// 	}
// }

// func executeCreated(widget Widget) {
// 	if created, ok := widget.(Created); ok {
// 		created.Created(widget)
// 	}

// 	for _, mixin := range Mixins(widget) {
// 		if created, ok := mixin.(Created); ok {
// 			created.Created(widget)
// 		}
// 	}
// }

// func executeDestroy(widget Widget) {
// 	GestureBinding().ClearGestureHandler(widget)

// 	for _, mixin := range Mixins(widget) {
// 		if f, ok := mixin.(Destroyed); ok {
// 			f.Destroyed()
// 		}
// 	}

// 	if f, ok := widget.(Destroyed); ok {
// 		f.Destroyed()
// 	}

// 	ClearMixins(widget)
// }
