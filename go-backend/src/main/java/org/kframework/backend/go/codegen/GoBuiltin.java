// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.ImmutableSet;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.kore.Sort;

import java.util.function.Function;

public class GoBuiltin {

    public static final ImmutableSet<String> HOOK_NAMESPACES;

    static {
        ImmutableSet.Builder<String> builder = ImmutableSet.builder();
        builder.add("BOOL").add("FLOAT").add("INT").add("IO").add("KEQUAL").add("KREFLECTION").add("LIST");
        builder.add("MAP").add("MINT").add("SET").add("STRING").add("ARRAY").add("BUFFER").add("BYTES");
        HOOK_NAMESPACES = builder.build();
    }

    /**
     * Some K framework function hook names are not suitable for Go. Here are the replacements.
     */
    public static final ImmutableMap<String, String> GO_FRIENDLY_HOOK_NAMES;

    static {
        ImmutableMap.Builder<String, String> builder = ImmutableMap.builder();
        builder.put("LIST.range", "listRange"); // range is a reserved word in Go
        builder.put("MAP.keys_list", "keysList"); // underscores
        builder.put("MAP.in_keys", "inKeys"); // underscores
        builder.put("ARRAY.in_keys", "inKeys"); // underscores
        GO_FRIENDLY_HOOK_NAMES = builder.build();
    }

    public static final ImmutableMap<String, Function<Sort, String>> OCAML_SORT_VAR_HOOKS;

    static {
        ImmutableMap.Builder<String, Function<Sort, String>> builder = ImmutableMap.builder();
        builder.put("BOOL.Bool", s -> "Bool _");
        builder.put("MINT.MInt", s -> "MInt _");
        builder.put("INT.Int", s -> "Int _");
        builder.put("FLOAT.Float", s -> "Float _");
        builder.put("STRING.String", s -> "String _");
        builder.put("BYTES.Bytes", s -> "Bytes _");
        builder.put("BUFFER.StringBuffer", s -> "StringBuffer _");
        builder.put("LIST.List", s -> "List (" + GoStringUtil.sortVariableName(s) + ",_,_)");
        builder.put("ARRAY.Array", s -> "Array (" + GoStringUtil.sortVariableName(s) + ",_,_)");
        builder.put("MAP.Map", s -> "Map (" + GoStringUtil.sortVariableName(s) + ",_,_)");
        builder.put("SET.Set", s -> "Set (" + GoStringUtil.sortVariableName(s) + ",_,_)");
        OCAML_SORT_VAR_HOOKS = builder.build();
    }

    public static final ImmutableSet<String> PREDICATE_HOOKS;

    static {
        ImmutableSet.Builder<String> builder = ImmutableSet.builder();
        builder.add("BOOL.Bool");
        builder.add("MINT.MInt");
        builder.add("INT.Int");
        builder.add("FLOAT.Float");
        builder.add("STRING.String");
        builder.add("BYTES.Bytes");
        builder.add("BUFFER.StringBuffer");
        builder.add("LIST.List");
        builder.add("ARRAY.Array");
        builder.add("MAP.Map");
        builder.add("SET.Set");
        PREDICATE_HOOKS = builder.build();
    }
}
