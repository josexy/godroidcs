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

import android.content.ContentUris;
import android.content.Context;
import android.database.Cursor;
import android.net.Uri;
import android.provider.MediaStore;
import android.util.Pair;

import com.joxrays.godroidsvr.base.SortOrderType;
import com.joxrays.godroidsvr.message.MediaStoreInfo;

import java.util.ArrayList;
import java.util.List;

public final class MediaStoreUtil {
    public static final Uri IMAGE_URI = MediaStore.Images.Media.EXTERNAL_CONTENT_URI;
    public static final Uri VIDEO_URI = MediaStore.Video.Media.EXTERNAL_CONTENT_URI;
    public static final Uri AUDIO_URI = MediaStore.Audio.Media.EXTERNAL_CONTENT_URI;
    public static final Uri DOWNLOAD_URI = MediaStore.Downloads.EXTERNAL_CONTENT_URI;

    // DCIM/ Pictures/
    public static final int IMAGE_TYPE = 0;
    // DCIM/ Movies/ Pictures/
    public static final int VIDEO_TYPE = 1;
    // Alarms/ Audiobooks/ Music/ Notifications/ Podcasts/ Ringtones/ Music/ Movies/
    public static final int AUDIO_TYPE = 2;
    // Download/
    public static final int DOWNLOAD_TYPE = 3;

    public static final String[] columns = new String[]{
            MediaStore.MediaColumns._ID,
            MediaStore.MediaColumns.DISPLAY_NAME,
            MediaStore.MediaColumns.SIZE,
            MediaStore.MediaColumns.DATE_ADDED,
            MediaStore.MediaColumns.DATE_MODIFIED,
    };

    private MediaStoreUtil() {

    }

    private static Pair<Cursor, Uri> openCursor(Context context, int type) {
        Uri uri = null;
        switch (type) {
            case IMAGE_TYPE:
                uri = IMAGE_URI;
                break;
            case AUDIO_TYPE:
                uri = AUDIO_URI;
                break;
            case VIDEO_TYPE:
                uri = VIDEO_URI;
                break;
            case DOWNLOAD_TYPE:
                uri = DOWNLOAD_URI;
                break;
        }
        return Pair.create(
                ContentResolverUtil.query(context, uri, columns, null, null, null, SortOrderType.Default),
                uri
        );
    }

    public static void deleteMediaFile(Context context, Uri uri) {
        ContentResolverUtil.delete(context, uri, null, null);
    }

    public static List<MediaStoreInfo> getAllMediaFilesInfo(Context context, int type) {
        Pair<Cursor, Uri> pair = openCursor(context, type);
        Cursor cursor = pair.first;

        List<MediaStoreInfo> list = new ArrayList<>();
        if (cursor != null) {
            while (cursor.moveToNext()) {
                String id = cursor.getString(0);
                String name = cursor.getString(1);
                String size = cursor.getString(2);
                String date_add = cursor.getString(3);
                String date_modify = cursor.getString(4);
                Uri uri = ContentUris.withAppendedId(pair.second, Long.parseLong(id));
                MediaStoreInfo info = MediaStoreInfo.newBuilder()
                        .setId(Integer.parseInt(id))
                        .setName(name)
                        .setSize(Long.parseLong(size != null ? size : "0"))
                        .setDateAdd(Long.parseLong(date_add != null ? date_add : "0"))
                        .setDateModify(Long.parseLong(date_modify != null ? date_modify : "0"))
                        .setUri(uri.toString())
                        .build();
                list.add(info);
            }
            cursor.close();
        }
        return list;
    }
}
