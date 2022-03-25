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

import android.util.Pair;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.joxrays.godroidsvr.message.FileInfo;
import com.joxrays.godroidsvr.message.FileInfoList;

import java.io.File;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.CopyOption;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.nio.file.StandardOpenOption;
import java.util.ArrayList;
import java.util.List;

public final class FilesUtil {

    private FilesUtil() {

    }

    public static class BaseTreeFileInfo {
        String id;
        String parentId;
        String text;
        String type;
        String data;
        List<BaseTreeFileInfo> children;

        public BaseTreeFileInfo(String id, String parentId, String text, String type, String data, List<BaseTreeFileInfo> children) {
            this.id = id;
            this.parentId = parentId;
            this.text = text;
            this.type = type;
            this.data = data;
            this.children = children;
        }
    }

    private static List<BaseTreeFileInfo> getBaseTreeFileInfo0(File dir, String parentId, boolean all, int level, int maxLevel) {
        File[] files = dir.listFiles();

        if (files == null) return null;
        List<BaseTreeFileInfo> list = new ArrayList<>();

        for (int i = 0; i < files.length; i++) {
            File file = files[i];
            if (!all && file.isHidden()) {
                continue;
            }
            String id = parentId + "_" + i;
            BaseTreeFileInfo info = new BaseTreeFileInfo(id, parentId, file.getName(), "file", file.getAbsolutePath(), null);

            if (file.isDirectory()) {
                info.type = "folder";
                if (level + 1 < maxLevel) {
                    info.children = getBaseTreeFileInfo0(file, id, all, level + 1, maxLevel);
                }
            }
            list.add(info);
        }
        return list;
    }

    public static BaseTreeFileInfo getBaseTreeFileInfo(File dir, String pid, boolean all, int maxLevel) {
        String id = pid + "_" + 0;
        List<BaseTreeFileInfo> list = getBaseTreeFileInfo0(dir, id, all, 0, maxLevel);
        return new BaseTreeFileInfo(id, pid, dir.getName(), "folder", dir.getAbsolutePath(), list);
    }

    public static String getBaseTreeFileInfoJson(File dir, String pid, boolean all, int maxLevel) {
        if (maxLevel > 2) maxLevel = 2;
        Gson gson = new GsonBuilder().serializeNulls().create();
        return gson.toJson(getBaseTreeFileInfo(dir, pid, all, maxLevel));
    }

    public static Pair<FileInfo, Exception> getFileInfo(Path path) {
        try {
            return Pair.create(FileInfo.newBuilder()
                    .setName(path.toString())
                    .setSize(Files.size(path))
                    .setReadable(Files.isReadable(path))
                    .setWritable(Files.isWritable(path))
                    .setExecutable(Files.isExecutable(path))
                    .setDir(Files.isDirectory(path))
                    .setLastModifiedTime(Files.getLastModifiedTime(path).toMillis())
                    .setOwner(Files.getOwner(path).getName())
                    .build(), null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static Pair<FileInfoList, Exception> listDir(String dirName, boolean all) {
        FileInfoList.Builder builder = FileInfoList.newBuilder();
        File dir = new File(dirName);
        File[] files = dir.listFiles();
        if (files == null)
            return Pair.create(null, ErrorExceptionUtil.ErrorCannotListDir);
        for (File f : files) {
            if (!all && f.isHidden())
                continue;
            Pair<FileInfo, Exception> pair = getFileInfo(f.toPath());
            if (pair.second != null) continue;
            builder.addValues(pair.first);
        }
        return Pair.create(builder.build(), null);
    }

    public static long size(File file) {
        return file.length();
    }

    public static Pair<Boolean, Exception> createFile(File file) {
        try {
            return Pair.create(file.createNewFile(), null);
        } catch (IOException ex) {
            return Pair.create(false, ex);
        }
    }

    public static boolean exist(File file) {
        return file.exists();
    }

    public static boolean deleteFile(File file) {
        return file.delete();
    }

    public static boolean mkDir(File file) {
        return file.mkdirs();
    }

    private static void rmDir0(File file, long[] status) {
        File[] list = file.listFiles();
        if (list == null) {
            return;
        }
        for (File f : list) {
            if (f.isFile()) {
                status[1] += size(f);
                status[0] = deleteFile(f) ? 1 : 0;
            } else if (f.isDirectory()) {
                rmDir0(f, status);
            }
        }
        status[0] = deleteFile(file) ? 1 : 0;
    }

    public static Pair<Boolean, Long> rmDir(File file) {
        if (file.isFile()) {
            return Pair.create(deleteFile(file), size(file));
        } else if (file.isDirectory()) {
            long[] status = {0, 0};
            rmDir0(file, status);
            return Pair.create(status[0] == 1, status[1]);
        }
        return Pair.create(false, 0L);
    }

    public static Pair<Path, Exception> move(File src, File dest) {
        try {
            Path path = Files.move(src.toPath(), dest.toPath(), StandardCopyOption.REPLACE_EXISTING, StandardCopyOption.ATOMIC_MOVE);
            return Pair.create(path, null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    private static void copy0(File src, File dest, CopyOption co) throws IOException {
        File[] list = src.listFiles();
        if (list == null) return;
        // copy parent directory
        Files.copy(src.toPath(), dest.toPath(), co);
        for (File file : list) {
            Path d = Paths.get(dest.getAbsolutePath(), file.getName());
            if (file.isDirectory()) {
                Files.copy(file.toPath(), d, co);
                copy0(file, d.toFile(), co);
            } else if (file.isFile()) {
                Files.copy(file.toPath(), d, co);
            }
        }
    }

    public static Pair<Path, Exception> copy(File src, File dest) {
        try {
            // nothing to do if file existed
            if (dest.exists()) {
                return Pair.create(dest.toPath(), null);
            }
            copy0(src, dest, StandardCopyOption.REPLACE_EXISTING);
            return Pair.create(dest.toPath(), null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static boolean rename(File src, File dest) {
        return src.renameTo(dest);
    }

    public static Pair<String, Exception> readText(File file) {
        // limit file size
        if (file.length() >= 1024 * 1024) {
            return Pair.create("", ErrorExceptionUtil.ErrorFileLarge);
        }
        try {
            String text = new String(Files.readAllBytes(file.toPath()));
            return Pair.create(text, null);
        } catch (Exception ex) {
            return Pair.create("", ex);
        }
    }

    public static Pair<Path, Exception> writeText(File file, String text, boolean append) {
        try {
            Path path = file.toPath();
            if (append) {
                Files.write(path, text.getBytes(StandardCharsets.UTF_8),
                        StandardOpenOption.CREATE, StandardOpenOption.APPEND,
                        StandardOpenOption.SYNC);
            } else {
                Files.write(path, text.getBytes(StandardCharsets.UTF_8),
                        StandardOpenOption.CREATE, StandardOpenOption.WRITE,
                        StandardOpenOption.TRUNCATE_EXISTING, StandardOpenOption.SYNC);
            }
            return Pair.create(path, null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }
}
