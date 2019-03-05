// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.gopackage;

public class GoPackage {

    private final String name;
    private String alias;
    private final String goPath;
    private final String relativePath;

    public GoPackage(String name, String goPath, String relativePath) {
        this.name = name;
        this.alias = name;
        this.goPath = goPath;
        this.relativePath = relativePath;
    }

    public String getName() {
        return name;
    }

    public String getAlias() {
        return alias;
    }

    public void setAlias(String alias) {
        this.alias = alias;
    }

    public String getGoPath() {
        return goPath;
    }

    public String getRelativePath() {
        return relativePath;
    }
}
