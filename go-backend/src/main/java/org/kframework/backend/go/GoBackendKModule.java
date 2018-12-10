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

}
