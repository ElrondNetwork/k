package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.kore.KVariable;

import java.util.Objects;

public class LhsKVarAlreadySeenNode extends LhsTreeNode {
    private final KVariable seenVar;

    public LhsKVarAlreadySeenNode(LhsTreeNode logicalParent, KVariable seenVar) {
        super(logicalParent);
        this.seenVar = seenVar;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsKVarAlreadySeenNode)) {
            return false;
        }
        LhsKVarAlreadySeenNode otherNode = (LhsKVarAlreadySeenNode)other;

        if (!Objects.equals(seenVar, otherNode.seenVar)) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = logicalParent.subject;
        String varName = getKVariableMVRef(seenVar);
        writer.sb.writeIndent();
        writer.sb.append("if i.Model.Equals(").append(subject).append(", ").append(varName).append(")");
        writer.sb.beginBlock("lhs KVariable, which reappears:" + seenVar.name());
    }
}
