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
    options := interpreter.ExecuteOptions{TraceToFile: false}
	for _, flag := range os.Args[2:] {
		if flag == "--trace" {
			options.TraceToFile = true
		}
	}

	interpreter.Execute("../", "tests/"+execFileName, options)
}
