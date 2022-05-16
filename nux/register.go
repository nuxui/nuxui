// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"reflect"
	"strings"

	"nuxui.org/nuxui/log"
)

type Creator func(Attr) any

var _creatorList = make(map[string]Creator)

func RegisterType(niltype any, creator Creator) {
	if debug_register {
		typeName := fmt.Sprintf("%T", niltype)
		ss := strings.Split(typeName, ".")
		sname := ss[len(ss)-1]
		if sname[0] < 'A' || sname[0] > 'Z' {
			log.Fatal("nuxui", "register type %s faild, this type is not exported", typeName)
		}
		if strings.Count(typeName, "*") != 1 {
			log.Fatal("nuxui", "register type %s faild, niltype argument should like `(*MyType)(nil)`", typeName)
		}
	}

	t := reflect.TypeOf(niltype).Elem()
	fullTypeName := t.PkgPath() + "." + t.Name()
	if _, ok := _creatorList[fullTypeName]; ok {
		log.Fatal("nuxui", "Type %s is already registed, do not register again", fullTypeName)
	} else {
		if creator == nil {
			log.Fatal("nuxui", "Type %s creator can not be nil", fullTypeName)
		}
		_creatorList[fullTypeName] = creator
	}
}

func FindTypeCreator(whichType any) Creator {
	if fullTypeName, ok := whichType.(string); ok {
		return findTypeCreatorByName(fullTypeName)
	}

	t := reflect.TypeOf(whichType).Elem()
	fullTypeName := t.PkgPath() + "." + t.Name()
	return findTypeCreatorByName(fullTypeName)
}

func findTypeCreatorByName(fullTypeName string) Creator {
	if c, ok := _creatorList[fullTypeName]; ok {
		return c
	}
	log.Fatal("nuxui", "Type '%s' can not find, make sure it was registed", fullTypeName)
	return nil
}
