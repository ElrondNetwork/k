// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.model;

public class RuleCounter {

    private int nextRuleIndex = 0;

    public int consumeRuleIndex() {
        return ++nextRuleIndex;
    }

}
