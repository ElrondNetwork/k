package org.kframework.backend.go;

import org.kframework.kore.KVariable;
import org.kframework.kore.VisitK;

import java.util.HashSet;
import java.util.Set;

public class AccumulateVariableNames extends VisitK {

    private final Set<String> varNames = new HashSet<>();

    public AccumulateVariableNames() {
    }

    public Set<String> getVarNames() {
        return varNames;
    }

    @Override
    public void apply(KVariable k) {
        varNames.add(k.name());
    }
}
