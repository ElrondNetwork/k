package org.kframework.backend.go.codegen;

import com.google.common.collect.ArrayListMultimap;
import com.google.common.collect.BiMap;
import com.google.common.collect.HashBiMap;
import com.google.common.collect.HashMultimap;
import com.google.common.collect.ListMultimap;
import com.google.common.collect.SetMultimap;
import com.google.common.collect.Sets;
import edu.uci.ics.jung.graph.DirectedGraph;
import edu.uci.ics.jung.graph.DirectedSparseGraph;
import org.kframework.attributes.Location;
import org.kframework.attributes.Source;
import org.kframework.backend.go.GoOptions;
import org.kframework.backend.go.GoPackageNameManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionHookName;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.processors.AccumulateRuleVars;
import org.kframework.backend.go.processors.PrecomputePredicates;
import org.kframework.backend.go.strings.GoStringBuilder;
import org.kframework.backend.go.strings.GoStringUtil;
import org.kframework.builtin.BooleanUtils;
import org.kframework.compile.ConvertDataStructureToLookup;
import org.kframework.compile.DeconstructIntegerAndFloatLiterals;
import org.kframework.compile.ExpandMacros;
import org.kframework.compile.GenerateSortPredicateRules;
import org.kframework.compile.RewriteToTop;
import org.kframework.definition.Module;
import org.kframework.definition.ModuleTransformer;
import org.kframework.definition.Production;
import org.kframework.definition.Rule;
import org.kframework.kil.Attribute;
import org.kframework.kompile.CompiledDefinition;
import org.kframework.kompile.KompileOptions;
import org.kframework.kore.InjectedKLabel;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KLabel;
import org.kframework.kore.KSequence;
import org.kframework.kore.KVariable;
import org.kframework.kore.Sort;
import org.kframework.kore.VisitK;
import org.kframework.main.GlobalOptions;
import org.kframework.utils.algorithms.SCCTarjan;
import org.kframework.utils.errorsystem.KEMException;
import org.kframework.utils.errorsystem.KExceptionManager;
import org.kframework.utils.file.FileUtil;
import scala.Function1;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Set;
import java.util.stream.Collectors;

import static org.kframework.Collections.*;
import static org.kframework.backend.go.CopiedStaticMethods.*;
import static org.kframework.kore.KORE.*;

public class DefinitionToGo {

    private transient final KExceptionManager kem;
    private transient final FileUtil files;
    private final GoPackageNameManager packageNameManager;
    private transient final GlobalOptions globalOptions;
    private transient final KompileOptions kompileOptions;
    private transient ExpandMacros expandMacros;
    private transient ConvertDataStructureToLookup convertDataStructure;
    private boolean threadCellExists;
    private transient Rule exitCodePattern;
    public GoOptions options;

    public DefinitionToGo(
            KExceptionManager kem,
            FileUtil files,
            GoPackageNameManager packageNameManager,
            GlobalOptions globalOptions,
            KompileOptions kompileOptions,
            GoOptions options) {
        this.kem = kem;
        this.files = files;
        this.packageNameManager = packageNameManager;
        this.globalOptions = globalOptions;
        this.kompileOptions = kompileOptions;
        this.options = options;
    }

    private Module mainModule;
    private KLabel topCellInitializer;

    public DefinitionData definitionData() {
        return new DefinitionData(mainModule, functions, anywhereKLabels, functionParams, topCellInitializer);
    }

