package %PACKAGE_INTERPRETER%

import (
	"bufio"
	"bytes"
	"fmt"
	m "%INCLUDE_MODEL%"
	koreparser "%INCLUDE_PARSER%"
	"os"
	"time"
)

// StartTraceReferenceComparer ... set the code to compare each step to list of steps exported from somewhere else
func StartTraceReferenceComparer(referenceFileName string) {
	trc := &traceReferenceComparer{referenceFileName: referenceFileName}
	traceHandlers = append(traceHandlers, trc)
}

type traceReferenceComparer struct {
	referenceFileName string
	referenceFile     *os.File
	fileReader        *bufio.Reader
	currentStep       int
	running           bool
}

func (t *traceReferenceComparer) initialize() {
	var err error
	t.referenceFile, err = os.Open(t.referenceFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.fileReader = bufio.NewReader(t.referenceFile)
	t.running = true

	//t.consumeStepLine()
	//t.consumeStepLine()
	//t.consumeStepLine()
}

func (t *traceReferenceComparer) closeTrace() {
	t.referenceFile.Close()
}

func (t *traceReferenceComparer) traceInitialState(state m.K) {
}

func (t *traceReferenceComparer) consumeStepLine() string {
	var err error
	expectedPrefix := []byte(fmt.Sprintf("step %d ", t.currentStep))
	t.currentStep++
	prefix := make([]byte, len(expectedPrefix))
	_, err = t.fileReader.Read(prefix)
	if err != nil {
		fmt.Println(err)
		t.running = false
		return ""
	}
	if !bytes.Equal(expectedPrefix, prefix) {
		fmt.Println("Wrong prefix in trace to compare with.")
		t.running = false
		return ""
	}

	var line string
	line, err = t.fileReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		t.running = false
		return ""
	}
	return line
}

func (t *traceReferenceComparer) traceStepStart(stepNr int, currentState m.K) {
	if !t.running {
		return
	}
	line := t.consumeStepLine()

	parsedLine := koreparser.Parse([]byte(line))
	lineAsK := convertParserModelToKModel(parsedLine)

	pure := m.CollectionsToK(currentState)
	if !pure.Equals(lineAsK) {
		fmt.Printf("Stopped at step %d.", stepNr)
		t.running = false
		formattedNow := time.Now().Format("20060102150405")
		writeStateToFile(pure, "traceRef_"+formattedNow+"_actual.log")
		writeStateToFile(lineAsK, "traceRef_"+formattedNow+"_expected.log")
	}
}

func writeStateToFile(state m.K, fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error while creating trace file:" + err.Error())
	}
	writer := bufio.NewWriter(file)
	writer.WriteString(m.PrettyPrint(state))
	writer.Flush()
	file.Close()
}

func (t *traceReferenceComparer) traceStepEnd(stepNr int, currentState m.K) {
}

func (t *traceReferenceComparer) traceNoStep(stepNr int, currentState m.K) {
}

func (t *traceReferenceComparer) traceRuleApply(ruleType string, stepNr int, ruleInfo string) {
}
