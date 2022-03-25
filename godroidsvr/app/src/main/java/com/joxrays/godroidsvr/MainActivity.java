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

package com.joxrays.godroidsvr;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.appcompat.widget.AppCompatEditText;
import androidx.appcompat.widget.AppCompatTextView;
import androidx.appcompat.widget.SwitchCompat;
import androidx.constraintlayout.widget.ConstraintLayout;
import androidx.core.app.ActivityCompat;
import androidx.work.Data;
import androidx.work.ExistingPeriodicWorkPolicy;
import androidx.work.PeriodicWorkRequest;
import androidx.work.WorkManager;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.pm.PackageManager;
import android.location.LocationManager;
import android.net.wifi.WifiManager;
import android.opengl.GLSurfaceView;
import android.os.Bundle;
import android.text.Editable;
import android.view.KeyEvent;
import android.view.WindowManager;
import android.widget.Toast;

import com.joxrays.godroidsvr.receiver.ScanWifiResultReceiver;
import com.joxrays.godroidsvr.service.MainNotificationService;
import com.joxrays.godroidsvr.service.ImageServerService;
import com.joxrays.godroidsvr.singleton.ActivitySingleton;
import com.joxrays.godroidsvr.singleton.GPUInfoSingleton;
import com.joxrays.godroidsvr.util.CommonUtil;
import com.joxrays.godroidsvr.base.MyLocationListener;
import com.joxrays.godroidsvr.util.NetworkUtil;
import com.joxrays.godroidsvr.util.PermissionUtil;
import com.joxrays.godroidsvr.service.MediaProjectionService;
import com.joxrays.godroidsvr.service.RpcServerWorker;
import com.joxrays.godroidsvr.util.LocationUtil;

import java.util.List;
import java.util.concurrent.TimeUnit;

import javax.microedition.khronos.egl.EGLConfig;
import javax.microedition.khronos.opengles.GL10;

public class MainActivity extends AppCompatActivity {

    public static final int SCREEN_CAP_REQ_CODE = 0x100;
    public static final int SCREEN_CORD_REQ_CODE = 0x200;

    public static final int RPC_SERVER_PORT = 9999;
    public static final int WS_SERVER_PORT = 10000;

    private boolean isStartScreenCapService;
    private boolean isStartScreenCordService;
    private boolean isDestroy;
    private long expiredTime;
    private int rpc_port;

    private ConstraintLayout constraintLayout;
    private GLSurfaceView surfaceView;

    private SwitchCompat switch_rpc;
    private AppCompatEditText edit_text_port;

    private MyLocationListener locationListener;
    private ScanWifiResultReceiver resultReceiver;

    private void registerLocationListener() {
        List<String> providers = LocationUtil.getAvailableProvides(this);
        if (locationListener == null && providers.contains(LocationManager.NETWORK_PROVIDER)) {
            locationListener = new MyLocationListener();
            if (ActivityCompat.checkSelfPermission(this,
                    android.Manifest.permission.ACCESS_FINE_LOCATION) != PackageManager.PERMISSION_GRANTED
                    && ActivityCompat.checkSelfPermission(this, android.Manifest.permission.ACCESS_COARSE_LOCATION)
                    != PackageManager.PERMISSION_GRANTED) {
                return;
            }
            CommonUtil.getLocationManager(this).requestLocationUpdates(LocationManager.NETWORK_PROVIDER,
                    5000, 1, locationListener);
        }
    }

    private void unregisterLocationListener() {
        if (locationListener != null) {
            CommonUtil.getLocationManager(this).removeUpdates(locationListener);
        }
    }

    private void registerWifiResultReceiver() {
        resultReceiver = new ScanWifiResultReceiver();
        IntentFilter intentFilter = new IntentFilter(WifiManager.SCAN_RESULTS_AVAILABLE_ACTION);
        registerReceiver(resultReceiver, intentFilter);
    }

    private void unregisterWifiResultReceiver() {
        if (resultReceiver != null)
            unregisterReceiver(resultReceiver);
    }

    // start rpc server
    public void startBackgroundRpcWorkerService() {
        Data data = new Data.Builder().putInt("PORT", rpc_port).build();
        PeriodicWorkRequest request = new PeriodicWorkRequest.Builder(
                RpcServerWorker.class,
                15, TimeUnit.MINUTES,
                15, TimeUnit.MINUTES)
                .setInputData(data)
                .build();

        WorkManager.getInstance(this)
                .enqueueUniquePeriodicWork(RpcServerWorker.TAG, ExistingPeriodicWorkPolicy.REPLACE, request);

    }

    // stop rpc server
    public void stopBackgroundRpcWorkerService() {
        WorkManager.getInstance(getBaseContext()).cancelUniqueWork(RpcServerWorker.TAG);
    }

    // start background service for media projection
    // actually the background service need to start a foreground service for media projection
    // the foreground service is used for capturing screen image data,
    // and the background service is used for starting a websocket server
    public void startServiceBy(Intent data, int code) {
        Intent intent = new Intent(this, ImageServerService.class);
        intent.putExtra("data", data);
        intent.putExtra("resultCode", code);
        intent.putExtra("port", rpc_port + 1);
        startService(intent);
    }

    // start foreground service for media projection
    public void startForegroundServiceBy(Intent data, int code) {
        Intent intent = new Intent(this, MediaProjectionService.class);
        intent.putExtra("data", data);
        intent.putExtra("resultCode", code);
        startForegroundService(intent);
    }

