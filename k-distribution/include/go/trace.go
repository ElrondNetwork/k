// Handles generation of traces
// (what rules were applied, in what order, what were the intermediate states).
// Generates a file and dumps trace there.

package %PACKAGE_INTERPRETER%

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var traceEnabled = false
var traceFile *os.File
var traceWriter *bufio.Writer

func initializeTrace() {
	traceEnabled = true
	formattedNow := time.Now().Format("20060102150405")
	var err error
	traceFile, err = os.OpenFile("trace_"+formattedNow+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error while creating trace file:" + err.Error())
	}

	traceWriter = bufio.NewWriter(traceFile)
}

func closeTrace() {
	traceEnabled = false
	traceWriter.Flush()
	traceFile.Close()
}

func traceInitialState(state K) {
    traceWriter.WriteString("initial state:\n\n")
	traceWriter.WriteString(state.PrettyTreePrint(0))
}

func traceStepStart(stepNr int, currentState K) {
	if !traceEnabled {
		return
	}
	traceWriter.WriteString(fmt.Sprintf("\n\nstep #%d begin\n\n", stepNr))
}

func traceStepEnd(stepNr int, currentState K) {
	if !traceEnabled {
		return
	}
	traceWriter.WriteString(fmt.Sprintf("\nstep #%d end; current state:\n\n", stepNr))
	traceWriter.WriteString(currentState.PrettyTreePrint(0))
}

func traceNoStep(stepNr int, currentState K) {
	if !traceEnabled {
		return
	}
	traceWriter.WriteString(fmt.Sprintf("\nstep #%d end, no more steps\n", stepNr))
}

func traceRuleApply(ruleType string, stepNr int, ruleInfo string) {
	if !traceEnabled {
		return
	}
	traceWriter.WriteString(fmt.Sprintf("rule %s #%-3d %s\n", ruleType, stepNr, ruleInfo))
}
