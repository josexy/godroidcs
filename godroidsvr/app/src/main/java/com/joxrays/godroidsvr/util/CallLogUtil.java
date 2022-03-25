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
import android.provider.CallLog;

import com.joxrays.godroidsvr.base.SortOrderType;
import com.joxrays.godroidsvr.message.CallLogInfo;
import com.joxrays.godroidsvr.message.CallLogMetaInfo;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public final class CallLogUtil {

    private CallLogUtil() {

    }

    public static String getCallLogType(int value) {
        String type;
        switch (value) {
            case CallLog.Calls.INCOMING_TYPE:
                type = "Incoming";
                break;
            case CallLog.Calls.OUTGOING_TYPE:
                type = "Outgoing";
                break;
            case CallLog.Calls.MISSED_TYPE:
                type = "Missed";
                break;
            case CallLog.Calls.VOICEMAIL_TYPE:
                type = "Voicemail";
                break;
            case CallLog.Calls.REJECTED_TYPE:
                type = "Rejected";
                break;
            case CallLog.Calls.BLOCKED_TYPE:
                type = "Blocked";
                break;
            case CallLog.Calls.ANSWERED_EXTERNALLY_TYPE:
                type = "Answered on another device";
                break;
            default:
                type = "Unknown";
                break;
        }
        return type;
    }

    public static List<CallLogInfo> getCallLogInfo(Context context, String number) {
        List<CallLogInfo> list = new ArrayList<>();
        Cursor cursor = ContentResolverUtil.query(context, CallLog.Calls.CONTENT_URI,
                new String[]{
                        CallLog.Calls._ID,
                        CallLog.Calls.DURATION,
                        CallLog.Calls.DATE,
                        CallLog.Calls.TYPE,
                }, CallLog.Calls.NUMBER + "=?", new String[]{number}, null, SortOrderType.Default);
        if (cursor != null) {
            while (cursor.moveToNext()) {
                String id = cursor.getString(0);
                String duration = cursor.getString(1);
                String date = cursor.getString(2);
                String type = cursor.getString(3);
                int value = Integer.parseInt(type);

                list.add(CallLogInfo.newBuilder()
                        .setId(Integer.parseInt(id))
                        .setNumber(CommonUtil.formatPhoneNumber(context, number))
                        .setDuration(Integer.parseInt(duration))
                        .setDate(Long.parseLong(date))
                        .setType(getCallLogType(value))
                        .build());
            }
            cursor.close();
        }
        return list;
    }

    public static List<CallLogMetaInfo> getAllCallLogInfo(Context context) {
        List<CallLogMetaInfo> list = new ArrayList<>();
        Cursor cursor = ContentResolverUtil.query(context, CallLog.Calls.CONTENT_URI,
                new String[]{
                        CallLog.Calls.NUMBER,
                        CallLog.Calls.DATE,
                }, null, null, CallLog.Calls.DATE, SortOrderType.ASC);
        Set<String> set = new HashSet<>();
        if (cursor != null) {
            while (cursor.moveToNext()) {
                String number = cursor.getString(0);
                if (set.contains(number)) {
                    continue;
                }
                set.add(number);
                String date = cursor.getString(1);

                list.add(CallLogMetaInfo.newBuilder()
                        .setNumber(number)
                        .setDate(Long.parseLong(date))
                        .build());
            }
            cursor.close();
            return list;
        }
        return null;
    }

    public static void deleteCallLog(Context context, String number) {
        ContentResolverUtil.delete(context, CallLog.Calls.CONTENT_URI, CallLog.Calls.NUMBER + "=?", new String[]{number});
    }
}

