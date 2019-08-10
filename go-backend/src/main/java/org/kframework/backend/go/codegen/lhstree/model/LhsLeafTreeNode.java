package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.backend.go.model.FunctionInfo;
import org.kframework.backend.go.model.Lookup;
import org.kframework.backend.go.model.RuleType;
import org.kframework.definition.Rule;
import org.kframework.kore.K;
import org.kframework.kore.KVariable;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class LhsLeafTreeNode extends LhsTreeNode {
    public final RuleType type;
    public final int ruleNum;
    public final FunctionInfo functionInfo;
    public final Rule rule;
    public final List<Lookup> lookups;
    public final K requires;
    public final K right;

    public final Map<KVariable, String> variablesDeclaredInLeaf = new HashMap<>();

    public LhsLeafTreeNode(LhsTreeNode logicalParent, RuleType type, int ruleNum, FunctionInfo functionInfo, Rule rule, List<Lookup> lookups, K requires, K right, Set<KVariable> alreadySeenLhsVariables) {
        super(logicalParent);
        this.type = type;
        this.ruleNum = ruleNum;
        this.functionInfo = functionInfo;
        this.rule = rule;
        this.lookups = lookups;
        this.requires = requires;
        this.right = right;
    }

    @Override
    public void findRulesBelow() {
        rulesBelow.clear();
        rulesBelow.add(ruleNum);
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        return false;
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        writer.writeLeaf(this);
    }

    @Override
    public String addKVariableMVRef(KVariable k) {
        if (super.getKVariableMVRef(k) != null) {
            throw new RuntimeException("cannot add to LhsLeafTreeNode a variable already declared in a previous node");
        }
        if (variablesDeclaredInLeaf.containsKey(k)) {
            throw new RuntimeException("cannot add the same KVariable twice in LhsLeafTreeNode");
        }
        String mvRef = nextVar(k.name());
        variablesDeclaredInLeaf.put(k, mvRef);
        return mvRef;
    }

    @Override
    public String getKVariableMVRef(KVariable k) {
        String mvRef = variablesDeclaredInLeaf.get(k);
        if (mvRef != null) {
            return mvRef;
        }
        return super.getKVariableMVRef(k);
    }
}
