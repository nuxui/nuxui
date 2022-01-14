// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (darwin || windows || linux) && !android

package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
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

	go func() {
		for {
			me.out.Write([]byte(<-me.logs))
		}
	}()

	return me
}

func (me *logger) output(depth int, color string, level Level, levelTag string, tag string, format string, msg ...interface{}) {
	if me.level > level {
		return
	}

	now := time.Now()

	me.mux.Lock()
	defer me.mux.Unlock()

	depth++
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

	var log string
	if me.flags&Lshortfile != 0 {
		_, file, _, ok := runtime.Caller(depth)
		if !ok {
			file = "???"
		}

		if me.flags&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}

		if me.out == os.Stdout {
			log = fmt.Sprintf(colorTemplate, now.Format(dformat), color, fmt.Sprintf("%s%s %s %s %s\n", levelTag, prefix, file, tag, fmt.Sprintf(format, msg...)))
		} else {
			log = fmt.Sprintf("%s %s%s %s %s %s\n", now.Format(dformat), levelTag, prefix, file, tag, fmt.Sprintf(format, msg...))
		}
	} else if me.flags&Llongfile != 0 {
		frame := getFrameName(depth)
		if me.out == os.Stdout {
			log = fmt.Sprintf(colorTemplate, now.Format(dformat), color, fmt.Sprintf("%s%s %s %s %s\n", levelTag, prefix, frame, tag, fmt.Sprintf(format, msg...)))
		} else {
			log = fmt.Sprintf("%s %s%s %s %s %s\n", now.Format(dformat), levelTag, prefix, frame, tag, fmt.Sprintf(format, msg...))
		}
	} else {
		if me.out == os.Stdout {
			log = fmt.Sprintf(colorTemplate, now.Format(dformat), color, fmt.Sprintf("%s%s %s %s\n", levelTag, prefix, tag, fmt.Sprintf(format, msg...)))
		} else {
			log = fmt.Sprintf("%s %s%s %s %s\n", now.Format(dformat), levelTag, prefix, tag, fmt.Sprintf(format, msg...))
		}
	}
	me.logs <- log

}

func getFrameName(skipFrames int) string {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrameName
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "? ? ?"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame.Function
}
