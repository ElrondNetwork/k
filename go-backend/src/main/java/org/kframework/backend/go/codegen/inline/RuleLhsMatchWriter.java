package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;

public interface RuleLhsMatchWriter {

    void appendKApplyMatch(GoStringBuilder sb, String subject, String labelName, int arity);

    void appendNonEmptyKSequenceMatch(GoStringBuilder sb, String subject);

    void appendNonEmptyKSequenceMinLengthMatch(GoStringBuilder sb, String subject, int minLength);

    void appendKTokenMatch(GoStringBuilder sb, String subject, String sortName);

    void appendPredicateMatch(String hookName, GoStringBuilder sb, String subject, String sortName);

    void appendBottomMatch(GoStringBuilder sb, String subject);
}
