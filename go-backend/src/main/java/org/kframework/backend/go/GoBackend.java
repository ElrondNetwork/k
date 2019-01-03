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

    @Inject
    public GoBackend(KExceptionManager kem, FileUtil files, GlobalOptions globalOptions, KompileOptions kompileOptions) {
        this.kem = kem;
        this.files = files;
        this.globalOptions = globalOptions;
        this.kompileOptions = kompileOptions;
    }

    @Override
    public void accept(CompiledDefinition def) {
        try {
            FileUtils.copyFile(files.resolveKBase("include/go/stringutil.go"), files.resolveKompiled("stringutil.go"));
            FileUtils.copyFile(files.resolveKBase("include/go/koremodel.go"), files.resolveKompiled("koremodel.go"));
            FileUtils.copyFile(files.resolveKBase("include/go/korelex.go"), files.resolveKompiled("korelex.go"));
            FileUtils.copyFile(files.resolveKBase("include/go/koreparser.y"), files.resolveKompiled("koreparser.y"));
            FileUtils.copyFile(files.resolveKBase("include/go/main.go"), files.resolveKompiled("main.go"));
        } catch (IOException e) {
            throw KEMException.criticalError("Error copying go files: " + e.getMessage(), e);
        }

        try {
            ProcessBuilder pb = files.getProcessBuilder();
            int exit;
            exit = pb.command("/usr/local/go/bin/go", "generate").directory(files.resolveKompiled(".")).inheritIO().start().waitFor();
            if (exit != 0) {
                throw KEMException.criticalError("go generate returned nonzero exit code: " + exit + "\nExamine output to see errors.");
            }
            exit = pb.command("go", "build").directory(files.resolveKompiled(".")).inheritIO().start().waitFor();
            if (exit != 0) {
                throw KEMException.criticalError("go build returned nonzero exit code: " + exit + "\nExamine output to see errors.");
            }

            exit = pb.command("./imp-kompiled").directory(files.resolveKompiled(".")).inheritIO().start().waitFor();
            if (exit != 0) {
                throw KEMException.criticalError("interpreter returned nonzero exit code: " + exit + "\nExamine output to see errors.");
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