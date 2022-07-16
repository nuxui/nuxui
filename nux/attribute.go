// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"nuxui.org/nuxui/log"
)

// TODO:: GetDimenArray GetIntArray...

type Attr map[string]any

// merge attrs to a new single attr, the back will override front
func MergeAttrs(attrs ...Attr) (attr Attr) {
	l := len(attrs)
	if l == 0 {
		return Attr{}
	}

	attr = Attr{}
	for _, item := range attrs {
		for k, v := range item {
			attr[k] = v
		}
	}
	return attr
}

func (me Attr) Has(key string) bool {
	_, ok := me[key]
	return ok
}

func (me Attr) Set(key string, value any) {
	me[key] = value
}

func (me Attr) Get(key string, defaultValue any) any {
	if attr, ok := me[key]; ok {
		return attr
	}
	return defaultValue
}

func (me Attr) GetAttr(key string, defaultValue Attr) Attr {
	ret := defaultValue
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case Attr:
			ret = t
		case map[string]any:
			ret = *(*Attr)(unsafe.Pointer(&t))
		default:
			log.E("nuxui", "unsupport convert %T to Attr, use default value instead", t)

		}
	}
	return ret
}

func (me Attr) GetString(key string, defaultValue string) string {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			return t
		case rune, byte, uint, uint16, uint32, uint64, int, int8, int16, int64, float32, float64, bool: //byte = uint8, rune = int32
			return fmt.Sprintf("%s", t)
		default:

			log.E("nuxui", "unsupport convert %s %T:%s to string, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetBool(key string, defaultValue bool) bool {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			b, e := strconv.ParseBool(t)
			if e != nil {
				log.E("nuxui", "unsupport convert %s %T:%s to bool, use default value instead", key, t, t)
				return defaultValue
			}
			return b
		case bool:
			return t
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to bool, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetColor(key string, defaultValue Color) (result Color) {
	result = defaultValue
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case Color:
			return t
		case string:
			t = strings.TrimSpace(t)
			if !strings.HasPrefix(t, "#") {
				log.E("nuxui", "color value of %s should start with #, use default value instead", key)
				result = defaultValue
			} else {
				c := strings.TrimPrefix(t, "#")
				if strlen(c) == 6 {
					c += "FF"
				}
				if strlen(c) != 8 {
					log.E("nuxui", "cannot convert %s: '%s' to Color, use default value instead", key, t)
					result = defaultValue
				}
				v, e := strconv.ParseUint(c, 16, 32)
				if e != nil {
					log.E("nuxui", "cannot convert %s: '%s' to Color, use default value instead", key, t)
					result = defaultValue
				} else {
					result = Color(v)
				}
			}
		case int:
			result = Color(t)
		case int64:
			result = Color(t)
		case int32:
			result = Color(t)
		case uint64:
			result = Color(t)
		case uint32:
			result = Color(t)
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to Color, use default value instead", key, t, t)
			result = defaultValue
		}
	}

	return result
}

func (me Attr) GetUint32(key string, defaultValue uint32) (result uint32) {
	result = defaultValue
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseUint(t, 10, 32)
			if e != nil {
				log.E("nuxui", "cannot convert %s %T:%s to uint32, use default value instead", key, t, t)
				result = defaultValue
			} else {
				result = uint32(v)
			}
		case int:
			result = uint32(t)
		case int64:
			result = uint32(t)
		case int32:
			result = uint32(t)
		case uint64:
			result = uint32(t)
		case uint32:
			result = uint32(t)
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to uint32, use default value instead", key, t, t)
			result = defaultValue
		}
	}

	return result
}

