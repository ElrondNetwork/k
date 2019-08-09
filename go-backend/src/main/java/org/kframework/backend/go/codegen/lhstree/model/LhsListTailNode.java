package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsListTailNode extends LhsTreeNode {
    public LhsListMatchSplitNode listParent;

    public LhsListTailNode(LhsListMatchSplitNode listParent) {
        super(listParent);
    }

    @Override
    protected void changeParent(LhsTreeNode logicalParent) {
        if (!(logicalParent instanceof LhsListMatchSplitNode)) {
            throw new RuntimeException("LhsListTailNode can only have a LhsListMatchSplitNode as logical parent");
        }
        listParent = (LhsListMatchSplitNode)logicalParent;
        super.changeParent(logicalParent);
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
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
