package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsBottomTreeNode extends LhsTreeNode {

    public LhsBottomTreeNode(LhsTreeNode logicalParent) {
        super(logicalParent);
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsBottomTreeNode)) {
            return false;
        }

        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        writer.sb.writeIndent().append("if ");
        writer.matchWriter.appendBottomMatch(writer.sb, logicalParent.subject);
    }
}
