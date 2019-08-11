package org.kframework.backend.go.codegen.lhstree;

import org.kframework.backend.go.codegen.inline.RuleLhsMatchWriter;
import org.kframework.backend.go.codegen.lhstree.model.LhsLeafTreeNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsTopTreeNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsTreeNode;
import org.kframework.backend.go.codegen.rules.RuleLhsWriter;
import org.kframework.backend.go.codegen.rules.RuleRhsWriter;
import org.kframework.backend.go.codegen.rules.RuleRhsWriterBase;
import org.kframework.backend.go.codegen.rules.RuleSideConditionWriter;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionInfo;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.Lookup;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.model.TempVarManager;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.builtin.BooleanUtils;
import org.kframework.compile.ConvertDataStructureToLookup;
import org.kframework.definition.Rule;
import org.kframework.kore.KApply;
import org.kframework.kore.VisitK;
import org.kframework.unparser.ToKast;
import org.kframework.utils.errorsystem.KEMException;

import java.util.List;

public class RuleLhsTreeWriter {
    public final GoStringBuilder sb;
    public final GoStringBuilder rhsSb;
    public final DefinitionData data;
    public final GoNameProvider nameProvider;
    public final RuleLhsMatchWriter matchWriter;

    public RuleLhsTreeWriter(GoStringBuilder sb, GoStringBuilder rhsSb, DefinitionData data, GoNameProvider nameProvider, RuleLhsMatchWriter matchWriter) {
        this.sb = sb;
        this.rhsSb = rhsSb;
        this.data = data;
        this.nameProvider = nameProvider;
        this.matchWriter = matchWriter;
    }

    public void writeLhsTree(LhsTopTreeNode top) {
        top.findRulesBelow();
        writeLhsNode(top);
    }

    private void writeLhsNode(LhsTreeNode node) {
        int currentIndent = sb.getCurrentIndent();
        node.writeRuleInfo(this);
        node.write(this);
        for (LhsTreeNode child : node.successors) {
            child.predecessor = node;
            writeLhsNode(child);
        }
        sb.endAllBlocks(currentIndent);
    }

    public void writeLeaf(LhsLeafTreeNode leafNode) {
        sb.appendIndentedLine("// rule #" + leafNode.ruleNum);
        appendSourceComment(sb, leafNode.rule);
        sb.writeIndent().append("// ");
        GoStringUtil.appendRuleComment(sb, leafNode.rule);
        sb.newLine();

        // if !matched
        sb.writeIndent().append("if !matched").beginBlock();

        // output lookups
        writeLookups(sb, leafNode.ruleNum,
                leafNode.functionInfo,
                leafNode.lookups,
                leafNode);

        // output leafNode.requires
        if (!leafNode.requires.equals(BooleanUtils.TRUE)) {
            sb.appendIndentedLine("// REQUIRES ", ToKast.apply(leafNode.requires));
            RuleSideConditionWriter sideCondVisitor = new RuleSideConditionWriter(data, nameProvider,
                    leafNode,
                    sb.getCurrentIndent());
            sideCondVisitor.apply(leafNode.requires);
            sideCondVisitor.writeEvalCalls(sb);
            sb.writeIndent().append("if ");
            sideCondVisitor.writeReturnValue(sb);
            sb.beginBlock();
        } else if (leafNode.requires.att().contains(PrecomputePredicates.COMMENT_KEY)) {
            // just a comment, so we know what happened
            sb.appendIndentedLine("// REQUIRES precomputed " + leafNode.requires.att().get(PrecomputePredicates.COMMENT_KEY));
        }

        // output RHS
        sb.appendIndentedLine("// RHS");
        if (rhsSb == null) {
            traceLine(sb, leafNode.type, leafNode.ruleNum, leafNode.rule);
            RuleRhsWriter rhsWriter = new RuleRhsWriter(data, nameProvider,
                    leafNode,
                    sb.getCurrentIndent(),
                    true, leafNode.functionInfo);
            rhsWriter.apply(leafNode.right);
            rhsWriter.writeEvalCalls(sb);
            sb.writeIndent();
            rhsWriter.writeReturnValue(sb);
            sb.newLine();
        } else {
            String rhsFuncName = "stepRHS" + leafNode.ruleNum;
            sb.appendIndentedLine("return i.", rhsFuncName, "(v, bv, config)");

            rhsSb.writeIndent().append("func (i *Interpreter) ").append(rhsFuncName);
            rhsSb.append("(v []m.KReference, bv []bool, config KReference) (m.KReference, error)").beginBlock();

            traceLine(rhsSb, leafNode.type, leafNode.ruleNum, leafNode.rule);
            RuleRhsWriter rhsWriter = new RuleRhsWriter(data, nameProvider,
                    leafNode,
                    rhsSb.getCurrentIndent(),
                    true, leafNode.functionInfo);
            rhsWriter.apply(leafNode.right);
            rhsWriter.writeEvalCalls(rhsSb);
            rhsSb.writeIndent();
            rhsWriter.writeReturnValue(rhsSb);
            rhsSb.newLine();
            rhsSb.endAllBlocks(0);
            rhsSb.newLine();
        }

    }

