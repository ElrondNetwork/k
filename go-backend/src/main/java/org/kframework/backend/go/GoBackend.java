// Copyright (c) 2015-2019 K Team. All Rights Reserved.
package org.kframework.backend.go;

import com.google.inject.Inject;
import org.apache.commons.io.FileUtils;
import org.kframework.backend.go.codegen.ConstantsGen;
import org.kframework.backend.go.codegen.DefinitionToGo;
import org.kframework.backend.go.codegen.EvalFunctionGen;
import org.kframework.backend.go.codegen.FreshFunctionGen;
import org.kframework.backend.go.codegen.GoBuiltin;
import org.kframework.backend.go.codegen.InterpreterDefGen;
import org.kframework.backend.go.codegen.KLabelsGen;
import org.kframework.backend.go.codegen.SortsGen;
import org.kframework.backend.go.codegen.StepFunctionGen;
import org.kframework.backend.go.codegen.StuckGen;
import org.kframework.backend.go.codegen.inline.InlineMatchGen;
import org.kframework.backend.go.codegen.inline.RuleLhsMatchDefaultWriter;
import org.kframework.backend.go.codegen.inline.RuleLhsMatchInlineManager;
import org.kframework.backend.go.codegen.inline.RuleLhsMatchWriter;
import org.kframework.backend.go.gopackage.GoPackageManager;
import org.kframework.backend.go.model.DefinitionData;
import org.kframework.backend.go.strings.GoNameProvider;
import org.kframework.backend.go.strings.GoNameProviderDebug;
import org.kframework.backend.go.strings.GoNameProviderProper;
import org.kframework.compile.Backend;
import org.kframework.definition.Definition;
import org.kframework.definition.Module;
import org.kframework.kompile.CompiledDefinition;
import org.kframework.kompile.Kompile;
import org.kframework.kompile.KompileOptions;
import org.kframework.main.GlobalOptions;
import org.kframework.utils.errorsystem.KEMException;
import org.kframework.utils.errorsystem.KExceptionManager;
import org.kframework.utils.file.FileUtil;

import java.io.File;
import java.io.IOException;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;
import java.util.function.Function;

public class GoBackend implements Backend {

    private final KExceptionManager kem;
    private final FileUtil files;
    private final GlobalOptions globalOptions;
    private final KompileOptions kompileOptions;
    private final GoOptions options;
    private GoPackageManager packageManager;

    @Inject
    public GoBackend(KExceptionManager kem, FileUtil files, GlobalOptions globalOptions, KompileOptions kompileOptions, GoOptions options) {
        this.kem = kem;
        this.files = files;
        this.globalOptions = globalOptions;
        this.kompileOptions = kompileOptions;
        this.options = options;
    }

