// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pango

type Weight int

const (
	WEIGHT_THIN       Weight = 100
	WEIGHT_ULTRALIGHT Weight = 200
	WEIGHT_LIGHT      Weight = 300
	WEIGHT_SEMILIGHT  Weight = 350
	WEIGHT_BOOK       Weight = 380
	WEIGHT_NORMAL     Weight = 400
	WEIGHT_MEDIUM     Weight = 500
	WEIGHT_SEMIBOLD   Weight = 600
	WEIGHT_BOLD       Weight = 700
	WEIGHT_ULTRABOLD  Weight = 800
	WEIGHT_HEAVY      Weight = 900
	WEIGHT_ULTRAHEAVY Weight = 1000
)

type WrapMode int

const (
	WRAP_WORD WrapMode = iota
	WRAP_CHAR
	WRAP_WORD_CHAR
)
