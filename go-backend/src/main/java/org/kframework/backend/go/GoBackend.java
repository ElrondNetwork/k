// Copyright (c) 2015-2018 Runtime Verification, Inc. (RV-Match team). All Rights Reserved.
package org.kframework.backend.go;

import com.google.inject.Inject;
import org.apache.commons.io.FileUtils;
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
    private GoPackageNameManager packageNameManager;

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
        packageNameManager = new GoPackageNameManager(files,mainModule.toLowerCase() + "interpreter");


//        DefinitionToGoCopied def = new DefinitionToGoCopied(kem, files, globalOptions, kompileOptions, options);
//        DefinitionToGoFresh fresh = new DefinitionToGoFresh(kem, files, packageNameManager, globalOptions, kompileOptions, options);
//        def.initialize(compiledDefinition);
//        fresh.initialize(compiledDefinition);
//
//        String ocaml = def.constants();
//        files.saveToKompiled("constants.ml", ocaml);
//
//        files.saveToKompiled("klabel.go", fresh.klabels());
//        files.saveToKompiled("sort.go", fresh.sorts());
//
//        try {
//            String ocamlDef = def.definition();
//            files.saveToKompiled("realdef.ml", ocamlDef);
//        } catch (Exception e) {
//            e.printStackTrace();
//        }
//        String goDef = fresh.definition();
//        files.saveToKompiled("definition.go", goDef);

        try {
            FileUtils.copyFile(files.resolveKBase("include/go/koreparser/stringutil.go"), files.resolveKompiled("koreparser/stringutil.go"));
            FileUtils.copyFile(files.resolveKBase("include/go/koreparser/model.go"), files.resolveKompiled("koreparser/model.go"));
            FileUtils.copyFile(files.resolveKBase("include/go/koreparser/korelex.go"), files.resolveKompiled("koreparser/korelex.go"));
            FileUtils.copyFile(files.resolveKBase("include/go/koreparser/koreparser.y"), files.resolveKompiled("koreparser/koreparser.y"));
            FileUtils.copyFile(files.resolveKBase("include/go/koreparser/gen.go"), files.resolveKompiled("koreparser/gen.go"));

            packageNameManager.copyFileAndReplaceGoPackages(
                    files.resolveKBase("include/go/main.go"), files.resolveKompiled("main.go"));
        } catch (IOException e) {
            throw KEMException.criticalError("Error copying go files: " + e.getMessage(), e);
        }

        try {
            ProcessBuilder pb = files.getProcessBuilder();
            int exit;
            exit = pb.command("go", "generate").directory(files.resolveKompiled("./koreparser")).inheritIO().start().waitFor();
            if (exit != 0) {
                throw KEMException.criticalError("go generate returned nonzero exit code: " + exit + "\nExamine output to see errors.");
            }
            if (!options.srcOnly) {
                exit = pb.command("go", "build").directory(files.resolveKompiled(".")).inheritIO().start().waitFor();
                if (exit != 0) {
                    throw KEMException.criticalError("go build returned nonzero exit code: " + exit + "\nExamine output to see errors.");
                }

                if (options.quickTest != null) {
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