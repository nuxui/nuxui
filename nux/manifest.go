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

var manifestInstance *manifest = &manifest{}

func NewManifest(attr Attr) Manifest {
	if attr == nil {
		attr = Attr{}
	}

	m := manifestInstance
	m.name = attr.GetString("name", "")
	m.appId = attr.GetString("appId", "")
	m.packageName = attr.GetString("package", "")
	m.version = attr.GetString("version", "")
	m.versionCode = attr.GetString("versionCode", "")
	m.versionId = attr.GetString("versionId", "")
	m.permissions = attr.GetStringArray("permissions", []string{})
	m.main = attr.GetString("main", "")
	m.multiWindow = attr.GetBool("multiWindow", false)

	return m
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
	multiWindow bool
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
