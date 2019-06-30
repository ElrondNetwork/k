// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen.rules;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.RuleVarContainer;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.model.TempVarCounters;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.builtin.Sorts;
import org.kframework.kil.Attribute;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KToken;
import org.kframework.unparser.ToKast;

public class RuleSideConditionWriter extends RuleRhsWriterBase {

    private enum ExpressionType {BOOLEAN, K}

    private ExpressionType expectedExprType = ExpressionType.BOOLEAN;
    private int depthFromFuncIsTrue = -1;

    public RuleSideConditionWriter(DefinitionData data,
                                   GoNameProvider nameProvider,
                                   RuleVarContainer vars,
                                   int tabsIndent) {
        super(data, nameProvider, vars, false, tabsIndent, "if ".length());
    }

    @Override
    protected RuleRhsWriterBase newInstanceWithSameConfig(int indent) {
        return new RuleSideConditionWriter(data, nameProvider,
                vars,
                indent);
    }

    @Override
    protected void start() {
        super.start();
        if (expectedExprType == ExpressionType.BOOLEAN) {
            currentSb.append("m.IsTrue(");
            expectedExprType = ExpressionType.K; // entering K territory
            depthFromFuncIsTrue = 0;
        } else {
            depthFromFuncIsTrue++;
        }
    }

    @Override
    protected void end() {
        assert depthFromFuncIsTrue != -1;

        if (expectedExprType == ExpressionType.K) {
            if (depthFromFuncIsTrue == 0) {
                //currentSb.newLine().writeIndent().append(")");
                currentSb.append(")");
                expectedExprType = ExpressionType.BOOLEAN; // back from K territory
            } else {
                depthFromFuncIsTrue--;
            }
        }
    }

    @Override
    public void apply(KApply k) {
        // we don't call start() just yet, because it would automatically generate a call to isTrue(...)
        // we try to replace some calls with native &&, ||, !
        if (expectedExprType == ExpressionType.BOOLEAN && data.isFunctionOrAnywhere(k.klabel())) {
            String hook = data.mainModule.attributesFor().apply(k.klabel()).<String>getOptional(Attribute.HOOK_KEY).orElse("");
            switch (hook) {
            case "BOOL.and":
            case "BOOL.andThen":
                assert k.klist().items().size() == 2;
                K arg1 = k.klist().items().get(0);
                K arg2 = k.klist().items().get(1);

                if (PrecomputePredicates.isTrueToken(arg1) &&
                        PrecomputePredicates.isTrueToken(arg2)) {
                    throw new RuntimeException("true && true should already have been collapsed in PrecomputePredicates.");
                } else if (PrecomputePredicates.isTrueToken(arg1)) {
                    // true && ...
                    // comment everything other than the second argument
                    appendKTokenComment((KToken) arg1);
                    currentSb.append("/* && */ ");
                    apply(arg2);
                } else if (PrecomputePredicates.isTrueToken(arg2)) {
                    // ... && true
                    // comment everything other than the first argument
                    apply(arg1);
                    currentSb.append(" /* && */ ");
                    appendKTokenComment((KToken) arg2);
                } else {
                    // we trick all nodes below to output to the eval call instead of the return by swapping the string builder
                    GoStringBuilder evalSb = new GoStringBuilder(topLevelIndent, 0);
                    GoStringBuilder backupSb = currentSb;
                    currentSb = evalSb;

                    // get arg1 evaluation first
                    String andVarName = "evalAnd" + vars.varCounters.consumeEvalVarIndex();
                    evalSb.appendIndentedLine("var ", andVarName, " bool // ", ToKast.apply(k));
                    evalSb.writeIndent().append(andVarName).append(" = ");
                    apply(arg1);
                    evalSb.newLine();
                    evalSb.writeIndent().append("if ").append(andVarName).beginBlock();

                    // all evaluations for arg2 need to happen in this block,
                    // to avoid executing any of them if arg1 is false
                    RuleRhsWriterBase arg2Writer = newInstanceWithSameConfig(evalSb.getCurrentIndent());
                    arg2Writer.apply(arg2);

                    arg2Writer.writeEvalCalls(evalSb);
                    evalSb.writeIndent().append(andVarName).append(" = ");
                    arg2Writer.writeReturnValue(evalSb);
                    evalSb.newLine();
                    evalSb.endOneBlock();

                    assert currentSb == evalSb;
                    currentSb = backupSb; // restore

                    evalCalls.add(evalSb.toString()); // eval call
                    currentSb.append(andVarName); // this is the actual result, we output the name of the variable
                }
                return;
            case "BOOL.or":
            case "BOOL.orElse":
                assert k.klist().items().size() == 2;

                // we trick all nodes below to output to the eval call instead of the return by swapping the string builder
                GoStringBuilder evalSb = new GoStringBuilder(topLevelIndent, 0);
                GoStringBuilder backupSb = currentSb;
                currentSb = evalSb;

                String orVarName = "evalOr" + vars.varCounters.consumeEvalVarIndex();
                evalSb.appendIndentedLine("var ", orVarName, " bool // ", ToKast.apply(k));

                // get arg1 evaluation first
                evalSb.writeIndent().append(orVarName).append(" = ");
                apply(k.klist().items().get(0));
                evalSb.newLine();
                evalSb.writeIndent().append("if !").append(orVarName).beginBlock();

                // all evaluations for arg2 have happen in this block,
                // to avoid executing any of them if arg1 is true
                RuleRhsWriterBase arg2Writer = newInstanceWithSameConfig(evalSb.getCurrentIndent());
                arg2Writer.apply(k.klist().items().get(1));

                arg2Writer.writeEvalCalls(evalSb);
                evalSb.writeIndent().append(orVarName).append(" = ");
                arg2Writer.writeReturnValue(evalSb);
                evalSb.newLine();
                evalSb.endOneBlock();

                // restore
                assert currentSb == evalSb;
                currentSb = backupSb;

                evalCalls.add(evalSb.toString()); // eval call
                currentSb.append(orVarName); // this is the actual result, we output the name of the variable
                return;
            case "BOOL.not":
                assert k.klist().items().size() == 1;

                currentSb.append("!(");
                apply(k.klist().items().get(0));
                currentSb.append(")");
                return;
            default:
                break;
            }
        }

        // no optimization could be performed, carry on
        // start() will be called, isTrue(...) appears, all calls down the stack are regular K output
        super.apply(k);
    }

    @Override
    public void apply(KToken k) {
        if (expectedExprType == ExpressionType.BOOLEAN && k.sort().equals(Sorts.Bool())) {
            appendKTokenComment(k);
            currentSb.append(k.s()); // true or false, directly
            return;
        }

        super.apply(k);
    }

}
