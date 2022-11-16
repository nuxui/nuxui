// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !wasm

package nux

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"nuxui.org/nuxui/log"
)

func ReadImageSVG(svgstr string) Image {
	return LoadImageSVGFromReader(strings.NewReader(svgstr))
}

func LoadImageSVGFromFile(fileName string) Image {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("nuxui", "%s", err.Error())
	}

	return LoadImageSVGFromReader(bytes.NewReader(data))
}

func LoadImageSVGFromReader(reader io.Reader) Image {
	svg := &imageSvg{}

	decoder := xml.NewDecoder(reader)
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("nuxui", "%s", err.Error())
		}
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "svg":
				svg.readSVGAttr(se.Attr)
			case "path":
				svg.readPathAttr(se.Attr)
			}
		case xml.EndElement:
		case xml.CharData:
		}
	}

	return svg
}

type imageSvg struct {
	viewBox Rect
	width   int32
	height  int32
	paths   []*svgpath
}

type svgpath struct {
	path  Path
	paint Paint
}

type pathparse struct {
	d      string
	len    int
	pos    int
	start  int
	c      uint8 // commond
	points [8]float32
}

func (me *imageSvg) readSVGAttr(attrs []xml.Attr) {
	for _, a := range attrs {
		switch a.Name.Local {
		case "viewBox":
			strs := strings.Split(strings.TrimSpace(a.Value), " ")
			if len(strs) == 4 {
				if v, e := strconv.ParseFloat(strs[0], 32); e == nil {
					me.viewBox.X = int32(v)
				}
				if v, e := strconv.ParseFloat(strs[1], 32); e == nil {
					me.viewBox.Y = int32(v)
				}
				if v, e := strconv.ParseFloat(strs[2], 32); e == nil {
					me.viewBox.Width = int32(v) - me.viewBox.X
				}
				if v, e := strconv.ParseFloat(strs[3], 32); e == nil {
					me.viewBox.Height = int32(v) - me.viewBox.Y
				}
			}
		case "width":
			if v, e := strconv.ParseFloat(a.Value, 32); e == nil {
				me.width = int32(v)
			}
		case "height":
			if v, e := strconv.ParseFloat(a.Value, 32); e == nil {
				me.height = int32(v)
			}
		}
	}
}

func (me *imageSvg) readPathAttr(attrs []xml.Attr) {
	p := &svgpath{
		paint: NewPaint(),
	}
	p.paint.SetColor(Black)

	for _, a := range attrs {
		switch a.Name.Local {
		case "d":
			p.path, _ = me.parsePath(a.Value)
		case "fill":
			if c, e := ParseColor(a.Value, Black); e == nil {
				p.paint.SetStyle(p.paint.Style() | PaintStyle_Fill)
				p.paint.SetColor(c)
			}
		case "stroke":
			if c, e := ParseColor(a.Value, Black); e == nil {
				p.paint.SetStyle(p.paint.Style() | PaintStyle_Stroke)
				p.paint.SetColor(c)
			}
		case "stroke-width":
			if width, e := strconv.ParseFloat(a.Value, 32); e == nil {
				p.paint.SetWidth(float32(width))
			} else {
				p.paint.SetWidth(1)
			}
		}
	}
	me.paths = append(me.paths, p)
}

