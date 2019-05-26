// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.ConstantKTokens;
import org.kframework.backend.go.strings.GoStringBuilder;

import java.util.Map;

public class ConstantsGen {

    private final GoPackageManager packageManager;
    private final ConstantKTokens constants;

    public ConstantsGen(GoPackageManager packageManager, ConstantKTokens constants) {
        this.packageManager = packageManager;
        this.constants = constants;
    }

    public String generate() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append(packageManager.goGeneratedFileComment).append("\n\n");
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append("\n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        sb.appendIndentedLine("// Int constants");
        for (Map.Entry<String, String> entry : constants.intConstants.entrySet()) {
            sb.appendIndentedLine("var ", entry.getKey(), " = ", entry.getValue());
        }
        sb.newLine();

        sb.appendIndentedLine("// String constants");
        for (Map.Entry<String, String> entry : constants.stringConstants.entrySet()) {
            sb.appendIndentedLine("var ", entry.getKey(), " = ", entry.getValue());
        }
        sb.newLine();

        sb.appendIndentedLine("// KToken constants");
        for (Map.Entry<String, String> entry : constants.tokenConstants.entrySet()) {
            sb.appendIndentedLine("var ", entry.getKey(), " = ", entry.getValue());
        }
        sb.newLine();

        return sb.toString();
    }

}
