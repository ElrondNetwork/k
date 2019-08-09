package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsListMatchSplitNode extends LhsTreeNode {
    public final String sortName;
    public final String labelName;
    public final String comment;

    public String headSubject;
    public String tailSubject;

    public LhsListMatchSplitNode(LhsTreeNode logicalParent, String sortName, String labelName, String comment) {
        super(logicalParent);
        this.sortName = sortName;
        this.labelName = labelName;
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsListMatchSplitNode)) {
            return false;
        }
        LhsListMatchSplitNode otherNode = (LhsListMatchSplitNode)other;
        if (!sortName.equals(otherNode.sortName)) {
            return false;
        }
        if (!labelName.equals(otherNode.labelName)) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        headSubject = oneTimeVariableMVRef("listHead");
        tailSubject = oneTimeVariableMVRef("listTail");

        writer.sb.writeIndent().append("if i.tempBoolVar, ");
        writer.sb.append(headSubject).append(", ").append(tailSubject);
        writer.sb.append(" = i.Model.ListSplitHeadTail(");
        writer.sb.append(logicalParent.subject).append(", ");
        writer.sb.append("m.").append(sortName).append(", ");
        writer.sb.append("m.").append(labelName);
        writer.sb.append("); i.tempBoolVar");
        writer.sb.beginBlock(comment);

    }
}
