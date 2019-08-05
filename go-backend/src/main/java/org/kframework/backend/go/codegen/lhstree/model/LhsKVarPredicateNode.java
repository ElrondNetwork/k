package org.kframework.backend.go.codegen.lhstree.model;

import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

public class LhsKVarPredicateNode extends LhsTreeNode {
    private final String hook;
    private final String sortName;
    private final String comment;

    public LhsKVarPredicateNode(LhsTreeNode logicalParent, String hook, String sortName, String comment) {
        super(logicalParent);
        this.hook = hook;
        this.sortName = sortName;
        this.comment = comment;
    }

    @Override
    public boolean matches(LhsTreeNode other) {
        if (!(other instanceof LhsKVarPredicateNode)) {
            return false;
        }
        LhsKVarPredicateNode otherNode = (LhsKVarPredicateNode)other;

        if (!hook.equals(otherNode.hook)) {
            return false;
        }
        if (!sortName.equals(otherNode.sortName)) {
            return false;
        }

        return otherNode.logicalParent.matches(logicalParent);
    }

    @Override
    public void write(RuleLhsTreeWriter writer) {
        subject = logicalParent.subject;
        writer.sb.writeIndent().append("if ");
        writer.matchWriter.appendPredicateMatch(hook, writer.sb, subject, sortName);
        writer.sb.beginBlock(comment);
    }
}
