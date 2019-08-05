package org.kframework.backend.go.codegen.rules;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionInfo;
import org.kframework.backend.go.model.VarContainer;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.unparser.ToKast;

public class RuleRhsWriter extends RuleRhsWriterBase {

    private boolean isTopLevelReturn;
    private FunctionInfo parentFunction; // we need it to handle tail recursion

    public RuleRhsWriter(
            DefinitionData data,
            GoNameProvider nameProvider,
            VarContainer vars,
            int tabsIndent,
            boolean isTopLevelReturn,
            FunctionInfo parentFunction) {
        super(data, nameProvider, vars, tabsIndent, 0);
        this.isTopLevelReturn = isTopLevelReturn;
        this.parentFunction = parentFunction;
    }

    @Override
    protected RuleRhsWriterBase newInstanceWithSameConfig(int indent) {
        return new RuleRhsWriter(data, nameProvider,
                vars,
                indent,
                false, parentFunction);
    }

    @Override
    protected void start() {
        isTopLevelReturn = false;
        super.start();
    }

    @Override
    public void apply(K k) {
        if (isTopLevelReturn) {
            if (parentFunction != null && !parentFunction.isSystemFunction()) {
                if (k instanceof KApply) {
                    KApply kapp = (KApply)k;
                    if (kapp.klabel() == parentFunction.label &&
                            kapp.klist().items().size() == parentFunction.arguments.arity()) {
                        // tail recursion detected!
                        writeTailRecursionEval(kapp);
                        return;
                    }
                }
            }

            currentSb.append("return ");
            super.apply(k);
            currentSb.append(", nil");
        } else {
            super.apply(k);
        }
    }

    private void writeTailRecursionEval(KApply kapp) {
        isTopLevelReturn = false;
        int arity = kapp.klist().items().size();
        currentSb.append("// tail recursion ").append(ToKast.apply(kapp));
        for (int i = 0; i < arity; i++) {
            currentSb.newLine().writeIndent();
            currentSb.append(parentFunction.arguments.varName(i)).append(" = ");
            apply(kapp.klist().items().get(i));
        }
        currentSb.newLine().writeIndent().append("matched = true");

    }

}
