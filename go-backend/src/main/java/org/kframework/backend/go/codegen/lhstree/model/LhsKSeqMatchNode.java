package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKSeqMatchNode extends LhsTreeNode {
    public final int minLength;
    public final String comment;

    public LhsKSeqMatchNode(LhsTreeNode logicalParent, int minLength, String comment) {
        super(logicalParent);
        this.minLength = minLength;
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsKSeqMatchNode)) {
            return false;
        }
        LhsKSeqMatchNode otherNode = (LhsKSeqMatchNode)other;
        if (minLength != otherNode.minLength) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = logicalParent.subject;

        writer.sb.writeIndent().append("if ");
        if (minLength == 1) {
            writer.matchWriter.appendNonEmptyKSequenceMatch(writer.sb, subject);
        } else {
            writer.matchWriter.appendNonEmptyKSequenceMinLengthMatch(writer.sb, subject, minLength);
        }
        writer.sb.beginBlock(comment);

    }
}
