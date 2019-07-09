package org.kframework.backend.go.codegen.inline;

import java.util.Objects;

public class KApplySignatureMatch {

    public final String labelName;
    public final int arity;
    public final String matchConstName;

    public KApplySignatureMatch(String labelName, int arity) {
        this.labelName = labelName;
        this.arity = arity;
        this.matchConstName = "kapplyMatch" + labelName + arity;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        KApplySignatureMatch that = (KApplySignatureMatch) o;
        return arity == that.arity &&
                Objects.equals(labelName, that.labelName);
    }

    @Override
    public int hashCode() {
        return Objects.hash(labelName, arity);
    }
}