func (me *imageSvg) parsePath(str string) (Path, error) {
	p := NewPath()
	var curX float32
	var curY float32
	var preCtrX float32
	var preCtrY float32
	var prevCurveX float32
	var prevCurveY float32
	var preCmd uint8 = 'm'
	var cmd uint8 = 'm'

	pp := &pathparse{
		d:      str,
		pos:    0,
		len:    len(str),
		c:      0,
		points: [8]float32{},
	}

	for pp.pos < pp.len {
		pp.skipBlank()
		if pp.pos >= pp.len {
			break
		}

		preCmd = cmd

		if pp.isalphabet(pp.d[pp.pos]) {
			cmd = pp.d[pp.pos]
		}

		switch cmd {
		case 'M':
			ps := pp.getPoints(2)
			p.MoveTo(ps[0], ps[1])
			curX = ps[0]
			curY = ps[1]
		case 'm':
			ps := pp.getPoints(2)
			p.MoveTo(curX+ps[0], curY+ps[1])
			curX += ps[0]
			curY += ps[1]
		case 'C':
			ps := pp.getPoints(6)
			p.CurveTo(ps[0], ps[1], ps[2], ps[3], ps[4], ps[5])
			preCtrX = ps[2]
			preCtrY = ps[3]
			curX = ps[4]
			curY = ps[5]
		case 'c':
			ps := pp.getPoints(6)
			p.CurveTo(curX+ps[0], curY+ps[1], curX+ps[2], curY+ps[3], curX+ps[4], curY+ps[5])
			preCtrX = curX + ps[2]
			preCtrY = curY + ps[3]
			curX += ps[4]
			curY += ps[5]
		case 'L':
			ps := pp.getPoints(2)
			p.LineTo(ps[0], ps[1])
			curX = ps[0]
			curY = ps[1]
			preCtrX = curX
			preCtrY = curY
		case 'l':
			ps := pp.getPoints(2)
			p.LineTo(curX+ps[0], curY+ps[1])
			curX += ps[0]
			curY += ps[1]
			preCtrX = curX
			preCtrY = curY
		case 'V':
			ps := pp.getPoints(1)
			p.LineTo(curX, ps[0])
			curY = ps[0]
			preCtrX = curX
			preCtrY = curY
		case 'v':
			ps := pp.getPoints(1)
			p.LineTo(curX, curY+ps[0])
			curY += ps[0]
			preCtrX = curX
			preCtrY = curY
		case 'H':
			ps := pp.getPoints(1)
			p.LineTo(ps[0], curY)
			curX = ps[0]
			preCtrX = curX
			preCtrY = curY
		case 'h':
			ps := pp.getPoints(1)
			p.LineTo(curX+ps[0], curY)
			curX += ps[0]
			preCtrX = curX
			preCtrY = curY
		case 'Q':
			ps := pp.getPoints(4)
			p.CurveToV(ps[0], ps[1], ps[2], ps[3])
			prevCurveX = ps[0]
			prevCurveY = ps[1]
			curX = ps[2]
			curY = ps[3]
		case 'q':
			ps := pp.getPoints(4)
			p.CurveToV(curX+ps[0], curY+ps[1], curX+ps[2], curY+ps[3])
			prevCurveX = curX + ps[0]
			prevCurveY = curY + ps[1]
			curX += ps[2]
			curY += ps[3]
		case 'T':
			ps := pp.getPoints(2)
			if preCmd != 'T' && preCmd != 't' && preCmd != 'Q' && preCmd != 'q' {
				prevCurveX = curX
				prevCurveY = curY
			} else {
				prevCurveX = curX + curX - prevCurveX
				prevCurveY = curY + curY - prevCurveY
			}
			p.CurveToV(prevCurveX, prevCurveY, ps[0], ps[1])
			curX = ps[0]
			curY = ps[1]
		case 't':
			ps := pp.getPoints(2)
			if preCmd != 'T' && preCmd != 't' && preCmd != 'Q' && preCmd != 'q' {
				prevCurveX = curX
				prevCurveY = curY
			} else {
				prevCurveX = curX + curX - prevCurveX
				prevCurveY = curY + curY - prevCurveY
			}
			p.CurveToV(prevCurveX, prevCurveY, curX+ps[0], curY+ps[1])
			curX += ps[0]
			curY += ps[1]
		case 'S':
			preCtrX = curX + curX - preCtrX
			preCtrY = curY + curY - preCtrY
			ps := pp.getPoints(4)
			p.CurveTo(preCtrX, preCtrY, ps[0], ps[1], ps[2], ps[3])
			preCtrX = ps[0]
			preCtrY = ps[1]
			curX = ps[2]
			curY = ps[3]
		case 's':
			preCtrX = curX + curX - preCtrX
			preCtrY = curY + curY - preCtrY
			ps := pp.getPoints(4)
			p.CurveTo(preCtrX, preCtrY, curX+ps[0], curY+ps[1], curX+ps[2], curY+ps[3])
			preCtrX = curX + ps[0]
			preCtrY = curY + ps[1]
			curX = curX + ps[2]
			curY = curY + ps[3]
		case 'a', 'A':
			ps := pp.getPoints(7)
			flagLarge := ps[3] > 0.001
			flagSweep := ps[4] > 0.001

			if pp.c == 'A' {
				me.arcToCurve(p, curX, curY, ps[0], ps[1], ps[2], flagLarge, flagSweep, ps[5], ps[6], &preCtrX, &preCtrY)
				curX = ps[5]
				curY = ps[6]
			} else {
				me.arcToCurve(p, curX, curY, ps[0], ps[1], ps[2], flagLarge, flagSweep, curX+ps[5], curY+ps[6], &preCtrX, &preCtrY)
				curX += ps[5]
				curY += ps[6]
			}
		case 'z', 'Z':
			p.Close()
			pp.pos++
		}
	}

	return p, nil
}

