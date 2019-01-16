package org.kframework.backend.go;

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
import org.kframework.builtin.BooleanUtils;
import org.kframework.builtin.Sorts;
import org.kframework.compile.ConvertDataStructureToLookup;
import org.kframework.compile.DeconstructIntegerAndFloatLiterals;
import org.kframework.compile.ExpandMacros;
import org.kframework.compile.GenerateSortPredicateRules;
import org.kframework.compile.LiftToKSequence;
import org.kframework.compile.RewriteToTop;
import org.kframework.definition.Constructors;
import org.kframework.definition.Module;
import org.kframework.definition.ModuleTransformer;
import org.kframework.definition.Production;
import org.kframework.definition.Rule;
import org.kframework.definition.Sentence;
import org.kframework.kil.Attribute;
import org.kframework.kompile.CompiledDefinition;
import org.kframework.kompile.Kompile;
import org.kframework.kompile.KompileOptions;
import org.kframework.kore.InjectedKLabel;
import org.kframework.kore.K;
import org.kframework.kore.KApply;
import org.kframework.kore.KLabel;
import org.kframework.kore.KORE;
import org.kframework.kore.KSequence;
import org.kframework.kore.KVariable;
import org.kframework.kore.Sort;
import org.kframework.kore.VisitK;
import org.kframework.main.GlobalOptions;
import org.kframework.parser.concrete2kore.ParserUtils;
import org.kframework.utils.StringUtil;
import org.kframework.utils.algorithms.SCCTarjan;
import org.kframework.utils.errorsystem.KEMException;
import org.kframework.utils.errorsystem.KExceptionManager;
import org.kframework.utils.file.FileUtil;
import scala.Function1;

import java.io.File;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.kframework.Collections.*;
import static org.kframework.backend.go.CopiedStaticMethods.*;
import static org.kframework.definition.Constructors.*;
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
    private Map<String, List<KLabel>> klabelsForEachPredicate;

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

    DefinitionData definitionData() {
        return new DefinitionData(mainModule);
    }
