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
}
