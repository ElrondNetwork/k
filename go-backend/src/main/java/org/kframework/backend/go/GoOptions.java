// Copyright (c) 2015-2018 Runtime Verification, Inc. (RV-Match team). All Rights Reserved.
package org.kframework.backend.go;

import com.beust.jcommander.Parameter;
import org.kframework.utils.inject.RequestScoped;
import org.kframework.utils.options.StringListConverter;

import java.io.Serializable;
import java.util.Collections;
import java.util.List;

@RequestScoped
public class GoOptions implements Serializable {

    public boolean profileRules;

    @Parameter(names="--go-src-only", description="Do not build definition; only generate .go files.")
    public boolean srcOnly;

    @Parameter(names="--go-hook-packages", listConverter=StringListConverter.class, description="<string> is a whitespace-separated list of paths to external Go packages required by the definition.")
    public List<String> hookPackagePaths = Collections.emptyList();

    @Parameter(names="--go-quick-test", description="After generating sources and build, also run interpreter with a test program.")
    public String quickTest;

    @Parameter(names="--go-verbose-vars", description="Generate more verbose variable names. They can be easier to read when debugging, but are not Go idiomatic.")
    public boolean verboseVars;

    @Parameter(names="--go-unit-tests", description="Also create some unit tests in target folder.")
    public boolean unitTests = true;

    // TODO: clean up
    //@Parameter(names="--no-link-prelude", description="Do not link interpreter binaries against constants.cmx and prelude.cmx. Do not use this if you don't know what you're doing.")
    public boolean noLinkPrelude;

    // TODO: clean up
    //@Parameter(names="--no-expand-macros", description="Do not expand macros on initial configurations at runtime. Will likely cause incorrect behavior if any macros are used in this term.")
    public boolean noExpandMacros;

    // TODO: clean up
    //@Parameter(names="--opaque-klabels", description="Declare all the klabels declared by the following secondary definition.")
    public String klabels;

    // TODO: clean up
    //@Parameter(names="--ocaml-dump-exit-code", description="Exit code which should trigger a dump of the configuration when using --ocaml-compile.")
    public Integer dumpExitCode;

    // TODO: clean up
    //@Parameter(names="--ocaml-serialize-config", listConverter=StringListConverter.class, description="<string> is a whitespace-separated list of configuration variables to precompute the value of")
    public List<String> serializeConfig = Collections.emptyList();

    // TODO: clean up
    //@Parameter(names="-O2", description="Optimize in ways that improve performance, but intere with debugging and increase compilation time and code size slightly.")
    public boolean optimize2;

    // TODO: clean up
    //@Parameter(names="-O3", description="Optimize aggressively in ways that significantly improve performance, but also increase compilation time and code size.")
    public boolean optimize3;

    // TODO: clean up
    //@Parameter(names="-Og", description="Optimize as much as possible without interfering with debugging experience.")
    public boolean optimizeG;

    // TODO: clean up
    //@Parameter(names="--reverse-rules", description="Reverse the order of rules as much as possible in order to find most nondeterminism without searching.")
    public boolean reverse;

    // TODO: clean up
    //@Parameter(names="--check-races", description="Checks for races among regular rules.")
    public boolean checkRaces;

    // TODO: clean up
    public boolean ocamlopt() { return optimize2 || optimize3; }

    // TODO: clean up
    public boolean optimizeStep() { return optimize3 || optimizeG; }
}
