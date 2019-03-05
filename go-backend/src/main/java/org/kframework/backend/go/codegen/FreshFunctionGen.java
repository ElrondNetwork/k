// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.kore.KLabel;
import org.kframework.kore.Sort;

import static org.kframework.Collections.*;

public class FreshFunctionGen {

    private final DefinitionData data;
    private final GoPackageManager packageManager;
    private final GoNameProvider nameProvider;

    public FreshFunctionGen(DefinitionData data, GoPackageManager packageManager, GoNameProvider nameProvider) {
        this.data = data;
        this.packageManager = packageManager;
        this.nameProvider = nameProvider;
    }

    public String generate() {
        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append("\n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")\n\n");

        sb.append("func freshFunction (s m.Sort, config m.K, counter int) (m.K, error) {\n");
        sb.append("\tswitch s {\n");
        for (Sort sort : iterable(data.mainModule.freshFunctionFor().keys())) {
            sb.append("\t\tcase m.").append(nameProvider.sortVariableName(sort));
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            KLabel freshFunction = data.mainModule.freshFunctionFor().apply(sort);
            sb.append(nameProvider.evalFunctionName(freshFunction));
            sb.append("(m.NewIntFromInt(counter), config, -1)\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Cannot find fresh function for sort \" + s.Name())\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }

}
