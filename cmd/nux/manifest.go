// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
)

type manifestXML struct {
	Activity activityXML `xml:"application>activity"`
}

type activityXML struct {
	Name     string        `xml:"name,attr"`
	MetaData []metaDataXML `xml:"meta-data"`
}

type metaDataXML struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// manifestLibName parses the AndroidManifest.xml and finds the library
// name of the NativeActivity.
func manifestLibName(data []byte) (string, error) {
	manifest := new(manifestXML)
	if err := xml.Unmarshal(data, manifest); err != nil {
		return "", err
	}
	if manifest.Activity.Name != "org.nuxui.app.NuxActivity" {
		return "", fmt.Errorf("can only build an .apk for NuxActivity, not %q", manifest.Activity.Name)
	}
	libName := ""
	for _, md := range manifest.Activity.MetaData {
		if md.Name == "android.app.lib_name" {
			libName = md.Value
			break
		}
	}
	if libName == "" {
		return "", errors.New("AndroidManifest.xml missing meta-data android.app.lib_name")
	}
	return libName, nil
}

type manifestTmplData struct {
	JavaPkgPath string
	Name        string
	LibName     string
}

var manifestTmpl = template.Must(template.New("manifest").Parse(`
<manifest 
    xmlns:android="http://schemas.android.com/apk/res/android" 
	package="{{.JavaPkgPath}}"
	android:versionCode="1"
	android:versionName="1.0"
    >
    <application 
        android:debuggable="true" 
        android:label="{{.Name}}" 
        android:name="org.nuxui.app.NuxApplication">
        <meta-data android:name="org.nuxui.app.libname" android:value="{{.LibName}}"/>
        <activity 
			android:configChanges="keyboardHidden|orientation" 
			android:label="{{.Name}}" 
			android:name="org.nuxui.app.NuxActivity">
            <intent-filter>
                <action android:name="android.intent.action.MAIN"/>
                <category android:name="android.intent.category.LAUNCHER"/>
            </intent-filter>
        </activity>
    </application>
</manifest>`))
