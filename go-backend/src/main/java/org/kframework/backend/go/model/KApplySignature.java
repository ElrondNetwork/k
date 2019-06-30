package org.kframework.backend.go.model;

import org.kframework.kore.KApply;
import org.kframework.kore.KLabel;

import java.util.Objects;

public class KApplySignature {

    public final KLabel label;
    public final int arity;

    private KApplySignature(KApply kapp) {
        this.label = kapp.klabel();
        this.arity = kapp.items().size();
    }

    public static KApplySignature of(KApply kapp) {
        return new KApplySignature(kapp);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        KApplySignature that = (KApplySignature) o;
        return arity == that.arity &&
                label.equals(that.label);
    }

    @Override
    public int hashCode() {
        return Objects.hash(label, arity);
    }
}
