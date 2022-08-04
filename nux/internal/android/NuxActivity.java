// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package org.nuxui.app;

import android.app.Activity;

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
    private static final String KEY_NATIVE_SAVED_STATE = "android:native_state";

    private native void    native_NuxActivity_onCreate(byte[] nativeSavedState);
    private native void    native_NuxActivity_onStart();
    private native void    native_NuxActivity_onRestart();
    private native void    native_NuxActivity_onResume();
    private native void    native_NuxActivity_onPause();
    private native void    native_NuxActivity_onStop();
    private native void    native_NuxActivity_onDestroy();
    private native void    native_NuxActivity_surfaceRedrawNeeded(SurfaceHolder holder);
    private native void    native_NuxActivity_surfaceCreated(SurfaceHolder holder);
    private native void    native_NuxActivity_surfaceChanged(SurfaceHolder holder, int format, int width, int height);
    private native void    native_NuxActivity_surfaceDestroyed(SurfaceHolder holder);
    private native boolean native_NuxActivity_onTouch(int deviceId, int pointerId, int action, float x, float y);

    // private Handler mHandler;
    private NuxView mNuxView;

    // private void load(){
    //     String libname = "nuxui";
    //     ActivityInfo ai;
    //     try {
    //         ai = getPackageManager().getActivityInfo(
    //                 getIntent().getComponent(), PackageManager.GET_META_DATA);
    //         if (ai.metaData != null) {
    //             String ln = ai.metaData.getString(META_DATA_LIB_NAME);
    //             if (ln != null) libname = ln;
    //         }
    //     } catch (PackageManager.NameNotFoundException e) {
    //         throw new RuntimeException("Error getting activity info", e);
    //     }

    //     System.loadLibrary(libname);
    // }

    @Override
    protected void onCreate( Bundle savedInstanceState) {
        Log.i("nuxui", "onCreate myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());

		requestWindowFeature(Window.FEATURE_NO_TITLE);
        
        super.onCreate(savedInstanceState);

        // load();

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            Window window = this.getWindow();
            window.addFlags(WindowManager.LayoutParams.FLAG_DRAWS_SYSTEM_BAR_BACKGROUNDS);
            window.setStatusBarColor(0x40000000);
            window.getDecorView().setSystemUiVisibility(View.SYSTEM_UI_FLAG_LAYOUT_STABLE | View.SYSTEM_UI_FLAG_LAYOUT_FULLSCREEN);
        }

        // mHandler = new Handler(getMainLooper(), this);

        mNuxView = new NuxView(this);
        mNuxView.getHolder().addCallback(this);
        ViewGroup.LayoutParams params = new ViewGroup.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.MATCH_PARENT);
        setContentView(mNuxView, params);
        mNuxView.requestFocus();

        byte[] nativeSavedState = savedInstanceState != null
                ? savedInstanceState.getByteArray(KEY_NATIVE_SAVED_STATE) : null;
        native_NuxActivity_onCreate(nativeSavedState);
    }

    @Override
    protected void onStart() {
        super.onStart();
        native_NuxActivity_onStart();
    }

    @Override
    protected void onRestart() {
        super.onRestart();
        native_NuxActivity_onRestart();
    }

    @Override
    protected void onResume() {
        super.onResume();
        native_NuxActivity_onResume();
    }

    @Override
    protected void onPause() {
        super.onPause();
        native_NuxActivity_onPause();
    }

    @Override
    protected void onStop() {
        super.onStop();
        native_NuxActivity_onStop();
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        native_NuxActivity_onDestroy();
    }

    @Override
    public void surfaceRedrawNeeded(SurfaceHolder holder) {
        Log.i("nuxui", "surfaceRedrawNeeded myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        native_NuxActivity_surfaceRedrawNeeded(holder);
    }

    @Override
    public void surfaceCreated(SurfaceHolder holder) {
        Log.i("nuxui", "surfaceCreated myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        native_NuxActivity_surfaceCreated(holder);
    }

    @Override
    public void surfaceChanged(SurfaceHolder holder, int format, int width, int height) {
        Log.i("nuxui", "surfaceChanged myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        native_NuxActivity_surfaceChanged(holder, format, width, height);
    }

    @Override
    public void surfaceDestroyed(SurfaceHolder holder) {
        Log.i("nuxui", "surfaceDestroyed myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        native_NuxActivity_surfaceDestroyed(holder);
    }

    @Override
    public boolean onTouchEvent(MotionEvent e) {
        // int pointerCounet = e.getPointerCount();
        // int maskedAction = e.getActionMasked();
        // int actionIndex = e.getActionIndex();
        // boolean handled = false;
        // if (maskedAction == MotionEvent.ACTION_DOWN || maskedAction == MotionEvent.ACTION_POINTER_DOWN){
        //     handled = native_NuxActivity_onTouch(e.getDeviceId(), e.getPointerId(actionIndex), MotionEvent.ACTION_DOWN, e.getX(actionIndex), e.getY(actionIndex));
        // }else if(maskedAction == MotionEvent.ACTION_UP || maskedAction == MotionEvent.ACTION_POINTER_UP){
        //     handled = native_NuxActivity_onTouch(e.getDeviceId(), e.getPointerId(actionIndex), MotionEvent.ACTION_UP, e.getX(actionIndex), e.getY(actionIndex));
        // }else if(maskedAction == MotionEvent.ACTION_MOVE){
        //     for(int i =0; i!= pointerCounet; i++){
        //         handled = native_NuxActivity_onTouch(e.getDeviceId(), e.getPointerId(i), MotionEvent.ACTION_MOVE, e.getX(i), e.getY(i));
        //     }
        // }else if(maskedAction == MotionEvent.ACTION_CANCEL||maskedAction == MotionEvent.ACTION_OUTSIDE){
        //     handled = native_NuxActivity_onTouch(e.getDeviceId(), e.getPointerId(actionIndex), maskedAction, e.getX(actionIndex), e.getY(actionIndex));
        // }

        // if (handled){
        //     return true;
        // }
        return super.onTouchEvent(e);
    }

    @Override
    public boolean onGenericMotionEvent(MotionEvent event) {
        // ACTION_HOVER_MOVE
        // ACTION_SCROLL
        // ACTION_HOVER_ENTER
        // ACTION_HOVER_EXIT
        // ACTION_BUTTON_PRESS
        // ACTION_BUTTON_RELEASE
        return super.onGenericMotionEvent(event);
    }

    public static void drawText(Canvas canvas ,String text, int width,  TextPaint paint){
        Log.i("nuxui", "drawText myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
        StaticLayout mTextLayout = new StaticLayout(text, paint, width, Layout.Alignment.ALIGN_NORMAL, 1.0f, 0, false);
        mTextLayout.draw(canvas);
    }

    public static StaticLayout createStaticLayout(String text, int width,  TextPaint paint){
        return new StaticLayout(text, paint, width, Layout.Alignment.ALIGN_NORMAL, 1.0f, 0, false);
    }

    public static Bitmap createBitmap(String fileName){
        BitmapFactory.Options options = new BitmapFactory.Options();
        Bitmap b = null;
        try{
            if (fileName.startsWith("assets/")){
                Log.i("nuxui", "createBitmap myPid=" + android.os.Process.myPid() + ", myTid="+android.os.Process.myTid()+", thread="+Thread.currentThread().getId());
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

    // public void backToUI(){
    //     Message msg = Message.obtain(mHandler, 100);
    //     mHandler.sendMessage(msg);
    // }

    // @Override
    // public boolean handleMessage(Message msg) {
    //     switch(msg.what){
    //         case 100:   // backToUI
    //             native_NuxActivity_onBackToUI();
    //             break;
    //     }
    //     return false;
    // }
}
