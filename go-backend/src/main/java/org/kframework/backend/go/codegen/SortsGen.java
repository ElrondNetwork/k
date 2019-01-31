package org.kframework.backend.go.codegen;

import org.kframework.backend.go.gopackage.GoPackageNameManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.builtin.Sorts;
import org.kframework.kore.Sort;

import java.util.Set;

import static org.kframework.Collections.*;

public class SortsGen {

    private final DefinitionData data;
    private final GoPackageNameManager packageNameManager;
    private final GoNameProvider nameProvider;

    public SortsGen(DefinitionData data, GoPackageNameManager packageNameManager, GoNameProvider nameProvider) {
        this.data = data;
        this.packageNameManager = packageNameManager;
        this.nameProvider = nameProvider;
    }

    public String generate() {
        Set<Sort> sorts = mutable(data.mainModule.definedSorts());
        sorts.add(Sorts.Bool());
        sorts.add(Sorts.MInt());
        sorts.add(Sorts.Int());
        sorts.add(Sorts.String());
        sorts.add(Sorts.Float());
        sorts.add(Sorts.StringBuffer());
        sorts.add(Sorts.Bytes());

        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageNameManager.modelPackage.getName()).append(" \n\n");
        sb.append("// Sort ... a K sort identifier\n");
        sb.append("type Sort int\n\n");

        // const declaration
        sb.append("const (\n");
        for (Sort s : sorts) {
            sb.append("\t// ").append(nameProvider.sortVariableName(s)).append(" ... ").append(s.name()).append("\n");
            sb.append("\t").append(nameProvider.sortVariableName(s));
            sb.append(" Sort = iota\n");
        }
        sb.append(")\n\n");

        // sort name method
        sb.append("// Name ... Sort name\n");
        sb.append("func (s Sort) Name () string {\n");
        sb.append("\tswitch s {\n");
        for (Sort sort : sorts) {
            sb.append("\t\tcase ").append(nameProvider.sortVariableName(sort));
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            sb.append(GoStringUtil.enquoteString(sort.name()));
            sb.append("\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Unexpected Sort.\")\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        // parse sort function
        sb.append("// ParseSort ... Yields the sort with the given name\n");
        sb.append("func ParseSort (name string) Sort {\n");
        sb.append("\tswitch name {\n");
        for (Sort sort : sorts) {
            sb.append("\t\tcase ");
            sb.append(GoStringUtil.enquoteString(sort.name()));
            sb.append(":\n");
            sb.append("\t\t\treturn ").append(nameProvider.sortVariableName(sort));
            sb.append("\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Parsing Sort failed. Unexpected Sort name:\" + name)\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }

}
