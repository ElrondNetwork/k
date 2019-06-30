// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go.codegen.rules;

import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.KApplySignature;
import org.kframework.backend.go.model.RuleVarContainer;
import org.kframework.backend.go.model.RuleVars;
import org.kframework.backend.go.model.TempVarCounters;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.builtin.Sorts;
import org.kframework.kil.Attribute;
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

public abstract class RuleRhsWriterBase extends VisitK {
    protected GoStringBuilder currentSb;
    protected final List<String> evalCalls = new ArrayList<>();

    protected final DefinitionData data;
    protected final GoNameProvider nameProvider;
    protected final RuleVarContainer vars;
    protected final boolean allowReuse;
    protected final int topLevelIndent;

    protected boolean newlineNext = false;

    protected void start() {
        if (newlineNext) {
            currentSb.newLine();
            currentSb.writeIndent();
            newlineNext = false;
        }
    }

    protected void end() {
    }

    public RuleRhsWriterBase(DefinitionData data,
                             GoNameProvider nameProvider,
                             RuleVarContainer vars,
                             boolean allowReuse,
                             int tabsIndent, int returnValSpacesIndent) {
        this.topLevelIndent = tabsIndent;
        this.currentSb = new GoStringBuilder(tabsIndent, returnValSpacesIndent);
        this.data = data;
        this.nameProvider = nameProvider;
        this.vars = vars;
        this.allowReuse = allowReuse;
    }