    public void initialize(CompiledDefinition def) {
        Function1<Module, Module> generatePredicates = new GenerateSortPredicateRules(false)::gen;
        this.convertDataStructure = new ConvertDataStructureToLookup(def.executionModule(), true);
        ModuleTransformer convertLookups = ModuleTransformer.fromSentenceTransformer(convertDataStructure::convert, "convertVars data structures to lookups");
        this.expandMacros = new ExpandMacros(def.executionModule(), files, kompileOptions, false);
        ModuleTransformer expandMacros = ModuleTransformer.fromSentenceTransformer(this.expandMacros::expand, "expand macro rules");
        ModuleTransformer deconstructInts = ModuleTransformer.fromSentenceTransformer(new DeconstructIntegerAndFloatLiterals()::convert, "remove matches on integer literals in left hand side");
        this.exitCodePattern = def.exitCodePattern;
        ModuleTransformer identity = ModuleTransformer.fromSentenceTransformer(s -> s, "identity function -- no transformation");

        Function1<Module, Module> pipeline = identity
                .andThen(convertLookups)
                .andThen(expandMacros)
                .andThen(deconstructInts)
                .andThen(generatePredicates);
        mainModule = pipeline.apply(def.executionModule());
        topCellInitializer = def.topCellInitializer;

        //TODO: warning! duplicate code with DefinitionToOcaml
        functionRules = HashMultimap.create();
        anywhereRules = ArrayListMultimap.create();
        anywhereKLabels = new HashSet<>();
        stream(mainModule.rules()).filter(r -> !r.att().contains(Attribute.MACRO_KEY) && !r.att().contains(Attribute.ALIAS_KEY)).forEach(r -> {
            K left = RewriteToTop.toLeft(r.body());
            if (left instanceof KApply) {
                KApply kapp = (KApply) left;
                if (mainModule.attributesFor().apply(kapp.klabel()).contains(Attribute.FUNCTION_KEY)) {
                    functionRules.put(kapp.klabel(), r);
                }
                if (r.att().contains("anywhere")) {
                    anywhereRules.put(kapp.klabel(), r);
                    anywhereKLabels.add(kapp.klabel());
                }
            }
        });

        //TODO: warning! duplicate code with DefinitionToOcaml
        functions = new HashSet<>(functionRules.keySet());
        for (Production p : iterable(mainModule.productions())) {
            if (p.att().contains(Attribute.FUNCTION_KEY)) {
                functions.add(p.klabel().get());
            }
        }

        //TODO: warning! duplicate code with DefinitionToOcaml
        klabelRuleMap = HashMultimap.create(functionRules);
        klabelRuleMap.putAll(anywhereRules);

        // prepare arity/params
        functionParams = new HashMap<>();
        for (KLabel label : Sets.union(functions, anywhereKLabels)) {
            int arity = getArity(label);
            FunctionParams functionVars = new FunctionParams(arity);
            functionParams.put(label, functionVars);
        }

    }

    SetMultimap<KLabel, Rule> functionRules;
    ListMultimap<KLabel, Rule> anywhereRules;
    Set<KLabel> functions;
    Set<KLabel> anywhereKLabels;
    SetMultimap<KLabel, Rule> klabelRuleMap;
    Map<KLabel, FunctionParams> functionParams;

