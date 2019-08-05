package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKApplyNode extends LhsTreeNode {
    public final String labelName;
    public final int arity;
    private final String comment;

    public LhsKApplyNode(LhsTreeNode logicalParent, String labelName, int arity, String comment) {
        super(logicalParent);
        this.labelName = labelName;
        this.arity = arity;
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsKApplyNode)) {
            return false;
        }
        LhsKApplyNode otherNode = (LhsKApplyNode)other;
        if (!labelName.equals(otherNode.labelName)) {
            return false;
        }
        if (arity != otherNode.arity) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        writer.sb.writeIndent().append("if ");
        writer.matchWriter.appendKApplyMatch(writer.sb, logicalParent.subject, labelName, arity);
        writer.sb.beginBlock(comment);
        subject = logicalParent.subject;
    }
}
