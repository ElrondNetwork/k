package org.kframework.backend.go.model;

import org.kframework.backend.go.codegen.GoBuiltin;

public class FunctionHookName {

    private final String originalName;
    private final String namespace;
    private final String goFuncName;
    private final String goHookObjName;

    public FunctionHookName(String hook) {
        this.originalName = hook;
        this.namespace = hook.substring(0, hook.indexOf('.'));
        this.goHookObjName = namespace.toLowerCase() + "Hooks";
        if (GoBuiltin.GO_FRIENDLY_HOOK_NAMES.containsKey(originalName)) {
            this.goFuncName = GoBuiltin.GO_FRIENDLY_HOOK_NAMES.get(originalName);
        } else {
            this.goFuncName = hook.substring(namespace.length() + 1);
        }
    }

    public String getOriginalName() {
        return originalName;
    }

    public String getNamespace() {
        return namespace;
    }

    public String getGoFuncName() {
        return goFuncName;
    }

    public String getGoHookObjName() {
        return goHookObjName;
    }

    public String getExternalGoPackageName() {
        return namespace.toLowerCase();
    }

    @Override
    public String toString() {
        return namespace + "." + goFuncName;
    }
}
