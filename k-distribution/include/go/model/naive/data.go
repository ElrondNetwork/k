%COMMENT%

package %PACKAGE%

import (
    "strings"
)

// K defines a K entity
type K interface {
	equals(other K) bool
	deepCopy() K
	prettyPrint(ms *ModelState, sb *strings.Builder, indent int)
	kprint(ms *ModelState, sb *strings.Builder)
	collectionsToK(ms *ModelState) K
}

// ModelState holds the state of the executor at a certain moment
type ModelState struct {
    initialized bool

	// memoTables is a structure containing all memoization maps.
	// Memoization tables are implemented as maps of maps of maps of ...
	memoTables map[MemoTable]interface{}
}

// Init prepares model for execution
func NewModel() *ModelState {
    ms := &ModelState{}
    ms.Init()
    return ms
}

// Init prepares model for execution
func (ms *ModelState) Init() {
    if ms.initialized {
        return
    }
}

// ClearModel ... clean up any data left from previous executions, to save memory
func (ms *ModelState) ClearModel() {
    ms.initialized = false
    ms.Init()
}

// PrintStats simply prints some statistics to the console.
func (ms *ModelState) PrintStats() {
    // nothing in this version
}