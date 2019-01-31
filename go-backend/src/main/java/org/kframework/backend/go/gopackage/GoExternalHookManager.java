package org.kframework.backend.go.gopackage;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class GoExternalHookManager {

    private final Map<String, GoPackage> nameToPackage;
    private final GoPackageManager packageManager;

    public GoExternalHookManager(List<String> hookPackagePaths, GoPackageManager packageManager) {
        this.packageManager = packageManager;
        nameToPackage = new HashMap<>();
        for (String hookPackagePath : hookPackagePaths) {
            GoPackage extHookPackage = packageManager.packageFromGoPath(hookPackagePath);
            nameToPackage.put(extHookPackage.getName(), extHookPackage);
        }
    }

    public boolean containsPackage(String packageName) {
        return nameToPackage.containsKey(packageName);
    }

    public GoPackage getPackage(String packageName) {
        return nameToPackage.get(packageName);
    }

}