    public String definition() {
        DirectedGraph<KLabel, Object> dependencies = new DirectedSparseGraph<>();

        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");

        sb.append("import \"fmt\"\n\n");

        // constants := arity == 0 && !impure
        List<List<KLabel>> functionOrder = sortFunctions(klabelRuleMap, functions, anywhereKLabels, dependencies); // result no longer required
        Set<KLabel> impurities = functions.stream().filter(lbl -> mainModule.attributesFor().apply(lbl).contains(Attribute.IMPURE_KEY)).collect(Collectors.toSet());
        impurities.addAll(ancestors(impurities, dependencies));
        Set<KLabel> constants = functions.stream().filter(lbl -> !impurities.contains(lbl) && stream(mainModule.productionsFor().apply(lbl)).filter(p -> p.arity() == 0).findAny().isPresent()).collect(Collectors.toSet());

        for (KLabel functionLabel : Sets.union(functions, anywhereKLabels)) {
            String hook = mainModule.attributesFor().get(functionLabel).getOrElse(() -> Att()).<String>getOptional(Attribute.HOOK_KEY).orElse(".");
            if (functions.contains(functionLabel)) {
                FunctionParams functionVars = functionParams.get(functionLabel);

                String functionName;
                if (mainModule.attributesFor().apply(functionLabel).contains("memo")) {
                    functionName = GoStringUtil.memoFunctionName(functionLabel);
                } else {
                    functionName = GoStringUtil.functionName(functionLabel);
                }

                assert sb.currentIndent() == 0;

                // start typing
                sb.append("func ");
                sb.append(functionName);
                sb.append("(").append(functionVars.parameterDeclaration()).append("config K) K");
                sb.beginBlock();


                // hook implementation
                FunctionHookName funcHook = new FunctionHookName(hook);
                if (GoBuiltin.HOOK_NAMESPACES.contains(funcHook.getNamespace()) ||
                        options.hookNamespaces.contains(funcHook.getNamespace())) {
                    sb.writeIndent().append("//hook: ").append(hook).newLine();

                    sb.writeIndent().append("lbl := ");
                    GoStringUtil.appendKlabelVariableName(sb.sb(), functionLabel);
                    sb.append(" // ").append(functionLabel.name()).newLine(); // just for readability
                    sb.writeIndent().append("sort := ");
                    GoStringUtil.appendSortVariableName(sb.sb(), mainModule.sortFor().apply(functionLabel));
                    sb.newLine();

                    sb.writeIndent().append("if hookRes, hookErr := ");
                    sb.append(funcHook.getGoHookObjName()).append(".").append(funcHook.getGoFuncName());
                    sb.append("(");
                    sb.append(functionVars.callParameters());
                    sb.append("lbl, sort, config); hookErr == nil");
                    sb.beginBlock();

                    if (mainModule.attributesFor().apply(functionLabel).contains("canTakeSteps")) {
                        sb.append("\t// eval ???\n");
                    }

                    sb.writeIndent().append("return hookRes").newLine();

                    sb.endOneBlockNoNewline().append(" else if _, isNotImpl := hookErr.(*hookNotImplementedError); isNotImpl ").beginBlock();
                    sb.writeIndent().append("fmt.Println(\"Warning! Call to hook ").append(hook).append(", which is not implemented.\")").newLine();
                    sb.endOneBlockNoNewline().append(" else").beginBlock();
                    sb.writeIndent().append("panic(\"Unexpected error occured while running hook function.\")").newLine();
                    sb.endOneBlock().newLine();
                } else if (!hook.equals(".")) {
                    kem.registerCompilerWarning("missing entry for hook " + hook);
                }

                // predicate
                if (mainModule.attributesFor().apply(functionLabel).contains(Attribute.PREDICATE_KEY, Sort.class)) {
                    Sort predicateSort = (mainModule.attributesFor().apply(functionLabel).get(Attribute.PREDICATE_KEY, Sort.class));
                    stream(mainModule.definedSorts()).filter(s -> mainModule.subsorts().greaterThanEq(predicateSort, s)).distinct()
                            .filter(sort -> mainModule.sortAttributesFor().contains(sort)).forEach(sort -> {
                        String sortHook = mainModule.sortAttributesFor().apply(sort).<String>getOptional("hook").orElse("");
                        if (GoBuiltin.PREDICATE_RULES.containsKey(sortHook)) {
                            sb.append("\t// predicate rule: ").append(sortHook).append("\n");
                            sb.append("\t");
                            sb.append(GoBuiltin.PREDICATE_RULES.get(sortHook).apply(sort));
                            sb.append("\n");
                        }
                    });
                }

                // main!
                List<Rule> rules = functionRules.get(functionLabel).stream().sorted(this::sortFunctionRules).collect(Collectors.toList());
                convertFunction(rules, sb, functionName, RuleType.FUNCTION, functionVars);

                // TODO: will have to convert to a nice returned error
                sb.writeIndent().append("panic(\"Stuck! Function: ").append(functionName).append(" Args:\"");
                for (String fv : functionVars.getVarNames()) {
                    sb.append(" + \"\\n\\t").append(fv).append(":\" + ").append(fv).append(".PrettyTreePrint(0)");
                }
                sb.append(")").newLine();
                sb.endAllBlocks(0);
                sb.append("\n");

                // not yet sure if we're keeping these
                if (constants.contains(functionLabel)) {
                    sb.append("//var ");
                    GoStringUtil.appendConstFunctionName(sb.sb(), functionLabel);
                    sb.append(" K = ");
                    GoStringUtil.appendFunctionName(sb.sb(), functionLabel);
                    sb.append("(internedBottom)\n\n");
                } else if (mainModule.attributesFor().apply(functionLabel).contains("memo")) {
                    sb.append("//memoization not yet implemented. Function name: ");
                    sb.append(functionLabel.toString());
                    sb.append("\n\n");
                    //encodeMemoizationOfFunction(sb, conn, functionLabel, functionName, arity);
                }
            }
        }

        return sb.toString();
    }

    private int getArity(KLabel functionLabel) {
        Set<Integer> arities = stream(mainModule.productionsFor().apply(functionLabel)).map(Production::arity).collect(Collectors.toSet());
        if (arities.size() > 1) {
            throw KEMException.compilerError("KLabel " + functionLabel + " has multiple productions with differing arities: " + mainModule.productionsFor().apply(functionLabel));
        }
        assert arities.size() == 1;
        return arities.iterator().next();
    }

