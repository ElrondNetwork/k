package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.kore.KVariable;

import java.util.Objects;

public class LhsKApplyNode extends LhsTreeNode {
    public final String labelName;
    public final int arity;
    public final KVariable alias;
    private final String comment;

    public LhsKApplyNode(LhsTreeNode logicalParent, String labelName, int arity, KVariable alias, String comment) {
        super(logicalParent);
        this.labelName = labelName;
        this.arity = arity;
        this. alias = alias;
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
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
        subject = logicalParent.subject;

        writer.sb.writeIndent().append("if ");
        writer.matchWriter.appendKApplyMatch(writer.sb, logicalParent.subject, labelName, arity);
        writer.sb.beginBlock(comment);
        if (alias != null) {
            writer.sb.appendIndentedLine("// KVariable ", alias.name(), " = ", subject);
        }
    }

    @Override
    public String getKVariableMVRef(KVariable k) {
        if (Objects.equals(k, this.alias)) {
            return subject;
        }
        return super.getKVariableMVRef(k);
    }
}
