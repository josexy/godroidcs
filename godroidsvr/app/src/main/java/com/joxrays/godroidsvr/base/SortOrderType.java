package com.joxrays.godroidsvr.base;

public enum SortOrderType {

    Default(0, null),
    ASC(1, "ASC"),
    DESC(2, "DESC");

    public int type;
    public String name;

    SortOrderType(int type, String name) {
        this.type = type;
        this.name = name;
    }
}

