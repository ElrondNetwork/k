// Handles generation of traces
// (what rules were applied, in what order, what were the intermediate states).
// Generates a file and dumps trace there.

package %PACKAGE_INTERPRETER%

import (
	"bufio"
	"fmt"
	m "%INCLUDE_MODEL%"
	"os"
	filepath "path/filepath"
	"time"
)

var traceEnabled = false
var traceName string
var traceFile *os.File
var traceWriter *bufio.Writer

func initializeTrace() {
	traceEnabled = true
	formattedNow := time.Now().Format("20060102150405")
	traceName = "trace_" + formattedNow
	var err error
	err = os.Mkdir(traceName, os.ModePerm)
	if err != nil {
		fmt.Println("Error while creating trace directory:" + err.Error())
	}
	newTraceFile(traceName + "_init.log")
}

func newTraceFile(fileName string) {
	if traceFile != nil {
		traceWriter.Flush()
		traceFile.Close()
	}
	var err error
	traceFile, err = os.OpenFile(filepath.Join(traceName, fileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error while creating trace file:" + err.Error())
	}

	traceWriter = bufio.NewWriter(traceFile)
}

func closeTrace() {
	if !traceEnabled {
		return
	}
	traceEnabled = false
	traceWriter.Flush()
	traceFile.Close()
}

func traceInitialState(state m.K) {
	if !traceEnabled {
		return
	}
	traceWriter.WriteString("initial state:\n\n")
	traceWriter.WriteString(m.PrettyPrint(state))
}

func traceStepStart(stepNr int, currentState m.K) {
	if !traceEnabled {
		return
	}
	newTraceFile(fmt.Sprintf("%s_step%d.log", traceName, stepNr))
	traceWriter.WriteString(fmt.Sprintf("\n\nstep #%d begin\n\n", stepNr))
}

func traceStepEnd(stepNr int, currentState m.K) {
	if !traceEnabled {
		return
	}
	traceWriter.WriteString(fmt.Sprintf("\nstep #%d end; current state:\n\n", stepNr))
	traceWriter.WriteString(m.PrettyPrint(currentState))
}

func traceNoStep(stepNr int, currentState m.K) {
	if !traceEnabled {
		return
	}
	traceWriter.WriteString(fmt.Sprintf("\nstep #%d end, no more steps\n\n", stepNr))
	traceWriter.WriteString(m.PrettyPrint(currentState))
}

func traceRuleApply(ruleType string, stepNr int, ruleInfo string) {
	if !traceEnabled {
		return
	}
	if ruleType == "STEP" {
		traceWriter.WriteString(fmt.Sprintf("rule %s #%-3d %s\n", ruleType, stepNr, ruleInfo))
	}
}
