// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "nuxui.org/nuxui/log"

// lifecycle

type OnCreate interface {
	OnCreate()
}

type Mount interface {
	OnMount()
}

type Eject interface {
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
	Measure(width, height MeasureDimen)
}

type Layout interface {
	Layout(x, y, width, height int32)
}

type Draw interface {
	Draw(Canvas)
}

func MountWidget(child Widget, parent Parent) {
	if child == nil {
		log.Fatal("nuxui", "child can not be nil when mount widget.")
		return
	}
	if parent == nil {
		log.Fatal("nuxui", "parent can not be nil when mount widget.")
		return
	}

	if (child.Info().Parent != nil) != child.Info().Mounted {
		log.Fatal("nuxui", "The widget %T:'%s' has wrong mount state with parent '%T'.", child, child.Info().ID, parent)
		return
	}
	mountWidget(child, parent)
}

func mountWidget(child Widget, parent Parent) {
	if child.Info().Mounted {
		log.Fatal("nuxui", "The widget '%T:%s' is already mounted to parent '%T'.", child, child.Info().ID, parent)
	}

	if parent != nil && !parent.Info().Mounted {
		// parent is not mounted, do nothing
		return
	}

	// mount child to parent
	child.Info().Parent = parent
	if f, ok := child.(Mount); ok {
		f.OnMount()
	}
	for _, m := range child.Info().Mixins {
		if mf, ok := m.(Mount); ok {
			mf.OnMount()
		}
	}
	child.Info().Mounted = true

	if p, ok := child.(Parent); ok {
		for _, c := range p.Children() {
			mountWidget(c, p)
		}
	}
	if compt, ok := child.(Component); ok {
		mountWidget(compt.Content(), parent)
	}
}

func EjectChild(child Widget) {
	if child == nil {
		log.Fatal("nuxui", "child can not be nil when eject widget.")
		return
	}

	if child.Info().Parent == nil || !child.Info().Mounted {
		log.Fatal("nuxui", "The widget '%s' is already ejected", child.Info().ID)
	}

	if f, ok := child.(Eject); ok {
		f.OnEject()
	}
	child.Info().Mounted = false
	child.Info().Parent = nil

	if p, ok := child.(Parent); ok {
		for _, c := range p.Children() {
			EjectChild(c)
		}
	}

	if compt, ok := child.(Component); ok {
		EjectChild(compt.Content())
	}
}
