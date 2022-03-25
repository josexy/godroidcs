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
import android.content.Intent;
import android.net.Uri;

public final class PhoneUtil {

    private PhoneUtil() {
    }

    /**
     * open dial phone UI
     *
     * @param context
     * @param phoneNumber
     */
    public static void dialPhone(Context context, String phoneNumber) {
        Intent dial_phone = new Intent(Intent.ACTION_DIAL);
        dial_phone.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        dial_phone.setData(Uri.parse("tel:" + phoneNumber));
        context.startActivity(dial_phone);
    }

    /**
     * call phone number
     *
     * @param context
     * @param phoneNumber
     */
    public static void callPhone(Context context, String phoneNumber) {
        Intent call_phone = new Intent(Intent.ACTION_CALL);
        call_phone.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        call_phone.setData(Uri.parse("tel:" + phoneNumber));
        context.startActivity(call_phone);
    }
}