func (me Attr) GetInt32(key string, defaultValue int32) int32 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseInt(t, 10, 32)
			if e != nil {
				log.E("nuxui", "cannot convert %s %T:%s to int32, use default value instead", key, t, t)
				return defaultValue
			}
			return int32(v)
		case int:
			return int32(t)
		case int64:
			return int32(t)
		case int32:
			return int32(t)
		case uint64:
			return int32(t)
		case uint32:
			return int32(t)
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to int32, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetInt(key string, defaultValue int) int {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseInt(t, 10, 32)
			if e != nil {
				log.E("nuxui", "cannot convert %s %T:%s to int32, use default value instead", key, t, t)
				return defaultValue
			}
			return int(v)
		case int:
			return t
		case int64:
			return int(t)
		case int32:
			return int(t)
		case uint64:
			return int(t)
		case uint32:
			return int(t)
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to int32, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetInt64(key string, defaultValue int64) int64 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseInt(t, 10, 64)
			if e != nil {
				log.E("nuxui", "cannot convert %s %T:%s to int64, use default value instead", key, t, t)
				return defaultValue
			}

			return int64(v)
		case int:
			return int64(t)
		case int64:
			return int64(t)
		case int32:
			return int64(t)
		case uint64:
			return int64(t)
		case uint32:
			return int64(t)
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to int64, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetFloat32(key string, defaultValue float32) float32 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 32)
			if e != nil {
				log.E("nuxui", "cannot convert %s %T:%s to float32, use default value instead", key, t, t)
				return defaultValue
			}

			return float32(v)
		case float64:
			return float32(t)
		case float32:
			return float32(t)
		case int:
			return float32(t)
		case int64:
			return float32(t)
		case int32:
			return float32(t)
		case uint64:
			return float32(t)
		case uint32:
			return float32(t)
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to float32, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetFloat64(key string, defaultValue float64) float64 {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			v, e := strconv.ParseFloat(t, 64)
			if e != nil {
				log.E("nuxui", "cannot convert %s %T:%s to float64, use default value instead", key, t, t)
				return defaultValue
			}

			return float64(v)
		case float64:
			return float64(t)
		case float32:
			return float64(t)
		case int:
			return float64(t)
		case int64:
			return float64(t)
		case int32:
			return float64(t)
		case uint64:
			return float64(t)
		case uint32:
			return float64(t)
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to float64, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetArray(key string, defaultValue []any) []any {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case []any:
			return t
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to array, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetAttrArray(key string, defaultValue []Attr) []Attr {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case []Attr:
			return t
		case []any:
			ret := make([]Attr, len(t))
			for i, a := range t {
				if atr, ok := a.(Attr); ok {
					ret[i] = atr
				} else {
					log.E("nuxui", "unsupport convert %s %T:%s to []Attr, use default value instead", key, t, t)
					return defaultValue
				}
			}
			return ret
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to Attr array, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) GetStringArray(key string, defaultValue []string) []string {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case []string:
			return t
		case []any:
			ret := make([]string, len(t))
			for i, a := range t {
				if str, ok := a.(string); ok {
					ret[i] = str
				} else {
					log.E("nuxui", "unsupport convert %s %T:%s to []string, use default value instead", key, t, t)
					return defaultValue
				}
			}
			return ret
		default:
			log.E("nuxui", "unsupport convert %s %T:%s to []string, use default value instead", key, t, t)
		}
	}
	return defaultValue
}

func (me Attr) Merge(attr ...Attr) {
	for _, a := range attr {
		for k, v := range a {
			me[k] = v
		}
	}
}

func (me Attr) String() string {
	return attrToString(me, 0)
}

func attrToString(attr Attr, depth int) string {
	var ret string
	var space string

	for i := 0; i != depth; i++ {
		space += "  "
	}

	ret += space + "{\n"

	for k, v := range attr {
		if arr, ok := v.([]any); ok {
			ret += fmt.Sprintf("%s  %s: [\n", space, k)
			for _, a := range arr {
				if child, ok := a.(Attr); ok {
					ret += attrToString(child, depth+1)
				} else {
					ret += fmt.Sprintf("%s    %s\n", space, a)
				}
			}
			ret += space + "  ]\n"
		} else if atr, ok := v.(Attr); ok {
			ret += fmt.Sprintf("%s  %s:\n", space, k)
			ret += attrToString(atr, depth+1)
		} else {
			ret += fmt.Sprintf("%s  %s: %s\n", space, k, v)
		}
	}
	ret += space + "}\n"

	return ret
}

