package org.kframework.backend.go.model;

import org.kframework.kore.K;
import org.kframework.unparser.ToKast;

public class Lookup {

    public enum Type {
        MATCH,
        SETCHOICE,
        MAPCHOICE,
        FILTERMAPCHOICE
    }

    private final Type type;
    private final K lhs;
    private final K rhs;

    public Lookup(Type type, K lhs, K rhs) {
        this.type = type;
        this.lhs = lhs;
        this.rhs = rhs;
    }

    public Type getType() {
        return type;
    }

    public K getLhs() {
        return lhs;
    }

    public K getRhs() {
        return rhs;
    }

    public String comment() {
        String lookupName;
        switch (type) {
        case MATCH:
            lookupName = "#match";
            break;
        case SETCHOICE:
            lookupName = "#setChoice";
            break;
        case MAPCHOICE:
            lookupName = "#mapChoice";
            break;
        case FILTERMAPCHOICE:
            lookupName = "#filterMapChoice";
            break;
        default:
            throw new RuntimeException("Unknown Lookup.Type");
        }

        return lookupName + "(" + ToKast.apply(lhs) + "," + ToKast.apply(rhs) + ")";
    }
}
