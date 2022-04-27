// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"reflect"

	"github.com/nuxui/nuxui/log"
)

func TypeName(a any) string {
	if a == nil {
		return ""
	}

	t := reflect.TypeOf(a)
	if t.Name() == "" {
		t = t.Elem()
	}

	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}

func Absi32(v int32) int32 {
	if v < 0 {
		v = 0 - v
	}
	return v
}

// TODO:: unit test ptr and args
func ReflectCall(ptr any, funcName string, args ...any) {
	if len(funcName) == 0 {
		log.E("nuxui", "ReflectCall receive a empty func name")
		return
	}

	if r := funcName[0]; r < 'A' || r > 'Z' {
		log.E("nuxui", "ReflectCall can not execute '%s' unexport function '%s'", TypeName(ptr), funcName)
		return
	}

	method := reflect.ValueOf(ptr).MethodByName(funcName)

	if method.Type().NumIn() != len(args) {
		log.Fatal("nuxui", "ReflectCall need %d arguments but receive %d aguments", method.Type().NumIn(), len(args))
	}

	in := make([]reflect.Value, method.Type().NumIn())
	for i := 0; i < method.Type().NumIn(); i++ {
		in[i] = reflect.ValueOf(args[i])
	}
	method.Call(in)
}

func Ref2Ptr(ptr any) uintptr {
	// TODO::
	return 0
}

// SameFunc compare two func and return true if equal
func SameFunc(a, b any) bool {
	// TODO:: check a, b is func
	if a == nil || b == nil {
		return a == b
	}
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

func MaxF(values ...float32) (max float32) {
	if len(values) > 0 {
		max = values[0]
	}
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return
}

func MinF(values ...float32) (min float32) {
	if len(values) > 0 {
		min = values[0]
	}
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return
}
