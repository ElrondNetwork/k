// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.model;

import com.google.common.collect.ListMultimap;
import com.google.common.collect.SetMultimap;
import org.kframework.definition.Module;
import org.kframework.definition.Rule;
import org.kframework.kore.KLabel;

import java.util.Map;
import java.util.Set;

public class DefinitionData {

    public final Module mainModule;
    public final Set<KLabel> functions;
    public final Set<KLabel> anywhereKLabels;
    public final SetMultimap<KLabel, Rule> functionRules;
    public final ListMultimap<KLabel, Rule> anywhereRules;
    public final Map<KLabel, FunctionParams> functionParams;
    public final KLabel topCellInitializer;

    public DefinitionData(
            Module mainModule,
            Set<KLabel> functions,
            Set<KLabel> anywhereKLabels,
            SetMultimap<KLabel, Rule> functionRules,
            ListMultimap<KLabel, Rule> anywhereRules,
            Map<KLabel, FunctionParams> functionParams,
            KLabel topCellInitializer) {
        this.mainModule = mainModule;
        this.functions = functions;
        this.anywhereKLabels = anywhereKLabels;
        this.functionRules = functionRules;
        this.anywhereRules = anywhereRules;
        this.functionParams = functionParams;
        this.topCellInitializer = topCellInitializer;
    }

    public boolean isFunctionOrAnywhere(KLabel klabel) {
        return functions.contains(klabel) || anywhereKLabels.contains(klabel);
    }
}
