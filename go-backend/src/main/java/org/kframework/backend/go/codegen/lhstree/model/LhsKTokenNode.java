package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.backend.go.codegen.rules.RuleRhsWriterBase;
import org.kframework.kore.KToken;
import org.kframework.unparser.ToKast;

public class LhsKTokenNode extends LhsTreeNode {
    private final KToken ktoken;

    public LhsKTokenNode(LhsTreeNode logicalParent, KToken ktoken) {
        super(logicalParent);
        this.ktoken = ktoken;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
        if (!(other instanceof LhsKTokenNode)) {
            return false;
        }
        LhsKTokenNode otherNode = (LhsKTokenNode)other;
        if (!ktoken.equals(otherNode.ktoken)) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = logicalParent.subject;
        writer.sb.writeIndent();
        writer.sb.append("if i.Model.Equals(").append(subject).append(", ");
        RuleRhsWriterBase.appendKTokenRepresentation(writer.sb, ktoken, writer.data, writer.nameProvider);
        writer.sb.append(")");
        writer.sb.beginBlock(ToKast.apply(ktoken));
    }
}
