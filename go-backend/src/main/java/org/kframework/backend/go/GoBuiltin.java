package org.kframework.backend.go;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.ImmutableSet;
import org.kframework.kore.Sort;

import java.util.function.Function;

public class GoBuiltin {

    public static final ImmutableSet<String> HOOK_NAMESPACES;

    static {
        ImmutableSet.Builder<String> builder = ImmutableSet.builder();
        builder.add("BOOL").add("FLOAT").add("INT").add("IO").add("K").add("KEQUAL").add("KREFLECTION").add("LIST");
        builder.add("MAP").add("MINT").add("SET").add("STRING").add("ARRAY").add("BUFFER").add("BYTES");
        HOOK_NAMESPACES = builder.build();
    }


    public static final ImmutableMap<String, Function<Sort, String>> PREDICATE_RULES;
    private static final String RETURN_BOOL_TRUE_BLOCK = " {\n\t\treturn Bool(true)\n\t}";

    static {
        ImmutableMap.Builder<String, Function<Sort, String>> builder = ImmutableMap.builder();
        builder.put("K.K", s -> "return Bool(true)");
        builder.put("K.KItem", s -> "// almost certainly incomplete, original: [_] -> [Bool true] | _ -> [Bool false] ???????");
        builder.put("INT.Int", s -> "if _, t := c.(Int); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("FLOAT.Float", s -> "if _, t := c.(Float); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("STRING.String", s -> "if _, t := c.(String); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("BYTES.Bytes", s -> "if _, t := c.(Bytes); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("BUFFER.StringBuffer", s -> "if _, t := c.(StringBuffer); t" + RETURN_BOOL_TRUE_BLOCK);
        builder.put("BOOL.Bool", s -> "if _, t := c.(Bool); t " + RETURN_BOOL_TRUE_BLOCK);
        builder.put("MINT.MInt", s -> "if _, t := c.(MInt); t " + RETURN_BOOL_TRUE_BLOCK);
        builder.put("MAP.Map", s -> "if mp, t := c.(Map); t && mp.Sort == " + GoStringUtil.sortVariableName(s) + RETURN_BOOL_TRUE_BLOCK);
        builder.put("SET.Set", s -> "if set, t := c.(Set); t && set.Sort == " + GoStringUtil.sortVariableName(s) + RETURN_BOOL_TRUE_BLOCK);
        builder.put("LIST.List", s -> "if list, t := c.(List); t && list.Sort == " + GoStringUtil.sortVariableName(s) + RETURN_BOOL_TRUE_BLOCK);
        builder.put("ARRAY.Array", s -> "if arr, t := c.(Array); t && arr.Sort == " + GoStringUtil.sortVariableName(s) + RETURN_BOOL_TRUE_BLOCK);
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

    public static final ImmutableMap<String, String> GO_SORT_VAR_HOOKS;
    static {
        ImmutableMap.Builder<String, String> builder = ImmutableMap.builder();
        builder.put("BOOL.Bool", "if %1$s, t := %2$s.(Bool); t");
        builder.put("MINT.MInt", "if %1$s, t := %2$s.(MInt); t.");
        builder.put("INT.Int",   "if %1$s, t := %2$s.(Int); t");
        builder.put("FLOAT.Float",  "if %1$s, t := %2$s.(Float); t");
        builder.put("STRING.String", "if %1$s, t := %2$s.(String);");
        builder.put("BYTES.Bytes", "if %1$s, t := %2$s.(Bytes); t.(Bytes); t");
        builder.put("BUFFER.StringBuffer",  "if %1$s, t := %2$s.(StringBuffer); t");
        builder.put("LIST.List", "if %1$s, t := %2$s.(List); t && %1$s.Sort == %3$s");
        builder.put("ARRAY.Array",  "if %1$s, t := %2$s.(Array); t && %1$s.Sort == %3$s");
        builder.put("MAP.Map", "if %1$s, t := %2$s.(Map); t && %1$s.Sort == %3$s");
        builder.put("SET.Set",  "if %1$s, t := %2$s.(Set); t && %1$s.Sort == %3$s");
        GO_SORT_VAR_HOOKS = builder.build();
    }
}
