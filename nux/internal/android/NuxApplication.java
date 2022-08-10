// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package org.nuxui.app;

import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.app.Application;
import android.os.Handler;
import android.os.Message;
import android.content.res.Resources;
import android.content.res.Configuration;
import android.util.DisplayMetrics;
import android.util.Log;

public class NuxApplication extends Application {
    public static final String META_DATA_LIB_NAME = "org.nuxui.app.libname";

    private native void native_NuxApplication_onConfigurationChanged(Configuration newConfig);
    private native void native_NuxApplication_onCreate(float density, int densityDpi, float scaledDensity, int widthPixels, int heightPixels, float xdpi, float ydpi);
    private native void native_NuxApplication_onLowMemory();
    private native void native_NuxApplication_onTerminate();
    private native void native_NuxApplication_onTrimMemory(int level);

    private static NuxApplication mInstance;
    public  static NuxApplication instance(){return mInstance;}

    private void load(){
        String libname = "nuxui";
        ApplicationInfo ai;
        try {
            ai = getPackageManager().getApplicationInfo(getPackageName(), PackageManager.GET_META_DATA);
            if (ai.metaData != null) {
                String ln = ai.metaData.getString(META_DATA_LIB_NAME);
                if (ln != null) libname = ln;
            }
        } catch (PackageManager.NameNotFoundException e) {
            throw new RuntimeException("Error getting activity info", e);
        }

        System.loadLibrary(libname);
    }

    @Override
    public void onConfigurationChanged(Configuration newConfig) {
        super.onConfigurationChanged(newConfig);
        native_NuxApplication_onConfigurationChanged(newConfig);
    }

    @Override
    public void onCreate() {
        super.onCreate();
        mInstance = this;

        load();

        DisplayMetrics dm = Resources.getSystem().getDisplayMetrics();
        native_NuxApplication_onCreate(dm.density, dm.densityDpi, dm.scaledDensity, dm.widthPixels, dm.heightPixels, dm.xdpi, dm.ydpi);
    }

    @Override
    public void onLowMemory() {
        super.onLowMemory();
        native_NuxApplication_onLowMemory();
    }

    @Override
    public void onTerminate() {
        super.onTerminate();
        native_NuxApplication_onTerminate();
    }

    @Override
    public void onTrimMemory(int level) {
        super.onTrimMemory(level);
        native_NuxApplication_onTrimMemory(level);
    }
}
