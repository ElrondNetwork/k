package org.kframework.backend.go.processors;

import org.apache.commons.lang3.NotImplementedException;
import org.kframework.backend.go.model.Lookup;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.TransformK;
import org.kframework.utils.errorsystem.KEMException;

import java.util.ArrayList;
import java.util.List;

/**
 * Extracts all the lookups into a list and replace them with Bool(true) where they were. <br />
 * To be called on the 'requires' part of the rule.
 */
public class LookupExtractor extends TransformK {

    private final List<Lookup> extractedLookups = new ArrayList<>();

    public LookupExtractor() {
    }

    public List<Lookup> getExtractedLookups() {
        return extractedLookups;
    }

    public boolean lookupsFound() {
        return !extractedLookups.isEmpty();
    }

    @Override
    public K apply(KApply k) {
        if (k.klabel().name().equals("#match")) {
            if (k.klist().items().size() != 2) {
                throw KEMException.internalError("Unexpected arity of lookup: " + k.klist().size(), k);
            }
            extractedLookups.add(new Lookup(Lookup.Type.MATCH, k));
            return PrecomputePredicates.trueTokenWithComment("lookup");
        } else if (k.klabel().name().equals("#setChoice")) {
            throw new NotImplementedException("#setChoice");
        } else if (k.klabel().name().equals("#mapChoice")) {
            throw new NotImplementedException("#mapChoice");
        } else if (k.klabel().name().equals("#filterMapChoice")) {
            throw new NotImplementedException("#filterMapChoice");
        }

        return super.apply(k);
    }

}
