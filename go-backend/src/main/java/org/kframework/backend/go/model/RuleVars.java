// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.model;

import org.kframework.kore.KApply;
import org.kframework.kore.KLabel;
import org.kframework.kore.KVariable;

import java.util.HashMap;
import java.util.Map;

public class RuleVars {
    private final Map<KVariable, String> kVarToName = new HashMap<>();
    private final Map<KVariable, Integer> kVarCount = new HashMap<>();
    public final Map<String, KLabel> listVars = new HashMap<>(); // TEMP

    private final Map<KApplySignature, Integer> signatureToKApplyCount = new HashMap<>();
    public final Map<KApplySignature, String> signatureToKApplyVarName = new HashMap<>();

    public RuleVars() {
    }

    public void putVar(KVariable kv, String varName) {
        kVarToName.put(kv, varName);
    }

    public boolean containsVar(KVariable kv) {
        return kVarToName.containsKey(kv);
    }

    public String getVarName(KVariable kv) {
        return kVarToName.get(kv);
    }

    public void incrementVarCount(KVariable kv) {
        Integer currentCount = kVarCount.get(kv);
        if (currentCount == null) {
            kVarCount.put(kv, 1);
        } else {
            kVarCount.put(kv, currentCount + 1);
        }
    }

    public int getVarCount(KVariable kv) {
        Integer currentCount = kVarCount.get(kv);
        if (currentCount == null) {
            return 0;
        } else {
            return currentCount;
        }
    }

    public void incrementKApplySignatureCount(KApply kapp) {
        KApplySignature signature = KApplySignature.of(kapp);
        if (!signatureToKApplyCount.containsKey(signature)) {
            signatureToKApplyCount.put(signature, 1);
        } else {
            Integer previousCout = signatureToKApplyCount.get(signature);
            signatureToKApplyCount.put(signature, previousCout + 1);
        }
    }

    public int getKApplySignatureCount(KApply kapp) {
        return signatureToKApplyCount.getOrDefault(KApplySignature.of(kapp), 0);
    }
}
