package org.kframework.backend.go.strings;

import org.kframework.definition.Rule;
import org.kframework.kore.Sort;
import org.kframework.unparser.ToKast;

public class GoStringUtil {


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
            } else if (codepoint == 0) {
                result.append("\\x00");
            } else if (codepoint <= 0xFF) {
                result.append("\\x");
                result.append(Integer.toHexString(codepoint));

            }
        }
        result.append(delimiter);
        return result.toString();
    }

    public static String packageName(String dirName) {
        return dirName.replaceAll("[\\-_.]", "").toLowerCase();
    }

    public static void appendRuleComment(GoStringBuilder sb, Rule r) {
        sb.append("{| rule ");
        sb.append(ToKast.apply(r.body()).replace("|}", "| )"));
        sb.append(" requires ");
        sb.append(ToKast.apply(r.requires()).replace("|)", "| )"));
        sb.append(" ensures ");
        sb.append(ToKast.apply(r.ensures()).replace("|)", "| )"));
        sb.append(" ");
        sb.append(r.att().toString().replace("|)", "| )"));
        sb.append(" |}");
    }

    public static String enquotedRuleComment(Rule r) {
        GoStringBuilder sb = new GoStringBuilder();
        appendRuleComment(sb, r);
        return enquoteString(sb.toString());
    }

    @Deprecated
    public static String sortVariableName(Sort sort) {
        return new GoNameProviderDebug().sortVariableName(sort);
    }

}
