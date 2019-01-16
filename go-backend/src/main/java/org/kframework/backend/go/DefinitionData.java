package org.kframework.backend.go;

import org.kframework.definition.Module;
import org.kframework.kore.KLabel;

import java.util.Set;

class DefinitionData {

    final Module mainModule;
    final Set<KLabel> functions;
    final Set<KLabel> anywhereKLabels;

    public DefinitionData(Module mainModule, Set<KLabel> functions, Set<KLabel> anywhereKLabels) {
        this.mainModule = mainModule;
        this.functions = functions;
        this.anywhereKLabels = anywhereKLabels;
    }
}
