// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.model;

import org.kframework.kore.KLabel;

import java.util.Arrays;
import java.util.Iterator;
import java.util.stream.Stream;

public class FunctionParams {
    private static final String EVAL_ARG_NAME = "c";
    private final String[] varNames;

    public FunctionParams(int arity) {
        varNames = new String[arity];
        if (arity == 0) {
            // nothing
        } else if (arity == 1) {
            varNames[0] = EVAL_ARG_NAME;
        } else {
            for (int i = 0; i < arity; i++) {
                varNames[i] = EVAL_ARG_NAME + (i + 1);
            }
        }
    }

    public int arity() {
        return varNames.length;
    }

    public String varName(int i) {
        return varNames[i];
    }

    public String parameterDeclaration() {
        if (varNames.length == 0) {
            return "";
        }
        StringBuilder sb = new StringBuilder();
        for (String v : varNames) {
            sb.append(v);
            sb.append(" m.KReference, ");
        }
        return sb.toString();
    }

    public String callParameters() {
        if (varNames.length == 0) {
            return "";
        }
        StringBuilder sb = new StringBuilder();
        for (String v : varNames) {
            sb.append(v);
            sb.append(", ");
        }
        return sb.toString();
    }

    public String paramNamesSeparatedByComma() {
        return String.join(", ", varNames);
    }

    public Iterable<String> getVarNames() {
        return new Iterable<String>() {
            @Override
            public Iterator<String> iterator() {
                return Arrays.stream(varNames).iterator();
            }
        };
    }

    public Stream<String> varNamesStream() {
        return Arrays.stream(varNames);
    }

}
