package org.kframework.backend.go.processors;

import org.kframework.backend.go.model.Lookup;

import java.util.List;

public class LookupVarExtractor {

    private final AccumulateRuleVars accumLhsVars;
    private final AccumulateRuleVars accumRhsVars;

    public LookupVarExtractor(AccumulateRuleVars accumLhsVars, AccumulateRuleVars accumRhsVars) {
        this.accumLhsVars = accumLhsVars;
        this.accumRhsVars = accumRhsVars;
    }

    public void apply(List<Lookup> lookups) {
        for (Lookup lookup : lookups) {
            accumLhsVars.apply(lookup.getLhs());
            accumRhsVars.apply(lookup.getRhs());
        }
    }
}
