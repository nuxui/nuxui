// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"os"
	"path/filepath"
	"strings"
)

type Image interface {
	PixelSize() (width, height int32)
	Draw(canvas Canvas)
}

func CreateImage(filename string) (Image, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	filename, err = filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".svg":
		return NewImageSVGFromFile(filename), nil
	case ".png", ".jpg", ".jpeg":
		return createImage(filename), nil
	}

	return nil, nil
}
