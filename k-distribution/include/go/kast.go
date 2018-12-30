package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("/home/andrei/elrond/k/k-distribution/target/release/k/bin/kast", "tests/sum.imp")
	cmd.Dir = "../"
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Kast: %s\n", out)


	
}