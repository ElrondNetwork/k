// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.model;

public class RuleInfo {

    private final boolean alwaysMatches;
    public int nrVars;
    public int nrBoolVars;

    public RuleInfo(boolean alwaysMatches, int nrVars, int nrBoolVars) {
        this.alwaysMatches = alwaysMatches;
        this.nrVars = nrVars;
        this.nrBoolVars = nrBoolVars;
    }

    public boolean alwaysMatches() {
        return alwaysMatches;
    }
}
