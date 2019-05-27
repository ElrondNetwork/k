%COMMENT%

package %PACKAGE_INTERPRETER%

import (
    "fmt"
	m "%INCLUDE_MODEL%"
)

func (i *Interpreter) warn(message string) {
    if i.Verbose {
        fmt.Printf("Warning: %s\n", message)
    }
}

// helps us deal with unused variables in some situations
func doNothing(c m.K) {
}

// DebugPrint ... prints a K item to console, useful for debugging
func (i *Interpreter) DebugPrint(c m.K) {
	fmt.Println(i.Model.PrettyPrint(c))
}
