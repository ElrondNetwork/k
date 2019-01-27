package org.kframework.backend.go.model;

public class RuleInfo {

    private final boolean topLevelIf;

    public RuleInfo(boolean topLevelIf) {
        this.topLevelIf = topLevelIf;
    }

    public boolean isTopLevelIf() {
        return topLevelIf;
    }
}
