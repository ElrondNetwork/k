package org.kframework.backend.go.codegen;

import org.kframework.backend.go.GoPackageNameManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.kore.KLabel;
import org.kframework.kore.Sort;

import static org.kframework.Collections.*;

public class FreshFunctionGen {

    private final DefinitionData data;
    private final GoPackageNameManager packageNameManager;
    private final GoNameProvider nameProvider;

    public FreshFunctionGen(DefinitionData data, GoPackageNameManager packageNameManager, GoNameProvider nameProvider) {
        this.data = data;
        this.packageNameManager = packageNameManager;
        this.nameProvider = nameProvider;
    }

    public String generate() {
        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");

        sb.append("func freshFunction (s Sort, config K, counter int) (K, error) {\n");
        sb.append("\tswitch s {\n");
        for (Sort sort : iterable(data.mainModule.freshFunctionFor().keys())) {
            sb.append("\t\tcase ").append(nameProvider.sortVariableName(sort));
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            KLabel freshFunction = data.mainModule.freshFunctionFor().apply(sort);
            sb.append(nameProvider.evalFunctionName(freshFunction));
            sb.append("(Int(counter), config)\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Cannot find fresh function for sort \" + s.name())\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }

}
