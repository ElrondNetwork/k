package org.kframework.backend.go.codegen.lhstree;

import org.kframework.attributes.Att;
import org.kframework.backend.go.codegen.GoBuiltin;
import org.kframework.backend.go.codegen.lhstree.model.LhsBottomTreeNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsHashKTokenTreeNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKApplyArgNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKApplyNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKSeqEmptyNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKSeqHeadNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKSeqMatchNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKSeqOneNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKSeqSplitNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKSeqTailNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKTokenNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKVarAlreadySeenNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKVarNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsKVarPredicateNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsListEmptyNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsListHeadNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsListMatchSplitNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsListTailNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsTopArgTreeNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsTopTreeNode;
import org.kframework.backend.go.codegen.lhstree.model.LhsTreeNode;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.kil.Attribute;
import org.kframework.kore.InjectedKLabel;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KAs;
import org.kframework.kore.KLabel;
import org.kframework.kore.KORE;
import org.kframework.kore.KRewrite;
import org.kframework.kore.KSequence;
import org.kframework.kore.KToken;
import org.kframework.kore.KVariable;
import org.kframework.kore.Sort;
import org.kframework.kore.VisitK;
import org.kframework.parser.outer.Outer;
import org.kframework.unparser.ToKast;
import org.kframework.utils.errorsystem.KEMException;

import java.util.Collection;
import java.util.Set;

public class RuleLhsTreeBuilder extends VisitK {
    private final DefinitionData data;
    private final GoNameProvider nameProvider;
    private final FunctionParams functionVars;

    public final LhsTopTreeNode topNode = new LhsTopTreeNode();
    private LhsTreeNode currentNode = topNode;

    /**
     * Whenever we see a variable more than once, instead of adding a variable declaration, we add a check that the two instances are equal.
     * This structure keeps track of that.
     */
    private final Set<KVariable> alreadySeenVariables;

    public RuleLhsTreeBuilder(
            DefinitionData data,
            GoNameProvider nameProvider,
            FunctionParams functionVars,
            Set<KVariable> alreadySeenVariables) {
        this.data = data;
        this.nameProvider = nameProvider;
        this.functionVars = functionVars;
        this.alreadySeenVariables = alreadySeenVariables;
    }

    private KVariable nextAlias = null;

    private KVariable consumeAlias() {
        if (nextAlias != null) {
            KVariable alias = nextAlias;
            nextAlias = null;
            return alias;
        }
        return null;
    }

    public void applyTopArgs(Collection<K> topArgs) {
        if (topArgs.size() != functionVars.arity()) {
            throw new RuntimeException("top-level arguments do not match rule structure");
        }
        int i = 0;
        for (K topArg : topArgs) {
            LhsTopArgTreeNode node = new LhsTopArgTreeNode(topNode, i, functionVars.varName(i));
            currentNode.chainNode(node);
            currentNode = node;
            i++;
            apply(topArg);
        }
    }

    public LhsTreeNode getLastNode() {
        return currentNode;
    }

    public void addNode(LhsTreeNode node) {
        currentNode.chainNode(node);
        currentNode = node;
    }

    @Override
    public void apply(KApply k) {
        if (k.klabel().name().equals("#KToken")) {
            assert k.klist().items().size() == 2;
            KToken ktoken = (KToken) k.klist().items().get(0);
            Sort sort = Outer.parseSort(ktoken.s());
            K value = k.klist().items().get(1);

            String sortName = nameProvider.sortVariableName(sort);
            LhsHashKTokenTreeNode node = new LhsHashKTokenTreeNode(currentNode, sortName);
            currentNode.chainNode(node);
            currentNode = node;
            apply(value);
        } else if (k.klabel().name().equals("#Bottom")) {
            LhsBottomTreeNode node = new LhsBottomTreeNode(currentNode);
            currentNode.chainNode(node);
            currentNode = node;
        } else if (data.functions.contains(k.klabel())) {
            if (data.collectionFor.containsKey(k.klabel())) {
                KLabel collectionLabel = data.collectionFor.get(k.klabel());
                Att attr = data.mainModule.attributesFor().apply(collectionLabel);
                if (attr.contains(Attribute.ASSOCIATIVE_KEY)
                        && !attr.contains(Attribute.COMMUTATIVE_KEY)
                        && !attr.contains(Attribute.IDEMPOTENT_KEY)) {

                    // list
                    Sort sort = data.mainModule.sortFor().apply(collectionLabel);
                    if (k.items().size() == 0) {
                        // empty list
                        LhsListEmptyNode node = new LhsListEmptyNode(currentNode,
                                nameProvider.sortVariableName(sort),
                                nameProvider.klabelVariableName(collectionLabel),
                                "empty list "+ ToKast.apply(k));
                        currentNode.chainNode(node);
                        currentNode = node;
                    } else if (k.items().size() == 2) {
                        LhsListMatchSplitNode listNode = new LhsListMatchSplitNode(currentNode,
                                nameProvider.sortVariableName(sort),
                                nameProvider.klabelVariableName(collectionLabel),
                                "list "+ ToKast.apply(k));
                        currentNode.chainNode(listNode);
                        currentNode = listNode;



                        KApply headListElem = (KApply) k.items().get(0);
                        boolean isElement = attr.contains("element") && headListElem.klabel().equals(KORE.KLabel(attr.get("element")));
                        // boolean isWrapElement = !isElement && attr.contains("wrapElement") && kapp.klabel().equals(KORE.KLabel(attr.get("wrapElement")));
                        if (!isElement) {
                            throw KEMException.internalError("First argument of list cons a list element type");
                        }
                        if (headListElem.klist().size() != 1) {
                            throw KEMException.internalError("List element should only have 1 argument");
                        }
                        LhsListHeadNode headNode = new LhsListHeadNode(listNode);
                        currentNode.chainNode(headNode);
                        currentNode = headNode;
                        apply(headListElem.items().get(0)); // not the element itself, but its contents

                        // tail
                        LhsListTailNode tailNode = new LhsListTailNode(listNode);
                        currentNode.chainNode(tailNode);
                        currentNode = tailNode;
                        if (k.items().get(1) instanceof KApply) {
                            KApply tail = (KApply) k.items().get(1);
                            apply(tail);
                        } else if (k.items().get(1) instanceof KVariable) {
                            KVariable tail = (KVariable) k.items().get(1);
                            apply(tail);
                        } else {
                            throw KEMException.internalError("Second argument of list cons should be either a KApply of a KVariable, representing the tail");
                        }
                    } else {
                        throw KEMException.internalError("List KApply should be either of length 0 (empty list), or length 2 (head-tail)");
                    }
                }
            }
        } else {
            int arity = k.klist().items().size();
            String comment = ToKast.apply(k);
            KVariable alias = consumeAlias();
            if (alias != null) {
                comment += " as " + alias.name();
            }

            LhsKApplyNode kappNode = new LhsKApplyNode(currentNode, nameProvider.klabelVariableName(k.klabel()), arity, comment);
            currentNode.chainNode(kappNode);
            currentNode = kappNode;

            int i = 0;
            for (K item : k.klist().items()) {
                LhsKApplyArgNode argNode = new LhsKApplyArgNode(kappNode, i);
                currentNode.chainNode(argNode);
                currentNode = argNode;
                apply(item);
                i++;
            }
        }
    }