    private void convertFunction(List<Rule> rules, GoStringBuilder sb, String functionName, RuleType type, FunctionParams functionVars) {
        int ruleNum = 0;
        for (Rule r : rules) {
            if (hasLookups(r)) {
                //ruleNum = convertVars(Collections.singletonList(r), sb, functionName, functionName, type, ruleNum);
                sb.append("\t// rule with lookups\n");
            } else {
                sb.append("\t// rule without lookups\n");
                ruleNum = convert(r, sb, type, ruleNum, functionName, functionVars);
            }
        }
    }

    private int convert(Rule r, GoStringBuilder sb, RuleType type, int ruleNum, String functionName, FunctionParams functionVars) {
        try {
            GoStringUtil.appendRuleComment(sb, r);

            K left = RewriteToTop.toLeft(r.body());
            K requires = r.requires();
            K right = RewriteToTop.toRight(r.body());

            // we need the variables beforehand, so we retrieve them here
            AccumulateRuleVars accumLhsVars = new AccumulateRuleVars();
            accumLhsVars.apply(left);

            // some evaluations can be precomputed
            PrecomputePredicates optimizeTransf = new PrecomputePredicates(
                    this.definitionData(), accumLhsVars.vars());
            requires = optimizeTransf.apply(requires);
            right = optimizeTransf.apply(right);

            // check which variables are actually used in requires or in rhs
            // note: this has to happen *after* PrecomputePredicates does its job
            AccumulateRuleVars accumRhsVars = new AccumulateRuleVars();
            accumRhsVars.apply(requires);
            accumRhsVars.apply(right);

            // output LHS
            GoLhsVisitor lhsVisitor = new GoLhsVisitor(sb, this.definitionData(), functionVars,
                    accumLhsVars.vars(),
                    accumRhsVars.vars());
            if (type == RuleType.ANYWHERE || type == RuleType.FUNCTION) {
                KApply kapp = (KApply) left;
                lhsVisitor.applyTuple(kapp.klist().items());
            } else {
                lhsVisitor.apply(left);
            }

//            boolean when = true;
//            if (type == RuleType.REGULAR && options.checkRaces) {
//                sb.append(" when start_after < ").append(ruleNum);
//                when = false;
//            }

            // output requires
            if (!requires.equals(BooleanUtils.TRUE)) {
                sb.writeIndent().append("/* REQUIRES */").newLine();
                sb.writeIndent().append("if ");
                sb.enableMiniIndent("if ");
                // condition starts here
                GoSideConditionVisitor sideCondVisitor = new GoSideConditionVisitor(sb, this.definitionData(),
                        accumLhsVars.vars());
                sideCondVisitor.apply(requires);
                // condition ends
                sb.disableMiniIndent();
                sb.beginBlock();
            } else if (requires.att().contains(PrecomputePredicates.COMMENT_KEY)) {
                // just a comment, so we know what happened
                sb.writeIndent().append("/* REQUIRES precomputed ");
                sb.append(requires.att().get(PrecomputePredicates.COMMENT_KEY));
                sb.append(" */").newLine();
            }

            // output RHS
            sb.writeIndent().append("// rhs here:").newLine();
            GoRhsVisitor rhsVisitor = new GoRhsVisitor(sb, this.definitionData(),
                    accumLhsVars.vars());
            sb.writeIndent();
            sb.append("return ");
            rhsVisitor.apply(right);
            sb.append("\n");

            // done
            sb.endAllBlocks(GoStringBuilder.FUNCTION_BODY_INDENT);
            sb.append("\n");
            return ruleNum + 1;
        } catch (NoSuchElementException e) {
            System.err.println(r);
            throw e;
        } catch (KEMException e) {
            e.exception.addTraceFrame("while compiling rule at " + r.att().getOptional(Source.class).map(Object::toString).orElse("<none>") + ":" + r.att().getOptional(Location.class).map(Object::toString).orElse("<none>"));
            throw e;
        }
    }

    private int numLookups(Rule r) {
        class Holder {
            int i;
        }
        Holder h = new Holder();
        new VisitK() {
            @Override
            public void apply(KApply k) {
                if (ConvertDataStructureToLookup.isLookupKLabel(k)) {
                    h.i++;
                }
                super.apply(k);
            }
        }.apply(r.requires());
        return h.i;
    }

