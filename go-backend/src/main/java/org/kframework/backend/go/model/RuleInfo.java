package org.kframework.backend.go.model;

public class RuleInfo {

    private final boolean alwaysMatches;

    public RuleInfo(boolean alwaysMatches) {
        this.alwaysMatches = alwaysMatches;
    }

    public boolean alwaysMatches() {
        return alwaysMatches;
    }
}
