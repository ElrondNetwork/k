%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"testing"
)

var inputBackup []m.K

// saves a copy of the arguments, so we can later check if they changed during the call
func (interpreter *Interpreter) backupInput(args ...m.K) {
	inputBackup = make([]m.K, len(args))
	for i := 0; i < len(args); i++ {
		inputBackup[i] = interpreter.Model.DeepCopy(args[i])
	}
}

// checks that arguments didn't change in the hook
func (interpreter *Interpreter) checkImmutable(t *testing.T, args ...m.K) {
	if len(args) != len(inputBackup) {
		t.Error("Test not set up properly. Should be the same number of parameters as the last backupInput call.")
	}
	for i := 0; i < len(args); i++ {
		if !interpreter.Model.Equals(args[i], inputBackup[i]) {
			t.Errorf("Input state changed! Got:%s Want:%s", interpreter.Model.PrettyPrint(args[i]), interpreter.Model.PrettyPrint(inputBackup[i]))

		}
	}
}
