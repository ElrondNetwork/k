// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen.rules;

import org.kframework.attributes.Location;
import org.kframework.attributes.Source;
import org.kframework.backend.go.codegen.inline.RuleLhsMatchWriter;
import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeBuilder;
import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.backend.go.codegen.lhstree.model.LhsLeafTreeNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsTopTreeNode;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionInfo;
import org.kframework.backend.go.model.Lookup;
import org.kframework.backend.go.model.RuleInfo;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.model.VarContainer;
import org.kframework.backend.go.processors.AccumulateRuleVars;
import org.kframework.backend.go.processors.LookupExtractor;
import org.kframework.backend.go.processors.LookupVarExtractor;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.compile.RewriteToTop;
import org.kframework.definition.Rule;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KVariable;
import org.kframework.utils.errorsystem.KEMException;

import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Set;

public class RuleWriter {

    private final DefinitionData data;
    private final GoNameProvider nameProvider;
    private final RuleLhsMatchWriter matchWriter;

    public RuleWriter(DefinitionData data, GoNameProvider nameProvider, RuleLhsMatchWriter matchWriter) {
        this.data = data;
        this.nameProvider = nameProvider;
        this.matchWriter = matchWriter;
    }

    public RuleInfo writeRule(Map<Integer, Rule> rules,
                              GoStringBuilder sb, GoStringBuilder rhsSb,
                              RuleType type,
                              FunctionInfo functionInfo) {

        int initialIndent = sb.getCurrentIndent();

        RuleInfo result = new RuleInfo();
        LhsTopTreeNode topNode = null;
        for (Map.Entry<Integer, Rule> entry : rules.entrySet()) {
            Integer ruleNum = entry.getKey();
            Rule r = entry.getValue();
            try {

                K left = RewriteToTop.toLeft(r.body());
                K requires = r.requires();
                K right = RewriteToTop.toRight(r.body());

                // we need the variables beforehand, so we retrieve them here
                AccumulateRuleVars accumLhsVars = new AccumulateRuleVars(nameProvider);
                accumLhsVars.apply(left);

                // lookups!
                LookupExtractor lookupExtractor = new LookupExtractor();
                requires = lookupExtractor.apply(requires); // also, lookups are eliminated from requires
                List<Lookup> lookups = lookupExtractor.getExtractedLookups();

                // some evaluations can be precomputed
                PrecomputePredicates optimizeTransf = new PrecomputePredicates(
                        data, accumLhsVars.vars());
                requires = optimizeTransf.apply(requires);
                right = optimizeTransf.apply(right);

                // check which variables are actually used in requires or in rhs
                // note: this has to happen *after* PrecomputePredicates does its job
                AccumulateRuleVars accumRhsVars = new AccumulateRuleVars(nameProvider);
                accumRhsVars.apply(requires);
                accumRhsVars.apply(right);

                // also collect vars from lookups
                new LookupVarExtractor(accumLhsVars, accumRhsVars).apply(lookups);

                // var indexes
                VarContainer vars = new VarContainer(
                        accumLhsVars.vars(),
                        accumRhsVars.vars());

                // output main LHS
                //sb.writeIndent().append("// LHS").newLine();
                Set<KVariable> alreadySeenLhsVariables = new HashSet<>(); // shared between main LHS and lookup LHS

                // LHS
                RuleLhsTreeBuilder treeBuilder = new RuleLhsTreeBuilder(
                        data,
                        nameProvider,
                        functionInfo.arguments,
                        alreadySeenLhsVariables);
                if (type == RuleType.ANYWHERE || type == RuleType.FUNCTION) {
                    KApply kapp = (KApply) left;
                    treeBuilder.applyTopArgs(kapp.klist().items());
                } else {
                    treeBuilder.applyTopArgs(Collections.singleton(left));
                }

                // leaf
                LhsLeafTreeNode leaf = new LhsLeafTreeNode(treeBuilder.getLastNode(),
                        type, ruleNum,
                        functionInfo,
                        r, lookups, requires, right,
                        alreadySeenLhsVariables);
                treeBuilder.addNode(leaf);

//            RuleLhsWriter lhsWriter = new RuleLhsWriter(sb, data,
//                    nameProvider, matchWriter,
//                    functionInfo.arguments,
//                    vars,
//                    alreadySeenLhsVariables,
//                    false);
                if (topNode == null) {
                    topNode = treeBuilder.topNode;
                } else {
                    topNode.mergeTree(treeBuilder.topNode);
                }
            } catch (NoSuchElementException e) {
                System.err.println(r);
                throw e;
            } catch (KEMException e) {
                e.exception.addTraceFrame("while compiling rule at " + r.att().getOptional(Source.class).map(Object::toString).orElse("<none>") + ":" + r.att().getOptional(Location.class).map(Object::toString).orElse("<none>"));
                throw e;
            }
        }

        if (topNode != null) {
            RuleLhsTreeWriter treeWriter = new RuleLhsTreeWriter(
                    sb, rhsSb,
                    data,
                    nameProvider, matchWriter);
            treeWriter.writeLhsTree(topNode);

            result.maxNrVars = topNode.maxNrVars();
            result.maxNrBoolVars = topNode.maxNrBoolVars();
        }

        // done
        sb.endAllBlocks(initialIndent);
        sb.newLine();

        // return some info regarding the written rules
        return result;
    }

}
