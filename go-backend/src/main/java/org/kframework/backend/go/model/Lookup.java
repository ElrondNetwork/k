package org.kframework.backend.go.model;

import org.kframework.kore.KApply;

public class Lookup {

    public enum Type { MATCH, SETCHOICE }

    private final Type type;
    private final KApply content;

    public Lookup(Type type, KApply content) {
        this.type = type;
        this.content = content;
    }

    public Type getType() {
        return type;
    }

    public KApply getContent() {
        return content;
    }
}