    @Override
    public void accept(CompiledDefinition compiledDefinition) {

        System.out.println("GoBackend.accept started.");

        String mainModule = kompileOptions.mainModule(files);
        packageManager = new GoPackageManager(files, mainModule.toLowerCase(), options);
        GoNameProvider nameProvider;
        if (options.verboseVars) {
            nameProvider = new GoNameProviderDebug();
        } else {
            nameProvider = new GoNameProviderProper();
        }
        RuleLhsMatchWriter matchWriter;
        if (options.naive) {
            matchWriter = new RuleLhsMatchDefaultWriter();
        } else {
            matchWriter = new RuleLhsMatchInlineManager();
        }

        DefinitionToOcamlTempCopy ocamlDef = new DefinitionToOcamlTempCopy(kem, files, globalOptions, kompileOptions, options);
        DefinitionToGo def = new DefinitionToGo(kem, files, packageManager, nameProvider, matchWriter, globalOptions, kompileOptions, options);
        ocamlDef.initialize(compiledDefinition);
        def.initialize(compiledDefinition);

        try {
            DefinitionData data = def.definitionData();

            // temporary, for convenience and comparison
//            files.saveToKompiled("constants.ml", ocamlDef.constants());
//            files.saveToKompiled("realdef.ml", ocamlDef.definition());
//            String execution_pmg_ocaml = ocamlDef.ocamlCompile(compiledDefinition.topCellInitializer, compiledDefinition.exitCodePattern, options.dumpExitCode);
//            files.saveToKompiled("execution_pgm.ml", execution_pmg_ocaml);

            // generate: model
            packageManager.saveToPackage(packageManager.modelPackage, "klabel.go",
                    new KLabelsGen(data, packageManager, nameProvider).klabels());
            packageManager.saveToPackage(packageManager.modelPackage, "sort.go",
                    new SortsGen(data, packageManager, nameProvider).generate());

            // generate: interpreter
            packageManager.saveToPackage(packageManager.interpreterPackage, "fresh.go",
                    new FreshFunctionGen(data, packageManager, nameProvider).generate());
            packageManager.saveToPackage(packageManager.interpreterPackage, "eval.go",
                    new EvalFunctionGen(data, packageManager, nameProvider).generate());
            StepFunctionGen stepFunctionGen = new StepFunctionGen(data, packageManager, nameProvider, matchWriter);
            stepFunctionGen.generateStepFunctionCode();
            packageManager.saveToPackage(packageManager.interpreterPackage, "step.go",
                    stepFunctionGen.outputStepFunctionCode());
            packageManager.saveToPackage(packageManager.interpreterPackage, "stepRhs.go",
                    stepFunctionGen.outputStepRhsCode());
            packageManager.saveToPackage(packageManager.interpreterPackage, "functions.go",
                    def.definition());
            packageManager.saveToPackage(packageManager.interpreterPackage, "stuck.go",
                    new StuckGen(data, packageManager, nameProvider, matchWriter).generateStuck());
            packageManager.saveToPackage(packageManager.interpreterPackage, "constants.go",
                    new ConstantsGen(packageManager, data.constants).generate());
            packageManager.saveToPackage(packageManager.interpreterPackage, "interpreterDef.go",
                    new InterpreterDefGen(data, packageManager).generate());
            if (matchWriter instanceof RuleLhsMatchInlineManager) {
                packageManager.saveToPackage(packageManager.interpreterPackage, "krefInline.go",
                        new InlineMatchGen(packageManager, (RuleLhsMatchInlineManager)matchWriter).generate());
            }

        } catch (Exception e) {
            e.printStackTrace();
            return;
        }

        try {
            // copy lexer and parser
            for (String fileName : Arrays.asList(
                    "stringutil.go",
                    "model.go", "korelex.go", "koreparser.y",
                    "gen.go")) {
                packageManager.copyFileToPackage(
                        files.resolveKBase("include/go/koreparser/" + fileName),
                        packageManager.koreParserPackage, fileName);
            }

            // copy: model
            for (String fileName : Arrays.asList(
                    "collectionsUtil.go", "dynArray.go",
                    "dataIntConvert.go", "dataIntBytes.go", "dataIntOperations.go",
                    "error.go", "memo.go",
                    "printUtil.go")) {
                // these ones are the same in all implementations
                packageManager.copyFileToPackage(
                        files.resolveKBase("include/go/model/" + fileName),
                        packageManager.modelPackage, fileName);

            }
            for (String fileName : Arrays.asList(
                    "collectionsToK.go",
                    "data.go",
                    "dataBool.go",
                    "dataCollections.go", "dataCollectionsLookups.go", "dataCollectionsMap.go",
                    "dataInt.go",
                    "dataKApply.go", "dataKSequence.go", "dataKToken.go", "dataKVariable.go",
                    "dataString.go", "dataOthers.go",
                    "deepCopy.go", "equals.go",
                    "kmapkey.go",
                    "kref.go",
                    "printK.go", "printPretty.go",
                    "referenceUsageInc.go", "referenceUsageDec.go", "referenceRecycle.go", "referencePreserve.go",
                    "transfer.go")) {
                if (options.naive) {
                    packageManager.copyFileToPackage(
                            files.resolveKBase("include/go/model/naive/" + fileName),
                            packageManager.modelPackage, fileName);
                } else {
                    packageManager.copyFileToPackage(
                            files.resolveKBase("include/go/model/" + fileName),
                            packageManager.modelPackage, fileName);
                }
            }

            // kref.go also gets copied into the interpreter
            for (String fileName : Arrays.asList(
                    "kref.go")) {
                if (options.naive) {
                    packageManager.copyFileToPackage(
                            files.resolveKBase("include/go/model/naive/" + fileName),
                            packageManager.interpreterPackage, fileName);
                } else {
                    packageManager.copyFileToPackage(
                            files.resolveKBase("include/go/model/" + fileName),
                            packageManager.interpreterPackage, fileName);
                }
            }

            // copy: interpreter
            for (String fileName : Arrays.asList(
                    "interpreterFunc.go",
                    "error.go", "global.go",
                    "kmodelconvert.go",
                    "run.go",
                    "trace.go", "tracepretty.go", "tracekprint.go", "tracecompare.go")) {
                packageManager.copyFileToPackage(
                        files.resolveKBase("include/go/interpreter/" + fileName),
                        packageManager.interpreterPackage, fileName);
            }

            // copy: builtin hook files
            for (String hookNamespace : GoBuiltin.HOOK_NAMESPACES) {
                String fileName = "hooks_" + hookNamespace.toLowerCase() + ".go";
                packageManager.copyFileToPackage(
                        files.resolveKBase("include/go/hooks/" + fileName),
                        packageManager.interpreterPackage, fileName);
            }

            // copy: unit tests
            if (options.unitTests) {
                for (String fileName : Arrays.asList(
                        "kref_test.go")) {
                    packageManager.copyFileToPackage(
                            files.resolveKBase("include/go/model/" + fileName),
                            packageManager.modelPackage, fileName);
                }
                for (String fileName : Arrays.asList(
                        "hooks_array_test.go",
                        "hooks_bool_test.go",
                        "hooks_buffer_test.go",
                        "hooks_bytes_test.go",
                        "hooks_int_test.go",
                        "hooks_kequal_test.go",
                        "hooks_map_test.go",
                        "hooks_string_test.go",
                        "ksequenceutil_test.go",
                        "testutil.go")) {
                    packageManager.copyFileToPackage(
                            files.resolveKBase("include/go/hooks/unittest/" + fileName),
                            packageManager.interpreterPackage, fileName);
                }
            }

            // files that enable us to run the interpreter directly on simple cases, like imp
            if (!options.srcOnly && options.quickTest != null) {
                // quickMain
                packageManager.copyFileAndReplaceGoPackages(
                        files.resolveKBase("include/go/quickMain.go"),
                        files.resolveKompiled("quickMain.go"),
                        null);

                //save .vscode config, it is convenient for VSCode users for debugging
                copyFileAndReplaceVsCodeConfig(
                        files.resolveKBase("include/go/vscode_launch.json"),
                        files.resolveKompiled(".vscode/launch.json"),
                        options.quickTest == null ? "" : '"' + options.quickTest + '"');
            }

        } catch (IOException e) {
            throw KEMException.criticalError("Error copying go files: " + e.getMessage(), e);
        } catch (Exception e) {
            e.printStackTrace();
            throw e;
        }

        try {
            ProcessBuilder pb = files.getProcessBuilder();
            int exit;
            exit = pb.command("go", "generate").directory(files.resolveKompiled("./koreparser")).inheritIO().start().waitFor();
            if (exit != 0) {
                throw KEMException.criticalError("go generate returned nonzero exit code: " + exit + "\nExamine output to see errors.");
            }

            if (!options.srcOnly) {
                System.out.println("Starting go build.");
                exit = pb.command("go", "build").directory(files.resolveKompiled(".")).inheritIO().start().waitFor();
                if (exit != 0) {
                    throw KEMException.criticalError("go build returned nonzero exit code: " + exit + "\nExamine output to see errors.");
                }

                System.out.println("Running model unit tests.");
                exit = pb.command("go", "test")
                        .directory(files.resolveKompiled(packageManager.interpreterPackage.getRelativePath()))
                        .inheritIO().start().waitFor();
                if (exit != 0) {
                    throw KEMException.criticalError("go test returned nonzero exit code: " + exit + "\nExamine output to see errors.");
                }

                System.out.println("Running interpreter unit tests.");
                exit = pb.command("go", "test")
                        .directory(files.resolveKompiled(packageManager.modelPackage.getRelativePath()))
                        .inheritIO().start().waitFor();
                if (exit != 0) {
                    throw KEMException.criticalError("go test returned nonzero exit code: " + exit + "\nExamine output to see errors.");
                }

                if (options.quickTest != null) {
                    // execute
                    String execCommand = "./" + files.getKompiledDirectoryName();
                    exit = pb.command(execCommand, options.quickTest).directory(files.resolveKompiled(".")).inheritIO().start().waitFor();
                    if (exit != 0) {
                        throw KEMException.criticalError("interpreter returned nonzero exit code: " + exit + "\nExamine output to see errors.");
                    }


                }

            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw KEMException.criticalError("Go process interrupted.", e);
        } catch (IOException e) {
            throw KEMException.criticalError("Error starting go build process: " + e.getMessage(), e);
        }


        System.out.println("GoBackend.accept completed successfully.");
    }

    private void copyFileAndReplaceVsCodeConfig(File srcFile, File destFile, String testProgramFileName) throws IOException {
        String contents = FileUtils.readFileToString(srcFile);
        contents = contents.replaceAll("%INPUT_PGM_FILE%", testProgramFileName);
        FileUtil.save(destFile, contents);
    }

    @Override
    public Function<Definition, Definition> steps() {
        return Kompile.defaultSteps(kompileOptions, kem, files);
    }

    @Override
    public Function<Module, Module> specificationSteps(Definition ignored) {
        throw new UnsupportedOperationException();
    }

    @Override
    public Set<String> excludedModuleTags() {
        return new HashSet<>(Arrays.asList("symbolic", "kore"));
    }
}