package org.kframework.backend.go.codegen.lhstree.model;


import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.backend.go.model.TempVarManager;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.kore.KVariable;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.TreeSet;

public abstract class LhsTreeNode implements TempVarManager {
    protected LhsTreeNode logicalParent;
    protected final List<LhsTreeNode> logicalChildren = new ArrayList<>();
    public LhsTreeNode predecessor;
    public final List<LhsTreeNode> successors = new ArrayList<>();
    protected String subject;

    protected final Set<Integer> rulesBelow = new TreeSet<>();
    private Integer nrVars = null;
    private Integer nrBoolVars = null;

    public LhsTreeNode(LhsTreeNode logicalParent) {
        changeParent(logicalParent);
    }

    protected void changeParent(LhsTreeNode logicalParent) {
        if (logicalParent == null) {
            throw new RuntimeException("logicalParent cannot be null");
        }
        this.logicalParent = logicalParent;
        logicalParent.logicalChildren.add(this);

        nrVars = logicalParent.nrVars;
        nrBoolVars = logicalParent.nrBoolVars;
    }

    public abstract boolean matches(LhsTreeNode other);

    public abstract void write(RuleLhsTreeWriter writer);

    public void chainNode(LhsTreeNode child) {
        child.predecessor = this;
        successors.add(child);
    }

    protected void merge(LhsTreeNode other) {
        for (LhsTreeNode otherChild : other.successors) {
            boolean matched = false;
            for (LhsTreeNode child : successors) {
                if (otherChild.matches(child)) {
                    matched = true;
                    child.merge(otherChild);
                    break;
                }
            }
            if (!matched) {
                chainNode(otherChild);
            }
        }
        for (LhsTreeNode logicalChild : other.logicalChildren) {
            logicalChild.changeParent(this);
        }

        // the merged node should not be usable after this
        other.successors.clear();
    }

    public void findRulesBelow() {
        rulesBelow.clear();
        for (LhsTreeNode succ : successors) {
            succ.findRulesBelow();
            rulesBelow.addAll(succ.rulesBelow);
        }
    }

    public void writeRuleInfo(RuleLhsTreeWriter writer) {
        if (rulesBelow.size() == predecessor.rulesBelow.size()) {
            return; // only print when rule set is getting split
        }
        if (rulesBelow.size() == 1) {
            writer.sb.writeIndent().append("// rule: ");
        } else {
            writer.sb.writeIndent().append("// rules: ");
        }
        writeCommaSeparatedRuleNumbers(writer.sb);
        writer.sb.newLine();
    }

    protected void writeCommaSeparatedRuleNumbers(GoStringBuilder sb) {
        boolean first = true;
        for (Integer ruleNum : rulesBelow) {
            if (first) {
                first = false;
            } else {
                sb.append(", ");
            }
            sb.append(ruleNum);
        }
    }

    @Override
    public String oneTimeVariableMVRef(String varName) {
        if (nrVars == null) {

        }
        //return "v[" + nextVarIndex() + " /*" + varName + "*/]";
        return "v[" + nextVarIndex() + "]";
    }

    @Override
    public String addKVariableMVRef(KVariable k) {
        throw new RuntimeException("Adding KVariable not supported by default in LHS tree nodes.");
    }

    @Override
    public String getKVariableMVRef(KVariable k) {
        if (predecessor == null) {
            return null;
        }
        return predecessor.getKVariableMVRef(k);
    }

    @Override
    public String evalBoolVarRef(String varName) {
        return "bv[" + nextBoolVar() + "]";
    }

    @Override
    public int getNrVars() {
        if (nrVars == null) {
            nrVars = predecessor.getNrVars();
        }
        return nrVars;
    }

    protected int nextVarIndex() {
        int varIndex = getNrVars();
        nrVars++;
        return varIndex;
    }

    protected String nextVar(String comment) {
        //return "v[" + nextVarIndex() + "]";
        return "v[" + nextVarIndex() + " /*" + comment + "*/]";
    }

    public int maxNrVars() {
        if (successors.size() == 0) {
            return getNrVars();
        }
        int maxVars = 0;
        for (LhsTreeNode child : successors) {
            int childMaxVars = child.maxNrVars();
            if (childMaxVars > maxVars) {
                maxVars = childMaxVars;
            }
        }
        return maxVars;
    }

    @Override
    public int getNrBoolVars() {
        if (nrBoolVars == null) {
            nrBoolVars = predecessor.getNrBoolVars();
        }
        return nrBoolVars;
    }

    protected int nextBoolVar() {
        int varIndex = getNrBoolVars();
        nrBoolVars++;
        return varIndex;
    }

    public int maxNrBoolVars() {
        if (successors.size() == 0) {
            return getNrBoolVars();
        }
        int maxVars = 0;
        for (LhsTreeNode child : successors) {
            int childMaxVars = child.maxNrBoolVars();
            if (childMaxVars > maxVars) {
                maxVars = childMaxVars;
            }
        }
        return maxVars;
    }
}
