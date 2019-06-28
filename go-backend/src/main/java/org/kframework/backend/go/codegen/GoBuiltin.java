// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.ImmutableSet;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.kore.Sort;

import java.util.function.BiFunction;
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

    public static final ImmutableMap<String, BiFunction<String, String, String>> PREDICATE_IFS;

    static {
        ImmutableMap.Builder<String, BiFunction<String, String, String>> builder = ImmutableMap.builder();
        builder.put("K.K", (c, s) -> "if true");
        builder.put("K.KItem", (c, s) -> "if true");
        builder.put("INT.Int", (c, s) -> "if m.IsInt(" + c + ")");
        builder.put("FLOAT.Float", (c, s) -> "if m.IsFloat(" + c + ")");
        builder.put("STRING.String", (c, s) -> "if m.IsString(" + c + ")");
        builder.put("BYTES.Bytes", (c, s) -> "if m.IsBytes(" + c + ")");
        builder.put("BUFFER.StringBuffer", (c, s) -> "if m.IsStringBuffer(" + c + ")");
        builder.put("BOOL.Bool", (c, s) -> "if m.IsBool(" + c + ")");
        builder.put("MINT.MInt", (c, s) -> "if m.IsMint(" + c + ")");
        builder.put("MAP.Map", (c, s) -> "if i.Model.IsMap(" + c + ", " + s + ")");
        builder.put("SET.Set", (c, s) -> "if i.Model.IsSet(" + c + ", " + s + ")");
        builder.put("LIST.List", (c, s) -> "if i.Model.IsList(" + c + ", " + s + ")");
        builder.put("ARRAY.Array", (c, s) -> "if i.Model.IsArray(" + c + ", " + s + ")");
        PREDICATE_IFS = builder.build();
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

    /**
     * Sort var hooks for basic types.
     */
    public static final ImmutableMap<String, String> SORT_VAR_HOOKS_1;

    /**
     * Sort var hooks for collection types, the sort needs to be checked too.
     */
    public static final ImmutableMap<String, String> SORT_VAR_HOOKS_2;

    static {
        ImmutableMap.Builder<String, String> builder = ImmutableMap.builder();
        builder.put("BOOL.Bool", "if %1$s, t := %2$s.(*m.Bool); t");
        builder.put("MINT.MInt", "if %1$s, t := %2$s.(m.MInt); t.");
        builder.put("INT.Int", "if %1$s, t := %2$s.(*m.Int); t");
        builder.put("FLOAT.Float", "if %1$s, t := %2$s.(*m.Float); t");
        builder.put("STRING.String", "if %1$s, t := %2$s.(*m.String); t");
        builder.put("BYTES.Bytes", "if %1$s, t := %2$s.(*m.Bytes); t");
        builder.put("BUFFER.StringBuffer", "if %1$s, t := %2$s.(*m.StringBuffer); t");
        SORT_VAR_HOOKS_1 = builder.build();
        builder = ImmutableMap.builder();
        builder.put("LIST.List", "if %1$s, t := %2$s.(*m.List); t && %1$s.Sort == m.%3$s");
        builder.put("ARRAY.Array", "if %1$s, t := %2$s.(*m.Array); t && %1$s.Sort == m.%3$s");
        builder.put("MAP.Map", "if %1$s, t := %2$s.(*m.Map); t && %1$s.Sort == m.%3$s");
        builder.put("SET.Set", "if %1$s, t := %2$s.(*m.Set); t && %1$s.Sort == m.%3$s");
        SORT_VAR_HOOKS_2 = builder.build();
    }

    public static final ImmutableSet<String> LHS_KVARIABLE_HOOKS;

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
        LHS_KVARIABLE_HOOKS = builder.build();
    }
}
