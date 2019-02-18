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
import org.apache.commons.lang3.NotImplementedException;
import org.kframework.backend.go.GoOptions;
import org.kframework.backend.go.codegen.rules.RuleWriter;
import org.kframework.backend.go.gopackage.GoExternalHookManager;
import org.kframework.backend.go.gopackage.GoPackage;
import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.model.FunctionHookName;
import org.kframework.backend.go.model.FunctionParams;
import org.kframework.backend.go.model.RuleCounter;
import org.kframework.backend.go.model.RuleInfo;
import org.kframework.backend.go.model.RuleType;
import org.kframework.backend.go.strings.GoNameProvider;
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
import scala.Option;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
import java.util.stream.Collectors;

import static org.kframework.Collections.*;
import static org.kframework.backend.go.CopiedStaticMethods.*;
import static org.kframework.kore.KORE.*;

public class DefinitionToGo {

    private transient final KExceptionManager kem;
    private transient final FileUtil files;
    private final GoPackageManager packageManager;
    private final GoNameProvider nameProvider;
    private transient final GlobalOptions globalOptions;
    private transient final KompileOptions kompileOptions;
    private transient ExpandMacros expandMacros;
    private transient ConvertDataStructureToLookup convertDataStructure;
    private boolean threadCellExists;
    private transient Rule exitCodePattern;
    public final GoOptions options;
    public final GoExternalHookManager extHookManager;

    public DefinitionToGo(
            KExceptionManager kem,
            FileUtil files,
            GoPackageManager packageManager,
            GoNameProvider nameProvider,
            GlobalOptions globalOptions,
            KompileOptions kompileOptions,
            GoOptions options) {
        this.kem = kem;
        this.files = files;
        this.packageManager = packageManager;
        this.nameProvider = nameProvider;
        this.globalOptions = globalOptions;
        this.kompileOptions = kompileOptions;
        this.options = options;
        this.extHookManager = new GoExternalHookManager(options.hookPackagePaths, packageManager);
    }

    private Module mainModule;
    private KLabel topCellInitializer;

    public DefinitionData definitionData() {
        return new DefinitionData(mainModule,
                functions, anywhereKLabels,
                functionRules, anywhereRules,
                functionParams, topCellInitializer);
    }

    RuleWriter ruleWriter;

