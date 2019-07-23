%COMMENT%

package %PACKAGE%

import (
	"fmt"
	"strings"
)

// KObject defines a K item object that is managed by the model
type KObject interface {
	referenceType() kreferenceType
	equals(ms *ModelState, other KObject) bool
	deepCopy(ms *ModelState) KObject
	transfer(from, to *ModelData) KObject
	prettyPrint(ms *ModelState, sb *strings.Builder, indent int)
	kprint(ms *ModelState, sb *strings.Builder)
	collectionsToK(ms *ModelState) KReference
	increaseUsage(ms *ModelState)
	decreaseUsage(ms *ModelState)
	recycleUnused(ms *ModelState)
	preserve(ms *ModelState)
}

type objectReuseStatus int

const (
	active objectReuseStatus = iota
	inRecycleBin
	preserved
)

// ModelData holds part of state of the executor at a certain moment.
type ModelData struct {
	// reference to self
	selfRef modelDataReference

	// allKs keeps all K sequence elements into one large structure
	// all K sequence element references point to this structure
	allKsElements []ksequenceElem

	// contains all KApply args, concatenated into one slice
	// KApply references contain the start position and arity,
	// so enough data to find their args in this slice
	allKApplyArgs []KReference

	// keeps big int objects, big int references point here
	bigInts []*bigInt

	// recycle bin for big ints
	// works as a stack
	bigIntRecycleBin []KReference

	// contains all bytes from types String, Bytes and KToken
	allBytes []byte

	// keeps object types mixed together
	allObjects []KObject
}

// ModelState holds the full state of the executor at a certain moment.
type ModelState struct {
	// mainData keeps the main data of the running interpreter
	mainData *ModelData

	// memoTables is a structure containing all memoization maps.
	// Memoization tables are implemented as maps of maps of maps of ... of references
	memoTables map[MemoTable]interface{}

	// memoData is where the actual data for the memoization tables resides
	// unlike the mainData, we never flush this one before execution is over
	memoData *ModelData

	// swapData is a small model used when collecting garbage.
	// It is needed to keep the current state while the main model is flushed.
	swapData *ModelData
}

// constantsData is the model data object that keeps constants defined statically
var constantsData = newSmallModelData(constDataRef)

func (ms *ModelState) getData(dataRef modelDataReference) *ModelData {
	switch dataRef {
	case mainDataRef:
		return ms.mainData
	case memoDataRef:
		return ms.memoData
	case constDataRef:
		return constantsData
	default:
		panic("unknown modelDataReference")
	}
}

func (md *ModelData) getReferencedObject(index uint64) KObject {
	if index >= uint64(len(md.allObjects)) {
		panic("trying to reference an object beyond allocated slice")
	}
	return md.allObjects[index]
}

func (md *ModelData) addObject(obj KObject) KReference {
	newIndex := len(md.allObjects)
	md.allObjects = append(md.allObjects, obj)
	return createKrefBasic(obj.referenceType(), md.selfRef, uint64(newIndex))
}

// NewModel creates a new blank model.
func NewModel() *ModelState {
	ms := &ModelState{}
	ms.mainData = newMainModelData()
	ms.memoTables = nil
	ms.memoData = newSmallModelData(memoDataRef)
	return ms
}

// newSmallModel creates a smaller model.
func newMainModelData() *ModelData {
	return &ModelData{
		selfRef:       mainDataRef,
		allKsElements: make([]ksequenceElem, 0, 100000),
		allKApplyArgs: make([]KReference, 0, 1000000),
		allBytes:      make([]byte, 0, 1<<20),
		allObjects:    make([]KObject, 0, 10000),
	}
}

// newSmallModel creates a smaller model.
func newSmallModelData(selfRef modelDataReference) *ModelData {
	return &ModelData{
		selfRef:       selfRef,
		allKsElements: make([]ksequenceElem, 0, 32),
		allKApplyArgs: make([]KReference, 0, 32),
		allBytes:      make([]byte, 0, 1024),
		allObjects:    make([]KObject, 0, 32),
	}
}

// Clear resets the model data as if it were new,
// but does not free the memory allocated by previous execution.
func (md *ModelData) Clear() {
	md.allKsElements = md.allKsElements[:0]
	md.allKApplyArgs = md.allKApplyArgs[:0]
	md.allObjects = md.allObjects[:0]
	md.allBytes = md.allBytes[:0]
	md.recycleAllInts()
}

// Clear resets the model as if it were new,
// but does not free the memory allocated by previous execution.
func (ms *ModelState) Clear() {
	ms.mainData.Clear()
	ms.memoTables = nil
	ms.memoData.Clear()
}

// Gc cleans up the model, but keeps the last state, given as argument.
// It does so by temporarily copying the last state to another model.
func (ms *ModelState) Gc(keepState KReference) KReference {
	if ms.swapData == nil {
		ms.swapData = newSmallModelData(noDataRef)
	} else {
		ms.swapData.Clear()
	}
	copiedState := transfer(ms.mainData, ms.swapData, keepState)
	ms.Clear()
	newState := transfer(ms.swapData, ms.mainData, copiedState)
	return newState
}

// PrintStats simply prints some statistics to the console.
// Useful for checking the size of the model data.
func (md *ModelData) PrintStats() {
	fmt.Printf("\n   KApply args:                           %d (cap: %d)", len(md.allKApplyArgs), cap(md.allKApplyArgs))
	fmt.Printf("\n   K sequence elements:                   %d (cap: %d)", len(md.allKsElements), cap(md.allKsElements))
	fmt.Printf("\n   BigInt objects:                        %d (cap: %d)", len(md.bigInts), cap(md.bigInts))
	fmt.Printf("\n   Bytes (strings, byte arrays, KTokens): %d (cap: %d)", len(md.allBytes), cap(md.allBytes))
	fmt.Printf("\n   Other objects:                         %d (cap: %d)", len(md.allObjects), cap(md.allObjects))
	fmt.Printf("\n   Recycle bin - BigInt                   %d (cap: %d)", len(md.bigIntRecycleBin), cap(md.bigIntRecycleBin))
}

// PrintStats simply prints some statistics to the console.
// Useful for checking the size of the model.
func (ms *ModelState) PrintStats() {
	fmt.Print("\nMain data container:")
	ms.mainData.PrintStats()
	fmt.Print("\nMemoization data container:")
	ms.memoData.PrintStats()
	if ms.swapData != nil {
		fmt.Print("\nSwap data container:")
		ms.swapData.PrintStats()
	}
	fmt.Print("\nConstants data container:")
	constantsData.PrintStats()
	fmt.Println()
}
