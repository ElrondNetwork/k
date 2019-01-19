package org.kframework.backend.go.model;

import org.kframework.definition.Module;
import org.kframework.kore.KLabel;

import java.util.Set;

public class DefinitionData {

    public final Module mainModule;
    public final Set<KLabel> functions;
    public final Set<KLabel> anywhereKLabels;

    public DefinitionData(Module mainModule, Set<KLabel> functions, Set<KLabel> anywhereKLabels) {
        this.mainModule = mainModule;
        this.functions = functions;
        this.anywhereKLabels = anywhereKLabels;
    }
}
