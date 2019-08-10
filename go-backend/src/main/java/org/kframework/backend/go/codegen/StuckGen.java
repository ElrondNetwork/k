// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import org.kframework.backend.go.codegen.inline.RuleLhsMatchWriter;
import org.kframework.backend.go.codegen.rules.RuleWriter;
import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionInfo;
import org.kframework.backend.go.model.RuleInfo;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;

import java.util.Collections;

public class StuckGen {

    private final DefinitionData data;
    private final GoPackageManager packageManager;
    private final RuleWriter ruleWriter;
    private static final int stuckRuleNumber = -1; // this only affects comments and traces

    public StuckGen(DefinitionData data, GoPackageManager packageManager, GoNameProvider nameProvider, RuleLhsMatchWriter matchWriter) {
        this.data = data;
        this.packageManager = packageManager;
        this.ruleWriter = new RuleWriter(data, nameProvider, matchWriter);
    }

    public String generateStuck() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append(packageManager.goGeneratedFileComment).append("\n\n");
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");

        int maxNrVars = 0;
        int maxNrBoolVars = 0;

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        sb.append("func (i *Interpreter) makeStuck(c m.KReference, config m.KReference) (m.KReference, error)").beginBlock();
        sb.appendIndentedLine("var v [makeStuckFuncNrVars]KReference");
        sb.appendIndentedLine("var bv [makeStuckFuncNrBoolVars]bool");

        if (data.makeStuck != null) {
            sb.appendIndentedLine("matched := false");

            RuleInfo ruleInfo = ruleWriter.writeRule(
                    Collections.singletonMap(stuckRuleNumber, data.makeStuck),
                    sb, RuleType.REGULAR,
                    FunctionInfo.systemFunctionInfo("makeStuck", 1));
            if (ruleInfo.maxNrVars > maxNrVars) {
                maxNrVars = ruleInfo.maxNrVars;
            }
            if (ruleInfo.maxNrBoolVars > maxNrBoolVars) {
                maxNrBoolVars = ruleInfo.maxNrBoolVars;
            }
        }
        sb.appendIndentedLine("doNothingWithVars(len(v), len(bv))"); // just to stop Go complaining about unused vars, never gets called
        sb.appendIndentedLine("return c, nil");
        sb.endOneBlock().newLine();

        sb.append("func (i *Interpreter) makeUnstuck(c m.KReference, config m.KReference) (m.KReference, error)").beginBlock();
        sb.appendIndentedLine("var v [makeStuckFuncNrVars]KReference");
        sb.appendIndentedLine("var bv [makeStuckFuncNrBoolVars]bool");

        if (data.makeUnstuck != null) {
            sb.appendIndentedLine("matched := false");
            RuleInfo ruleInfo = ruleWriter.writeRule(
                    Collections.singletonMap(stuckRuleNumber, data.makeUnstuck),
                    sb, RuleType.REGULAR,
                    FunctionInfo.systemFunctionInfo("makeUnstuck", 1));
            if (ruleInfo.maxNrVars > maxNrVars) {
                maxNrVars = ruleInfo.maxNrVars;
            }
            if (ruleInfo.maxNrBoolVars > maxNrBoolVars) {
                maxNrBoolVars = ruleInfo.maxNrBoolVars;
            }
        }
        sb.appendIndentedLine("doNothingWithVars(len(v), len(bv))"); // just to stop Go complaining about unused vars, never gets called
        sb.appendIndentedLine("return c, nil");
        sb.endOneBlock().newLine();

        sb.appendIndentedLine("const makeStuckFuncNrVars = ", Integer.toString(maxNrVars));
        sb.appendIndentedLine("const makeStuckFuncNrBoolVars = ", Integer.toString(maxNrBoolVars));
        sb.newLine();

        return sb.toString();
    }

}