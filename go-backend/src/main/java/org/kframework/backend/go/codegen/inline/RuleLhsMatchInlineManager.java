package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;

import java.util.HashSet;
import java.util.Set;

public class RuleLhsMatchInlineManager implements RuleLhsMatchWriter {

    public final Set<KApplySignatureMatch> kapplySignatures = new HashSet<>();

    @Override
    public void appendKApplyMatch(GoStringBuilder sb, String subject, String labelName, int arity) {
        KApplySignatureMatch signature = new KApplySignatureMatch(labelName, arity);
        kapplySignatures.add(signature);

        sb.append(subject).append("&kapplyMatchMask == ").append(signature.matchConstName);
    }

}
