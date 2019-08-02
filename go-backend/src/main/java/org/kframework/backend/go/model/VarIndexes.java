package org.kframework.backend.go.model;

import org.kframework.kore.KVariable;

import java.util.HashMap;
import java.util.Map;

public class VarIndexes {
    private int nrVars = 0;
    private int nrBoolVars = 0;
    public final Map<KVariable, Integer> kvariableIndexes = new HashMap<>();
    public final Map<String, Integer> varNameToIndex = new HashMap<>();

    public VarIndexes() {
    }

    public String variableMVRef(String varName) {
        Integer index = varNameToIndex.get(varName);
        if (index == null) {
            index = nrVars;
            nrVars++;
            varNameToIndex.put(varName, index);
        }
        return "v[" + index + "]";
    }

    public String oneTimeVariableMVRef(String varName) {
        Integer index = nrVars;
        nrVars++;
        //return "mv[" + index + "/*" + varName + "*/]";
        return "v[" + index + "]";
    }

    public String kvariableMVRef(KVariable k) {
        Integer index = kvariableIndexes.get(k);
        if (index == null) {
            index = nrVars;
            nrVars++;
            kvariableIndexes.put(k, index);
        }
        return "v[" + index + " /*" + k.name() + "*/]";
    }

    public String evalBoolVarRef(String varName) {
        return "bv[" + (nrBoolVars++) + "]";
    }

    public int getNrVars() {
        return nrVars;
    }

    public int getNrBoolVars() {
        return nrBoolVars;
    }

}
