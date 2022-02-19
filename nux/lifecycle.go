// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

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
	mountWidget(child, parent)
}

func mountWidget(child Widget, parent Parent) {
	if parent == nil {
		log.I("nuxui", "mountWidget child=%s, parent=nil", child.Info().ID)
	} else {
		log.I("nuxui", "mountWidget child=%s, parent=%s", child.Info().ID, parent.Info().ID)
	}

	if (child.Info().Parent != nil) != child.Info().Mounted {
		log.E("nuxui", "The widget %T:'%s' has wrong mount state with parent '%T'.", child, child.Info().ID, parent)
		return
	}

	if child.Info().Mounted {
		log.Fatal("nuxui", "The widget '%s' is already mounted to parent '%T'.", child.Info().ID, parent)
	}

	if parent != nil && !parent.Info().Mounted {
		// parent is not mounted, do nothing
		// log.I("nuxui", " return  == child.(OnMount)=%s,  %T", child.Info().ID, child)
		return
	}

	child.Info().Parent = parent
	// log.I("nuxui", "child.(OnMount)=%s,  %T", child.Info().ID, child)
	if f, ok := child.(OnMount); ok {
		// log.I("nuxui", "child.(OnMount)=true", child.Info().ID)
		f.OnMount()
	}
	child.Info().Mounted = true
	for _, m := range child.Info().Mixins {
		if mf, ok := m.(OnMount); ok {
			mf.OnMount()
		}
	}

	if p, ok := child.(Parent); ok {
		for _, c := range p.Children() {
			// log.I("nuxui", "child.(OnMount)=%s,  %T", child.Info().ID, child)
			if f, ok := c.(OnMount); ok {
				f.OnMount()
			}
			c.Info().Mounted = true
			if c.Info().ID == "xxx" {
				log.I("nuxui", "xxx mixins %s", c.Info().Mixins)
			}
			for _, m := range c.Info().Mixins {
				if mf, ok := m.(OnMount); ok {
					mf.OnMount()
				}
			}

			if compt, ok := c.(Component); ok {
				c = compt.Content()
			}
			mountWidget(c, p)
		}
	}
}

func EjectChild(child Widget) {
	if child == nil {
		log.Fatal("nuxui", "child can not be nil when eject widget.")
		return
	}
	ejectChild(child)
}

func ejectChild(child Widget) {
	if child.Info().Parent == nil {
		log.Fatal("nuxui", "The widget '%s' is already ejected", child.Info().ID)
	}

	child.Info().Parent = nil
	if f, ok := child.(OnEject); ok {
		f.OnEject()
		child.Info().Mounted = false
	}

	if p, ok := child.(Parent); ok {
		for _, c := range p.Children() {
			if f, ok := c.(OnEject); ok {
				f.OnEject()
				c.Info().Mounted = false
			}
			if compt, ok := c.(Component); ok {
				c = compt.Content()
			}
			ejectChild(c)
		}
	}
}
