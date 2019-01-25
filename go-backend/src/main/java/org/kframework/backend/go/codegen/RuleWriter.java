package org.kframework.backend.go.codegen;

import org.kframework.attributes.Location;
import org.kframework.attributes.Source;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.Lookup;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.processors.AccumulateRuleVars;
import org.kframework.backend.go.processors.LookupExtractor;
import org.kframework.backend.go.processors.LookupVarExtractor;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.builtin.BooleanUtils;
import org.kframework.compile.ConvertDataStructureToLookup;
import org.kframework.compile.RewriteToTop;
import org.kframework.definition.Rule;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.VisitK;
import org.kframework.utils.errorsystem.KEMException;

import java.util.List;
import java.util.NoSuchElementException;

public class RuleWriter {

    private final DefinitionData data;

    public RuleWriter(DefinitionData data) {
        this.data = data;
    }

    public void convertFunction(List<Rule> rules, GoStringBuilder sb, RuleType type, FunctionParams functionVars) {
        int ruleNum = 0;
        for (Rule r : rules) {
            sb.appendIndentedLine("// rule");
            ruleNum = convert(r, sb, type, ruleNum, functionVars);
        }
    }

    public int convert(Rule r, GoStringBuilder sb, RuleType type, int ruleNum,
                       FunctionParams functionVars) {
        try {
            GoStringUtil.appendRuleComment(sb, r);

            K left = RewriteToTop.toLeft(r.body());
            K requires = r.requires();
            K right = RewriteToTop.toRight(r.body());

            // we need the variables beforehand, so we retrieve them here
            AccumulateRuleVars accumLhsVars = new AccumulateRuleVars();
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
            AccumulateRuleVars accumRhsVars = new AccumulateRuleVars();
            accumRhsVars.apply(requires);
            accumRhsVars.apply(right);

            // also collect vars from lookups
            new LookupVarExtractor(accumLhsVars, accumRhsVars).apply(lookups);

            // output LHS
            sb.writeIndent().append("// LHS").newLine();
            GoLhsVisitor lhsVisitor = new GoLhsVisitor(sb, data, functionVars,
                    accumLhsVars.vars(),
                    accumRhsVars.vars());
            if (type == RuleType.ANYWHERE || type == RuleType.FUNCTION) {
                KApply kapp = (KApply) left;
                lhsVisitor.applyTuple(kapp.klist().items());
            } else {
                lhsVisitor.apply(left);
            }

            // output lookups
            writeLookups(sb, ruleNum, lookups,
                    accumLhsVars.vars(),
                    accumRhsVars.vars());

            // output requires
            if (!requires.equals(BooleanUtils.TRUE)) {
                sb.appendIndentedLine("// REQUIRES");
                GoSideConditionVisitor sideCondVisitor = new GoSideConditionVisitor(data,
                        accumLhsVars.vars(), sb.getCurrentIndent(), "if ".length());
                sideCondVisitor.apply(requires);
                sideCondVisitor.writeEvalCalls(sb);
                sb.writeIndent().append("if ");
                sideCondVisitor.writeReturnValue(sb);
                sb.beginBlock();
            } else if (requires.att().contains(PrecomputePredicates.COMMENT_KEY)) {
                // just a comment, so we know what happened
                sb.appendIndentedLine("// REQUIRES precomputed " + requires.att().get(PrecomputePredicates.COMMENT_KEY));
            }

            // output RHS
            sb.appendIndentedLine("// RHS");
            GoRhsVisitor rhsVisitor = new GoRhsVisitor(data,
                    accumLhsVars.vars(), sb.getCurrentIndent(), 0);
            rhsVisitor.apply(right);
            rhsVisitor.writeEvalCalls(sb);
            sb.writeIndent();
            sb.append("return ");
            rhsVisitor.writeReturnValue(sb);
            sb.append(", nil").newLine();

            // done
            sb.endAllBlocks(GoStringBuilder.FUNCTION_BODY_INDENT);
            sb.newLine();
            return ruleNum + 1;
        } catch (NoSuchElementException e) {
            System.err.println(r);
            throw e;
        } catch (KEMException e) {
            e.exception.addTraceFrame("while compiling rule at " + r.att().getOptional(Source.class).map(Object::toString).orElse("<none>") + ":" + r.att().getOptional(Location.class).map(Object::toString).orElse("<none>"));
            throw e;
        }
    }

    private void writeLookups(GoStringBuilder sb, int ruleNum, List<Lookup> lookups, RuleVars lhsVars, RuleVars rhsVars) {
        if (lookups.isEmpty()) {
            return;
        }
        sb.appendIndentedLine("// LOOKUPS");
        sb.writeIndent().append("if guard > ").append(ruleNum).beginBlock();

        int matchIndex = 0;
        for (Lookup lookup : lookups) {
            switch (lookup.getType()) {
            case MATCH:
                KApply k = lookup.getContent();
                String matchVar = "e" + matchIndex;
                matchIndex++;
                String reapply = "return stepLookups(c, config, " + ruleNum + ") // reapply";

                GoRhsVisitor rhsVisitor = new GoRhsVisitor(data, rhsVars, sb.getCurrentIndent(), 0);
                rhsVisitor.apply(k.klist().items().get(1));
                rhsVisitor.writeEvalCalls(sb);
                sb.writeIndent().append(matchVar).append(" := ");
                rhsVisitor.writeReturnValue(sb);

                sb.newLine();
                sb.writeIndent().append("if _, isBottom := ").append(matchVar).append(".(Bottom); isBottom").beginBlock();
                sb.writeIndent().append(reapply).newLine();
                sb.endOneBlockNoNewline().append(" else").beginBlock();

                int ourElseIndent = sb.getCurrentIndent();

                GoLhsVisitor lhsVisitor = new GoLhsVisitor(sb, data, new FunctionParams(0),
                        lhsVars,
                        rhsVars);
                lhsVisitor.setNextSubject(matchVar);
                lhsVisitor.apply(k.klist().items().get(0));

                if (lhsVisitor.getTopExpressionType() == GoLhsVisitor.ExpressionType.IF) {
                    // defer final else for when we close the corresponding block
                    sb.addCallbackWhenReturningFromBlock(ourElseIndent, sbRef -> {
                        sbRef.append(" else").beginBlock();
                        sbRef.writeIndent().append(reapply).newLine();
                        sbRef.endOneBlockNoNewline();
                    });
                }
                break;
            default:
                throw KEMException.internalError("Unexpected lookup type");
            }

        }

    }

    static int numLookups(Rule r) {
        class Holder {
            int i;
        }
        Holder h = new Holder();
        new VisitK() {
            @Override
            public void apply(KApply k) {
                if (ConvertDataStructureToLookup.isLookupKLabel(k)) {
                    h.i++;
                }
                super.apply(k);
            }
        }.apply(r.requires());
        return h.i;
    }

    static boolean hasLookups(Rule r) {
        return numLookups(r) > 0;
    }

}
