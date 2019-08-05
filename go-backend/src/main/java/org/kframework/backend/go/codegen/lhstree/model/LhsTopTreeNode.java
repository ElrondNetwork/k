package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsTopTreeNode extends LhsTreeNode {

    public LhsTopTreeNode() {
        super(null);
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
}
