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

import android.os.AsyncTask;

import java.util.concurrent.TimeUnit;
import java.util.function.Function;

public class TimeoutTask<Param, Result> extends AsyncTask<Param, Void, Result> {

    private final Function<Param, Result> runnable;

    public TimeoutTask(Function<Param, Result> runnable) {
        this.runnable = runnable;
    }

    @Override
    protected Result doInBackground(Param[] params) {
        Param param = null;
        if (params != null && params.length > 0) {
            param = params[0];
        }
        return runnable.apply(param);
    }

    // execute task and get result in given timeout
    public Result executeAndGet(int timeout, Param[] params) {
        try {
            execute(params);
            return get(timeout, TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            return null;
        }
    }
}
