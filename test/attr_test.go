// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"fmt"
	"testing"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

var template = `
{
  import: {
    ui: github.com/nuxui/nuxui/ui,
  },

  layout: {
	id: "root",
	widget: ui.Column,
	width: 1wt,
	height: 1wt,
	background: #215896,
	padding: {left: 10px, top: 10px, right: 10px, bottom: 10px},
	children:[
	{
		id: "edit",
		widget: ui.Editor,
		width: 1wt,
		height: 30px,
		background: #982368,
		text: "nuxui.org example",
		font: {family: "Menlo, Monaco, Courier New, monospace", size: 14, color: #ffffff }
	},{
		id: "header",
		widget: ui.Column,
		width: 1wt,
		height: 100px,
		background: #123098,
		padding: {left: 10px, top: 10px, right: 10px, bottom: 10px},
		children:[
		{
			id: "xxx",
			widget: ui.Text,
			width: auto,
			height: auto,
			background: #982368,
			text: "{{me.name}}",
		}
		]
	}
	]
  }
}
  
  `

func TestAttr(t *testing.T) {
	defer log.Close()
	// attr := nux.ParseAttr(template)
	// log.V("test", "%s", attr)

	a := nux.Attr{
		"age":  1,
		"name": "nihao",
	}
	fmt.Println(a.GetString("name", "1"))
	printMapTemplate(a, 0)
}

func printMapTemplate(data nux.Attr, depth int) {
	s := ""
	for i := 0; i != depth; i++ {
		s += "  "
	}

	s2 := s + "  "
	fmt.Println(s + "{")
	for k, v := range data {
		switch t := v.(type) {
		case nux.Attr:
			fmt.Printf(s2+"%s: ", k)
			printMapTemplate(t, depth+1)
		default:
			fmt.Printf(s2+"%s: %s,\n", k, t)
		}
	}
	fmt.Println(s + "}")
}
