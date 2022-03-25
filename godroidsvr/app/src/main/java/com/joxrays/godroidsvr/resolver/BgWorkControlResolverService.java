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

package com.joxrays.godroidsvr.resolver;

import android.app.Activity;
import android.util.Pair;

import com.joxrays.godroidsvr.MainActivity;
import com.joxrays.godroidsvr.singleton.ActivitySingleton;
import com.joxrays.godroidsvr.singleton.FinalMediaProjectionSingleton;
import com.joxrays.godroidsvr.util.ClipboardUtil;
import com.joxrays.godroidsvr.util.ControllerUtil;
import com.joxrays.godroidsvr.message.Empty;
import com.joxrays.godroidsvr.message.Boolean;
import com.joxrays.godroidsvr.message.Integer;
import com.joxrays.godroidsvr.message.String;

import io.grpc.stub.StreamObserver;

public class BgWorkControlResolverService extends ControlResolverGrpc.ControlResolverImplBase {

    private final BgWorkBaseResolverGroup group;

    public BgWorkControlResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @Override
    public void getScreenBrightness(Empty request, StreamObserver<Integer> responseObserver) {
        int value = ControllerUtil.getScreenBrightness(this.group.getContext());
        responseObserver.onNext(Integer.newBuilder().setValue(value).build());
        responseObserver.onCompleted();
    }

    @Override
    public void setScreenBrightness(Integer request, StreamObserver<Empty> responseObserver) {
        Exception ex = ControllerUtil.setScreenBrightness(this.group.getContext(), request.getValue());
        if (ex != null) {
            responseObserver.onError(ex);
            return;
        }
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void getScreenBrightnessMode(Empty request, StreamObserver<Integer> responseObserver) {
        Pair<java.lang.Integer, Exception> pair = ControllerUtil.getScreenBrightnessMode(this.group.getContext());
        if (pair.second != null) {
            responseObserver.onError(pair.second);
            return;
        }
        responseObserver.onNext(Integer.newBuilder()
                .setValue(pair.first)
                .build());
        responseObserver.onCompleted();
    }

    // set screen brightness mode (0: manual, 1: auto)
    @Override
    public void setScreenBrightnessMode(Boolean request, StreamObserver<Empty> responseObserver) {
        ControllerUtil.setScreenBrightnessMode(this.group.getContext(), request.getValue() ? 1 : 0);
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void getClipboardText(Empty request, StreamObserver<String> responseObserver) {
        responseObserver.onNext(String.newBuilder()
                .setValue(ClipboardUtil.getText(this.group.getContext()))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void setClipboardText(String request, StreamObserver<Empty> responseObserver) {
        ClipboardUtil.setText(this.group.getContext(), request.getValue());
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void getVolume(Empty request, StreamObserver<Integer> responseObserver) {
        responseObserver.onNext(Integer.newBuilder()
                .setValue(ControllerUtil.getAudioMusicValue(this.group.getContext()))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void setVolume(Integer request, StreamObserver<Empty> responseObserver) {
        // show ui when set volume
        ControllerUtil.setAudioMusicValue(this.group.getContext(), request.getValue(), true);
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void increaseVolume(Empty request, StreamObserver<Empty> responseObserver) {
        // show ui when increase volume
        ControllerUtil.increaseAudioMusic(this.group.getContext(), true);
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void decreaseVolume(Empty request, StreamObserver<Empty> responseObserver) {
        // show ui when decrease volume
        ControllerUtil.decreaseAudioMusic(this.group.getContext(), true);
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    /* start screen mirroring. Use ImageReader */
    @Override
    public void startScreenCapture(Empty request, StreamObserver<Empty> responseObserver) {
        Activity activity = ActivitySingleton.getInstance().getActivity();
        if (activity instanceof MainActivity) {
            ((MainActivity) (activity)).startScreenCapService();
        }
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    /* stop screen mirroring */
    @Override
    public void stopScreenCapture(Empty request, StreamObserver<Empty> responseObserver) {
        Activity activity = ActivitySingleton.getInstance().getActivity();
        if (activity instanceof MainActivity) {
            ((MainActivity) (activity)).stopScreenCapService();
        }
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    /* start screen record and save to file. Use MediaRecorder */
    @Override
    public void startScreenRecord(Empty request, StreamObserver<Empty> responseObserver) {
        Activity activity = ActivitySingleton.getInstance().getActivity();
        if (activity instanceof MainActivity) {
            ((MainActivity) (activity)).startScreenCordService();
        }
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    /* stop screen record */
    @Override
    public void stopScreenRecord(Empty request, StreamObserver<String> responseObserver) {
        Activity activity = ActivitySingleton.getInstance().getActivity();
        if (activity instanceof MainActivity) {
            ((MainActivity) (activity)).stopScreenCordService();
        }

        // return the video file path to client finally
        java.lang.String imageFile = FinalMediaProjectionSingleton.getInstance().getImageFile();
        FinalMediaProjectionSingleton.getInstance().setImageFile(null);
        responseObserver.onNext(String.newBuilder().setValue(imageFile != null ? imageFile : "").build());
        responseObserver.onCompleted();
    }
}
