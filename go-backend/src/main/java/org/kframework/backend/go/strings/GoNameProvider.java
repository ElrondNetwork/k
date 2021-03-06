// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.strings;

import org.kframework.kore.KLabel;
import org.kframework.kore.Sort;

public interface GoNameProvider {

    String klabelVariableName(KLabel klabel);

    String sortVariableName(Sort sort);

    String evalFunctionName(KLabel lbl);

    String constFunctionName(KLabel lbl);

    String memoFunctionName(KLabel lbl);

    String memoTableName(KLabel lbl);

    String ruleVariableName(String varName);
}
