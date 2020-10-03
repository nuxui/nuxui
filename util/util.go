// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"reflect"

	"github.com/nuxui/nuxui/log"
)

func GetTypeName(a interface{}) string {
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
func ReflectCall(ptr interface{}, funcName string, args ...interface{}) {
	if len(funcName) == 0 {
		log.E("nuxui", "ReflectCall receive a empty func name")
		return
	}

	if r := funcName[0]; r < 'A' || r > 'Z' {
		log.E("nuxui", "ReflectCall can not execute '%s' unexport function '%s'", GetTypeName(ptr), funcName)
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

func Ref2Ptr(ptr interface{}) uintptr {
	// TODO::
	return 0
}
