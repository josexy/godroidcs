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

package com.joxrays.godroidsvr.base;

import android.location.Location;
import android.location.LocationListener;
import android.util.Log;

import androidx.annotation.NonNull;

import com.joxrays.godroidsvr.singleton.LocationSingleton;

public class MyLocationListener implements LocationListener {

    @Override
    public void onProviderEnabled(@NonNull String provider) {
        Log.d("", provider + " enabled");
    }

    @Override
    public void onProviderDisabled(@NonNull String provider) {
        Log.d("", provider + " disabled");
    }

    @Override
    public void onLocationChanged(Location loc) {
        Log.d("", "" + loc);
        LocationSingleton.getInstance().setLocation(loc);
    }
}