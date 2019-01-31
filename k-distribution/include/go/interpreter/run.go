package %PACKAGE_INTERPRETER%

import (
	"fmt"
	koreparser "%INCLUDE_PARSER%"
	"log"
	m "%INCLUDE_MODEL%"
	"os/exec"
)

func callKast(kdir string, programPath string) []byte {
	cmd := exec.Command("/home/andrei/elrond/k/k-distribution/target/release/k/bin/kast", programPath)
	cmd.Dir = kdir
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Kast error: " + err.Error())
	}
	return out
}

// Execute ... interprets the program in the file given at input
func Execute(kdir string, execFile string) {
	kast := callKast(kdir, execFile)
	fmt.Printf("Kast: %s\n\n", kast)

	parserK := koreparser.Parse(kast)
	kinput := convertParserModelToKModel(parserK)
	fmt.Println("input:")
	fmt.Println(kinput.PrettyTreePrint(0))

	// top cell initialization
	initMap := make(map[m.K]m.K)
	initMap[m.KToken{Sort: m.SortKConfigVar, Value: "$PGM"}] = kinput
	kmap := m.Map{Sort: m.SortMap, Label: m.KLabelForMap, Data: initMap}
	evalK := m.KApply{Label: topCellInitializer, List: []m.K{kmap}}
	kinit, err := eval(evalK, m.Bottom{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("\n\ntop level init:")
	fmt.Println(kinit.PrettyTreePrint(0))

    // prepare trace
    initializeTrace()
    defer closeTrace()

	// execute
	final, stepsMade := takeStepsNoThread(kinit, 10000)
	fmt.Println("\n\nresult:")
	fmt.Println(final.PrettyTreePrint(0))

	fmt.Printf("\n\nsteps made: %d\n", stepsMade)

}

func takeStepsNoThread(k m.K, maxSteps int) (m.K, int) {
	n := 1
	current := k
	traceInitialState(k)

	var err error
	for n < maxSteps {
	    traceStepStart(n, current)
		current, err = step(current)
		if err != nil {
			if _, t := err.(*noStepError); t {
			    traceNoStep(n, current)
				return current, n
			}
			panic(err.Error())
		}

        traceStepEnd(n, current)

		n++
	}
	return current, n
}
