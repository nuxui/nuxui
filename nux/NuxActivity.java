// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package org.nuxui.app;

import android.app.Activity;
import android.app.NativeActivity;
import android.content.Context;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.util.Log;
import android.view.KeyCharacterMap;
import android.view.WindowManager;
import android.view.Window;

public class NuxActivity extends NativeActivity {
	// private static NuxActivity nuxActivity;

	// public NuxActivity() {
	// 	super();
	// 	nuxActivity = this;
	// }

	// String getTmpdir() {
	// 	return getCacheDir().getAbsolutePath();
	// }

	// static int getRune(int deviceId, int keyCode, int metaState) {
	// 	try {
	// 		int rune = KeyCharacterMap.load(deviceId).get(keyCode, metaState);
	// 		if (rune == 0) {
	// 			return -1;
	// 		}
	// 		return rune;
	// 	} catch (KeyCharacterMap.UnavailableException e) {
	// 		return -1;
	// 	} catch (Exception e) {
	// 		Log.e("nuxui", "exception reading KeyCharacterMap", e);
	// 		return -1;
	// 	}
	// }

	@Override
	public void onCreate(Bundle savedInstanceState) {
		Log.info("nuxui", "activity onCreate")
		// requestWindowFeature(Window.FEATURE_NO_TITLE);
		// getWindow().setFlags(WindowManager.LayoutParams.FLAG_FULLSCREEN,WindowManager.LayoutParams.FLAG_FULLSCREEN);

		setTranslucentStatus(0);


		super.onCreate(savedInstanceState);
		Log.info("nuxui", "activity onCreate after super.onCreate(savedInstanceState)")
	}

	private void setTranslucentStatus(int on) {
		Window win = getWindow();  
		WindowManager.LayoutParams winParams = win.getAttributes();  
		final int bits = WindowManager.LayoutParams.FLAG_TRANSLUCENT_STATUS;  
		if (on <= 0) {  
			winParams.flags &= ~bits;  
		} else {  
			winParams.flags |= bits;  
		}  
		win.setAttributes(winParams);  
	}

	public int getStatusBarHeight() { 
		int result = 0;
		int resourceId = getResources().getIdentifier("status_bar_height", "dimen", "android");
		if (resourceId > 0) {
			result = getResources().getDimensionPixelSize(resourceId);
		} 
		return result;
  	} 
}
