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

        sb.appendIndentedLine("const kapplyMatchMask uint64 = ((((refTypeMask << 1 << refKApplyLabelShift) | refKApplyLabelMask) << refKApplyArityShift) | refKApplyArityMask) << refKApplyIndexShift");
        sb.newLine();

        for(KApplySignatureMatch signature: inlineManager.kapplySignatures) {
            sb.append("const ").append(signature.matchConstName);
            sb.append(" = ((((uint64(kapplyRef) << 1 << refKApplyLabelShift) | uint64(m.");
            sb.append(signature.labelName);
            sb.append(")) << refKApplyArityShift) | ");
            sb.append(signature.arity);
            sb.append(") << refKApplyIndexShift");
            sb.newLine();
        }

        sb.newLine();
        sb.appendIndentedLine("const ktokenMatchMask uint64 = ((refTypeMask << 1 << refKTokenSortShift) | refKTokenSortMask) << refKTokenLengthShift << refKTokenIndexShift");
        sb.newLine();

        for (String ktokenSort : inlineManager.ktokenSortNames) {
            sb.append("const ").append(inlineManager.ktokenMatchConstName((ktokenSort)));
            sb.append(" = ((uint64(ktokenRef) << 1 << refKTokenSortShift) | uint64(m.");
            sb.append(ktokenSort);
            sb.append(")) << refKTokenLengthShift << refKTokenIndexShift");
            sb.newLine();
        }

        sb.newLine();
        return sb.toString();
    }

}
