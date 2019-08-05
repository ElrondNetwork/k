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
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsKApplyArgNode)) {
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

    public void mergeTree(LhsTreeNode other) {
        this.merge(other);
    }
}
