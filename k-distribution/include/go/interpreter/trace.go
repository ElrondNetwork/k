%COMMENT%

// Handles generation of traces
// (what rules were applied, in what order, what were the intermediate states).
// Multiple trace handlers supported.

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

var traceHandlers []traceHandler

// we can have multiple writers to write program execution traces in various formats
// they are all intended for easier debugging
type traceHandler interface {
	initialize()
	closeTrace()
	traceInitialState(state m.K)
	traceStepStart(stepNr int, currentState m.K)
	traceStepEnd(stepNr int, currentState m.K)
	traceNoStep(stepNr int, currentState m.K)
	traceRuleApply(ruleType string, stepNr int, ruleInfo string)
}

func (i *Interpreter) initializeTrace() {
	for _, t := range i.traceHandlers {
		t.initialize()
	}
}

func (i *Interpreter) closeTrace() {
	for _, t := range i.traceHandlers {
		t.closeTrace()
	}
}

func (i *Interpreter) traceInitialState(state m.K) {
	for _, t := range i.traceHandlers {
		t.traceInitialState(state)
	}
}

func (i *Interpreter) traceStepStart(stepNr int, currentState m.K) {
	for _, t := range i.traceHandlers {
		t.traceStepStart(stepNr, currentState)
	}
}

func (i *Interpreter) traceStepEnd(stepNr int, currentState m.K) {
	for _, t := range i.traceHandlers {
		t.traceStepEnd(stepNr, currentState)
	}
}

func (i *Interpreter) traceNoStep(stepNr int, currentState m.K) {
	for _, t := range i.traceHandlers {
		t.traceNoStep(stepNr, currentState)
	}
}

func (i *Interpreter) traceRuleApply(ruleType string, stepNr int, ruleInfo string) {
	for _, t := range i.traceHandlers {
		t.traceRuleApply(ruleType, stepNr, ruleInfo)
	}
}
