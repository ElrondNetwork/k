package main

import (
	interpreter "%INCLUDE_INTERPRETER%"
	"os"
)

func main() {
    testArg := os.Args[1]
	interpreter.Execute("../", "tests/" + testArg)
}
