package org.kframework.backend.go.model;

import org.kframework.kore.KVariable;

public interface TempVarManager {

    String oneTimeVariableMVRef(String varName);

    String addKVariableMVRef(KVariable k);

    String getKVariableMVRef(KVariable k);

    String evalBoolVarRef(String varName);

    int getNrVars();

    int getNrBoolVars();
}
