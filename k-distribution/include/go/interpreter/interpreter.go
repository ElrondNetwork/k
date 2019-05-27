%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

// Interpreter is a container with a reference to model and basic options
type Interpreter struct {
	Model *m.ModelState

	TracePretty bool
	TraceKPrint bool
	Verbose     bool
	MaxSteps    int
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter() *Interpreter {
	model := &m.ModelState{}
	model.Init()

	return &Interpreter{
		Model:       model,
		TracePretty: false,
		TraceKPrint: false,
		Verbose:     false,
		MaxSteps:    0,
	}
}
