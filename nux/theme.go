// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
)

type ThemeKind int

const (
	ThemeLight ThemeKind = iota
	ThemeDark
)

type Theme interface {
	Kind() ThemeKind
	GetStyle(names ...string) Attr
}

func AppTheme() Theme {
	return appTheme
}

func ApplyTheme(kind ThemeKind, themes ...string) {
	// TODO:: call at main func, but app is not running, isMainThread is not work
	if false && !IsMainThread() {
		log.Fatal("nuxui", "ApplyTheme can only run at main thread")
	}

	styles := make([]Attr, len(themes))
	for i, theme := range themes {
		styles[i] = InflateStyle(theme)
	}

	appTheme = &theme{kind: kind, styles: styles}
}

var appTheme Theme = &theme{kind: ThemeLight, styles: []Attr{}}

type theme struct {
	kind   ThemeKind
	styles []Attr
}

func (me *theme) Kind() ThemeKind {
	return me.kind
}

func (me *theme) GetStyle(names ...string) (attr Attr) {
	for _, name := range names {
		for _, style := range me.styles {
			if style.Has(name) {
				attr = MergeAttrs(attr, style.GetAttr(name, nil))
			}
		}
	}
	return
}
