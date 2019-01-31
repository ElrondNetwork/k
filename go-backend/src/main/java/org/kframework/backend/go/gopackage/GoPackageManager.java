package org.kframework.backend.go.gopackage;

import org.apache.commons.io.FileUtils;
import org.kframework.backend.go.GoOptions;
import org.kframework.utils.errorsystem.KEMException;
import org.kframework.utils.file.FileUtil;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;

public class GoPackageManager {

    private final FileUtil files;
    private final Path goSrcPath;
    public final GoPackage koreParserPackage;
    public final GoPackage interpreterPackage;
    public final GoPackage modelPackage;

    public GoPackageManager(
            FileUtil files,
            String languageName,
            GoOptions options) {
        this.files = files;

        String goSrc = System.getenv("GOSRC");
        if (goSrc == null) {
            throw KEMException.criticalError("GOSRC environment variable not set. This should point to the $GOPATH/src/");
        }

        try {
            goSrcPath = new File(goSrc).getCanonicalFile().toPath();

            // TODO: make package output path configurable
            this.koreParserPackage = packageFromRelativePath("koreparser", "./koreparser");

            String modelPackageName = languageName + "model";
            this.modelPackage = packageFromRelativePath(modelPackageName, "./" + modelPackageName);
            this.modelPackage.setAlias("m");

            String interpreterPackageName = languageName + "interpreter";
            this.interpreterPackage = packageFromRelativePath(interpreterPackageName, "./" + interpreterPackageName);
        } catch (IOException e) {
            throw KEMException.criticalError("Failed to initialize GoPackageManager, error: " + e.getMessage(), e);
        }
    }

    public GoPackage packageFromRelativePath(String pkgName, String relativePath) {
        try {
            Path absPath = files.resolveKompiled("./" + relativePath).getCanonicalFile().toPath();
            Path relPath = goSrcPath.relativize(absPath);
            String goPath = relPath.toString();
            return new GoPackage(pkgName, goPath, relativePath);
        } catch (IOException e) {
            throw KEMException.criticalError("Failed to initialize GoPackage, error computing relative paths: " + e.getMessage(), e);
        }
    }

    public GoPackage packageFromGoPath(String goPath) {
        String pkgName = new File(goPath).getName().toLowerCase();
        return new GoPackage(pkgName, goPath, null);
    }

    public void saveToPackage(GoPackage pkg, String fileName, String contents) throws IOException {
        files.saveToKompiled(
                pkg.getRelativePath() + "/" + fileName,
                contents);
    }

    public void copyFileToPackage(File srcFile, GoPackage pkg, String fileName) throws IOException {
        copyFileAndReplaceGoPackages(
                srcFile,
                files.resolveKompiled(pkg.getRelativePath() + "/" + fileName));
    }

    private static final String PACKAGE_INTERPRETER = "%PACKAGE_INTERPRETER%";
    private static final String INCLUDE_INTERPRETER = "%INCLUDE_INTERPRETER%";
    private static final String PACKAGE_MODEL = "%PACKAGE_MODEL%";
    private static final String INCLUDE_MODEL = "%INCLUDE_MODEL%";
    private static final String PACKAGE_PARSER = "%PACKAGE_PARSER%";
    private static final String INCLUDE_PARSER = "%INCLUDE_PARSER%";

    public void copyFileAndReplaceGoPackages(File srcFile, File destFile) throws IOException {
        // TODO: optimize, by doing all regex in one go
        // stream directly to file
        // solution here: https://stackoverflow.com/questions/1326682/java-replacing-multiple-different-substring-in-a-string-at-once-or-in-the-most
        String contents = FileUtils.readFileToString(srcFile);
        contents = contents.replaceAll(PACKAGE_INTERPRETER, interpreterPackage.getName());
        contents = contents.replaceAll(INCLUDE_INTERPRETER, interpreterPackage.getGoPath());
        contents = contents.replaceAll(PACKAGE_MODEL, modelPackage.getName());
        contents = contents.replaceAll(INCLUDE_MODEL, modelPackage.getGoPath());
        contents = contents.replaceAll(PACKAGE_PARSER, koreParserPackage.getName());
        contents = contents.replaceAll(INCLUDE_PARSER, koreParserPackage.getGoPath());
        FileUtil.save(destFile, contents);
    }

}
