// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.processors;

import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.kore.KApply;
import org.kframework.kore.KVariable;
import org.kframework.kore.VisitK;

public class AccumulateRuleVars extends VisitK {

    private final GoNameProvider nameProvider;
    private final RuleVars varInfo = new RuleVars();

    public AccumulateRuleVars(GoNameProvider nameProvider) {
        this.nameProvider = nameProvider;
    }

    public RuleVars vars() {
        return varInfo;
    }

    @Override
    public void apply(KVariable k) {
        String varName = nameProvider.ruleVariableName(k.name());
        if (varInfo.containsVar(k)) {
            assert varInfo.getVarName(k).equals(varName);
        }
        varInfo.putVar(k, varName);
        varInfo.incrementVarCount(k);
    }

    @Override
    public void apply(KApply k) {
        varInfo.incrementKApplySignatureCount(k);
        super.apply(k);
    }
}
