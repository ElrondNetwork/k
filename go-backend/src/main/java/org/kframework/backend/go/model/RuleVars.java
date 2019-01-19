package org.kframework.backend.go.model;

import org.kframework.kore.KLabel;
import org.kframework.kore.KVariable;

import java.util.HashMap;
import java.util.Map;

public class RuleVars {
    private final Map<KVariable, String> kVarToName = new HashMap<>();
    public final Map<String, KLabel> listVars = new HashMap<>(); // TEMP

    public RuleVars() {
    }

    public void putVar(KVariable kv, String varName) {
        kVarToName.put(kv, varName);
    }

    public boolean containsVar(KVariable kv) {
        return kVarToName.containsKey(kv);
    }

    public String getVarName(KVariable kv) {
        return kVarToName.get(kv);
    }

}
