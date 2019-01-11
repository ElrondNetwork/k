
package $INTERPRETER_PACKAGE$

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"$INCLUDE_KORE_PARSER$"
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

	k := koreparser.Parse(kast)
	fmt.Println(k.PrettyTreePrint(0))
}

func main() {
	kastParseAndPrint()
}