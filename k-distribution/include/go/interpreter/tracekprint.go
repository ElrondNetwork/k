%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	"bufio"
	"fmt"
	m "%INCLUDE_MODEL%"
	"os"
	"time"
)

// this one writes the state at each step, all in one file
type traceKPrint struct {
	file       *os.File
	fileWriter *bufio.Writer
	interpreter *Interpreter
}

func (t *traceKPrint) initialize() {
	formattedNow := time.Now().Format("20060102150405")
	fileName := "traceKPrint_" + formattedNow + ".log"
	var err error
	t.file, err = os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error while creating trace file:" + err.Error())
	}

	t.fileWriter = bufio.NewWriter(t.file)
}

func (t *traceKPrint) closeTrace() {
	t.fileWriter.Flush()
	t.file.Close()
}

func (t *traceKPrint) traceInitialState(state m.K) {
}

func (t *traceKPrint) traceStepStart(stepNr int, currentState m.K) {
	kast := t.interpreter.Model.KPrint(currentState)
	t.fileWriter.WriteString(fmt.Sprintf("\nstep %d %s", stepNr, kast))
}

func (t *traceKPrint) traceStepEnd(stepNr int, currentState m.K) {
}

func (t *traceKPrint) traceNoStep(stepNr int, currentState m.K) {
}

func (t *traceKPrint) traceRuleApply(ruleType string, stepNr int, ruleInfo string) {
}
