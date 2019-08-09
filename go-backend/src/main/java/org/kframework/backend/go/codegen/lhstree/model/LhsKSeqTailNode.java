package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKSeqTailNode extends LhsTreeNode {
    public LhsKSeqSplitNode kseqSplit;

    public LhsKSeqTailNode(LhsKSeqSplitNode kseqSplit) {
        super(kseqSplit);
    }

    @Override
    protected void changeParent(LhsTreeNode logicalParent) {
        if (!(logicalParent instanceof LhsKSeqSplitNode)) {
            throw new RuntimeException("LhsKSeqTailNode can only have a LhsKSeqSplitNode as logical parent");
        }
        kseqSplit = (LhsKSeqSplitNode)logicalParent;
        super.changeParent(logicalParent);
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsKSeqTailNode)) {
            return false;
        }

        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = kseqSplit.tailSubject;
    }
}