func (me Attr) GetDimen(key string, defaultValue string) Dimen {
	if attr, ok := me[key]; ok {
		switch t := attr.(type) {
		case string:
			ret, e := ParseDimen(t)
			if e != nil {
				log.E("nuxui", "can't convert %s to Dimen, use default value instead, %s", t, e.Error())
				ret, e = ParseDimen(defaultValue)
				if e != nil {
					log.E("nuxui", "can't convert default value %s to Dimen, %s", defaultValue, e.Error())
					goto end
				}
				return ret
			}
			return ret
		case int:
			return Dimen(t)
		default:
			log.E("nuxui", "unsupport convert %T to Dimen, use default value instead", t)
		}
	} else {
		ret, e := ParseDimen(defaultValue)
		if e != nil {
			log.E("nuxui", "can't convert default value %s to Dimen, %s", defaultValue, e.Error())
			goto end
		}
		return ret
	}

end:
	ret := ADimen(0, Auto)
	return ret
}

///////////////////////////////////// Parse /////////////////////////////////////

// ParseAttr parse attr string to Attr
// print error if has error and return empty Attr
func ParseAttr(attrtext string) Attr {
	ret, err := ParseAttrWitthError(attrtext)
	if err != nil {
		if debug_attr {
			log.Fatal("nuxui", err.Error())
		} else {
			log.E("nuxui", err.Error())
		}
		return Attr{}
	}
	return ret
}

func ParseAttrWitthError(attrtext string) (Attr, error) {
	return (&attr{}).parse(attrtext)
}

const (
	_EOF    = -1
	_STRING = iota
	_ARRAY
	_STRUCT
)

type attr struct {
	data []rune
	len  int
	pos  int
	line int
	r    rune
	i    Attr // import
}

func (me *attr) parse(template string) (Attr, error) {
	me.data = []rune(template)
	me.len = len(me.data)
	me.pos = -1
	me.line = 1
	me.r = _EOF

	me.nextRune()
	me.skipBlank()
	if me.r != '{' {
		return nil, fmt.Errorf("line:%d first element of Attr is not '{'", me.line)
	}
	return me.nextStruct()
}

func (me *attr) nextStruct() (Attr, error) {
	data := Attr{}
	var key string
	var value any
	var err error

	for {
		me.nextRune()
		me.skipBlank()
		if me.r == _EOF {
			return nil, fmt.Errorf("line:%d unclosed attribute, not end with '}'", me.line)
		}
		if me.r == '}' {
			break
		}
		if me.r == ',' {
			continue
		}

		key, err = me.nextKey()
		if err != nil {
			return nil, err
		}
		me.nextRune()
		me.skipBlank()
		value, err = me.nextValue(_STRUCT)
		if err != nil {
			return nil, err
		}
		if debug_attr {
			if _, ok := data[key]; ok {
				log.E("nuxui", "line:%d attribute key %s is already exist", me.line, key)
			}
		}
		data[key] = value

		if me.i == nil && key == "import" {
			if i, ok := value.(Attr); ok {
				me.i = i
			}
		}
	}

	return data, nil
}

func (me *attr) nextValue(parent int) (any, error) {
	switch me.r {
	case '"', '`':
		return me.nextString(me.r)
	case '{':
		return me.nextStruct()
	case '[':
		return me.nextArray()
	default:
		str, err := me.nextString(',')
		if err != nil {
			return nil, err
		}

		if parent == _ARRAY && strings.Contains(str, " ") {
			return nil, fmt.Errorf("line:%d strings which has space in the array should be quoted", me.line)
		}

		if me.i != nil {
			strs := strings.Split(str, ".")
			if len(strs) == 2 {
				if v, ok := me.i[strs[0]]; ok {
					return fmt.Sprintf("%s.%s", v, strs[1]), nil
				}
			}
		}
		return str, nil
	}
	return nil, nil
}

func (me *attr) nextArray() ([]any, error) {
	arr := make([]any, 0)
	var value any
	var err error
	for {
		me.nextRune()
		me.skipBlank()
		if me.r == ']' {
			break
		}
		if me.r == ',' {
			continue
		}
		value, err = me.nextValue(_ARRAY)
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
	}
	return arr, nil
}

func (me *attr) skipBlank() {
	for {
		switch me.r {
		case ' ', '\r', '\n', '\t', '\f', '\b', 65279:
			if me.r == '\n' {
				me.line++
			}
			me.nextRune()
			continue
		case '/':
			me.skipComment()
			continue
		}
		break
	}
}

func (me *attr) nextRune() rune {
	me.pos++
	if me.pos >= me.len {
		me.r = _EOF
	} else {
		me.r = me.data[me.pos]
	}
	return me.r
}

