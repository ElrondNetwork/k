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
        top.populateDFOrderIndex(0);
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
        if (leafNode.rule.att().contains("owise")) {
            sb.appendIndentedLine("// [owise]");
        }
        if (leafNode.rule.att().contains("structural")) {
            sb.appendIndentedLine("// [structural]");
        }

        // output lookups
        writeLookups(sb,
                leafNode.getDFOrderIndex(),
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

    private void writeLookups(GoStringBuilder mainSb,
                              int guardIndex,
                              FunctionInfo functionInfo,
                              List<Lookup> lookups,
                              TempVarManager varManager) {
        if (lookups.isEmpty()) {
            return;
        }
        mainSb.appendIndentedLine("// LOOKUPS");
        mainSb.writeIndent().append("if guard < ").append(guardIndex).beginBlock();

        int lookupIndex = 0;
        for (Lookup lookup : lookups) {
            String reapply = "return i." + functionInfo.goName + "(" + functionInfo.arguments.callParameters() + "config, " + guardIndex + ") // reapply";

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
                        "setChoice" + lookupIndex,
                        "IsSet", "SetChoiceLookup",
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
                        "mapChoice" + lookupIndex,
                        "IsMap", "MapKeyChoiceLookup",
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
            String varPrefix, String predicateFunc, String iteratorCall,
            TempVarManager varManager,
            String reapply,
            RuleLhsWriter lhsWriter, RuleRhsWriterBase rhsWriter) {

        String choiceSubject = varPrefix + "Subj";
        String choiceKeyVar = varPrefix + "Key";
        String choiceResult = varPrefix + "Result";
        String errVar = varPrefix + "Err";

        rhsWriter.writeEvalCalls(mainSb);
        mainSb.writeIndent().append(choiceSubject).append(" := ");
        rhsWriter.writeReturnValue(mainSb);
        mainSb.newLine();

        mainSb.writeIndent().append("if i.Model.").append(predicateFunc).append("(").append(choiceSubject).append(")").beginBlock();

        int foreachIndent = mainSb.getCurrentIndent();
        mainSb.writeIndent().append(choiceResult).append(", ").append(errVar).append(" := ");
        mainSb.append("i.Model.").append(iteratorCall);
        mainSb.append("(").append(choiceSubject).append(", func (").append(choiceKeyVar).append(" KReference) (KReference, error)").beginBlock();

        mainSb.addCallbackBeforeReturningFromBlock(foreachIndent, s -> {
            s.newLine();
            s.appendIndentedLine("return m.InternedBottom, nil // #choice end");
            s.endOneBlock();
        });
        mainSb.addCallbackAfterReturningFromBlock(foreachIndent, s -> {
            s.append(") // choice foreach end");
            s.newLine();
            s.writeIndent().append("if ").append(errVar).append(" != nil").beginBlock();
            s.appendIndentedLine("return m.NoResult, ", errVar);
            s.endOneBlock();
            s.writeIndent().append("if ").append(choiceResult).append(" != m.InternedBottom").beginBlock();
            s.appendIndentedLine("return ", choiceResult, ", nil");
            s.endOneBlock();
        });

        lhsWriter.setNextSubject(choiceKeyVar);
        lhsWriter.apply(lookup.getLhs());
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
