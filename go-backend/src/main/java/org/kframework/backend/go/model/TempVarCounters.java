package org.kframework.backend.go.model;

public class TempVarCounters {

    private int evalVarCounter = 0;

    public TempVarCounters() {
    }

    public int consumeEvalVarIndex() {
        return evalVarCounter++;
    }
}
