// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Manifest interface {
	Name() string
	AppID() string
	Package() string
	Version() string
	VersionCode() string
	VersionID() string
	Permissions() []string

	// Main main widget
	Main() string
}

func NewManifest() Manifest {
	return newManifest()
}

type manifest struct {
	name        string
	appId       string
	packageName string
	version     string
	versionCode string
	versionId   string
	permissions []string
	main        string
}

func newManifest() *manifest {
	return &manifest{}
}

func (me *manifest) Creating(attr Attr) {
	me.name = attr.GetString("name", "")
	me.appId = attr.GetString("appId", "")
	me.packageName = attr.GetString("package", "")
	me.version = attr.GetString("version", "")
	me.versionCode = attr.GetString("versionCode", "")
	me.versionId = attr.GetString("versionId", "")
	me.permissions = attr.GetStringArray("permissions", []string{})
	me.main = attr.GetString("main", "")
}

func (me *manifest) Name() string {
	return me.name
}

func (me *manifest) AppID() string {
	return me.appId
}

func (me *manifest) Package() string {
	return me.packageName
}

func (me *manifest) Version() string {
	return me.version
}

func (me *manifest) VersionCode() string {
	return me.versionCode
}

func (me *manifest) VersionID() string {
	return me.versionId
}

func (me *manifest) Permissions() []string {
	return me.permissions
}

func (me *manifest) Main() string {
	return me.main
}
