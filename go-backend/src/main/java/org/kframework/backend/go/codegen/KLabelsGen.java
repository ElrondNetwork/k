package org.kframework.backend.go.codegen;

import org.kframework.backend.go.GoPackageNameManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.compile.ConvertDataStructureToLookup;
import org.kframework.kil.Attribute;
import org.kframework.kore.KLabel;
import org.kframework.kore.KORE;

import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

import static org.kframework.Collections.*;

public class KLabelsGen {

    private final DefinitionData data;
    private final GoPackageNameManager packageNameManager;
    private final GoNameProvider nameProvider;
    private final Map<KLabel, KLabel> collectionFor;

    public KLabelsGen(DefinitionData data, GoPackageNameManager packageNameManager, GoNameProvider nameProvider) {
        this.data = data;
        this.packageNameManager = packageNameManager;
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
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");
        sb.append("// KLabel ... a k label identifier").newLine();
        sb.append("type KLabel int\n\n");

        // const declaration
        sb.append("const (\n");
        sb.increaseIndent();
        for (KLabel klabel : klabels) {
            sb.writeIndent().append(nameProvider.klabelVariableName(klabel));
            sb.append(" KLabel = iota\n");
        }
        sb.decreaseIndent();
        sb.append(")").newLine().newLine();

        sb.append("const klabelForMap KLabel = ").append(nameProvider.klabelVariableName(KORE.KLabel("_Map_"))).newLine();
        sb.append("const klabelForSet KLabel = ").append(nameProvider.klabelVariableName(KORE.KLabel("_Set_"))).newLine();
        sb.append("const klabelForList KLabel = ").append(nameProvider.klabelVariableName(KORE.KLabel("_List_"))).newLine();
        sb.newLine();

        // klabel name method
        sb.append("func (kl KLabel) name () string").beginBlock();
        sb.writeIndent().append("switch kl").beginBlock();
        for (KLabel klabel : klabels) {
            sb.writeIndent().append("case ").append(nameProvider.klabelVariableName(klabel));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            sb.append(GoStringUtil.enquoteString(klabel.name()));
            sb.append("\n");
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Unexpected KLabel.\")\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // parse klabel function
        sb.append("func parseKLabel (name string) KLabel").beginBlock();
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
        sb.append("func (kl KLabel) collectionFor() KLabel").beginBlock();
        sb.append("\tswitch kl").beginBlock();
        for (Map.Entry<KLabel, KLabel> entry : collectionFor.entrySet()) {
            sb.writeIndent().append("case ");
            sb.append(nameProvider.klabelVariableName(entry.getKey()));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            sb.append(nameProvider.klabelVariableName(entry.getValue()));
            sb.append("\n");
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Cannot call method collectionFor for KLabel \" + kl.name())\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // unit for
        sb.append("func (kl KLabel) unitFor() KLabel").beginBlock();
        sb.append("\tswitch kl").beginBlock();
        for (KLabel label : collectionFor.values().stream().collect(Collectors.toSet())) {
            sb.writeIndent().append("case ");
            sb.append(nameProvider.klabelVariableName(label));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            KLabel unitLabel = KORE.KLabel(data.mainModule.attributesFor().apply(label).get(Attribute.UNIT_KEY));
            sb.append(nameProvider.klabelVariableName(unitLabel));
            sb.append("\n");
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Cannot call method unitFor for KLabel \" + kl.name())\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        // el for
        sb.append("func (kl KLabel) elFor() KLabel").beginBlock();
        sb.append("\tswitch kl").beginBlock();
        for (KLabel label : collectionFor.values().stream().collect(Collectors.toSet())) {
            sb.writeIndent().append("case ").append(nameProvider.klabelVariableName(label));
            sb.append(":\n");
            sb.writeIndent().append("\treturn ");
            KLabel elLabel = KORE.KLabel(data.mainModule.attributesFor().apply(label).get("element"));
            sb.append(nameProvider.klabelVariableName(elLabel));
            sb.append("\n");
        }
        sb.writeIndent().append("default:\n");
        sb.writeIndent().append("\tpanic(\"Cannot call method elFor for KLabel \" + kl.name())\n");
        sb.endOneBlock();
        sb.endOneBlock().newLine();

        return sb.toString();
    }

}
