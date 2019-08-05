package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.kore.KVariable;

import java.util.Objects;

public class LhsKVarNode extends LhsTreeNode {
    protected final KVariable kvar;

    public LhsKVarNode(LhsTreeNode logicalParent, KVariable kvar) {
        super(logicalParent);
        this.kvar = kvar;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsKVarNode)) {
            return false;
        }
        LhsKVarNode otherNode = (LhsKVarNode)other;

        if (!Objects.equals(kvar, otherNode.kvar)) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = writer.vars.varIndexes.kvariableMVRef(kvar);
        writer.sb.writeIndent();
        writer.sb.append(subject).append(" = ").append(logicalParent.subject);
        writer.sb.append(" // lhs KVariable ").append(kvar.name()).newLine();
        //writer.sb.append(" // ").append(kvar.name()).newLine();
    }
}
