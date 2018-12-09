// Copyright (c) 2015-2018 Runtime Verification, Inc. (RV-Match team). All Rights Reserved.
package org.kframework.backend.go;

import com.google.inject.AbstractModule;
import com.google.inject.Module;
import com.google.inject.multibindings.MapBinder;
import org.apache.commons.lang3.tuple.Pair;
import org.kframework.main.AbstractKModule;

import java.util.Collections;
import java.util.List;

/**
 * Created by dwightguth on 5/27/15.
 */
public class GoBackendKModule extends AbstractKModule {

    @Override
    public List<Pair<Class<?>, Boolean>> kompileOptions() {
        //return Collections.singletonList(Pair.of(OcamlOptions.class, true));
        return Collections.emptyList();
    }

    @Override
    public List<Pair<Class<?>, Boolean>> krunOptions() {
        //return Collections.singletonList(Pair.of(OcamlKRunOptions.class, true));
        return Collections.emptyList();
    }

    @Override
    public List<Module> getKompileModules() {
        List<Module> mods = super.getKompileModules();
        mods.add(new AbstractModule() {
            @Override
            protected void configure() {

                MapBinder<String, org.kframework.compile.Backend> mapBinder = MapBinder.newMapBinder(
                        binder(), String.class, org.kframework.compile.Backend.class);
                mapBinder.addBinding("go").to(GoBackend.class);
            }
        });
        return mods;
    }

    @Override
    public List<Module> getKRunModules() {
        return Collections.emptyList();
        /*return Collections.singletonList(new AbstractModule() {
            @Override
            protected void configure() {

                bindOptions(OcamlBackendKModule.this::krunOptions, binder());

                MapBinder<String, Function<org.kframework.definition.Module, Rewriter>> rewriterBinder = MapBinder.newMapBinder(
                        binder(), TypeLiteral.get(String.class), new TypeLiteral<Function<org.kframework.definition.Module, Rewriter>>() {
                        });
                rewriterBinder.addBinding("ocaml").to(OcamlRewriter.class);

                MapBinder<String, Integer> checkPointBinder = MapBinder.newMapBinder(
                        binder(), String.class, Integer.class, Names.named("checkpointIntervalMap"));
                checkPointBinder.addBinding("ocaml").toInstance(5000); //TODO(dwightguth): finesse this number

                MapBinder<ToolActivation, ExecutionMode> executionBinder = MapBinder.newMapBinder(binder(),
                        ToolActivation.class, ExecutionMode.class);

                executionBinder.addBinding(new ToolActivation.OptionActivation("--ocaml-compile")).to(OcamlCompileExecutionMode.class);
                executionBinder.addBinding(new ToolActivation.OptionActivation("--interpret")).to(InterpreterExecutionMode.class);
            }
        });*/
    }
}