    public void initialize(CompiledDefinition def) {
        Function1<Module, Module> generatePredicates = new GenerateSortPredicateRules(false)::gen;
        this.convertDataStructure = new ConvertDataStructureToLookup(def.executionModule(), true);
        ModuleTransformer convertLookups = ModuleTransformer.fromSentenceTransformer(convertDataStructure::convert, "convert data structures to lookups");
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

        ruleWriter = new RuleWriter(this.definitionData(), nameProvider);
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
        sb.append("package ").append(packageManager.interpreterPackage.getName()).append("\n\n");

        // generating imports
        Set<GoPackage> importsSorted = new TreeSet<>(Comparator.comparing(GoPackage::getName));
        importsSorted.add(new GoPackage("fmt", null, null));
        importsSorted.add(packageManager.modelPackage);
        for (KLabel functionLabel : functions) {
            Option<String> hook = mainModule.attributesFor().get(functionLabel).getOrElse(() -> Att()).getOption(Attribute.HOOK_KEY);
            if (hook.nonEmpty()) {
                FunctionHookName funcHook = new FunctionHookName(hook.get());
                GoPackage externalHookPkg = extHookManager.getPackage(funcHook.getExternalGoPackageName());
                if (externalHookPkg != null) {
                    importsSorted.add(externalHookPkg);
                }
            }
        }
        sb.append("import (").newLine();
        for (GoPackage pkg : importsSorted) {
            if (pkg.getGoPath() == null) {
                sb.append("\t\"").append(pkg.getName()).append("\"").newLine(); // "fmt"
            } else {
                sb.append("\t").append(pkg.getAlias()).append(" \"").append(pkg.getGoPath()).append("\"").newLine();
            }
        }
        sb.append(")\n\n");

        // constants := arity == 0 && !impure
        List<List<KLabel>> functionOrder = sortFunctions(klabelRuleMap, functions, anywhereKLabels, dependencies); // result no longer required
        Set<KLabel> impurities = functions.stream().filter(lbl -> mainModule.attributesFor().apply(lbl).contains(Attribute.IMPURE_KEY)).collect(Collectors.toSet());
        impurities.addAll(ancestors(impurities, dependencies));
        Set<KLabel> constants = functions.stream().filter(lbl -> !impurities.contains(lbl) && stream(mainModule.productionsFor().apply(lbl)).filter(p -> p.arity() == 0).findAny().isPresent()).collect(Collectors.toSet());

        RuleCounter ruleCounter = new RuleCounter();

        for (KLabel functionLabel : Sets.union(functions, anywhereKLabels)) {
            String hook = mainModule.attributesFor().get(functionLabel).getOrElse(() -> Att()).<String>getOptional(Attribute.HOOK_KEY).orElse(".");
            if (functions.contains(functionLabel)) {
                FunctionParams functionVars = functionParams.get(functionLabel);

                String functionName;
                if (mainModule.attributesFor().apply(functionLabel).contains("memo")) {
                    functionName = nameProvider.memoFunctionName(functionLabel);
                } else {
                    functionName = nameProvider.evalFunctionName(functionLabel);
                }

                assert sb.getCurrentIndent() == 0;

                // start typing
                sb.append("func ");
                sb.append(functionName);
                sb.append("(").append(functionVars.parameterDeclaration()).append("config m.K, guard int) (m.K, error)");
                sb.beginBlock();

                // if we print a return under no if, all code that follows is unreachable
                boolean unreachableCode = false;

                // hook implementation
                FunctionHookName funcHook = new FunctionHookName(hook);
                String hookCall = null;
                if (GoBuiltin.HOOK_NAMESPACES.contains(funcHook.getNamespace())) {
                    hookCall = funcHook.getGoHookObjName() + "." + funcHook.getGoFuncName();
                } else if (extHookManager.containsPackage(funcHook.getExternalGoPackageName())) {
                    hookCall = funcHook.getExternalGoPackageName() + "." + funcHook.getExternalGoFuncName();
                } else if (!hook.equals(".")) {
                    kem.registerCompilerWarning("missing entry for hook " + hook);
                }
                if (hookCall != null) {
                    sb.appendIndentedLine("//hook: ", hook);

                    sb.appendIndentedLine("lbl := m.", nameProvider.klabelVariableName(functionLabel), " // ", functionLabel.name());
                    Sort sort = mainModule.sortFor().apply(functionLabel);
                    sb.appendIndentedLine("sort := m.", nameProvider.sortVariableName(sort));

                    sb.writeIndent().append("if hookRes, hookErr := ");
                    sb.append(hookCall).append("(");
                    sb.append(functionVars.callParameters());
                    sb.append("lbl, sort, config); hookErr == nil");
                    sb.beginBlock();

                    if (mainModule.attributesFor().apply(functionLabel).contains("canTakeSteps")) {
                        throw new NotImplementedException("'canTakeSteps' attribute not implemented in Go backend.");
                    }

                    sb.writeIndent().append("return hookRes, nil").newLine();

                    sb.endOneBlockNoNewline().append(" else if _, isNotImpl := hookErr.(*hookNotImplementedError); isNotImpl ").beginBlock();
                    sb.writeIndent().append("fmt.Println(\"Warning! Call to hook ").append(hook).append(", which is not implemented.\")").newLine();
                    sb.endOneBlockNoNewline().append(" else").beginBlock();
                    sb.writeIndent().append("return m.NoResult, hookErr").newLine();
                    sb.endOneBlock().newLine();
                }

                // predicate
                if (mainModule.attributesFor().apply(functionLabel).contains(Attribute.PREDICATE_KEY, Sort.class)) {
                    Sort predicateSort = (mainModule.attributesFor().apply(functionLabel).get(Attribute.PREDICATE_KEY, Sort.class));

                    List<Sort> sorts = stream(mainModule.definedSorts())
                            .filter(s -> mainModule.subsorts().greaterThanEq(predicateSort, s))
                            .distinct()
                            .filter(sort -> mainModule.sortAttributesFor().contains(sort))
                            .collect(Collectors.toList());

                    for(Sort sort : sorts) {
                        String sortHook = mainModule.sortAttributesFor().apply(sort).<String>getOptional("hook").orElse("");
                        if (GoBuiltin.PREDICATE_RULES.containsKey(sortHook)) {
                            String sortName = "m." + nameProvider.sortVariableName(sort);
                            if (unreachableCode) {
                                sb.writeIndent().append("// unreachable predicate rule: ").append(sortHook).newLine();
                                sb.writeIndent().append("/* ");
                                sb.append(GoBuiltin.PREDICATE_RULES.get(sortHook).apply(sortName));
                                sb.append(" */").newLine();
                            } else {
                                sb.writeIndent().append("// predicate rule: ").append(sortHook).newLine();
                                sb.writeIndent();
                                sb.append(GoBuiltin.PREDICATE_RULES.get(sortHook).apply(sortName));
                                sb.newLine();

                                if (sortHook.equals("K.K") || sortHook.equals("K.KItem")) {
                                    unreachableCode = true;
                                }
                            }
                        }
                    };
                }

                // main!
                List<Rule> rules = functionRules.get(functionLabel).stream().sorted(this::sortFunctionRules).collect(Collectors.toList());
                for (Rule r : rules) {
                    if (unreachableCode) {
                        sb.appendIndentedLine("// unreachable");
                        sb.writeIndent().append("// ");
                        GoStringUtil.appendRuleComment(sb, r);
                        sb.newLine();
                    } else {
                        RuleInfo ruleInfo = ruleWriter.writeRule(r, sb, RuleType.FUNCTION, ruleCounter, functionName, functionVars);
                        if (ruleInfo.alwaysMatches()) {
                            unreachableCode = true;
                        }
                    }
                }

                if (!unreachableCode) {
                    // stuck!
                    sb.writeIndent().append("return m.NoResult, &stuckError{funcName: \"").append(functionName).append("\", args: ");
                    if (functionVars.arity() == 0) {
                        sb.append("nil");
                    } else {
                        sb.append("[]m.K{");
                        sb.append(functionVars.paramNamesSeparatedByComma());
                        sb.append("}");
                    }
                    sb.append("}").newLine();
                }

                sb.endAllBlocks(0);
                sb.newLine();

                // not yet sure if we're keeping these
                if (constants.contains(functionLabel)) {
                    sb.append("//var ").append(nameProvider.constFunctionName(functionLabel));
                    sb.append(" K = ").append(nameProvider.evalFunctionName(functionLabel));
                    sb.append("(m.InternedBottom)\n\n");
                } else if (mainModule.attributesFor().apply(functionLabel).contains("memo")) {
                    writeMemoTableAndEval(sb, functionLabel, functionVars);
                }
            } else if (anywhereKLabels.contains(functionLabel)) {
                FunctionParams functionVars = functionParams.get(functionLabel);
                String functionName = nameProvider.evalFunctionName(functionLabel);

                assert sb.getCurrentIndent() == 0;

                // start typing
                sb.appendIndentedLine("// ANYWHERE");
                sb.append("func ");
                sb.append(functionName);
                sb.append("(").append(functionVars.parameterDeclaration()).append("config m.K, guard int) (m.K, error)");
                sb.beginBlock();

                // main!
                List<Rule> rules = anywhereRules.get(functionLabel);
                for (Rule r : rules) {
                    RuleInfo ruleInfo = ruleWriter.writeRule(r, sb, RuleType.ANYWHERE, ruleCounter, functionName, functionVars);
                }

                // final return
                sb.appendIndentedLine("lbl := m.", nameProvider.klabelVariableName(functionLabel), " // ", functionLabel.name());
                sb.writeIndent().append("return m.KApply{Label: lbl, List: ");
                if (functionVars.arity() == 0) {
                    sb.append("nil");
                } else {
                    sb.append("[]m.K{");
                    sb.append(functionVars.paramNamesSeparatedByComma());
                    sb.append("}");
                }
                sb.append("}, nil").newLine();

                sb.endAllBlocks(0);
                sb.newLine();
            }
        }

        return sb.toString();
    }

