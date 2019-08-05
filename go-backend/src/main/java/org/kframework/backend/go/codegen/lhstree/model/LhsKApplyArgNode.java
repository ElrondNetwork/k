package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKApplyArgNode extends LhsTreeNode {
    public final LhsKApplyNode kappParent;
    public final int argIndex;

    public LhsKApplyArgNode(LhsKApplyNode kappParent, int argIndex) {
        super(kappParent);
        this.kappParent = kappParent;
        this.argIndex = argIndex;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsKApplyArgNode)) {
            return false;
        }
        LhsKApplyArgNode otherNode = (LhsKApplyArgNode)other;
        if (argIndex != otherNode.argIndex) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = writer.vars.varIndexes.oneTimeVariableMVRef("kappVar");
        writer.sb.appendIndentedLine(subject, " = i.Model.KApplyArg(", logicalParent.subject, ", ", Integer.toString(argIndex), ")");

    }
}
