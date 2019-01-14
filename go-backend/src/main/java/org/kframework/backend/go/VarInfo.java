package org.kframework.backend.go;

import com.google.common.collect.HashMultimap;
import com.google.common.collect.SetMultimap;
import org.kframework.kore.K;
import org.kframework.kore.KLabel;
import org.kframework.kore.KVariable;

import java.util.HashMap;
import java.util.Map;

class VarInfo {
    final SetMultimap<KVariable, String> vars;
    final Map<String, KLabel> listVars;
    final Map<K, String> termCache;

    VarInfo() { this(HashMultimap.create(), new HashMap<>(), new HashMap<>()); }

    VarInfo(VarInfo vars) {
        this(HashMultimap.create(vars.vars), new HashMap<>(vars.listVars), new HashMap<>(vars.termCache));
    }

    VarInfo(SetMultimap<KVariable, String> vars, Map<String, KLabel> listVars, Map<K, String> termCache) {
        this.vars = vars;
        this.listVars = listVars;
        this.termCache = termCache;
    }
}
