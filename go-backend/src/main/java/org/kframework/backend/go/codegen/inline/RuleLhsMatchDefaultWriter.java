package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.utils.errorsystem.KEMException;

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

    @Override
    public void appendPredicateMatch(String hookName, GoStringBuilder sb, String subject, String sortName) {
        switch(hookName) {
        case "INT.Int":
            sb.append("m.IsInt(").append(subject).append(")");
            return;
        case "FLOAT.Float":
            sb.append("m.IsFloat(").append(subject).append(")");
            return;
        case "STRING.String":
            sb.append("m.IsString(").append(subject).append(")");
            return;
        case "BYTES.Bytes":
            sb.append("m.IsBytes(").append(subject).append(")");
            return;
        case "BUFFER.StringBuffer":
            sb.append("m.IsStringBuffer(").append(subject).append(")");
            return;
        case "BOOL.Bool":
            sb.append("m.IsBool(").append(subject).append(")");
            return;
        case "MINT.MInt":
            sb.append("m.IsMint(").append(subject).append(")");
            return;
        case "MAP.Map":
            sb.append("i.Model.IsMapWithSort(").append(subject).append(", m.").append(sortName).append(")");
            return;
        case "SET.Set":
            sb.append("i.Model.IsSetWithSort(").append(subject).append(", m.").append(sortName).append(")");
            return;
        case "LIST.List":
            sb.append("i.Model.IsList(").append(subject).append(", m.").append(sortName).append(")");
            return;
        case "ARRAY.Array":
            sb.append("i.Model.IsArray(").append(subject).append(", m.").append(sortName).append(")");
            return;
        default:
            throw KEMException.internalError("Unknown predicate hook: " + hookName);
        }
    }

    @Override
    public void appendBottomMatch(GoStringBuilder sb, String subject) {
        sb.append("m.IsBottom(").append(subject).append(")");
    }
}
