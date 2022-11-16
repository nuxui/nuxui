// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package log


import (
	"io"
	"time"
	"fmt"
)

func new(out io.Writer, prefix string, flags int, depth int) Logger {
	me := &logger{
		depth:  depth,
		out:    out,
		flags:  flags,
		prefix: prefix,
		level:  VERBOSE,
		logs:   make(chan string, lBufferSize),
		timer:  map[uint32]time.Time{},
	}

	return me
}

func (me *logger) output(depth int, color string, level Level, levelTag string, tag string, format string, msg ...any) {
	if me.level > level {
		return
	}

	now := time.Now()

	var prefix string
	var dformat string

	if me.prefix != "" {
		prefix = " " + me.prefix
	}

	if me.flags&LUTC != 0 {
		now = now.UTC()
	}
	if me.flags&Ldate != 0 {
		dformat += "2006-01-02"
	} else {
		dformat += "01-02"
	}

	if me.flags&Ltime != 0 {
		dformat += " 15:04:05"

		if me.flags&Lmicroseconds != 0 {
			dformat += ".000000"
		} else {
			dformat += ".000"
		}
	}

	str := fmt.Sprintf(format, msg...)
	log := fmt.Sprintf("%s %s%s %s %s", now.Format(dformat), levelTag, prefix, tag, str)
	println(log)
}