    protected abstract RuleRhsWriterBase newInstanceWithSameConfig(int indent);

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
            currentSb.append("i.Model.NewKToken(m.").append(nameProvider.sortVariableName(sort));
            currentSb.append(", ");
            apply(value);
            currentSb.append(")");
        } else if (k.klabel().name().equals("#Bottom")) {
            currentSb.append("m.InternedBottom");
        } else if (data.functions.contains(k.klabel()) || data.anywhereKLabels.contains(k.klabel())) {
            applyKApplyExecute(k);
        } else {
            applyKApplyAsIs(k);
        }
        end();
    }

    private void applyKApplyAsIs(KApply k) {
        if (allowReuse &&
                vars.lhsVars.getKApplySignatureCount(k) == 1 &&
                vars.rhsVars.getKApplySignatureCount(k) == 1 ) {
            // reuse KApply!
            currentSb.append("i.Model.ReuseKApply(");
            currentSb.append(vars.lhsVars.signatureToKApplyVarName.get(KApplySignature.of(k)));
            currentSb.append(", ");
        } else {
            currentSb.append("i.Model.NewKApply(");
        }
        currentSb.append("m.").append(nameProvider.klabelVariableName(k.klabel()));
        currentSb.append(", // as-is ").append(k.klabel().name());
        currentSb.increaseIndent();
        for (K item : k.klist().items()) {
            newlineNext = true;
            apply(item);
            currentSb.append(",");
        }
        currentSb.decreaseIndent();
        currentSb.newLine();
        currentSb.writeIndent();
        currentSb.append(")");
    }

    protected void applyKApplyExecute(KApply k) {
        int evalVarIndex = vars.varCounters.consumeEvalVarIndex();
        String evalVarName = "eval" + evalVarIndex;
        String errVarName = "err" + evalVarIndex;

        // return the eval variable
        currentSb.append(evalVarName);

        // also add an eval call to the eval calls
        GoStringBuilder evalSb = new GoStringBuilder(topLevelIndent, 0);
        GoStringBuilder backupSb = currentSb;
        currentSb = evalSb; // we trick all nodes below to output to the eval call instead of the return by changing the string builder

        String hook = data.mainModule.attributesFor().apply(k.klabel()).<String>getOptional(Attribute.HOOK_KEY).orElse("");
        switch (hook) {
        case "KEQUAL.ite":
            shortCircuitedIfThenElse(k, evalVarName, evalSb);
            break;
        case "BOOL.and":
        case "BOOL.andThen":
            shortCircuitedAnd(k, evalVarName, evalSb);
            break;
        case "BOOL.or":
        case "BOOL.orElse":
            shortCircuitedOr(k, evalVarName, evalSb);
            break;
        default:
            // no hook
            // regular eval"KEQUAL.ite"
            // evalX, errX := func <funcName> (...)
            evalSb.writeIndent().append(evalVarName).append(", ").append(errVarName).append(" := ");
            evalSb.append("i.");
            evalSb.append(nameProvider.evalFunctionName(k.klabel())); // func name
            if (k.items().size() == 0) { // call parameters
                evalSb.append("(config, -1) // ").append(ToKast.apply(k)).newLine();
            } else {
                evalSb.append("( // ").append(ToKast.apply(k));
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
            break;
        }

        evalCalls.add(evalSb.toString());
        assert currentSb == evalSb;
        currentSb = backupSb; // restore
    }

    private void shortCircuitedIfThenElse(KApply k, String evalVarName, GoStringBuilder evalSb) {
        // if-then-else hook
        // if arg0 { arg1 } else { arg2 }
        // optimization aside, it is important that we avoid evaluating the branch we don't need
        // because there are cases when evaluating this branch causes the entire execution to fail
        assert k.klist().items().size() == 3;

        // arg0 (the condition) is treated normally
        evalSb.appendIndentedLine("var ", evalVarName, " m.KReference // ", ToKast.apply(k));
        evalSb.writeIndent().append("if m.IsTrue(");
        apply(k.klist().items().get(0)); // 1st argument, the condition
        evalSb.append(")").beginBlock("rhs if-then-else");

        // we separate the evaluation of arg1 from the rest
        // and place its entire evaluation tree in the if block
        // we need a new writer to avoid mixing arg1's evalCalls with the current ones
        RuleRhsWriterBase ifBranchWriter = newInstanceWithSameConfig(evalSb.getCurrentIndent());
        ifBranchWriter.apply(k.klist().items().get(1));

        ifBranchWriter.writeEvalCalls(evalSb);
        evalSb.writeIndent().append(evalVarName).append(" = ");
        ifBranchWriter.writeReturnValue(evalSb);

        // } else {
        evalSb.newLine().endOneBlockNoNewline().append(" else").beginBlock();

        // we separate the evaluation of arg2 from the rest
        // and place its entire evaluation tree in the else block
        // we need a new writer to avoid mixing arg2's evalCalls with the current ones
        RuleRhsWriterBase elseBranchWriter = newInstanceWithSameConfig(evalSb.getCurrentIndent());
        elseBranchWriter.apply(k.klist().items().get(2));

        elseBranchWriter.writeEvalCalls(evalSb);
        evalSb.writeIndent().append(evalVarName).append(" = ");
        elseBranchWriter.writeReturnValue(evalSb);

        // }
        evalSb.newLine().endOneBlock();
    }

    private void shortCircuitedAnd(KApply k, String evalVarName, GoStringBuilder evalSb) {
        // and hook, implemented as:
        // result = arg1
        // if result {
        //   result = arg2
        // }
        // optimization aside, it is important that we avoid to unnecessarily evaluate the second argument
        // because there are cases when evaluating it causes the entire execution to fail
        assert k.klist().items().size() == 2;

        evalSb.appendIndentedLine("var ", evalVarName, " m.KReference // ", ToKast.apply(k));
        evalSb.writeIndent().append(evalVarName).append(" = ");
        apply(k.klist().items().get(0));
        evalSb.newLine();
        evalSb.writeIndent().append("if m.IsTrue(").append(evalVarName).append(")").beginBlock();

        // all evaluations for arg2 need to happen in this block,
        // to avoid executing any of them if arg1 is false
        RuleRhsWriterBase arg2Writer = newInstanceWithSameConfig(evalSb.getCurrentIndent());
        arg2Writer.apply(k.klist().items().get(1));

        arg2Writer.writeEvalCalls(evalSb);
        evalSb.writeIndent().append(evalVarName).append(" = ");
        arg2Writer.writeReturnValue(evalSb);
        evalSb.newLine();
        evalSb.endOneBlock();
    }

    private void shortCircuitedOr(KApply k, String evalVarName, GoStringBuilder evalSb) {
        // or hook, implemented as:
        // result = arg1
        // if !result {
        //   result = arg2
        // }
        // optimization aside, it is important that we avoid to unnecessarily evaluate the second argument
        // because there are cases when evaluating it causes the entire execution to fail
        assert k.klist().items().size() == 2;

        evalSb.appendIndentedLine("var ", evalVarName, " m.KReference // ", ToKast.apply(k));
        evalSb.writeIndent().append(evalVarName).append(" = ");
        apply(k.klist().items().get(0));
        evalSb.newLine();
        evalSb.writeIndent().append("if !m.IsTrue(").append(evalVarName).append(")").beginBlock();

        // all evaluations for arg2 need to happen in this block,
        // to avoid executing any of them if arg1 is true
        RuleRhsWriterBase arg2Writer = newInstanceWithSameConfig(evalSb.getCurrentIndent());
        arg2Writer.apply(k.klist().items().get(1));

        arg2Writer.writeEvalCalls(evalSb);
        evalSb.writeIndent().append(evalVarName).append(" = ");
        arg2Writer.writeReturnValue(evalSb);
        evalSb.newLine();
        evalSb.endOneBlock();
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
                if (k.s().equals("0")) {
                    sb.append("m.IntZero");
                    return;
                }
                String ivarName = nameProvider.constVariableName("Int", k.s());
                String ivalue = "m.NewIntConstant(\"" + k.s() + "\")";
                data.constants.intConstants.put(ivarName, ivalue);
                sb.append(ivarName);
                return;
            case "FLOAT.Float":
                sb.append("&m.Float{Value: ").append(k.s()).append("}");
                return;
            case "STRING.String":
                if (k.s().equals("")) {
                    sb.append("m.StringEmpty");
                    return;
                }
                String unquotedStr = StringUtil.unquoteKString(k.s());
                String goStr = GoStringUtil.enquoteString(unquotedStr);
                String svarName = nameProvider.constVariableName("String", k.s());
                String svalue = "m.NewStringConstant(" + k.s() + ")";
                data.constants.stringConstants.put(svarName, svalue);
                sb.append(svarName);
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

        String ktvarName = nameProvider.constVariableName(
                "KToken",
                nameProvider.sortVariableName(k.sort()) + k.s());
        String ktvalue = "m.NewKTokenConstant(m." + nameProvider.sortVariableName(k.sort()) +
                ", " + GoStringUtil.enquoteString(k.s()) + ")";
        data.constants.tokenConstants.put(ktvarName, ktvalue);
        sb.append(ktvarName);
    }

    @Override
    public void apply(KVariable v) {
        start();
        String varName = vars.lhsVars.getVarName(v);
        if (varName == null) {
            currentSb.append("/* varName=null */ m.InternedBottom");
            end();
            return;
        }

        if (!vars.lhsVars.containsVar(v) && varName.startsWith("?")) {
            throw KEMException.internalError("Failed to compile rule due to unmatched variable on right-hand-side. This is likely due to an unsupported collection pattern: " + varName, v);
        } else if (!vars.lhsVars.containsVar(v)) {
            currentSb.append("panic(\"Stuck!\")");
        } else {
            KLabel listVar = vars.lhsVars.listVars.get(varName);
            if (listVar != null) {
                Sort sort = data.mainModule.sortFor().apply(listVar);
                currentSb.append("&m.List{Sort: m.").append(nameProvider.sortVariableName(sort));
                currentSb.append(", Label: m.").append(nameProvider.klabelVariableName(listVar));
                //currentSb.append(", ");
                //currentSb.append(varOccurrance);
                currentSb.append(" /* ??? */}");
            } else {
                if (vars.kapplyVariableNames.contains(varName)) {
                    currentSb.append(varName).append(".GetReference()");
                } else {
                    currentSb.append(varName);
                }
            }
        }
        end();
    }

    @Override
    public void apply(KSequence k) {
        int size = k.items().size();
        switch (size) {
        case 0:
            start();
            currentSb.append("m.EmptyKSequence");
            end();
            return;
        case 1:
            currentSb.append("/* rhs KSequence size=1 */ ");
            apply(k.items().get(0));
            return;
        default:
            start();
            currentSb.append("i.Model.AssembleKSequence(");
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
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        start();
        currentSb.append("m.NewInjectedKLabel(");
        currentSb.append(nameProvider.klabelVariableName(k.klabel()));
        currentSb.append(")");
        end();
    }

}