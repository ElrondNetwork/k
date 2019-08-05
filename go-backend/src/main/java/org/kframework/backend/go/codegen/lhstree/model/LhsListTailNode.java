package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsListTailNode extends LhsTreeNode {
    public final LhsListMatchSplitNode listParent;

    public LhsListTailNode(LhsListMatchSplitNode listParent) {
        super(listParent);
        this.listParent = listParent;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsListTailNode)) {
            return false;
        }

        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = listParent.tailSubject;
    }
}
