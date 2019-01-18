package org.kframework.backend.go;

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

import java.util.List;
import java.util.Set;

class GoLhsVisitor extends VisitK {
    private final GoStringBuilder sb;
    private final VarInfo vars;
    private final DefinitionData data;

    private final FunctionVars functionVars;
    private final Set<String> usedVariableNames;

    public GoLhsVisitor(GoStringBuilder sb, VarInfo vars, DefinitionData data, FunctionVars functionVars, Set<String> usedVariableNames) {
        this.sb = sb;
        this.vars = vars;
        this.data = data;
        this.functionVars = functionVars;
        this.usedVariableNames = usedVariableNames;
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

    private void lhsTypeIf(String castVar, String subject, String type) {
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
            lhsTypeIf("kt", consumeSubject(), "KToken");
            sb.append(" && kt.Sort == ");
            GoStringUtil.appendSortVariableName(sb.sb(), sort);
            sb.beginBlock("lhs KApply #KToken");
            nextSubject = "kt.Value";
            apply(value);
        } else if (k.klabel().name().equals("#Bottom")) {
            lhsTypeIf("_", consumeSubject(), "Bottom");
        } else {
            lhsTypeIf("kapp", consumeSubject(), "KApply");
            sb.append(" && kapp.Label == ");
            GoStringUtil.appendKlabelVariableName(sb.sb(), k.klabel());
            sb.beginBlock("lhs KApply");
            int i = 0;
            for (K item : k.klist().items()) {
                nextSubject = "kapp.List[" + i + "]";
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
        sb.writeIndent();
        sb.writeIndent();
        sb.append("if ").append(consumeSubject()).append(" == ");
        GoRhsVisitor.appendKTokenRepresentation(sb.sb(), k, data);
        sb.beginBlock("lhs KToken");
    }

    @Override
    public void apply(KVariable k) {
        String varName = GoStringUtil.variableName(k.name());

        if (vars.vars.containsKey(k)) {
            sb.writeIndent();
            sb.append("if ").append(consumeSubject()).append(" == ").append(varName);
            sb.beginBlock("apply KVariable, which reappears:" + k.name());
            vars.vars.put(k, varName);
            return;
        }

        vars.vars.put(k, varName);
        Sort s = k.att().getOptional(Sort.class).orElse(KORE.Sort(""));
        if (data.mainModule.sortAttributesFor().contains(s)) {
            String hook = data.mainModule.sortAttributesFor().apply(s).<String>getOptional("hook").orElse("");
            if (GoBuiltin.GO_SORT_VAR_HOOKS.containsKey(hook)) {
                sb.writeIndent();
                String pattern = GoBuiltin.GO_SORT_VAR_HOOKS.get(hook);
                sb.append(String.format(pattern,
                        varName, consumeSubject(), GoStringUtil.sortVariableName(s)));
                sb.beginBlock("apply KVariable with hook:" + hook);
                return;
            }
        }

        if (varName.equals("_")) {
            sb.writeIndent();
            sb.append("// "); // no code here, it is redundant
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable _\n");
        } else if (!usedVariableNames.contains(k.name())) {
            sb.writeIndent();
            sb.append("// "); // no code here, go will complain that the variable is not used, and will refuse to compile
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable not used\n");
        } else {
            sb.writeIndent();
            sb.append(varName).append(" := ").append(consumeSubject()).append(" // lhs KVariable ok\n");
        }
    }

    @Override
    public void apply(KSequence k) {
        sb.writeIndent();
        sb.append("// apply KSequence size:" + k.items().size() + "\n");
        if (k.items().size() == 1) {
            apply(k.items().get(0));
            return;
        }

        int i = 1;
        for (K item : k.items()) {
            apply(item);
            i++;
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        lhsTypeIf("ikl", consumeSubject(), "InjectedKLabel");
        sb.append(" && ikl.Sort == ");
        GoStringUtil.appendKlabelVariableName(sb.sb(), k.klabel());
        sb.beginBlock();
    }

}