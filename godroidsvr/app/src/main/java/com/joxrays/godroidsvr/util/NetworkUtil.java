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

import android.annotation.SuppressLint;
import android.content.Context;
import android.net.ConnectivityManager;
import android.net.DhcpInfo;
import android.net.LinkAddress;
import android.net.LinkProperties;
import android.net.Network;
import android.net.NetworkCapabilities;
import android.net.ProxyInfo;
import android.net.wifi.WifiManager;
import android.os.Build;
import android.util.Pair;

import androidx.annotation.RequiresApi;

import com.google.gson.Gson;
import com.joxrays.godroidsvr.base.PublicAddressInfo;
import com.joxrays.godroidsvr.base.TimeoutTask;
import com.joxrays.godroidsvr.message.DetailActiveNetworkInfo;
import com.joxrays.godroidsvr.message.InetAddr;
import com.joxrays.godroidsvr.message.NetInterfaceInfo;
import com.joxrays.godroidsvr.message.NetInterfaceInfoList;
import com.joxrays.godroidsvr.message.DetailWifiInfo;
import com.joxrays.godroidsvr.message.PublicNetworkInfo;

import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.InterfaceAddress;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.List;

public final class NetworkUtil {

    private NetworkUtil() {

    }

    public static String getWifiStatus(Context context) {
        WifiManager wm = CommonUtil.getWifiManager(context);
        if (wm != null) {
            switch (wm.getWifiState()) {
                case WifiManager.WIFI_STATE_ENABLED: // Enable
                    return "Enabled";
                case WifiManager.WIFI_STATE_DISABLED: // Disabled
                    return "Disabled";
                case WifiManager.WIFI_STATE_ENABLING: // Enabling
                    return "Enabling";
                case WifiManager.WIFI_STATE_DISABLING: // Disabling
                    return "Disabling";
            }
        }
        return "Unknown";
    }

    public static String getAddress(int address) {
        return (address & 0xFF) + "." +
                ((address >> 8) & 0xFF) + "." +
                ((address >> 16) & 0xFF) + "." +
                ((address >> 24) & 0xFf);
    }

    @RequiresApi(api = Build.VERSION_CODES.R)
    @SuppressLint({"MissingPermission", "HardwareIds"})
    public static Pair<DetailWifiInfo, Exception> getWifiInfo(Context context) {
        WifiManager wm = CommonUtil.getWifiManager(context);
        if (wm == null) return Pair.create(null, ErrorExceptionUtil.ErrorNotFoundWifiInfo);
        android.net.wifi.WifiInfo info = wm.getConnectionInfo();
        if (info == null) return Pair.create(null, ErrorExceptionUtil.ErrorNotFoundWifiInfo);
        DhcpInfo dhcpInfo = wm.getDhcpInfo();
        DetailWifiInfo.Builder builder = DetailWifiInfo.newBuilder();
        String ssid = info.getSSID();
        if (ssid == null) ssid = "";
        ssid = ssid.replaceAll("^\"|\"$", "");

        builder.setSsid(ssid)
                .setBssid(info.getBSSID() != null ? info.getBSSID() : "")
                .setNetworkId(info.getNetworkId())
                .setMac(info.getMacAddress() != null ? info.getMacAddress() : "")
                .setSignal(CommonUtil.getSignalStrengthFromLevel((info.getRssi())))
                .setLinkSpeed(info.getLinkSpeed())     // Mbps
                .setTxSpeed(info.getTxLinkSpeedMbps()) // Mbps
                .setRxSpeed(info.getRxLinkSpeedMbps()) // Mbps
                .setFrequency(info.getFrequency())
                .setIp(getAddress(dhcpInfo.ipAddress))
                .setGateway(getAddress(dhcpInfo.gateway))
                .setNetmask(getAddress(dhcpInfo.netmask))
                .setServerAddress(getAddress(dhcpInfo.serverAddress))
                .setDns1(getAddress(dhcpInfo.dns1))
                .setDns2(getAddress(dhcpInfo.dns2))
                .setStatus(getWifiStatus(context));
        return Pair.create(builder.build(), null);
    }

