package org.kframework.backend.go.codegen.rules;

import org.kframework.backend.go.codegen.GoBuiltin;
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
import org.kframework.unparser.ToKast;

import java.util.List;
import java.util.Set;

public class RuleLhsWriter extends VisitK {
    private final GoStringBuilder sb;
    private final DefinitionData data;
    private final GoNameProvider nameProvider;
    private final FunctionParams functionVars;
    private final RuleVars lhsVars;
    private final RuleVars rhsVars;

    /**
     * Whenever we see a variable more than once, instead of adding a variable declaration, we add a check that the two instances are equal.
     * This structure keeps track of that.
     */
    private final Set<KVariable> alreadySeenVariables;
    private final boolean startWithScopeBlockIfNecessary;

    private int kitemIndex = 0;

    public enum ExpressionType { IF, STATEMENT, NOTHING }

    private ExpressionType topExpressionType = null;
    private boolean containsIf = false;

    private void handleExpressionType(ExpressionType et) {
        if (topExpressionType == null) {
            topExpressionType = et;
            if (startWithScopeBlockIfNecessary && et != ExpressionType.IF) {
                sb.scopingBlock("scoping block, to avoid variable name collisions");
            }
        }
        if (et == ExpressionType.IF) {
            containsIf = true;
        }
    }

    public ExpressionType getTopExpressionType() {
        if (topExpressionType == null) {
            topExpressionType = ExpressionType.NOTHING;
        }
        return topExpressionType;
    }

    public boolean containsIf() {
        return containsIf;
    }

    public RuleLhsWriter(GoStringBuilder sb,
                         DefinitionData data,
                         GoNameProvider nameProvider,
                         FunctionParams functionVars,
                         RuleVars lhsVars, RuleVars rhsVars,
                         Set<KVariable> alreadySeenVariables,
                         boolean startWithScopeBlockIfNecessary) {
        this.sb = sb;
        this.data = data;
        this.nameProvider = nameProvider;
        this.functionVars = functionVars;
        this.lhsVars = lhsVars;
        this.rhsVars = rhsVars;
        this.alreadySeenVariables = alreadySeenVariables;
        this.startWithScopeBlockIfNecessary = startWithScopeBlockIfNecessary;
    }

    private String nextSubject = null;
    private int nextFunctionVarIndex = 0;
    private KVariable nextAlias = null;

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

    private KVariable consumeAlias() {
        if (nextAlias != null) {
            KVariable alias = nextAlias;
            nextAlias = null;
            return alias;
        }
        return null;
    }

    private void lhsTypeIf(String castVar, String subject, String type) {
        handleExpressionType(ExpressionType.IF);
        sb.writeIndent();
        sb.append("if ").append(castVar).append(", t := ");
        sb.append(subject).append(".(*m.").append(type).append("); t");
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
            sb.append(" && ").append(ktVar).append(".Sort == m.").append(nameProvider.sortVariableName(sort));
            sb.beginBlock("lhs KApply #KToken");
            nextSubject = ktVar + ".Value";
            apply(value);
        } else if (k.klabel().name().equals("#Bottom")) {
            lhsTypeIf("_", consumeSubject(), "Bottom");
        } else {
            KVariable alias  = consumeAlias();
            String kappVar;
            String aliasComment = "";
            if (alias != null) {
                kappVar = lhsVars.getVarName(alias);
                aliasComment = " as " + alias.name();
            } else {
                kappVar = "kapp" + kitemIndex;
                kitemIndex++;
            }

            lhsTypeIf(kappVar, consumeSubject(), "KApply");
            sb.append(" && ").append(kappVar).append(".Label == m.").append(nameProvider.klabelVariableName(k.klabel()));
            sb.append(" && len(").append(kappVar).append(".List) == ").append(k.klist().items().size());
            sb.beginBlock(ToKast.apply(k), aliasComment);
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

        if (!(k.alias() instanceof KVariable)) {
            throw new IllegalArgumentException("KAs alias is not a KVariable.");
        }
        nextAlias = (KVariable)k.alias();
        apply(k.pattern());

        if (nextAlias != null) {
            throw new RuntimeException("KAs alias was not consumed. This scenario was not handled. An alias will be missing ");
        }
    }

    @Override
    public void apply(KRewrite k) {
        throw new AssertionError("unexpected rewrite");
    }