func (me *attr) previousRune() rune {
	me.pos--
	if me.pos >= me.len || me.pos < 0 {
		me.r = _EOF
	} else {
		me.r = me.data[me.pos]
	}
	return me.r
}

func (me *attr) skipComment() error {
	me.nextRune()
	if me.r == '/' {
		for {
			me.nextRune()
			if me.r == '\n' {
				me.line++
				me.nextRune()
				return nil
			} else if me.r == _EOF {
				return nil
			}
		}
	} else if me.r == '*' {
		me.nextRune()
		for me.r != _EOF {
			if me.r == '*' {
				me.nextRune()
				if me.r == '/' {
					me.nextRune()
					return nil
				}
				continue
			}
			me.nextRune()
		}
	} else {
		return fmt.Errorf("line:%d invalid comment", me.line)
	}
	return nil
}

// obtain next key, and last rune is ':'
func (me *attr) nextKey() (string, error) {
	p := me.pos
	if me.r < 'A' || (me.r > 'Z' && me.r < 'a') || me.r > 'z' {
		return "", fmt.Errorf("line:%d invalid key name, the first letter must be [A-Za-z]", me.line)
	}
	me.nextRune()
	for (me.r >= 'a' && me.r <= 'z') || (me.r >= 'A' && me.r <= 'Z') || (me.r >= '0' && me.r <= '9') || me.r == '_' {
		me.nextRune()
	}
	key := string(me.data[p:me.pos])
	me.skipBlank()
	if me.r != ':' {
		return "", fmt.Errorf(`line:%d invalid format, missing ':' after key name "%s", `, me.line, key)
	}
	return key, nil
}

func (me *attr) nextString(end rune) (string, error) {
	ret := []rune{}
	quot := (end == '"' || end == '`')

	if end == ',' && me.r == end {
		return "", nil
	}

	if !quot {
		me.previousRune()
	}
	p := me.pos + 1

	for {
		me.nextRune()
		if me.r == '\\' {
			if me.pos+1 < me.len {
				ret = append(ret, me.data[p:me.pos]...)
				me.nextRune()
				switch me.r {
				case 'a':
					ret = append(ret, '\a')
				case 'b':
					ret = append(ret, '\b')
				case '\\':
					ret = append(ret, '\\')
				case 't':
					ret = append(ret, '\t')
				case 'n':
					ret = append(ret, '\n')
				case 'f':
					ret = append(ret, '\f')
				case 'r':
					ret = append(ret, '\r')
				case 'v':
					ret = append(ret, '\v')
				case '\'':
					ret = append(ret, '\'')
				case '"':
					ret = append(ret, '"')
				case 'x', 'u', 'U':
					n := 2
					switch me.r {
					case 'x':
						n = 2
					case 'u':
						n = 4
					case 'U':
						n = 8
					}
					if me.pos+n+1 < me.len {
						v, err := strconv.ParseUint(string(me.data[me.pos+1:me.pos+n+1]), 16, 64)
						if err != nil {
							return "", fmt.Errorf("line:%d parsing '\\%s' error, invalid syntax", me.line, string(me.data[me.pos:me.pos+n+1]))
						}
						ret = append(ret, rune(v))
						me.pos += n
					} else {
						return "", fmt.Errorf("line:%d parsing '\\%s' error, invalid syntax", me.line, string(me.data[me.pos:me.len-1]))
					}
				}
				p = me.pos + 1
			} else { // _EOF
				if quot {
					return "", fmt.Errorf("line:%d unclosed string", me.line)
				}
				return "", nil
			}
		} else {
			if quot {
				switch me.r {
				case '"', '`':
					if me.r == end && me.data[me.pos-1] != '\\' {
						return string(append(ret, me.data[p:me.pos]...)), nil
					}
				case _EOF:
					return "", fmt.Errorf("line:%d unclosed string", me.line)
				}
			} else {
				switch me.r {
				case '\n':
					return "", fmt.Errorf("line:%d missing comma at line end or quot whole string", me.line)
				case ',', '}', ']':
					ret := strings.TrimSpace(string(append(ret, me.data[p:me.pos]...)))
					me.previousRune()
					return ret, nil
				case _EOF:
					return "", fmt.Errorf("line:%d invalid value, syntax error", me.line)
				}
			}
		}

	}

	return "", nil
}

func strlen(text string) int {
	return len([]rune(text))
}
