package org.kframework.backend.go;

import com.google.inject.Inject;
import org.kframework.compile.Backend;
import org.kframework.definition.Definition;
import org.kframework.definition.Module;
import org.kframework.kompile.CompiledDefinition;
import org.kframework.kompile.Kompile;
import org.kframework.kompile.KompileOptions;
import org.kframework.main.GlobalOptions;
import org.kframework.utils.errorsystem.KExceptionManager;
import org.kframework.utils.file.FileUtil;

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
        StringBuilder sb = new StringBuilder();
        sb.append("package main\n");
        sb.append("\n");
        sb.append("import \"fmt\"\n");
        sb.append("\n");
        sb.append("func main() {\n");
        sb.append("\tfmt.Printf(\"hello, world from kgo!\\n\")\n");
        sb.append("}\n");

        files.saveToKompiled("hello.go", sb.toString());
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