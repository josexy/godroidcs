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

import android.content.ContentProviderOperation;
import android.content.Context;
import android.database.Cursor;
import android.net.Uri;
import android.provider.ContactsContract;

import com.joxrays.godroidsvr.base.SortOrderType;
import com.joxrays.godroidsvr.message.ContactInfo;
import com.joxrays.godroidsvr.message.ContactMetaInfo;

import java.util.ArrayList;
import java.util.List;

public final class ContactUtil {

    private ContactUtil() {
    }

    private static int getPhoneTypeInt(String s) {
        switch (s) {
            case "Home":
                return ContactsContract.CommonDataKinds.Phone.TYPE_HOME;
            case "Work":
                return ContactsContract.CommonDataKinds.Phone.TYPE_WORK;
            case "Mobile":
                return ContactsContract.CommonDataKinds.Phone.TYPE_MOBILE;
            default:
                return ContactsContract.CommonDataKinds.Phone.TYPE_OTHER;
        }
    }

    private static int getEmailTypeInt(String s) {
        switch (s) {
            case "Home":
                return ContactsContract.CommonDataKinds.Email.TYPE_HOME;
            case "Work":
                return ContactsContract.CommonDataKinds.Email.TYPE_WORK;
            case "Mobile":
                return ContactsContract.CommonDataKinds.Email.TYPE_MOBILE;
            default:
                return ContactsContract.CommonDataKinds.Email.TYPE_OTHER;
        }
    }

    private static String getPhoneTypeString(int type) {
        switch (type) {
            case ContactsContract.CommonDataKinds.Phone.TYPE_HOME:
                return "Home";
            case ContactsContract.CommonDataKinds.Phone.TYPE_WORK:
                return "Work";
            case ContactsContract.CommonDataKinds.Phone.TYPE_MOBILE:
                return "Mobile";
            default:
                return "Other";
        }
    }

    private static String getEmailTypeString(int type) {
        switch (type) {
            case ContactsContract.CommonDataKinds.Email.TYPE_HOME:
                return "Home";
            case ContactsContract.CommonDataKinds.Email.TYPE_WORK:
                return "Work";
            case ContactsContract.CommonDataKinds.Email.TYPE_MOBILE:
                return "Mobile";
            default:
                return "Other";
        }
    }

    public static ContactInfo getContactInfo(Context context, String id) {
        ContactInfo.Builder builder = ContactInfo.newBuilder();
        builder.setId(Integer.parseInt(id));

        Cursor name = ContentResolverUtil.query(context, ContactsContract.Contacts.CONTENT_URI,
                new String[]{ContactsContract.Contacts.DISPLAY_NAME,},
                ContactsContract.Contacts._ID + "=?",
                new String[]{id}, null, SortOrderType.Default);
        if (name != null) {
            if (name.moveToNext()) {
                builder.setName(name.getString(0));
            }
            name.close();
        }

        Cursor phone = ContentResolverUtil.query(context, ContactsContract.CommonDataKinds.Phone.CONTENT_URI,
                new String[]{
                        ContactsContract.CommonDataKinds.Phone.NUMBER,
                        ContactsContract.CommonDataKinds.Phone.TYPE
                },
                ContactsContract.CommonDataKinds.Phone.CONTACT_ID + "=?",
                new String[]{id}, null, SortOrderType.Default);

        // query phone numbers information
        if (phone != null) {
            List<ContactInfo.PhoneInfo> list = new ArrayList<>();
            while (phone.moveToNext()) {
                String number = CommonUtil.formatPhoneNumber(context, phone.getString(0));
                String type = getPhoneTypeString(Integer.parseInt(phone.getString(1)));
                list.add(ContactInfo.PhoneInfo.newBuilder()
                        .setNumber(number)
                        .setType(type)
                        .build());
            }

            builder.addAllPhones(list);
            phone.close();
        }

        // query emails information
        Cursor email = ContentResolverUtil.query(context, ContactsContract.CommonDataKinds.Email.CONTENT_URI,
                new String[]{
                        ContactsContract.CommonDataKinds.Email.ADDRESS,
                        ContactsContract.CommonDataKinds.Email.TYPE
                },
                ContactsContract.CommonDataKinds.Email.CONTACT_ID + "=?",
                new String[]{id}, null, SortOrderType.Default);
        if (email != null) {
            List<ContactInfo.EmailInfo> list = new ArrayList<>();
            while (email.moveToNext()) {
                String email_ = email.getString(0);
                String type = getEmailTypeString(Integer.parseInt(email.getString(1)));
                list.add(ContactInfo.EmailInfo.newBuilder()
                        .setEmail(email_)
                        .setType(type)
                        .build());
            }
            builder.addAllEmails(list);
            email.close();
        }
        return builder.build();
    }

