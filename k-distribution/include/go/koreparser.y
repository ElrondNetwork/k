// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is an example of a goyacc program.
// To build it:
// goyacc -p "expr" expr.y (produces y.go)
// go build -o expr y.go
// expr
// > <type an expression>

%{

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"bufio"
)

%}

%union {
	str []byte
}

%type	<str>	expr

%token KSEQ DOTK '(' ')' ',' DOTKLIST TOKENLABEL KLABELLABEL KLABEL ID_KLABEL KVARIABLE STRING

%token	<str>	NUM

%%

top:
	expr
	{
		fmt.Println("Hello")
	}

expr:
	KSEQ 
	{
	
	}


%%



func oldMain() {
	in := bufio.NewReader(os.Stdin)
	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			log.Fatalf("WriteString: %s", err)
		}
		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}

		koreParse(&koreLex{line: line})
	}
}
