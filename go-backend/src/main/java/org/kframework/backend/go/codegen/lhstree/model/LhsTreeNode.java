package org.kframework.backend.go.codegen.lhstree.model;


import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;

import java.util.ArrayList;
import java.util.List;

public abstract class LhsTreeNode {
    public final LhsTreeNode logicalParent;
    public final List<LhsTreeNode> children = new ArrayList<>();
    protected String subject;

    public LhsTreeNode(LhsTreeNode logicalParent) {
        this.logicalParent = logicalParent;
    }

    public abstract boolean matches(LhsTreeNode other);

    public abstract void write(RuleLhsTreeWriter writer);

    public void addChild(LhsTreeNode child) {
        children.add(child);
    }

    protected void merge(LhsTreeNode other) {
        for (LhsTreeNode otherChild : other.children) {
            boolean matched = false;
            for (LhsTreeNode child:children) {
                if (otherChild.matches(child)) {
                    matched = true;
                    child.merge(otherChild);
                    break;
                }
            }
            if (!matched) {
                children.add(otherChild);
            }
        }
    }

}
