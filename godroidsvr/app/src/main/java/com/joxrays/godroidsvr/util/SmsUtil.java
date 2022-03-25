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

import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.database.Cursor;
import android.provider.Telephony;
import android.telephony.SmsManager;

import com.joxrays.godroidsvr.base.SortOrderType;
import com.joxrays.godroidsvr.message.SmsInfo;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public final class SmsUtil {

    private SmsUtil() {
    }

    public static List<SmsInfo> getSmsInfo(Context context, String address) {

        List<SmsInfo> list = new ArrayList<>();

        // inbox and outbox
        Cursor cursor = ContentResolverUtil.query(context, Telephony.Sms.CONTENT_URI,
                new String[]{
                        Telephony.Sms._ID,
                        Telephony.Sms.DATE_SENT,
                        Telephony.Sms.DATE,
                        Telephony.Sms.READ,
                        Telephony.Sms.BODY,
                        Telephony.Sms.TYPE,
                },
                Telephony.Sms.ADDRESS + "=?", new String[]{address},
                Telephony.Sms._ID, SortOrderType.ASC);
        if (cursor != null) {
            while (cursor.moveToNext()) {
                String id = cursor.getString(0);
                String sent = cursor.getString(1);
                String received = cursor.getString(2);
                String read = cursor.getString(3);
                String body = cursor.getString(4);
                int type = Integer.parseInt(cursor.getString(5));

                list.add(SmsInfo.newBuilder()
                        .setId(Integer.parseInt(id))
                        .setAddress(CommonUtil.formatPhoneNumber(context, address))
                        .setSentDate(Long.parseLong(sent))
                        .setReceivedDate(Long.parseLong(received))
                        .setRead(Boolean.parseBoolean(read))
                        .setBody(body)
                        .setSentReceived(type == Telephony.Sms.MESSAGE_TYPE_SENT)
                        .build());
            }
            cursor.close();
        }
        return list;
    }

    public static List<String> getAllSmsMetaInfo(Context context) {
        Set<String> set = new HashSet<>();
        Cursor cursor = ContentResolverUtil.query(context, Telephony.Sms.CONTENT_URI,
                new String[]{Telephony.Sms.ADDRESS}, null, null,
                null, SortOrderType.Default);
        if (cursor != null) {
            while (cursor.moveToNext()) {
                String address = cursor.getString(0);
                set.add(address);
            }
            cursor.close();
        }
        return new ArrayList<>(set);
    }

    public static void sendSms(Context context, String dest, String message) {
        SmsManager sm = SmsManager.getDefault();
        PendingIntent pi = PendingIntent.getBroadcast(context, 0, new Intent("SMS_SENT"), PendingIntent.FLAG_IMMUTABLE);
        sm.sendTextMessage(dest, null, message, pi, null);
    }
}
