package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;

public class RuleLhsMatchDefaultWriter implements RuleLhsMatchWriter {

    @Override
    public void appendKApplyMatch(GoStringBuilder sb, String subject, String labelName, int arity) {
        sb.append("m.MatchKApply(").append(subject).append(", ");
        sb.append("uint64(m.").append(labelName).append("), ");
        sb.append(arity).append(")");
    }
}
