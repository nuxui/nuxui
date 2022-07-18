// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package theme

const BootstrapLight = `
{
	import: {
		ui: nuxui.org/nuxui/ui,
	},

	style: {
		btn: {
			type: ui.Button,
			textSize: 14,
			padding: {left: 16px, top: 8px, right: 16px, bottom: 8px},
		},
		btn_default: {
			textColor: #262626,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{
						shape: rect,
						solid: #e0e0e0,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
					{state:"pressed|hovered", shape:{
						shape: rect,
						solid: #9e9e9e,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
					{state:"hovered", shape:{
						shape: rect,
						solid: #c9c6c6,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
				]
			}
		},
		btn_primary: {
			textColor: #ffffff,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{
						shape: rect,
						solid: #3f51b5,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
					{state:"pressed|hovered", shape:{
						shape: rect,
						solid: #2b397e,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
					{state:"hovered", shape:{
						shape: rect,
						solid: #334296,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
				]
			}
		},
		btn_secondary: {
			textColor: #ffffff,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{
						shape: rect,
						solid: #f50057,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
					{state:"hovered|pressed", shape:{
						shape: rect,
						solid: #ab003c,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
					{state:"hovered", shape:{
						shape: rect,
						solid: #dc0350,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
				],
			},
		},
		btn_disable: {
			textColor: #888888,
			disable: true,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{
						shape: rect,
						solid: #494949,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
				],
			},
		},
		btn_default_text: {
			textColor: #ffffff,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{}},
					{state:"pressed|hovered", shape:{
						shape: rect,
						solid: #797979,,
						stroke: #ffffff3a,
						strokeWidth: 1px,
						cornerRadius: 4px,
					}},
					{state:"hovered", shape:{
						shape: rect,
						stroke: #ffffff3a,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
				]
			}
		},
		btn_primary_text: {
			textColor: #6578db,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{}},
					{state:"pressed|hovered", shape:{
						shape: rect,
						solid: #3a4683,,
						stroke: #3f51b5,
						strokeWidth: 1px,
						cornerRadius: 4px,
					}},
					{state:"hovered", shape:{
						shape: rect,
						stroke: #3f51b5,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
				]
			}
		},
		btn_secondary_text: {
			textColor: #e30044,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{}},
					{state:"pressed|hovered", shape:{
						shape: rect,
						solid: #782b3f,,
						stroke: #952744,
						strokeWidth: 1px,
						cornerRadius: 4px,
					}},
					{state:"hovered", shape:{
						shape: rect,
						stroke: #952744,
						cornerRadius: 4px,
						shadow:{color: #00000088, x: 0, y: 1px, blur: 3px},
					}},
				]
			}
		},
		btn_disable_text: {
			textColor: #575757,
			disable: true,
		},
		btn_default_outline: {
			textColor: #ffffff,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{
						shape: rect,
						stroke: #ffffff3a,
						strokeWidth: 1px,
						cornerRadius: 4px,
					}},
					{state:"pressed|hovered", shape:{
						shape: rect,
						solid: #797979,,
						stroke: #ffffff3a,
						strokeWidth: 1px,
						cornerRadius: 4px,
					}},
					{state:"hovered", shape:{
						shape: rect,
						solid: #7979793a,,
						stroke: #ffffff3a,
						strokeWidth: 1px,
						cornerRadius: 4px,
					}},
				],
			},
		},
		btn_primary_outline: {
			textColor: #6578db,
			background: {
				type: ui.ShapeDrawable,
				states: [
					{state:"default", shape: {
						shape: rect,
						stroke: #4659c4,
						cornerRadius: 4px,
					}},
					{state:"pressed|hovered", shape: {
						shape: rect,
						solid: #3a4683,
						stroke: #4659c4,
						cornerRadius: 4px,
					}},
					{state:"hovered", shape: {
						shape: rect,
						solid: #4659c43a,
						stroke: #4659c4,
						cornerRadius: 4px,
					}},
				],
			},
		},
		btn_secondary_outline: {
			textColor: #f50357,
			background: {
				type: ui.ShapeDrawable,
				states:[
					{state:"default", shape:{
						shape: rect,
						stroke: #952744,
						cornerRadius: 4px,
					}},
					{state:"pressed|hovered", shape:{
						shape: rect,
						solid: #782b3f,
						stroke: #952744,
						cornerRadius: 4px,
					}},
					{state:"hovered", shape:{
						shape: rect,
						solid: #782b3f3a,
						stroke: #952744,
						cornerRadius: 4px,
					}},
				],
			},
		},
		btn_disable_outline:{
			disable: true,
			textColor: #575757,
			background: {
				type: ui.ShapeDrawable,
				shape:{
					shape: rect,
					stroke: #494949,
					cornerRadius: 4px,
				},
			},
		}
	},
}
`
