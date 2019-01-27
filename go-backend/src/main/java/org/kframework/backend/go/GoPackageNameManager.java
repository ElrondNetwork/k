package org.kframework.backend.go;

import org.apache.commons.io.FileUtils;
import org.kframework.utils.errorsystem.KEMException;
import org.kframework.utils.file.FileUtil;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;

public class GoPackageNameManager {

    private final FileUtil files;
    private final String interpreterPackageName;
    private final String interpreterInclude;

    private final String koreParserPackageName;
    private final String koreParserInclude;

    public GoPackageNameManager(
            FileUtil files,
            String interpreterPackageName) {
        this.files = files;
        this.interpreterPackageName = interpreterPackageName;
        this.koreParserPackageName = "koreparser";

        String goSrc = System.getenv("GOSRC");
        if (goSrc == null) {
            throw KEMException.criticalError("GOSRC environment variable not set. This should point to the $GOPATH/src/");
        }

        try {
            Path goSrcPath = new File(goSrc).getCanonicalFile().toPath();

            // koreparser path init
            // TODO: make package output path configurable
            Path koreParserAbsPath = files.resolveKompiled("./" + koreParserPackageName).getCanonicalFile().toPath();
            Path koreParserRelPath = goSrcPath.relativize(koreParserAbsPath);
            koreParserInclude = koreParserRelPath.toString();

            Path interpreterAbsPath = files.resolveKompiled("./" + interpreterPackageName).getCanonicalFile().toPath();
            Path interpreterRelPath = goSrcPath.relativize(interpreterAbsPath);
            interpreterInclude = interpreterRelPath.toString();
        } catch (IOException e) {
            throw KEMException.criticalError("Failed to initialize GoPackageNameManager, error computing relative paths: " + e.getMessage(), e);
        }
    }

    public String getInterpreterPackageName() {
        return interpreterPackageName;
    }

    public String getKoreParserPackageName() {
        return koreParserPackageName;
    }


    public void saveToInterpreterDir(String fileName, String contents) throws IOException {
        files.saveToKompiled(
                interpreterPackageName + "/" + fileName,
                contents);
    }

    public void copyFileToInterpreterDir(File srcFile, String fileName) throws IOException {
        copyFileAndReplaceGoPackages(
                srcFile,
                files.resolveKompiled(interpreterPackageName + "/" + fileName));
    }

    private static final String PACKAGE_INTERPRETER = "%PACKAGE_INTERPRETER%";
    private static final String INCLUDE_INTERPRETER = "%INCLUDE_INTERPRETER%";
    private static final String PACKAGE_PARSER = "%PACKAGE_PARSER%";
    private static final String INCLUDE_PARSER = "%INCLUDE_PARSER%";

    public void copyFileAndReplaceGoPackages(File srcFile, File destFile) throws IOException {
        // TODO: optimize, by doing all regex in one go
        // stream directly to file
        // solution here: https://stackoverflow.com/questions/1326682/java-replacing-multiple-different-substring-in-a-string-at-once-or-in-the-most
        String contents = FileUtils.readFileToString(srcFile);
        contents = contents.replaceAll(PACKAGE_INTERPRETER, interpreterPackageName);
        contents = contents.replaceAll(INCLUDE_INTERPRETER, interpreterInclude);
        contents = contents.replaceAll(PACKAGE_PARSER, koreParserPackageName);
        contents = contents.replaceAll(INCLUDE_PARSER, koreParserInclude);
        FileUtil.save(destFile, contents);
    }

}
