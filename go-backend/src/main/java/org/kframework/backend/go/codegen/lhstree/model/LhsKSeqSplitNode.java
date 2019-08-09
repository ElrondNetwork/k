package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKSeqSplitNode extends LhsTreeNode {
    public String headSubject;
    public String tailSubject;
    public final String comment;

    public LhsKSeqSplitNode(LhsTreeNode logicalParent, String comment) {
        super(logicalParent);
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsKSeqSplitNode)) {
            return false;
        }

        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        headSubject = oneTimeVariableMVRef("kseqHead");
        tailSubject = oneTimeVariableMVRef("kseqTail");

        writer.sb.writeIndent().append("_, ");
        writer.sb.append(headSubject).append(", ");
        writer.sb.append(tailSubject).append(" = i.Model.KSequenceSplitHeadTail(").append(logicalParent.subject).append(") // ");
        writer.sb.append(comment);
        writer.sb.newLine();
    }
}
