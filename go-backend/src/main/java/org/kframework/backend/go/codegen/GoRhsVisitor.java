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

import java.util.ArrayList;
import java.util.List;

public class GoRhsVisitor extends VisitK {
    protected GoStringBuilder currentSb;
    protected final List<String> evalCalls = new ArrayList<>();

    protected final DefinitionData data;
    private final RuleVars lhsVars;
    private final int topLevelIndent;

    private boolean newlineNext = false;
    private int evalVarIndex = 0;

    protected void start() {
        if (newlineNext) {
            currentSb.newLine();
            currentSb.writeIndent();
            newlineNext = false;
        }
    }

    protected void end() {
    }

    public GoRhsVisitor(DefinitionData data, RuleVars lhsVars, int tabsIndent, int returnValSpacesIndent) {
        this.topLevelIndent = tabsIndent;
        this.currentSb = new GoStringBuilder(tabsIndent, returnValSpacesIndent);
        this.data = data;
        this.lhsVars = lhsVars;
    }

    public void writeEvalCalls(GoStringBuilder sb) {
        for (String evalCall : evalCalls) {
            sb.append(evalCall);
        }
    }

    public void writeReturnValue(GoStringBuilder sb) {
        sb.append(currentSb.toString());
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
            currentSb.append("KToken{Sort: ");
            GoStringUtil.appendSortVariableName(currentSb.sb(), sort);
            currentSb.append(", Value:");
            apply(value);
            currentSb.append("}");
        } else if (k.klabel().name().equals("#Bottom")) {
            currentSb.append("Bottom{}");
        } else if (data.functions.contains(k.klabel()) || data.anywhereKLabels.contains(k.klabel())) {
            applyKApplyExecute(k);
        } else {
            applyKApplyAsIs(k);
        }
        end();
    }

    private void applyKApplyAsIs(KApply k) {
        currentSb.append("KApply{Label: ");
        GoStringUtil.appendKlabelVariableName(currentSb.sb(), k.klabel());
        currentSb.append(", List: []K{ // as-is ").append(k.klabel().name());
        currentSb.increaseIndent();
        for (K item : k.klist().items()) {
            newlineNext = true;
            apply(item);
            currentSb.append(",");
        }
        currentSb.decreaseIndent();
        currentSb.newLine();
        currentSb.writeIndent();
        currentSb.append("}}");
    }

    protected void applyKApplyExecute(KApply k) {
        String evalVarName = "eval" + evalVarIndex;
        String errVarName = "err" + evalVarIndex;
        evalVarIndex++;

        // return the eval variable
        currentSb.append(evalVarName);

        // also add an eval call to the eval calls
        GoStringBuilder evalSb = new GoStringBuilder(topLevelIndent, 0);
        GoStringBuilder backupSb = currentSb;
        currentSb = evalSb; // we trick all nodes below to output to the eval call instead of the return by changing the string builder

        evalSb.writeIndent().append(evalVarName).append(", ").append(errVarName).append(" := ");
        GoStringUtil.appendFunctionName(evalSb.sb(), k.klabel()); // func name
        if (k.items().size() == 0) { // call parameters
            evalSb.append("(config)").newLine();
        } else {
            evalSb.append("(");
            evalSb.increaseIndent();
            for(K item : k.items()) {
                newlineNext = true;
                apply(item);
                evalSb.append(",");
            }
            evalSb.newLine().writeIndent().append("config)");
            evalSb.decreaseIndent();
            evalSb.newLine();
        }
        evalSb.writeIndent().append("if ").append(errVarName).append(" != nil").beginBlock();
        evalSb.writeIndent().append("return noResult, ").append(errVarName).newLine();
        evalSb.endOneBlock();

        evalCalls.add(evalSb.toString());
        assert currentSb == evalSb;
        currentSb = backupSb; // restore
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
        appendKTokenRepresentation(currentSb, k, data);
        end();
    }

    protected void appendKTokenComment(KToken k) {
        if (k.sort().equals(Sorts.Bool()) && k.att().contains(PrecomputePredicates.COMMENT_KEY)) {
            currentSb.append("/* rhs precomputed ").append(k.att().get(PrecomputePredicates.COMMENT_KEY)).append(" */ ");
        } else{
            currentSb.append("/* rhs KToken */ ");
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

        sb.append("KToken{Sort: ");
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
            currentSb.append("/* varName=null */ internedBottom");
            end();
            return;
        }

        if (!lhsVars.containsVar(v) && varName.startsWith("?")) {
            throw KEMException.internalError("Failed to compile rule due to unmatched variable on right-hand-side. This is likely due to an unsupported collection pattern: " + varName, v);
        } else if (!lhsVars.containsVar(v)) {
            currentSb.append("panic(\"Stuck!\")");
        } else {
            KLabel listVar = lhsVars.listVars.get(varName);
            if (listVar != null) {
                currentSb.append("List{Sort: ");
                GoStringUtil.appendSortVariableName(currentSb.sb(), data.mainModule.sortFor().apply(listVar));
                currentSb.append(", Label:");
                GoStringUtil.appendKlabelVariableName(currentSb.sb(), listVar);
                //currentSb.append(", ");
                //currentSb.append(varOccurrance);
                currentSb.append(" /* ??? */}");
            } else {
                currentSb.append(varName);
            }
        }
        end();
    }

    @Override
    public void apply(KSequence k) {
        int size = k.items().size();
        switch (k.items().size()) {
        case 1:
            currentSb.append("/* rhs KSequence size=1 */ ");
            apply(k.items().get(0));
            return;
        case 2:
            start();
            currentSb.append("assembleFromHeadTail(");
            currentSb.increaseIndent();
            for (K item : k.items()) {
                newlineNext = true;
                apply(item);
                currentSb.append(",");
            }
            currentSb.decreaseIndent();
            currentSb.newLine().writeIndent().append(")");
            return;
        default:
            start();
            currentSb.append("KSequence { ks: append([]K{ ");
            currentSb.increaseIndent();
            // heads
            for (int i = 0; i < k.items().size() - 1; i++) {
                newlineNext = true;
                apply(k.items().get(i));
                currentSb.append(",");
            }
            // tail
            currentSb.decreaseIndent();
            currentSb.newLine().writeIndent().append("}, ");
            apply(k.items().get(k.items().size() - 1));
            currentSb.append(".ks...)}");
            end();
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        start();
        currentSb.append("InjectedKLabel{Sort: ");
        GoStringUtil.appendKlabelVariableName(currentSb.sb(), k.klabel());
        currentSb.append("}");
        end();
    }

}