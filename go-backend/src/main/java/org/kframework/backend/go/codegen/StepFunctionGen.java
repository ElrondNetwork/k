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

    public String generateStep() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append(packageManager.goGeneratedFileComment).append("\n\n");
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        sb.append("func (i *Interpreter) step(c m.KReference) (m.KReference, error)").beginBlock();
        sb.appendIndentedLine("config := c");
        sb.appendIndentedLine("matched := false");
        sb.appendIndentedLine("v := i.stepTempVars");
        sb.appendIndentedLine("bv := i.stepTempBoolVars");

        RuleInfo ruleInfo = ruleWriter.writeRule(
                stepRules, sb, RuleType.REGULAR,
                FunctionInfo.systemFunctionInfo("step", 1));

        sb.writeIndent().append("return i.stepLookups(c, config, -1)\n");
        sb.endOneBlock().newLine();

        sb.append("func (i *Interpreter) stepLookups(c m.KReference, config m.KReference, guard int) (m.KReference, error)").beginBlock();
        sb.appendIndentedLine("matched := false");
        sb.appendIndentedLine("v := i.stepTempVars");
        sb.appendIndentedLine("bv := i.stepTempBoolVars");

        RuleInfo lookupRuleInfo = ruleWriter.writeRule(
                lookupRules, sb, RuleType.REGULAR,
                FunctionInfo.systemFunctionInfo("stepLookups", 1));

        if (lookupRuleInfo.maxNrVars > ruleInfo.maxNrVars) {
            ruleInfo.maxNrVars = lookupRuleInfo.maxNrVars;
        }
        if (lookupRuleInfo.maxNrBoolVars > ruleInfo.maxNrBoolVars) {
            ruleInfo.maxNrBoolVars = lookupRuleInfo.maxNrBoolVars;
        }

        sb.appendIndentedLine("return c, noStep");
        sb.endOneBlock().newLine();

        sb.appendIndentedLine("// indicates the maximum number of variables required by a rule");
        sb.appendIndentedLine("// needed to initialize the matched variables (mv) slice");
        sb.appendIndentedLine("const stepMaxVarCount = " + ruleInfo.maxNrVars);
        sb.appendIndentedLine("const stepMaxBoolVarCount = " + ruleInfo.maxNrBoolVars);
        sb.newLine();

        return sb.toString();
    }

    private int sortRules(Rule r1, Rule r2) {
        return ComparisonChain.start()
                .compareTrueFirst(r1.att().contains("structural"), r2.att().contains("structural"))
                .compareFalseFirst(r1.att().contains("owise"), r2.att().contains("owise"))
                //.compareFalseFirst(indexesPoorly(r1), indexesPoorly(r2))
                .result();
    }


}
