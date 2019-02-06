package org.kframework.backend.go.codegen;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.ImmutableSet;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.kore.Sort;
import org.kframework.utils.StringUtil;

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

    public static final ImmutableMap<String, Function<String, String>> PREDICATE_RULES;
    private static final String RETURN_BOOL_TRUE_BLOCK = " {\n\t\treturn m.Bool(true), nil\n\t}";

    static {
        ImmutableMap.Builder<String, Function<String, String>> builder = ImmutableMap.builder();
        builder.put("K.K", s -> "return m.Bool(true), nil");
        builder.put("K.KItem", s -> "return m.Bool(true), nil");
        builder.put("INT.Int", s -> "if _, t := c.(m.Int); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("FLOAT.Float", s -> "if _, t := c.(m.Float); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("STRING.String", s -> "if _, t := c.(m.String); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("BYTES.Bytes", s -> "if _, t := c.(m.Bytes); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("BUFFER.StringBuffer", s -> "if _, t := c.(m.StringBuffer); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("BOOL.Bool", s -> "if _, t := c.(m.Bool); t " + RETURN_BOOL_TRUE_BLOCK);
        builder.put("MINT.MInt", s -> "if _, t := c.(m.MInt); t " + RETURN_BOOL_TRUE_BLOCK);
        builder.put("MAP.Map", s -> "if mp, t := c.(m.Map); t && mp.Sort == " + s + RETURN_BOOL_TRUE_BLOCK);
        builder.put("SET.Set", s -> "if set, t := c.(m.Set); t && set.Sort == " + s + RETURN_BOOL_TRUE_BLOCK);
        builder.put("LIST.List", s -> "if list, t := c.(m.List); t && list.Sort == " + s + RETURN_BOOL_TRUE_BLOCK);
        builder.put("ARRAY.Array", s -> "if arr, t := c.(m.Array); t && arr.Sort == " + s + RETURN_BOOL_TRUE_BLOCK);
        PREDICATE_RULES = builder.build();
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


    public static final ImmutableMap<String, Function<String, String>> GO_SORT_TOKEN_HOOKS;
    static {
        ImmutableMap.Builder<String, Function<String, String>> builder = ImmutableMap.builder();
        builder.put("BOOL.Bool", s -> "m.Bool(" + s + ")");
        builder.put("MINT.MInt", s -> "m.MInt(" + s + ")");
//        builder.put("MINT.MInt", s -> {
//            MIntBuiltin m = MIntBuiltin.of(s);
//            return "(MInt (" + m.precision() + ", Z.of_string \"" + m.value() + "))";
//        });
        builder.put("INT.Int", s -> "m.Int(" + s + ")");
        builder.put("FLOAT.Float", s -> "m.Float(" + s + ")");
//        builder.put("FLOAT.Float", s -> {
//            FloatBuiltin f = FloatBuiltin.of(s);
//            return "(round_to_range(Float ((Gmp.FR.from_string_prec_base " + f.precision() + " Gmp.GMP_RNDN 10 \"" + f.value() + "\"), " + f.exponent() + ", " + f.precision() + ")))";
//        });
        builder.put("STRING.String", s -> "m.String (" + GoStringUtil.enquoteString(StringUtil.unquoteKString(s)) + ")");
        builder.put("BYTES.Bytes", s -> "m.Bytes([]byte(" + GoStringUtil.enquoteString(StringUtil.unquoteKString(s)) + "))");
        builder.put("BUFFER.StringBuffer", s -> "m.StringBuffer{}");
        GO_SORT_TOKEN_HOOKS = builder.build();
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
        builder.put("BOOL.Bool", "if %1$s, t := %2$s.(m.Bool); t");
        builder.put("MINT.MInt", "if %1$s, t := %2$s.(m.MInt); t.");
        builder.put("INT.Int",   "if %1$s, t := %2$s.(m.Int); t");
        builder.put("FLOAT.Float",  "if %1$s, t := %2$s.(m.Float); t");
        builder.put("STRING.String", "if %1$s, t := %2$s.(m.String); t");
        builder.put("BYTES.Bytes", "if %1$s, t := %2$s.(m.Bytes); t");
        builder.put("BUFFER.StringBuffer",  "if %1$s, t := %2$s.(m.StringBuffer); t");
        SORT_VAR_HOOKS_1 = builder.build();
        builder = ImmutableMap.builder();
        builder.put("LIST.List", "if %1$s, t := %2$s.(m.List); t && %1$s.Sort == m.%3$s");
        builder.put("ARRAY.Array",  "if %1$s, t := %2$s.(m.Array); t && %1$s.Sort == m.%3$s");
        builder.put("MAP.Map", "if %1$s, t := %2$s.(m.Map); t && %1$s.Sort == m.%3$s");
        builder.put("SET.Set",  "if %1$s, t := %2$s.(m.Set); t && %1$s.Sort == m.%3$s");
        SORT_VAR_HOOKS_2 = builder.build();
    }
}
