package org.kframework.backend.go.model;

public class VarContainer {
    public final RuleVars lhsVars;
    public final RuleVars rhsVars;

    public final TempVarIndexes varIndexes = new TempVarIndexes();

    public VarContainer(RuleVars lhsVars, RuleVars rhsVars) {
        this.lhsVars = lhsVars;
        this.rhsVars = rhsVars;
    }

}
