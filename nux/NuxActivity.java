// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package org.nuxui.app;

import android.app.Activity;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageManager;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Canvas;
import android.graphics.Paint;
import android.os.Build;
import android.os.Bundle;
import android.text.Layout;
import android.text.StaticLayout;
import android.text.TextPaint;
import android.util.Log;
import android.view.MotionEvent;
import android.view.SurfaceHolder;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.view.WindowManager;
import android.content.res.Resources;
import android.graphics.RectF;
import android.graphics.Rect;
import java.io.IOException;

public class NuxActivity extends Activity implements SurfaceHolder.Callback2 {
    public static final String META_DATA_LIB_NAME = "android.app.lib_name";
    private static final String KEY_NATIVE_SAVED_STATE = "android:native_state";

    private native void onCreateNative(byte[] nativeSavedState, float density);
    private native void onStartNative();
    private native void onRestartNative();
    private native void onResumeNative();
    private native void onPauseNative();
    private native void onStopNative();
    private native void onDestroyNative();
    private native void surfaceRedrawNeededNative(SurfaceHolder holder);
    private native void surfaceCreatedNative(SurfaceHolder holder);
    private native void surfaceChangedNative(SurfaceHolder holder, int format, int width, int height);
    private native void surfaceDestroyedNative(SurfaceHolder holder);

    private NuxView mNuxView;

    private void load(){
        String libname = "nuxui";
        ActivityInfo ai;
        try {
            ai = getPackageManager().getActivityInfo(
                    getIntent().getComponent(), PackageManager.GET_META_DATA);
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
    protected void onCreate( Bundle savedInstanceState) {
        Log.i("nuxui", "onCreate myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());

		requestWindowFeature(Window.FEATURE_NO_TITLE);
        
        super.onCreate(savedInstanceState);

        load();

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            Window window = this.getWindow();
            window.addFlags(WindowManager.LayoutParams.FLAG_DRAWS_SYSTEM_BAR_BACKGROUNDS);
            window.setStatusBarColor(0x40000000);
            window.getDecorView().setSystemUiVisibility(View.SYSTEM_UI_FLAG_LAYOUT_STABLE | View.SYSTEM_UI_FLAG_LAYOUT_FULLSCREEN);
        }

        mNuxView = new NuxView(this);
        mNuxView.getHolder().addCallback(this);
        ViewGroup.LayoutParams params = new ViewGroup.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.MATCH_PARENT);
        setContentView(mNuxView, params);
        mNuxView.requestFocus();

        byte[] nativeSavedState = savedInstanceState != null
                ? savedInstanceState.getByteArray(KEY_NATIVE_SAVED_STATE) : null;
        onCreateNative(nativeSavedState, Resources.getSystem().getDisplayMetrics().density);
    }

    @Override
    protected void onStart() {
        super.onStart();
        onStartNative();
    }

    @Override
    protected void onRestart() {
        super.onRestart();
        onRestartNative();
    }

    @Override
    protected void onResume() {
        super.onResume();
        onResumeNative();
    }

    @Override
    protected void onPause() {
        super.onPause();
        onPauseNative();
    }

    @Override
    protected void onStop() {
        super.onStop();
        onStopNative();
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        onDestroyNative();
    }

    @Override
    public void surfaceRedrawNeeded(SurfaceHolder holder) {
        Log.i("nuxui", "surfaceRedrawNeeded myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        surfaceRedrawNeededNative(holder);
    }

    @Override
    public void surfaceCreated(SurfaceHolder holder) {
        Log.i("nuxui", "surfaceCreated myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        surfaceCreatedNative(holder);
    }

    @Override
    public void surfaceChanged(SurfaceHolder holder, int format, int width, int height) {
        Log.i("nuxui", "surfaceChanged myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        surfaceChangedNative(holder, format, width, height);
    }

    @Override
    public void surfaceDestroyed(SurfaceHolder holder) {
        Log.i("nuxui", "surfaceDestroyed myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        surfaceDestroyedNative(holder);
    }

    @Override
    public boolean onTouchEvent(MotionEvent event) {
        Log.i("nuxui", "onTouchEvent myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        return super.onTouchEvent(event);
    }

    public static void drawText(Canvas canvas ,String text, int width,  TextPaint paint){
        Log.i("nuxui", "drawText myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        StaticLayout mTextLayout = new StaticLayout(text, paint, width, Layout.Alignment.ALIGN_NORMAL, 1.0f, 0, false);
        mTextLayout.draw(canvas);
    }

    public static StaticLayout createStaticLayout(String text, int width,  TextPaint paint){
        return new StaticLayout(text, paint, width, Layout.Alignment.ALIGN_NORMAL, 1.0f, 0, false);
    }

    public static Bitmap createImage(String fileName){
        BitmapFactory.Options options = new BitmapFactory.Options();
        Bitmap b = null;
        try{
            if (fileName.startsWith("assets/")){
                Log.i("nuxui", "createImage myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
                String s = fileName.substring(7, fileName.length());
                b = BitmapFactory.decodeStream(NuxApplication.instance().getAssets().open(s));
                Log.i("nuxui", "startsWith assets: " + s + " b=" + b);
                return b;
            }else{
                Log.i("nuxui", "not startsWith assets: " + fileName);
                b = BitmapFactory.decodeFile(fileName, options);
            }
        } catch (IOException e) {
            Log.i("nuxui", "startsWith assets " + e);
            e.printStackTrace();
        }
        return b;
    }
}
