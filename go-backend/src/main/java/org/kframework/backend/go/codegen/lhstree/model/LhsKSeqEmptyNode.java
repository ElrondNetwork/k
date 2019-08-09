package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKSeqEmptyNode extends LhsTreeNode {
    public final String comment;

    public LhsKSeqEmptyNode(LhsTreeNode logicalParent, String comment) {
        super(logicalParent);
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsKSeqEmptyNode)) {
            return false;
        }
        return other.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        writer.sb.appendIndentedLine("// KSequence, size 0:", comment);
    }
}
