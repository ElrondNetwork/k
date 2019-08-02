package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.utils.errorsystem.KEMException;

import java.util.HashSet;
import java.util.Set;

public class RuleLhsMatchInlineManager implements RuleLhsMatchWriter {

    public final Set<KApplySignatureMatch> kapplySignatures = new HashSet<>();
    public final Set<String> ktokenSortNames = new HashSet<>();
    public final Set<String> mapSortNames = new HashSet<>();
    public final Set<String> listSortNames = new HashSet<>();
    public final Set<String> setSortNames = new HashSet<>();
    public final Set<String> arraySortNames = new HashSet<>();

    @Override
    public void appendKApplyMatch(GoStringBuilder sb, String subject, String labelName, int arity) {
        KApplySignatureMatch signature = new KApplySignatureMatch(labelName, arity);
        kapplySignatures.add(signature);

        sb.append(subject).append("&kapplyMatchMask == ").append(signature.matchConstName);
    }

    @Override
    public void appendNonEmptyKSequenceMatch(GoStringBuilder sb, String subject) {
        sb.append(subject).append(">>refTypeShift != refEmptyKseqTypeAsUint");
    }

    @Override
    public void appendNonEmptyKSequenceMinLengthMatch(GoStringBuilder sb, String subject, int minLength) {
        sb.append(subject).append(">>refTypeShift == refNonEmptyKseqTypeAsUint && (");
        sb.append(subject).append(">>refNonEmptyKseqIndexShift&refNonEmptyKseqLengthMask) >= ");
        sb.append(minLength);
    }

    public String ktokenMatchConstName(String sortName) {
        return "ktokenMatch" + sortName;
    }

    @Override
    public void appendKTokenMatch(GoStringBuilder sb, String subject, String sortName) {
        ktokenSortNames.add(sortName);
        sb.append(subject).append("&ktokenMatchMask == ").append(ktokenMatchConstName(sortName));
    }

    public String mapMatchConstName(String sortName) {
        return "mapMatch" + sortName;
    }

    public String setMatchConstName(String sortName) {
        return "setMatch" + sortName;
    }

    public String listMatchConstName(String sortName) {
        return "listMatch" + sortName;
    }

    public String arrayMatchConstName(String sortName) {
        return "arrayMatch" + sortName;
    }

    @Override
    public void appendPredicateMatch(String hookName, GoStringBuilder sb, String subject, String sortName) {
        switch(hookName) {
        case "INT.Int":
            sb.append("i.tempTypeVar = ").append(subject).append(" >> refTypeShift; i.tempTypeVar == uint64(smallPositiveIntRef) || i.tempTypeVar == uint64(smallNegativeIntRef) || i.tempTypeVar == uint64(bigIntRef)");
            return;
        case "FLOAT.Float":
            sb.append(subject).append(">>refTypeShift == uint64(floatRef)");
            return;
        case "STRING.String":
            sb.append(subject).append(">>refTypeShift == uint64(stringRef)");
            return;
        case "BYTES.Bytes":
            sb.append(subject).append(">>refTypeShift == uint64(bytesRef)");
            return;
        case "BUFFER.StringBuffer":
            sb.append(subject).append(">>refTypeShift == uint64(stringBufferRef)");
            return;
        case "BOOL.Bool":
            sb.append(subject).append(">>refTypeShift == uint64(boolRef)");
            return;
        case "MINT.MInt":
            sb.append(subject).append(">>refTypeShift == uint64(mintRef)");
            return;
        case "MAP.Map":
            mapSortNames.add(sortName);
            sb.append(subject).append("&collectionMatchMask == ").append(mapMatchConstName(sortName));
            return;
        case "SET.Set":
            setSortNames.add(sortName);
            sb.append(subject).append("&collectionMatchMask == ").append(setMatchConstName(sortName));
            return;
        case "LIST.List":
            listSortNames.add(sortName);
            sb.append(subject).append("&collectionMatchMask == ").append(listMatchConstName(sortName));
            return;
        case "ARRAY.Array":
            arraySortNames.add(sortName);
            sb.append(subject).append("&collectionMatchMask == ").append(arrayMatchConstName(sortName));
            return;
        default:
            throw KEMException.internalError("Unknown predicate hook: " + hookName);
        }
    }

    @Override
    public void appendBottomMatch(GoStringBuilder sb, String subject) {
        sb.append(subject).append(" == m.InternedBottom");
    }
}
