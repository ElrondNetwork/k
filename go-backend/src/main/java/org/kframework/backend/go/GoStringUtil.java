package org.kframework.backend.go;

import com.google.common.collect.ImmutableMap;
import org.kframework.definition.Rule;
import org.kframework.kore.KLabel;
import org.kframework.kore.Sort;
import org.kframework.unparser.ToKast;

import java.util.regex.Pattern;

class GoStringUtil {

    public static final Pattern identChar = Pattern.compile("[A-Za-z0-9_]");

    private static final String[] ASCII_READABLE_ENCODING = new String[]{
            null,// 00
            null,// 01
            null,// 02
            null,// 03
            null,// 04
            null,// 05
            null,// 06
            null,// 07
            null,// 08
            null,// 09
            null,// 0a
            null,// 0b
            null,// 0c
            null,// 0d
            null,// 0e
            null,// 0f
            null,// 10
            null,// 11
            null,// 12
            null,// 13
            null,// 14
            null,// 15
            null,// 16
            null,// 17
            null,// 18
            null,// 19
            null,// 1a
            null,// 1b
            null,// 1c
            null,// 1d
            null,// 1e
            null,// 1f
            "_Spce_",// 20
            "_Bang_",// 21
            "_Quot_",// 22
            "_Hash_",// 23
            "_Dolr_",// 24
            "_Perc_",// 25
            "_Amps_",// 26
            "_Apos_",// 27
            "_LPar_",// 28
            "_RPar_",// 29
            "_Star_",// 2a
            "_Plus_",// 2b
            "_Comm_",// 2c
            "_Hyph_",// 2d
            "_Stop_",// 2e
            "_Slsh_",// 2f
            "0",// 30
            "1",// 31
            "2",// 32
            "3",// 33
            "4",// 34
            "5",// 35
            "6",// 36
            "7",// 37
            "8",// 38
            "9",// 39
            "_Coln_",// 3a
            "_SCln_",// 3b
            "_Lt_",// 3c
            "_Eqls_",// 3d
            "_Gt_",// 3e
            "_Ques_",// 3f
            "_At_",// 40
            "A",// 41
            "B",// 42
            "C",// 43
            "D",// 44
            "E",// 45
            "F",// 46
            "G",// 47
            "H",// 48
            "I",// 49
            "J",// 4a
            "K",// 4b
            "L",// 4c
            "M",// 4d
            "N",// 4e
            "O",// 4f
            "P",// 50
            "Q",// 51
            "R",// 52
            "S",// 53
            "T",// 54
            "U",// 55
            "V",// 56
            "W",// 57
            "X",// 58
            "Y",// 59
            "Z",// 5a
            "_LSqB_",// 5b
            "_Bash_",// 5c
            "_RSqB_",// 5d
            "_Xor_",// 5e
            "_",// 5f
            "_BQuo_",// 60
            "a",// 61
            "b",// 62
            "c",// 63
            "d",// 64
            "e",// 65
            "f",// 66
            "g",// 67
            "h",// 68
            "i",// 69
            "j",// 6a
            "k",// 6b
            "l",// 6c
            "m",// 6d
            "n",// 6e
            "o",// 6f
            "p",// 70
            "q",// 71
            "r",// 72
            "s",// 73
            "t",// 74
            "u",// 75
            "v",// 76
            "w",// 77
            "x",// 78
            "y",// 79
            "z",// 7a
            "_LBra_",// 7b
            "_Pipe_",// 7c
            "_RBra_",// 7d
            "_Tild_",// 7e
            null// 7f
    };

    static void appendAlphanumericEncodedString(StringBuilder sb, String name) {
        for (int i = 0; i < name.length(); i++) {
            int charAt = (int) name.charAt(i);
            if (charAt < 128 && ASCII_READABLE_ENCODING[charAt] != null) {
                sb.append(ASCII_READABLE_ENCODING[charAt]);
            } else {
                sb.append(String.format("_%04x_", charAt));
            }
        }
    }

    static void appendKlabelVariableName(StringBuilder sb, KLabel klabel) {
        sb.append("lbl");
        appendAlphanumericEncodedString(sb, klabel.name());
    }

    static void appendSortVariableName(StringBuilder sb, Sort sort) {
        sb.append("sort");
        appendAlphanumericEncodedString(sb, sort.name());
    }

    static String sortVariableName(Sort sort) {
        StringBuilder sb = new StringBuilder();
        appendSortVariableName(sb, sort);
        return sb.toString();
    }

