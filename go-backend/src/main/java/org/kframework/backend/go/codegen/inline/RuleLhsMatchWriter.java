package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.strings.GoStringBuilder;

public interface RuleLhsMatchWriter {

    void appendKApplyMatch(GoStringBuilder sb, String subject, String labelName, int arity);

}
