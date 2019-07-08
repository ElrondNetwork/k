%COMMENT%

package %PACKAGE%

import (
	m "%INCLUDE_MODEL%"
)

// GetNrSteps yields how many steps were executed until now from the start of the last execution
func (i *Interpreter) GetNrSteps() int {
    return i.currentStep
}

// GetLastState yields the current (last) state of the interpreter
func (i *Interpreter) GetState() m.KReference {
     return i.state
}