    private void writeLookups(GoStringBuilder mainSb, int ruleNum,
                              FunctionInfo functionInfo,
                              List<Lookup> lookups,
                              TempVarManager varManager) {
        if (lookups.isEmpty()) {
            return;
        }
        mainSb.appendIndentedLine("// LOOKUPS");
        mainSb.writeIndent().append("if guard < ").append(ruleNum).beginBlock();

        int lookupIndex = 0;
        for (Lookup lookup : lookups) {
            String reapply = "return i." + functionInfo.goName + "(" + functionInfo.arguments.callParameters() + "config, " + ruleNum + ") // reapply";

            mainSb.appendIndentedLine("// lookup:", lookup.comment());

            RuleLhsWriter lhsWriter = new RuleLhsWriter(mainSb, data,
                    nameProvider, matchWriter,
                    new FunctionParams(0),
                    varManager,
                    false);
            RuleRhsWriter rhsWriter = new RuleRhsWriter(data, nameProvider,
                    varManager,
                    mainSb.getCurrentIndent(),
                    false, functionInfo);
            rhsWriter.apply(lookup.getRhs());

            switch (lookup.getType()) {
            case MATCH:
                String matchVar = varManager.oneTimeVariableMVRef("matchEval" + lookupIndex);
                rhsWriter.writeEvalCalls(mainSb);
                mainSb.writeIndent().append(matchVar).append(" = ");
                rhsWriter.writeReturnValue(mainSb);

                mainSb.newLine();
                mainSb.writeIndent().append("if ");
                matchWriter.appendBottomMatch(mainSb, matchVar);
                mainSb.beginBlock();
                mainSb.writeIndent().append(reapply).newLine();
                mainSb.endOneBlock();

                int ourElseIndent = mainSb.getCurrentIndent();

                lhsWriter.setNextSubject(matchVar);
                lhsWriter.apply(lookup.getLhs());

                if (lhsWriter.getTopExpressionType() == RuleLhsWriter.ExpressionType.IF) {
                    // defer final else for when we close the corresponding block
                    mainSb.addCallbackAfterReturningFromBlock(ourElseIndent, sbRef -> {
                        sbRef.append(" else").beginBlock();
                        sbRef.writeIndent().append(reapply).newLine();
                        sbRef.endOneBlockNoNewline();
                    });
                }
                break;
            case SETCHOICE:
                lhsWriter = new RuleLhsWriter(mainSb, data,
                        nameProvider, matchWriter,
                        new FunctionParams(0),
                        varManager,
                        false);
                writeChoiceLookup(
                        mainSb, lookup,
                        "setChoice" + lookupIndex, "GetSetObject",
                        varManager,
                        reapply,
                        lhsWriter, rhsWriter);
                break;
            case MAPCHOICE:
                lhsWriter = new RuleLhsWriter(mainSb, data,
                        nameProvider, matchWriter,
                        new FunctionParams(0),
                        varManager,
                        false);
                writeChoiceLookup(
                        mainSb, lookup,
                        "mapChoice" + lookupIndex, "GetMapObject",
                        varManager,
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
            GoStringBuilder mainSb, Lookup lookup,
            String varPrefix, String expectedKType,
            TempVarManager varManager,
            String reapply,
            RuleLhsWriter lhsWriter, RuleRhsWriterBase rhsWriter) {

        String setChoiceVar = varPrefix + "Eval";
        String setVar = varPrefix + "Obj";
        String isSetVar = varPrefix + "TypeOk";
        String choiceKeyVar = varPrefix + "Key";
        String choiceKeyKItem = varPrefix + "Elem";
        String choiceVar = varPrefix + "Result";
        String errVar = varPrefix + "Err";

        rhsWriter.writeEvalCalls(mainSb);
        mainSb.writeIndent().append(setChoiceVar).append(" := ");
        rhsWriter.writeReturnValue(mainSb);
        mainSb.newLine();

        mainSb.writeIndent()
                .append(setVar).append(", ").append(isSetVar).append(" := i.Model.")
                .append(expectedKType).append("(")
                .append(setChoiceVar).append(")").newLine();
        mainSb.writeIndent().append("if !").append(isSetVar).beginBlock();
        mainSb.writeIndent().append(reapply).newLine();
        mainSb.endOneBlock();

        mainSb.appendIndentedLine("var ", choiceVar, " m.KReference = m.InternedBottom");
        int forIndent = mainSb.getCurrentIndent();
        mainSb.writeIndent().append("for ").append(choiceKeyVar).append(" := range ").append(setVar).append(".Data").beginBlock();
        mainSb.appendIndentedLine("var ", errVar, " error");
        mainSb.appendIndentedLine(choiceKeyKItem, ", ", errVar, " := i.Model.ToKItem(", choiceKeyVar, ")");
        mainSb.writeIndent().append("if ").append(errVar).append(" != nil").beginBlock();
        mainSb.appendIndentedLine("return m.NoResult, ", errVar);
        mainSb.endOneBlock();

        // this will be after the end of the for, reapply if we didn't hit return in the for loop
        mainSb.addCallbackAfterReturningFromBlock(forIndent, s -> {
            s.newLine();
            s.writeIndent().append("if ").append(choiceVar).append(" == m.InternedBottom").beginBlock();
            s.appendIndentedLine(reapply);
            s.endOneBlock();
            s.appendIndentedLine("return ", choiceVar, ", nil");
        });

        lhsWriter.setNextSubject(choiceKeyKItem);
        lhsWriter.apply(lookup.getLhs());

        // the function goes inside
        // I made it a function just because I didn't feel like changing the RHS returns
        // it can also be done without this function, but then the RHS must assign the result to choice var instead of returning
        int funcIndent = mainSb.getCurrentIndent();
        mainSb.writeIndent().append(choiceVar).append(", ").append(errVar).append(" = ");
        mainSb.append("func() (m.KReference, error)").beginBlock();

        mainSb.addCallbackBeforeReturningFromBlock(funcIndent, s -> {
            s.newLine();
            s.appendIndentedLine("return m.InternedBottom, nil // #setChoice end");
            s.endOneBlock();
        });

        mainSb.addCallbackAfterReturningFromBlock(funcIndent, s -> {
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

    private static void traceLine(GoStringBuilder mainSb, RuleType ruleType, int ruleNum, Rule r) {
        mainSb.appendIndentedLine(
                "i.traceRuleApply(\"",
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

    private static void appendSourceComment(GoStringBuilder mainSb, Rule r) {
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
        mainSb.appendIndentedLine("// source: ", source, " @", startLine);
    }
}
