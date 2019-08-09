package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsHashKTokenTreeNode extends LhsTreeNode {
    public final String sortName;

    public LhsHashKTokenTreeNode(LhsTreeNode logicalParent, String sortName) {
        super(logicalParent);
        this.sortName = sortName;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsHashKTokenTreeNode)) {
            return false;
        }
        LhsHashKTokenTreeNode otherNode = (LhsHashKTokenTreeNode) other;
        if (!otherNode.sortName.equals(sortName)) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = logicalParent.subject;
        writer.sb.writeIndent();
        writer.sb.append("if ");
        writer.matchWriter.appendKTokenMatch(writer.sb, subject, sortName);
        writer.sb.beginBlock("lhs KApply #KToken");
    }
}
