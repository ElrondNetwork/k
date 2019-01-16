package org.kframework.backend.go;

import org.kframework.kore.InjectedKLabel;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KAs;
import org.kframework.kore.KLabel;
import org.kframework.kore.KRewrite;
import org.kframework.kore.KSequence;
import org.kframework.kore.KToken;
import org.kframework.kore.KVariable;
import org.kframework.kore.Sort;
import org.kframework.kore.VisitK;
import org.kframework.parser.outer.Outer;
import org.kframework.utils.errorsystem.KEMException;

class GoRhsVisitor extends VisitK {
    protected final GoStringBuilder sb;
    protected final VarInfo vars;
    protected final DefinitionData data;

    private boolean newlineNext = false;

    protected void start() {
        if (newlineNext) {
            sb.newLine();
            sb.writeIndent();
            newlineNext = false;
        }
    }

    protected void end() {
    }

    public GoRhsVisitor(GoStringBuilder sb, VarInfo vars, DefinitionData data) {
        this.sb = sb;
        this.vars = vars;
        this.data = data;
    }

    @Override
    public void apply(KApply k) {
        start();
        if (k.klabel().name().equals("#KToken")) {
            //magic down-ness
            sb.append("KToken{Sort: ");
            Sort sort = Outer.parseSort(((KToken) ((KSequence) k.klist().items().get(0)).items().get(0)).s());
            GoStringUtil.appendSortVariableName(sb.sb(), sort);
            sb.append(", Value:");
            apply(((KSequence) k.klist().items().get(1)).items().get(0));
            sb.append("}");
        } else if (k.klabel().name().equals("#Bottom")) {
            sb.append("Bottom{}");
        } else if (data.functions.contains(k.klabel()) || data.anywhereKLabels.contains(k.klabel())) {
            applyKApplyExecute(k);
        } else {
            applyKApplyAsIs(k);
        }
        end();
    }

    private void applyKApplyAsIs(KApply k) {
        sb.append("/* as-is */ KApply{Label: ");
        GoStringUtil.appendKlabelVariableName(sb.sb(), k.klabel());
        sb.append(", List: []K{");
        sb.increaseIndent();
        for (K item : k.klist().items()) {
            newlineNext = true;
            apply(item);
            sb.append(",");
        }
        sb.decreaseIndent();
        sb.newLine();
        sb.writeIndent();
        sb.append("}}");
    }

    protected void applyKApplyExecute(KApply k) {
        sb.append("/* execute: */"); // comment
        GoStringUtil.appendFunctionName(sb.sb(), k.klabel()); // func name
        if (k.items().size() == 0) { // call parameters
            sb.append("(config)");
        } else {
            sb.append("(");
            sb.increaseIndent();
            for(K item : k.items()) {
                newlineNext = true;
                apply(item);
                sb.append(",");
            }
            sb.newLine().writeIndent().append("config)");
            sb.decreaseIndent();
        }
    }

    @Override
    public void apply(KAs k) {
        start();
        sb.append("/* rhs KAs */");
        end();
    }

    @Override
    public void apply(KRewrite k) {
        throw new AssertionError("unexpected rewrite");
    }

    @Override
    public void apply(KToken k) {
        start();
        sb.append("/* KToken */");
        appendKTokenRepresentation(sb.sb(), k, data);
        end();
    }

    public static void appendKTokenRepresentation(StringBuilder sb, KToken k, DefinitionData data) {
        if (data.mainModule.sortAttributesFor().contains(k.sort())) {
            String hook = data.mainModule.sortAttributesFor().apply(k.sort()).<String>getOptional("hook").orElse("");
            if (GoBuiltin.GO_SORT_TOKEN_HOOKS.containsKey(hook)) {
                sb.append(GoBuiltin.GO_SORT_TOKEN_HOOKS.get(hook).apply(k.s()));
                return;
            }
        }

        sb.append(" KToken{Sort: ");
        GoStringUtil.appendSortVariableName(sb, k.sort());
        sb.append(", Value: ");
        sb.append(GoStringUtil.enquoteString(k.s()));
        sb.append("}");
    }

    @Override
    public void apply(KVariable v) {
        start();
        String varName = GoStringUtil.variableName(v.name());

        if (vars.vars.get(v).isEmpty() && varName.startsWith("?")) {
            throw KEMException.internalError("Failed to compile rule due to unmatched variable on right-hand-side. This is likely due to an unsupported collection pattern: " + varName, v);
        } else if (vars.vars.get(v).isEmpty()) {
            sb.append("panic(\"Stuck!\")");
        } else {
            String varOccurrence = vars.vars.get(v).iterator().next();
            KLabel listVar = vars.listVars.get(vars.vars.get(v).iterator().next());
            if (listVar != null) {
                sb.append("List{Sort: ");
                GoStringUtil.appendSortVariableName(sb.sb(), data.mainModule.sortFor().apply(listVar));
                sb.append(", Label:");
                GoStringUtil.appendKlabelVariableName(sb.sb(), listVar);
                //sb.append(", ");
                //sb.append(varOccurrance);
                sb.append(" /* ??? */}");
            } else {
                sb.append(varOccurrence);
            }
        }
        end();
    }

    @Override
    public void apply(KSequence k) {
        int size = k.items().size();
        sb.append("/* KSequence size=").append(size).append(" */ ");
        if (size == 0) {
        } else if (size == 1) {
            apply(k.items().get(0));
        } else {
            start();
            boolean first = true;
            for (K item : k.items()) {
                if (first) {
                    first = false;
                } else {
                    sb.append(", ");
                }
                apply(item);
            }
            end();
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        start();
        sb.append("InjectedKLabel{Sort: ");
        GoStringUtil.appendKlabelVariableName(sb.sb(), k.klabel());
        sb.append("}");
        end();
    }

}