package main

import (
	"fmt"
)

type noStepError struct {
}

func (e *noStepError) Error() string {
	return "No step could be performed."
}

type stuckError struct {
	funcName string
	args     []K
}

func (e *stuckError) Error() string {
	if len(e.args) == 0 {
		return "Stuck! Function name: " + e.funcName + ". No args."
	}
	s := "Stuck! Function name: " + e.funcName + ". Args:"
	for i, arg := range e.args {
		s += fmt.Sprintf("\n\t%d: %s", i, arg.PrettyTreePrint(1))
	}
	return s
}

type evalArityViolatedError struct {
	funcName      string
	expectedArity int
	actualArity   int
}

func (e *evalArityViolatedError) Error() string {
	return fmt.Sprintf(
		"Eval function arity violated. Function name: %s. Expected arity: %d. Actual arity: %d.",
		e.funcName, e.expectedArity, e.actualArity)
}

type hookNotImplementedError struct {
}

func (e *hookNotImplementedError) Error() string {
	return "Hook not implemented."
}

type hookInvalidArgsError struct {
}

func (e *hookInvalidArgsError) Error() string {
	return "Invalid argument(s) provided to hook."
}
