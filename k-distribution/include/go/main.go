// This file holds the go generate command to run yacc on the grammar in koreparser.y.
// To build koreparser:
//	% go generate
//	% go build

//go:generate goyacc -o koreparser.go -p "kore" koreparser.y

package main

import (
	"fmt"
	"log"
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

func main() {
	kast := callKast("tests/sum.imp")
	fmt.Printf("Kast: %s\n", kast)

	//testStr := "Aaa(#token\"#token\\\"0\" `abc` `qw\\`er\"\"` .::K  .K~>.K"
	//testStr := "#abc"
	x := koreLexerImpl{line: []byte(kast)}
	//yylval := koreSymType{}

	result := koreParse(&x)
	fmt.Println(result)
	fmt.Printf("%#v\n", lastResult)

}
