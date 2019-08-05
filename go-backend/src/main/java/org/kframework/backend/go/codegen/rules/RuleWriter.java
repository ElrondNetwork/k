// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen.rules;

import org.kframework.attributes.Location;
import org.kframework.attributes.Source;
import org.kframework.backend.go.codegen.inline.RuleLhsMatchWriter;
import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeBuilder;
import org.kframework.backend.go.codegen.lhstree.RuleLhsTreeWriter;
import org.kframework.backend.go.codegen.lhstree.model.LhsLeafTreeNode;
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
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.compile.RewriteToTop;
import org.kframework.definition.Rule;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KVariable;
import org.kframework.utils.errorsystem.KEMException;

import java.util.Collections;
import java.util.HashSet;
import java.util.List;
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

    public RuleInfo writeRule(Rule r, GoStringBuilder sb, RuleType type, int ruleNum,
                              FunctionInfo functionInfo) {
        try {
            int ruleIndent = sb.getCurrentIndent();

            sb.appendIndentedLine("// rule #" + ruleNum);
            appendSourceComment(sb, r);
            sb.writeIndent().append("// ");
            GoStringUtil.appendRuleComment(sb, r);
            sb.newLine();

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

            // if !matched
            sb.writeIndent().append("if !matched").beginBlock();

            // output main LHS
            sb.writeIndent().append("// LHS").newLine();
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
                    vars,
                    alreadySeenLhsVariables);
            treeBuilder.addNode(leaf);

            RuleLhsTreeWriter treeWriter = new RuleLhsTreeWriter(sb, data,
                    nameProvider, matchWriter, vars);
            treeWriter.writeLhsTree(treeBuilder.topNode);

//            RuleLhsWriter lhsWriter = new RuleLhsWriter(sb, data,
//                    nameProvider, matchWriter,
//                    functionInfo.arguments,
//                    vars,
//                    alreadySeenLhsVariables,
//                    false);


            // done
            sb.endAllBlocks(ruleIndent);
            sb.newLine();

            // return some info regarding the written rule
            //boolean alwaysMatches = !lhsWriter.containsIf() && !requiresContainsIf;
            return new RuleInfo(
                    false,
                    vars.varIndexes.getNrVars(), vars.varIndexes.getNrBoolVars());
        } catch (NoSuchElementException e) {
            System.err.println(r);
            throw e;
        } catch (KEMException e) {
            e.exception.addTraceFrame("while compiling rule at " + r.att().getOptional(Source.class).map(Object::toString).orElse("<none>") + ":" + r.att().getOptional(Location.class).map(Object::toString).orElse("<none>"));
            throw e;
        }
    }

    private static void appendSourceComment(GoStringBuilder sb, Rule r) {
        String source;
        if (r.source().isPresent()) {
            source = r.source().get().source();
            if (source.contains("/")) {
                source = source.substring(source.lastIndexOf("/") + 1);
            }
        } else {
            source = "?";
        }
        String startLine;
        if (r.location().isPresent()) {
            startLine = Integer.toString(r.location().get().startLine());
        } else {
            startLine = "?";
        }
        sb.appendIndentedLine("// source: ", source, " @", startLine);
    }
}