    public static List<ContactMetaInfo> getAllContactMetaInfo(Context context) {
        List<ContactMetaInfo> list = new ArrayList<>();
        Cursor cursor = ContentResolverUtil.query(context, ContactsContract.Contacts.CONTENT_URI,
                new String[]{
                        ContactsContract.Contacts._ID,
                        ContactsContract.Contacts.DISPLAY_NAME,
                        ContactsContract.Contacts.LOOKUP_KEY
                }, null, null, ContactsContract.Contacts._ID, SortOrderType.ASC);
        if (cursor != null) {
            while (cursor.moveToNext()) {
                String id = cursor.getString(0);
                String name = cursor.getString(1);
                // for delete operation, we need a Uri
                Uri uri = Uri.withAppendedPath(ContactsContract.Contacts.CONTENT_LOOKUP_URI, cursor.getString(2));

                list.add(ContactMetaInfo.newBuilder()
                        .setId(Integer.parseInt(id))
                        .setName(name)
                        .setUri(uri.toString())
                        .build());
            }
            cursor.close();
        }
        return list;
    }

    public static void deleteContact(Context context, String uri) {
        try {
            ContentResolverUtil.delete(context, Uri.parse(uri), null, null);
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }

    public static void addContact(Context context, ContactInfo info) {

        ArrayList<ContentProviderOperation> ops = new ArrayList<>();

        ContentProviderOperation op = ContentProviderOperation.newInsert(ContactsContract.RawContacts.CONTENT_URI)
                .withValue(ContactsContract.RawContacts.ACCOUNT_TYPE, null)
                .withValue(ContactsContract.RawContacts.ACCOUNT_NAME, null)
                .build();

        ops.add(op);

        // Name
        ops.add(
                ContentProviderOperation.newInsert(ContactsContract.Data.CONTENT_URI)
                        .withValueBackReference(ContactsContract.Data.RAW_CONTACT_ID, 0)
                        .withValue(ContactsContract.Data.MIMETYPE, ContactsContract.CommonDataKinds.StructuredName.CONTENT_ITEM_TYPE)
                        .withValue(ContactsContract.CommonDataKinds.StructuredName.DISPLAY_NAME, info.getName())
                        .build()
        );

        // Phone Numbers
        for (ContactInfo.PhoneInfo phoneInfo : info.getPhonesList()) {
            ops.add(
                    ContentProviderOperation.newInsert(ContactsContract.Data.CONTENT_URI)
                            .withValueBackReference(ContactsContract.Data.RAW_CONTACT_ID, 0)
                            .withValue(ContactsContract.Data.MIMETYPE, ContactsContract.CommonDataKinds.Phone.CONTENT_ITEM_TYPE)
                            .withValue(ContactsContract.CommonDataKinds.Phone.NUMBER, phoneInfo.getNumber())
                            .withValue(ContactsContract.CommonDataKinds.Phone.TYPE, getPhoneTypeInt(phoneInfo.getType()))
                            .build()
            );
        }

        // Emails
        for (ContactInfo.EmailInfo emailInfo : info.getEmailsList()) {
            ops.add(
                    ContentProviderOperation.newInsert(ContactsContract.Data.CONTENT_URI)
                            .withValueBackReference(ContactsContract.Data.RAW_CONTACT_ID, 0)
                            .withValue(ContactsContract.Data.MIMETYPE, ContactsContract.CommonDataKinds.Email.CONTENT_ITEM_TYPE)
                            .withValue(ContactsContract.CommonDataKinds.Email.DATA, emailInfo.getEmail())
                            .withValue(ContactsContract.CommonDataKinds.Email.TYPE, getEmailTypeInt(emailInfo.getType()))
                            .build()
            );
        }
        try {
            context.getContentResolver().applyBatch(ContactsContract.AUTHORITY, ops);
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }
}
