// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen.rules;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.model.TempVarCounters;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.builtin.Sorts;
import org.kframework.kil.Attribute;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KToken;

public class RuleSideConditionWriter extends RuleRhsWriter {

    private enum ExpressionType {BOOLEAN, K}

    private ExpressionType expectedExprType = ExpressionType.BOOLEAN;
    private int depthFromFuncIsTrue = -1;

    public RuleSideConditionWriter(DefinitionData data,
                                   GoNameProvider nameProvider,
                                   RuleVars lhsVars,
                                   TempVarCounters tempVarCounters,
                                   int tabsIndent, int returnValSpacesIndent) {
        super(data, nameProvider, lhsVars, tempVarCounters, tabsIndent, returnValSpacesIndent);
    }

    @Override
    protected void start() {
        super.start();
        if (expectedExprType == ExpressionType.BOOLEAN) {
            currentSb.append("isTrue(");
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
                    appendKTokenComment((KToken)arg1);
                    currentSb.append("/* && */").newLine().writeIndent();
                    apply(arg2);
                } else if (PrecomputePredicates.isTrueToken(arg2)) {
                    // ... && true
                    // comment everything other than the first argument
                    apply(arg1);
                    currentSb.append(" /* && */ ");
                    appendKTokenComment((KToken)arg2);
                } else {
                    // output && with both arguments
                    currentSb.append("(");
                    apply(arg1);
                    //currentSb.append(") && (");
                    currentSb.append(") &&").newLine().writeIndent().append("(");
                    apply(arg2);
                    currentSb.append(")");
                }
                return;
            case "BOOL.or":
            case "BOOL.orElse":
                assert k.klist().items().size() == 2;

                currentSb.append("(");
                apply(k.klist().items().get(0));
                currentSb.append(") || (");
                apply(k.klist().items().get(1));
                currentSb.append(")");
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
