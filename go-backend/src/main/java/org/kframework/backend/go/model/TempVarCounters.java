// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.model;

public class TempVarCounters {

    private int lhsVarCounter = 0;
    private int evalVarCounter = 0;

    public TempVarCounters() {
    }

    public int consumeLhsVarIndex() {
        return lhsVarCounter++;
    }

    public int consumeEvalVarIndex() {
        return evalVarCounter++;
    }
}
