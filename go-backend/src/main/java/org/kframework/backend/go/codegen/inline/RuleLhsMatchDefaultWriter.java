package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;

public class RuleLhsMatchDefaultWriter implements RuleLhsMatchWriter {

    @Override
    public void appendKApplyMatch(GoStringBuilder sb, String subject, String labelName, int arity) {
        sb.append("m.MatchKApply(").append(subject).append(", ");
        sb.append("uint64(m.").append(labelName).append("), ");
        sb.append(arity).append(")");
    }

    @Override
    public void appendNonEmptyKSequenceMatch(GoStringBuilder sb, String subject) {
        sb.writeIndent().append("m.MatchNonEmptyKSequence(");
        sb.append(subject).append(")");
    }

    @Override
    public void appendNonEmptyKSequenceMinLengthMatch(GoStringBuilder sb, String subject, int minLength) {
        sb.writeIndent().append("m.MatchNonEmptyKSequenceMinLength(");
        sb.append(subject).append(", ").append(minLength).append(")");
    }

    @Override
    public void appendKTokenMatch(GoStringBuilder sb, String subject, String sortName) {
        sb.append("m.MatchKToken(").append(subject).append(", ");
        sb.append("uint64(m.").append(sortName).append("))");

    }
}
