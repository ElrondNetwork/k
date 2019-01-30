package main

import (
	"os"
)

func main() {
    testArg := os.Args[1]
	Execute("../", "tests/" + testArg)
}
