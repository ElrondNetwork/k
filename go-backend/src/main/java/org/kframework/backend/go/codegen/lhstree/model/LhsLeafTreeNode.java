package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.backend.go.model.FunctionInfo;
import org.kframework.backend.go.model.Lookup;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.model.VarContainer;
import org.kframework.definition.Rule;
import org.kframework.kore.K;
import org.kframework.kore.KVariable;

import java.util.List;
import java.util.Set;

public class LhsLeafTreeNode extends LhsTreeNode {
    public final RuleType type;
    public final int ruleNum;
    public final FunctionInfo functionInfo;
    public final Rule rule;
    public final List<Lookup> lookups;
    public final K requires;
    public final K right;

    public final VarContainer vars;
    public final Set<KVariable> alreadySeenLhsVariables;

    public LhsLeafTreeNode(LhsTreeNode logicalParent, RuleType type, int ruleNum, FunctionInfo functionInfo, Rule rule, List<Lookup> lookups, K requires, K right, VarContainer vars, Set<KVariable> alreadySeenLhsVariables) {
        super(logicalParent);
        this.type = type;
        this.ruleNum = ruleNum;
        this.functionInfo = functionInfo;
        this.rule = rule;
        this.lookups = lookups;
        this.requires = requires;
        this.right = right;
        this.vars = vars;
        this.alreadySeenLhsVariables = alreadySeenLhsVariables;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        return false;
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        writer.writeLeaf(this);
    }
}
