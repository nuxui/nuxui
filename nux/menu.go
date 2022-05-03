// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Menu interface {
	ID() string
	SetID(id string)
	Title() string
	SetTitle(title string)
	Icon() Image
	SetIcon(image Image)
	Action() func()
	SetAction(action func())
	Key() string
	SetKey(key string)
	Parent() Menu
	Content() Widget    // custom view
	SetContent() Widget // custom view
}

type menu struct {
}

type MenuItem struct {
}
