// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/nuxui/nuxui/log"
)

/*
width: auto 100px 100dp 50% 16:9 1wt unlimit
margin: !auto 10px 10dp 1wt 5% !ratio !unlimit
padding: !auto 10px 10dp 5% !wt !ratio !unlimit
*/

type Mode int32

const (
	// EXACTLY
	Pixel   Mode = 0 // 1111 or 0000 is px or 7 1px
	Auto    Mode = 1 // wrap content
	Unspec  Mode = 2
	Ratio   Mode = 3 // 16:9, no zero, no negative
	Percent Mode = 4 // 50% => match_parent * 50% => parent.size - parent.padding
	Weight  Mode = 5 // 1wt >= 0
	Default Mode = 6 // default nux unit
)

// TODO:: interface or struct
type Dimen interface {
	Mode() Mode
	Value() float32
	Equal(dimen Dimen) bool
	ModeName() string
	String() string
}

type dimen struct {
	mode  Mode
	value float32
}

func NewDimen(value float32, mode Mode) Dimen {
	d := &dimen{mode, value}
	if !d.valid() {
		log.Fatal("nux", "invalid dimension "+d.String())
	}
	d.fix()
	return d
}

func ParseDimen(s string) (Dimen, error) {
	d := &dimen{}
	if s == "" || strings.Compare(s, "auto") == 0 {
		d.mode = Auto
		d.value = 0
	} else if strings.HasSuffix(s, "%") {
		d.mode = Percent
		v, e := strconv.ParseFloat(string(s[0:len(s)-1]), 32)
		if e != nil {
			return d, fmt.Errorf(`invalid percent dimension format: "%s"`, s)
		}
		d.value = float32(v)
	} else if strings.Compare(s, "0") == 0 {
		d.mode = Pixel
		d.value = 0
	} else if strings.HasSuffix(s, "px") {
		d.mode = Pixel
		v, e := strconv.ParseFloat(string(s[0:len(s)-2]), 32)
		if e != nil {
			return d, fmt.Errorf(`invalid pixel dimension format: "%s"`, s)
		}
		if v > 0 { // math.round
			d.value = float32(int32(v + 0.5))
		} else {
			d.value = float32(int32(v - 0.5))
		}

	} else if strings.HasSuffix(s, "wt") {
		d.mode = Weight
		v, e := strconv.ParseFloat(string(s[0:len(s)-2]), 32)
		if e != nil || v <= 0 {
			return d, fmt.Errorf(`invalid weight dimension format: "%s"`, s)
		}
		d.value = float32(v)
	} else if strings.Contains(s, ":") {
		d.mode = Ratio
		v := strings.Split(s, ":")
		w, e1 := strconv.ParseFloat(v[0], 32)
		h, e2 := strconv.ParseFloat(v[1], 32)
		if len(v) != 2 || e1 != nil || e2 != nil || w <= 0 || h <= 0 {
			return d, fmt.Errorf(`invalid ratio dimension format: "%s"`, s)
		}
		d.value = float32(w / h)
	} else {
		d.mode = Default
		v, e := strconv.ParseFloat(s, 32)
		if e != nil {
			return d, fmt.Errorf(`invalid default dimension format: "%s"`, s)
		}
		d.value = float32(v)
	}
	if !d.valid() {
		log.Fatal("nux", "invalid dimension "+d.String())
	}
	d.fix()
	return d, nil
}

func (me *dimen) Mode() Mode {
	return me.mode
}

func (me *dimen) Value() float32 {
	return me.value
}

func (me *dimen) Equal(d Dimen) bool {
	return me.mode == d.Mode() && math.Abs(float64(me.value-d.Value())) <= math.SmallestNonzeroFloat32
}

func (me *dimen) ModeName() string {
	switch me.mode {
	case Pixel:
		return "Pixel"
	case Auto:
		return "Auto"
	case Percent:
		return "Percent"
	case Weight:
		return "Weight"
	case Ratio:
		return "Ratio"
	case Default:
		return "Default"
	}
	return "Unlimit"
}

func (me *dimen) String() string {
	if me == nil {
		return "nil"
	}

	switch me.mode {
	case Pixel:
		return fmt.Sprintf(`%dpx`, int32(me.value))
	case Auto:
		return "auto"
	case Percent:
		return fmt.Sprintf(`%.2f%%`, me.value)
	case Weight:
		return fmt.Sprintf(`%dwt`, int32(me.value))
	case Ratio:
		return fmt.Sprintf(`%.2f`, me.value)
	case Default:
		return "default"
	}
	return "unlimit"
}

func (me *dimen) fix() {
	if me.mode == Pixel {
		//TODO if me.value == 0 || math.IsNaN(me.value) || math.IsInf(me.value, 0) {
		// }
		if me.value >= 0 {
			me.value = float32(int32(me.value + 0.5))
		} else {
			me.value = float32(int32(me.value - 0.5))
		}
	}
}

func (me *dimen) valid() bool {
	switch me.mode {
	case Percent, Weight, Ratio, Unspec:
		return me.value >= 0
	default:
		return true
	}
}

// Non-negative and only support Pixel, Auto, Unlimit
type MeasuredDimen int32

const MaxMeasuredDimen int32 = 0x3FFFFFFF
const measuredDimenModeMask int32 = -1073741824 //0xC0000000
const measuredDimenValueMask int32 = 0x3FFFFFFF

func MeasureSpec(value int32, mode Mode) int32 {
	if mode < 0 || mode > 2 {
		log.Fatal("nux", "error for MeasureSpec mdoe, only support Pixel, Auto, Unspec")
	}
	if value > MaxMeasuredDimen {
		log.Fatal("nux", "error for MeasureSpec value, cannot be negative")
	} else if value < 0 {
		log.Fatal("nux", "error for MeasureSpec value, cannot be negative")
	}
	return int32(mode)<<30 | value
}

func MeasureSpecMode(size int32) Mode {
	return Mode((size & measuredDimenModeMask) >> 30 & 3)
}

func MeasureSpecValue(size int32) int32 {
	return size & measuredDimenValueMask
}

func MeasureSpecString(size int32) string {
	switch MeasureSpecMode(size) {
	case Pixel:
		return fmt.Sprintf(`%dpx`, MeasureSpecValue(size))
	case Auto:
		return fmt.Sprintf(`%dat`, MeasureSpecValue(size))
	}
	return "Unspec"
}
