package org.kframework.backend.go.processors;

import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.kore.KVariable;
import org.kframework.kore.VisitK;

public class AccumulateRuleVars extends VisitK {

    private final RuleVars varInfo = new RuleVars();

    public AccumulateRuleVars() {
    }

    public RuleVars vars() {
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
