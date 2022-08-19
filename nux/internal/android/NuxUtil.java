// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package org.nuxui.app;

import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Canvas;
import android.graphics.Paint;
import android.text.Layout;
import android.text.StaticLayout;
import android.text.TextPaint;
import android.util.Log;
import java.io.IOException;

public class NuxUtil {

    public static void setPaintStyle(Paint paint, int style) {
        if (style == 0) {
            paint.setStyle(Paint.Style.FILL);
        } else if (style == 1) {
            paint.setStyle(Paint.Style.STROKE);
        } else {
            paint.setStyle(Paint.Style.FILL_AND_STROKE);
        }
    }

    public static int getPaintStyle(Paint paint) {
        Paint.Style style = paint.getStyle();
        switch (style) {
            case FILL: return 0;
            case STROKE: return 1;
            case FILL_AND_STROKE: return 2;
        }
        return 0;
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

}