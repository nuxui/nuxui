// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package org.nuxui.app;

import android.content.Context;
import android.util.AttributeSet;
import android.view.SurfaceView;

public class NuxView extends SurfaceView {
    public NuxView(Context context)
    {
        super(context);
        onCreate(context, null, 0);
    }

    public NuxView(Context context, AttributeSet attrs)
    {
        super(context, attrs);
        onCreate(context, attrs, 0);
    }

    public NuxView(Context context, AttributeSet attrs, int defStyle)
    {
        super(context, attrs, defStyle);
        onCreate(context, attrs, defStyle);
    }

    protected void onCreate(Context context, AttributeSet attrs, int defStyle)
    {
    }

}
