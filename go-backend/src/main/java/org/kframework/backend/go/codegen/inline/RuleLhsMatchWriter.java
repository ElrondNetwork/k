package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;

public interface RuleLhsMatchWriter {

    void appendKApplyMatch(GoStringBuilder sb, String subject, String labelName, int arity);

    void appendNonEmptyKSequenceMatch(GoStringBuilder sb, String subject);

    void appendNonEmptyKSequenceMinLengthMatch(GoStringBuilder sb, String subject, int minLength);

}
