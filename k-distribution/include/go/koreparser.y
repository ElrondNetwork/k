// This is the source for the generation of the kore parser.
// To build it:
// goyacc -o koreparser.go -p "kore" koreparser.y (produces koreparser.go)
// go build
//

%{

package main

%}

%union {
	str []byte
	k K
	kseq KSequence
	klist []K
}

%type	<k>	k kitem
%type	<klist>	klist
%type	<kseq>	ksequence

%token KSEQ DOTK '(' ')' ',' DOTKLIST TOKENLABEL KLABELLABEL KVARIABLE KLABEL STRING

%token	<str> KLABEL STRING KVARIABLE

%%

top:
	k
	{
		lastResult = $1
	}

ksequence:
	kitem KSEQ kitem
	{
		$$ = KSequence { ks:[]K{$1, $3} }
	}
|	ksequence KSEQ kitem
	{
		$$ = KSequence { ks:append($1.ks, $3) }
	}
|	DOTK
	{
		$$ = KSequence { ks:nil }
	}

k:
	kitem
	{
		$$ = $1
	}
|	ksequence
	{
		$$ = $1
	}

kitem:
	KLABEL '(' klist ')'
	{
		$$ = KApply { label:string($1), list:$3 }
	}
|   KLABELLABEL '(' KLABEL ')'
    {
        $$ = InjectedKLabel { label:string($3) }
    }
| 	TOKENLABEL '(' STRING ',' STRING ')'
	{
		$$ = KToken { value: string($3), sort: string($5) }
	}
|   KVARIABLE
    {
        $$ = KVariable { name: string($1) }
    }

klist:
	k
	{
		$$ = []K{$1}
	}
|	klist ',' k
	{
		$$ = append($1, $3)
	}
|	klist ',' ',' k
	{
		$$ = append($1, $4)
	}
|	DOTKLIST
	{
		$$ = nil
	}

%%

var lastResult K

