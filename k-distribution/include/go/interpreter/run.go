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
	TracePretty bool
	TraceKPrint bool
	Verbose     bool
}

// ExecuteSimple ... interprets the program in the file given at input
func ExecuteSimple(kdir string, execFile string, options ExecuteOptions) {
	verbose = options.Verbose

	kast := callKast(kdir, execFile)
	if verbose {
		fmt.Printf("Kast: %s\n\n", kast)
	}

	data := make(map[string][]byte)
	data["PGM"] = kast
	final, stepsMade, err := Execute(&data, options)

	if err != nil {
		panic(err)
	}

	if verbose {
		fmt.Println("\n\npretty print:")
		fmt.Println(m.PrettyPrint(final))
		fmt.Println("\n\nk print:")
        fmt.Println(m.KPrint(final))
		fmt.Printf("\n\nsteps made: %d\n", stepsMade)
	}
}

// Execute ... interprets the program with the structure
func Execute(kastMap *map[string][]byte, options ExecuteOptions) (finalState m.K, stepsMade int, err error) {
	verbose = options.Verbose

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
		return kinit, 0, err
	}

	if verbose {
		fmt.Println("\n\ntop level init:")
		fmt.Println(m.PrettyPrint(kinit))
	}

	// prepare trace
	if options.TracePretty {
		traceHandlers = append(traceHandlers, &tracePrettyDebug{})
	}
	if options.TraceKPrint {
		traceHandlers = append(traceHandlers, &traceKPrint{})
	}
	initializeTrace()
	defer closeTrace()

	// execute
	return takeStepsNoThread(kinit, 10000)
}

func takeStepsNoThread(k m.K, maxSteps int) (finalState m.K, stepsMade int, err error) {
	stepsMade = 0
	traceInitialState(k)

	finalState = k
	err = nil

	for stepsMade < maxSteps {
		traceStepStart(stepsMade, finalState)
		finalState, err = step(finalState)
		if err != nil {
			if _, t := err.(*noStepError); t {
				traceNoStep(stepsMade, finalState)
				err = nil
			}
			return
		}

		traceStepEnd(stepsMade, finalState)
		stepsMade++
	}
	err = errMaxStepsReached
	return
}
