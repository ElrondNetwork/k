// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import org.kframework.backend.go.gopackage.GoPackage;
import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoStringBuilder;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.Set;
import java.util.TreeSet;
import java.util.stream.Collectors;

public class InterpreterDefGen {

    private final DefinitionData data;
    private final GoPackageManager packageManager;

    public InterpreterDefGen(DefinitionData data, GoPackageManager packageManager) {
        this.data = data;
        this.packageManager = packageManager;
    }

    public String generate() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append(packageManager.goGeneratedFileComment).append("\n\n");
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append("\n\n");

        // generating imports
        Set<GoPackage> importsSorted = new TreeSet<>(Comparator.comparing(GoPackage::getName));
        importsSorted.add(packageManager.modelPackage);
        importsSorted.addAll(data.extHookManager.nameToPackage.values());
        sb.append("import (").newLine();
        for (GoPackage pkg : importsSorted) {
            if (pkg.getGoPath() == null) {
                sb.append("\t\"").append(pkg.getName()).append("\"").newLine(); // "fmt"
            } else {
                sb.append("\t").append(pkg.getAlias()).append(" \"").append(pkg.getGoPath()).append("\"").newLine();
            }
        }
        sb.append(")\n\n");

        // generate fields
        List<FieldDefinition> fields = new ArrayList<>();
        for (String extHookName : data.extHookManager.nameToPackage.keySet()) {
            String lower = extHookName.toLowerCase();
            String capitalized = Character.toUpperCase(lower.charAt(0)) + lower.substring(1);

            String fieldName = lower + "Ref";
            String fieldDef = fieldName + " *" + extHookName + "." + capitalized;
            fields.add(new FieldDefinition(fieldName, fieldDef));
        }

        sb.append("// Interpreter is a container with a reference to model and basic options").newLine();
        sb.append("type Interpreter struct").beginBlock();

        sb.appendIndentedLine("Model         *m.ModelState");
        sb.appendIndentedLine("traceHandlers []traceHandler");
        sb.appendIndentedLine("Verbose       bool");
        sb.appendIndentedLine("MaxSteps      int");

        // hook references here
        if (fields.size() > 0) {
            sb.newLine();
            for (FieldDefinition field : fields) {
                sb.appendIndentedLine(field.fieldDef);
            }
        }

        sb.endOneBlock();

        sb.append("// NewInterpreter creates a new interpreter instance").newLine();
        sb.append("func NewInterpreter(");
        sb.append(fields.stream().map(fd -> fd.fieldDef).collect(Collectors.joining(", ")));
        sb.append(") *Interpreter").beginBlock();

        sb.appendIndentedLine("model := &m.ModelState{}");
        sb.appendIndentedLine("model.Init()");
        sb.newLine();
        sb.writeIndent().append("return &Interpreter").beginBlock();
        sb.appendIndentedLine("Model:         model,");
        sb.appendIndentedLine("traceHandlers: nil,");
        sb.appendIndentedLine("Verbose:       false,");
        sb.appendIndentedLine("MaxSteps:      0,");
        for (FieldDefinition field : fields) {
            sb.appendIndentedLine(field.fieldName, ": ", field.fieldName, ",");
        }
        sb.endOneBlock();
        sb.endOneBlock();

        return sb.toString();
    }

    private class FieldDefinition {
        public final String fieldName;
        public final String fieldDef;

        public FieldDefinition(String fieldName, String fieldDef) {
            this.fieldName = fieldName;
            this.fieldDef = fieldDef;
        }
    }

}
