package org.kframework.backend.go.codegen;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.builtin.Sorts;
import org.kframework.kil.Attribute;
import org.kframework.kore.KApply;
import org.kframework.kore.KToken;

public class GoSideConditionVisitor extends GoRhsVisitor {

    private enum ExpressionType { BOOLEAN, K }

    private ExpressionType expectedExprType = ExpressionType.BOOLEAN;
    private int depthFromFuncIsTrue = -1;

    public GoSideConditionVisitor(GoStringBuilder sb, DefinitionData data, RuleVars lhsVars) {
        super(sb, data, lhsVars);
    }

    @Override
    protected void start() {
        super.start();
        if (expectedExprType == ExpressionType.BOOLEAN) {
            sb.append("isTrue(");
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
                //sb.newLine().writeIndent().append(")");
                sb.append(")");
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
        if (expectedExprType == ExpressionType.BOOLEAN) {
            String hook = data.mainModule.attributesFor().apply(k.klabel()).<String>getOptional(Attribute.HOOK_KEY).orElse("");
            switch (hook) {
            case "BOOL.and":
            case "BOOL.andThen":
                assert k.klist().items().size() == 2;

                sb.append("(");
                apply(k.klist().items().get(0));
                //sb.append(") && (");
                sb.append(") &&").newLine().writeIndent().append("(");
                apply(k.klist().items().get(1));
                sb.append(")");
                return;
            case "BOOL.or":
            case "BOOL.orElse":
                assert k.klist().items().size() == 2;

                sb.append("(");
                apply(k.klist().items().get(0));
                sb.append(") || (");
                apply(k.klist().items().get(1));
                sb.append(")");
                return;
            case "BOOL.not":
                assert k.klist().items().size() == 1;

                sb.append("!(");
                apply(k.klist().items().get(0));
                sb.append(")");
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
            sb.append(k.s()); // true or false, directly
            return;
        }

        super.apply(k);
    }
}
