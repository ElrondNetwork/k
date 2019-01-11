package main

import (
	"fmt"
	koreparser "kgoimp/imp-kompiled/koreparser"
	"log"
	"os"
	"os/exec"
)

func callKast(programPath string) []byte {
	cmd := exec.Command("/home/andrei/elrond/k/k-distribution/target/release/k/bin/kast", programPath)
	cmd.Dir = "../"
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func kastParseAndPrint() {
	testArg := os.Args[1]
	kast := callKast("tests/" + testArg)
	fmt.Printf("Kast: %s\n\n", kast)

	parserK := koreparser.Parse(kast)
	k := convertParserModelToKModel(parserK)
	fmt.Println(k.PrettyTreePrint(0))
}

func main() {
	kastParseAndPrint()

}
