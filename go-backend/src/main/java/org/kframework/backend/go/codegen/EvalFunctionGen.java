package org.kframework.backend.go.codegen;

import com.google.common.collect.Sets;
import org.kframework.backend.go.GoPackageNameManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.kore.KLabel;

public class EvalFunctionGen {

    private final DefinitionData data;
    private final GoPackageNameManager packageNameManager;

    public EvalFunctionGen(DefinitionData data, GoPackageNameManager packageNameManager) {
        this.data = data;
        this.packageNameManager = packageNameManager;
    }

    public String generate() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");

        sb.append("const topCellInitializer KLabel = ");
        GoStringUtil.appendKlabelVariableName(sb.sb(), data.topCellInitializer);
        sb.append("\n\n");

        sb.append("func eval(c K, config K) (K, error)").beginBlock();
        sb.writeIndent().append("kapp, isKapply := c.(KApply)\n");
        sb.writeIndent().append("if !isKapply").beginBlock();
        sb.writeIndent().append("return c, nil").newLine();
        sb.endOneBlock();
        sb.writeIndent().append("switch kapp.Label").beginBlock();
        for (KLabel label : Sets.union(data.functions, data.anywhereKLabels)) {
            sb.writeIndent().append("case ");
            GoStringUtil.appendKlabelVariableName(sb.sb(), label);
            sb.append(":").newLine().increaseIndent();

            // arity check
            int arity = data.functionParams.get(label).arity();
            sb.writeIndent().append("if len(kapp.List) != ").append(arity).beginBlock();
            sb.writeIndent().append("return noResult, &evalArityViolatedError{funcName:\"");
            GoStringUtil.appendFunctionName(sb.sb(), label);
            sb.append("\", expectedArity: ").append(arity);
            sb.append(", actualArity: len(kapp.List)}").newLine();
            sb.endOneBlock();

            // function call
            sb.writeIndent().append("return ");
            GoStringUtil.appendFunctionName(sb.sb(), label);
            sb.append("(");
            for (int i = 0; i < arity; i++) {
                sb.append("kapp.List[").append(i).append("], ");
            }
            sb.append("config)").newLine();
            sb.decreaseIndent();
        }
        sb.writeIndent().append("default:").newLine().increaseIndent();
        sb.writeIndent().append("return c, nil").newLine();
        sb.decreaseIndent();
        sb.endOneBlock();
        sb.endOneBlock();
        sb.newLine();

        return sb.toString();
    }
}
