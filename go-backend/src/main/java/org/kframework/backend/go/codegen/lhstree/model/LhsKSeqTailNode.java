package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKSeqTailNode extends LhsTreeNode {
    public final LhsKSeqSplitNode kseqSplit;

    public LhsKSeqTailNode(LhsKSeqSplitNode kseqSplit) {
        super(kseqSplit);
        this.kseqSplit = kseqSplit;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
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