    @Override
    public void apply(KToken k) {
        handleExpressionType(ExpressionType.IF);
        sb.writeIndent();
        sb.append("if ").append(consumeSubject()).append(".Equals(");
        RuleRhsWriter.appendKTokenRepresentation(sb, k, data, nameProvider);
        sb.append(")");
        sb.beginBlock(ToKast.apply(k));
    }

    @Override
    public void apply(KVariable k) {
        String varName = lhsVars.getVarName(k);

        if (alreadySeenVariables.contains(k)) {
            handleExpressionType(ExpressionType.IF);
            sb.writeIndent();
            sb.append("if ").append(consumeSubject()).append(".Equals(").append(varName).append(")");
            sb.beginBlock("lhs KVariable, which reappears:" + k.name());
            return;
        }
        alreadySeenVariables.add(k);

        Sort s = k.att().getOptional(Sort.class).orElse(KORE.Sort(""));
        if (data.mainModule.sortAttributesFor().contains(s)) {
            String hook = data.mainModule.sortAttributesFor().apply(s).<String>getOptional("hook").orElse("");
            if (GoBuiltin.SORT_VAR_HOOKS_1.containsKey(hook)) {
                // these ones don't need to get passed a sort name
                // but if the variable doesn't appear on the RHS, we must make it '_'
                handleExpressionType(ExpressionType.IF);
                sb.writeIndent();
                String pattern = GoBuiltin.SORT_VAR_HOOKS_1.get(hook);
                boolean varNeeded = rhsVars.containsVar(k) // needed in RHS
                                || lhsVars.getVarCount(k) > 1; // needed in LHS, when it reappears
                String declarationVarName = varNeeded ? varName : "_";
                sb.append(String.format(pattern,
                        declarationVarName, consumeSubject()));
                sb.beginBlock("lhs KVariable with hook:" + hook);
                return;
            } else if (GoBuiltin.SORT_VAR_HOOKS_2.containsKey(hook)) {
                // these ones need to get passed a sort name
                // since the variable is used in the condition, we never have to make it '_'
                handleExpressionType(ExpressionType.IF);
                sb.writeIndent();
                String pattern = GoBuiltin.SORT_VAR_HOOKS_2.get(hook);
                sb.append(String.format(pattern,
                        varName, consumeSubject(), nameProvider.sortVariableName(s)));
                sb.beginBlock("lhs KVariable with hook:" + hook);
                return;
            }
        }

        if (varName == null) {
            handleExpressionType(ExpressionType.NOTHING);
            sb.writeIndent();
            sb.append("// varName=null").newLine();
        } else if (varName.equals("_")) {
            handleExpressionType(ExpressionType.NOTHING);
            sb.writeIndent();
            sb.append("// "); // no code here, it is redundant
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable _\n");
        } else if (!rhsVars.containsVar(k)) {
            handleExpressionType(ExpressionType.NOTHING);
            String subject = consumeSubject();
            sb.writeIndent();
            sb.append("doNothing(").append(subject).append(") ");
            sb.append("// "); // no code here, go will complain that the variable is not used, and will refuse to compile
            sb.append(varName).append(" := ").append(subject).append(" // lhs KVariable not used").newLine();
        } else {
            handleExpressionType(ExpressionType.STATEMENT);
            sb.writeIndent();
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable ").append(k.name()).newLine();
        }
    }

    @Override
    public void apply(KSequence k) {
        switch (k.items().size()) {
        case 0:
            sb.appendIndentedLine("// KSequence, size 0:", ToKast.apply(k));
            return;
        case 1:
            // no KSequence, go straight to the only item
            sb.appendIndentedLine("// KSequence, size 1:", ToKast.apply(k));
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
            sb.beginBlock(ToKast.apply(k));
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
            sb.append(" && len(").append(kseqVar).append(".Ks) >= ").append(nrHeads);
            sb.beginBlock("lhs KSequence size:" + k.items().size());
            // heads
            for (int i = 0; i < nrHeads; i++) {
                nextSubject = kseqVar + ".Ks[" + i + "]";
                apply(k.items().get(i));
            }
            // tail
            nextSubject = "&m.KSequence{Ks:" + kseqVar + ".Ks[" + nrHeads + ":]}"; // slice with the rest, can be empty
            apply(k.items().get(nrHeads)); // last element
            return;
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        lhsTypeIf("ikl", consumeSubject(), "InjectedKLabel");
        sb.append(" && ikl.Label == m.").append(nameProvider.klabelVariableName(k.klabel()));
        sb.beginBlock();
    }

}