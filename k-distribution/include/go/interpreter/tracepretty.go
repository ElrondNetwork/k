%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	"bufio"
	"fmt"
	m "%INCLUDE_MODEL%"
	"os"
	"path/filepath"
	"time"
)

// creates a folder with the timestamp and writes each step in a separate file
// this way it is easier to follow changes from one step to the next
type tracePrettyDebug struct {
	dirName     string
	currentFile *os.File
	fileWriter  *bufio.Writer
	interpreter *Interpreter
}

func (t *tracePrettyDebug) initialize() {
	formattedNow := time.Now().Format("20060102150405")
	t.dirName = "trace_" + formattedNow
	var err error
	err = os.Mkdir(t.dirName, os.ModePerm)
	if err != nil {
		fmt.Println("Error while creating trace directory:" + err.Error())
	}
	t.newTraceFile(t.dirName + "_init.log")
}

func (t *tracePrettyDebug) newTraceFile(fileName string) {
	if t.currentFile != nil {
		t.fileWriter.Flush()
		t.currentFile.Close()
	}
	var err error
	t.currentFile, err = os.OpenFile(filepath.Join(t.dirName, fileName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error while creating trace file:" + err.Error())
	}

	t.fileWriter = bufio.NewWriter(t.currentFile)
}

func (t *tracePrettyDebug) closeTrace() {
	t.fileWriter.Flush()
	t.currentFile.Close()
}

func (t *tracePrettyDebug) traceInitialState(state m.K) {
	t.fileWriter.WriteString("initial state:\n\n")
	t.fileWriter.WriteString(t.interpreter.Model.PrettyPrint(state))
}

func (t *tracePrettyDebug) traceStepStart(stepNr int, currentState m.K) {
	t.newTraceFile(fmt.Sprintf("%s_step%d.log", t.dirName, stepNr))
	t.fileWriter.WriteString(fmt.Sprintf("\n\nstep #%d begin\n\n", stepNr))
}

func (t *tracePrettyDebug) traceStepEnd(stepNr int, currentState m.K) {
	t.fileWriter.WriteString(fmt.Sprintf("\nstep #%d end; current state:\n\n", stepNr))
	t.fileWriter.WriteString(t.interpreter.Model.PrettyPrint(currentState))
}

func (t *tracePrettyDebug) traceNoStep(stepNr int, currentState m.K) {
	t.fileWriter.WriteString(fmt.Sprintf("\nstep #%d end, no more steps\n\n", stepNr))
	t.fileWriter.WriteString(t.interpreter.Model.PrettyPrint(currentState))
}

func (t *tracePrettyDebug) traceRuleApply(ruleType string, stepNr int, ruleInfo string) {
	if ruleType == "STEP" {
		t.fileWriter.WriteString(fmt.Sprintf("rule %s #%-3d %s\n", ruleType, stepNr, ruleInfo))
	}
}
