// Copyright [2021] [josexy]
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com.joxrays.godroidsvr.util;

import android.content.ClipData;
import android.content.Context;

public final class ClipboardUtil {

    private ClipboardUtil() {

    }

    public static void setText(Context context, String text) {
        CommonUtil.getClipboardManager(context).setPrimaryClip(
                ClipData.newPlainText(context.getPackageName(), text)
        );
    }

    public static String getText(Context context) {
        ClipData data = CommonUtil.getClipboardManager(context).getPrimaryClip();
        if (data != null && data.getItemCount() > 0) {
            return data.getItemAt(0).coerceToText(context).toString();
        }
        return "";
    }
}
