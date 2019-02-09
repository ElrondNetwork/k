package org.kframework.backend.go.codegen.rules;

import org.kframework.attributes.Location;
import org.kframework.attributes.Source;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.Lookup;
import org.kframework.backend.go.model.RuleCounter;
import org.kframework.backend.go.model.RuleInfo;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.model.TempVarCounters;
import org.kframework.backend.go.processors.AccumulateRuleVars;
import org.kframework.backend.go.processors.LookupExtractor;
import org.kframework.backend.go.processors.LookupVarExtractor;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.builtin.BooleanUtils;
import org.kframework.compile.ConvertDataStructureToLookup;
import org.kframework.compile.RewriteToTop;
import org.kframework.definition.Rule;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KVariable;
import org.kframework.kore.VisitK;
import org.kframework.utils.errorsystem.KEMException;

import java.util.HashSet;
import java.util.List;
import java.util.NoSuchElementException;
import java.util.Set;

public class RuleWriter {

    private final DefinitionData data;
    private final GoNameProvider nameProvider;
    private final TempVarCounters tempVarCounters = new TempVarCounters();

    public RuleWriter(DefinitionData data, GoNameProvider nameProvider) {
        this.data = data;
        this.nameProvider = nameProvider;
    }

    public RuleInfo writeRule(Rule r, GoStringBuilder sb, RuleType type, RuleCounter ruleCounter,
                              String functionName, FunctionParams functionVars) {
        try {
            int ruleNum = ruleCounter.consumeRuleIndex();
            sb.appendIndentedLine("// rule #" + ruleNum);
            sb.append("\t// ");
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

            // output main LHS
            sb.writeIndent().append("// LHS").newLine();
            Set<KVariable> alreadySeenLhsVariables = new HashSet<>(); // shared between main LHS and lookup LHS
            RuleLhsWriter lhsWriter = new RuleLhsWriter(sb, data, nameProvider, functionVars,
                    accumLhsVars.vars(),
                    accumRhsVars.vars(),
                    alreadySeenLhsVariables,
                    true);
            if (type == RuleType.ANYWHERE || type == RuleType.FUNCTION) {
                KApply kapp = (KApply) left;
                lhsWriter.applyTuple(kapp.klist().items());
            } else {
                lhsWriter.apply(left);
            }

            // output lookups
            writeLookups(sb, ruleNum,
                    functionName, functionVars,
                    lookups,
                    accumLhsVars.vars(),
                    accumRhsVars.vars(),
                    alreadySeenLhsVariables);

            // output requires
            boolean requiresContainsIf = false;
            if (!requires.equals(BooleanUtils.TRUE)) {
                sb.appendIndentedLine("// REQUIRES");
                RuleSideConditionWriter sideCondVisitor = new RuleSideConditionWriter(data, nameProvider,
                        accumLhsVars.vars(), tempVarCounters,
                        sb.getCurrentIndent(), "if ".length());
                sideCondVisitor.apply(requires);
                sideCondVisitor.writeEvalCalls(sb);
                sb.writeIndent().append("if ");
                requiresContainsIf = true;
                sideCondVisitor.writeReturnValue(sb);
                sb.beginBlock();
            } else if (requires.att().contains(PrecomputePredicates.COMMENT_KEY)) {
                // just a comment, so we know what happened
                sb.appendIndentedLine("// REQUIRES precomputed " + requires.att().get(PrecomputePredicates.COMMENT_KEY));
            }

            // output RHS
            sb.appendIndentedLine("// RHS");
            traceLine(sb, type, ruleNum, r);
            RuleRhsWriter rhsWriter = new RuleRhsWriter(data, nameProvider,
                    accumLhsVars.vars(), tempVarCounters,
                    sb.getCurrentIndent(), 0);
            rhsWriter.apply(right);
            rhsWriter.writeEvalCalls(sb);
            sb.writeIndent();
            sb.append("return ");
            rhsWriter.writeReturnValue(sb);
            sb.append(", nil").newLine();

            // done
            sb.endAllBlocks(GoStringBuilder.FUNCTION_BODY_INDENT);
            sb.newLine();

            // return some info regarding the written rule
            boolean alwaysMatches = !lhsWriter.containsIf() && !requiresContainsIf;
            return new RuleInfo(alwaysMatches);
        } catch (NoSuchElementException e) {
            System.err.println(r);
            throw e;
        } catch (KEMException e) {
            e.exception.addTraceFrame("while compiling rule at " + r.att().getOptional(Source.class).map(Object::toString).orElse("<none>") + ":" + r.att().getOptional(Location.class).map(Object::toString).orElse("<none>"));
            throw e;
        }
    }

