package org.kframework.backend.go.model;

import org.kframework.definition.Module;
import org.kframework.kore.KLabel;

import java.util.Map;
import java.util.Set;

public class DefinitionData {

    public final Module mainModule;
    public final Set<KLabel> functions;
    public final Set<KLabel> anywhereKLabels;
    public final Map<KLabel, FunctionParams> functionParams;
    public final KLabel topCellInitializer;

    public DefinitionData(
            Module mainModule,
            Set<KLabel> functions,
            Set<KLabel> anywhereKLabels,
            Map<KLabel, FunctionParams> functionParams,
            KLabel topCellInitializer) {
        this.mainModule = mainModule;
        this.functions = functions;
        this.anywhereKLabels = anywhereKLabels;
        this.functionParams = functionParams;
        this.topCellInitializer = topCellInitializer;
    }
}
