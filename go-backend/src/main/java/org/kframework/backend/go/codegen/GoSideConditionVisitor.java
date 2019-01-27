package org.kframework.backend.go.codegen;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.builtin.Sorts;
import org.kframework.kil.Attribute;
import org.kframework.kore.KApply;
import org.kframework.kore.KToken;

public class GoSideConditionVisitor extends GoRhsVisitor {

    private enum ExpressionType { BOOLEAN, K }

    private ExpressionType expectedExprType = ExpressionType.BOOLEAN;
    private int depthFromFuncIsTrue = -1;

    public GoSideConditionVisitor(DefinitionData data,
                                   GoNameProvider nameProvider,
                                   RuleVars lhsVars,
                                   int tabsIndent, int returnValSpacesIndent) {
        super(data, nameProvider, lhsVars, tabsIndent, returnValSpacesIndent);
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

                currentSb.append("(");
                apply(k.klist().items().get(0));
                //currentSb.append(") && (");
                currentSb.append(") &&").newLine().writeIndent().append("(");
                apply(k.klist().items().get(1));
                currentSb.append(")");
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
