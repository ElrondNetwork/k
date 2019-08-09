package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsTopTreeNode extends LhsTreeNode {

    public LhsTopTreeNode() {
        super(null);
    }

    @Override
    protected void changeParent(LhsTreeNode logicalParent) {
        if (logicalParent != null) {
            throw new RuntimeException("LhsTopTreeNode cannot have a logical parent");
        }
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        return other instanceof LhsTopTreeNode;
    }


    @Override
    public void write(RuleLhsTreeWriter writer) {
    }

    public void mergeTree(LhsTreeNode other) {
        this.merge(other);
    }

    @Override
    public int getNrVars() {
        return 0;
    }

    @Override
    public int getNrBoolVars() {
        return 0;
    }
}