    private void writeLookups(GoStringBuilder sb, int ruleNum,
                              String functionName, FunctionParams functionArgs,
                              List<Lookup> lookups,
                              RuleVars lhsVars, RuleVars rhsVars,
                              Set<KVariable> alreadySeenLhsVariables) {
        if (lookups.isEmpty()) {
            return;
        }
        sb.appendIndentedLine("// LOOKUPS");
        sb.writeIndent().append("if guard < ").append(ruleNum).beginBlock();

        int lookupIndex = 0;
        for (Lookup lookup : lookups) {
            String reapply = "return " + functionName + "("+ functionArgs.callParameters() +"config, " + ruleNum + ") // reapply";

            sb.appendIndentedLine("// lookup:", lookup.comment());

            RuleLhsWriter lhsWriter = new RuleLhsWriter(sb, data, nameProvider, new FunctionParams(0),
                    lhsVars, rhsVars,
                    alreadySeenLhsVariables,
                    false);
            RuleRhsWriter rhsWriter = new RuleRhsWriter(data, nameProvider,
                    rhsVars, tempVarCounters,
                    sb.getCurrentIndent(), 0);
            rhsWriter.apply(lookup.getRhs());

            switch (lookup.getType()) {
            case MATCH:
                String matchVar = "matchEval" + lookupIndex;
                rhsWriter.writeEvalCalls(sb);
                sb.writeIndent().append(matchVar).append(" := ");
                rhsWriter.writeReturnValue(sb);

                sb.newLine();
                sb.writeIndent().append("if _, isBottom := ").append(matchVar).append(".(m.Bottom); isBottom").beginBlock();
                sb.writeIndent().append(reapply).newLine();
                sb.endOneBlock();

                int ourElseIndent = sb.getCurrentIndent();

                lhsWriter.setNextSubject(matchVar);
                lhsWriter.apply(lookup.getLhs());

                if (lhsWriter.getTopExpressionType() == RuleLhsWriter.ExpressionType.IF) {
                    // defer final else for when we close the corresponding block
                    sb.addCallbackAfterReturningFromBlock(ourElseIndent, sbRef -> {
                        sbRef.append(" else").beginBlock();
                        sbRef.writeIndent().append(reapply).newLine();
                        sbRef.endOneBlockNoNewline();
                    });
                }
                break;
            case SETCHOICE:
                lhsWriter = new RuleLhsWriter(sb, data, nameProvider, new FunctionParams(0),
                        lhsVars, rhsVars,
                        alreadySeenLhsVariables,
                        false);
                writeChoiceLookup(
                        sb, lookup,
                        "setChoice" + lookupIndex, "m.Set",
                        reapply,
                        lhsWriter, rhsWriter);
                break;
            case MAPCHOICE:
                lhsWriter = new RuleLhsWriter(sb, data, nameProvider, new FunctionParams(0),
                        lhsVars, rhsVars,
                        alreadySeenLhsVariables,
                        false);
                writeChoiceLookup(
                        sb, lookup,
                        "mapChoice" + lookupIndex, "m.Map",
                        reapply,
                        lhsWriter, rhsWriter);
                break;
            default:
                throw KEMException.internalError("Unexpected lookup type");
            }

            lookupIndex++;
        }
    }

    private void writeChoiceLookup(
            GoStringBuilder sb, Lookup lookup,
            String varPrefix, String expectedKType,
            String reapply,
            RuleLhsWriter lhsWriter, RuleRhsWriter rhsWriter) {

        String setChoiceVar = varPrefix + "Eval";
        String setVar = varPrefix + "Obj";
        String isSetVar = varPrefix + "TypeOk";
        String setElemVar = varPrefix + "Elem";
        String choiceVar = varPrefix + "Result";
        String errVar = varPrefix + "Err";

        rhsWriter.writeEvalCalls(sb);
        sb.writeIndent().append(setChoiceVar).append(" := ");
        rhsWriter.writeReturnValue(sb);
        sb.newLine();

        sb.writeIndent()
                .append(setVar).append(", ").append(isSetVar).append(" := ")
                .append(setChoiceVar).append(".(").append(expectedKType).append(")").newLine();
        sb.writeIndent().append("if !").append(isSetVar).beginBlock();
        sb.writeIndent().append(reapply).newLine();
        sb.endOneBlock();

        sb.appendIndentedLine("var ", choiceVar, " = m.InternedBottom");
        int forIndent = sb.getCurrentIndent();
        sb.writeIndent().append("for ").append(setElemVar).append(" := range ").append(setVar).append(".Data").beginBlock();
        sb.appendIndentedLine("var ", errVar, " error");

        // this will be after the end of the for, reapply if we didn't hit return in the for loop
        sb.addCallbackAfterReturningFromBlock(forIndent, s -> {
            s.newLine();
            s.writeIndent().append("if ").append(choiceVar).append(" == m.InternedBottom").beginBlock();
            s.appendIndentedLine(reapply);
            s.endOneBlock();
            s.appendIndentedLine("return ", choiceVar, ", nil");
        });

        lhsWriter.setNextSubject(setElemVar);
        lhsWriter.apply(lookup.getLhs());

        // the function goes inside
        // I made it a function just because I didn't feel like changing the RHS returns
        // it can also be done without this function, but then the RHS must assign the result to choice var instead of returning
        int funcIndent = sb.getCurrentIndent();
        sb.writeIndent().append(choiceVar).append(", ").append(errVar).append(" = ");
        sb.append("func() (m.K, error)").beginBlock();

        sb.addCallbackBeforeReturningFromBlock(funcIndent, s -> {
            s.newLine();
            s.appendIndentedLine("return m.InternedBottom, nil // #setChoice end");
            s.endOneBlock();
        });

        sb.addCallbackAfterReturningFromBlock(funcIndent, s -> {
            s.append("()").newLine(); // function call
            s.writeIndent().append("if ").append(errVar).append(" != nil").beginBlock();
            s.appendIndentedLine("return m.NoResult, ", errVar);
            s.endOneBlock();
        });
    }

    private static String traceRuleTypeString(RuleType ruleType) {
        switch (ruleType) {
        case FUNCTION:
            return "FUNC";
        case ANYWHERE:
            return "ANYW";
        case REGULAR:
            return "STEP";
        case PATTERN:
            return "PATT";
        default:
            return "????";
        }
    }

    private static void traceLine(GoStringBuilder sb, RuleType ruleType, int ruleNum, Rule r) {
        sb.appendIndentedLine(
                "traceRuleApply(\"",
                traceRuleTypeString(ruleType),
                "\", ",
                Integer.toString(ruleNum),
                ", ",
                GoStringUtil.enquotedRuleComment(r),
                ")");
    }

    public static int numLookups(Rule r) {
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

    public static boolean hasLookups(Rule r) {
        return numLookups(r) > 0;
    }

}
