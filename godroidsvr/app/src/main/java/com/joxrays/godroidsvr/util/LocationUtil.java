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

import android.Manifest;
import android.content.Context;
import android.content.pm.PackageManager;
import android.location.Address;
import android.location.Geocoder;
import android.location.Location;
import android.location.LocationManager;
import android.util.Pair;

import androidx.core.app.ActivityCompat;

import com.joxrays.godroidsvr.message.LocationInfo;
import com.joxrays.godroidsvr.singleton.LocationSingleton;

import java.io.IOException;
import java.util.List;
import java.util.Locale;

public final class LocationUtil {

    private LocationUtil() {

    }

    public static List<String> getAvailableProvides(Context context) {
        return CommonUtil.getLocationManager(context).getProviders(true);
    }

    public static Pair<LocationInfo, Exception> getAddressInfo(Context context, String provider) {
        LocationManager lm = CommonUtil.getLocationManager(context);
        if (ActivityCompat.checkSelfPermission(context, Manifest.permission.ACCESS_FINE_LOCATION) != PackageManager.PERMISSION_GRANTED
                && ActivityCompat.checkSelfPermission(context, Manifest.permission.ACCESS_COARSE_LOCATION)
                != PackageManager.PERMISSION_GRANTED) {
            return Pair.create(null, ErrorExceptionUtil.ErrorCannotGetLocation);
        }
        Location loc = lm.getLastKnownLocation(provider);
        Location updatedLoc = LocationSingleton.getInstance().getLocation();

        if (updatedLoc != null) {
            loc = updatedLoc;
        }
        return getAddressInfoByLocation(context, loc);
    }

    public static Pair<LocationInfo, Exception> getAddressInfoByLocation(Context context, Location location) {
        if (location == null) return Pair.create(null, ErrorExceptionUtil.ErrorCannotGetLocation);

        Pair<Address, Exception> pair = getAddressByLocation(context, location);
        if (pair.second != null) {
            return Pair.create(null, pair.second);
        }
        if (pair.first == null) {
            return Pair.create(null, ErrorExceptionUtil.ErrorCannotGetLocation);
        }
        LocationInfo.Builder builder = LocationInfo.newBuilder();
        Address address = pair.first;
        String line = "";
        if (address.getMaxAddressLineIndex() >= 0) {
            line = address.getAddressLine(0);
        }
        String countryName = address.getCountryName();
        String countryCode = address.getCountryCode();
        String adminArea = address.getAdminArea();
        String locality = address.getLocality();
        String subLocality = address.getSubLocality();

        builder.setLongitude(address.getLongitude())
                .setLatitude(address.getLatitude())
                .setCountryName(countryName != null ? countryName : "")
                .setCountryCode(countryCode != null ? countryCode : "")
                .setAdminArea(adminArea != null ? adminArea : "")
                .setLocality(locality != null ? locality : "")
                .setSubLocality(subLocality != null ? subLocality : "")
                .setAddressLine(line);
        return Pair.create(builder.build(), null);
    }

    public static Pair<Address, Exception> getAddressByLocation(Context context, Location location) {
        Geocoder geocoder = new Geocoder(context, Locale.getDefault());
        List<Address> addresses;
        try {
            addresses = geocoder.getFromLocation(location.getLatitude(), location.getLongitude(), 1);
            Address address = null;
            if (addresses != null && addresses.size() > 0) {
                address = addresses.get(0);
            }
            return Pair.create(address, null);

        } catch (IOException ex) {
            return Pair.create(null, ex);
        }
    }
}
