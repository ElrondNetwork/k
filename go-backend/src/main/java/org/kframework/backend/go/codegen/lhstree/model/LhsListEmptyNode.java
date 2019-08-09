package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsListEmptyNode extends LhsTreeNode {
    public final String sortName;
    public final String labelName;
    public final String comment;

    public LhsListEmptyNode(LhsTreeNode logicalParent, String sortName, String labelName, String comment) {
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
        if (!(other instanceof LhsListEmptyNode)) {
            return false;
        }
        LhsListEmptyNode otherNode = (LhsListEmptyNode)other;
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
        writer.sb.writeIndent().append("if i.Model.IsEmptyList(");
        writer.sb.append(logicalParent.subject).append(", ");
        writer.sb.append("m.").append(sortName).append(", ");
        writer.sb.append("m.").append(labelName);
        writer.sb.append(")");
        writer.sb.beginBlock(comment);

    }
}