    static void appendFunctionName(StringBuilder sb, KLabel lbl) {
        sb.append("eval");
        appendAlphanumericEncodedString(sb, lbl.name());
    }

    static void appendMemoFunctionName(StringBuilder sb, KLabel lbl) {
        sb.append("memo");
        appendAlphanumericEncodedString(sb, lbl.name());
    }

    static String functionName(KLabel lbl) {
        StringBuilder sb = new StringBuilder();
        appendFunctionName(sb, lbl);
        return sb.toString();
    }

    static String memoFunctionName(KLabel lbl) {
        StringBuilder sb = new StringBuilder();
        appendMemoFunctionName(sb, lbl);
        return sb.toString();
    }

    static void appendConstFunctionName(StringBuilder sb, KLabel lbl) {
        sb.append("const");
        appendAlphanumericEncodedString(sb, lbl.name());
    }

    static String variableName(String varName) {
        if (varName.equals("_")) {
            return "_";
        }
        StringBuilder sb = new StringBuilder();
        sb.append("var");
        appendAlphanumericEncodedString(sb, varName);
        return sb.toString();
    }

    public static String enquoteString(String value) {
        char delimiter = '"';
        final int length = value.length();
        StringBuilder result = new StringBuilder();
        result.append(delimiter);
        for (int offset = 0, codepoint; offset < length; offset += Character.charCount(codepoint)) {
            codepoint = value.codePointAt(offset);
            if (codepoint > 0xFF) {
                // unicode
                result.append((char) codepoint);
            } else if (codepoint == delimiter) {
                result.append("\\").append(delimiter);
            } else if (codepoint == '\\') {
                result.append("\\\\");
            } else if (codepoint == '\n') {
                result.append("\\n");
            } else if (codepoint == '\t') {
                result.append("\\t");
            } else if (codepoint == '\r') {
                result.append("\\r");
            } else if (codepoint == '\b') {
                result.append("\\b");
            } else if (codepoint >= 32 && codepoint < 127) {
                result.append((char) codepoint);
            } else if (codepoint <= 0xff) {
                if (codepoint < 10) {
                    result.append("\\00");
                    result.append(codepoint);
                } else if (codepoint < 100) {
                    result.append("\\0");
                    result.append(codepoint);
                } else {
                    result.append("\\");
                    result.append(codepoint);
                }
            }
        }
        result.append(delimiter);
        return result.toString();
    }

    private static final ImmutableMap<String, String> GO_RESERVED_TO_REPLACEMENTS;

    static {
        ImmutableMap.Builder<String, String> builder = ImmutableMap.builder();
        builder.put("range", "lrange");
        GO_RESERVED_TO_REPLACEMENTS = builder.build();
    }

    static void appendHookMethodName(StringBuilder sb, String objectName, String methodName) {
        sb.append(objectName.toLowerCase());
        sb.append("Hooks.");
        if (GO_RESERVED_TO_REPLACEMENTS.containsKey(methodName)) {
            // replace Go reserved words with some non-conflicting ones
            methodName = GO_RESERVED_TO_REPLACEMENTS.get(methodName);
        }
        sb.append(methodName);
    }

    private static final String EVAL_ARG_NAME = "c";

    /**
     * Helps create function definitions and function calls with number of arguments given by the arity parameter.
     *
     * @param arity             number of parameters
     * @param cParamDeclaration output section of function declaration here
     * @param cParamCall        output section of function call here
     */
    static void createParameterDeclarationAndCall(int arity, StringBuilder cParamDeclaration, StringBuilder cParamCall) {
        if (arity == 0) {
            // nothing needs to be added
        } else if (arity == 1) {
            cParamDeclaration.append(EVAL_ARG_NAME).append(" K,");
            cParamCall.append(EVAL_ARG_NAME).append(", ");
        } else if (arity > 1) {
            for (int i = 1; i <= arity; i++) {
                cParamDeclaration.append(EVAL_ARG_NAME).append(i).append(" K, ");
                cParamCall.append(EVAL_ARG_NAME).append(i).append(", ");
            }
        }
    }

    static void appendRuleComment(StringBuilder sb, Rule r) {
        sb.append("\t// {| rule ");
        sb.append(ToKast.apply(r.body()).replace("|}", "| )"));
        sb.append(" requires ");
        sb.append(ToKast.apply(r.requires()).replace("|)", "| )"));
        sb.append(" ensures ");
        sb.append(ToKast.apply(r.ensures()).replace("|)", "| )"));
        sb.append(" ");
        sb.append(r.att().toString().replace("|)", "| )"));
        sb.append(" |}\n");
    }

}
