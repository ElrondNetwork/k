// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import com.google.common.collect.ComparisonChain;
import org.kframework.backend.go.codegen.rules.RuleWriter;
import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionParams;
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

    public StepFunctionGen(DefinitionData data, GoPackageManager packageManager, GoNameProvider nameProvider) {
        this.data = data;
        this.packageManager = packageManager;
        this.ruleWriter = new RuleWriter(data, nameProvider);
        List<Rule> unsortedRules = stream(data.mainModule.rules()).collect(Collectors.toList());
//        if (options.reverse) {
//            Collections.reverse(unsortedRules);
//        }
        sortedRules = unsortedRules.stream()
                .sorted(this::sortRules)
                .filter(r -> !data.functionRules.values().contains(r) && !r.att().contains(Attribute.MACRO_KEY) && !r.att().contains(Attribute.ALIAS_KEY) && !r.att().contains(Attribute.ANYWHERE_KEY))
                .collect(Collectors.toList());

        Map<Boolean, List<Rule>> groupedByLookup = sortedRules.stream()
                .collect(Collectors.groupingBy(RuleWriter::hasLookups));
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
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        writeStepFunction(sb, sortedRules);

        return sb.toString();
    }

    public String generateLookupsStep() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        writeLookupsStepFunction(sb, sortedRules);

        return sb.toString();
    }

    public String generateStepRules() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append(" \n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        writeStepRules(sb, sortedRules);

        return sb.toString();
    }


    private void writeStepRules(GoStringBuilder sb, List<Rule> sortedRules) {
        for (Map.Entry<Integer, Rule> entry : stepRules.entrySet()) {
            int ruleNum = entry.getKey();
            Rule r = entry.getValue();

            String funcName = "stepRule" + ruleNum;

            sb.append("func ").append(funcName).append("(c m.K, config m.K) (m.K, error)").beginBlock();

            RuleInfo ruleInfo = ruleWriter.writeRule(
                    r, sb, RuleType.REGULAR, ruleNum,
                    "step", new FunctionParams(1));
            assert !ruleInfo.alwaysMatches();

            sb.appendIndentedLine("return c, noStep");
            sb.endOneBlock().newLine();
        }

        for (Map.Entry<Integer, Rule> entry : lookupRules.entrySet()) {
            int ruleNum = entry.getKey();
            Rule r = entry.getValue();

            String funcName = "stepLookupRule" + ruleNum;

            sb.append("func ").append(funcName).append("(c m.K, config m.K, guard int) (m.K, error)").beginBlock();

            RuleInfo ruleInfo = ruleWriter.writeRule(
                    r, sb, RuleType.REGULAR, ruleNum,
                    "stepLookups", new FunctionParams(1));
            assert !ruleInfo.alwaysMatches();

            sb.appendIndentedLine("return c, noStep");
            sb.endOneBlock().newLine();
        }
    }

    private void writeStepFunction(GoStringBuilder sb, List<Rule> sortedRules) {
        sb.append("func step(c m.K) (m.K, error)").beginBlock();
        sb.writeIndent().append("config := c").newLine();
        sb.appendIndentedLine("var result m.K");
        sb.appendIndentedLine("var err error");
        for (Map.Entry<Integer, Rule> entry : stepRules.entrySet()) {
            int ruleNum = entry.getKey();
            String funcName = "stepRule" + ruleNum;
            sb.appendIndentedLine("result, err = ", funcName, "(c, config)");
            sb.writeIndent().append("if err == nil").beginBlock();
            sb.appendIndentedLine("return result, nil");
            sb.endOneBlock();
            sb.writeIndent().append("if _, isNoStep := err.(*noStepError); !isNoStep").beginBlock();
            sb.appendIndentedLine("return result, err");
            sb.endOneBlock();
        }

        sb.writeIndent().append("return stepLookups(c, config, -1)\n");
        sb.endOneBlock().newLine();
    }

    private void writeLookupsStepFunction(GoStringBuilder sb, List<Rule> sortedRules) {
        sb.append("func stepLookups(c m.K, config m.K, guard int) (m.K, error)").beginBlock();
        sb.appendIndentedLine("var result m.K");
        sb.appendIndentedLine("var err error");
        for (Map.Entry<Integer, Rule> entry : lookupRules.entrySet()) {
            int ruleNum = entry.getKey();
            String funcName = "stepLookupRule" + ruleNum;
            sb.appendIndentedLine("result, err = ", funcName, "(c, config, guard)");
            sb.writeIndent().append("if err == nil").beginBlock();
            sb.appendIndentedLine("return result, nil");
            sb.endOneBlock();
            sb.writeIndent().append("if _, isNoStep := err.(*noStepError); !isNoStep").beginBlock();
            sb.appendIndentedLine("return result, err");
            sb.endOneBlock();
        }

        sb.appendIndentedLine("return c, noStep");
        sb.endOneBlock().newLine();
    }


    private int sortRules(Rule r1, Rule r2) {
        return ComparisonChain.start()
                .compareTrueFirst(r1.att().contains("structural"), r2.att().contains("structural"))
                .compareFalseFirst(r1.att().contains("owise"), r2.att().contains("owise"))
                //.compareFalseFirst(indexesPoorly(r1), indexesPoorly(r2))
                .result();
    }


}
