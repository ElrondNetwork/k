package %PACKAGE_INTERPRETER%

import (
    "fmt"
	m "%INCLUDE_MODEL%"
)

var freshCounter int

// helps us deal with unused variables in some situations
func doNothing(c m.K) {
}

// can be handy when debugging
func debugPrint(c m.K) {
	fmt.Println(m.PrettyPrint(c))
}