func (me *imageSvg) addArcToCurve(path Path, startX, startY, radiusX, radiusY, angle, endX, endY, startAngle, endAngle float32, endControlX, endControlY *float32) {
	// trigonometry
	t := float32(math.Tan(float64((endAngle - startAngle) / 4)))
	hx := radiusX * t * 4 / 3
	hy := radiusY * t * 4 / 3

	// calculate control points
	startCPX := startX + hx*float32(math.Sin(float64(startAngle)))
	startCPY := startY - hy*float32(math.Cos(float64(startAngle)))

	startCPX = 2*startX - startCPX
	startCPY = 2*startY - startCPY

	endCPX := endX + hx*float32(math.Sin(float64(endAngle)))
	endCPY := endY - hy*float32(math.Cos(float64(endAngle)))

	// !!! don't forget to rotate points back to the original reference frame

	// add curve
	path.CurveTo(startCPX, startCPY, endCPX, endCPY, endX, endY)

	*endControlX = endCPX
	*endControlY = endCPY
}

func (me *imageSvg) arcToCurve2(path Path, startX, startY, radiusX, radiusY, angle float32, sweep bool,
	endX, endY, startAngle, endAngle, centerX, centerY float32, endControlX,
	endControlY *float32) {
	angleDiff := endAngle - startAngle
	if math.Abs(float64(angleDiff)) > (math.Pi * 120.0 / 180.0) {
		// too wide so we need to break it up into multiple curves
		endAngleNext := endAngle
		endXNext := endX
		endYNext := endY

		if sweep && (endAngle > startAngle) {
			endAngle = startAngle + (math.Pi * 120.0 / 180.0)
		} else {
			endAngle = startAngle + (math.Pi*120.0/180.0)*-1
		}
		endX = centerX + radiusX*float32(math.Cos(float64(endAngle)))
		endY = centerY + radiusY*float32(math.Sin(float64(endAngle)))

		me.addArcToCurve(
			path, startX, startY, radiusX, radiusY, angle, endX, endY, startAngle, endAngle, endControlX, endControlY)
		me.arcToCurve2(path, endX, endY, radiusX, radiusY, angle, sweep, endXNext, endYNext, endAngle, endAngleNext,
			centerX, centerY, endControlX, endControlY)
	} else {
		me.addArcToCurve(
			path, startX, startY, radiusX, radiusY, angle, endX, endY, startAngle, endAngle, endControlX, endControlY)
	}
}

