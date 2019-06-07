// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import org.kframework.backend.go.codegen.rules.RuleWriter;
import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.RuleInfo;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;

public class StuckGen {

    private final DefinitionData data;
    private final GoPackageManager packageManager;
    private final RuleWriter ruleWriter;
    private static final int stuckRuleNumber = -1; // this only affects comments and traces

    public StuckGen(DefinitionData data, GoPackageManager packageManager, GoNameProvider nameProvider) {
        this.data = data;
        this.packageManager = packageManager;
        this.ruleWriter = new RuleWriter(data, nameProvider);
    }

    public String generateStuck() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append(packageManager.goGeneratedFileComment).append("\n\n");
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        sb.append("func (i *Interpreter) makeStuck(c m.K, config m.K) (m.K, error)").beginBlock();
        if (data.makeStuck != null) {
            RuleInfo ruleInfo = ruleWriter.writeRule(
                    data.makeStuck, sb, RuleType.REGULAR, stuckRuleNumber,
                    "makeStuck", new FunctionParams(1));
            assert !ruleInfo.alwaysMatches();
        }
        sb.appendIndentedLine("return c, nil");
        sb.endOneBlock().newLine();

        sb.append("func (i *Interpreter) makeUnstuck(c m.K, config m.K) (m.K, error)").beginBlock();
        if (data.makeUnstuck != null) {
            RuleInfo ruleInfo = ruleWriter.writeRule(
                    data.makeUnstuck, sb, RuleType.REGULAR, stuckRuleNumber,
                    "makeUnstuck", new FunctionParams(1));
            assert !ruleInfo.alwaysMatches();
        }
        sb.appendIndentedLine("return c, nil");
        sb.endOneBlock().newLine();

        return sb.toString();
    }

}
