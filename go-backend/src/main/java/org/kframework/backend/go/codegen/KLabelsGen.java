// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen;

import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.compile.ConvertDataStructureToLookup;
import org.kframework.kil.Attribute;
import org.kframework.kore.KLabel;
import org.kframework.kore.KORE;
import org.kframework.kore.Sort;

import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

import static org.kframework.Collections.*;
import static org.kframework.kore.KORE.*;

public class KLabelsGen {

    private final DefinitionData data;
    private final GoPackageManager packageManager;
    private final GoNameProvider nameProvider;
    private final Map<KLabel, KLabel> collectionFor;

    public KLabelsGen(DefinitionData data, GoPackageManager packageManager, GoNameProvider nameProvider) {
        this.data = data;
        this.packageManager = packageManager;
        this.nameProvider = nameProvider;
        collectionFor = ConvertDataStructureToLookup.collectionFor(data.mainModule);
    }

    public String klabels() {
        Set<KLabel> klabels = mutable(data.mainModule.definedKLabels());
        klabels.add(KORE.KLabel("#Bottom"));
        klabels.add(KORE.KLabel("littleEndianBytes"));
        klabels.add(KORE.KLabel("bigEndianBytes"));
        klabels.add(KORE.KLabel("signedBytes"));
        klabels.add(KORE.KLabel("unsignedBytes"));
        //addOpaqueKLabels(klabels);

        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageManager.modelPackage.getName()).append("\n\n");
        sb.append("// KLabel ... a k label identifier").newLine();
        sb.append("type KLabel int\n\n");

        int i = 0;
        for (KLabel klabel : klabels) {
            sb.appendIndentedLine("// ", nameProvider.klabelVariableName(klabel), " ... ", klabel.name());
            sb.appendIndentedLine("const ", nameProvider.klabelVariableName(klabel), " KLabel = ", Integer.toString(i++));
            sb.newLine();
        }
        sb.appendIndentedLine("// LblDummy ... dummy label used in tests");
        sb.appendIndentedLine("const LblDummy KLabel = ", Integer.toString(i++));
        sb.newLine().newLine();

        sb.append("//KLabelForMap ... The KLabel that identifies maps").newLine();
        sb.append("const KLabelForMap KLabel = ").append(nameProvider.klabelVariableName(KORE.KLabel("_Map_"))).newLine();
        sb.append("//KLabelForSet ... The KLabel that identifies sets").newLine();
        sb.append("const KLabelForSet KLabel = ").append(nameProvider.klabelVariableName(KORE.KLabel("_Set_"))).newLine();
        sb.append("//KLabelForList ... The KLabel that identifies lists").newLine();
        sb.append("const KLabelForList KLabel = ").append(nameProvider.klabelVariableName(KORE.KLabel("_List_"))).newLine();
        sb.newLine();

        // klabel name method
        sb.append("// Name ... KLabel name").newLine();
        sb.append("func (kl KLabel) Name () string").beginBlock();
        sb.writeIndent().append("switch kl").beginBlock();
        for (KLabel klabel : klabels) {
            sb.writeIndent().append("case ").append(nameProvider.klabelVariableName(klabel));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            sb.append(GoStringUtil.enquoteString(klabel.name()));
            sb.newLine();
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Unexpected KLabel.\")\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // parse klabel function
        sb.append("// ParseKLabel ... Yields the KLabel with the given name").newLine();
        sb.append("func ParseKLabel (name string) KLabel").beginBlock();
        sb.append("\tswitch name").beginBlock();
        for (KLabel klabel : klabels) {
            sb.writeIndent().append("case ");
            sb.append(GoStringUtil.enquoteString(klabel.name()));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ").append(nameProvider.klabelVariableName(klabel));
            sb.append("\n");
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Parsing KLabel failed. Unexpected KLabel name:\" + name)\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // collection for
        sb.append("// CollectionFor ... TODO: document").newLine();
        sb.append("func CollectionFor(kl KLabel) KLabel").beginBlock();
        sb.append("\tswitch kl").beginBlock();
        for (Map.Entry<KLabel, KLabel> entry : collectionFor.entrySet()) {
            sb.writeIndent().append("case ");
            sb.append(nameProvider.klabelVariableName(entry.getKey()));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            sb.append(nameProvider.klabelVariableName(entry.getValue()));
            sb.newLine();
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Cannot call method collectionFor for KLabel \" + kl.Name())\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // unit for
        sb.append("// UnitFor ... TODO: document").newLine();
        sb.append("func UnitFor(kl KLabel) KLabel").beginBlock();
        sb.append("\tswitch kl").beginBlock();
        for (KLabel label : collectionFor.values().stream().collect(Collectors.toSet())) {
            sb.writeIndent().append("case ");
            sb.append(nameProvider.klabelVariableName(label));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            KLabel unitLabel = KORE.KLabel(data.mainModule.attributesFor().apply(label).get(Attribute.UNIT_KEY));
            sb.append(nameProvider.klabelVariableName(unitLabel));
            sb.newLine();
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Cannot call method UnitFor for KLabel \" + kl.Name())\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // element for
        sb.append("// ElementFor ... TODO: document").newLine();
        sb.append("func ElementFor(kl KLabel) KLabel").beginBlock();
        sb.append("\tswitch kl").beginBlock();
        for (KLabel label : collectionFor.values().stream().collect(Collectors.toSet())) {
            sb.writeIndent().append("case ").append(nameProvider.klabelVariableName(label));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            KLabel elLabel = KORE.KLabel(data.mainModule.attributesFor().apply(label).get("element"));
            sb.append(nameProvider.klabelVariableName(elLabel));
            sb.newLine();
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Cannot call method ElementFor for KLabel \" + kl.Name())\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // unit for array
        Set<Sort> sorts = mutable(data.mainModule.definedSorts());
        Set<Sort> arraySorts = sorts.stream().filter(s -> data.mainModule.sortAttributesFor().get(s).getOrElse(() -> Att()).<String>getOptional("hook")
                .orElse(".").equals("ARRAY.Array")).collect(Collectors.toSet());
        sb.append("// UnitForArray ... TODO: document").newLine();
        sb.append("func UnitForArray(s Sort) KLabel").beginBlock();
        sb.append("\tswitch s").beginBlock();
        for (Sort sort : arraySorts) {
            sb.writeIndent().append("case ").append(nameProvider.sortVariableName(sort));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            KLabel elLabel = KLabel(data.mainModule.sortAttributesFor().apply(sort).get("unit"));
            sb.append(nameProvider.klabelVariableName(elLabel));
            sb.newLine();
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Unexpected Sort in method UnitForArray.\")\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // element for array
        sb.append("// ElementForArray ... TODO: document").newLine();
        sb.append("func ElementForArray(s Sort) KLabel").beginBlock();
        sb.append("\tswitch s").beginBlock();
        for (Sort sort : arraySorts) {
            sb.writeIndent().append("case ").append(nameProvider.sortVariableName(sort));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            KLabel elLabel = KLabel(data.mainModule.sortAttributesFor().apply(sort).get("element"));
            sb.append(nameProvider.klabelVariableName(elLabel));
            sb.newLine();
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Unexpected Sort in method ElementForArray.\")\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        return sb.toString();
    }

}
