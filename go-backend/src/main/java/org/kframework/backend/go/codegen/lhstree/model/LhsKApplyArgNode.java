package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKApplyArgNode extends LhsTreeNode {
    public LhsKApplyNode kappParent;
    public final int argIndex;

    public LhsKApplyArgNode(LhsKApplyNode kappParent, int argIndex) {
        super(kappParent);
        this.kappParent = kappParent;
        this.argIndex = argIndex;
    }

    @Override
    protected void changeParent(LhsTreeNode logicalParent) {
        if (!(logicalParent instanceof LhsKApplyNode)) {
            throw new RuntimeException("LhsKApplyArgNode can only have a LhsKApplyNode as logical parent");
        }
        kappParent = (LhsKApplyNode)logicalParent;
        if (argIndex >= kappParent.arity) {
            throw new RuntimeException("LhsKApplyArgNode cannot have argument index that exceeds parent KApply arity");
        }
        super.changeParent(logicalParent);
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (other == this) {
            return true;
        }
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
        subject = oneTimeVariableMVRef("kappVar");
        writer.sb.appendIndentedLine(subject, " = i.Model.KApplyArg(", logicalParent.subject, ", ", Integer.toString(argIndex), ")");

    }
}
