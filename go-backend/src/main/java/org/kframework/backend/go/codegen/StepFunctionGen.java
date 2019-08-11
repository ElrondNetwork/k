// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import com.google.common.collect.ComparisonChain;
import org.kframework.backend.go.codegen.inline.RuleLhsMatchWriter;
import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.backend.go.codegen.rules.RuleWriter;
import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionInfo;
import org.kframework.backend.go.model.RuleCounter;
import org.kframework.backend.go.model.RuleInfo;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.definition.Rule;
import org.kframework.kil.Attribute;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import static org.kframework.Collections.*;

public class StepFunctionGen {

    private final DefinitionData data;
    private final GoPackageManager packageManager;
    private final RuleWriter ruleWriter;

    private final List<Rule> sortedRules;
    private final Map<Integer, Rule> stepRules = new HashMap<>();
    private final Map<Integer, Rule> lookupRules = new HashMap<>();

    private final RuleCounter ruleCounter = new RuleCounter();

    GoStringBuilder stepFuncSb = new GoStringBuilder();
    GoStringBuilder stepRhsSb = new GoStringBuilder();

    public StepFunctionGen(DefinitionData data, GoPackageManager packageManager, GoNameProvider nameProvider, RuleLhsMatchWriter matchWriter) {
        this.data = data;
        this.packageManager = packageManager;
        this.ruleWriter = new RuleWriter(data, nameProvider, matchWriter);
        List<Rule> unsortedRules = stream(data.mainModule.rules()).collect(Collectors.toList());
//        if (options.reverse) {
//            Collections.reverse(unsortedRules);
//        }
        sortedRules = unsortedRules.stream()
                .sorted(this::sortRules)
                .filter(r -> !data.functionRules.values().contains(r) && !r.att().contains(Attribute.MACRO_KEY) && !r.att().contains(Attribute.ALIAS_KEY) && !r.att().contains(Attribute.ANYWHERE_KEY))
                .collect(Collectors.toList());

        Map<Boolean, List<Rule>> groupedByLookup = sortedRules.stream()
                .collect(Collectors.groupingBy(RuleLhsTreeWriter::hasLookups));
        if (groupedByLookup.containsKey(false)) {
            for (Rule r : groupedByLookup.get(false)) {
                int ruleNum = ruleCounter.consumeRuleIndex();
                stepRules.put(ruleNum, r);
            }
        }
        if (groupedByLookup.containsKey(true)) {
            for (Rule r : groupedByLookup.get(true)) {
                int ruleNum = ruleCounter.consumeRuleIndex();
                lookupRules.put(ruleNum, r);
            }
        }
    }

    public void generateStepFunctionCode() {
        stepFuncSb.append(packageManager.goGeneratedFileComment).append("\n\n");
        stepFuncSb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");
        stepFuncSb.append("import (\n");
        stepFuncSb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        stepFuncSb.append(")\n\n");

        stepRhsSb.append(packageManager.goGeneratedFileComment).append("\n\n");
        stepRhsSb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");
        stepRhsSb.append("import (\n");
        stepRhsSb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        stepRhsSb.append(")\n\n");

        stepFuncSb.append("func (i *Interpreter) step(c m.KReference) (m.KReference, error)").beginBlock();
        stepFuncSb.appendIndentedLine("config := c");
        stepFuncSb.appendIndentedLine("matched := false");
        stepFuncSb.appendIndentedLine("v := i.stepTempVars");
        stepFuncSb.appendIndentedLine("bv := i.stepTempBoolVars");

        RuleInfo ruleInfo = ruleWriter.writeRule(
                stepRules,
                stepFuncSb, stepRhsSb,
                RuleType.REGULAR,
                FunctionInfo.systemFunctionInfo("step", 1));

        stepFuncSb.writeIndent().append("return i.stepLookups(c, config, -1)\n");
        stepFuncSb.endOneBlock().newLine();

        stepFuncSb.append("func (i *Interpreter) stepLookups(c m.KReference, config m.KReference, guard int) (m.KReference, error)").beginBlock();
        stepFuncSb.appendIndentedLine("matched := false");
        stepFuncSb.appendIndentedLine("v := i.stepTempVars");
        stepFuncSb.appendIndentedLine("bv := i.stepTempBoolVars");

        RuleInfo lookupRuleInfo = ruleWriter.writeRule(
                lookupRules,
                stepFuncSb, stepRhsSb,
                RuleType.REGULAR,
                FunctionInfo.systemFunctionInfo("stepLookups", 1));

        if (lookupRuleInfo.maxNrVars > ruleInfo.maxNrVars) {
            ruleInfo.maxNrVars = lookupRuleInfo.maxNrVars;
        }
        if (lookupRuleInfo.maxNrBoolVars > ruleInfo.maxNrBoolVars) {
            ruleInfo.maxNrBoolVars = lookupRuleInfo.maxNrBoolVars;
        }

        stepFuncSb.appendIndentedLine("return c, noStep");
        stepFuncSb.endOneBlock().newLine();

        stepFuncSb.appendIndentedLine("// stepMaxVarCount indicates the maximum number of variables required by a rule");
        stepFuncSb.appendIndentedLine("// needed to initialize the step variables (i.stepTempVars) slice");
        stepFuncSb.appendIndentedLine("const stepMaxVarCount = " + ruleInfo.maxNrVars);
        stepFuncSb.newLine();

        stepFuncSb.appendIndentedLine("// stepMaxBoolVarCount indicates the maximum number of boolean variables required by a rule");
        stepFuncSb.appendIndentedLine("// needed to initialize the boolean variables (i.stepTempBoolVars) slice");
        stepFuncSb.appendIndentedLine("const stepMaxBoolVarCount = " + ruleInfo.maxNrBoolVars);
        stepFuncSb.newLine();
    }

    public String outputStepFunctionCode() {
        return stepFuncSb.toString();
    }

    public String outputStepRhsCode() {
        return stepRhsSb.toString();
    }

    private int sortRules(Rule r1, Rule r2) {
        return ComparisonChain.start()
                .compareTrueFirst(r1.att().contains("structural"), r2.att().contains("structural"))
                .compareFalseFirst(r1.att().contains("owise"), r2.att().contains("owise"))
                //.compareFalseFirst(indexesPoorly(r1), indexesPoorly(r2))
                .result();
    }

}
