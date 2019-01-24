package org.kframework.backend.go.processors;

import org.kframework.backend.go.model.Lookup;
import org.kframework.kore.KApply;
import org.kframework.utils.errorsystem.KEMException;

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
            switch (lookup.getType()) {
            case MATCH:
                KApply k = lookup.getContent();
                if (k.klist().items().size() != 2) {
                    throw KEMException.internalError("Unexpected arity of lookup: " + k.klist().size(), k);
                }
                accumLhsVars.apply(k.klist().items().get(0));
                accumRhsVars.apply(k.klist().items().get(1));
                break;
            default:
                throw KEMException.internalError("Unexpected lookup type");
            }
        }
    }
}