    private static String getMacAddress(byte[] mac) {
        if (mac == null) return "";
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < mac.length; i++) {
            sb.append(String.format("%02X", mac[i]));
            if (i + 1 != mac.length) {
                sb.append(":");
            }
        }
        return sb.toString();
    }

    public static Pair<List<NetInterfaceInfo>, Exception> getAllNetInterfaceInfo() {
        List<NetInterfaceInfo> list = new ArrayList<>();
        try {
            // get all network interfaces
            for (Enumeration<NetworkInterface> interfaces = NetworkInterface.getNetworkInterfaces(); interfaces.hasMoreElements(); ) {
                NetworkInterface network = interfaces.nextElement();
                NetInterfaceInfo.Builder info = NetInterfaceInfo.newBuilder()
                        .setUp(network.isUp())
                        .setMtu(network.getMTU())
                        .setMacAddress(getMacAddress(network.getHardwareAddress()))
                        .setName(network.getName());

                List<InetAddr> inetAddrs = new ArrayList<>();

                for (InterfaceAddress address : network.getInterfaceAddresses()) {
                    InetAddr.Builder builder = InetAddr.newBuilder();
                    builder.setIp(address.getAddress().getHostAddress()).setPrefixLength(address.getNetworkPrefixLength());

                    // ipv4
                    if (address.getAddress() instanceof Inet4Address) {
                        builder.setBroadcastIp(address.getBroadcast() != null
                                ? address.getBroadcast().getHostAddress() : "")
                                .setIpv4(true);
                    } else {
                        // ipv6 have not broadcast address, but it has multicast address
                        builder.setBroadcastIp("").setIpv4(false);
                    }
                    inetAddrs.add(builder.build());
                }
                list.add(info.addAllInetAddrs(inetAddrs).build());
            }
        } catch (SocketException ex) {
            return Pair.create(null, ex);
        }
        return Pair.create(list, null);
    }

    public static boolean hasActiveNetwork(Context context) {
        ConnectivityManager cm = CommonUtil.getConnectivityManager(context);
        Network network = cm.getActiveNetwork();
        if (network == null) return false;
        NetworkCapabilities nc = cm.getNetworkCapabilities(network);
        return nc != null && nc.hasCapability(NetworkCapabilities.NET_CAPABILITY_VALIDATED);
    }

    @RequiresApi(api = Build.VERSION_CODES.R)
    public static List<DetailActiveNetworkInfo> getActivityNetworkDetailInfo(Context context) {
        ConnectivityManager cm = CommonUtil.getConnectivityManager(context);
        List<DetailActiveNetworkInfo> list = new ArrayList<>();

        Network[] networks = cm.getAllNetworks();
        if (networks == null) {
            return null;
        }

        for (Network network : networks) {
            DetailActiveNetworkInfo.Builder builder = DetailActiveNetworkInfo.newBuilder();
            LinkProperties lp = cm.getLinkProperties(network);
            NetworkCapabilities nc = cm.getNetworkCapabilities(network);

            builder.setName(lp.getInterfaceName());
            builder.setMtu(lp.getMtu());

            ProxyInfo pi = lp.getHttpProxy();
            if (pi != null) {
                builder.setProxy(com.joxrays.godroidsvr.message.ProxyInfo.newBuilder()
                        .setHost(pi.getHost())
                        .setPac(pi.getPacFileUrl().toString())
                        .setPort(pi.getPort())
                        .build());
            }

            List<InetAddress> addresses = lp.getDnsServers();
            if (addresses != null) {
                String dns = "";
                for (InetAddress address : addresses) {
                    if (address instanceof Inet4Address) {
                        dns = address.getHostAddress();
                        break;
                    }
                }
                builder.setDns(dns);
            }

            List<LinkAddress> linkAddresses = lp.getLinkAddresses();
            if (linkAddresses != null) {
                String link = "";
                for (LinkAddress address : linkAddresses) {
                    if (address.getAddress() instanceof Inet4Address) {
                        link = address.getAddress().getHostAddress();
                        break;
                    }
                }
                builder.setIp(link);
            }

            String type = "Unknown";
            if (nc.hasTransport(NetworkCapabilities.TRANSPORT_WIFI)) {
                type = "WIFI";
            } else if (nc.hasTransport(NetworkCapabilities.TRANSPORT_CELLULAR)) {
                type = "Cellular";
            } else if (nc.hasTransport(NetworkCapabilities.TRANSPORT_VPN)) {
                type = "VPN";
            }

            builder.setType(type);
            builder.setSignal(CommonUtil.getSignalStrengthFromLevel(nc.getSignalStrength()));
            builder.setHasNetwork(nc.hasCapability(NetworkCapabilities.NET_CAPABILITY_VALIDATED));
            list.add(builder.build());
        }
        return list;
    }

    @RequiresApi(api = Build.VERSION_CODES.R)
    public static Pair<NetInterfaceInfoList, Exception> getNetworkInfo(Context context) {
        Pair<DetailWifiInfo, Exception> pair1 = getWifiInfo(context);
        if (pair1.second != null) {
            return Pair.create(null, pair1.second);
        }
        Pair<List<NetInterfaceInfo>, Exception> pair2 = getAllNetInterfaceInfo();
        if (pair2.second != null) {
            return Pair.create(null, pair2.second);
        }
        return Pair.create(NetInterfaceInfoList.newBuilder()
                .addAllValues(pair2.first)
                .build(), null);
    }

    public static Pair<PublicNetworkInfo, Exception> getPublicAddressInfo() {
        TimeoutTask<Void, Pair<String, Exception>> task = new TimeoutTask<>((param) -> {
            try {
                URL url = new URL("https://ipinfo.io/json");
                HttpURLConnection connection = (HttpURLConnection) url.openConnection();
                connection.setRequestMethod("GET");
                connection.setInstanceFollowRedirects(true);
                connection.connect();
                int code = connection.getResponseCode();
                if (code != 200)
                    return null;
                try (InputStream in = connection.getInputStream()) {
                    byte[] buffer = new byte[4096];
                    StringBuilder sb = new StringBuilder();
                    while (true) {
                        int n = in.read(buffer);
                        if (n <= 0) break;
                        sb.append(new String(buffer, 0, n));
                    }

                    return Pair.create(sb.toString(), null);
                } finally {
                    connection.disconnect();
                }
            } catch (Exception ex) {
                return Pair.create(null, ex);
            }
        });

        Pair<String, Exception> result = task.executeAndGet(4000, null);
        if (result == null) {
            return Pair.create(null, ErrorExceptionUtil.ErrorExecuteTimeout);
        }
        if (result.second != null) {
            return Pair.create(null, result.second);
        }
        PublicAddressInfo info = new Gson().fromJson(result.first, PublicAddressInfo.class);
        if (info == null) {
            return Pair.create(null, ErrorExceptionUtil.ErrorParseJson);
        }
        PublicNetworkInfo.Builder builder = PublicNetworkInfo.newBuilder();
        builder.setIp(info.ip != null ? info.ip : "")
                .setHostname(info.hostname != null ? info.hostname : "")
                .setCity(info.city != null ? info.city : "")
                .setRegion(info.region != null ? info.region : "")
                .setCountry(info.country != null ? info.country : "")
                .setLocation(info.loc != null ? info.loc : "")
                .setIsp(info.org != null ? info.org : "")
                .setTimezone(info.timezone != null ? info.timezone : "");
        return Pair.create(builder.build(), null);
    }
}
