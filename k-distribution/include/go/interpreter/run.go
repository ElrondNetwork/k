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

// ExecuteSimple ... interprets the program in the file given at input
func ExecuteSimple(kdir string, execFile string, options ExecuteOptions) {
	kast := callKast(kdir, execFile)
	fmt.Printf("Kast: %s\n\n", kast)

	data := make(map[string][]byte)
	data["PGM"] = kast
	Execute(&data, options)
}

// Execute ... interprets the program with the structure
func Execute(kastMap *map[string][]byte, options ExecuteOptions) {

	kConfigMap := make(map[m.KMapKey]m.K)
	for key, kastValue := range *kastMap {
		ktoken := &m.KToken{Sort: m.SortKConfigVar, Value: "$" + key}
		ktokenAsKey, _ := m.MapKey(ktoken)
		parsedValue := koreparser.Parse(kastValue)
		kValue := convertParserModelToKModel(parsedValue)
		kConfigMap[ktokenAsKey] = kValue
	}

	// top cell initialization
	kmap := &m.Map{Sort: m.SortMap, Label: m.KLabelForMap, Data: kConfigMap}
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
