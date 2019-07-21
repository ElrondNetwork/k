// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen.inline;

import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.strings.GoStringBuilder;

public class InlineMatchGen {

    private final GoPackageManager packageManager;
    private final RuleLhsMatchInlineManager inlineManager;

    public InlineMatchGen(GoPackageManager packageManager, RuleLhsMatchInlineManager inlineManager) {
        this.packageManager = packageManager;
        this.inlineManager = inlineManager;
    }

    public String generate() {
        GoStringBuilder sb = new GoStringBuilder();
        sb.append(packageManager.goGeneratedFileComment).append("\n\n");
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append("\n\n");

        sb.append("import (\n");
        sb.append("\tm \"").append(packageManager.modelPackage.getGoPath()).append("\"\n");
        sb.append(")").newLine().newLine();

        sb.appendIndentedLine("const kapplyMatchMask uint64 = ((((refTypeMask << refModelShift << refKApplyLabelShift) | refKApplyLabelMask) << refKApplyArityShift) | refKApplyArityMask) << refKApplyIndexShift");
        sb.newLine();

        for(KApplySignatureMatch signature: inlineManager.kapplySignatures) {
            sb.append("const ").append(signature.matchConstName);
            sb.append(" = ((((uint64(kapplyRef) << refModelShift << refKApplyLabelShift) | uint64(m.");
            sb.append(signature.labelName);
            sb.append(")) << refKApplyArityShift) | ");
            sb.append(signature.arity);
            sb.append(") << refKApplyIndexShift");
            sb.newLine();
        }

        sb.newLine();
        sb.appendIndentedLine("const ktokenMatchMask uint64 = ((refTypeMask << refModelShift << refKTokenSortShift) | refKTokenSortMask) << refKTokenLengthShift << refKTokenIndexShift");
        sb.newLine();

        for (String ktokenSort : inlineManager.ktokenSortNames) {
            sb.append("const ").append(inlineManager.ktokenMatchConstName((ktokenSort)));
            sb.append(" = ((uint64(ktokenRef) << refModelShift << refKTokenSortShift) | uint64(m.");
            sb.append(ktokenSort);
            sb.append(")) << refKTokenLengthShift << refKTokenIndexShift");
            sb.newLine();
        }

        sb.newLine();
        sb.appendIndentedLine("const collectionMatchMask uint64 = ((refTypeMask << refModelShift << refCollectionSortShift) | refCollectionSortMask) << refCollectionLabelShift << refCollectionIndexShift");
        sb.newLine();

        for (String sortName : inlineManager.mapSortNames) {
            sb.append("const ").append(inlineManager.mapMatchConstName((sortName)));
            sb.append(" = ((uint64(mapRef) << refModelShift << refCollectionSortShift) | uint64(m.");
            sb.append(sortName);
            sb.append(")) << refCollectionLabelShift << refCollectionIndexShift");
            sb.newLine();
        }
        for (String sortName : inlineManager.setSortNames) {
            sb.append("const ").append(inlineManager.setMatchConstName((sortName)));
            sb.append(" = ((uint64(setRef) << refModelShift << refCollectionSortShift) | uint64(m.");
            sb.append(sortName);
            sb.append(")) << refCollectionLabelShift << refCollectionIndexShift");
            sb.newLine();
        }
        for (String sortName : inlineManager.listSortNames) {
            sb.append("const ").append(inlineManager.listMatchConstName((sortName)));
            sb.append(" = ((uint64(listRef) << refModelShift << refCollectionSortShift) | uint64(m.");
            sb.append(sortName);
            sb.append(")) << refCollectionLabelShift << refCollectionIndexShift");
            sb.newLine();
        }
        for (String sortName : inlineManager.arraySortNames) {
            sb.append("const ").append(inlineManager.arrayMatchConstName((sortName)));
            sb.append(" = ((uint64(arrayRef) << refModelShift << refCollectionSortShift) | uint64(m.");
            sb.append(sortName);
            sb.append(")) << refCollectionLabelShift << refCollectionIndexShift");
            sb.newLine();
        }

        return sb.toString();
    }

}
