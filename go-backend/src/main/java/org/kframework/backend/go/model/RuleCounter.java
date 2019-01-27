package org.kframework.backend.go.model;

public class RuleCounter {

    private int nextRuleIndex = 0;

    public int consumeRuleIndex() {
        return ++nextRuleIndex;
    }

}
