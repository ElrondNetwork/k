%COMMENT%

package %PACKAGE_INTERPRETER%

import (
    "fmt"
	m "%INCLUDE_MODEL%"
)

var verbose bool = true

func warn(message string) {
    if verbose {
        fmt.Printf("Warning: %s\n", message)
    }
}

// helps us deal with unused variables in some situations
func doNothing(c m.K) {
}

// can be handy when debugging
func debugPrint(c m.K) {
	fmt.Println(m.PrettyPrint(c))
}

// DebugPrint ... prints a K item to console, useful for debugging
func DebugPrint(c m.K) {
	fmt.Println(m.PrettyPrint(c))
}
