// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.strings;

import org.kframework.kore.KLabel;
import org.kframework.kore.Sort;

/**
 * Provides proper names for variables, functions, etc., according to Go naming guidelines
 * (no underscores in names).
 */
public class GoNameProviderProper implements GoNameProvider {

    @Override
    public String klabelVariableName(KLabel klabel) {
        StringBuilder sb = new StringBuilder();
        sb.append("Lbl");
        appendAlphanumericEncodedString(sb, klabel.name());
        return sb.toString();
    }

    @Override
    public String sortVariableName(Sort sort) {
        StringBuilder sb = new StringBuilder();
        sb.append("Sort");
        appendAlphanumericEncodedString(sb, sort.name());
        return sb.toString();
    }

    @Override
    public String evalFunctionName(KLabel lbl) {
        StringBuilder sb = new StringBuilder();
        sb.append("eval");
        appendAlphanumericEncodedString(sb, lbl.name());
        return sb.toString();
    }

    @Override
    public String constFunctionName(KLabel lbl) {
        StringBuilder sb = new StringBuilder();
        sb.append("const");
        appendAlphanumericEncodedString(sb, lbl.name());
        return sb.toString();
    }

    @Override
    public String memoFunctionName(KLabel lbl) {
        StringBuilder sb = new StringBuilder();
        sb.append("memo");
        appendAlphanumericEncodedString(sb, lbl.name());
        return sb.toString();
    }

    @Override
    public String memoTableName(KLabel lbl) {
        StringBuilder sb = new StringBuilder();
        sb.append("memoTable");
        appendAlphanumericEncodedString(sb, lbl.name());
        return sb.toString();
    }

    @Override
    public String ruleVariableName(String varName) {
        if (varName.equals("_")) {
            return "_";
        }
        StringBuilder sb = new StringBuilder();
        sb.append("var");
        appendAlphanumericEncodedString(sb, varName);
        return sb.toString();
    }

    private static void appendAlphanumericEncodedString(StringBuilder sb, String name) {
        for (int i = 0; i < name.length(); i++) {
            int charAt = (int) name.charAt(i);
            if (i == 0) {
                charAt = Character.toUpperCase(charAt);
            }
            if (i > 0 && name.charAt(i-1) == 'I' && charAt == 'd') {
                charAt = 'D'; // Id => ID
            }
            if (charAt < 128 && ENCODING[charAt] != null) {
                sb.append(ENCODING[charAt]);
            } else {
                sb.append(String.format("X%04x", charAt));
            }
        }
    }

    private static final String[] ENCODING = new String[]{
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
            "Xspace",// 20
            "Xbang",// 21
            "Xquote",// 22
            "Xhash",// 23
            "Xdolr",// 24
            "Xpercent",// 25
            "Xamps",// 26
            "Xapos",// 27
            "Xlparen",// 28
            "Xrparen",// 29
            "Xstar",// 2a
            "Xplus",// 2b
            "Xcomma",// 2c
            "Xhyphen",// 2d
            "Xdot",// 2e
            "Xslash",// 2f
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
            "Xcolon",// 3a :
            "Xscolon",// 3b ;
            "Xlt",// 3c <
            "Xeq",// 3d =
            "Xgt",// 3e >
            "Xques",// 3f ?
            "Xat",// 40 @
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
            "Xlsqb",// 5b
            "Xbash",// 5c
            "Xrsqb",// 5d
            "Xxor",// 5e
            "Xu",// 5f _
            "Xbquote",// 60 "
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
            "Xlbracket",// 7b
            "Xpipe",// 7c
            "Xrbracket",// 7d
            "Xtilde",// 7e
            null// 7f
    };



}

