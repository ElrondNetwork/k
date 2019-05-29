%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	"fmt"
	koreparser "%INCLUDE_PARSER%"
	"log"
	m "%INCLUDE_MODEL%"
	"math"
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

// ExecuteSimple ... interprets the program in the file given at input
func (i *Interpreter) ExecuteSimple(kdir string, execFile string) {
	kast := callKast(kdir, execFile)
	if i.Verbose {
		fmt.Printf("Kast: %s\n\n", kast)
	}

	data := make(map[string][]byte)
	data["PGM"] = kast
	final, stepsMade, err := i.Execute(data)

	if err != nil {
		panic(err)
	}

	if i.Verbose {
		fmt.Println("\n\npretty print:")
		fmt.Println(i.Model.PrettyPrint(final))
		fmt.Println("\n\nk print:")
		fmt.Println(i.Model.KPrint(final))
		fmt.Printf("\n\nsteps made: %d\n", stepsMade)
	}
}

// Execute ... interprets the program with the structure
func (i *Interpreter) Execute(kastMap map[string][]byte) (finalState m.K, stepsMade int, err error) {
	kConfigMap := make(map[m.KMapKey]m.K)
	for key, kastValue := range kastMap {
		ktoken := m.KToken{Sort: m.SortKConfigVar, Value: "$" + key}
		parsedValue := koreparser.Parse(kastValue)
		kValue := i.convertParserModelToKModel(parsedValue)
		kConfigMap[ktoken] = kValue
	}

	// top cell initialization
	kmap := &m.Map{Sort: m.SortMap, Label: m.KLabelForMap, Data: kConfigMap}
	evalK := &m.KApply{Label: TopCellInitializer, List: []m.K{kmap}}
	kinit, err := i.Eval(evalK, m.InternedBottom)
	if err != nil {
		fmt.Println(err.Error())
		return kinit, 0, err
	}

	if i.Verbose {
		fmt.Println("\n\ntop level init:")
		fmt.Println(i.Model.PrettyPrint(kinit))
	}

	// execute
	return i.TakeStepsNoThread(kinit)
}

// TakeStepsNoThread ... executes as many steps as possible given the starting configuration
func (i *Interpreter) TakeStepsNoThread(k m.K) (finalState m.K, stepsMade int, err error) {
	i.initializeTrace()
	defer i.closeTrace()

	// start
	stepsMade = 0
	i.traceInitialState(k)

	finalState = k
	err = nil

	maxSteps := i.MaxSteps
	if maxSteps == 0 {
		// not set, it means we don't limit the number of steps
		// except when it overflows an int ... not yet sure if we need uint64, might be overkill
		maxSteps = math.MaxInt32
	}

	for stepsMade < maxSteps {
		i.traceStepStart(stepsMade, finalState)
		finalState, err = i.step(finalState)
		if err != nil {
			if _, t := err.(*noStepError); t {
				i.traceNoStep(stepsMade, finalState)
				err = nil
			}
			return
		}

		i.traceStepEnd(stepsMade, finalState)
		stepsMade++
	}
	err = errMaxStepsReached
	return
}
