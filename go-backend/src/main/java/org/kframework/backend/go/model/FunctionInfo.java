package org.kframework.backend.go.model;

import org.kframework.kore.KLabel;

public class FunctionInfo {

    public final KLabel label;
    public final String goName;
    public final boolean isMemo;
    public final FunctionParams arguments;

    private FunctionInfo(KLabel label, String goName, boolean isMemo, int arity) {
        this.label = label;
        this.goName = goName;
        this.isMemo = isMemo;
        this.arguments = new FunctionParams(arity);
    }

    public static FunctionInfo definitionFunctionInfo(KLabel label, String goName, boolean isMemo, int arity) {
        return new FunctionInfo(label, goName, isMemo, arity);
    }

    public static FunctionInfo systemFunctionInfo(String name, int arity) {
        return new FunctionInfo(null, name, false, arity);
    }

    public boolean isSystemFunction() {
        return goName == null;
    }
}
