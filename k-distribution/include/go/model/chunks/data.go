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
	prettyPrint(ms *ModelState, sb *strings.Builder, indent int)
	kprint(ms *ModelState, sb *strings.Builder)
	collectionsToK(ms *ModelState) KReference
}

// ModelState holds the state of the executor at a certain moment
type ModelState struct {
	initialized bool

	// allKs keeps all KSequences into one large structure
	// all KSequences point to this structure
	allKs *ksequenceSliceContainer

	// keeps object types mixed together
	// objects are grouped in chunks
	allObjects [][]KObject

	// memoTables is a structure containing all memoization maps.
	// Memoization tables are implemented as maps of maps of maps of ...
	memoTables map[MemoTable]interface{}
}

// constantsModel is another instance of the model, but which only contains a few constants.
var constantsModel = NewModel()

const itemIndexBits = 8
const chunkCapacity = 1 << itemIndexBits
const itemIndexMask = chunkCapacity - 1

func (ms *ModelState) getObjectAtIndex(index uint32) KObject {
	chunkIndex := index >> itemIndexBits
	itemIndex := index & itemIndexMask
	if chunkIndex >= uint32(len(ms.allObjects)) {
		panic("trying to reference object beyond allocated objects")
	}
	chunk := ms.allObjects[chunkIndex]
	if itemIndex >= uint32(len(chunk)) {
		panic("trying to reference object beyond allocated objects")
	}
	return chunk[itemIndex]
}

func (ms *ModelState) getReferencedObject(ref KReference) KObject {
	if ref.constantObject {
		return constantsModel.getObjectAtIndex(uint32(ref.value1))
	}
	return ms.getObjectAtIndex(uint32(ref.value1))
}

func (ms *ModelState) addObject(obj KObject) KReference {
	var lastChunk []KObject
	var lastChunkIndex int
	if len(ms.allObjects) == 0 {
		lastChunkIndex = 0
		lastChunk = make([]KObject, 0, chunkCapacity)
		ms.allObjects = append(ms.allObjects, nil)
	} else {
		lastChunkIndex = len(ms.allObjects) - 1
		lastChunk = ms.allObjects[lastChunkIndex]
		if len(lastChunk) == chunkCapacity {
			lastChunkIndex++
			lastChunk = make([]KObject, 0, chunkCapacity)
			ms.allObjects = append(ms.allObjects, nil)
		}
	}

	newItemIndex := len(lastChunk)
	lastChunk = append(lastChunk, obj)
	ms.allObjects[lastChunkIndex] = lastChunk

	newIndex := (lastChunkIndex << itemIndexBits) | newItemIndex
	return KReference{refType: obj.referenceType(), constantObject: false, value1: newIndex, value2: 0}
}

func addConstantObject(obj KObject) KReference {
	ref := constantsModel.addObject(obj)
	ref.constantObject = true
	return ref
}

func (ms *ModelState) countObjects() int {
	if len(ms.allObjects) == 0 {
		return 0
	}
	nrChunks := len(ms.allObjects)
	lastChunk := ms.allObjects[nrChunks-1]
	return ((nrChunks - 1) << itemIndexBits) | len(lastChunk)
}

// NewModel creates a new blank model.
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
	ms.initialized = true
	ms.allKs = &ksequenceSliceContainer{}
	ms.allObjects = nil
	ms.memoTables = nil
}

// ClearModel ... clean up any data left from previous executions, to save memory
func (ms *ModelState) ClearModel() {
	ms.initialized = false
	ms.Init()
}

// PrintStats simply prints some statistics to the console.
// Useful for checking the size of the model data.
func (ms *ModelState) PrintStats() {
	fmt.Printf("Nr. objects: %d (%d chunks)\n", ms.countObjects(), len(ms.allObjects))
	fmt.Printf("Nr. K sequence slices: %d\n", len(ms.allKs.allSlices))
}
