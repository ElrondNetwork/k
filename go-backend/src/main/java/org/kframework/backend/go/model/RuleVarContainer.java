package org.kframework.backend.go.model;

import org.kframework.kore.KVariable;

import java.util.HashSet;
import java.util.Set;
import java.util.TreeSet;

public class RuleVarContainer {
    public final RuleVars lhsVars;
    public final RuleVars requiresVars;
    public final RuleVars rhsVars;

    /**
     * Whenever we see a variable more than once, instead of adding a variable declaration, we add a check that the two instances are equal.
     * This structure keeps track of that.
     * It is shared between main LHS and lookup LHS
     */
    public final Set<KVariable> alreadySeenLhsVariables = new HashSet<>();

    /**
     * Easy way to figure out on the RHS if a variable is a KApply.
     * We populate this from the LHS.
     */
    public final Set<String> kapplyVariableNames = new TreeSet<>();

    /**
     * Helps generate unique variable names.z
     */
    public final TempVarCounters varCounters = new TempVarCounters();

    public RuleVarContainer(RuleVars lhsVars, RuleVars requiresVars, RuleVars rhsVars) {
        this.lhsVars = lhsVars;
        this.requiresVars = requiresVars;
        this.rhsVars = rhsVars;
    }
}
