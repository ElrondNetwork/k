package org.kframework.backend.go.codegen;

import com.google.common.collect.Sets;
import org.kframework.backend.go.GoPackageNameManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.kore.KLabel;

public class EvalFunctionGen {

    private final DefinitionData data;
    private final GoPackageNameManager packageNameManager;

    public EvalFunctionGen(DefinitionData data, GoPackageNameManager packageNameManager) {
        this.data = data;
        this.packageNameManager = packageNameManager;
    }

    /**
     * WARNING: depends on fields functions and anywhereKLabels, only run after definition()
     * TODO: untangle this dependency
     */
    public String generate() {
        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");

        sb.append("import \"fmt\"\n\n");

        sb.append("const topCellInitializer KLabel = ");
        GoStringUtil.appendKlabelVariableName(sb, data.topCellInitializer);
        sb.append("\n\n");

        sb.append("func eval(c K, config K) K {\n");
        sb.append("\tkApply, typeOk := c.(KApply)\n");
        sb.append("\tif !typeOk {\n");
        sb.append("\t\treturn c\n");
        sb.append("\t}\n");
        sb.append("\tswitch kApply.Label {\n");
        for (KLabel label : Sets.union(data.functions, data.anywhereKLabels)) {
            sb.append("\t\tcase ");
            GoStringUtil.appendKlabelVariableName(sb, label);
            sb.append(":\n");

            // arity check
            int arity = data.functionParams.get(label).arity();
            sb.append("\t\t\tif len(kApply.List) != ").append(arity).append(" {\n");
            sb.append("\t\t\t\tpanic(fmt.Sprintf(\"");
            GoStringUtil.appendFunctionName(sb, label);
            sb.append(" function arity violated. Expected arity: ").append(arity);
            sb.append(". Nr. params provided: %d\", len(kApply.List)))\n");
            sb.append("\t\t\t}\n");

            // function call
            sb.append("\t\t\treturn ");
            GoStringUtil.appendFunctionName(sb, label);
            sb.append("(");
            for (int i = 0; i < arity; i++) {
                sb.append("kApply.List[").append(i).append("], ");
            }
            sb.append("config)\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\treturn c\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }
}
