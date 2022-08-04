// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package org.nuxui.app;

import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.app.Application;
import android.os.Handler;
import android.os.Message;
import android.util.Log;
import android.content.res.Resources;

public class NuxApplication extends Application implements Handler.Callback {
    public static final String META_DATA_LIB_NAME = "org.nuxui.app.libname";

    private native void native_NuxApplication_onCreate(float density);
    private native void native_NuxApplication_onBackToUI();

    private static NuxApplication mInstance;
    private Handler mHandler;

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
    public void onCreate() {
        Log.i("nuxui", "NuxApplication onCreate");

        super.onCreate();
        mInstance = this;

        Log.i("nuxui", "NuxApplication load");
        load();

        mHandler = new Handler(getMainLooper(), this);

        Log.i("nuxui", "NuxApplication lonative_NuxApplication_onCreatead");
        native_NuxApplication_onCreate(Resources.getSystem().getDisplayMetrics().density);
    }

    public static NuxApplication instance(){
        return mInstance;
    }

    public void backToUI(){
        Message msg = Message.obtain(mHandler, 100);
        mHandler.sendMessage(msg);
    }

    @Override
    public boolean handleMessage(Message msg) {
        switch(msg.what){
            case 100:   // backToUI
                native_NuxApplication_onBackToUI();
                break;
        }
        return false;
    }
}
