package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKSeqOneNode extends LhsTreeNode {
    public final String comment;

    public LhsKSeqOneNode(LhsTreeNode logicalParent, String comment) {
        super(logicalParent);
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsKSeqOneNode)) {
            return false;
        }
        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = logicalParent.subject;
        writer.sb.appendIndentedLine("// KSequence, size 1:", comment);
    }
}
