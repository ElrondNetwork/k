package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsListHeadNode extends LhsTreeNode {
    public final LhsListMatchSplitNode listParent;

    public LhsListHeadNode(LhsListMatchSplitNode listParent) {
        super(listParent);
        this.listParent = listParent;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsListHeadNode)) {
            return false;
        }

        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = listParent.headSubject;
    }
}