//
//    public void initialize(CompiledDefinition def) {
//        //TODO: add module preprocessing
//        mainModule = def.executionModule();
//    }

    public void initialize(CompiledDefinition def) {
        Function1<Module, Module> generatePredicates = new GenerateSortPredicateRules(false)::gen;
        this.convertDataStructure = new ConvertDataStructureToLookup(def.executionModule(), true);
        ModuleTransformer convertLookups = ModuleTransformer.fromSentenceTransformer(convertDataStructure::convert, "convert data structures to lookups");
        ModuleTransformer liftToKSequence = ModuleTransformer.fromSentenceTransformer(new LiftToKSequence()::lift, "lift K into KSequence");
        this.expandMacros = new ExpandMacros(def.executionModule(), files, kompileOptions, false);
        ModuleTransformer expandMacros = ModuleTransformer.fromSentenceTransformer(this.expandMacros::expand, "expand macro rules");
        ModuleTransformer deconstructInts = ModuleTransformer.fromSentenceTransformer(new DeconstructIntegerAndFloatLiterals()::convert, "remove matches on integer literals in left hand side");
        // this.threadCellExists = containsThreadCell(def);
        this.exitCodePattern = def.exitCodePattern;
        ModuleTransformer splitThreadCell = this.threadCellExists ?
                ModuleTransformer.fromSentenceTransformer(new
                        SplitThreadsCell(def.executionModule())::convert, "split threads cell into thread local and global") :
                ModuleTransformer.fromSentenceTransformer(s -> s, "identity function -- no transformation");
        ModuleTransformer preprocessKLabelPredicates = ModuleTransformer.fromSentenceTransformer(new PreprocessKLabelPredicates(def.executionModule())::convert, "preprocess klabel predicates");
        Sentence thread = Constructors.Production(KLabel("#Thread"), Sorts.KItem(), Seq(
                Terminal("#Thread"), Terminal("("),
                NonTerminal(Sorts.K()), Terminal(","),
                NonTerminal(Sorts.K()), Terminal(","),
                NonTerminal(Sorts.K()), Terminal(","),
                NonTerminal(Sorts.K()), Terminal(")")));
        Sentence bottom = Constructors.Production(KLabel("#Bottom"), Sorts.KItem(), Seq(Terminal("#Bottom")));
        Sentence threadLocal = Constructors.Production(KLabel("#ThreadLocal"), Sorts.KItem(), Seq(Terminal("#ThreadLocal")));
        Function1<Module, Module> pipeline = preprocessKLabelPredicates
                .andThen(splitThreadCell)
                .andThen(mod -> Constructors.Module(mod.name(), mod.imports(),
                        Stream.concat(stream(mod.localSentences()),
                                Stream.<Sentence>of(thread, bottom, threadLocal)).collect(org.kframework.Collections.toSet()), mod.att()))
                .andThen(convertLookups)
                .andThen(expandMacros)
                .andThen(deconstructInts)
                .andThen(generatePredicates)
                .andThen(liftToKSequence);
        mainModule = pipeline.apply(def.executionModule());
        //mainModule = def.executionModule();
        //collectionFor = ConvertDataStructureToLookup.collectionFor(mainModule);
        //filteredMapConstructors = ConvertDataStructureToLookup.filteredMapConstructors(mainModule);
    }

    public String klabels() {
        Set<KLabel> klabels = mutable(mainModule.definedKLabels());
        klabels.add(KORE.KLabel("#Bottom"));
        klabels.add(KORE.KLabel("littleEndianBytes"));
        klabels.add(KORE.KLabel("bigEndianBytes"));
        klabels.add(KORE.KLabel("signedBytes"));
        klabels.add(KORE.KLabel("unsignedBytes"));
        addOpaqueKLabels(klabels);

        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append(" \n\n");
        sb.append("type KLabel int\n\n");

        // const declaration
        sb.append("const (\n");
        for (KLabel klabel : klabels) {
            sb.append("\t");
            GoStringUtil.appendKlabelVariableName(sb, klabel);
            sb.append(" KLabel = iota\n");
        }
        sb.append(")\n");

        // klabel name method
        sb.append("func (s KLabel) name () string {\n");
        sb.append("\tswitch s {\n");
        for (KLabel klabel : klabels) {
            sb.append("\t\tcase ");
            GoStringUtil.appendKlabelVariableName(sb, klabel);
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            sb.append(GoStringUtil.enquoteString(klabel.name()));
            sb.append("\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Unexpected KLabel.\")\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        // parse klabel function
        sb.append("func parseKLabel (name string) KLabel {\n");
        sb.append("\tswitch name {\n");
        for (KLabel klabel : klabels) {
            sb.append("\t\tcase ");
            sb.append(GoStringUtil.enquoteString(klabel.name()));
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            GoStringUtil.appendKlabelVariableName(sb, klabel);
            sb.append("\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Parsing KLabel failed. Unexpected KLabel name:\" + name)\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }

    public String sorts() {
        Set<Sort> sorts = mutable(mainModule.definedSorts());
        sorts.add(Sorts.Bool());
        sorts.add(Sorts.MInt());
        sorts.add(Sorts.Int());
        sorts.add(Sorts.String());
        sorts.add(Sorts.Float());
        sorts.add(Sorts.StringBuffer());
        sorts.add(Sorts.Bytes());

        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append(" \n\n");
        sb.append("type Sort int\n\n");

        // const declaration
        sb.append("const (\n");
        for (Sort s : sorts) {
            sb.append("\t");
            GoStringUtil.appendSortVariableName(sb, s);
            sb.append(" Sort = iota\n");
        }
        sb.append(")\n\n");

        // sort name method
        sb.append("func (s Sort) name () string {\n");
        sb.append("\tswitch s {\n");
        for (Sort sort : sorts) {
            sb.append("\t\tcase ");
            GoStringUtil.appendSortVariableName(sb, sort);
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            sb.append(GoStringUtil.enquoteString(sort.name()));
            sb.append("\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Unexpected Sort.\")\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        // parse sort function
        sb.append("func parseSort (name string) Sort {\n");
        sb.append("\tswitch name {\n");
        for (Sort sort : sorts) {
            sb.append("\t\tcase ");
            sb.append(GoStringUtil.enquoteString(sort.name()));
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            GoStringUtil.appendSortVariableName(sb, sort);
            sb.append("\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Parsing Sort failed. Unexpected Sort name:\" + name)\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }

    /**
     * TODO: not sure if needed
     */
    private void addOpaqueKLabels(Set<KLabel> klabels) {
        if (options.klabels == null)
            return;
        File definitionFile = files.resolveWorkingDirectory(options.klabels).getAbsoluteFile();
        List<File> lookupDirectories = kompileOptions.outerParsing.includes.stream().map(files::resolveWorkingDirectory).collect(Collectors.toList());
        lookupDirectories.add(Kompile.BUILTIN_DIRECTORY);
        Set<Module> mods = new ParserUtils(files::resolveWorkingDirectory, kem, globalOptions).loadModules(
                new HashSet<>(),
                "require " + StringUtil.enquoteCString(definitionFile.getPath()),
                Source.apply(definitionFile.getAbsolutePath()),
                definitionFile.getParentFile(),
                lookupDirectories,
                new HashSet<>(), false);
        mods.stream().forEach(m -> klabels.addAll(mutable(m.definedKLabels())));
    }

    public String freshDefinition() {
        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");

        sb.append("func freshFunction (s Sort, config K, counter int) K {\n");
        sb.append("\tswitch s {\n");
        for (Sort sort : iterable(mainModule.freshFunctionFor().keys())) {
            sb.append("\t\tcase ");
            GoStringUtil.appendSortVariableName(sb, sort);
            sb.append(":\n");
            sb.append("\t\t\treturn ");
            KLabel freshFunction = mainModule.freshFunctionFor().apply(sort);
            GoStringUtil.appendFunctionName(sb, freshFunction);
            sb.append("(Int(counter), config)\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\tpanic(\"Cannot find fresh function for sort \" + s.name())\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }

    /**
     * WARNING: depends on fields functions and anywhereKLabels, only run after definition()
     * TODO: untangle this dependency
     */
    public String evalDefinition() {
        StringBuilder sb = new StringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");

        sb.append("import \"fmt\"\n\n");

        sb.append("func eval(c K, config K) K {\n");
        sb.append("\tkApply, typeOk := c.(KApply)\n");
        sb.append("\tif !typeOk {\n");
        sb.append("\t\treturn c\n");
        sb.append("\t}\n");
        sb.append("\tswitch kApply.Label {\n");
        for (KLabel label : Sets.union(functions, anywhereKLabels)) {
            sb.append("\t\tcase ");
            GoStringUtil.appendKlabelVariableName(sb, label);
            sb.append(":\n");

            // arity check
            int arity = getArity(label);
            sb.append("\t\t\tif len(kApply.List) != ").append(arity).append(" {\n");
            sb.append("\t\t\t\tpanic(fmt.Sprintf(\"");
            GoStringUtil.appendFunctionName(sb, label);
            sb.append(" function arity violated. Expected arity: ").append(arity);
            sb.append(". Nr. params provided: %d\", len(kApply.List)))\n");
            sb.append("\t\t\t}\n");

            // function call
            sb.append("\t\t\treturn ");
            GoStringUtil.appendFunctionName(sb, label);
            sb.append("(");
            for (int i = 0; i < arity; i++) {
                sb.append("kApply.List[").append(i).append("], ");
            }
            sb.append("config)\n");
        }
        sb.append("\t\tdefault:\n");
        sb.append("\t\t\treturn c\n");
        sb.append("\t}\n");
        sb.append("}\n\n");

        return sb.toString();
    }

    Set<KLabel> functions;
    Set<KLabel> anywhereKLabels;

    public String definition() {
        DirectedGraph<KLabel, Object> dependencies = new DirectedSparseGraph<>();

        //TODO: warning! duplicate code with DefinitionToOcaml
        SetMultimap<KLabel, Rule> functionRules = HashMultimap.create();
        ListMultimap<KLabel, Rule> anywhereRules = ArrayListMultimap.create();
        anywhereKLabels = new HashSet<>();
        stream(mainModule.rules()).filter(r -> !r.att().contains(Attribute.MACRO_KEY) && !r.att().contains(Attribute.ALIAS_KEY)).forEach(r -> {
            K left = RewriteToTop.toLeft(r.body());
            if (left instanceof KSequence) {
                KSequence kseq = (KSequence) left;
                if (kseq.items().size() == 1 && kseq.items().get(0) instanceof KApply) {
                    KApply kapp = (KApply) kseq.items().get(0);
                    if (mainModule.attributesFor().apply(kapp.klabel()).contains(Attribute.FUNCTION_KEY)) {
                        functionRules.put(kapp.klabel(), r);
                    }
                    if (r.att().contains("anywhere")) {
                        anywhereRules.put(kapp.klabel(), r);
                        anywhereKLabels.add(kapp.klabel());
                    }
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
        SetMultimap<KLabel, Rule> klabelRuleMap = HashMultimap.create(functionRules);
        klabelsForEachPredicate = new HashMap<>();
        for (KLabel functionLabel : klabelRuleMap.keySet()) {
            if (mainModule.attributesFor().get(functionLabel).getOrElse(() -> Att()).contains("klabelPredicate")) {
                klabelsForEachPredicate.put(functionLabel.name(), computeKLabelsForPredicate(functionRules.get(functionLabel)));
            }
        }
        klabelRuleMap.putAll(anywhereRules);

        GoStringBuilder sb = new GoStringBuilder();
        sb.append("package ").append(packageNameManager.getInterpreterPackageName()).append("\n\n");

        // constants := arity == 0 && !impure
        List<List<KLabel>> functionOrder = sortFunctions(klabelRuleMap, functions, anywhereKLabels, dependencies); // result no longer required
        Set<KLabel> impurities = functions.stream().filter(lbl -> mainModule.attributesFor().apply(lbl).contains(Attribute.IMPURE_KEY)).collect(Collectors.toSet());
        impurities.addAll(ancestors(impurities, dependencies));
        Set<KLabel> constants = functions.stream().filter(lbl -> !impurities.contains(lbl) && stream(mainModule.productionsFor().apply(lbl)).filter(p -> p.arity() == 0).findAny().isPresent()).collect(Collectors.toSet());

        for (KLabel functionLabel : Sets.union(functions, anywhereKLabels)) {
            String hook = mainModule.attributesFor().get(functionLabel).getOrElse(() -> Att()).<String>getOptional(Attribute.HOOK_KEY).orElse(".");
            if (functions.contains(functionLabel)) {
                // prepare arity/params
                int arity = getArity(functionLabel);
                FunctionVars functionVars = new FunctionVars(arity);

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
                String namespace = hook.substring(0, hook.indexOf('.'));
                String function = hook.substring(namespace.length() + 1);
                if (GoBuiltin.HOOK_NAMESPACES.contains(namespace) || options.hookNamespaces.contains(namespace)) {
                    sb.writeIndent().append("//hook: ").append(hook).newLine();

                    sb.writeIndent().append("lbl := ");
                    GoStringUtil.appendKlabelVariableName(sb.sb(), functionLabel);
                    sb.append(" // ").append(functionLabel.name()).newLine(); // just for readability
                    sb.writeIndent().append("sort := ");
                    GoStringUtil.appendSortVariableName(sb.sb(), mainModule.sortFor().apply(functionLabel));
                    sb.newLine();

                    sb.writeIndent().append("if hookRes, hookErr := ");
                    GoStringUtil.appendHookMethodName(sb.sb(), namespace, function);
                    sb.append("(");
                    sb.append(functionVars.callParameters());
                    sb.append("lbl, sort, config); hookErr == nil");
                    sb.beginBlock();

                    if (mainModule.attributesFor().apply(functionLabel).contains("canTakeSteps")) {
                        sb.append("\t// eval ???\n");
                    }

                    sb.writeIndent().append("return hookRes").newLine();

                    sb.endOneBlockNoNewline().append(" else if _, errTypeOk := hookErr.(*hookNotImplementedError); !errTypeOk ").beginBlock();
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

                sb.append("\tpanic(\"Stuck!\")\n");
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

    private void convertFunction(List<Rule> rules, GoStringBuilder sb, String functionName, RuleType type, FunctionVars functionVars) {
        int ruleNum = 0;
        for (Rule r : rules) {
            if (hasLookups(r)) {
                //ruleNum = convert(Collections.singletonList(r), sb, functionName, functionName, type, ruleNum);
                sb.append("\t// rule with lookups\n");
            } else {
                sb.append("\t// rule without lookups\n");
                ruleNum = convert(r, sb, type, ruleNum, functionName, functionVars);
            }
        }
    }

    private int convert(Rule r, GoStringBuilder sb, RuleType type, int ruleNum, String functionName, FunctionVars functionVars) {
        try {
            GoStringUtil.appendRuleComment(sb, r);

            K left = RewriteToTop.toLeft(r.body());
            K requires = r.requires();
            VarInfo vars = new VarInfo();

            // convertLHS
            GoLhsVisitor lhsVisitor = new GoLhsVisitor(sb, vars, this.definitionData(), functionVars);
            if (type == RuleType.ANYWHERE || type == RuleType.FUNCTION) {
                KApply kapp = (KApply) ((KSequence) left).items().get(0);
                lhsVisitor.applyTuple(kapp.klist().items());
            } else {
                lhsVisitor.apply(left);
            }
            //String result = convert(vars);
            String suffix = "";
            boolean when = true;
            if (type == RuleType.REGULAR && options.checkRaces) {
                sb.append(" when start_after < ").append(ruleNum);
                when = false;
            }
//            if (!requires.equals(KSequence(BooleanUtils.TRUE)) || !result.equals("true")) {
//                suffix = convertSideCondition(sb, requires, vars, Collections.emptyList(), when, type, ruleNum);
//            }
            sb.writeIndent();
            sb.append("// rhs here!\n");

            for (String varName : vars.vars.values()) {
                if (!varName.equals("_")) {
                    sb.writeIndent();
                    sb.append("doNothingWithVar(").append(varName).append(") // temp\n");
                }
            }
            GoRhsVisitor rhsVisitor = new GoRhsVisitor(sb, vars, this.definitionData());
            K right = RewriteToTop.toRight(r.body());
            rhsVisitor.apply(right);
            sb.append("\n");

            //convertRHS(sb, type, r, vars, suffix, ruleNum, functionName, true);

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
            if (rhs.equals(KSequence(BooleanUtils.FALSE)) && r.att().contains("owise")) {
                continue;
            }
            if (!rhs.equals(KSequence(BooleanUtils.TRUE))) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            if (!(lhs instanceof KSequence)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            KSequence kseq = (KSequence) lhs;
            if (kseq.items().size() != 1 || !(kseq.items().get(0) instanceof KApply)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            KApply function = (KApply) kseq.items().get(0);
            if (function.items().size() != 1 || !(function.items().get(0) instanceof KSequence)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            kseq = (KSequence) function.items().get(0);
            if (kseq.items().size() != 1 || !(kseq.items().get(0) instanceof InjectedKLabel)) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            InjectedKLabel injection = (InjectedKLabel) kseq.items().get(0);
            if (injection.klabel() instanceof KVariable) {
                throw KEMException.compilerError("Unexpected form for klabel predicate rule, expected predicate(_) => false [owise] or predicate(#klabel(`klabel`)) => true.", r);
            }
            labels.add(injection.klabel());
        }
        return labels;
    }

}