    private boolean hasLookups(Rule r) {
        return numLookups(r) > 0;
    }

    private int sortFunctionRules(Rule a1, Rule a2) {
        return Boolean.compare(a1.att().contains("owise"), a2.att().contains("owise"));
    }

    private List<List<KLabel>> sortFunctions(SetMultimap<KLabel, Rule> functionRules,
                                             Set<KLabel> functions,
                                             Set<KLabel> anywhereKLabels,
                                             final DirectedGraph<KLabel, Object> dependencies) {
        BiMap<KLabel, Integer> mapping = HashBiMap.create();
        int counter = 0;
        for (KLabel lbl : functions) {
            mapping.put(lbl, counter++);
        }
        for (KLabel lbl : anywhereKLabels) {
            mapping.put(lbl, counter++);
        }
        mapping.put(KLabel(""), counter++); //use blank klabel to simulate dependencies on eval
        List<Integer>[] predecessors = new List[functions.size() + anywhereKLabels.size() + 1];
        for (int i = 0; i < predecessors.length; i++) {
            predecessors[i] = new ArrayList<>();
        }

        class GetPredecessors extends VisitK {
            private final KLabel current;

            public GetPredecessors(KLabel current) {
                this.current = current;
            }

            @Override
            public void apply(KApply k) {
                if (functions.contains(k.klabel()) || anywhereKLabels.contains(k.klabel())) {
                    predecessors[mapping.get(current)].add(mapping.get(k.klabel()));
                    dependencies.addEdge(new Object(), current, k.klabel());
                }
                if (k.klabel() instanceof KVariable) {
                    // this function requires a call to eval, so we need to add the dummy dependency
                    predecessors[mapping.get(current)].add(mapping.get(KLabel("")));
                    dependencies.addEdge(new Object(), current, KLabel(""));
                }
                super.apply(k);
            }
        }

        for (Map.Entry<KLabel, Rule> entry : functionRules.entries()) {
            GetPredecessors visitor = new GetPredecessors(entry.getKey());
            visitor.apply(entry.getValue().body());
            visitor.apply(entry.getValue().requires());
        }

        for (KLabel label : Sets.union(functions, anywhereKLabels)) {
            String hook = mainModule.attributesFor().apply(label).<String>getOptional(Attribute.HOOK_KEY).orElse(".");

            if (mainModule.attributesFor().apply(label).contains("canTakeSteps")) {
                // this function requires a call to eval, so we need to add the dummy dependency
                predecessors[mapping.get(label)].add(mapping.get(KLabel("")));
                dependencies.addEdge(new Object(), label, KLabel(""));
            }

            if (hook.equals("KREFLECTION.fresh")) {
                for (KLabel freshFunction : iterable(mainModule.freshFunctionFor().values())) {
                    predecessors[mapping.get(label)].add(mapping.get(freshFunction));
                    dependencies.addEdge(new Object(), label, freshFunction);
                }
            }
            //eval depends on everything
            predecessors[mapping.get(KLabel(""))].add(mapping.get(label));
            dependencies.addEdge(new Object(), KLabel(""), label);
        }

        List<List<Integer>> components = new SCCTarjan().scc(predecessors);

        return components.stream().map(l -> l.stream()
                .map(i -> mapping.inverse().get(i)).collect(Collectors.toList()))
                .collect(Collectors.toList());
    }

    private static List<KLabel> computeKLabelsForPredicate(Set<Rule> rules) {
        List<KLabel> labels = new ArrayList<>();
        for (Rule r : rules) {
            K body = r.body();
            K lhs = RewriteToTop.toLeft(body);
            K rhs = RewriteToTop.toRight(body);
            if (rhs.equals(BooleanUtils.FALSE) && r.att().contains("owise")) {
                continue;
            }
            if (!rhs.equals(BooleanUtils.TRUE)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            if (!(lhs instanceof KSequence)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            if (!(lhs instanceof KApply)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            KApply function = (KApply) lhs;
            if (!(function.items().get(0) instanceof InjectedKLabel)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            InjectedKLabel injection = (InjectedKLabel) function.items().get(0);
            if (injection.klabel() instanceof KVariable) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            labels.add(injection.klabel());
        }
        return labels;
    }

}