    private void writeMemoTableAndEval(GoStringBuilder sb, KLabel functionLabel, FunctionParams functionArgs) {
        // table declaration
        String tableName = nameProvider.memoTableName(functionLabel);
        sb.writeIndent().append("var ").append(tableName).append(" map[");
        switch (functionArgs.arity()) {
        case 0:
        case 1:
            sb.append("m.K");
            break;
        default:
            sb.append("[").append(functionArgs.arity()).append("]m.K");
        }
        sb.append("]m.K").newLine().newLine();

        // eval function
        sb.append("func ");
        sb.append(nameProvider.evalFunctionName(functionLabel));
        sb.append("(").append(functionArgs.parameterDeclaration()).append("config m.K, guard int) (m.K, error)");
        sb.beginBlock();

        sb.writeIndent().append("memoKey := ");
        switch (functionArgs.arity()) {
        case 0:
            sb.append("m.InternedBottom");
            break;
        case 1:
            sb.append(functionArgs.varName(0));
            break;
        default:
            sb.append("[").append(functionArgs.arity()).append("]m.K{").append(functionArgs.paramNamesSeparatedByComma()).append("}");
            break;
        }
        sb.newLine();

        sb.writeIndent().append("if result, found := ").append(tableName).append("[memoKey]; found").beginBlock();
        sb.appendIndentedLine("return result, nil");
        sb.endOneBlock();
        sb.writeIndent().append("computation, err := ");
        sb.append(nameProvider.memoFunctionName(functionLabel)).append("(");
        sb.append(functionArgs.callParameters()).append("config, guard)").newLine();
        sb.writeIndent().append("if err != nil").beginBlock();
        sb.appendIndentedLine("return m.NoResult, err");
        sb.endOneBlock();
        sb.writeIndent().append(tableName).append("[memoKey] = computation").newLine();
        sb.appendIndentedLine("return computation, nil");
        sb.endOneBlock();
        sb.newLine();
    }

    private int getArity(KLabel functionLabel) {
        Set<Integer> arities = stream(mainModule.productionsFor().apply(functionLabel)).map(Production::arity).collect(Collectors.toSet());
        if (arities.size() > 1) {
            throw KEMException.compilerError("KLabel " + functionLabel + " has multiple productions with differing arities: " + mainModule.productionsFor().apply(functionLabel));
        }
        assert arities.size() == 1;
        return arities.iterator().next();
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
