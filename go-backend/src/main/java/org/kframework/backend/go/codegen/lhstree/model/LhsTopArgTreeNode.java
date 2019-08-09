package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsTopArgTreeNode extends LhsTreeNode {
    public final int argIndex;
    public final String argVarName;

    public LhsTopArgTreeNode(LhsTopTreeNode topParent, int argIndex, String argVarName) {
        super(topParent);
        this.argIndex = argIndex;
        this.argVarName = argVarName;
    }

    @Override
    protected void changeParent(LhsTreeNode logicalParent) {
        if (!(logicalParent instanceof LhsTopTreeNode)) {
            throw new RuntimeException("LhsTopArgTreeNode can only have a LhsTopTreeNode as logical parent");
        }
        super.changeParent(logicalParent);
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsTopArgTreeNode)) {
            return false;
        }
        LhsTopArgTreeNode otherNode = (LhsTopArgTreeNode) other;
        if (argIndex != otherNode.argIndex) {
            return false;
        }

        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = argVarName;
    }

}