    @Override
    public void apply(KAs k) {
        if (!(k.alias() instanceof KVariable)) {
            throw new IllegalArgumentException("KAs alias is not a KVariable.");
        }
        nextAlias = (KVariable) k.alias();
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
        LhsKTokenNode node = new LhsKTokenNode(currentNode, k);
        currentNode.chainNode(node);
        currentNode = node;
    }

    @Override
    public void apply(KVariable k) {
        if (alreadySeenVariables.contains(k)) {
            LhsKVarAlreadySeenNode node = new LhsKVarAlreadySeenNode(currentNode, k);
            currentNode.chainNode(node);
            currentNode = node;
            return;
        }
        alreadySeenVariables.add(k);

        Sort s = k.att().getOptional(Sort.class).orElse(KORE.Sort(""));
        if (data.mainModule.sortAttributesFor().contains(s)) {
            String hook = data.mainModule.sortAttributesFor().apply(s).<String>getOptional("hook").orElse("");
            if (GoBuiltin.PREDICATE_HOOKS.contains(hook)) {
                //String comment = ToKast.apply((K) k) + ": lhs KVariable with hook:" + hook;
                String comment = "lhs KVariable with hook:" + hook;
                LhsKVarPredicateNode predicateNode = new LhsKVarPredicateNode(currentNode, hook,
                        nameProvider.sortVariableName(s),
                        comment);
                currentNode.chainNode(predicateNode);
                currentNode = predicateNode;
            }
        }

        LhsKVarNode node = new LhsKVarNode(currentNode, k);
        currentNode.chainNode(node);
        currentNode = node;
    }

    @Override
    public void apply(KSequence k) {
        switch (k.items().size()) {
        case 0:
            // comment only
            LhsKSeqEmptyNode emptyNode = new LhsKSeqEmptyNode(currentNode, ToKast.apply(k));
            currentNode.chainNode(emptyNode);
            currentNode = emptyNode;
            return;
        case 1:
            // no KSequence, go straight to the only item
            LhsKSeqOneNode node = new LhsKSeqOneNode(currentNode, ToKast.apply(k));
            currentNode.chainNode(node);
            currentNode = node;
            apply(k.items().get(0));
            return;
        default:
            int nrHeads = k.items().size() - 1;

            LhsKSeqMatchNode matchNode = new LhsKSeqMatchNode(currentNode, nrHeads, ToKast.apply(k));
            currentNode.chainNode(matchNode);
            currentNode = matchNode;

            LhsTreeNode kseqNode = matchNode;
            LhsKSeqHeadNode[] headNodes = new LhsKSeqHeadNode[nrHeads];
            LhsKSeqTailNode tailNode = null;
            for (int i = 0; i < nrHeads; i++) {
                // split into head :: tail, if subject is KSequence; subject :: emptySequence otherwise
                // if multiple heads required, split repeatedly
                String splitComment = ToKast.apply(k.items().get(i)) + " ~> ...";
                LhsKSeqSplitNode splitNode = new LhsKSeqSplitNode(kseqNode, splitComment);
                currentNode.chainNode(splitNode);
                currentNode = splitNode;

                headNodes[i] = new LhsKSeqHeadNode(splitNode);


                tailNode = new LhsKSeqTailNode(splitNode);
                if (i < nrHeads-1) {
                    // add all tails except the last
                    currentNode.chainNode(tailNode);
                    currentNode = tailNode;
                }

                // the tail becomes the next subject to split
                kseqNode = tailNode;
            }

            // process heads
            for (int i = 0; i < nrHeads; i++) {
                currentNode.chainNode(headNodes[i]);
                currentNode = headNodes[i];
                apply(k.items().get(i));
            }

            // process tail
            currentNode.chainNode(tailNode);
            currentNode = tailNode;
            apply(k.items().get(nrHeads));

            return;
        }
    }

    @Override
    public void apply(InjectedKLabel k) {
        throw new RuntimeException("not implemented");
    }

}