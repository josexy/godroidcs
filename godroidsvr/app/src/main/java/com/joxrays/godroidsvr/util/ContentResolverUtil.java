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

import android.content.Context;
import android.database.Cursor;
import android.graphics.Bitmap;
import android.net.Uri;
import android.util.Pair;
import android.util.Size;

import com.joxrays.godroidsvr.base.SortOrderType;

import java.io.InputStream;

public final class ContentResolverUtil {

    private ContentResolverUtil() {

    }

    public static Cursor query(Context context, Uri uri, String[] projection,
                               String selection, String[] selectionArgs,
                               String sortColumn, SortOrderType type) {
        if (sortColumn != null) {
            sortColumn += " " + type.name;
        }
        return context.getContentResolver().query(uri, projection, selection, selectionArgs, sortColumn);
    }

    public static int delete(Context context, Uri uri, String where, String[] selectionArgs) {
        return context.getContentResolver().delete(uri, where, selectionArgs);
    }

    public static Pair<InputStream, Exception> openFile(Context context, Uri uri) {
        try {
            return Pair.create(context.getContentResolver().openInputStream(uri), null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static Pair<Bitmap, Exception> getThumbnail(Context context, Uri uri, int width, int height) {
        try {
            return Pair.create(
                    context.getContentResolver().loadThumbnail(uri, new Size(width, height), null),
                    null
            );
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }
}
