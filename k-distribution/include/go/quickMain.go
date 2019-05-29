%COMMENT%

// File provided for convenience, to quickly evaluate simple languages, like IMP.
// Will not work if there are external hooks.

package main

import (
	interpreter "%INCLUDE_INTERPRETER%"
	"os"
)

func main() {
	if len(os.Args) == 0 {
		panic("Argument expected. First argument should be the program to execute.")
	}
	execFileName := os.Args[1]

    i := interpreter.NewInterpreter()
    i.Verbose = true
	for _, flag := range os.Args[2:] {
		if flag == "--trace" {
			i.SetTracePretty()
		}
	}

	i.ExecuteSimple("../", execFileName)
}
