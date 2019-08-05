package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKSeqHeadNode extends LhsTreeNode {
    public final LhsKSeqSplitNode kseqSplit;

    public LhsKSeqHeadNode(LhsKSeqSplitNode kseqSplit) {
        super(kseqSplit);
        this.kseqSplit = kseqSplit;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsKSeqHeadNode)) {
            return false;
        }

        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = kseqSplit.headSubject;
    }
}
