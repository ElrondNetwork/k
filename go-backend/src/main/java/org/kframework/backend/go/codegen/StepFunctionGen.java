package org.kframework.backend.go.codegen;

import com.google.common.collect.ComparisonChain;
import org.kframework.backend.go.GoPackageNameManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.definition.Rule;
import org.kframework.kil.Attribute;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import static org.kframework.Collections.*;

public class StepFunctionGen {

    private final DefinitionData data;
    private final GoPackageNameManager packageNameManager;
    private final RuleWriter ruleWriter;

    private final List<Rule> sortedRules;

    private int ruleNum = 0;

    public StepFunctionGen(DefinitionData data, GoPackageNameManager packageNameManager) {
        this.data = data;
        this.packageNameManager = packageNameManager;
        this.ruleWriter = new RuleWriter(data);
        List<Rule> unsortedRules = stream(data.mainModule.rules()).collect(Collectors.toList());
//        if (options.reverse) {
//            Collections.reverse(unsortedRules);
//        }
        sortedRules = unsortedRules.stream()
                .sorted(this::sortRules)
                .filter(r -> !data.functionRules.values().contains(r) && !r.att().contains(Attribute.MACRO_KEY) && !r.att().contains(Attribute.ALIAS_KEY) && !r.att().contains(Attribute.ANYWHERE_KEY))
                .collect(Collectors.toList());
    }

    public String generateStep() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append(" \n\n");

        writeStepFunction(sb, sortedRules, "step");

        return sb.toString();
    }

    public String generateLookupsStep() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append(" \n\n");

        writeLookupsStepFunction(sb, sortedRules, "step");

        return sb.toString();
    }

    private void writeStepFunction(GoStringBuilder sb, List<Rule> sortedRules, String funcName) {
        String lookupsFuncName = funcName + "Lookups";
        Map<Boolean, List<Rule>> groupedByLookup = sortedRules.stream()
                .collect(Collectors.groupingBy(RuleWriter::hasLookups));

        sb.append("func ").append(funcName).append("(c K) (K, error)").beginBlock();
        sb.writeIndent().append("config := c").newLine();
        if (groupedByLookup.containsKey(false)) {
            for (Rule r : groupedByLookup.get(false)) {
                ruleNum = ruleWriter.convert(r, sb, RuleType.REGULAR, ruleNum, new FunctionParams(1));
            }
        }

        sb.writeIndent().append("return ").append(lookupsFuncName).append("(c, config, -1)\n");
        sb.endOneBlock().newLine();
    }

    private void writeLookupsStepFunction(GoStringBuilder sb, List<Rule> sortedRules, String funcName) {
        String lookupsFuncName = funcName + "Lookups";
        Map<Boolean, List<Rule>> groupedByLookup = sortedRules.stream()
                .collect(Collectors.groupingBy(RuleWriter::hasLookups));

        sb.append("func ").append(lookupsFuncName).append("(c K, config K, guard int) (K, error)").beginBlock("k * step_function");

        List<Rule> lookupRules = groupedByLookup.getOrDefault(true, Collections.emptyList());
        for (Rule r : lookupRules) {
            ruleNum = ruleWriter.convert(r, sb, RuleType.REGULAR,ruleNum, new FunctionParams(1));
        }

        sb.appendIndentedLine("return c, &noStepError{}");
        sb.endAllBlocks(0);
        sb.append("\n");
    }


    private int sortRules(Rule r1, Rule r2) {
        return ComparisonChain.start()
                .compareTrueFirst(r1.att().contains("structural"), r2.att().contains("structural"))
                .compareFalseFirst(r1.att().contains("owise"), r2.att().contains("owise"))
                //.compareFalseFirst(indexesPoorly(r1), indexesPoorly(r2))
                .result();
    }


}