    // stop main notification service
    public void stopMainNotification() {
        Intent intent = new Intent(this, MainNotificationService.class);
        stopService(intent);
    }

    // start screen mirroring service
    public void startScreenCapService() {
        startActivityForResult(CommonUtil.getMediaProjectionManager(this).createScreenCaptureIntent(), SCREEN_CAP_REQ_CODE);
    }

    // stop screen mirroring service
    public void stopScreenCapService() {
        if (!isStartScreenCapService) return;
        isStartScreenCapService = false;
        Intent intent = new Intent(this, ImageServerService.class);
        stopService(intent);
    }

    // start screen record service
    public void startScreenCordService() {
        startActivityForResult(CommonUtil.getMediaProjectionManager(this).createScreenCaptureIntent(), SCREEN_CORD_REQ_CODE);
    }

    // stop screen record service
    public void stopScreenCordService() {
        if (!isStartScreenCordService) return;
        isStartScreenCordService = false;
        Intent intent = new Intent(this, MediaProjectionService.class);
        stopService(intent);
    }

    @SuppressLint("WrongConstant")
    @Override
    protected void onActivityResult(int requestCode, int resultCode, @Nullable Intent data) {
        super.onActivityResult(requestCode, resultCode, data);

        if (requestCode == SCREEN_CAP_REQ_CODE && resultCode == RESULT_OK) {
            stopScreenCordService();
            startServiceBy(data, resultCode);
            isStartScreenCapService = true;
        }

        if (requestCode == SCREEN_CORD_REQ_CODE && resultCode == RESULT_OK) {
            stopScreenCapService();
            startForegroundServiceBy(data, resultCode);
            isStartScreenCordService = true;
        }
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        exitApp();
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        switch_rpc = findViewById(R.id.switch_rpc_server);
        edit_text_port = findViewById(R.id.edit_text_port);

        switch_rpc.setOnCheckedChangeListener((bv, is) -> {

            if (is) {
                Editable editable = edit_text_port.getText();
                rpc_port = RPC_SERVER_PORT;
                if (editable != null) {
                    String value = editable.toString();
                    if (!value.trim().isEmpty())
                        try {
                            rpc_port = Integer.parseInt(value);
                            if (rpc_port <= 0 || rpc_port >= 65536) {
                                switch_rpc.setChecked(false);
                                return;
                            }
                        } catch (Exception ex) {
                            switch_rpc.setChecked(false);
                            return;
                        }
                }

                PermissionUtil.requireNeedPermissions(this);
                startBackgroundRpcWorkerService();
            } else {
                stopBackgroundRpcWorkerService();
                stopScreenCapService();
                stopScreenCordService();
            }
        });

        findViewById(R.id.btn_exit).setOnClickListener(v -> {
            exitApp();
        });

        // initialize
        ActivitySingleton.getInstance().setActivity(this);
        getWindow().addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);
        startForegroundService(new Intent(this, MainNotificationService.class));
        initSurfaceView();
        registerLocationListener();
        registerWifiResultReceiver();

        WifiManager wm = CommonUtil.getWifiManager(this);
        if (wm != null) {
            android.net.wifi.WifiInfo info = wm.getConnectionInfo();
            if (info != null) {
                int value = info.getIpAddress();
                if (value != 0) {
                    String address = NetworkUtil.getAddress(info.getIpAddress());
                    String text = "WIFI address: " + address;
                    ((AppCompatTextView) findViewById(R.id.txt_message)).setText(text);
                }
            }
        } else {
            ((AppCompatTextView) findViewById(R.id.txt_message)).setText("");
        }

    }

    public void initSurfaceView() {
        constraintLayout = findViewById(R.id.constraint_layout);
        surfaceView = new GLSurfaceView(this);
        surfaceView.setRenderer(new GLSurfaceView.Renderer() {
            @Override
            public void onSurfaceCreated(GL10 gl, EGLConfig config) {
                GPUInfoSingleton.getInstance().setInfo(
                        gl.glGetString(GL10.GL_VENDOR),
                        gl.glGetString(GL10.GL_RENDERER),
                        gl.glGetString(GL10.GL_VERSION));
                runOnUiThread(() -> constraintLayout.removeView(surfaceView));
            }

            @Override
            public void onSurfaceChanged(GL10 gl, int width, int height) {
            }

            @Override
            public void onDrawFrame(GL10 gl) {
            }
        });
        constraintLayout.addView(surfaceView);
    }

    private void exitApp() {
        if (isDestroy) {
            return;
        }
        isDestroy = true;

        CommonUtil.getNotificationManager(this).cancelAll();
        unregisterLocationListener();
        unregisterWifiResultReceiver();

        stopBackgroundRpcWorkerService();
        stopScreenCordService();
        stopScreenCapService();
        stopMainNotification();

        finishAffinity();
        finish();
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        if (keyCode == KeyEvent.KEYCODE_BACK) {
            if (System.currentTimeMillis() - expiredTime >= 2000) {
                expiredTime = System.currentTimeMillis();
                Toast.makeText(this, "Press back again to exit", Toast.LENGTH_SHORT).show();
                return false;
            } else {
                exitApp();
            }
        }
        return super.onKeyDown(keyCode, event);
    }
}