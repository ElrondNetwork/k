package org.kframework.backend.go.strings;

import org.kframework.kore.KLabel;
import org.kframework.kore.Sort;

/**
 * Provides names for variables, functions, etc. that are as close as possible to the original.
 * Easier for debugging, will result in some warnings in Go.
 */
public class GoNameProviderDebug implements GoNameProvider {

    @Override
    public String klabelVariableName(KLabel klabel) {
        StringBuilder sb = new StringBuilder();
        sb.append("lbl");
        appendAlphanumericEncodedString(sb, klabel.name());
        return sb.toString();
    }


    @Override
    public String sortVariableName(Sort sort) {
        StringBuilder sb = new StringBuilder();
        sb.append("sort");
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
            if (charAt < 128 && ENCODING[charAt] != null) {
                sb.append(ENCODING[charAt]);
            } else {
                sb.append(String.format("_%04x_", charAt));
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



}
