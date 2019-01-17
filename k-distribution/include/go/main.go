package main

import (
	"fmt"
	koreparser "$INCLUDE_KORE_PARSER$"
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
	kinput := convertParserModelToKModel(parserK)
	fmt.Println("input:")
	fmt.Println(kinput.PrettyTreePrint(0))

	m := make(map[K]K)
	m[KToken{Sort: sortKConfigVar, Value: "$PGM"}] = kinput
	kmap := Map{Sort: sortMap, Label: lbl_Map_, data: m}
	evalK := KApply{Label: topCellInitializer, List: []K{kmap}}
	kresult := eval(evalK, Bottom{})
	fmt.Println("\n\noutput:")
	fmt.Println(kresult.PrettyTreePrint(0))
}

func main() {
	kastParseAndPrint()

}
