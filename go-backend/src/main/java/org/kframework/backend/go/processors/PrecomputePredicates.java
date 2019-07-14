// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.processors;

import org.kframework.backend.go.codegen.GoBuiltin;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.builtin.Sorts;
import org.kframework.kil.Attribute;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KORE;
import org.kframework.kore.KToken;
import org.kframework.kore.KVariable;
import org.kframework.kore.Sort;
import org.kframework.kore.TransformK;

import java.util.Optional;

/**
 * Replaces certain function calls with the result, where that result is known from compile time.
 */
public class PrecomputePredicates extends TransformK {

    public static final String COMMENT_KEY = "precomputeComment";

    private final DefinitionData data;
    private final RuleVars lhsVars;

    public PrecomputePredicates(DefinitionData data, RuleVars lhsVars) {
        this.data = data;
        this.lhsVars = lhsVars;
    }

    public static KToken trueTokenWithComment(String precomputeComment) {
        return KORE.KToken(
                "true",
                Sorts.Bool(),
                KORE.Att().add(COMMENT_KEY, precomputeComment));
    }

    public static boolean isTrueToken(K k) {
        if (k instanceof KToken) {
            KToken kt = (KToken)k;
            return kt.sort().equals(Sorts.Bool()) && kt.s().equals("true");
        }
        return false;
    }

    @Override
    public K apply(KApply k) {
        if (data.isFunctionOrAnywhere(k.klabel())) {

            // precompute true scenarios
            if (data.mainModule.attributesFor().apply(k.klabel()).contains(Attribute.PREDICATE_KEY, Sort.class)) {
                Sort s = data.mainModule.attributesFor().apply(k.klabel()).get(Attribute.PREDICATE_KEY, Sort.class);

                // is K
                if (s.equals(Sorts.K()) && k.klist().items().size() == 1) {
                    return trueTokenWithComment("isK");
                }

                // is KItem
                if (s.equals(Sorts.KItem()) && k.klist().items().size() == 1) {
                    return trueTokenWithComment("isKItem");
                }

                // is <Sort>
                if (data.mainModule.sortAttributesFor().contains(s)) {
                    String hook2 = data.mainModule.sortAttributesFor().apply(s).<String>getOptional("hook").orElse("");
                    if (GoBuiltin.PREDICATE_HOOKS.contains(hook2)) {
                        if (k.klist().items().size() == 1 && k.klist().items().get(0) instanceof KVariable) {
                            KVariable kvar = (KVariable) k.klist().items().get(0);
                            if (lhsVars.containsVar(kvar)) {
                                Optional<Sort> varSort = kvar.att().getOptional(Sort.class);
                                if (varSort.isPresent() && varSort.get().equals(s)) {
                                    return trueTokenWithComment("is" + s.name() + "(" + kvar.name() + ")");
                                }
                            }
                        }
                    }
                }
            }

            // collapse (precomputed true) && (precomputed true)
            String hook = data.mainModule.attributesFor().apply(k.klabel()).<String>getOptional(Attribute.HOOK_KEY).orElse("");
            if (hook.equals("BOOL.and") || hook.equals("BOOL.andThen")) {
                KApply kappTransf = (KApply) super.apply(k);
                assert kappTransf.klist().items().size() == 2;
                if (isTrueToken(kappTransf.klist().items().get(0)) &&
                        isTrueToken(kappTransf.klist().items().get(1))) {
                    String comm1 = kappTransf.klist().items().get(0).att().getOption(COMMENT_KEY).getOrElse(() -> "KToken");
                    String comm2 = kappTransf.klist().items().get(1).att().getOption(COMMENT_KEY).getOrElse(() -> "KToken");
                    return trueTokenWithComment(comm1 + " && " + comm2);
                }

                return kappTransf;
            }
        }

        return super.apply(k);
    }
}
