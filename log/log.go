// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	VERBOSE = iota + 2
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type Logger interface {
	V(tag string, format string, msg ...interface{})
	D(tag string, format string, msg ...interface{})
	I(tag string, format string, msg ...interface{})
	W(tag string, format string, msg ...interface{})
	E(tag string, format string, msg ...interface{})
	Fatal(tag string, format string, msg ...interface{}) // TODO:: print log, print stack, exit
	SetOutput(w io.Writer)
	Flags() int
	SetFlags(flag int)
	Prefix() string
	SetPrefix(prefix string)
	Close() // defer call at main function, cause panic will drop logs
}

const lBufferSize = 20

const (
	Ldate         = 1 << iota // the date in the local time zone: 2009/01/23
	Ltime                     // the time in the local time zone: 01:23:23
	Lmicroseconds             // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                 // full file name and line number: /a/b/c/d.go:23
	Lshortfile                // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                      // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ltime     // initial values for the standard logger
)

type logger struct {
	size   uint
	flags  int
	depth  int
	prefix string
	mux    sync.Mutex
	out    io.Writer
	logs   chan string
}

func New(out io.Writer, prefix string, flags int) Logger {
	return new(out, prefix, flags, 1)
}

func (me *logger) V(tag string, format string, msg ...interface{}) {
	me.output(me.depth, VERBOSE, "V", tag, format, msg...)
}

func (me *logger) D(tag string, format string, msg ...interface{}) {
	me.output(me.depth, DEBUG, "D", tag, format, msg...)
}

func (me *logger) I(tag string, format string, msg ...interface{}) {
	me.output(me.depth, INFO, "I", tag, format, msg...)
}

func (me *logger) W(tag string, format string, msg ...interface{}) {
	me.output(me.depth, WARN, "W", tag, format, msg...)
}

func (me *logger) E(tag string, format string, msg ...interface{}) {
	me.output(me.depth, ERROR, "E", tag, format, msg...)
}

func (me *logger) Fatal(tag string, format string, msg ...interface{}) {
	me.output(me.depth, ERROR, "Fatal", tag, format, msg...)
	panic(fmt.Sprintf(format, msg...))
}

func (me *logger) SetOutput(w io.Writer) {
	me.out = w
}

// Flags returns the output flags for the logger.
func (me *logger) Flags() int {
	return me.flags
}

// SetFlags sets the output flags for the logger.
func (me *logger) SetFlags(flags int) {
	me.flags = flags
}

// Prefix returns the output prefix for the logger.
func (me *logger) Prefix() string {
	return me.prefix
}

// SetPrefix sets the output prefix for the logger.
func (me *logger) SetPrefix(prefix string) {
	me.prefix = prefix
}

func (me *logger) Close() {
	for len(me.logs) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
}

var std = new(os.Stdout, "", LstdFlags, 2)

func V(tag string, format string, msg ...interface{}) {
	std.V(tag, format, msg...)
}

func D(tag string, format string, msg ...interface{}) {
	std.D(tag, format, msg...)
}

func I(tag string, format string, msg ...interface{}) {
	std.I(tag, format, msg...)
}

func W(tag string, format string, msg ...interface{}) {
	std.W(tag, format, msg...)
}

func E(tag string, format string, msg ...interface{}) {
	std.E(tag, format, msg...)
}

func Fatal(tag string, format string, msg ...interface{}) {
	std.Fatal(tag, format, msg...)
}

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// Flags returns the output flags for the standard logger.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the output flags for the standard logger.
func SetFlags(flags int) {
	std.SetFlags(flags)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func Close() {
	std.Close()
}
