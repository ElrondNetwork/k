package %PACKAGE_INTERPRETER%

import (
	"fmt"
	koreparser "%INCLUDE_PARSER%"
	"log"
	m "%INCLUDE_MODEL%"
	"os/exec"
)

func callKast(kdir string, programPath string) []byte {
	cmd := exec.Command("kast", programPath)
	cmd.Dir = kdir
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Kast error: " + err.Error())
	}
	return out
}

// ExecuteOptions ... options for executing programs
type ExecuteOptions struct {
	TraceToFile bool
}

// Execute ... interprets the program in the file given at input
func Execute(kdir string, execFile string, options ExecuteOptions) {
	kast := callKast(kdir, execFile)
	fmt.Printf("Kast: %s\n\n", kast)

	parserK := koreparser.Parse(kast)
	kinput := convertParserModelToKModel(parserK)
	fmt.Println("input:")
	fmt.Println(m.PrettyPrint(kinput))

	// top cell initialization
	initMap := make(map[m.KMapKey]m.K)
    var pgmToken = &m.KToken{Sort: m.SortKConfigVar, Value: "$PGM"}
    initMap[pgmToken.AsMapKey()] = kinput
	kmap := &m.Map{Sort: m.SortMap, Label: m.KLabelForMap, Data: initMap}
	evalK := &m.KApply{Label: topCellInitializer, List: []m.K{kmap}}
	kinit, err := eval(evalK, m.InternedBottom)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("\n\ntop level init:")
	fmt.Println(m.PrettyPrint(kinit))

	// prepare trace
	if options.TraceToFile {
		initializeTrace()
		defer closeTrace()
	}

	// execute
	final, stepsMade := takeStepsNoThread(kinit, 10000)
	fmt.Println("\n\nresult:")
	fmt.Println(m.PrettyPrint(final))

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
