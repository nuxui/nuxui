// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

/*
 var a dimen = 10 // 10px
 var a dimen = -10 // -10px
*/

type Dimen int32
type Mode byte

const (
	Pixel   Mode = iota // 1 means 1px
	Auto                // wrap content
	Ems                 // 1em = 1 x font-size
	Weight              // 1wt >= 0
	Ratio               // 16:9, no zero, no negative
	Percent             // 50% means (parent.size - parent.padding)*0.5
	spare               // default nux unit
	pixel               // negative pixel 111
	Unlimit = 2         // MeasuredDimen used Pixel, Auto, Unlimit
)

const (
	MaxDimen    Dimen  = 0x0FFFFFFF // 00001111111111111111111111111111
	MinDimen    Dimen  = -268435456 // 11110000000000000000000000000000
	dimenSMask  uint32 = 0x80000000 // 10000000000000000000000000000000
	dimenMMask  uint32 = 0x70000000 // 01110000000000000000000000000000
	dimenVMask  uint32 = 0x8FFFFFFF // 10001111111111111111111111111111
	dimenFMask  uint32 = 0x7FFFFFF8 // 01111111111111111111111111111000
	dimenFRMask uint32 = 0x0FFFFFFF // 00001111111111111111111111111111
)

func ADimen(value float32, mode Mode) Dimen {
	if mode < Pixel || mode > pixel {
		log.Fatal("nuxui", "Invalid dimen mode")
	}

	switch mode {
	case Weight:
		if value == 0 {
			mode = Pixel
		} else if value < 0 {
			log.Fatal("nuxui", "Invalid dimen value, the value of Weight can not be negative")
		}
	case Ratio:
		if value <= 0 {
			log.Fatal("nuxui", "Invalid dimen value, the value of Ratio must be positive")
		}
	case Percent:
		if value == 0 {
			mode = Pixel
		}
	}

	switch mode {
	case Ratio, Percent:
		v := *(*uint32)(unsafe.Pointer(&value))
		return Dimen(dimenSMask&v | uint32(mode)<<28 | (dimenFMask&v)>>3)
	default:
		var v int32
		if value >= 0 {
			v = int32(value + 0.5)
		} else {
			v = int32(value - 0.5)
		}
		return Dimen(dimenSMask&uint32(v) | uint32(mode)<<28 | (dimenVMask & uint32(v)))
	}
}

func (me Dimen) Mode() Mode {
	return Mode(dimenMMask & uint32(me) >> 28)
}

func (me Dimen) Value() float32 {
	switch me.Mode() {
	case Percent, Ratio:
		v := dimenSMask&uint32(me) | (dimenFRMask&uint32(me))<<3
		return *(*float32)(unsafe.Pointer(&v))
	default:
		if dimenSMask&uint32(me) == dimenSMask {
			return float32(int32(dimenSMask&uint32(me) | dimenMMask | (dimenVMask & uint32(me))))
		}
		return float32(int32(dimenSMask&uint32(me) | (dimenVMask & uint32(me))))
	}
}

func (me Mode) String() string {
	switch me {
	case Pixel, pixel:
		return "Pixel"
	case Auto:
		return "Auto"
	case Ems:
		return "Ems"
	case Weight:
		return "Weight"
	case Ratio:
		return "Ratio"
	case Percent:
		return "Percent"
	case spare:
		return "spare"
	}
	log.Fatal("nuxui", "can not run here.")
	return ""
}

func (me Dimen) String() string {
	switch me.Mode() {
	case Pixel, pixel:
		return fmt.Sprintf(`%dpx`, int32(me.Value()))
	case Auto:
		return "auto"
	case Weight:
		return fmt.Sprintf(`%dwt`, int32(me.Value()))
	case Ems:
		return fmt.Sprintf(`%dem`, int32(me.Value()))
	case Ratio:
		return fmt.Sprintf(`%.2f`, me.Value())
	case Percent:
		return fmt.Sprintf(`%.2f%%`, me.Value())
	case spare:
		return "spare"
	}
	log.Fatal("nuxui", "can not run here.")
	return ""
}

func SDimen(s string) Dimen {
	d, e := ParseDimen(s)
	if e != nil {
		log.Fatal("nuxui", e.Error())
	}

	return d
}
func ParseDimen(s string) (Dimen, error) {
	if s == "" || s == "auto" {
		return ADimen(0, Auto), nil
	} else if s == "unlimit" {
		return ADimen(0, Unlimit), nil
	} else if strings.HasSuffix(s, "%") {
		v, e := strconv.ParseFloat(string(s[0:strlen(s)-1]), 32)
		if e != nil {
			return 0, fmt.Errorf(`invalid Percent dimension format: "%s"`, s)
		}
		return ADimen(float32(v), Percent), nil
	} else if s == "0" {
		return 0, nil
	} else if strings.HasSuffix(s, "px") {
		v, e := strconv.ParseFloat(string(s[0:strlen(s)-2]), 32)
		if e != nil {
			return 0, fmt.Errorf(`invalid Pixel dimension format: "%s"`, s)
		}
		return ADimen(float32(v), Pixel), nil
	} else if strings.HasSuffix(s, "em") {
		v, e := strconv.ParseFloat(string(s[0:strlen(s)-2]), 32)
		if e != nil {
			return 0, fmt.Errorf(`invalid Ems dimension format: "%s"`, s)
		}
		return ADimen(float32(v), Ems), nil
	} else if strings.HasSuffix(s, "wt") {
		v, e := strconv.ParseFloat(string(s[0:strlen(s)-2]), 32)
		if e != nil {
			return 0, fmt.Errorf(`invalid Weight dimension format: "%s"`, s)
		}
		if v < 0 {
			return 0, fmt.Errorf(`invalid Weight dimension format: "%s", the weight must be positive number`, s)
		}
		return ADimen(float32(v), Weight), nil
	} else if strings.Contains(s, ":") {
		v := strings.Split(s, ":")
		w, e1 := strconv.ParseFloat(v[0], 32)
		h, e2 := strconv.ParseFloat(v[1], 32)
		if e1 != nil || e2 != nil || len(v) != 2 {
			return 0, fmt.Errorf(`invalid Ratio dimension format: "%s"`, s)
		}
		if w <= 0 || h <= 0 {
			return 0, fmt.Errorf(`invalid Ratio dimension format: "%s", the ratio must be positive`, s)
		}
		return ADimen(float32(w/h), Ratio), nil
	}
	return 0, fmt.Errorf(`invalid dimension format: "%s"`, s)
}

//------------------  MeasuredDimen ------------------------------------

// Non-negative and only support Pixel, Auto, Unlimit
type MeasuredDimen int32

const MaxMeasuredDimen int32 = 0x3FFFFFFF
const measuredDimenModeMask int32 = -1073741824 //0xC0000000
const measuredDimenValueMask int32 = 0x3FFFFFFF

func MeasureSpec(value int32, mode Mode) int32 {
	if !(mode == Pixel || mode == Auto || mode == Unlimit) {
		log.Fatal("nuxui", "error for MeasureSpec mdoe, only support Pixel, Auto, Unlimit")
	}
	if value > MaxMeasuredDimen {
		log.Fatal("nuxui", "error for MeasureSpec value, out of range")
	} else if value < 0 {
		// log.Fatal("nuxui", "error for MeasureSpec value, cannot be negative")
		value = 0
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
	return fmt.Sprintf(`%d %s`, MeasureSpecValue(size), MeasureSpecMode(size))
}
