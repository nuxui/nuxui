// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"path/filepath"
	"strings"
)

var _ Image = (*nativeImage)(nil)

type Image interface {
	PixelSize() (width, height int32)
	Draw(canvas Canvas)
}

func LoadImageFromFile(filename string) (Image, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".svg":
		return LoadImageSVGFromFile(filename), nil
	case ".png", ".jpg", ".jpeg":
		return loadImageFromFile(filename), nil
	}

	return nil, nil
}
