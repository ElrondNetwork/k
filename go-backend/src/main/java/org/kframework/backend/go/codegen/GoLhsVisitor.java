package org.kframework.backend.go.codegen;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.kore.InjectedKLabel;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KAs;
import org.kframework.kore.KORE;
import org.kframework.kore.KRewrite;
import org.kframework.kore.KSequence;
import org.kframework.kore.KToken;
import org.kframework.kore.KVariable;
import org.kframework.kore.Sort;
import org.kframework.kore.VisitK;
import org.kframework.parser.outer.Outer;
import org.kframework.utils.errorsystem.KEMException;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class GoLhsVisitor extends VisitK {
    private final GoStringBuilder sb;
    private final DefinitionData data;
    private final GoNameProvider nameProvider;
    private final FunctionParams functionVars;
    private final RuleVars lhsVars;
    private final RuleVars rhsVars;

    private int kitemIndex = 0;

    public enum ExpressionType { IF, STATEMENT, NOTHING }

    private ExpressionType topExpressionType = null;

    private void initTopExpressionType(ExpressionType et) {
        if (topExpressionType == null) {
            topExpressionType = et;
        }
    }

    public ExpressionType getTopExpressionType() {
        if (topExpressionType == null) {
            throw new RuntimeException("Generated expression type not initialized.");
        }
        return topExpressionType;
    }

    /**
     * Whenever we see a variable more than once, instead of adding a variable declaration, we add a check that the two instances are equal.
     * This structure keeps track of that.
     */
    private final Set<KVariable> alreadySeenVariables = new HashSet<>();

    public GoLhsVisitor(GoStringBuilder sb,
                        DefinitionData data,
                        GoNameProvider nameProvider,
                        FunctionParams functionVars,
                        RuleVars lhsVars, RuleVars rhsVars) {
        this.sb = sb;
        this.data = data;
        this.nameProvider = nameProvider;
        this.functionVars = functionVars;
        this.lhsVars = lhsVars;
        this.rhsVars = rhsVars;
    }

    private String nextSubject = null;
    private int nextFunctionVarIndex = 0;

    private String consumeSubject() {
        if (nextSubject != null) {
            String subj = nextSubject;
            nextSubject = null;
            return subj;
        }
        if (nextFunctionVarIndex < functionVars.arity()) {
            String subj = functionVars.varName(nextFunctionVarIndex);
            nextFunctionVarIndex++;
            return subj;
        }
        return "BAD_SUBJ";
    }

    public void setNextSubject(String subj) {
        nextSubject = subj;
    }

    private void lhsTypeIf(String castVar, String subject, String type) {
        initTopExpressionType(ExpressionType.IF);
        sb.writeIndent();
        sb.append("if ").append(castVar).append(", t := ");
        sb.append(subject).append(".(").append(type).append("); t");
    }

    @Override
    public void apply(KApply k) {
        if (k.klabel().name().equals("#KToken")) {
            assert k.klist().items().size() == 2;
            KToken ktoken = (KToken) k.klist().items().get(0);
            Sort sort = Outer.parseSort(ktoken.s());
            K value = k.klist().items().get(1);

            //magic down-ness
            String ktVar = "kt" + kitemIndex;
            kitemIndex++;
            lhsTypeIf(ktVar, consumeSubject(), "KToken");
            sb.append(" && ").append(ktVar).append(".Sort == ").append(nameProvider.sortVariableName(sort));
            sb.beginBlock("lhs KApply #KToken");
            nextSubject = ktVar + ".Value";
            apply(value);
        } else if (k.klabel().name().equals("#Bottom")) {
            lhsTypeIf("_", consumeSubject(), "Bottom");
        } else {
            String kappVar = "kapp" + kitemIndex;
            kitemIndex++;
            lhsTypeIf(kappVar, consumeSubject(), "KApply");
            sb.append(" && ").append(kappVar).append(".Label == ").append(nameProvider.klabelVariableName(k.klabel()));
            sb.append(" && len(").append(kappVar).append(".List) == ").append(k.klist().items().size());
            sb.beginBlock("lhs KApply " + k.klabel().name());
            int i = 0;
            for (K item : k.klist().items()) {
                nextSubject = kappVar + ".List[" + i + "]";
                apply(item);
                i++;
            }
        }
    }

    void applyTuple(List<K> items) {
        for (K item : items) {
            apply(item);
        }
    }

    @Override
    public void apply(KAs k) {
        throw KEMException.internalError("GoLhsVisitor.apply(KAs) not implemented.");
    }

    @Override
    public void apply(KRewrite k) {
        throw new AssertionError("unexpected rewrite");
    }

    @Override
    public void apply(KToken k) {
        initTopExpressionType(ExpressionType.IF);
        sb.writeIndent();
        sb.append("if ").append(consumeSubject()).append(" == (");
        GoRhsVisitor.appendKTokenRepresentation(sb, k, data, nameProvider);
        sb.append(")");
        sb.beginBlock("lhs KToken");
    }

    @Override
    public void apply(KVariable k) {
        String varName = lhsVars.getVarName(k);

        if (alreadySeenVariables.contains(k)) {
            initTopExpressionType(ExpressionType.IF);
            sb.writeIndent();
            sb.append("if ").append(consumeSubject()).append(" == ").append(varName);
            sb.beginBlock("lhs KVariable, which reappears:" + k.name());
            alreadySeenVariables.add(k);
            return;
        }

        Sort s = k.att().getOptional(Sort.class).orElse(KORE.Sort(""));
        if (data.mainModule.sortAttributesFor().contains(s)) {
            String hook = data.mainModule.sortAttributesFor().apply(s).<String>getOptional("hook").orElse("");
            if (GoBuiltin.SORT_VAR_HOOKS_1.containsKey(hook)) {
                // these ones don't need to get passed a sort name
                // but if the variable doesn't appear on the RHS, we must make it '_'
                initTopExpressionType(ExpressionType.IF);
                sb.writeIndent();
                String pattern = GoBuiltin.SORT_VAR_HOOKS_1.get(hook);
                String declarationVarName = rhsVars.containsVar(k) ? varName : "_";
                sb.append(String.format(pattern,
                        declarationVarName, consumeSubject()));
                sb.beginBlock("lhs KVariable with hook:" + hook);
                return;
            } else if (GoBuiltin.SORT_VAR_HOOKS_2.containsKey(hook)) {
                // these ones need to get passed a sort name
                // since the variable is used in the condition, we never have to make it '_'
                initTopExpressionType(ExpressionType.IF);
                sb.writeIndent();
                String pattern = GoBuiltin.SORT_VAR_HOOKS_2.get(hook);
                sb.append(String.format(pattern,
                        varName, consumeSubject(), nameProvider.sortVariableName(s)));
                sb.beginBlock("lhs KVariable with hook:" + hook);
                return;
            }
        }

        if (varName == null) {
            initTopExpressionType(ExpressionType.NOTHING);
            sb.writeIndent();
            sb.append("// varName=null").newLine();
        } else if (varName.equals("_")) {
            initTopExpressionType(ExpressionType.NOTHING);
            sb.writeIndent();
            sb.append("// "); // no code here, it is redundant
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable _\n");
        } else if (!rhsVars.containsVar(k)) {
            initTopExpressionType(ExpressionType.NOTHING);
            sb.writeIndent();
            sb.append("// "); // no code here, go will complain that the variable is not used, and will refuse to compile
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable not used\n");
        } else {
            initTopExpressionType(ExpressionType.STATEMENT);
            sb.writeIndent();
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable ").append(k.name()).newLine();
        }
    }

    @Override
    public void apply(KSequence k) {
        switch (k.items().size()) {
        case 1:
            // no KSequence, go straight to the only item
            sb.appendIndentedLine("// lhs KSequence size:" + k.items().size());
            apply(k.items().get(0));
            return;
        case 2:
            // split into head :: tail, if subject is KSequence; subject :: emptySequence otherwise
            String kseqHead = "kseq" + kitemIndex + "Head";
            String kseqTail = "kseq" + kitemIndex + "Tail";
            kitemIndex++;
            sb.writeIndent().append("if ok, ");
            sb.append(kseqHead).append(", ");
            sb.append(kseqTail).append(" := trySplitToHeadTail(").append(consumeSubject()).append("); ok");
            sb.beginBlock("lhs KSequence size:" + k.items().size());
            nextSubject = kseqHead;
            apply(k.items().get(0));
            nextSubject = kseqTail;
            apply(k.items().get(1));
            return;
        default:
            // must match KSequence
            String kseqVar = "kseq" + kitemIndex;
            kitemIndex++;
            lhsTypeIf(kseqVar, consumeSubject(), "KSequence");
            int nrHeads = k.items().size() - 1;
            sb.append(" && len(").append(kseqVar).append(".ks) >= ").append(nrHeads);
            sb.beginBlock("lhs KSequence size:" + k.items().size());
            // heads
            for (int i = 0; i < nrHeads; i++) {
                nextSubject = kseqVar + ".ks[" + i + "]";
                apply(k.items().get(i));
            }
            // tail
            nextSubject = "KSequence{ks:" + kseqVar + ".ks[" + nrHeads + ":]}"; // slice with the rest, can be empty
            apply(k.items().get(nrHeads)); // last element
            return;
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        lhsTypeIf("ikl", consumeSubject(), "InjectedKLabel");
        sb.append(" && ikl.Label == ").append(nameProvider.klabelVariableName(k.klabel()));
        sb.beginBlock();
    }

}