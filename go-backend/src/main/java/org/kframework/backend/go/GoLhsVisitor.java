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

import java.util.List;

class GoLhsVisitor extends VisitK {
    private final StringBuilder sb;
    private final VarInfo vars;
    private final DefinitionData data;

    private boolean inBooleanExp;
    private boolean topAnywherePre;
    private boolean topAnywherePost;

    private int indentDepth = 1;
    private String subject = "c";

    public void writeIndent() {
        for (int i = 0; i < indentDepth; i++) {
            sb.append('\t');
        }
    }

    public void beginBlock() {
        sb.append(" {\n");
        indentDepth++;
    }

    public void endAllBlocks() {
        while (indentDepth > 1) {
            indentDepth--;
            writeIndent();
            sb.append("}\n");
        }
    }

    public GoLhsVisitor(StringBuilder sb, VarInfo vars, DefinitionData data, boolean useNativeBooleanExp, boolean anywhereRule) {
        this.sb = sb;
        this.vars = vars;
        this.data = data;
        this.inBooleanExp = useNativeBooleanExp;
        this.topAnywherePre = anywhereRule;
        this.topAnywherePost = anywhereRule;
    }

    private void lhsTypeIf(String castVar, String type) {
        writeIndent();
        sb.append("if ").append(castVar).append(", t := ");
        sb.append(subject).append(".(").append(type).append("); t");
    }

    @Override
    public void apply(KApply k) {
        if (k.klabel().name().equals("#KToken")) {
            //magic down-ness
            lhsTypeIf("kt", "KToken");
            sb.append(" && kt.Sort == ");
            Sort sort = Outer.parseSort(((KToken) ((KSequence) k.klist().items().get(0)).items().get(0)).s());
            GoStringUtil.appendSortVariableName(sb, sort);
            beginBlock();
            subject = "kt.Value";
            apply(((KSequence) k.klist().items().get(1)).items().get(0));
        } else if (k.klabel().name().equals("#Bottom")) {
            lhsTypeIf("_", "Bottom");
        } else {
            topAnywherePre = false;
            lhsTypeIf("kapp", "KApply");
            sb.append(" && kapp.Label == ");
            GoStringUtil.appendKlabelVariableName(sb, k.klabel());
            beginBlock();
            applyTuple(k.klist().items());
        }
    }


    void applyTuple(List<K> items) {
        for (K item : items) {
            apply(item);
        }
    }

    @Override
    public void apply(KAs k) {
        writeIndent();
        sb.append("// apply KAs\n");
    }

    @Override
    public void apply(KRewrite k) {
        throw new AssertionError("unexpected rewrite");
    }

    @Override
    public void apply(KToken k) {
        writeIndent();
        sb.append("// apply KToken\n");
    }

    @Override
    public void apply(KVariable k) {
        String varName = GoStringUtil.variableName(k.name());

        vars.vars.put(k, varName);
        Sort s = k.att().getOptional(Sort.class).orElse(KORE.Sort(""));
        if (data.mainModule.sortAttributesFor().contains(s)) {
            String hook = data.mainModule.sortAttributesFor().apply(s).<String>getOptional("hook").orElse("");
            if (GoBuiltin.OCAML_SORT_VAR_HOOKS.containsKey(hook)) {
                // comment
                writeIndent();
                sb.append("// apply KVariable with hook:").append(hook);
                sb.append("\n");

                // code
                writeIndent();
                String pattern = GoBuiltin.GO_SORT_VAR_HOOKS.get(hook);
                sb.append(String.format(pattern,
                        varName, subject, GoStringUtil.sortVariableName(s)));
                beginBlock();
                return;
            }
        }

        if (varName.equals("_")) {
            // no code here, it is redundant
            writeIndent();
            sb.append("//");
            sb.append(varName).append(" := ").append(subject).append(" // apply KVariable\n");
        } else {
            writeIndent();
            sb.append(varName).append(" := ").append(subject).append(" // apply KVariable\n");
        }
    }

    @Override
    public void apply(KSequence k) {
        writeIndent();
        sb.append("// apply KSequence size:" + k.items().size() +"\n");
        if (k.items().size() == 1) {
            subject = "c";
            apply(k.items().get(0));
            return;
        }

        int i = 1;
        for (K item : k.items()) {
            subject = "c" + i;
            apply(item);
            i++;
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        lhsTypeIf("ikl", "InjectedKLabel");
        sb.append(" && ikl.Sort == ");
        GoStringUtil.appendKlabelVariableName(sb, k.klabel());
        beginBlock();
    }

    private void apply(List<K> items, boolean klist) {
    }
}