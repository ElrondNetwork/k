package org.kframework.backend.go;

import org.kframework.kore.KVariable;
import org.kframework.kore.VisitK;

public class AccumulateVars extends VisitK {

    private final VarInfo varInfo = new VarInfo();

    public AccumulateVars() {
    }

    public VarInfo vars() {
        return varInfo;
    }

    @Override
    public void apply(KVariable k) {
        String varName = GoStringUtil.variableName(k.name());
        if (varInfo.containsVar(k)) {
            assert varInfo.getVarName(k).equals(varName);
        }
        varInfo.putVar(k, varName);
    }
}