func (me *imageSvg) arcToCurve(path Path, startX, startY, radiusX, radiusY, angle float32, large, sweep bool, endX, endY float32, endControlX, endControlY *float32) {
	if radiusX == 0 || radiusY == 0 {
		path.LineTo(endX, endY)
		return
	}

	angle = math.Pi / 180.0 * angle
	hx := (startX - endX) / 2.0
	hy := (startY - endY) / 2.0

	validateRadii := (hx*hx)/(radiusX*radiusX) + (hy*hy)/(radiusY*radiusY)
	if validateRadii > 1 {
		radiusX = radiusX * float32(math.Sqrt(float64(validateRadii)))
		radiusY = radiusY * float32(math.Sqrt(float64(validateRadii)))
	}

	radiusX2 := (radiusX * radiusX)
	radiusY2 := (radiusY * radiusY)
	HX2 := (hx * hx)
	HY2 := (hy * hy)

	k := (radiusX2*radiusY2 - radiusX2*HY2 - radiusY2*HX2) / (radiusX2*HY2 + radiusY2*HX2)
	if large == sweep {
		k = float32(math.Sqrt(math.Abs(float64(k))) * -1)
	} else {
		k = float32(math.Sqrt(math.Abs(float64(k))))
	}

	centerX := k * (radiusX * hy / radiusY)
	centerY := k * (-radiusY * hx / radiusX)

	// F.6.5.3 - center of ellipse
	centerX += (startX + endX) / 2
	centerY += (startY + endY) / 2

	// calculate angles
	// F.6.5.4 to F.6.5.6
	aS := (startY - centerY) / radiusY
	aE := (endY - centerY) / radiusY

	// preventing out of range errors with asin due to floating point errors
	aS = float32(math.Min(float64(aS), 1.0))
	aS = float32(math.Max(float64(aS), -1.0))

	aE = float32(math.Min(float64(aE), 1.0))
	aE = float32(math.Max(float64(aE), -1.0))

	// get the angle
	startAngle := float32(math.Asin(float64(aS)))
	endAngle := float32(math.Asin(float64(aE)))

	if startX < centerX {
		startAngle = math.Pi - startAngle
	}

	if endX < centerX {
		endAngle = math.Pi - endAngle
	}

	if startAngle < 0 {
		startAngle = math.Pi*2 + startAngle
	}

	if endAngle < 0 {
		endAngle = math.Pi*2 + endAngle
	}

	if sweep && startAngle > endAngle {
		startAngle = startAngle - math.Pi*2.0
	}

	if !sweep && endAngle > startAngle {
		endAngle = endAngle - math.Pi*2.0
	}

	me.arcToCurve2(path, startX, startY, radiusX, radiusY, angle, sweep, endX, endY, startAngle, endAngle, centerX, centerY, endControlX, endControlY)
}

func (me *imageSvg) PixelSize() (width, height int32) {
	return int32(me.width), int32(me.height)
}

func (me *imageSvg) Draw(canvas Canvas) {
	canvas.Save()
	canvas.Scale(float32(me.width)/float32(me.viewBox.Width), float32(me.height)/float32(me.viewBox.Height))
	for _, p := range me.paths {
		canvas.DrawPath(p.path, p.paint)
	}
	canvas.Restore()
}

// --------------- pathparse ----------------------------
func (me *pathparse) isdigit(c uint8) bool {
	return (c >= '0' && c <= '9') || c == '.' || c == '-'
}

func (me *pathparse) isalphabet(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func (me *pathparse) skipBlank() {
	for {
		switch me.d[me.pos] {
		case ' ', '\r', '\n', '\t', '\f', '\b', ',':
			me.pos++
			if me.pos >= me.len {
				break
			}
			continue
		}
		break
	}
}

func (me *pathparse) getPoints(count int) [8]float32 {
	me.skipBlank()

	if me.isalphabet(me.d[me.pos]) {
		me.pos++
	}

	begin := false
	var c uint8 = 0
	for i := 0; i != count; i++ {
		me.skipBlank()

		me.start = me.pos
		begin = true
		for me.pos < me.len {
			c = me.d[me.pos]
			if (c >= '0' && c <= '9') || c == '.' || (c == '-' && begin) {
				begin = false
				me.pos++
				continue
			} else {
				break
			}
		}

		n, e := strconv.ParseFloat(me.d[me.start:me.pos], 32)
		if e != nil {
			log.Fatal("nuxui", "parse svg path error, %s", e.Error())
		}
		me.points[i] = float32(n)
	}

	return me.points
}
