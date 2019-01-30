package org.kframework.backend.go.strings;

import org.kframework.kore.KLabel;

import java.util.HashMap;
import java.util.Map;

public class GoNameProviderProperMinKLabel extends GoNameProviderProper {

    private int nextKlIndex = 0;
    private final Map<String, Integer> klabelIndexes = new HashMap<>();

    @Override
    public String klabelVariableName(KLabel klabel) {
        Integer klIndex = klabelIndexes.get(klabel.name());
        if (klIndex == null) {
            klIndex = nextKlIndex;
            nextKlIndex++;
            klabelIndexes.put(klabel.name(), klIndex);
        }
        return "kl" + klIndex + " /* " + klabel.name() + " */";
    }
}
