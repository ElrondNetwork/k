package org.kframework.backend.go.codegen;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.builtin.Sorts;
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

public class GoRhsVisitor extends VisitK {
    protected final GoStringBuilder sb;
    protected final DefinitionData data;
    private final RuleVars lhsVars;

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

    public GoRhsVisitor(GoStringBuilder sb, DefinitionData data, RuleVars lhsVars) {
        this.sb = sb;
        this.data = data;
        this.lhsVars = lhsVars;
    }

    @Override
    public void apply(KApply k) {
        start();
        if (k.klabel().name().equals("#KToken")) {
            assert k.klist().items().size() == 2;
            KToken ktoken = (KToken) k.klist().items().get(0);
            Sort sort = Outer.parseSort(ktoken.s());
            K value = k.klist().items().get(1);

            //magic down-ness
            sb.append("KToken{Sort: ");
            GoStringUtil.appendSortVariableName(sb.sb(), sort);
            sb.append(", Value:");
            apply(value);
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
        sb.append("KApply{Label: ");
        GoStringUtil.appendKlabelVariableName(sb.sb(), k.klabel());
        sb.append(", List: []K{ // as-is ").append(k.klabel().name());
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
        sb.append("/* execute: */ "); // comment
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
        throw KEMException.internalError("GoRhsVisitor.apply(KAs) not implemented.");
    }

    @Override
    public void apply(KRewrite k) {
        throw new AssertionError("unexpected rewrite");
    }

    @Override
    public void apply(KToken k) {
        start();
        appendKTokenComment(k);
        appendKTokenRepresentation(sb, k, data);
        end();
    }

    protected void appendKTokenComment(KToken k) {
        if (k.sort().equals(Sorts.Bool()) && k.att().contains(PrecomputePredicates.COMMENT_KEY)) {
            sb.append("/* rhs precomputed ").append(k.att().get(PrecomputePredicates.COMMENT_KEY)).append(" */ ");
        } else{
            sb.append("/* rhs KToken */ ");
        }
    }

    /**
     * This one is also used by the GoLhsVisitor.
     * */
    public static void appendKTokenRepresentation(GoStringBuilder sb, KToken k, DefinitionData data) {
        if (data.mainModule.sortAttributesFor().contains(k.sort())) {
            String hook = data.mainModule.sortAttributesFor().apply(k.sort()).<String>getOptional("hook").orElse("");
            if (GoBuiltin.GO_SORT_TOKEN_HOOKS.containsKey(hook)) {
                sb.append(GoBuiltin.GO_SORT_TOKEN_HOOKS.get(hook).apply(k.s()));
                return;
            }
        }

        sb.append(" KToken{Sort: ");
        GoStringUtil.appendSortVariableName(sb.sb(), k.sort());
        sb.append(", Value: ");
        sb.append(GoStringUtil.enquoteString(k.s()));
        sb.append("}");
    }

    @Override
    public void apply(KVariable v) {
        start();
        String varName = lhsVars.getVarName(v);
        if (varName == null) {
            sb.append("/* varName=null */ internedBottom");
            end();
            return;
        }

        if (!lhsVars.containsVar(v) && varName.startsWith("?")) {
            throw KEMException.internalError("Failed to compile rule due to unmatched variable on right-hand-side. This is likely due to an unsupported collection pattern: " + varName, v);
        } else if (!lhsVars.containsVar(v)) {
            sb.append("panic(\"Stuck!\")");
        } else {
            KLabel listVar = lhsVars.listVars.get(varName);
            if (listVar != null) {
                sb.append("List{Sort: ");
                GoStringUtil.appendSortVariableName(sb.sb(), data.mainModule.sortFor().apply(listVar));
                sb.append(", Label:");
                GoStringUtil.appendKlabelVariableName(sb.sb(), listVar);
                //sb.append(", ");
                //sb.append(varOccurrance);
                sb.append(" /* ??? */}");
            } else {
                sb.append(varName);
            }
        }
        end();
    }

    @Override
    public void apply(KSequence k) {
        int size = k.items().size();
        switch (k.items().size()) {
        case 1:
            sb.append("/* rhs KSequence size=1 */ ");
            apply(k.items().get(0));
            return;
        case 2:
            start();
            sb.append("assembleFromHeadTail(");
            sb.increaseIndent();
            for (K item : k.items()) {
                newlineNext = true;
                apply(item);
                sb.append(",");
            }
            sb.decreaseIndent();
            sb.newLine().writeIndent().append(")");
            return;
        default:
            start();
            sb.append("KSequence { ks: append([]K{ ");
            sb.increaseIndent();
            // heads
            for (int i = 0; i < k.items().size() - 1; i++) {
                newlineNext = true;
                apply(k.items().get(i));
                sb.append(",");
            }
            // tail
            sb.decreaseIndent();
            sb.newLine().writeIndent().append("}, ");
            apply(k.items().get(k.items().size() - 1));
            sb.append(".ks...)}");
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