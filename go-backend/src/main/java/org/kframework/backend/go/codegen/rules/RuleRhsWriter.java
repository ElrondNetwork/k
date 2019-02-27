package org.kframework.backend.go.codegen.rules;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.model.TempVarCounters;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoNameProvider;
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
import org.kframework.unparser.ToKast;
import org.kframework.utils.StringUtil;
import org.kframework.utils.errorsystem.KEMException;

import java.util.ArrayList;
import java.util.List;

public class RuleRhsWriter extends VisitK {
    protected GoStringBuilder currentSb;
    protected final List<String> evalCalls = new ArrayList<>();

    protected final DefinitionData data;
    protected final GoNameProvider nameProvider;
    private final RuleVars lhsVars;
    private final TempVarCounters tempVarCounters;
    private final int topLevelIndent;

    private boolean newlineNext = false;

    protected void start() {
        if (newlineNext) {
            currentSb.newLine();
            currentSb.writeIndent();
            newlineNext = false;
        }
    }

    protected void end() {
    }

    public RuleRhsWriter(DefinitionData data,
                         GoNameProvider nameProvider,
                         RuleVars lhsVars,
                         TempVarCounters tempVarCounters,
                         int tabsIndent, int returnValSpacesIndent) {
        this.topLevelIndent = tabsIndent;
        this.currentSb = new GoStringBuilder(tabsIndent, returnValSpacesIndent);
        this.data = data;
        this.nameProvider = nameProvider;
        this.lhsVars = lhsVars;
        this.tempVarCounters = tempVarCounters;
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
            currentSb.append("&m.KToken{Sort: m.").append(nameProvider.sortVariableName(sort));
            currentSb.append(", Value:");
            apply(value);
            currentSb.append("}");
        } else if (k.klabel().name().equals("#Bottom")) {
            currentSb.append("&m.Bottom{}");
        } else if (data.functions.contains(k.klabel()) || data.anywhereKLabels.contains(k.klabel())) {
            applyKApplyExecute(k);
        } else {
            applyKApplyAsIs(k);
        }
        end();
    }

    private void applyKApplyAsIs(KApply k) {
        currentSb.append("&m.KApply{Label: m.").append(nameProvider.klabelVariableName(k.klabel()));
        currentSb.append(", List: []m.K{ // as-is ").append(k.klabel().name());
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
        int evalVarIndex = tempVarCounters.consumeEvalVarIndex();
        String evalVarName = "eval" + evalVarIndex;
        String errVarName = "err" + evalVarIndex;

        // return the eval variable
        currentSb.append(evalVarName);

        // also add an eval call to the eval calls
        GoStringBuilder evalSb = new GoStringBuilder(topLevelIndent, 0);
        GoStringBuilder backupSb = currentSb;
        currentSb = evalSb; // we trick all nodes below to output to the eval call instead of the return by changing the string builder

        String comment = ToKast.apply(k);

        evalSb.writeIndent().append(evalVarName).append(", ").append(errVarName).append(" := ");
        evalSb.append(nameProvider.evalFunctionName(k.klabel())); // func name
        if (k.items().size() == 0) { // call parameters
            evalSb.append("(config, -1) // ").append(comment).newLine();
        } else {
            evalSb.append("( // ").append(comment);
            evalSb.increaseIndent();
            for (K item : k.items()) {
                newlineNext = true;
                apply(item);
                evalSb.append(",");
            }
            evalSb.newLine().writeIndent().append("config, -1)");
            evalSb.decreaseIndent();
            evalSb.newLine();
        }
        evalSb.writeIndent().append("if ").append(errVarName).append(" != nil").beginBlock();
        evalSb.writeIndent().append("return m.NoResult, ").append(errVarName).newLine();
        evalSb.endOneBlock();

        evalCalls.add(evalSb.toString());
        assert currentSb == evalSb;
        currentSb = backupSb; // restore
    }

    @Override
    public void apply(KAs k) {
        throw KEMException.internalError("RuleRhsWriter.apply(KAs) not implemented.");
    }

    @Override
    public void apply(KRewrite k) {
        throw new AssertionError("unexpected rewrite");
    }

    @Override
    public void apply(KToken k) {
        start();
        appendKTokenComment(k);
        appendKTokenRepresentation(currentSb, k, data, nameProvider);
        end();
    }

    protected void appendKTokenComment(KToken k) {
        if (k.sort().equals(Sorts.Bool()) && k.att().contains(PrecomputePredicates.COMMENT_KEY)) {
            currentSb.append("/* rhs precomputed ").append(k.att().get(PrecomputePredicates.COMMENT_KEY)).append(" */ ");
        } else {
            currentSb.append("/* rhs KToken */ ");
        }
    }

    /**
     * This one is also used by the RuleLhsWriter.
     */
    public static void appendKTokenRepresentation(GoStringBuilder sb, KToken k, DefinitionData data, GoNameProvider nameProvider) {
        if (data.mainModule.sortAttributesFor().contains(k.sort())) {
            String hook = data.mainModule.sortAttributesFor().apply(k.sort()).<String>getOptional("hook").orElse("");
            switch (hook) {
            case "BOOL.Bool":
                if (k.s().equals("true")) {
                    sb.append("m.BoolTrue");
                } else if (k.s().equals("false")) {
                    sb.append("m.BoolFalse");
                } else {
                    throw new RuntimeException("Unexpected Bool token value: " + k.s());
                }
                return;
            case "MINT.MInt":
                sb.append("&m.MInt{Value: ").append(k.s()).append("}");
                return;
            case "INT.Int":
                sb.append("m.NewIntFromString(\"").append(k.s()).append("\")");
                return;
            case "FLOAT.Float":
                sb.append("&m.Float{Value: ").append(k.s()).append("}");
                return;
            case "STRING.String":
                String unquotedStr = StringUtil.unquoteKString(k.s());
                String goStr = GoStringUtil.enquoteString(unquotedStr);
                sb.append("m.NewString(").append(goStr).append(")");
                return;
            case "BYTES.Bytes":
                String unquotedBytes = StringUtil.unquoteKString(k.s());
                String goBytes = GoStringUtil.enquoteString(unquotedBytes);
                sb.append("&m.Bytes{Value: ").append(goBytes).append("}");
                return;
            case "BUFFER.StringBuffer":
                sb.append("&m.StringBuffer{}");
                return;
            }
        }

        sb.append("&m.KToken{Sort: m.").append(nameProvider.sortVariableName(k.sort()));
        sb.append(", Value: ");
        sb.append(GoStringUtil.enquoteString(k.s()));
        sb.append("}");
    }

    @Override
    public void apply(KVariable v) {
        start();
        String varName = lhsVars.getVarName(v);
        if (varName == null) {
            currentSb.append("/* varName=null */ m.InternedBottom");
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
                Sort sort = data.mainModule.sortFor().apply(listVar);
                currentSb.append("&m.List{Sort: m.").append(nameProvider.sortVariableName(sort));
                currentSb.append(", Label: m.").append(nameProvider.klabelVariableName(listVar));
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
        case 0:
            start();
            currentSb.append("m.EmptyKSequence");
            end();
            return;
        case 1:
            currentSb.append("/* rhs KSequence size=1 */ ");
            apply(k.items().get(0));
            return;
        case 2:
            start();
            currentSb.append("assembleFromHeadAndTail(");
            currentSb.increaseIndent();
            for (K item : k.items()) {
                newlineNext = true;
                apply(item);
                currentSb.append(",");
            }
            currentSb.decreaseIndent();
            currentSb.newLine().writeIndent().append(")");
            end();
            return;
        default:
            start();
            currentSb.append("assembleFromHeadSliceAndTail([]m.K{");
            currentSb.increaseIndent();
            // head slice
            for (int i = 0; i < k.items().size() - 1; i++) {
                newlineNext = true;
                apply(k.items().get(i));
                currentSb.append(",");
            }
            // tail
            currentSb.decreaseIndent();
            currentSb.newLine().writeIndent().append("}, ");
            K tail = k.items().get(k.items().size() - 1);
            apply(tail);
            currentSb.append(")");
            end();
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        start();
        currentSb.append("&m.InjectedKLabel{Label: m.");
        currentSb.append(nameProvider.klabelVariableName(k.klabel()));
        currentSb.append("}");
        end();
    }

}